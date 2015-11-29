package main

import (
	"flag"
	"net/http"
	"github.com/go_files_api/files"
)

func main() {
	var (
		mongoUri = flag.String("mongoUri", "mongodb://127.0.0.1:27017", "MongoDB Connection String")
		apiUri = flag.String("apiUri", "localhost:8080", "Web API Hosting URI")
	)
	flag.Parse()

	session := files.Session(*mongoUri)
	defer session.Close()

	http.ListenAndServe(*apiUri, files.AppHandlers(session))
}
