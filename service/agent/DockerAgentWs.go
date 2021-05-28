package agent

import (
	"docker-agent/service/conf"
	"docker-agent/utils"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type DockerAgentWs struct {
}

var wsConn *utils.WsConn

func StartWs() {
	endpoint := conf.DockerWsServer
	wsConn = utils.NewWsBuilder().
		WsUrl(endpoint).
		AutoReconnect().
		ProtoHandleFunc(WsMsgHandle).
		ReconnectInterval(time.Millisecond * 5).
		Build()
	go exitHandler(wsConn)
}

func exitHandler(c *utils.WsConn) {
	pingTicker := time.NewTicker(10 * time.Minute)
	pongTicker := time.NewTicker(time.Second)
	defer pingTicker.Stop()
	defer pongTicker.Stop()
	defer c.CloseWs()

	for {
		select {
		case <-c.CloseChan():
			return
		case t := <-pingTicker.C:
			c.SendPingMessage([]byte(strconv.Itoa(int(t.UnixNano() / int64(time.Millisecond)))))
		case t := <-pongTicker.C:
			c.SendPongMessage([]byte(strconv.Itoa(int(t.UnixNano() / int64(time.Millisecond)))))
		}
	}
}

func SendWsMsg(ch string, data interface{}) {
	msg := map[string]interface{}{
		"ch": ch,
		"ts": time.Now().UnixNano() / 1e6,
		"d":  data,
	}
	wsConn.SendJsonMessage(msg)
}
func WsMsgHandle(msg []byte) error {
	datamap := make(map[string]interface{})
	err := json.Unmarshal(msg, &datamap)
	if err != nil {
		log.Println("json unmarshal error for ", string(msg))
		return err
	}

	ch, isOk := datamap["ch"].(string)
	if !isOk {
		log.Println("no message ch, msg:" + string(msg))
		return err
	}

	log.Println("recv:" + string(msg))
	data := map[string]interface{}{}
	if datamap["d"] != nil {
		data = datamap["d"].(map[string]interface{})
	}
	err, resp := MsgHandle(ch, data)
	SendWsMsg(ch+".ack", map[string]interface{}{"err": err, "resp": resp})
	return err
}
func MsgHandle(ch string, data map[string]interface{}) (error, map[string]interface{}) {
	switch ch {
	case "base.ht.ping":
		break
	case "base.ht.pong":
		break
	case "docker.image.pull":
		image := data["image"].(string)
		err := PullImage(image)
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
