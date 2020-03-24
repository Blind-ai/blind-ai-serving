package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/blind-ai-serving/pkg/fall/api"
)

// read the key files before starting http handlers
func init() {
	log.Info("init finished")
}

func main() {
	api.HandleRequests()
}