package cfg

import (
	"fmt"

	"testing")

func TestLoad(t *testing.T) {
	c := Load()
	fmt.Println(c.RSA)
	fmt.Println(c.AwsDB)
}