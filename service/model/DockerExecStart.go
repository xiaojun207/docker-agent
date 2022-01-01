package model

type DockerExecStart struct {
	Id  string
	In  chan []byte
	Out func(d []byte) error
}
