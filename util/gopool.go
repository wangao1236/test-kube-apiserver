package util

import (
	"log"
	"sync"
)

// Task represents an in-process Goroutine task.
type Task interface {
	// Run method corresponds to Run method of Java's Runnable interface.
	Run()
}

// Executor defines the actions associated with the Goroutine pool.
type Executor interface {
	// Execute method corresponds to Execute method of Java's ExecutorService interface.
	Execute(task Task)
	// Wait waits for all the tasks to complete.
	Wait()
	// Done returns a channel which is closed if all the tasks completed.
	Done() chan struct{}
}

type executor struct {
	lock             sync.Mutex
	waitingTasks     []chan struct{}
	activeTasks      int64
	concurrencyLimit int64
	done             chan struct{}
}

func (ex *executor) Execute(task Task) {
	ex.start(task)
}

func (ex *executor) Wait() {
	<-ex.done
}

func (ex *executor) Done() chan struct{} {
	return ex.done
}

func (ex *executor) start(task Task) {
	startCh := make(chan struct{})
	stopCh := make(chan struct{})

	go startTask(startCh, stopCh, task)
	ex.enqueue(startCh)
	go ex.waitDone(stopCh)

}

// NewExecutor returns a new Executor.
func NewExecutor(concurrencyLimit int64) Executor {
	ex := &executor{
		waitingTasks:     make([]chan struct{}, 0),
		activeTasks:      0,
		concurrencyLimit: concurrencyLimit,
		done:             make(chan struct{}),
	}
	return ex
}

func startTask(startCh, stopCh chan struct{}, task Task) {
	defer close(stopCh)

	<-startCh
	log.Printf("task: %p is running", startCh)
	task.Run()
}

func (ex *executor) enqueue(startCh chan struct{}) {
	ex.lock.Lock()
	defer ex.lock.Unlock()

	if ex.concurrencyLimit == 0 || ex.activeTasks < ex.concurrencyLimit {
		log.Printf("Task: %p start executing", startCh)
		close(startCh)
		ex.activeTasks++
	} else {
		log.Printf("Task: %p start waitting", startCh)
		ex.waitingTasks = append(ex.waitingTasks, startCh)
	}
}

func (ex *executor) waitDone(stopCh chan struct{}) {
	<-stopCh

	ex.lock.Lock()
	defer ex.lock.Unlock()

	if len(ex.waitingTasks) == 0 {
		ex.activeTasks--
		if ex.activeTasks == 0 {
			close(ex.done)
		}
	} else {
		close(ex.waitingTasks[0])
		ex.waitingTasks = ex.waitingTasks[1:]
	}
}
