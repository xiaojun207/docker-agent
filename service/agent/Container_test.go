package agent

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestContainerLogs(t *testing.T) {
	log.Println(time.Now().String())
	containerId := "302c6c35f81699fd8d3c02e903cea24b4979b51080abd8dcbbb06bf3de45016b"
	logs, err := ContainerLogs(containerId, "100", "2021-05-30T04:07:11Z")
	log.Println(err)
	fmt.Println(logs)
}
