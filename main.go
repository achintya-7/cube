package main

import (
	"log"
	"time"

	"github.com/achintya-7/cube/task"
	"github.com/achintya-7/cube/worker"
	workerApi "github.com/achintya-7/cube/worker/api"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	workerServer := workerApi.NewServer("localhost", 8080, &w)

	go runTasks(&w)
	workerServer.Start()
}

func runTasks(w *worker.Worker) {
	log.Println("Starting task runner")

	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v", result.Error)
			}
		} else {
			log.Println("No tasks in queue")
		}

		log.Println("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}
