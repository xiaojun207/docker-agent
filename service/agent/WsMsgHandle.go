package agent

import (
	"docker-agent/service/conf"
	"docker-agent/utils"
	"log"
)

func MsgHandle(ch string, data map[string]interface{}) (error, map[string]interface{}) {
	switch ch {
	case "base.ht.ping":
		break
	case "base.ht.pong":
		break
	case "docker.image.pull":
		image := data["image"].(string)
		err := ImagePull(image)
		log.Println("ws: " + ch + " image:" + image + " is pull complate")
		return err, map[string]interface{}{"image": image}
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
		imageName := data["imageName"].(string)
		containerName := data["containerName"].(string)
		ports := utils.MapInterfaceToString(data["ports"].(map[string]interface{}))
		env := utils.ArrInterfaceToStr(data["env"].([]interface{}))
		volumes := utils.ArrInterfaceToStr(data["volumes"].([]interface{}))

		containerId, err := ContainerCreate(imageName, containerName, ports, env, volumes)
		log.Println("ws: docker.container.run taskId:" + taskId)
		PostContainers()
		return err, map[string]interface{}{"taskId": taskId, "containerId": containerId}
	case "docker.container.run":
		taskId := data["taskId"].(string)
		imageName := data["imageName"].(string)
		containerName := data["containerName"].(string)
		ports := utils.MapInterfaceToString(data["ports"].(map[string]interface{}))
		env := utils.ArrInterfaceToStr(data["env"].([]interface{}))
		volumes := utils.ArrInterfaceToStr(data["volumes"].([]interface{}))

		err, containerId := RunContainer(imageName, containerName, ports, env, volumes)
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
		conf.LogsFollow.Delete(containerId)
		PostContainersStats()
		return nil, map[string]interface{}{"containerId": containerId}
	case "docker.container.log.follow":
		containerId := data["containerId"].(string)
		isFollow, _ := conf.LogsFollow.LoadBool(containerId)
		if isFollow {
			log.Println("ws: " + ch + " containerId:" + containerId + ", log is follow")
			return nil, map[string]interface{}{"containerId": containerId}
		}

		conf.LogsFollow.Store(containerId, true)

		go ContainerLogFollow(containerId, func(timestamps int64, line string) bool {
			d := map[string]interface{}{
				"cId":  containerId,
				"ts":   timestamps,
				"line": line,
			}
			err := SendWsMsg("docker.container.log.line", d)
			if err != nil {
				conf.LogsFollow.Delete(containerId)
				return false
			}
			follow, _ := conf.LogsFollow.LoadBool(containerId)
			return follow
		})
		PostContainersStats()
		log.Println("ws: " + ch + " containerId:" + containerId)
		return nil, map[string]interface{}{"containerId": containerId}
	default:
		log.Println("unknown message "+ch, data)
		return nil, nil
	}
	return nil, nil
}
