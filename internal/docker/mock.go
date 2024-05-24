package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClientMock struct {
	client.APIClient
}

func NewDockerClientMock() *DockerClientMock {
	return &DockerClientMock{}
}

// Pretend a container is being restarted
func (d *DockerClientMock) ContainerRestart(ctx context.Context, container string, opts container.StopOptions) error {
	return nil
}
