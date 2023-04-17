package internal

import (
	"encoding/json"
	"os"
	"sync"
)

var queueMux sync.Mutex
var fileMux sync.Mutex

type Queue struct {
	queue []QueueItem
	path  string
}

const QUEUE_FILE = "./queue.json"

func NewQueue() *Queue {
	fileMux.Lock()
	defer fileMux.Unlock()
	queue := &Queue{path: QUEUE_FILE}
	HandleFatalErrorf(queue.load(), "Could not load Queue File. %v")
	return queue
}

func (q *Queue) load() error {
	data, err := os.ReadFile(q.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	HandleErrorB(json.Unmarshal(data, &q.queue))
	return nil
}

func (q *Queue) save() {
	fileMux.Lock()
	defer fileMux.Unlock()
	data := HandleError(json.Marshal(q.queue))
	HandleErrorB(os.WriteFile(q.path, data, 0644))
}

func (q *Queue) Push(item QueueItem) {
	queueMux.Lock()
	defer queueMux.Unlock()
	q.queue = append(q.queue, item)
	q.save()
}

func (q *Queue) Pop() QueueItem {

	queueMux.Lock()
	defer queueMux.Unlock()
	if len(q.queue) == 0 {
		return emptyQueueItem()
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	q.save()
	return item
}

func (q *Queue) Peek() QueueItem {
	queueMux.Lock()
	defer queueMux.Unlock()
	if len(q.queue) == 0 {
		return emptyQueueItem()
	}
	return q.queue[0]
}

func (q *Queue) GetAll() []QueueItem {
	queueMux.Lock()
	defer queueMux.Unlock()
	items := make([]QueueItem, len(q.queue))
	copy(items, q.queue)
	return items
}

func (q *Queue) Length() int {
	queueMux.Lock()
	defer queueMux.Unlock()
	return len(q.queue)
}

func (q *Queue) clear() {
	queueMux.Lock()
	defer queueMux.Unlock()
	q.queue = []QueueItem{}
	q.save()
}
