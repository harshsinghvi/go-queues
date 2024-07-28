package queues

import (
	"fmt"
	"sync"
)

type Input[IN any, OUT any] struct {
	Value         IN
	OutputChannel chan OUT
}

type Queue[IN any, OUT any] struct {
	InputChannel  chan Input[IN, OUT]
	OutputChannel chan OUT
	WorkerCount   int
	Workload      func(int, IN, ...interface{}) OUT
	WaitGroup     *sync.WaitGroup
	// Publish     func(IN) bool
	// Job           func(value IN) OUT
	// Close         func()
	// Wait          func()
}

func NewQueue[IN any, OUT any](workerCount int, workload func(int, IN, ...interface{}) OUT, buffer int) *Queue[IN, OUT] {
	var queue Queue[IN, OUT]
	queue.InputChannel = make(chan Input[IN, OUT], buffer)
	queue.WorkerCount = workerCount
	queue.Workload = workload
	queue.WaitGroup = &sync.WaitGroup{}

	queue.WaitGroup.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func(workerId int) {
			defer fmt.Printf("Stopped worker %d\n", workerId)
			defer queue.WaitGroup.Done()
			for input := range queue.InputChannel {
				output := queue.Workload(workerId, input.Value)
				if input.OutputChannel != nil {
					input.OutputChannel <- output
				}
			}
		}(i)
	}

	fmt.Printf("Started %d workers\n", workerCount)
	return &queue
}

func (queue *Queue[IN, OUT]) Wait() {
	queue.WaitGroup.Wait()
}

func (queue *Queue[IN, OUT]) Close() {
	close(queue.InputChannel)
}

func (queue *Queue[IN, OUT]) Publish(value IN) bool {
	queue.InputChannel <- Input[IN, OUT]{
		Value:         value,
		OutputChannel: nil,
	}
	return true
}

func (queue *Queue[IN, OUT]) Job(value IN) OUT {
	ch := make(chan OUT)
	defer close(ch)
	queue.InputChannel <- Input[IN, OUT]{
		Value:         value,
		OutputChannel: ch,
	}
	return <-ch
}
