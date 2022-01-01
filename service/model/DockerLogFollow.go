package model

type LogFollow struct {
	Id  string
	Out func(timestamps int64, line string) bool
}
