package utils

import (
	"bytes"
	"docker-agent/service/conf"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func PostData(uri string, v interface{}) (error, map[string]interface{}) {
	bytesData, _ := json.Marshal(v)
	//log.Println(string(bytesData))
	reader := bytes.NewReader(bytesData)
	return request("POST", uri, reader)
}

func GetData(uri string) (error, map[string]interface{}) {
	return request("GET", uri, nil)
}

func request(method, uri string, body io.Reader) (error, map[string]interface{}) {
	url := conf.DockerServer + uri
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(err.Error())
		return err, map[string]interface{}{}
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("authorization", conf.Token)
	request.Header.Set("AppId", conf.AppId)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
		return err, map[string]interface{}{}
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return err, map[string]interface{}{}
	}
	//byte数组直接转成string，优化内存
	res := map[string]interface{}{}
	json.Unmarshal([]byte(respBytes), &res)

	return nil, res
}
