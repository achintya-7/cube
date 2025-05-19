package main

import (
	"fmt"
	"time"

	"github.com/achintya-7/cube/task"
	"github.com/achintya-7/cube/worker"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container",
		Image: "postgres:latest",
		Env:   []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres"},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Println("Error creating container:", result.Error)
		return nil, nil
	}

	fmt.Printf("Container %s created successfully\n", result.ContainerID)

	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Println("Error stopping container:", result.Error)
		return nil
	}

	fmt.Printf("Container %s stopped successfully\n", id)

	return &result
}

func main() {
	db := make(map[uuid.UUID]*task.Task)

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "test-task",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	fmt.Println("Starting task")
	w.AddTask(t)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.ContainerID = result.ContainerID
	
	fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 30)
	
	fmt.Printf("stopping task %s\n", t.ID)
	
	t.State = task.Completed
	w.AddTask(t)	
	result = w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
}
