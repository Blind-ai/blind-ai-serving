package util

import (
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Fatal(err error, level ...int) {
	if err != nil {
		if len(level) == 0 || level[0] == 0  {
			panic(err)
		} else if level[0] == 1 {
			log.Error(err)
		} else if level[0] > 1 {
			log.Info(err)
		}
	}
}

func Encode(data string) string {
	return  base64.StdEncoding.EncodeToString([]byte(data))
}

func Decode(id string) (string, error) {
	dec, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}
	return string(dec), err
}

func AccessControlAllow(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")
	return w
}