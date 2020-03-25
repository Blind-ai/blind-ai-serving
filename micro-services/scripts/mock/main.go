package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"github.com/google/uuid"
)



func GetDoctors(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	fmt.Fprintf(w, "ID:0;FirstName:Cabaillero;LastName:JeanSass;Age:21;Adress:xxxxxxxxxxx")
}

func GetPatients(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	fmt.Fprintf(w, "ID:0;FirstName:Louis;LastName:Vuitton;Age:21;Adress:xxxxxxxxxxx")

}

func fallHandler(filename string) (int, error) {
	return rand.Int() % 100, nil
}

func lunghHandler(filename string) (int, error) {
	return rand.Int() % 100, nil
}

func skinHandler(filename string) (int, error) {
	return rand.Int() % 100, nil
}

func checkFile(r *http.Request) (bool, string, string) {
	fType := r.FormValue("type")

	// verify that the file type is linked to one of our three AI
	if fType != "skin" && fType != "lungh" && fType != "skin" {
		return false, "field \"type\" is wrong, should be : \"skin\", \"lungh\" or \"skin\"", ""
	}
	// open the file contained in the html request
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println(err)
		return false, err.Error(), ""
	} ; defer file.Close()
	//verify that the file is either an image or a video
	content := handler.Header["Content-Type"]
	if content[0][:5] != "image" && content[0][:5] != "video" {
		return false, "file sent is neither an image or a video", ""
	}
	// load the file bytes
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return false, err.Error(), ""
	}
	// put our file in a tmp folder with a uuid name + file extension
	ioutil.WriteFile("tmp/" + uuid.New().String() + handler.Filename[len(handler.Filename)-4:], fileBytes, 0644)
	log.Info("Successfully Uploaded File\n")
	return true, content[0][:5], handler.Filename
}

func UploadFile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")

	fmt.Println(r.Body)
	fmt.Println(r)

	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(20 << 20)

	// Create a temporary file within our temp-images directory that follows
	valid, info, filename := checkFile(r)
	if valid != true {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, info)
		return
	}
	/**
	In this part we call the different AI handlers depending on the type of the file
	*/
	var err error
	var result int
	if info == "image" {
		if r.FormValue("type") == "lungh" {
			result, err = lunghHandler(filename)
		} else {
			result, err = skinHandler(filename)
		}
	} else {
		result, err = fallHandler(filename)
	}
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, err.Error())
		return
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, strconv.Itoa(result))
		return
	}
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/api/mock/uploadfile",UploadFile)
	router.HandleFunc("/api/mock/getdoctors",GetDoctors)
	router.HandleFunc("/api/mock/getpatients",GetPatients)
	fmt.Println("listen and serve ...")
	log.Fatal(http.ListenAndServe(":8001", router))

}