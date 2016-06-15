package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	setupCLI()

	flag.Parse()

	http.HandleFunc("/", terraState)
	log.Fatal(http.ListenAndServe(":12345", nil))
}
