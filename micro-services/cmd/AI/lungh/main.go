package main

import (
	"os"

	"github.com/blind-ai-serving/pkg/lungh/api"
	log "github.com/sirupsen/logrus"
)

// read the key files before starting http handlers
func init() {
	log.Info("init lungh finished")
}

func main() {
	endpoint := os.Getenv("RESNET_ENDPOINT")

	log.Info(endpoint)

	api.HandleRequests()
}
