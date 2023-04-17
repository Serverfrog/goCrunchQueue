package internal

// reads config, starts the queue, api and hold

import (
	log "github.com/sirupsen/logrus"
)

var queue *Queue

func ValidateSetup() error {
	err := crunchyValidation()
	if err != nil {
		return nil
	}
	return nil
}

func StartApplication(version string, configPath *string) {
	ReadConfig(configPath)
	queue = NewQueue()
	StartServer(version)

}

func StopApplication() {
	queue.save()
}

type Worker struct {
	currentItem QueueItem
	EventListener
}

var worker = Worker{currentItem: emptyQueueItem()}

func (w *Worker) onEvent(event Event) {
	switch event.id {
	case Added:
		worker.process()
	default:
	}
}

func (w *Worker) process() {
	if !w.currentItem.isEmpty() {
		return
	}
	peeked := queue.Peek()
	if peeked.isEmpty() {
		return
	}
	go w.work()
}
func (w *Worker) work() {
	if w.currentItem.isEmpty() {
		return
	}
	for !w.currentItem.isEmpty() {
		log.Infof("Work on item %v | url:%v", w.currentItem.Id, w.currentItem.CrunchyrollUrl)
		execCrunchy(w.currentItem)
		log.Infof("Finished Work on item %v | url:%v", w.currentItem.Id, w.currentItem.CrunchyrollUrl)
		w.currentItem = queue.Pop()
	}
}
