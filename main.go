package main

import (
	"flag"
	"net/http"
	"log"
)

var (
	listen string
	file string
)

func init() {
	flag.StringVar(&listen, "listen", ":8080", "Host and port to listen to")
	flag.StringVar(&file, "file", "emails", "Filename to save emails to")
	flag.Parse()
}

func main() {
	log.Print("Listening on ", listen)
	http.ListenAndServe(listen, mux())
}