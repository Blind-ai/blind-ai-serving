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

func RunTFS(w http.ResponseWriter, r *http.Request) {
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
			r.FormValue("tty"))

	fmt.Println("docker",
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
		r.FormValue("tty"))
	_, err := exec.LookPath("docker")
	if err != nil {
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	go execTFS(cmd)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running " + r.FormValue("name")))
}

func RemoveTFS(w http.ResponseWriter, r *http.Request) {
	var out bytes.Buffer
	cmd := exec.Command("docker",
		"rm",
		r.FormValue("name"))
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	fmt.Println(out.String())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("removed " + r.FormValue("name")))
}


func StartTFS(w http.ResponseWriter, r *http.Request) {
	var out bytes.Buffer
	cmd := exec.Command("docker",
		"start",
		r.FormValue("name"))
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("started " + r.FormValue("name")))
}

func StopTFS(w http.ResponseWriter, r *http.Request) {
	var out bytes.Buffer
	cmd := exec.Command("docker",
		"kill",
		r.FormValue("name"))
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusForbidden) ; fmt.Fprintln(w, fmt.Sprintf("%v", err)) ; return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("stopped " + r.FormValue("name")))
}