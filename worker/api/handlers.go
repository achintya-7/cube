package api

import (
	"log"

	"github.com/achintya-7/cube/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) getTasks(c *gin.Context) {
	tasks := s.worker.GetTasks()
	c.JSON(200, tasks)
}

func (s *Server) startTask(c *gin.Context) {
	var taskEvent task.TaskEvent
	err := c.ShouldBindJSON(&taskEvent)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	s.worker.AddTask(taskEvent.Task)
	log.Println("Task added:", taskEvent.Task.ID)
	c.JSON(200, taskEvent.Task)
}

func (s *Server) stopTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(400, gin.H{"error": "Task ID is required"})
		return
	}

	taskUUID, _ := uuid.Parse(taskID)
	taskToStop, ok := s.worker.Db[taskUUID]
	if !ok {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	s.worker.AddTask(taskCopy)

	log.Println("Task stopped:", taskCopy.ID)
	c.JSON(204, nil)
}
