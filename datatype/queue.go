package datatype

import (
	log "github.com/Sirupsen/logrus"
	"sync"
)

var m *sync.Mutex = new(sync.Mutex)

type Node struct {
	Value Task
	Prev  *Node
	Next  *Node
}

type LinkedQueue struct {
	Head  *Node
	Tail  *Node
	Size  int
	LockQ *sync.Mutex
}

func (queue *LinkedQueue) Lock() {
	queue.LockQ.Lock()
}

func (queue *LinkedQueue) Unlock() {
	queue.LockQ.Unlock()
}

func (queue *LinkedQueue) SizeOf() int {
	queue.Lock()
	defer queue.Unlock()
	return queue.Size
}

//out from head
func (queue *LinkedQueue) Peek() Task {
	queue.Lock()
	defer queue.Unlock()
	if queue.Head == nil {
		log.Errorln("datatype.Peek() Empty queue.")
	}
	return queue.Head.Value
}

func (queue *LinkedQueue) Add(value Task) {
	queue.Lock()
	defer queue.Unlock()
	new_node := &Node{value, queue.Tail, nil}
	if queue.Head == nil {
		queue.Head = new_node
		queue.Tail = nil
	} else {
		queue.Tail.Next = new_node
		queue.Tail = new_node
	}
	queue.Size++
	new_node = nil
}

func (queue *LinkedQueue) Remove() {
	queue.Lock()
	defer queue.Unlock()
	log.Infoln("************* datatype.Remove() ***************")
	if queue.Head == nil {
		panic("Empty queue.:")
	}
	first_node := queue.Head
	queue.Head = first_node.Next
	first_node.Next = nil
	queue.Size--
	first_node = nil
}
