package main

import (
	"context"
	"docker-agent/utils"
	"flag"
	"github.com/docker/docker/client"
	"log"
	"time"
)

//https://docs.docker.com/engine/api/sdk/examples/

func main() {
	flag.StringVar(&utils.DockerServer, "DockerServer", "127.0.0.1", "DockerServer服务地址")
	flag.StringVar(&utils.Token, "Token", "12345678", "token")
	flag.Parse()

	log.Println("Start docker agent, AppId:", utils.AppId)
	log.Println("utils.DockerServer:", utils.DockerServer)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	for true {
		go work(cli, ctx)
		time.Sleep(time.Minute * 1)
	}
}

func work(cli *client.Client, ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work:", err)
		}
	}()
	log.Println("-----------------work----------------------------------")
	utils.RegDocker(cli, ctx)
	//getTask(cli, ctx)
	utils.PostContainers(cli, ctx)
	//StopAllRunningContainers(cli, ctx)

}
