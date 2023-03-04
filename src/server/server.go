package main

import (
	"log"
	"net/http"

	"github.com/colorodoxyz/serve-rest/src/helper"
)

var store = make(map[string]helper.KeyValue)

func keyValueApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
	} else if r.Method == "POST" {
	} else if r.Method == "DELETE" {
	} else {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

}

func main() {
	http.HandleFunc("/api/key", keyValueApi)
	err := http.ListenAndServe(":5001", nil)
	log.Fatal(err)
}
