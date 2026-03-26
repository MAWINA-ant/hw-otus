package hw05parallelexecution

import (
	"errors"
	"sync"
)

type LimitErrorCounter struct {
	mu    sync.Mutex
	value int
	err   error
}

func (ec *LimitErrorCounter) Increment() {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.value++
}

func (ec *LimitErrorCounter) GetValue() int {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	return ec.value
}

func (ec *LimitErrorCounter) SetError(err error) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.err = err
}

func (ec *LimitErrorCounter) GetError() error {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	return ec.err
}

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task)
	chDone := make(chan struct{})
	go func() {
		defer close(ch)
		for _, t := range tasks {
			select {
			case ch <- t:
			case <-chDone:
				return
			}
		}
	}()

	wg := sync.WaitGroup{}
	var errorCounter LimitErrorCounter
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range ch {
				err := task()
				if m > 0 {
					if err != nil {
						errorCounter.Increment()
					}
					if errorCounter.GetValue() >= m {
						errorCounter.SetError(ErrErrorsLimitExceeded)
						return
					}
				}
			}
		}()
	}
	wg.Wait()
	close(chDone)
	return errorCounter.err
}
