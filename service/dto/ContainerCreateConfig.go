package dto

import (
	"github.com/docker/go-connections/nat"
)

type ContainerCreateConfig struct {
	TaskId        string `json:"taskId"`
	ImageName     string `json:"imageName"`
	ContainerName string `json:"containerName"`
	Ports         []struct {
		IP          string `json:"IP"`
		PrivatePort string `json:"PrivatePort"`
		PublicPort  string `json:"PublicPort"`
		Type        string `json:"Type"`
	} `json:"ports"`
	Env     []string `json:"env"`
	Volumes []string `json:"volumes"`
	Mounts  []struct {
		Destination string `json:"Destination"`
		Mode        string `json:"Mode"`
		Propagation string `json:"Propagation"`
		RW          bool   `json:"RW"`
		Source      string `json:"Source"`
		Type        string `json:"Type"`
	} `json:"Mounts"`
	Memory       int64             `json:"memory"`
	LogType      string            `json:"logType"`
	LogConfigMap map[string]string `json:"logConfig"`
	Cover        bool              `json:"cover"`
}

func (e *ContainerCreateConfig) PortsToSet() (nat.PortSet, nat.PortMap) {
	ExposedPorts := nat.PortSet{}
	PortBindings := nat.PortMap{}

	for _, p := range e.Ports {
		var openPort = nat.Port(p.PrivatePort)
		ExposedPorts[openPort] = struct{}{}
		PortBindings[openPort] = []nat.PortBinding{
			{
				HostIP:   p.IP,
				HostPort: p.PublicPort,
			},
		}
	}
	return ExposedPorts, PortBindings
}

func (e ContainerCreateConfig) GetVolumes() (res []string) {
	//for _, v := range e.Volumes {
	//	tmp := v["Source"].(string) + ":" + v["Destination"].(string) + ":" + utils.If(v["RW"].(bool), "rw", "ro").(string)
	//	res = append(res, tmp)
	//}
	return e.Volumes
}
