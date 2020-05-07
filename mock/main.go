package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type (
	Response struct {
		Probability float32
		ProcessingTime string
	}
)

func receiveVideo(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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

func evaluateVideo(w http.ResponseWriter, r *http.Request) {
	var response Response

	t := time.Now()
	response.Probability = (rand.Float32() * 100)

	//process image
	if _, err := receiveVideo(w, r) ; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}

	response.ProcessingTime = time.Since(t).String()

	if toSend, err := json.Marshal(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Write(toSend)
	}
}


func receiveImage(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		return nil, err
	}

	// open the file contained in the html request
	file, handler, err := r.FormFile("file")
	if err != nil {
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
	var response Response

	t := time.Now()
	response.Probability = (rand.Float32() * 100)

	//process image
	if _, err := receiveImage(w, r) ; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) ; return
	}

	response.ProcessingTime = time.Since(t).String()

	if toSend, err := json.Marshal(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) ; return
	} else {
		w.Write(toSend)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/fall/evaluate/video", evaluateVideo).Methods("POST")
	router.HandleFunc("/api/lungh/evaluate/image", evaluateImage).Methods("POST")
	router.HandleFunc("/api/skin/evaluate/image", evaluateImage).Methods("POST")

	credentials := handlers.AllowCredentials()
	exposed := handlers.ExposedHeaders([]string{"X-Csrf-Token"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Csrf-Token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Println("Listen and serve...")
	log.Fatal(http.ListenAndServe(":8001", handlers.CORS(exposed, headers, methods, origins, credentials)(router)))
}