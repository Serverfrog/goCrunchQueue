package internal

import log "github.com/sirupsen/logrus"

const (
	Added          = 1
	Removed        = 2
	Process        = 3
	Processed      = 4
	ErrLogUpdated  = 5
	InfoLogUpdated = 6
)

type Event struct {
	Id      int
	Item    QueueItem
	Message string
}

type EventListener interface {
	onEvent(event Event)
	getName() string
}

var eventHandler = &EventHandler{listener: []EventListener{}}

type EventHandler struct {
	listener []EventListener
}

func (e *EventHandler) handleEvent(event Event) {
	log.Debugf("Event Fired: %v", event)
	for _, eventListener := range e.listener {
		eventListener.onEvent(event)
	}
}

func (e *EventHandler) addListener(listener EventListener) {
	log.Debugf("Added Listener %v", listener.getName())
	e.listener = append(e.listener, listener)
}
