package workerpool

import (
	"sync"
)

type Task struct {
	closer *sync.Once
	action func()
	done   chan struct{}
}

func NewTask(action func()) Task {
	return Task{action: action, closer: new(sync.Once), done: make(chan struct{})}
}

func (t Task) Do() {
	t.action()
	t.Stop()
}

func (t Task) Wait() {
	<-t.done
}

func (t Task) Stop() {
	t.closer.Do(func() {
		close(t.done)
	})
}

func (t Task) IsCompleted() bool {
	select {
	case <-t.done:
		return true
	default:
		return false
	}
}
