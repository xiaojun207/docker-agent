package conf

import (
	"docker-agent/service/model"
	"github.com/docker/docker/api/types"
)

var DockerInfo types.Info

var (
	LogsFollow model.SyncMap
)
