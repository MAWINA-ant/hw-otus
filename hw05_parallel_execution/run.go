package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task, len(tasks))
	for _, t := range tasks {
		ch <- t
	}
	close(ch)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	countErrors := 0
	var resErr error
	for range n {
		wg.Add(1)
		go func(ch <-chan Task) {
			defer wg.Done()
			for {
				task, ok := <-ch
				if !ok {
					break
				}
				if task() != nil {
					mu.Lock()
					countErrors++
					mu.Unlock()
				}
				if countErrors >= m && m > 0 {
					resErr = ErrErrorsLimitExceeded
					break
				}
			}
		}(ch)
	}
	wg.Wait()
	return resErr
}
