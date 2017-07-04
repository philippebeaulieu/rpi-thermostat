package statequeue

import "github.com/philippebeaulieu/rpi-thermostat/thermostat"

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		states:   make([]thermostat.State, size),
		size:     size,
		index:    0,
		count:    0,
		head:     0,
		overflow: false,
	}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	states   []thermostat.State
	size     int
	index    int
	count    int
	head     int
	overflow bool
}

// Push adds a node to the queue.
func (q *Queue) Push(n thermostat.State) {
	q.states[q.index] = n

	q.index++

	if q.count < q.size {
		q.count++
	}

	if q.index >= q.size {
		q.overflow = true
		q.index = 0
	}

	if q.overflow == true {
		q.head++

		if q.head >= q.size {
			q.head = 0
		}
	}

}

// ToArray returns a copy of the internal queue array
func (q *Queue) ToArray() []thermostat.State {
	states := make([]thermostat.State, q.count)
	i := q.head

	for n := 0; n < q.count; n++ {
		states[n] = q.states[i]
		i = q.nextIndex(i)
	}

	return states
}

func (q *Queue) nextIndex(i int) int {
	if i >= q.size {
		return 0
	}

	return i + 1
}
