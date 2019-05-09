package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

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
	res := client.Ping()
	if res.Err() != nil {
		fmt.Printf("Failed to connect to redis : %v \n", res.Err())
		os.Exit(1)
	}
}

type response struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Pong")
}

func respondError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	respond(w, message, "")
}
func respond(w io.Writer, msg, data string) {
	response := response{
		Message: msg,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
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
		respondError(w, "Not a valid URL")
		return
	}
	urlHash, err := Hash(payload.URL)

	if err != nil {
		respondError(w, "Failure")
		return
	}
	res := client.HSet("urlmaps", urlHash, payload.URL)
	if res.Err() != nil {
		respondError(w, "Failed to save value in redis")
		return
	}
	respond(w, "Success", urlHash)

}

func lookup(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	url, err := client.HGet("urlmaps", hash).Result()
	if err != nil {
		fmt.Printf("Error : %v \n", err)
		respondError(w, "Failure")
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
