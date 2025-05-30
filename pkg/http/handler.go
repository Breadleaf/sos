package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/breadleaf/sos/pkg/storage"
	"github.com/gorilla/mux"
)

func NewHandler(store storage.Backend) http.Handler {
	r := mux.NewRouter()

	// create or delete bucket
	r.HandleFunc("/buckets/{bucket}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["bucket"]

		switch r.Method {
		case http.MethodPut: // create
			if err := store.CreateBucket(name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case http.MethodDelete: // delete
			if err := store.DeleteBucket(name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default: // invalid
			http.Error(w, fmt.Sprintf("method not allowed: %s", r.Method), http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodPut, http.MethodDelete)

	// list all buckets
	r.HandleFunc("/buckets", func(w http.ResponseWriter, r *http.Request) {
		buckets, err := store.ListBuckets()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(buckets)
	}).Methods(http.MethodGet)

	// upload, download, or delete a single object
	r.HandleFunc("/buckets/{bucket}/object/{key:.*}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucket, key := vars["bucket"], vars["key"]

		switch r.Method {
		case http.MethodPut: // upload
			defer r.Body.Close()
			if err := store.PutObject(bucket, key, r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		case http.MethodGet: // download
			reader, err := store.GetObject(bucket, key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			defer reader.Close()
			io.Copy(w, reader)
		case http.MethodDelete: // delete
			if err := store.DeleteObject(bucket, key); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default: // invalid
			http.Error(w, fmt.Sprintf("method not allowed: %s", r.Method), http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodPut, http.MethodGet, http.MethodDelete)

	// list all objects in a bucket
	r.HandleFunc("/buckets/{bucket}/objects", func(w http.ResponseWriter, r *http.Request) {
		bucket := mux.Vars(r)["bucket"]
		objects, err := store.ListObjects(bucket)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(objects)
	}).Methods(http.MethodGet)

	return r
}
