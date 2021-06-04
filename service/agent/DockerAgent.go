package agent

import (
	"context"
	"docker-agent/service/conf"
	"docker-agent/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
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
	conf.DockerInfo, err = cli.Info(ctx)
	if err != nil {
		log.Println("DockerInfo.err:", err)
	}
}

func PostDockerInfo() {
	log.Println("PostDockerInfo.info:", conf.DockerInfo.Name)
	utils.PostData("/reg", conf.DockerInfo)
	log.Println("PostDockerInfo.post success:", conf.DockerInfo.Name)
}

func PostContainers() {
	containers, err := ContainerList()
	if err != nil {
		log.Println("PostContainers.err:", err)
		return
	}

	data := map[string]interface{}{
		"ID":        conf.DockerInfo.ID,
		"Name":      conf.DockerInfo.Name,
		"conainers": containers,
	}
	utils.PostData("/containers", data)
	log.Println("PostContainers size:", len(containers))
}

func PostContainersStats() {
	err, stats := ContainersStats()
	if err != nil {
		log.Println("PostContainersStats.err:", err)
		return
	}

	data := map[string]interface{}{
		"ID":    conf.DockerInfo.ID,
		"Name":  conf.DockerInfo.Name,
		"Stats": stats,
		"Time":  time.Now().Unix(),
	}
	utils.PostData("/containers/stats", data)
	log.Println("PostContainersStats size:", len(stats))
}

func PostContainerStats(containerId string) {
	stats, err := ContainerStats(containerId)
	if err != nil {
		log.Println("PostContainerStats.err:", err)
		return
	}
	utils.PostData("/container/stats", stats)
	log.Println("PostContainerStats size:", len(stats))
}

func PostImageList() {
	images, err := ImageList()
	if err != nil {
		log.Println("PostImageList.err:", err)
		return
	}

	data := map[string]interface{}{
		"ID":     conf.DockerInfo.ID,
		"Name":   conf.DockerInfo.Name,
		"Images": images,
	}
	utils.PostData("/images", data)
	log.Println("PostImageList size:", len(images))
}
