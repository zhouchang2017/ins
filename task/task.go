package task

import (
	"log"
	"ins/client"
	"ins/handlers"
)

type Taskable interface {
	Run()
	String() string
}

type task struct {
	client  *client.Client
	url     string
	handler handlers.Handler
}

func NewTask(client *client.Client, url string, handler handlers.Handler) *task {
	return &task{
		client:  client,
		url:     url,
		handler: handler,
	}
}

func (w *task) Run() {
	log.Printf("访问	%s\n", w.url)
	w.handler.Read(w.client, w.url)

	w.client.Wait()
}

func (w *task) String() string {
	return w.url
}
