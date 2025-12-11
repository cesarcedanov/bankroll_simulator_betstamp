package service

import (
	"bankroll_simulator_betstamp/model"
	"math/rand"
	"sync"
	"time"
)

// Worker Pool Job
type Job struct {
	IterationID int
	Request     model.SimulationRequest
}

// Job Result
type JobResult struct {
	FinalBankroll float64
	Busted        bool
}

// Worker Pool
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan JobResult
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, numWorkers*2),
		results:    make(chan JobResult, numWorkers*2),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for job := range wp.jobs {
		result := RunSingleIteration(job.Request, rng)
		wp.results <- result
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}
