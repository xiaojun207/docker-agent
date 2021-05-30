package main

import (
	"docker-agent/service/agent"
	"docker-agent/service/conf"
	"flag"
	"fmt"
	"log"
	"time"
)

//https://docs.docker.com/engine/api/sdk/examples/

func PrintParamInfo() {

	fmt.Println("Parameter description:")
	fmt.Println("\tDockerServer\trequired\thttp server address，like：http://127.0.0.1:8080/dockerApi")
	fmt.Println("\tDockerWsServer\trequired\twebsocket server address，like：http://127.0.0.1:8080/dockerApi/ws")
	fmt.Println("\tToken\t\tOptional\tThe http and websocket header authorization for dockerserver auth")
}

func main() {
	PrintParamInfo()

	flag.StringVar(&conf.DockerServer, "DockerServer", "http://127.0.0.1:8080/dockerApi", "DockerServer服务地址，如：http://127.0.0.1:8080/dockerApi")
	flag.StringVar(&conf.DockerWsServer, "DockerWsServer", "", "DockerWsServer服务地址，如：http://127.0.0.1:8068/dockerApi/ws")
	flag.StringVar(&conf.Token, "Token", "", "The http and websocket header authorization for dockerserver auth")
	flag.Parse()

	log.Println("Start docker agent, AppId:", conf.AppId)
	log.Println("conf.DockerServer:", conf.DockerServer)
	log.Println("conf.DockerWsServer:", conf.DockerWsServer)

	if conf.DockerServer == "" {
		log.Panic("DockerServer must be set. like: -DockerServer http://127.0.0.1:8080/dockerApi")
	}

	if conf.DockerWsServer != "" {
		log.Println("Start connect to DockerWsServer :", conf.DockerWsServer)
		agent.StartWs()
	} else {
		log.Println("DockerWsServer not set, that websocket client not start, if you want websocket, you can add param,like :  -DockerWsServer http://127.0.0.1:8080/dockerApi/ws ")
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
