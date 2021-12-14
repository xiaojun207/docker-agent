package main

import (
	"docker-agent/service/agent"
	"docker-agent/service/conf"
	"docker-agent/utils"
	"flag"
	"log"
	"time"
)

//https://docs.docker.com/engine/api/sdk/examples/

func main() {
	flag.StringVar(&conf.DockerServer, "DockerServer", "", "DockerServer服务地址，如：http://127.0.0.1:8068/dockerMgrApi/agent")
	flag.StringVar(&conf.Username, "Username", "agent", "The username for dockerserver auth, default : agent")
	flag.StringVar(&conf.Password, "Password", "", "The password for dockerserver auth. You can get the token from DockerServer first start console logs")
	flag.Parse()
	flag.Usage()

	log.Println("Start docker agent, AppId:", conf.AppId)
	log.Println("conf.DockerServer:", conf.DockerServer)
	log.Println("conf.Username:", conf.Username)
	log.Println("conf.Password:", utils.SubStr(conf.Password, 0, 10)+"*****")

	if conf.DockerServer == "" {
		log.Panic("DockerServer must be set. like: -DockerServer http://127.0.0.1:8068/dockerMgrApi/agent")
	}

	log.Println("Start connect to DockerWsServer :", conf.GetDockerWsUrl())

	utils.Login()

	agent.StartWs()
	agent.GetAgentConfig()
	log.Println("PostDockerInfo")
	agent.PostDockerInfo()
	for true {
		go work()
		log.Println("time.Sleep", conf.TaskFrequency)
		time.Sleep(conf.TaskFrequency)
	}
}

func work() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work.err:", err)
		}
	}()
	log.Println("-----------------work----------------------------------")
	log.Println("PostContainers")
	agent.PostContainers()
	log.Println("PostContainersStats")
	agent.PostContainersStats()
	log.Println("PostImageList")
	agent.PostImageList()
}
