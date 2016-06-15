package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var once sync.Once
var storage Storer

const (
	ok            = 200
	internalError = 500
)

func getStorage() Storer {
	once.Do(func() {
		t := &MongoDB{}
		err := t.Setup(getConfigPath())
		if err != nil {
			panic(err)
		}

		storage = t
	})

	return storage
}

// Terra State, the web server
func terraState(w http.ResponseWriter, r *http.Request) {
	engine := getStorage()

	ident := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/"), "/")

	switch r.Method {
	case "GET":
		state, err := engine.Get(ident)
		if err != nil {
			log.Println("Error Getting State", err.Error())
			w.WriteHeader(internalError)
		}
		w.Write(state)
	case "POST":
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, r.Body); err != nil {
			w.WriteHeader(internalError)
		} else if err := engine.Save(ident, buf.Bytes()); err != nil {
			log.Println("Error Saving State", err.Error())
		}
	case "DELETE":
		if err := engine.Delete(ident); err != nil {
			w.WriteHeader(ok)
		} else {
			log.Println("Error Deleting State", err.Error())
			w.WriteHeader(internalError)
		}
	default:
		w.WriteHeader(internalError)
		w.Write([]byte(fmt.Sprintf("Unknown method: %s", r.Method)))
	}
}
