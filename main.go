package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imgix/imgix-go"
)

// Response struct
type Response struct {
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

var imgIXClient imgix.Client

func main() {

	imgIXClient = imgix.NewClient("serviceprofile.imgix.net")

	r := mux.NewRouter()
	r.HandleFunc("/api/profiles", uploadHandler).Methods("POST")
	r.HandleFunc("/api/profiles/{url}", imgixHandler).Methods("GET")

	fmt.Printf("Running app in: %v\n", PORT)
	http.ListenAndServe(":"+PORT, r)
}
