package main

import (
	"docker-agent/service/agent"
	"docker-agent/service/conf"
	"flag"
	"log"
	"time"
)

//https://docs.docker.com/engine/api/sdk/examples/

func main() {
	flag.StringVar(&conf.DockerServer, "DockerServer", "http://127.0.0.1:8080/dockerApi", "DockerServer服务地址，如：http://127.0.0.1:8080/dockerApi")
	flag.StringVar(&conf.DockerWsServer, "DockerWsServer", "", "DockerWsServer服务地址，如：http://127.0.0.1:8068/dockerApi/ws")
	flag.StringVar(&conf.Token, "Token", "", "The http and websocket header authorization for dockerserver auth")
	flag.Parse()

	log.Println("Start docker agent, AppId:", conf.AppId)
	log.Println("conf.DockerServer:", conf.DockerServer)
	log.Println("conf.DockerWsServer:", conf.DockerWsServer)

	if conf.DockerServer == "" && conf.DockerWsServer == "" {
		log.Panic("DockerServer and DockerWsServer, must one of not empty\nlike: -DockerServer http://127.0.0.1:8080/dockerApi, or  -DockerWsServer http://127.0.0.1:8080/dockerApi/ws ")
	}

	if conf.DockerWsServer != "" {
		log.Println("Start connect to DockerWsServer :", conf.DockerWsServer)
		agent.StartWs()
	} else {
		log.Println("DockerWsServer not set, that websocket client not start")
	}

	for true {
		go work()
		time.Sleep(time.Minute * 1)
	}
}

func work() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work.err:", err)
		}
	}()
	log.Println("-----------------work----------------------------------")
	agent.PostDockerInfo()
	//getTask(cli)
	agent.PostContainers()
	//StopAllRunningContainers(cli)

}
