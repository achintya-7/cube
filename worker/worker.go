package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/achintya-7/cube/task"
)

func (w *Worker) CollectStats() {
	fmt.Println("I will collect stats")
}

func (w *Worker) RunTask() task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("No task to run")
		return task.DockerResult{Error: nil}
	}

	taskQueued := t.(task.Task)

	taskPersisted := w.Db[taskQueued.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = taskPersisted
	}

	var result task.DockerResult
	if !task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		err := fmt.Errorf("invalid state transition from %v to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
	}

	switch taskQueued.State {
	case task.Scheduled:
		result = w.StartTask(taskQueued)

	case task.Completed:
		result = w.StopTask(taskQueued)

	default:
		result.Error = fmt.Errorf("invalid state %v for task %v", taskQueued.State, taskQueued.ID)
	}

	return result
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	t.StartTime = time.Now().UTC()
	config := task.NewConfig(&t)
	dckr := task.NewDocker(config)

	result := dckr.Run()
	if result.Error != nil {
		log.Printf("Error starting task %s: %v", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}

	t.ContainerID = result.ContainerID
	t.State = task.Running
	w.Db[t.ID] = &t

	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	dckr := task.NewDocker(config)

	result := dckr.Stop(t.ContainerID)
	if result.Error != nil {
		log.Printf("Error stopping task %s: %v", t.ID, result.Error)
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed

	w.Db[t.ID] = &t

	log.Printf("Stopped and removed container %s from task %s", t.ContainerID, t.ID.String())

	return result
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) GetTasks() []*task.Task {
	tasks := make([]*task.Task, 0, len(w.Db))

	for _, t := range w.Db {
		tasks = append(tasks, t)
	}

	return tasks
}
