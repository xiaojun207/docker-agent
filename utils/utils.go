package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-basic/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"unsafe"
)

var DockerServer = "http://localhost:8080"
var Token = "12345678"
var AppId = uuid.New()

func PostData(uri string, v interface{}) {
	url := DockerServer + uri
	bytesData, err := json.Marshal(v)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//log.Println(string(bytesData))
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("authorization", Token)
	request.Header.Set("AppId", AppId)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	log.Println(*str)
}

func GetData(uri string) {
	url := DockerServer + uri
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("authorization", Token)
	request.Header.Set("AppId", AppId)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	log.Println(*str)
}
