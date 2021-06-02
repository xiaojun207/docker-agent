package agent

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/docker/docker/api/types"
	"log"
	"strconv"
	"testing"
)

func TestContainerLogs(t *testing.T) {
	containerId := "drone-runner"
	//logs, err := ContainerLogs(containerId, "100", "2021-05-30T04:07:11Z")
	//log.Println(err)
	//fmt.Println(logs)

	i, err := cli.ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "40",
		Details:    true,
	})
	if err != nil {
		log.Fatal(err)
	}
	header := make([]byte, 8)
	Follow := true
	for Follow {
		_, err := i.Read(header)
		if err != nil {
			log.Fatal(err)
		}
		//var w io.Writer
		//switch header[0] {
		//case 1:
		//	w = os.Stdout
		//default:
		//	w = os.Stderr
		//}
		count := binary.BigEndian.Uint32(header[4:])
		dat := make([]byte, count)
		_, err = i.Read(dat)
		//fmt.Fprint(w, string(dat))
		fmt.Print(string(dat))
	}
}

func TestContainerLogs2(t *testing.T) {
	containerId := "drone-runner"
	ContainerLogFollow(containerId, func(timestamps int64, line string) bool {
		fmt.Println("timestamps:" + strconv.FormatInt(timestamps, 10))
		fmt.Println("line:" + line)
		return true
	})
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
