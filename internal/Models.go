package internal

import "github.com/google/uuid"

// Contains the Datamodels

// QueueItem is what we'll be storing in the queue.  It can be any struct
// as long as the fields you want stored are public.
type QueueItem struct {
	Name           string
	Id             string
	CrunchyrollUrl string
}

// QueueItemBuilder creates a new item and returns a pointer to it.
// This is used when we load a segment of the queue from disk.
func QueueItemBuilder() interface{} {
	return &QueueItem{}
}

func createQueueItem(name, url string) QueueItem {
	return QueueItem{
		Name:           name,
		Id:             uuid.New().String(),
		CrunchyrollUrl: url,
	}
}

func (q QueueItem) isEmpty() bool {
	return q.Name == "" && q.Id == "" && q.CrunchyrollUrl == ""
}
func emptyQueueItem() QueueItem {
	return QueueItem{}
}
