package agent

import (
	"docker-agent/service/conf"
	"docker-agent/utils"
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type DockerAgentWs struct {
}

var wsConn *utils.WsConn

func StartWs() {
	endpoint := conf.GetDockerWsUrl()
	wsConn = utils.NewWsBuilder().
		WsUrl(endpoint).
		AutoReconnect().
		Dump().
		ProtoHandleFunc(wsMsgHandle).
		ErrorHandleFunc(ErrorHandleFunc).
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

func SendWsMsg(ch string, data interface{}) error {
	msg := map[string]interface{}{
		"ch": ch,
		"ts": time.Now().UnixNano() / 1e6,
		"d":  data,
	}
	return wsConn.SendJsonMessage(msg)
}

func wsMsgHandle(msg []byte) error {
	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Println("wsMsgHandle.err:", err)
	//	}
	//}()
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

	if datamap["d"] != nil && reflect.TypeOf(data) == reflect.TypeOf(datamap["d"]) {
		data = datamap["d"].(map[string]interface{})
	}
	err, resp := MsgHandle(ch, data)

	respCh := "base.ack"
	if strings.HasPrefix(ch, "docker.") {
		respCh = "docker.task.ack"
	}

	taskId, has := data["taskId"]
	code := "100200"
	errMsg := ""
	if err != nil {
		code = "100100"
		errMsg = err.Error()
	}
	if has {
		SendWsMsg(respCh, map[string]interface{}{"code": code, "msg": errMsg, "taskId": taskId, "resp": resp})
	}

	return err
}

func ErrorHandleFunc(err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("ErrorHandleFunc.err:", err)
		}
	}()
	if err != nil {
		text := err.Error()
		log.Println("ErrorHandleFunc:", text)
		r := map[string]interface{}{}
		json.Unmarshal([]byte(text), &r)
		if r["code"] == "105101" {
			utils.Login()
		}
	}
}
