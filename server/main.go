package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	sosHTTP "github.com/breadleaf/sos/pkg/http"
	"github.com/breadleaf/sos/pkg/storage"
)

func main() {
	var rootDir string
	var bindAddr string

	flag.StringVar(&rootDir, "data", "./data", "filesystem root for buckets")
	flag.StringVar(&bindAddr, "listen", ":8080", "address to bind HTTP server")
	flag.Parse()

	store, err := storage.NewDiskBackend(rootDir)
	if err != nil {
		log.Fatalf("failed to init DiskBackend: %v", err)
	}

	handler := sosHTTP.NewHandler(store)
	fmt.Printf("SOS server listening on %s, data in %q\n", bindAddr, rootDir)
	log.Fatal(http.ListenAndServe(bindAddr, handler))
}
