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
	flag.StringVar(&conf.DockerServer, "DockerServer", "127.0.0.1", "DockerServer服务地址")
	flag.StringVar(&conf.DockerWsServer, "DockerWsServer", "", "DockerServer服务地址")
	flag.StringVar(&conf.Token, "Token", "12345678", "token")
	flag.Parse()

	log.Println("Start docker agent, AppId:", conf.AppId)
	log.Println("utils.DockerServer:", conf.DockerServer)

	if conf.DockerServer == "" && conf.DockerWsServer == "" {
		log.Panic("DockerServer and DockerWsServer, must one of not empty")
	}

	agent.RegDocker()
	if conf.DockerWsServer != "" {
		log.Println("Start connect DockerWsServer")
		agent.StartWs()
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
	agent.RegDocker()
	//getTask(cli)
	agent.PostContainers()
	//StopAllRunningContainers(cli)

}
