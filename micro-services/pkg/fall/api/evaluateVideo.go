package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const SERVER_URL = "http://localhost:8501/v1/models/resnet:predict"

func receiveVideo(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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
	fmt.Println(content[0])
	if content[0] != "video/mp4" {
		return nil, errors.New("file sent is not a video")
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
		instances map[string][]map[string]interface{}

	}
)

func evaluateVideo(w http.ResponseWriter, r *http.Request) {
	log.Println("received request")
	var mp TensorFlowResp
	imgBytes, err := receiveVideo(w, r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}

	b64Img := base64.StdEncoding.EncodeToString(imgBytes)
	request := fmt.Sprintf("{\"instances\" : [{\"b64\": \"%s\"}]}", b64Img)


	TFSRes, err := http.Post(SERVER_URL, "application/json", bytes.NewBuffer([]byte(request))) ; if err != nil {
		w.WriteHeader(http.StatusBadRequest) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	//we dont know the exact content of the json so we create a map
	body, err := ioutil.ReadAll(TFSRes.Body)
	err = json.Unmarshal(body, &mp.instances)
	if err != nil {
		fmt.Println("ai eroor")
		w.WriteHeader(http.StatusBadRequest) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}

	// we need to resolve the right types and then get our value
	probabilities := mp.instances["predictions"][0]["probabilities"].([]interface{})
	class := int(mp.instances["predictions"][0]["classes"].(float64))
	probability := probabilities[class].(float64)

	// return the result to the user
	w.Write([]byte(fmt.Sprintf("class: %d,  probability: %.8f", class, probability)))
}