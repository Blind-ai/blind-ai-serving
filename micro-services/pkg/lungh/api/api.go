package api

import (
	commonApi "github.com/blind-ai-serving/pkg/common/api"
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
	Router.HandleFunc("/api/lungh/run", commonApi.RunTFS).Methods("GET")
	Router.HandleFunc("/api/lungh/remove", commonApi.RemoveTFS).Methods("GET")
	Router.HandleFunc("/api/lungh/start", commonApi.StartTFS).Methods("GET")
	Router.HandleFunc("/api/lungh/stop", commonApi.StopTFS).Methods("GET")

	Router.HandleFunc("/api/lungh/evaluate/image", evaluateImage).Methods("POST")

	credentials := handlers.AllowCredentials()
	exposed := handlers.ExposedHeaders([]string{"X-Csrf-Token"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Csrf-Token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8001", handlers.CORS(exposed, headers, methods, origins, credentials)(Router)))
}