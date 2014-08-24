package eq

import (
	"fmt"
	"sync"
	"testing"
)

func ExampleQueue() {

	q := NewEQ(3)

	q.Enqueue <- 1
	q.Enqueue <- 2
	q.Enqueue <- 3
	q.Enqueue <- 4
	q.Enqueue <- 5
	q.Enqueue <- 6

	fmt.Println(<-q.Dequeue)
	fmt.Println(<-q.Dequeue)
	fmt.Println(<-q.Dequeue)
	fmt.Println(<-q.Dequeue)
	fmt.Println(<-q.Dequeue)
	fmt.Println(<-q.Dequeue)

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6

}

func ExampleConcurrentQueue() {

	q := NewEQ(3)

	var workers = 3
	var wg sync.WaitGroup
	wg.Add(workers * 2)

	for i := 0; i < workers; i++ {
		go func() {
			for i := 0; i < 10; i++ {
				fmt.Println(<-q.Dequeue)
			}
			wg.Done()
		}()
	}
	for i := 0; i < workers; i++ {
		go func() {
			for i := 0; i < 10; i++ {
				q.Enqueue <- i
			}
			wg.Done()
		}()
	}

	wg.Wait()

	// The order is preseved when GOMAXPROCS is 1
	//

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9

}

func BenchmarkEnqueues(b *testing.B) {

	q := NewEQ(1 << 10)

	for i := 0; i < b.N; i++ {
		q.Enqueue <- i
	}
}

func BenchmarkDequeues(b *testing.B) {

	q := NewEQ(1 << 10)

	for i := 0; i < 5<<20; i++ {
		q.Enqueue <- i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		<-q.Dequeue
	}

}

func BenchmarkEnqueueDequeue(b *testing.B) {

	q := NewEQ(1 << 10)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Enqueue <- i
		<-q.Dequeue
	}

}
