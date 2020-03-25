package main

import (
	"github.com/blind-ai-serving/pkg/skin/api"
	log "github.com/sirupsen/logrus"
)

// read the key files before starting http handlers
func init() {
	log.Info("init finished")
}

func main() {
	api.HandleRequests()
}