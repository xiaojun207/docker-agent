package utils

import (
	"bytes"
	"docker-agent/service/conf"
	"encoding/json"
	"errors"
	"github.com/xiaojun207/go-base-utils/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const RSAPublicKey = "-----BEGIN RSA Public Key-----\n" +
	"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv1sWwK9qQDett8Sw4NL4\n" +
	"22RBT4MZiPY9Y/BCdl8KOgWQ+lDcnPVgn355HdZE+h6ADduEtVn2C2N6uD5nLOE3\n" +
	"eDpsUPKny5djLjYfsFJxQMTKG4vU9AMFItsHK45r3xA1ByMjsA9Ti7Gwzgd9w+tw\n" +
	"+ZE67JGNNtr00aIvEMyZlc7grG/4mT51AovJTt+ThBS9IWzemhSVDP8dHEI+XHK5\n" +
	"pL88MXPTkEXO1WheTlL7iDPXvo4/50a2osr17EFoGzC/aEG3PaBdPVNEFl7izH7J\n" +
	"F6T2/V1QnwpgKHojnCReB7A+vP3VCqZaMTXoBROnVrRuKsm1AS2KBrygwg1GvZwH\n" +
	"FwIDAQAB\n" +
	"-----END RSA Public Key-----\n"

func Login() {
	if conf.Token == "" {
		err, resp := PostData("/login", map[string]string{
			"username": conf.Username,
			"password": utils.RSAEncrypt(conf.Password, RSAPublicKey),
		})
		if err != nil {
			log.Println("Login Server err:", err)
			log.Panic(errors.New("login Server Fail"))
		}
		if resp["code"] == "100200" {
			log.Println("Login Server Success")
			conf.Token = resp["data"].(string)
		} else {
			log.Println("Login Server Fail:", resp)
			log.Panic(errors.New("login Server Fail"))
		}
	}
}

var cacheMap = map[string]string{}

func PostData(uri string, v interface{}) (error, map[string]interface{}) {
	bytesData, _ := json.Marshal(v)
	//log.Println(string(bytesData))
	d, has := cacheMap[uri]
	text := string(bytesData)
	if has && d == text {
		log.Println("数据已经提交过，避免重复提交, uri:", uri)
		return nil, nil
	}
	reader := bytes.NewReader(bytesData)
	err, resp := request("POST", uri, reader)
	cacheMap[uri] = text
	return err, resp
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
