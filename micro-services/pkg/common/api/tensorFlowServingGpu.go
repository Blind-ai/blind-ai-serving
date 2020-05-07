package api

import (
	"fmt"
	"net/http"
	"os/exec"
)

func RunTFSGpu(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("docker",
		"run",
		"--publish",
		r.FormValue("publish"),
		"--name",
		r.FormValue("name"),
		"--mount",
		r.FormValue("mount"),
		"--env",
		r.FormValue("env"),
		"--tty",
		r.FormValue("tty"),
		"--runtime=nvidia")
	_, err := exec.LookPath("docker")
	if err != nil {
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	go execTFS(cmd)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running " + r.FormValue("name")))
}