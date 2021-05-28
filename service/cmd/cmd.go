package cmd

import (
	"bytes"
	"github.com/docker/docker/pkg/reexec"
	"log"
)

func CmdRun(args ...string) {
	log.Println("CmdRun.args:", args)

	var outInfo bytes.Buffer
	cmd := reexec.Command(args...)
	log.Println("Run")
	//cmd.Stdout = nil
	//out, err := cmd.StdoutPipe()
	//if err != nil {
	//	log.Println("err1:", err)
	//	return
	//}
	// 设置接收
	cmd.Stdout = &outInfo
	err := cmd.Run()
	if err != nil {
		log.Println("err2:", err)
		return
	}

	log.Println("res:", outInfo.String())
}
