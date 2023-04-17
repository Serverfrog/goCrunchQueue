package internal

const (
	Added     = 1
	Removed   = 2
	Process   = 3
	Processed = 4
)

type Event struct {
	id   int
	item QueueItem
}

type EventListener interface {
	onEvent(event Event)
}

var listener []EventListener

func handleEvent(event Event) {
	for _, eventListener := range listener {
		eventListener.onEvent(event)
	}
}
