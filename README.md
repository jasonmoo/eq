#eq
elastic, *threadsafe* queue

A simple queue that expands via container/list and uses channels for `Enqueue` and `Dequeue` actions.

###Example:

	// channels are buffered to provided size
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

`Close` closes the `Enqueue` channel and ensures all entries are put in the
`Dequeue` channel before closing it too.  Issuing `Close` followed by `range`
over the `Dequeue` channel will drain it.

	q := NewEQ(3)

	q.Enqueue <- 1
	q.Enqueue <- 2
	q.Enqueue <- 3
	q.Enqueue <- 4
	q.Enqueue <- 5
	q.Enqueue <- 6

	q.Close()

	for v := range q.Dequeue {
		fmt.Println(<-q.Dequeue)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6

