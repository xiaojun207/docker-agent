package utils

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"log"
	"os"
)

var dockerInfo types.Info

func RegDocker(cli *client.Client, ctx context.Context) {
	var err error
	dockerInfo, err = cli.Info(ctx)
	if err != nil {
		panic(err)
	}

	PostData("/dockerMgrApi/reg", dockerInfo)
	//log.Println(dockerInfo.ID)
}

func getTask(cli *client.Client, ctx context.Context) {
	GetData("/dockerMgrApi/task?ID=" + dockerInfo.ID)
}

//Run a container
func RunContainer(cli *client.Client, ctx context.Context) {
	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

//Run a container in the background
func RunContainer2(cli *client.Client, ctx context.Context) {
	imageName := "bfirsh/reticulate-splines"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	log.Println(resp.ID)
}

//List and manage containers
func PostContainers(cli *client.Client, ctx context.Context) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	//for _, container := range containers {
	//log.Println(container.Names[0], container.Ports)
	//}
	data := map[string]interface{}{
		"ID":        dockerInfo.ID,
		"conainers": containers,
	}
	PostData("/dockerMgrApi/containers", data)
	//log.Println(dockerInfo.ID)
}

//Stop all running containers
func StopAllRunningContainers(cli *client.Client, ctx context.Context) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		log.Println("Stopping container ", container.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		log.Println("Success")
	}
}

//Stop all running containers
func RmAllRunningContainers(cli *client.Client, ctx context.Context) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		log.Println("Stopping container ", container.ID[:10], "... ")

		if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
		log.Println("Success")
	}
}
