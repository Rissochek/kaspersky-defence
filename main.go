package main

import (
	"log"
	"time"

	"github.com/Rissochek/kaspersky-sandbox/utils"
	"github.com/Rissochek/kaspersky-sandbox/pool"
)

type Server struct {
	pool pool.WorkerPool
}

func NewServer(pool pool.WorkerPool) *Server {
	return &Server{pool: pool}
}

func main() {
	workers := utils.GetKeyFromEnv("WORKERS")
	queueSize := utils.GetKeyFromEnv("QUEUE_SIZE")
	workerPool := pool.InitPool(workers, queueSize)

	//количество задач которое нужно выполнить
	tasksTODO := utils.GetKeyFromEnv("TASKS")

	for i := range tasksTODO {
		log.Printf("%v task is submitting", i)
		time.Sleep(200 * time.Millisecond)
		workerPool.Submit(pool.SomeTask)
	}
	//проверяем корректную отработку логики остановки
	workerPool.StopWait()
}