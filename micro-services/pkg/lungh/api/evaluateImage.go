package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func receiveImage(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		return nil, err
	}

	// open the file contained in the html request
	file, handler, err := r.FormFile("file")
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
type (
	TensorFlowResp struct {
		instances map[string]string
	}
	Response struct {
		Type			string
		Probability   	float32
		ProcessingTime	string
	}
)

func evaluateImage(w http.ResponseWriter, r *http.Request) {
	log.Println("received request")
	var bodyBytes []byte
	var imgBytes []byte
	var err error
	var endpoint = "http://localhost:8501/api/lungh/compute"//var endpoint = os.Getenv("RESNET_ENDPOINT")
	fmt.Println("envpoint:", endpoint)
	var response Response
	var t = time.Now()

	if imgBytes, err = receiveImage(w, r) ; err != nil {
		w.WriteHeader(http.StatusBadRequest) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	resp, err := http.Post(endpoint, "image", bytes.NewReader(imgBytes)) ; if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadRequest);fmt.Fprintln(w, fmt.Sprintf("%v", err));return
	}
	if bodyBytes, err = ioutil.ReadAll(resp.Body) ; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) ; return
	}
	fmt.Println(string(bodyBytes))
	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		fmt.Println("unmarshallError")
		http.Error(w, err.Error(), http.StatusInternalServerError) ; return
	}

	response.ProcessingTime = time.Since(t).String()
	if toSend, err := json.Marshal(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(toSend)
	}
}