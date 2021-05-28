package conf

import "github.com/go-basic/uuid"

var DockerServer = ""   //eg.: http://localhost:8080/docker/
var DockerWsServer = "" //eg.: ws://localhost:8080/docker/ws
var Token = "12345678"
var AppId = uuid.New()
