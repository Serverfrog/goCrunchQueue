package internal

// reads config, starts the queue, api and hold

import (
	log "github.com/sirupsen/logrus"
	"sync"
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
	eventHandler.addListener(&worker)
	eventHandler.addListener(&wsEventListener)
	worker.process()

	StartServer(version)
}

func StopApplication() {
	queue.save()
}

type Worker struct {
	currentItem QueueItem
	mux         sync.Mutex
	EventListener
}

var worker = Worker{currentItem: emptyQueueItem()}

func (w *Worker) onEvent(event Event) {
	switch event.Id {
	case Added:
		worker.process()
	default:
	}
}

func (w *Worker) getName() string {
	return "Worker"
}

func (w *Worker) process() {
	w.mux.Lock()
	if !w.currentItem.isEmpty() {
		w.mux.Unlock()
		return
	}
	peeked := queue.Peek()
	if peeked.isEmpty() {
		w.mux.Unlock()
		return
	}
	w.currentItem = queue.Pop()
	w.mux.Unlock()
	go w.work()
}
func (w *Worker) work() {
	w.mux.Lock()
	if w.currentItem.isEmpty() {
		w.mux.Unlock()
		return
	}
	w.mux.Unlock()
	for !w.currentItem.isEmpty() {
		log.Infof("Work on Item %v | url:%v", w.currentItem.Id, w.currentItem.CrunchyrollUrl)
		err := execCrunchy(w.currentItem)
		w.mux.Lock()
		defer w.mux.Unlock()
		if err != nil {
			log.Warnf("Error when Executing crunchy-cli. %v", err)
			w.currentItem = queue.Pop()
		} else {
			log.Infof("Finished Work on Item %v | url:%v", w.currentItem.Id, w.currentItem.CrunchyrollUrl)
			w.currentItem = queue.Pop()
		}
	}

}
