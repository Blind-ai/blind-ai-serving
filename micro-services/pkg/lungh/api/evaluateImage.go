package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

const SERVER_URL = "http://localhost:8501/v1/models/resnet:predict"

func receiveImage(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		return nil, err
	}

	// open the file contained in the html request
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println(err)
		return nil, err
	} ; defer file.Close()

	//verify that the file is either an image or a video
	content := handler.Header["Content-Type"]
	if content[0] != "image/png" && content[0] != "image/jpeg" {
		return nil, errors.New("file sent is neither an png or jpeg image")
	}
	// load the file bytes
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// put our file in a tmp folder with a uuid name + file extension
	//ioutil.WriteFile(uuid.New().String() + handler.Filename[len(handler.Filename)-4:], fileBytes, 0644)
	log.Info("Successfully Received File\n")
	return fileBytes, nil
}


func evaluateImage(w http.ResponseWriter, r *http.Request) {
	nb := rand.Int63n(100)
	s1 := strconv.FormatInt(int64(nb), 10)
	imgBytes, err := receiveImage(w, r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	encoded := base64.StdEncoding.EncodeToString(imgBytes)
	res, err := http.Post(SERVER_URL, "application/json", bytes.NewBuffer([]byte(encoded))) ; if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	//we dont know the exact content of the json so we create a map
	mp := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&mp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Println(mp)
	// Declare a new Person struct.
	fmt.Println(res.Header, res.Body)
	w.Write([]byte(s1))
}