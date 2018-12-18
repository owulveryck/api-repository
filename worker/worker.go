package worker

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/owulveryck/api-repository/dao"
	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/session"
)

var jobQueue chan Job

// NewJobQueue ...
func NewJobQueue(queueLength int) chan<- Job {
	jobQueue = make(chan Job, queueLength)
	return jobQueue
}

// Job represents the job to be run
// START_JOB OMIT
type Job struct {
	Payload       object.IDer // HL
	TransactionID uuid.UUID   // OMIT
	Path          string
}

// END_JOB OMIT

// Worker represents the worker that executes the job
type worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

// NewWorker ...
func newWorker(workerPool chan chan Job) worker {
	return worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
// START_POOL OMIT
func (w worker) Start() {
	go func() {
		for {
			// ...
			// register the current worker into the worker queue. // OMIT
			w.WorkerPool <- w.JobChannel // OMIT
			select {
			case job := <-w.JobChannel:
				session.Upsert(context.TODO(), job.TransactionID, session.Element{ // OMIT
					ID:     job.Payload.ID(),   // OMIT
					Status: http.StatusCreated, // OMIT
				}) // OMIT
				if err := dao.Save(context.TODO(), job.Payload, job.Path); err != nil { // HL
					//...
					session.Upsert(context.TODO(), job.TransactionID, session.Element{ // OMIT
						ID:     job.Payload.ID(),               // OMIT
						Status: http.StatusInternalServerError, // OMIT
						Err:    err.Error(),                    // OMIT
					}) // OMIT
				} else { // OMIT
					session.Upsert(context.TODO(), job.TransactionID, session.Element{ // OMIT
						ID:     job.Payload.ID(), // OMIT
						Status: http.StatusOK,    // OMIT
					}) // OMIT
				}
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// END_POOL OMIT

// Stop signals the worker to stop listening for work requests.
func (w worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
