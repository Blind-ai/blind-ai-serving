package cfg

import (
	"encoding/json"
	"io/ioutil"
)

type RSA struct {
	PathPrivate		string `json:"PathPrivate"`
	PathPublic		string `json:"PathPublic"`
}

type Conf struct {
	ContractorRSA		RSA		`json:"ContractorRSA"`
	SignatoryRSA		RSA		`json:"SignatoryRSA"`
}


var CFG Conf

func init() {
	CFG = Load()
}


func Load() Conf {
	var conf Conf

	data, err := ioutil.ReadFile("/app/cfg/cfg.json")
	util.Fatal(err)
	err = json.Unmarshal(data, &conf)
	util.Fatal(err)
	return conf
}