package agent

import (
	"context"
	"docker-agent/service/conf"
	"docker-agent/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"log"
	"reflect"
	"sort"
	"time"
)

type DockerAgent struct {
	ctx        context.Context
	DockerInfo types.Info
	cli        *client.Client
}

var ctx = context.Background()
var cli *client.Client

var lastInfo types.Info
var lastContainers []types.Container
var lastContainersStats []map[string]interface{}
var lastImageList []types.ImageSummary

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
	//arrNetwork, err :=  cli.NetworkList(ctx, types.NetworkListOptions{})
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
	defer func() {
		if err := recover(); err != nil {
			log.Println("PostDockerInfo.err:", err)
		}
	}()

	var err error
	conf.DockerInfo, err = cli.Info(ctx)
	if err != nil {
		log.Println("DockerInfo.err:", err)
	}
	if !reflect.ValueOf(lastInfo).IsZero() {
		lastInfo.SystemTime = conf.DockerInfo.SystemTime
	}
	SendWsMsg("docker.info.systemTime", conf.DockerInfo.SystemTime)

	if !reflect.DeepEqual(lastInfo, conf.DockerInfo) {
		log.Println("PostDockerInfo.info:", conf.DockerInfo.Name)
		utils.PostData("/reg", conf.DockerInfo)
		//SendWsMsg("docker.info", conf.DockerInfo)
		log.Println("PostDockerInfo.post success:", conf.DockerInfo.Name)
	} else {
		log.Println("PostDockerInfo.DeepEqual, ignore it")
	}
	lastInfo = conf.DockerInfo
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
	defer func() {
		if err := recover(); err != nil {
			log.Println("PostContainers.err:", err)
		}
	}()
	containers, err := ContainerList()
	if err != nil {
		log.Println("PostContainers.list.err:", err)
		return
	}
	for _, container := range containers {
		//log.Println("container",i,",Ports:", container.Ports)
		p := PortSlice{container.Ports}
		sort.Sort(&p)
		//log.Println("container",i, ",Arr:", p.Arr)
		container.Ports = p.Arr
		m := MountSlice{container.Mounts}
		sort.Sort(&m)
		container.Mounts = m.Arr
	}

	if !reflect.DeepEqual(lastContainers, containers) {
		data := map[string]interface{}{
			"ID":         conf.DockerInfo.ID,
			"Name":       conf.DockerInfo.Name,
			"containers": containers,
		}
		//utils.PostData("/containers", data)
		SendWsMsg("docker.container.list", data)
		log.Println("PostContainers size:", len(containers))
	} else {
		log.Println("PostContainers.DeepEqual, ignore it")
	}
	lastContainers = containers
}

func PostContainersStats() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("PostContainersStats.err:", err)
		}
	}()
	err, stats := ContainersStats()
	if err != nil {
		log.Println("PostContainersStats.err:", err)
		return
	}
	if !reflect.DeepEqual(stats, lastContainersStats) {
		data := map[string]interface{}{
			"ID":    conf.DockerInfo.ID,
			"Name":  conf.DockerInfo.Name,
			"Stats": stats,
			"Time":  time.Now().Unix(),
		}
		//utils.PostData("/containers/stats", data)
		SendWsMsg("docker.container.stats", data)
		log.Println("PostContainersStats size:", len(stats))
	} else {
		log.Println("PostContainersStats.DeepEqual, ignore it")
	}
	lastContainersStats = stats
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
	defer func() {
		if err := recover(); err != nil {
			log.Println("PostImageList.err:", err)
		}
	}()
	images, err := ImageList()
	if err != nil {
		log.Println("PostImageList.err:", err)
		return
	}
	// 排序，确保数据不因顺序不同，而重复提交
	images = ImageSlice{images}.Sort()

	if !reflect.DeepEqual(images, lastImageList) {
		if !reflect.ValueOf(lastImageList).IsZero() {
			utils.CompareInter(images, lastImageList)
		}

		data := map[string]interface{}{
			"ID":     conf.DockerInfo.ID,
			"Name":   conf.DockerInfo.Name,
			"Images": images,
		}

		//utils.PostData("/images", data)
		SendWsMsg("docker.image.list", data)
		log.Println("PostImageList size:", len(images))
	} else {
		log.Println("PostImageList.DeepEqual, ignore it")
	}
	lastImageList = images
}
