package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var urls = map[string]string{}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Pong")
}

func respond(w io.Writer, msg, data string) {
	response := struct {
		Message string `json:"message"`
		Data    string `json:"string"`
	}{
		Message: msg,
		Data:    data,
	}
	enc := json.NewEncoder(w)
	enc.Encode(response)
}

func Hash(a string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, a)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func shotern(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	urlHash, err := Hash(payload.URL)
	if err != nil {
		respond(w, "Failure", "")
	} else {
		respond(w, "Success", urlHash)
		urls[urlHash] = payload.URL
	}
}

func lookup(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	url, ok := urls[hash]
	if !ok {
		respond(w, "Failure", "")
	} else {
		respond(w, "Success", url)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", ping)
	r.HandleFunc("/shortern", shotern)
	r.HandleFunc("/{hash}", lookup)
	http.Handle("/", r)
	fmt.Println("Starting the server... ")
	http.ListenAndServe("localhost:8000", nil)
}
