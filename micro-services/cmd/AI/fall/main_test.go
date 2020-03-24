package main

import (
	"fmt"
	"testing"
)

func TestMyInit(t *testing.T) {
	handlers := MyInit()
	fmt.Println(handlers.RSA)
	handlers.Router.HandleFunc("/api/test", nil).Methods("POST")
}