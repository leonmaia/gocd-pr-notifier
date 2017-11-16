package main

import (
	"fmt"
	"time"
)

func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				fmt.Printf("worker%d: Received work request\n", w.ID)
				for {
					time.Sleep(work.Delay)
					fmt.Printf("worker%d: Checking status of the pipeline\n", w.ID)
					result := isPipelineAvailable(work.PipelineName, work.StatusCheckURL, work.Auth)
					switch result {
					case true:
						work.Request.Do()
						return
					default:
						fmt.Printf("worker%d: Pipeline is busy trying again in a few seconds\n", w.ID)
					}
				}
			case <-w.QuitChan:
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}
