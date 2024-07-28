package main

import (
	"fmt"
	queues "harshsinghvi/go-queues/queue"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	Square = queues.NewQueue[int, int](100, func(i1, i2 int, i3 ...interface{}) int {
		time.Sleep(100 * time.Millisecond)
		result := i2 * i2
		fmt.Printf("Consumer %d input: %d result: %d \n", i1, i2, result)
		return result
	}, 10)

	SquarePriority = queues.NewQueue[int, int](100, func(i1, i2 int, i3 ...interface{}) int {
		time.Sleep(100 * time.Millisecond)
		result := i2 * i2
		fmt.Printf("Consumer %d input: %d result: %d \n", i1, i2, result)
		return result
	}, 10)
)

// var arr = map[string]queues.Queue[any, any]{
// 	"Square": queues.NewQueue[int, int](100, func(i1, i2 int, i3 ...interface{}) int {
// 		time.Sleep(100 * time.Millisecond)
// 		result := i2 * i2
// 		fmt.Printf("Consumer %d input: %d result: %d \n", i1, i2, result)
// 		return result
// 	}, 10),

// 	SquarePriority: queues.NewQueue[int, int](100, func(i1, i2 int, i3 ...interface{}) int {
// 		time.Sleep(100 * time.Millisecond)
// 		result := i2 * i2
// 		fmt.Printf("Consumer %d input: %d result: %d \n", i1, i2, result)
// 		return result
// 	}, 10),
// }

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("GOT SIG Stopping Queues")
			Square.Close()
		}
	}()
}

func main() {
	for i := 0; i < 1000; i++ {
		Square.Publish(i)
	}

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Println("Job =>", Square.Job(i))
			wg.Done()
		}(i)
	}


	// Async
	Square.Publish(10)
	Square.Publish(10)
	Square.Publish(10)
	Square.Publish(10)
	Square.Publish(10)

	// Sync
	fmt.Println(Square.Job(10))
	fmt.Println(Square.Job(10))
	fmt.Println(Square.Job(10))
	fmt.Println(Square.Job(10))
	fmt.Println(Square.Job(10))

	wg.Wait()
	Square.Close()
	Square.Wait()
}
