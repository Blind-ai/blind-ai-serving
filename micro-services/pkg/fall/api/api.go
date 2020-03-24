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


/**

POST	/api/auth/fall/create/profile	username password firstname lastname email phonenumber country address city postalcode birthday sex amazon paypal
POST	/api/auth/fall/signin			username password
POST 	/api/auth/fall/signout
POST	/api/fall/create/review		contractid rating body
POST	/api/fall/edit/profile			[phonenumber country address city postalcode password wechat paypal]
GET		/api/fall/search/contracts		[page limit]
POST	/api/fall/take/contract		contractid
GET		/api/fall/search/lungh/contracts	username [page limit]
GET		/api/fall/search/lungh/profile	username
GET		/api/fall/search/lungh/reviews	username [page limit]
GET		/api/fall/search/self/profile
GET		/api/fall/search/self/reviews			[page limit]

 */

var Router *mux.Router

func init() {
	Router = mux.NewRouter()
}

func HandleRequests() {
	//TODO delete account

	Router.HandleFunc("/api/fall/evaluate/image", evaluateImage).Methods("POST")

	credentials := handlers.AllowCredentials()
	exposed := handlers.ExposedHeaders([]string{"X-Csrf-Token"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Csrf-Token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	log.Fatal(http.ListenAndServe("localhost:8001", handlers.CORS(exposed, headers, methods, origins, credentials)(Router)))
}