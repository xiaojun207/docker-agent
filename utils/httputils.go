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
	"strconv"
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
	log.Println("login...")
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

func PostData(uri string, v interface{}) (error, map[string]interface{}) {
	bytesData, _ := json.Marshal(v)
	//log.Println(string(bytesData))
	reader := bytes.NewReader(bytesData)
	err, resp := request("POST", uri, reader)
	return err, resp
}

func GetData(uri string) (error, map[string]interface{}) {
	return request("GET", uri, nil)
}

func request(method, uri string, body io.Reader) (error, map[string]interface{}) {
	url := conf.DockerServer + uri
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println("request.New.err:", err)
		return err, map[string]interface{}{}
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("authorization", conf.Token)
	req.Header.Set("AppId", conf.AppId)
	req.Header.Set("HostIp", conf.HostIp)
	req.Header.Set("PrivateIp", conf.PrivateIp)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("request.Do.err:", err)
		return err, map[string]interface{}{}
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("request.readAll.err:", err)
		return err, map[string]interface{}{}
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status), map[string]interface{}{
			"code": strconv.FormatInt(int64(int(resp.StatusCode)), 10),
			"msg":  string(respBytes),
		}
	}

	//byte数组直接转成string，优化内存
	res := map[string]interface{}{}
	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		log.Println("request.Unmarshal.err:", err)
	}

	return nil, res
}
