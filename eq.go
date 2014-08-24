package eq

import "container/list"

type (
	EQ struct {
		Enqueue chan interface{}
		Dequeue chan interface{}
		closed  chan struct{}
		buffer  *list.List
	}
)

func NewEQ(size int) *EQ {

	q := &EQ{
		Enqueue: make(chan interface{}, size),
		Dequeue: make(chan interface{}, size),
		closed:  make(chan struct{}),
		buffer:  list.New(),
	}

	go func() {

		var next, v interface{}

	READ:
		select {
		case v = <-q.Enqueue:
			if q.buffer.Len() == 0 {
				next = v
			} else {
				q.buffer.PushBack(v)
				next = q.buffer.Remove(q.buffer.Front())
			}
			for {
				select {
				case v = <-q.Enqueue:
					q.buffer.PushBack(v)
				case q.Dequeue <- next:
					if q.buffer.Len() == 0 {
						goto READ
					}
					next = q.buffer.Remove(q.buffer.Front())
				case <-q.closed:
					goto CLOSE
				}
			}
		case <-q.closed:
			goto CLOSE
		}

	CLOSE:
		close(q.Enqueue)
		for v = range q.Enqueue {
			q.buffer.PushBack(v)
		}
		q.Dequeue <- next
		for q.buffer.Len() > 0 {
			q.Dequeue <- q.buffer.Remove(q.buffer.Front())
		}
		close(q.Dequeue)

	}()

	return q

}

func (q *EQ) Close() {
	q.closed <- struct{}{}
}
