package main

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	id         string
	n          int
	result     int
	workerPool *WorkerPool
	cache      *map[int]int
}

func (w *Worker) launch() {

	var wg sync.WaitGroup

	wg.Add(1)

	t := time.Now()

	go func() {
		defer func() {
			wg.Done()
		}()
		result := Fibonacci(w.n, *w.cache, &w.workerPool.mutex)
		w.result = result
	}()

	wg.Wait()

	elapsedTime := time.Since(t)

	go func() {
		w.workerPool.mutex.Lock()
		w.workerPool.cache[w.n] = w.result
		w.workerPool.mutex.Unlock()
	}()

	w.workerPool.quitChan <- w.id
	log.Printf(
		"[%v] Computing Fib(%d) :::: Time %v :::: Result %d ::::",
		w.id, w.n, elapsedTime, w.result,
	)

}

type WorkerPool struct {
	maxWorkers int
	cache      map[int]int
	workers    map[string]*Worker
	quitChan   chan string
	jobQueue   chan int
	mutex      sync.RWMutex
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		maxWorkers: maxWorkers,
		workers:    make(map[string]*Worker),
		quitChan:   make(chan string),
		jobQueue:   make(chan int),
		cache:      make(map[int]int),
		mutex:      sync.RWMutex{},
	}
}

func (wp *WorkerPool) startNewWorker(n int) {
	if wp.maxWorkers < len(wp.workers) {
		go func() {
			wp.jobQueue <- n
		}() // Not the best option in my opinion. A refactor may be necessary.
		return
	}

	log.Printf("[main] Starting a Worker to compute Fib(%d)", n)

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

			log.Printf(
				"[main] Deleting Worker with ID %v computing Fib(%d)",
				id, wp.workers[id].n,
			)

			delete(wp.workers, id)

		}
	}
}

func Fibonacci(n int, cache map[int]int, mutex *sync.RWMutex) int {

	mutex.RLock()
	res, exists := cache[n]
	mutex.RUnlock()

	if exists {
		return res
	}

	if n <= 1 {
		return n
	}

	return Fibonacci(n-1, cache, mutex) + Fibonacci(n-2, cache, mutex)

}
