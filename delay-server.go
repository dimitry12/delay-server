package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// limitNumClients is HTTP handling middleware that ensures no more than
// maxClients requests are passed concurrently to the given handler f.
func limitNumClients(f http.HandlerFunc, maxClients int) http.HandlerFunc {
  // Counting semaphore using a buffered channel
  sema := make(chan struct{}, maxClients)

  return func(w http.ResponseWriter, req *http.Request) {
    sema <- struct{}{}
    defer func() { <-sema }()
    f(w, req)
  }
}

func handler(w http.ResponseWriter, r *http.Request) {

        setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

        maxMs, parseError := strconv.ParseInt(r.URL.Query().Get("max"), 0, 32)
	if parseError != nil {
		maxMs = 1
	}

	minMs, parseError := strconv.ParseInt(r.URL.Query().Get("min"), 0, 32)
	if parseError != nil {
		minMs = 0
	}

	failureChance, parseError := strconv.ParseInt(r.URL.Query().Get("failure"), 0, 32)
	if parseError != nil {
		failureChance = 0
	}

	if maxMs < 0 || maxMs > 30000 {
		http.Error(w, "invalid 'maxMs' query param. must be >= 0 and <= 30000", http.StatusBadRequest)
	}

	if minMs < 0 {
		http.Error(w, "invalid 'minMs' query param. must be >= 0", http.StatusBadRequest)
	}

	if maxMs < minMs {
		http.Error(w, "invalid 'maxMs' & 'minMs' query params. maxMs must be greater than or equal to minMs", http.StatusBadRequest)
	}

	time.Sleep(time.Duration(rand.Intn(int(maxMs-minMs))+int(minMs)) * time.Millisecond)

	if failureChance > 0 && rand.Intn(int(failureChance)) == 0 {
		http.Error(w, "Mock error", http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "<h1>welcome to the go delay server</h1>"+
			"<h2> supported query params</h2>"+
			"<ul>"+
			"<li>max : max delay in milliseconds. defaults to 1</li>"+
			"<li>min : min delay in milliseconds. defaults to 0</li>"+
			"<li>failure : 1 in X failure chance. defaults to 0 (off)</li>"+
			"</ul>")
	}

}

func main() {
	http.HandleFunc("/", limitNumClients(handler, 1))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
	fmt.Println("listing on " + port)
}
