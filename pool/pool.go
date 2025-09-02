package pool

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type WorkerPool interface {
	Submit(task func())
	SubmitWait(task func())
	Stop()
	StopWait()
}

type Pool struct {
	numberOfWorkers int
	taskQueue       chan *WorkTask
	wg              *sync.WaitGroup
	initWg          *sync.WaitGroup
	forceStopChan   chan struct{}
}

type WorkTask struct {
	doneChan chan struct{}
	task     func()
}

func NewWorkTask(doneChan chan struct{}, task func()) *WorkTask {
	return &WorkTask{doneChan: doneChan, task: task}
}

// инициализируем пул
func InitPool(workers int, queueSize int) WorkerPool {
	taskQueue := make(chan *WorkTask, queueSize)
	wg := &sync.WaitGroup{}
	initWg := &sync.WaitGroup{}
	forceStopChan := make(chan struct{})
	pool := NewPool(workers, taskQueue, wg, initWg, forceStopChan)

	log.Println("waiting for all workers are ready")
	initWg.Wait()
	log.Println("all workers are ready to work")

	return pool
}

func NewPool(numberOfWorkers int, taskQueue chan *WorkTask, wg *sync.WaitGroup, initWg *sync.WaitGroup, forceStopChan chan struct{}) *Pool {
	pool := &Pool{numberOfWorkers: numberOfWorkers, taskQueue: taskQueue, wg: wg, initWg: initWg, forceStopChan: forceStopChan}
	for i := range numberOfWorkers {
		wg.Add(1)
		initWg.Add(1)
		go pool.HandleWorker(i)
	}
	return pool
}

func (pool *Pool) Submit(task func()) {
	workTask := NewWorkTask(nil, task)
	pool.taskQueue <- workTask
}

func (pool *Pool) SubmitWait(task func()) {
	workTask := NewWorkTask(make(chan struct{}), task)
	pool.taskQueue <- workTask
	log.Println("waiting till task done")
	<-workTask.doneChan
	log.Println("task done")
}

func (pool *Pool) Stop() {
	log.Printf("stopping")
	close(pool.forceStopChan)
	pool.wg.Wait()
	log.Printf("current tasks is done, tasks from query is undone")
}

func (pool *Pool) StopWait() {
	log.Printf("stopping accepting new tasks")
	close(pool.taskQueue)
	pool.wg.Wait()
}

func (pool *Pool) HandleWorker(workerId int) {
	log.Printf("worker %v is starting", workerId)
	pool.initWg.Done()

	for task := range pool.taskQueue {
		select{ 
		case <-pool.forceStopChan:
			log.Printf("worker %v finished", workerId)
			pool.wg.Done()
			return
		default:
			pool.HandleTask(task)
		}
	}

	log.Printf("worker %v finished", workerId)
	pool.wg.Done()
}

// обертка для обработки функции. Также здесь вызывает функция из интерфейса хука
func (pool *Pool) HandleTask(workTask *WorkTask) {
	workTask.task()
	if workTask.doneChan != nil {
		close(workTask.doneChan)
	}
}

// некая логика обработки задачи
func SomeTask() {
	executingTime := time.Duration((rand.Intn(2) + 3)) * time.Second
	time.Sleep(executingTime)
	fmt.Println("Hello from task")
}
