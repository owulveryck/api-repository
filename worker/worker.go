package worker

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/repository"
	"github.com/owulveryck/api-repository/session"
)

var jobQueue chan Job

// NewJobQueue ...
func NewJobQueue(queueLength int) chan Job {
	return make(chan Job, queueLength)
}

// Job represents the job to be run
type Job struct {
	Payload       object.IDer
	TransactionID uuid.UUID
	Path          string
}

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
func (w worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				/*
					err := w.saver.Save(context.TODO(), job.Payload, job.path)
						job.ret <- returnCode{
							Err:  err,
							Code: job.Payload.ID(),
						}
				*/
				session.Upsert(context.TODO(), job.TransactionID, session.Element{
					ID:     job.Payload.ID(),
					Status: http.StatusCreated,
				})
				if err := repository.Save(context.TODO(), job.Payload, job.Path); err != nil {
					session.Upsert(context.TODO(), job.TransactionID, session.Element{
						ID:     job.Payload.ID(),
						Status: http.StatusInternalServerError,
						Err:    err.Error(),
					})
				} else {
					session.Upsert(context.TODO(), job.TransactionID, session.Element{
						ID:     job.Payload.ID(),
						Status: http.StatusOK,
					})
				}
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
