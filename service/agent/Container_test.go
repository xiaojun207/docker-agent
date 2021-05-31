package agent

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestContainerLogs(t *testing.T) {
	log.Println(time.Now().String())
	containerId := "302c6c35f81699fd8d3c02e903cea24b4979b51080abd8dcbbb06bf3de45016b"
	logs, err := ContainerLogs(containerId, "100", "2021-05-30T04:07:11Z")
	log.Println(err)
	fmt.Println(logs)
}

func TestContainerCreate(t *testing.T) {
	imageName := "redis:6.0"
	containerName := "redis6"

	ContainerRemove(containerName)

	ports := map[string]string{
		"6389": "6379",
	}
	env := []string{}

	volumes := []string{
		"/tmp:/tmp",
	}
	resp, err := ContainerCreate(imageName, containerName, ports, env, volumes)
	if err != nil {
		log.Println("ContainerCreate.err:", err)
	}

	log.Println(resp)
	err = ContainerStart(containerName)

	if err != nil {
		log.Println("ContainerStart.err:", err)
	}
}
