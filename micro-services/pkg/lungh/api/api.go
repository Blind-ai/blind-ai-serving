package api

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ToAssign(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/myjwt")
}

var Router *mux.Router

func init() {
	Router = mux.NewRouter()
}

func HandleRequests() {
	Router.HandleFunc("/api/fall/evaluate/image", evaluateImage).Methods("POST")


	credentials := handlers.AllowCredentials()
	exposed := handlers.ExposedHeaders([]string{"X-Csrf-Token"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Csrf-Token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe("localhost:8001", handlers.CORS(exposed, headers, methods, origins, credentials)(Router)))
}