package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"
)

// https://gist.github.com/moraes/2141121
type Node struct {
	Value []byte
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *CycleQueue {
	return &CycleQueue{
		nodes: make([]*Node, size),
		size:  size,
		stop:  make(chan bool),
	}
}

// CycleQueue is a basic FIFO queue based on a circular list that resizes as needed.
type CycleQueue struct {
	mux   sync.Mutex
	cmd   *exec.Cmd
	nodes []*Node
	size  int
	tail  int
	stop  chan bool
}

// Push adds a node to the queue.
func (q *CycleQueue) Add(n *Node) {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.nodes[q.tail%q.size] = n
	fmt.Print(",Add:", q.tail)
	q.tail++
}

func (q *CycleQueue) Get(index int) ([]byte, int) {
	q.mux.Lock()
	defer q.mux.Unlock()

	if index == -1 {
		index = q.tail - 1
	}
	if index < 0 {
		index = 0
	}

	var tailBytes bytes.Buffer

	total := 0
	for index < q.tail && total < q.size {
		node := q.nodes[index%q.size]
		fmt.Print(",Get:", index)
		tailBytes.Write(node.Value)
		index++
		total++
	}

	return tailBytes.Bytes(), index
}
