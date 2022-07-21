package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	workerPool := NewWorkerPool(5)

	go workerPool.StartListen()

	log.Println("[main] Worker Pool started to listen. Please, send an HTTP request to localhost:8000/number")

	router.HandleFunc("/{n}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		n, err := vars["n"]

		if !err {
			http.Error(w, "n not given", http.StatusBadRequest)
			return
		}

		parsedNumber, _ := strconv.Atoi(n)
		workerPool.jobQueue <- parsedNumber
	})

	errStartingServer := http.ListenAndServe(":8080", router)
	if errStartingServer != nil {
		log.Fatalln("Failed to start up the HTTP service")
	}

}
