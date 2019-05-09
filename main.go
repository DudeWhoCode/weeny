package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var client *redis.Client

func init() {
	client = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Pong")
}

func respond(w io.Writer, msg, data string) {
	response := struct {
		Message string `json:"message"`
		Data    string `json:"data"`
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
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	return true

}

func shotern(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	if valid := isValidURL(payload.URL); !valid {
		w.WriteHeader(http.StatusBadRequest)
		respond(w, "Not a valid URL", "")
		return
	}
	urlHash, err := Hash(payload.URL)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respond(w, "Failure", "")
		return
	}
	respond(w, "Success", urlHash)
	client.HSet("urlmaps", urlHash, payload.URL)

}

func lookup(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	url, err := client.HGet("urlmaps", hash).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respond(w, "Failure", "")
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	respond(w, "Success", url)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", ping).Methods("GET")
	r.HandleFunc("/shortern", shotern).Methods("POST")
	r.HandleFunc("/{hash}", lookup).Methods("GET")
	fmt.Println("Starting the server... ")
	if err := http.ListenAndServe("localhost:8000", r); err != nil {
		fmt.Printf("Failed to start server %v \n", err)
	}
}
