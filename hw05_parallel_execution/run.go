package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n == 0 {
		return nil
	}

	var errCount uint64
	taskChan := make(chan Task)

	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Wait()
	}()

	// run workers
	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case task := <-taskChan:
					// start task
					err := task()
					if err != nil {
						atomic.AddUint64(&errCount, 1)
					}
				case <-ctx.Done():
					wg.Done()
					return
				}
			}
		}()
		wg.Add(1)
	}
	// start tasks
	for _, task := range tasks {
		currErrCount := atomic.LoadUint64(&errCount)
		if currErrCount >= uint64(m) {
			return ErrErrorsLimitExceeded
		}
		// give a new task to worker
		taskChan <- task
	}

	return nil
}
