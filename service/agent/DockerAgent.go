package agent

import (
	"context"
	"docker-agent/service/conf"
	"docker-agent/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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
	//os.Setenv("DOCKER_CERT_PATH", "")
	//os.Setenv("DOCKER_HOST", "")
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	conf.DockerInfo, err = cli.Info(ctx)
	if err != nil {
		log.Println("DockerInfo.err:", err)
	}
	//cli.NetworkList(ctx, types.NetworkListOptions{})
}

func GetAgentConfig() {
	err, res := utils.GetData("/config")
	if err != nil {
		log.Println("GetAgentConfig.err:", err)
	}
	config := res["data"].(map[string]interface{})
	log.Println("GetAgentConfig.res:", res)
	agentConfig := (config["agentConfig"]).(map[string]interface{})

	conf.TaskFrequency = time.Duration(agentConfig["TaskFrequency"].(float64)) * time.Second

	log.Println("conf.TaskFrequency:", conf.TaskFrequency)
}

func PostDockerInfo() {
	log.Println("PostDockerInfo.info:", conf.DockerInfo.Name)
	utils.PostData("/reg", conf.DockerInfo)
	log.Println("PostDockerInfo.post success:", conf.DockerInfo.Name)
}

func SystemPrune() error {
	resp, err := cli.ImagesPrune(ctx, filters.NewArgs())
	if err != nil {
		return err
	}
	log.Println("ImagesPrune ", resp)
	cache, err := cli.BuildCachePrune(ctx, types.BuildCachePruneOptions{})
	if err != nil {
		return err
	}
	log.Println("BuildCachePrune", cache)
	cresp, err := cli.ContainersPrune(ctx, filters.NewArgs())
	if err != nil {
		return err
	}
	log.Println("ContainersPrune", cresp)
	nresp, err := cli.NetworksPrune(ctx, filters.NewArgs())
	if err != nil {
		return err
	}
	log.Println("NetworksPrune", nresp)
	vresp, err := cli.VolumesPrune(ctx, filters.NewArgs())
	if err != nil {
		return err
	}
	log.Println("VolumesPrune", vresp)

	return err
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
	//utils.PostData("/containers", data)
	SendWsMsg("docker.container.list", data)
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
	//utils.PostData("/containers/stats", data)
	SendWsMsg("docker.container.stats", data)
	log.Println("PostContainersStats size:", len(stats))
}

func PostContainerStats(containerId string) {
	stats, err := ContainerStats(containerId)
	if err != nil {
		log.Println("PostContainerStats.err:", err)
		return
	}
	utils.PostData("/container/stats", stats)
	SendWsMsg("docker.container.stats", stats)
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

	//utils.PostData("/images", data)
	SendWsMsg("docker.image.list", data)
	log.Println("PostImageList size:", len(images))
}
