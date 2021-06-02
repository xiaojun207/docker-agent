package conf

import (
	"github.com/go-basic/uuid"
	"strings"
)

var DockerServer = "" //eg.: http://localhost:8080/docker/
var Token = "12345678"
var AppId = uuid.New()

func GetDockerWsUrl() string {
	//var DockerWsServer = "" //eg.: ws://localhost:8080/docker/ws
	if strings.HasPrefix(DockerServer, "http://") {
		return "ws://" + strings.TrimLeft(DockerServer, "http://") + "/ws"
	} else {
		return "wss://" + strings.TrimLeft(DockerServer, "https://") + "/ws"
	}
}
