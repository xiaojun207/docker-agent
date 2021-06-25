package agent

import (
	"context"
	"docker-agent/service/dto"
	"encoding/binary"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/daemon/logger/jsonfilelog"
	"log"
	"strconv"
	"testing"
)

func TestName(t *testing.T) {

}

func TestContainerLogs(t *testing.T) {
	containerId := "redis6"
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
	conf := dto.ContainerCreateConfig{}
	conf.ImageName = "redis:6.0"
	conf.ContainerName = "redis6"

	ContainerRemove(conf.ContainerName)

	conf.Ports = []struct {
		IP          string `json:"IP"`
		PrivatePort string `json:"PrivatePort"`
		PublicPort  string `json:"PublicPort"`
		Type        string `json:"Type"`
	}{
		{
			PrivatePort: "6379",
			PublicPort:  "6379",
		},
	}

	conf.Env = []string{}

	conf.Volumes = []string{
		"/tmp:/tmp",
	}

	conf.LogType = jsonfilelog.Name // jsonfilelog.Name, fluentd.Name
	conf.LogConfigMap = map[string]string{
		"mode": string(container.LogModeNonBlock),
	}

	resp, err := ContainerCreate(conf)
	if err != nil {
		log.Println("ContainerCreate.err:", err)
	}

	log.Println(resp)
	err = ContainerStart(conf.ContainerName)

	if err != nil {
		log.Println("ContainerStart.err:", err)
	}
}
