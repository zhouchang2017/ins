package worker

import (
	"github.com/labstack/gommon/log"
	"ins/task"
	"sync"
)

var wg sync.WaitGroup

var worker *Worker

type Worker struct {
	channel chan task.Taskable
}

func NewWorker() *Worker {
	worker = &Worker{channel: make(chan task.Taskable)}
	return worker
}

func (w *Worker) Add(t task.Taskable) {
	w.channel <- t
}

func (w *Worker) Run() {

	// defer wg.Done()
	for job := range w.channel {
		log.Printf("run job fetch %s", job)
		job.Run()
	}
}

func (w *Worker) Wait() {
	wg.Wait()
}

func (w *Worker) GetChannel() chan task.Taskable {
	return w.channel
}
