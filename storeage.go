package main

type Storer interface {
	Setup(cfgpath string) error
	Get(ident string) ([]byte, error)
	Save(ident string, data []byte) error
	Delete(ident string) error
}
