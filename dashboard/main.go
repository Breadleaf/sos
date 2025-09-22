package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	sos "github.com/breadleaf/sos/pkg/http/client"
)

var sosClient *sos.Client
var serverBind string

func init() {
	var bindAddr string
	flag.StringVar(
		&bindAddr,
		"listen",
		":8000",
		"address to bind HTTP server",
	)

	var serverURL string
	flag.StringVar(
		&serverURL,
		"server url",
		":8080",
		"url of the sos server to communicate to",
	)

	flag.Parse()

	sosClient = sos.NewClient(serverURL)

	serverBind = fmt.Sprintf("localhost%s", bindAddr)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reroute := "/static/dashboard.html"
		log.Printf("routing user to: %s\n", reroute)
		http.Redirect(w, r, reroute, http.StatusSeeOther)
	})

	mux.HandleFunc("/static/{file}", func(w http.ResponseWriter, r *http.Request) {
		file := r.PathValue("file")
		filePath := fmt.Sprintf("./static/%s", file)

		log.Printf("user has requested file: %s, stating...\n", filePath)

		// only send file if it exists and is valid
		if stat, err := os.Stat(filePath); err != nil && stat.IsDir() {
			log.Printf(
				"requested file: %s, either did not exist or was a directory..., error %v\n",
				filePath,
				err,
			)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("serving file: %s\n", filePath)
			http.ServeFile(w, r, filePath)
		}
	})

	log.Fatal(
		http.ListenAndServe(serverBind, mux),
	)
}
