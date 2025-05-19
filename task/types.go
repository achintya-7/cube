package task

import (
	"time"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID            uuid.UUID         `json:"id" binding:"required"`
	ContainerID   string            `json:"container_id"`
	Name          string            `json:"name" binding:"required"`
	State         State             `json:"state" binding:"required"`
	Image         string            `json:"image" binding:"required"`
	Memory        int               `json:"memory" `
	Disk          int               `json:"disk" `
	ExposedPorts  nat.PortSet       `json:"exposed_ports"`
	PortBindings  map[string]string `json:"port_bindings"`
	RestartPolicy string            `json:"restart_policy"`
	StartTime     time.Time         `json:"start_time"`
	FinishTime    time.Time         `json:"finish_time"`
}

type TaskEvent struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	State     State     `json:"state" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
	Task      Task      `json:"task" binding:"required"`
}

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	ExposedPorts  nat.PortSet
	Cmd           []string
	Image         string
	Cpu           float64
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy string
	ContainerID   string
}

func NewConfig(t *Task) *Config {
	return &Config{
		Name:          t.Name,
		Image:         t.Image,
		RestartPolicy: t.RestartPolicy,
	}
}

type Docker struct {
	Client *client.Client
	Config Config
}

func NewDocker(c *Config) *Docker {
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	return &Docker{
		Client: dc,
		Config: *c,
	}
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerID string
	Result      string
}
