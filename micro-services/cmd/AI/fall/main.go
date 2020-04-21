package main

import (
	"github.com/blind-ai-serving/pkg/fall/api"
	log "github.com/sirupsen/logrus"
)

// read the key files before starting http handlers
func init() {
	log.Info("init fall finished")
}

func main() {
	api.HandleRequests()
}