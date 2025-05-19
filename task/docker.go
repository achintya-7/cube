package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stdcopy"
)

func (d *Docker) Run() DockerResult {
	ctx := context.TODO()

	reader, err := d.Client.ImagePull(ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image: %s with error: %v", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}

	r := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
	}

	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container: %s with error: %v", d.Config.Name, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Error starting container: %s with error: %v", d.Config.Name, err)
		return DockerResult{Error: err}
	}

	// dd.Config.Runtime.ContainerID = resp.ID

	// out, err := cli.ContainerLogs(
	// ctx,
	// resp.ID,
	// types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
	// )
	// if err != nil {
	// log.Printf("Error getting logs for container %s: %v\n"
	// return DockerResult{Error: err}
	// , resp.ID, err)
	// }
	// stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	// return DockerResult{ContainerId: resp.ID, Action: "start"
	// ➥ Result: "success"}

	stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
	return DockerResult{ContainerID: resp.ID, Action: "start", Result: "success"}
}

func (d *Docker) Stop(id string) DockerResult {
	log.Println("Attempting to stop container: ", id)
	ctx := context.TODO()

	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container: %s with error: %v", id, err)
		return DockerResult{Error: err}
	}

	containerRemoveOptions := container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}
	err = d.Client.ContainerRemove(ctx, id, containerRemoveOptions)
	if err != nil {
		log.Printf("Error removing container: %s with error: %v", id, err)
		return DockerResult{Error: err}
	}

	return DockerResult{ContainerID: id, Action: "stop", Result: "success", Error: nil}
}
