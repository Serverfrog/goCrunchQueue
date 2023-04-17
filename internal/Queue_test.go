package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestQueue_GetAll(t *testing.T) {
	HandleErrorB(os.Remove(QUEUE_FILE))
	queue := NewQueue()
	queue.clear()
	items := []QueueItem{
		createQueueItem("test1", "test1"),
		createQueueItem("test2", "test2"),
		createQueueItem("test3", "test3"),
		createQueueItem("test4", "test4"),
		createQueueItem("test5", "test5"),
		createQueueItem("test6", "test6"),
	}
	for _, item := range items {
		queue.Push(item)
	}
	assert.Equal(t, len(items), queue.Length())

	allItems := queue.GetAll()
	assert.Equal(t, len(items), len(allItems))
	for i, item := range allItems {
		assert.Equal(t, items[i], item)
	}

	queue = NewQueue()
	assert.Equal(t, items[0], queue.Pop())
	assert.Equal(t, items[1], queue.Pop())
	assert.Equal(t, items[2], queue.Pop())
	assert.Equal(t, items[3], queue.Pop())
	assert.Equal(t, items[4], queue.Pop())
	assert.Equal(t, items[5], queue.Pop())

	assert.Equal(t, 0, queue.Length())

	assert.Equal(t, QueueItem{}, queue.Pop())
	queue.clear()
}
