package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestTimeName(t *testing.T) {
	s := "2021-06-02T10:51:12.725599800Z"
	fmt.Println(time.Now().UnixNano() / 1e6)
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	t1, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	fmt.Println("t1:", t1)
	t2, err := time.Parse("2006-01-02T15:04:05.000000000Z", s)
	fmt.Println("length:", len("2006-01-02T15:04:05.000000000Z"))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(t2)
	fmt.Println(s)
}

func TestFindContainer(t *testing.T) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	name := "/" + "redis"
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
			log.Println(container.Names[0], container.Ports)
		}
	}
}

func TestContainerStatus(t *testing.T) {

}

func TestContainerStatus2(t *testing.T) {
	//获取ctx
	ctx := context.Background()

	//cli客户端对象
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	//获取容器id 这个其实docker ps那个命令，不过我们只需要容器id
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	//遍历获取到的容器
	for _, container := range containers {
		fmt.Println("--------容器ID-------")
		fmt.Println(container.ID)
		fmt.Println(container.Image)
		fmt.Println(container.ImageID)
		fmt.Println("根据容器id获取容器的stats")
		//通过cli的ContainerStats方法可以获取到 docker stats命令的详细信息，其实是一个容器监控的方法
		//这个方法返回了容器使用CPU 内存 网络 磁盘等诸多信息
		containerStats, err := cli.ContainerStats(ctx, container.ID, false)
		if err != nil {
			panic(err)
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
		stats := map[string]interface{}{}
		json.Unmarshal(buf.Bytes(), &stats)
		fmt.Printf(newStr)
		log.Println(stats)
	}
}

func TestRunContainerTest(t *testing.T) {
	//Run a container

	//reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//io.Copy(os.Stdout, reader)
	var openPort nat.Port = "6379/tcp"
	//container开放端口
	var hostPort string = "7072"
	config := &container.Config{
		Image: "redis:6.0",
		ExposedPorts: nat.PortSet{
			openPort: struct{}{},
		},
	}

	// -m ${memory} --memory-swap -1
	// --log-driver fluentd --log-opt fluentd-address=${FluentdHost} --log-opt tag="{{.Name}}"
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			openPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			},
		},
		LogConfig: container.LogConfig{
			Type: "fluentd", // jsonfilelog.Name
			Config: map[string]string{
				"mode":            string(container.LogModeNonBlock),
				"fluentd-address": "172.10.1.2:24224",
				"tag":             "{{.Name}}",
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
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

func TestDisk(t *testing.T) {
	disk, err := cli.DiskUsage(ctx)
	log.Println(err)
	log.Println("LayersSize:", disk.LayersSize, float64(disk.LayersSize*1.0)/1024.0/1024.0/1024.0)
	log.Println("BuilderSize:", disk.BuilderSize, float64(disk.BuilderSize*1.0)/1024.0/1024.0/1024.0)

}

func TestCompare(t *testing.T) {
	s1 := "/data/web"
	s2 := "/usr/local/openresty/extr_lua"
	i := strings.Compare(s2, s1)
	log.Println(i)

}
