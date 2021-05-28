package cmd

import "testing"

func TestCmd(t *testing.T) {
	args := []string{"docker", "stop", "test"}
	CmdRun(args...)
}
