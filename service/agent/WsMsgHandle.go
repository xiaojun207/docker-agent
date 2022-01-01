package agent

import (
	"docker-agent/service/conf"
	"docker-agent/service/dto"
	"docker-agent/service/model"
	utils2 "github.com/xiaojun207/go-base-utils/utils"
	"log"
)

func MsgHandle(ch string, data map[string]interface{}) (error, map[string]interface{}) {
	switch ch {
	case "base.ht.ping":
		break
	case "base.ht.pong":
		break
	case "base.config.update":
		GetAgentConfig()
		return nil, map[string]interface{}{}
	case "docker.image.pull":
		image := data["image"].(string)
		err := ImagePull(image)
		log.Println("ws: " + ch + " image:" + image + " is pull complate")
		PostImageList()
		return err, map[string]interface{}{"image": image}
	case "docker.image.remove":
		imageId := data["imageId"].(string)
		err := ImageRemove(imageId)
		log.Println("ws: "+ch+" imageId:"+imageId+" is remove complate, err:", err)
		PostImageList()
		return err, map[string]interface{}{"imageId": imageId}
	case "docker.image.prune":
		err := ImagePrune()
		PostImageList()
		return err, map[string]interface{}{}
	case "docker.system.prune":
		err := SystemPrune()
		PostImageList()
		PostContainers()
		return err, map[string]interface{}{}
	case "docker.container.list":
		PostContainers()
		return nil, map[string]interface{}{}
	case "docker.container.restart":
		containerId := data["containerId"].(string)
		err := ContainerRestart(containerId)
		log.Println("ws: "+ch+" containerId:"+containerId+", err:", err)
		PostContainers()
		return err, map[string]interface{}{"containerId": containerId}
	case "docker.container.stop":
		containerId := data["containerId"].(string)
		err := ContainerStop(containerId)
		log.Println("ws: " + ch + " containerId:" + containerId)
		PostContainers()
		return err, map[string]interface{}{"containerId": containerId}
	case "docker.container.remove":
		containerId := data["containerId"].(string)
		err := ContainerRemove(containerId)
		log.Println("ws: " + ch + " containerId:" + containerId)
		PostContainers()
		return err, map[string]interface{}{"containerId": containerId}
	case "docker.container.start":
		containerId := data["containerId"].(string)
		err := ContainerStart(containerId)
		log.Println("ws: " + ch + " containerId:" + containerId)
		PostContainers()
		return err, map[string]interface{}{"containerId": containerId}
	case "docker.container.create":
		taskId := data["taskId"].(string)
		conf := dto.ContainerCreateConfig{}
		utils2.MapToStruct(data, &conf)
		containerId, err := ContainerCreate(conf)
		log.Println("ws: docker.container.run taskId:" + taskId)
		PostContainers()
		return err, map[string]interface{}{"taskId": taskId, "containerId": containerId}
	case "docker.container.run":
		taskId := data["taskId"].(string)
		conf := dto.ContainerCreateConfig{}
		utils2.MapToStruct(data, &conf)
		err, containerId := RunContainer(conf)
		log.Println("ws: docker.container.run taskId:" + taskId)
		PostContainers()
		return err, map[string]interface{}{"containerId": containerId}
	case "docker.container.stats":
		containerId := data["containerId"].(string)
		stats, err := ContainerStats(containerId)
		log.Println("ws: " + ch + " containerId:" + containerId)
		return err, map[string]interface{}{"containerId": containerId, "stats": stats}
	case "docker.containers.stats":
		PostContainersStats()
		return nil, map[string]interface{}{}
	case "docker.container.logs":
		containerId := data["containerId"].(string)
		tail := data["tail"].(string)
		since := data["since"].(string)
		logs, err := ContainerLogs(containerId, tail, since)
		log.Println("ws: " + ch + " containerId:" + containerId)
		return err, map[string]interface{}{"containerId": containerId, "logs": logs}
	case "docker.container.log.follow.close":
		containerId := data["containerId"].(string)
		conf.RemoveLogFollow(containerId)
		PostContainersStats()
		return nil, map[string]interface{}{"containerId": containerId}
	case "docker.container.log.follow":
		containerId := data["containerId"].(string)
		logFollow(containerId)
		PostContainersStats()
		log.Println("ws: " + ch + " containerId:" + containerId)
		return nil, map[string]interface{}{"containerId": containerId}
	case "docker.container.exec":
		containerId := data["cId"].(string)
		cmd := ""
		if data["cmd"] != nil {
			cmd = data["cmd"].(string)
		}
		d := ([]byte)(data["d"].(string))
		dockerExec(containerId, cmd, d)
		break
	case "docker.container.exec.close":
		containerId := data["cId"].(string)
		conf.RemoveDockerExecStart(containerId)
		break
	default:
		log.Println("unknown message "+ch, data)
		return nil, nil
	}
	return nil, nil
}

func logFollow(containerId string) {
	Init := func() model.LogFollow {
		out := func(timestamps int64, line string) bool {
			d := map[string]interface{}{
				"cId":  containerId,
				"ts":   timestamps,
				"line": line,
			}
			err := SendWsMsg("docker.container.log.line", d)
			if err != nil {
				conf.RemoveLogFollow(containerId)
				return false
			}
			return conf.IsLogFollow(containerId)
		}
		go ContainerLogFollow(containerId, out)
		return model.LogFollow{
			Id:  containerId,
			Out: out,
		}
	}
	conf.GetLogFollowStart(containerId, Init)
}

func dockerExec(containerId, cmd string, d []byte) {
	Init := func() model.DockerExecStart {
		in := make(chan []byte)
		out := func(d []byte) error {
			msg := map[string]interface{}{
				"cId": containerId,
				"d":   string(d),
			}
			err := SendWsMsg("docker.container.exec", msg)
			if err != nil {
				conf.RemoveDockerExecStart(containerId)
			}
			return err
		}

		go ContainerExec(containerId, cmd, out, in)
		return model.DockerExecStart{
			Id:  containerId,
			In:  in,
			Out: out,
		}
	}

	execModel := conf.GetDockerExecStart(containerId, Init)
	if d != nil && len(d) > 0 {
		execModel.In <- d
	}
}
