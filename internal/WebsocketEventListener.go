package internal

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"sync"
)

var wsEventListener = WebsocketEventListener{
	events: []Event{},
	mux:    sync.Mutex{},
}

type WebsocketEventListener struct {
	events []Event
	mux    sync.Mutex
}

func (ws *WebsocketEventListener) onEvent(event Event) {
	ws.mux.Lock()
	defer ws.mux.Unlock()
	log.Debugf("Added Event to Websocket Slice %v, %v", event.Id, event.Item.Id)
	ws.events = append(ws.events, Event{
		Id: event.Id,
		Item: QueueItem{
			Name:           event.Item.Name,
			Id:             event.Item.Id,
			CrunchyrollUrl: event.Item.CrunchyrollUrl,
		},
		Message: event.Message,
	})
}
func (ws *WebsocketEventListener) getName() string {
	return "WebsocketEventListener"
}

func (ws *WebsocketEventListener) getEvents() map[Event][]byte {
	ws.mux.Lock()
	defer ws.mux.Unlock()
	returnEvents := make([]Event, len(ws.events))

	copy(returnEvents, ws.events)
	ws.events = []Event{}
	contentMapping := map[Event][]byte{}
	for _, event := range returnEvents {
		contentMapping[event] = ws.getContentForEvent(event)
	}

	return contentMapping
}

func (ws *WebsocketEventListener) getContentForEvent(event Event) []byte {
	//if event.Id == InfoLogUpdated {
	//	path := fmt.Sprintf("%v/%v-out.txt", configuration.LogDestination, event.Item.Id)
	//	event.Message = string(HandleError(os.ReadFile(path)))
	//	return HandleError(json.Marshal(event))
	//}
	//if event.Id == ErrLogUpdated {
	//	path := fmt.Sprintf("%v/%v-err.txt", configuration.LogDestination, event.Item.Id)
	//	event.Message = string(HandleError(os.ReadFile(path)))
	//	return HandleError(json.Marshal(event))
	//}
	return HandleError(json.Marshal(event))
}
