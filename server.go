package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var workerPool = NewWorkerPool(5)

func main() {

	router := mux.NewRouter()

	go workerPool.StartListen()

	log.Println(
		`
		[main] Worker Pool started to listen. 
		Please, send an HTTP request to localhost:8000/number
		`,
	)

	router.HandleFunc("/{n}", newJob)

	errStartingServer := http.ListenAndServe(":8080", router)
	if errStartingServer != nil {
		log.Fatalln("Failed to start up the HTTP service")
	}

}

func newJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := vars["n"]

	if !err {
		http.Error(w, "n not given", http.StatusBadRequest)
		return
	}

	parsedNumber, _ := strconv.Atoi(n)
	workerPool.jobQueue <- parsedNumber
}
