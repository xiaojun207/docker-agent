package conf

import (
	"docker-agent/service/model"
	"github.com/docker/docker/api/types"
	"time"
)

var DockerInfo types.Info

var (
	LogFollowMap model.SyncMap
	ExecMap      model.SyncMap
)

var (
	TaskFrequency = time.Second * 60
)

func GetLogFollowStart(id string, Init func() model.LogFollow) model.LogFollow {
	m, _ := LogFollowMap.LoadInit(id, func() interface{} {
		return Init()
	})
	return m.(model.LogFollow)
}

func IsLogFollow(id string) bool {
	return LogFollowMap.ContainKey(id)
}

func RemoveLogFollow(id string) {
	LogFollowMap.LoadAndDelete(id)
}

///////////////////////
func GetDockerExecStart(id string, Init func() model.DockerExecStart) model.DockerExecStart {
	m, _ := ExecMap.LoadInit(id, func() interface{} {
		return Init()
	})
	return m.(model.DockerExecStart)
}

func RemoveDockerExecStart(id string) {
	m, has := ExecMap.LoadAndDelete(id)
	if has {
		close(m.(model.DockerExecStart).In)
	}
}
