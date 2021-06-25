package conf

import (
	"docker-agent/service/model"
	"github.com/docker/docker/api/types"
	"time"
)

var DockerInfo types.Info

var (
	LogsFollow model.SyncMap
)

var (
	TaskFrequency = time.Second * 60
)
