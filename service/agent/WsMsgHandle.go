package agent

import "log"

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
	case "docker.container.restart":
		containerId := data["containerId"].(string)
		err := ContainerRestart(containerId)
		log.Println("ws: " + ch + " containerId:" + containerId)
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
		id := data["id"].(string)
		imageName := data["imageName"].(string)
		containerName := data["containerName"].(string)
		hostPort := data["hostPort"].(string)
		appPort := data["appPort"].(string)
		containerId, err := ContainerCreate(imageName, containerName, hostPort, appPort)
		log.Println("ws: docker.container.run id:" + id)
		PostContainers()
		return err, map[string]interface{}{"id": id, "containerId": containerId}
	case "docker.container.run":
		id := data["id"].(string)
		imageName := data["imageName"].(string)
		containerName := data["containerName"].(string)
		hostPort := data["hostPort"].(string)
		appPort := data["appPort"].(string)
		err := RunContainer(imageName, containerName, hostPort, appPort)
		log.Println("ws: docker.container.run id:" + id)
		PostContainers()
		return err, map[string]interface{}{"id": id}
	case "docker.container.stats":
		containerId := data["containerId"].(string)
		stats, err := ContainerStats(containerId)
		log.Println("ws: " + ch + " containerId:" + containerId)
		return err, map[string]interface{}{"containerId": containerId, "stats": stats}
	default:
		log.Println("unknown message "+ch, data)
		return nil, nil
	}
	return nil, nil
}
