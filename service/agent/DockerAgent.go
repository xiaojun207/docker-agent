package agent

import (
	"bytes"
	"context"
	"docker-agent/service/conf"
	"docker-agent/utils"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"os"
	"time"
)

type DockerAgent struct {
	ctx        context.Context
	DockerInfo types.Info
	cli        *client.Client
}

var ctx = context.Background()
var cli *client.Client

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
}

func RegDocker() {
	var err error
	conf.DockerInfo, err = cli.Info(ctx)
	if err != nil {
		panic(err)
	}

	utils.PostData("/reg", conf.DockerInfo)
	log.Println("RegDocker:", conf.DockerInfo.Name)
}

func PullImage(imageName string) error {
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, out)
	return err
}

//Run a container in the background
func ContainerCreate(imageName, containerName string, hostPort, appPort string) (string, error) {
	var openPort = nat.Port(appPort)

	config := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			openPort: struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			openPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		log.Println("ContainerCreate.err:", err)
		return "", err
	}
	return resp.ID, err
}

//Run a container in the background
func RunContainer(imageName, containerName string, hostPort, appPort string) error {
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	CleanOldContainer(containerName)

	containerId, err := ContainerCreate(imageName, containerName, hostPort, appPort)
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	log.Println("containerId:", containerId)
	return err
}

// 根据容器名称，停止并删除容器
func CleanOldContainer(containerName string) {
	contain, err := FindContainer(containerName)
	if err != nil {
		return
	}
	containerId := contain.ID
	timeout := 0 * time.Second
	cli.ContainerStop(ctx, containerId, &timeout)
	cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
}

func ContainerStart(containerId string) error {
	return cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
}

func ContainerRestart(containerId string) error {
	timeout := 0 * time.Second
	return cli.ContainerRestart(ctx, containerId, &timeout)
}

func ContainerStop(containerId string) error {
	timeout := 0 * time.Second
	return cli.ContainerStop(ctx, containerId, &timeout)
}

func ContainerRemove(containerId string) error {
	timeout := 0 * time.Second
	cli.ContainerStop(ctx, containerId, &timeout)
	return cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
}

func FindContainer(containerName string) (types.Container, error) {
	name := "/" + containerName
	option := types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", name),
		),
	}

	containers, err := cli.ContainerList(ctx, option)
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		log.Println("[" + container.Names[0] + "]")
		if container.Names[0] == name {
			return container, nil
		}
	}
	return types.Container{}, errors.New("Not found container by name of " + containerName)
}

//List and manage containers
func PostContainers() {
	option := types.ContainerListOptions{
		All: true,
	}
	containers, err := cli.ContainerList(ctx, option)
	if err != nil {
		panic(err)
	}
	//for _, container := range containers {
	//log.Println(container.Names[0], container.Ports)
	//}
	data := map[string]interface{}{
		"ID":        conf.DockerInfo.ID,
		"Name":      conf.DockerInfo.Name,
		"conainers": containers,
	}
	utils.PostData("/containers", data)
	log.Println("PostContainers size:", len(containers))
}

//Stop all running containers
func StopRunningContainer(containerId string) {
	log.Println("Stopping container ", containerId, "... ")
	if err := cli.ContainerStop(ctx, containerId, nil); err != nil {
		panic(err)
	}
	log.Println("Success")
}

//Stop all running containers
func StopAllRunningContainers() {
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
func RmAllRunningContainers() {
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

func ContainerStats(containerId string) (map[string]interface{}, error) {
	//通过cli的ContainerStats方法可以获取到 docker stats命令的详细信息，其实是一个容器监控的方法
	//这个方法返回了容器使用CPU 内存 网络 磁盘等诸多信息
	containerStats, err := cli.ContainerStats(ctx, containerId, false)
	if err != nil {
		return nil, err
	}
	/**
	ContainerStats的返回的结构如下 注意这个Body的类型是io.ReadCloser 好奇怪的类型 下面我们给他转成json
	type ContainerStats struct {
		Body   io.ReadCloser `json:"body"`
		OSType string        `json:"ostype"`
	}
	*/
	//fmt.Println(containerStats)
	//fmt.Println("containerStats.Body的内容是: ", containerStats.Body)
	buf := new(bytes.Buffer)
	//io.ReadCloser 转换成 Buffer 然后转换成json字符串
	buf.ReadFrom(containerStats.Body)
	newStr := buf.String()
	//fmt.Printf(newStr)
	res := make(map[string]interface{})
	err = json.Unmarshal([]byte(newStr), &res)

	return res, err
}
