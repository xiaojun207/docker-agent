package agent

import (
	"github.com/docker/docker/api/types"
	"io"
	"log"
)

func exec(container string, workdir, cmd string) (hr types.HijackedResponse, err error) {
	// 执行/bin/bash命令
	ir, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   workdir,
		Cmd:          []string{cmd},
		Tty:          true,
	})
	if err != nil {
		return
	}

	// 附加到上面创建的/bin/bash进程中
	hr, err = cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return
	}
	return
}

func ContainerExec(containerId, cmd string, out func(d []byte) error, inChan chan []byte) error {
	log.Println("ContainerExec.containerId:", containerId, ",cmd:", cmd)
	// 执行exec，获取到容器终端的连接
	hr, err := exec(containerId, "/", cmd)
	if err != nil {
		log.Println(err)
		return err
	}
	// 关闭I/O流
	defer hr.Close()
	// 退出进程
	defer func() {
		hr.Conn.Write([]byte("exit\r"))
	}()

	// copy from docker to ws
	go wsWriterCopy(hr.Conn, out)
	// copy from ws to docker
	wsReaderCopy(inChan, hr.Conn)
	return err
}

// copy from docker to ws
func wsWriterCopy(reader io.Reader, outFunc func(d []byte) error) {
	buf := make([]byte, 8192)
	for {
		nr, err := reader.Read(buf)
		if nr > 0 {
			if outFunc(buf[0:nr]) != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

// copy from ws to docker
func wsReaderCopy(inChan chan []byte, writer io.Writer) {
	for {
		p, ok := <-inChan
		if !ok {
			log.Println("wsReaderCopy.break")
			break
		}
		writer.Write(p)
	}
}
