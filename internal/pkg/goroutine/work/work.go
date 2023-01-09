package work

import (
	"sync"
)

/*
	工作池
*/

type Worker struct {
	l          sync.RWMutex
	taskChan   chan func()
	workerChan chan func()
}

type Config struct {
	TaskChanCapacity   int
	WorkerChanCapacity int
	WorkerNum          int
}

func Init(config Config) *Worker {
	worker := &Worker{
		l:          sync.RWMutex{},
		taskChan:   make(chan func(), config.TaskChanCapacity),
		workerChan: make(chan func(), config.WorkerChanCapacity),
	}
	for i := 0; i < config.WorkerNum; i++ {
		go worker.work()
	}
	return worker
}

func (worker *Worker) SendTask(task func()) {
	worker.taskChan <- task
}

func (worker *Worker) work() {
	for task := range worker.taskChan {
		task()
	}
}
