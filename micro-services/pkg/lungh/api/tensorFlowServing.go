package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func execTFS(cmd *exec.Cmd) {
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run() ; err != nil {
		log.Println(err)
	}
	fmt.Println(out.String())
}

func runTFS(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("docker",
		"run",
			"-p",
			"8501:8501",
			"--name",
			"tfserving_resnet",
			"--mount",
			"type=bind,source=/tmp/resnet,target=/models/resnet",
			"-e",
			"MODEL_NAME=resnet",
			"-t",
			"tensorflow/serving",
			"--rm")

	_, err := exec.LookPath("docker")
	if err != nil {
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	go execTFS(cmd)
	w.WriteHeader(http.StatusOK)
}

func removeTFS(w http.ResponseWriter, r *http.Request) {
	var out bytes.Buffer
	cmd := exec.Command("docker",
		"rm",
		"tfserving_resnet")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	fmt.Println(out.String())
	w.WriteHeader(http.StatusOK)
	w.Write(out.Bytes())
}


func startTFS(w http.ResponseWriter, r *http.Request) {
	var out bytes.Buffer
	cmd := exec.Command("docker",
		"start", "tfserving_resnet")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(out.Bytes())
}

func stopTFS(w http.ResponseWriter, r *http.Request) {
	var out bytes.Buffer
	cmd := exec.Command("docker",
		"kill", "tfserving_resnet")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(out.Bytes())
}