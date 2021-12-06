package utils

import (
	"github.com/ricochet2200/go-disk-usage/du"
	"log"
	"testing"
)

func TestPrintDisk(t *testing.T) {
	usage := du.NewDiskUsage("/")
	log.Println("Usage:", usage.Usage()/1024.0/1024.0/1024.0)
	log.Println("Free:", usage.Free()/1024.0/1024.0/1024.0)
	log.Println("Size:", usage.Size()/1024.0/1024.0/1024.0)
	log.Println("Used:", usage.Used()/1024.0/1024.0/1024.0)
	log.Println("Available:", usage.Available()/1024.0/1024.0/1024.0)
}
