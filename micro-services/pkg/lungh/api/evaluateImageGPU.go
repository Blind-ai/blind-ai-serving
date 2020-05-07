package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func EvaluateImageGpu(w http.ResponseWriter, r *http.Request) {
	log.Println("received request")
	const endpoint_gpu = "http://localhost:8501/v1/models/resnet:predict"
	var mp TensorFlowResp
	imgBytes, err := receiveImage(w, r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}

	b64Img := base64.StdEncoding.EncodeToString(imgBytes)
	request := fmt.Sprintf("{\"instances\" : [{\"b64\": \"%s\"}]}", b64Img)

	t := time.Now()
	TFSRes, err := http.Post(endpoint_gpu, "application/json", bytes.NewBuffer([]byte(request))) ; if err != nil {
		w.WriteHeader(http.StatusBadRequest) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	result := time.Since(t)
	//we dont know the exact content of the json so we create a map
	body, err := ioutil.ReadAll(TFSRes.Body)
	err = json.Unmarshal(body, &mp.instances)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}

	// we need to resolve the right types and then get our value
	probabilities := mp.instances["predictions"][0]["probabilities"].([]interface{})
	class := int(mp.instances["predictions"][0]["classes"].(float64))
	probability := probabilities[class].(float64)

	// return the result to the user
	w.Write([]byte(fmt.Sprintf("class: %d,  probability: %.8f, processing time: %s", class, probability, result.String())))
}