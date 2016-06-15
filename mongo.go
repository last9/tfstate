package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	conf *dbConf
}

type dbConf struct {
	Mongo struct {
		Host     string
		Database string
		Username string
		Password string
	}
}

type tfState struct {
	Version uint32
	Serial  uint32
	Modules []map[string]interface{}
}

// Dumb wrapper around conf.
func (s *MongoDB) getConfig() *dbConf {
	return s.conf
}

// Get a new Session to MongoDB. Do not need to cache it.
// mgo has an internal pool.
func (s *MongoDB) getSession() (*mgo.Session, error) {
	cfg := s.getConfig()

	session, err := mgo.Dial(cfg.Mongo.Host)
	if err != nil {
		return nil, err
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

// Try to Get the MongoDB and yeah save it also.
func (md *MongoDB) getDb(s *mgo.Session) (*mgo.Database, error) {
	cfg := md.getConfig()
	db := s.DB(cfg.Mongo.Database)

	if cfg.Mongo.Username == "" && cfg.Mongo.Password == "" {
		return db, nil
	}

	if err := db.Login(cfg.Mongo.Username, cfg.Mongo.Password); err != nil {
		return nil, err
	}

	return db, nil
}

var defaultYaml = `
mongo:
host: "127.0.0.1:27017"
database: terraform
`

// Do not do anything in MongoDB Setup.
func (s *MongoDB) Setup(cfgpath string) error {
	t := &dbConf{}

	if cfgpath == "" {
		content := []byte(defaultYaml)
		tmpfile, err := ioutil.TempFile("", "default_tfstate.yaml")
		if err != nil {
			return err
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(content); err != nil {
			return err
		}

		if err := tmpfile.Close(); err != nil {
			return err
		}

		cfgpath = tmpfile.Name()
	}

	if err := parseConfig(cfgpath, t); err != nil {
		return err
	}

	s.conf = t
	return nil
}

// Get the latest state for a ident saved in the database.
func (s *MongoDB) Get(ident string) ([]byte, error) {
	log.Println("Getting ident", ident)

	session, err := s.getSession()
	if err != nil {
		return nil, err
	}

	defer session.Close()

	db, err := s.getDb(session)
	if err != nil {
		return nil, err
	}

	coll := db.C(ident)

	state := &tfState{}
	if err := coll.Find(nil).Sort("$natural").One(&state); err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
	}

	return []byte{}, nil
}

// Write the state to a database, even if the same version/serial exists, new
// entry will be created. This will help maintain versions and never loose
// history.
func (s *MongoDB) Save(ident string, data []byte) error {
	log.Println("Saving ident", ident)

	m := &tfState{}
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}

	session, err := s.getSession()
	if err != nil {
		return err
	}

	defer session.Close()

	db, err := s.getDb(session)
	if err != nil {
		return err
	}

	coll := db.C(ident)

	return coll.Insert(m)
}

// Delete the state from database? How do we care. Just archive the collection
// by renaming it to <coll>-archive-<timestamp>.
func (s *MongoDB) Delete(ident string) error {
	log.Println("Deleting ident", ident)
	return nil
}
