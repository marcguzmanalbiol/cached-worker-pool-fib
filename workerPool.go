package main

import (
	"log"

	"github.com/google/uuid"
)

type Cache struct{}

type Worker struct {
	id         string
	n          int
	result     int
	workerPool *WorkerPool
	cache      *Cache
}

func (w *Worker) launch() {

	w.result = Fibonacci(w.n)
	w.workerPool.quitChan <- w.id
	log.Printf("Fib(%d) Worker Result %d", w.n, w.result)

}

type WorkerPool struct {
	maxWorkers int
	cache      Cache
	workers    map[string]*Worker
	quitChan   chan string
	jobQueue   chan int
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		maxWorkers: maxWorkers,
		cache:      Cache{},
		workers:    make(map[string]*Worker),
		quitChan:   make(chan string),
		jobQueue:   make(chan int),
	}
}

func (wp *WorkerPool) startNewWorker(n int) {
	if wp.maxWorkers < len(wp.workers) {
		go func() {
			wp.jobQueue <- n
		}() // Not the best option in my opinion. A refactor may be necessary.
		return
	}

	log.Println("Number of Working Workers", len(wp.workers))
	log.Printf("Starting a Worker to compute Fib(%d)", n)

	worker := Worker{
		id:         uuid.New().String(),
		n:          n,
		cache:      &wp.cache,
		workerPool: wp,
	}

	wp.workers[worker.id] = &worker

	go worker.launch()
}

func (wp *WorkerPool) StartListen() {

	for {
		select {
		case n := <-wp.jobQueue:
			wp.startNewWorker(n)

		case id := <-wp.quitChan:
			log.Printf("Deleting Worker with ID %v computing Fib(%d)", id, wp.workers[id].n)
			delete(wp.workers, id)

		}
	}
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}

	return Fibonacci(n-1) + Fibonacci(n-2)
}
