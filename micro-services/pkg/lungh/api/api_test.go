package api

import (
	"bytes"
	"encoding/json"
	createC "github.com/jackline/pkg/api/lungh/create"
	delC "github.com/jackline/pkg/api/lungh/delete"
	editC "github.com/jackline/pkg/api/lungh/edit"
	searchCSelf "github.com/jackline/pkg/api/lungh/search/self"
	searchCSign "github.com/jackline/pkg/api/lungh/search/fall"
	createS "github.com/jackline/pkg/api/fall/create"
	editS "github.com/jackline/pkg/api/fall/edit"
	searchS "github.com/jackline/pkg/api/fall/search"
	searchSCont "github.com/jackline/pkg/api/fall/search/lungh"
	searchSSelf "github.com/jackline/pkg/api/fall/search/self"
	"github.com/jackline/pkg/api/signatory/take"
	"github.com/jackline/pkg/util"
	neturl "net/url"

	"github.com/gorilla/mux"
	"github.com/jackline/pkg/api/auth"
	"github.com/jackline/pkg/database"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func myInit() *mux.Router {
	// INIT Router
	xrouter := mux.NewRouter()
	// INIT Keys for contractors
	database.DeleteAllTables()
	database.CreateAllTables()

	//contractors routes
	xrouter.HandleFunc("/api/auth/lungh/signup", auth.CreateContractor).Methods("POST")
	xrouter.HandleFunc("/api/auth/lungh/signin", auth.SignInContractor).Methods("POST")

	xrouter.HandleFunc("/api/lungh/create/contract", createC.Contract).Methods("POST")
	xrouter.HandleFunc("/api/lungh/create/review", createC.Review).Methods("POST")


	xrouter.HandleFunc("/api/lungh/edit/profile", editC.Profile).Methods("POST")
	xrouter.HandleFunc("/api/lungh/edit/contract", editC.Contract).Methods("POST")

	xrouter.HandleFunc("/api/lungh/search/fall/profile", searchCSign.Profile).Methods("GET")
	xrouter.HandleFunc("/api/lungh/search/fall/reviews", searchCSign.Reviews).Methods("GET")

	xrouter.HandleFunc("/api/lungh/search/self/profile", searchCSelf.Profile).Methods("GET")
	xrouter.HandleFunc("/api/lungh/search/self/reviews", searchCSelf.Reviews).Methods("GET")

	xrouter.HandleFunc("/api/lungh/delete/contract", delC.Contract).Methods("POST")

	//signatories routes
	xrouter.HandleFunc("/api/auth/fall/signup", auth.CreateSignatory).Methods("POST")
	xrouter.HandleFunc("/api/auth/fall/signin", auth.SignInSignatory).Methods("POST")

	xrouter.HandleFunc("/api/fall/create/review", createS.Review).Methods("POST")

	xrouter.HandleFunc("/api/fall/edit/profile", editS.Profile).Methods("POST")

	xrouter.HandleFunc("/api/fall/search/contracts", searchS.Contracts).Methods("GET")
	xrouter.HandleFunc("/api/fall/search/lungh/contracts", searchSCont.Contracts).Methods("GET")
	xrouter.HandleFunc("/api/fall/search/lungh/profile", searchSCont.Profile).Methods("GET")
	xrouter.HandleFunc("/api/fall/search/lungh/reviews", searchSCont.Reviews).Methods("GET")
	xrouter.HandleFunc("/api/fall/search/self/profile", searchSSelf.Profile).Methods("GET")
	xrouter.HandleFunc("/api/fall/search/self/reviews", searchSSelf.Reviews).Methods("GET")

	xrouter.HandleFunc("/api/fall/take/contract", take.Contract).Methods("POST")
	return xrouter
}

var gRouter = myInit()
var gContractID string
var gCookiesContractor []*http.Cookie
var gCsrfContractor string
var gCookiesSignatory []*http.Cookie
var gCsrfSignatory string

func editRequest(request *http.Request, cookies []*http.Cookie, csrf string) {
	for i := range cookies {
		request.AddCookie(cookies[i])
	}

	request.Header.Set("X-Csrf-Token", csrf)
}

/**
			CONTRACTOR API TESTS
 */
func TestValidContractorCreateProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/auth/lungh/signup"
	requestBody, _ := json.Marshal(map[string]string {
		"Username"	:"testcreatecontractor",
		"Password"	:"password",
		"FirstName" :"cjdcoy",
		"LastName" 	:"lastname",
		"Email" 	:"testcreatecontractor@gmail.com",
		"PhoneNumber":"0000000000",
		"Country"	:"france",
		"Address"	:"address",
		"City"		:"city",
		"PostalCode":"postalcode",
		"BirthDay"	:"2000-01-01",
		"Sex"		:"M",
	})
	request, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
	gCookiesContractor = util.ReadSetCookies(rec.Header())
	gCsrfContractor = rec.Header()["X-Csrf-Token"][0]
}
func TestWrongContractorCreateProfile(t *testing.T) {
	route := "http://localhost:8001/api/auth/lungh/signup"
	var requestBody []byte

	for i := 0 ; i < 14 ;i++ {
		switch i {
			case 0: //already existing username
			requestBody, _ = json.Marshal(map[string]string {"username" : "testcreatecontractor", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong1@gmail.com", "phonenumber" : "0000000001", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 1: //already existing email
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong1", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testcreatecontractor@gmail.com", "phonenumber" : "0000000002", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 2: //already existing phone number
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong2", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong2@gmail.com", "phonenumber" : "0000000000", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 3: //phone number too short
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong3", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong3@gmail.com", "phonenumber" : "0123", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 14: //username too short
			requestBody, _  = json.Marshal(map[string]string {"username" : "te", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong4@gmail.com", "phonenumber" : "0000000003", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 4: //password too short
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong4", "password" : "five", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong5@gmail.com", "phonenumber" : "0000000004", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 5: //username not alphanum
			requestBody, _  = json.Marshal(map[string]string {"username" : "te$twrong5", "password" : "five", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong6@gmail.com", "phonenumber" : "0000000005", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 6: //unprintable in password
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong6", "password" : "pass万word", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong7@gmail.com", "phonenumber" : "0000000006", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 7: //too young (must >18)
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong7", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong8@gmail.com", "phonenumber" : "0000000007", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2010-01-01", "sex" : "M"})
			case 8: //wrong email
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong8", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong9gmail.com", "phonenumber" : "0000000008", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			/*case 9: //no postal code
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong9", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong9@gmail.com", "phonenumber" : "0000000009", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "", "birthday" : "2000-01-01", "sex" : "M"})
			case 10: //no address
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong10", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong10@gmail.com", "phonenumber" : "0000000010", "country" : "france", "address" : "", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 11: //no city
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong11", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong11@gmail.com", "phonenumber" : "0000000011", "country" : "france", "address" : "address", "city" : "", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 12: //already existing paypal
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong12", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong12@gmail.com", "phonenumber" : "0000000012", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M", "paypal" : "testpaypal"})
			case 13: //already existing wechat
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong13", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong13@gmail.com", "phonenumber" : "0000000013", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M", "wechat" : "testwechat"})
		*/default:
			requestBody, _ = json.Marshal(map[string]string {"username" : "testcreatecontractor", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong1@gmail.com", "phonenumber" : "0000000001", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		}
		rec := httptest.NewRecorder()

		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))

		gRouter.ServeHTTP(rec, request)
		if rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
		}
	}
}

func TestValidContractorSignIn(t *testing.T){
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/auth/lungh/signin"
	requestBody, _ := json.Marshal(map[string]string {"username" : "testcreatecontractor", "password" : "password"})
	request, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

}
func TestWrongContractorSignIn(t *testing.T) {
	var requestBody []byte
	route := "http://localhost:8001/api/auth/lungh/signin"

	for i := 0; i < 5; i++ {
		switch i {
		case 0: //wrong username
			requestBody, _ = json.Marshal(map[string]string{"username": "wrongusername", "password": "password"})
		case 1: //wrong password
			requestBody, _ = json.Marshal(map[string]string{"username": "username", "password": "wrongpassword"})
		case 2: //no username
			requestBody, _ = json.Marshal(map[string]string{"username": "", "password": "password"})
		case 3: //no password
			requestBody, _ = json.Marshal(map[string]string{"username": "username", "password": ""})
		case 4: //nousername no password
			requestBody, _ = json.Marshal(map[string]string{"username": "", "password": ""})
		}
		rec := httptest.NewRecorder()

		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))

		gRouter.ServeHTTP(rec, request)
		if rec.Result().StatusCode != http.StatusForbidden {
			log.Info(rec.Body)
			t.Errorf("expected StatusForbidden, got %v, iteration: "+strconv.FormatInt(int64(i), 10), rec.Result().Status)
		}
	}
}

func TestValidContractorEditProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/edit/profile?"
	url := route + "phonenumber=1000000000&country=germany&address=testedit&city=testedit&postalcode=testedit&password=testedit&wechat=testedit&paypal=testedit"
	request, _ := http.NewRequest("POST",  url, nil)
	editRequest(request, gCookiesContractor, gCsrfContractor)
	//this header refer to the user

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
	gCsrfContractor = rec.Header()["X-Csrf-Token"][0]
}
/*func TestWrongContractorEditProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/auth/contractor/signup?"
	requestBody, _ := json.Marshal(map[string]string {
		"Username"	:"testwrongedit",
		"Password"	:"password",
		"FirstName" :"cjdcoy",
		"LastName" 	:"lastname",
		"Email" 	:"testwrongedit@gmail.com",
		"PhoneNumber":"1100000000",
		"Country"	:"france",
		"Address"	:"address",
		"City"		:"city",
		"PostalCode":"postalcode",
		"BirthDay"	:"2000-01-01",
		"Sex"		:"M",
	})
	request, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))

gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	var urls []string
	route = "http://localhost:8001/api/contractor/edit/profile?"

	urls = append(urls,
		//already used paypal
		route+"email=diseaz91@gmail.com&phonenumber=666666666&country=germany&address=testedit&"+
			"city=testedit&postalcode=testedit&password=testedit"+
			"&paypal=testwrongeditpaypal&wechat=",
		//already used wechat
		route+"email=diseaz91@gmail.com&phonenumber=666666666&country=germany&address=testedit&city=testedit&"+
			"postalcode=testedit&password=testedit&"+
			"paypal=&wechat=testwrongeditwechat",
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesContractor, gCsrfContractor, "lungh)

gRouter.ServeHTTP(rec, request)
		if rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: "+strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}*/

func TestValidContractorCreateContract(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/create/contract?"
	url := route + "quantity=99&number=99&value=1&currency=yuan&status=open&title=" + neturl.QueryEscape("[URGENT] vend kamas & dofus pourpre mp prix") + "&body=" + neturl.QueryEscape("create contract \"test\"")
	request, _ := http.NewRequest("POST",  url, nil)
	//this header refer to the user
	editRequest(request, gCookiesContractor, gCsrfContractor)

	gRouter.ServeHTTP(rec, request)

	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
	contractID := rec.Body.String()
	gContractID = ""
	//invisible character add up at the end of the body, this is for preventing it to appear in our uuid key
	for _, v := range contractID {
		if v > 32 { //32 is the first printable character in ascii table
			gContractID += string(v)
		}
	}
	gCsrfContractor = rec.Header()["X-Csrf-Token"][0]
}
func TestWrongContractorCreateContract(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/lungh/create/contract?"

	urls = append(urls,
		//invalid char in title #1
		route + "quantity=99&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] \"vend kamas mp prix",
		//invalid char in title #2
		route + "quantity=99&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] 'vend kamas mp prix",
		//invalid quantity (too much)
		route + "quantity=100&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix",
		//invalid quantity (not enough)
		route + "quantity=0&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix",
		//invalid number (too much)
		route + "quantity=99&number=100&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix",
		//invalid number (not enough)
		route + "quantity=99&number=0&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix",
		//invalid value (not too much)
		route + "quantity=99&number=99&value=10000&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix",
		//invalid value (not enough)
		route + "quantity=99&number=99&value=0&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix",
	)


	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesContractor, gCsrfContractor)

		gRouter.ServeHTTP(rec, request)
		if rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}

func TestValidContractorEditContract(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/edit/contract?"
	url := route + "quantity=1&number=1&value=9999&currency=euro&status=open&title=" +  neturl.QueryEscape("[URGENT] vend kamas mp prix") + "&id=" + gContractID + "&body=editcontracttest"
	request, err := http.NewRequest("POST",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesContractor, gCsrfContractor)

	gRouter.ServeHTTP(rec, request)

	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
	gCsrfContractor = rec.Header()["X-Csrf-Token"][0]
}
func TestWrongContractorEditContract(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/lungh/edit/contract?"

	urls = append(urls,
		//invalid contract id
		route + "quantity=99&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix&id=invalid",
		//invalid char in title #1
		route + "quantity=99&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] \"vend kamas mp prix&id=" + gContractID,
		//invalid char in title #2
		route + "quantity=99&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] 'vend kamas mp prix&id=" + gContractID,
		//invalid quantity (too much)
		route + "quantity=100&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix&id=" + gContractID,
		//invalid quantity (not enough)
		route + "quantity=0&number=99&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix&id=" + gContractID,
		//invalid number (too much)
		route + "quantity=99&number=100&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix&id=" + gContractID,
		//invalid number (not enough)
		route + "quantity=99&number=0&value=1&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix&id=" + gContractID,
		//invalid value (not too much)
		route + "quantity=99&number=99&value=10000&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix&id=" + gContractID,
		//invalid value (not enough)
		route + "quantity=99&number=99&value=0&currency=yuan&status=waiting for approval&title=[URGENT] vend kamas mp prix&id=" + gContractID,
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesContractor, gCsrfContractor)

		gRouter.ServeHTTP(rec, request)
		if (i == 0 && rec.Result().StatusCode != http.StatusForbidden) && rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}

func TestValidContractorDeleteContract(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/delete/contract?"
	url := route + "contractid=" + gContractID
	request, err := http.NewRequest("POST",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesContractor, gCsrfContractor)

	gRouter.ServeHTTP(rec, request)

	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
	gCsrfContractor = rec.Header()["X-Csrf-Token"][0]
}
func TestWrongContractorDeleteContract(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/lungh/delete/contract?"

	urls = append(urls,
		//invalid contract id
		route + "contractid=invalid",
		//no contract id
		route + "contractid=",
		//already deleted contract ID
		route + "contractid=" + gContractID,
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesContractor, gCsrfContractor)

		gRouter.ServeHTTP(rec, request)
		if (i == 0 && rec.Result().StatusCode != http.StatusForbidden) && rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}


/**
				SIGNATORY TEST PART
*/


func TestValidSignatoryCreateProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/auth/fall/signup"
	requestBody, _ := json.Marshal(map[string]string {
		"Username"	:"testcreatesignatory",
		"Password"	:"password",
		"FirstName" :"cjdcoy",
		"LastName" 	:"lastname",
		"Email" 	:"testcreatesignatory@gmail.com",
		"PhoneNumber":"0000000000",
		"Country"	:"france",
		"Address"	:"address",
		"City"		:"city",
		"PostalCode":"postalcode",
		"BirthDay"	:"2000-01-01",
		"Sex"		:"M",
	})
	request, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))

	gRouter.ServeHTTP(rec, request)

	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
	gCsrfSignatory = rec.Header()["X-Csrf-Token"][0]
	gCookiesSignatory = util.ReadSetCookies(rec.Header())

}
func TestWrongSignatoryCreateProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/auth/fall/signup"
	var requestBody []byte

	for i := 0 ; i < 14 ;i++ {
		switch i {
		case 0: //already existing username
			requestBody, _ = json.Marshal(map[string]string {"username" : "testcreatesignatory", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong1@gmail.com", "phonenumber" : "0000000001", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		case 1: //already existing email
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong1", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testcreatesignatory@gmail.com", "phonenumber" : "0000000002", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		case 2: //already existing phone number
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong2", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong2@gmail.com", "phonenumber" : "0000000000", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		case 3: //phone number too short
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong3", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong3@gmail.com", "phonenumber" : "0123", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		case 14: //username too short
			requestBody, _  = json.Marshal(map[string]string {"username" : "te", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong4@gmail.com", "phonenumber" : "0000000003", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		case 4: //password too short
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong4", "password" : "five", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong5@gmail.com", "phonenumber" : "0000000004", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		case 5: //username not alphanum
			requestBody, _  = json.Marshal(map[string]string {"username" : "te$twrong5", "password" : "five", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong6@gmail.com", "phonenumber" : "0000000005", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		case 6: //unprintable in password
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong6", "password" : "pass万word", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong7@gmail.com", "phonenumber" : "0000000006", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		case 7: //too young (must >18)
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong7", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong8@gmail.com", "phonenumber" : "0000000007", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2010-01-01", "sex" : "M"})
		case 8: //wrong email
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong8", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong9gmail.com", "phonenumber" : "0000000008", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			/*case 9: //no postal code
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong9", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong9@gmail.com", "phonenumber" : "0000000009", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "", "birthday" : "2000-01-01", "sex" : "M"})
			case 10: //no address
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong10", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong10@gmail.com", "phonenumber" : "0000000010", "country" : "france", "address" : "", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 11: //no city
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong11", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong11@gmail.com", "phonenumber" : "0000000011", "country" : "france", "address" : "address", "city" : "", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
			case 12: //already existing paypal
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong12", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong12@gmail.com", "phonenumber" : "0000000012", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M", "paypal" : "testpaypal"})
			case 13: //already existing wechat
			requestBody, _  = json.Marshal(map[string]string {"username" : "testwrong13", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong13@gmail.com", "phonenumber" : "0000000013", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M", "wechat" : "testwechat"})
			*/default:
			requestBody, _ = json.Marshal(map[string]string {"username" : "testcreatesignatory", "password" : "password", "firstname" : "cjdcoy", "lastname" : "lastname", "email" : "testwrong1@gmail.com", "phonenumber" : "0000000001", "country" : "france", "address" : "address", "city" : "city", "postalcode" : "postalcode", "birthday" : "2000-01-01", "sex" : "M"})
		}

		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))

		gRouter.ServeHTTP(rec, request)
		if rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
		}
	}
}


func TestValidSignatoryEditProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/edit/profile?"
	url := route + "phonenumber=1000000000&country=germany&address=testedit&city=testedit&postalcode=testedit&password=testedit"
	request, _ := http.NewRequest("POST",  url, nil)

	//This header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)

	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
}
/*func TestWrongSignatoryEditProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/auth/signatory/signup?"
	requestBody, _ := json.Marshal(map[string]string {
		"Username"	:"testwrongedit",
		"Password"	:"password",
		"FirstName" :"cjdcoy",
		"LastName" 	:"lastname",
		"Email" 	:"testwrongedit@gmail.com",
		"PhoneNumber":"1100000000",
		"Country"	:"france",
		"Address"	:"address",
		"City"		:"city",
		"PostalCode":"postalcode",
		"BirthDay"	:"2000-01-01",
		"Sex"		:"M",
	})
	request, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))

gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	var urls []string
	route = "http://localhost:8001/api/signatory/edit/profile?"

	urls = append(urls,
		//already used amazon
		route + "email=diseaz91@gmail.com&phonenumber=666666666&country=germany&address=testedit&" +
			"city=testedit&postalcode=testedit&password=testedit" +
			"&amazon=testwrongeditamazon&paypal=",
		//already used paypal
		route + "email=diseaz91@gmail.com&phonenumber=666666666&country=germany&address=testedit&city=testedit&" +
			"postalcode=testedit&password=testedit&" +
			"amazon=&paypal=testwrongeditpaypal",
	)


	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesSignatory, gCsrfSignatory, "fall)

gRouter.ServeHTTP(rec, request)
		if rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}*/

func TestValidSignatorySearchContract(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/search/contracts?"
	url := route + "research=vend&page=0&limit=10"
	request, _ := http.NewRequest("GET",  url, nil)

	//This header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
}

func TestValidSignatoryTakeContract(t *testing.T) {
	TestValidContractorCreateContract(t)
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/take/contract?"
	url := route + "contractid=" + gContractID
	request, _ := http.NewRequest("POST",  url, nil)

	//This header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
}
func TestWrongSignatoryTakeContract(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/fall/take/contract?"

	urls = append(urls,
		//already taken contract
		route + "contractid=" + gContractID,
		//no contract id
		route + "contractid=",
		//wrong contract id
		route + "contractid=test",
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesSignatory, gCsrfSignatory)

		gRouter.ServeHTTP(rec, request)
		if (i == 0 && rec.Result().StatusCode != http.StatusForbidden) &&  rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}


/**
		/!\ FROM THIS POINT CONTRACTOR AND SIGNATORIES FUNCTIONS ARE MIXED /!\
				(due to the complexity of the different interactions)
 */


func TestValidContractorCreateReview(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/create/review?"
	url := route + "rating=4.2&body=review from a lungh&contractid=" + gContractID
	request, err := http.NewRequest("POST",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesContractor, gCsrfContractor)

	gRouter.ServeHTTP(rec, request)

	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
}
func TestWrongContractorCreateReview(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/lungh/create/review?"

	urls = append(urls,
		//already created review
		route + "rating=4.2&body=great review&contractid=" + gContractID,
		//no contract id
		route + "rating=4.2&body=&contractid=",
		//wrong rating number (too big)
		route + "rating=5.1&body=great review&contractid=" + gContractID,
		//wrong rating number (too little)
		route + "rating=-0.9&body=great review&contractid=" + gContractID,
		//no rating number
		route + "rating=&body=great review&contractid=" + gContractID,
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesContractor, gCsrfContractor)

		gRouter.ServeHTTP(rec, request)
		if rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}

func TestValidSignatoryCreateReview(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/create/review?"
	url := route + "rating=2.4&body=review from a fall&contractid=" + gContractID
	request, err := http.NewRequest("POST",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)

	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
}
func TestWrongSignatoryCreateReview(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/fall/create/review?"

	urls = append(urls,
		//already created review
		route + "rating=4.2&body=great review&contractid=" + gContractID,
		//no contract id
		route + "rating=4.2&body=&contractid=",
		//wrong rating number (too big)
		route + "rating=5.1&body=great review&contractid=" + gContractID,
		//wrong rating number (too little)
		route + "rating=-0.9&body=great review&contractid=" + gContractID,
		//no rating number
		route + "rating=&body=great review&contractid=" + gContractID,
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesSignatory, gCsrfSignatory)

		gRouter.ServeHTTP(rec, request)
		if rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}

func TestValidContractorSearchSelfProfile(t *testing.T) {
	var contractor database.Contractor
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/search/self/profile?"
	url := route
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesContractor, gCsrfContractor)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &contractor) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &contractor))
	}
	if contractor.Username != "testcreatecontractor" || contractor.PhoneNumber != "1000000000" || contractor.Email != "testcreatecontractor@gmail.com" {
		log.Info(rec.Body)
		t.Errorf("Wrong informations contained in the json")
	}
}
func TestWrongContractorSearchSelfProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/search/self/profile?"
	url := route
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//putting wrong Csfs to generate error
	editRequest(request, gCookiesContractor, "")
	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusInternalServerError {
		log.Info(rec.Body)
		t.Errorf("expected status StatusForbidden, got %v", rec.Result().Status)
	}
}

func TestValidContractorSearchSelfReviews(t *testing.T) {
	var reviews []database.SignatoryReview
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/search/self/reviews?"
	url := route
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesContractor, gCsrfContractor)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &reviews) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &reviews))
	}
	if reviews[0].Body != "review from a fall" || reviews[0].Rating != 2.4 {
		log.Info(rec.Body)
		t.Errorf("wrong data contained in json")
	}
}
func TestWrongContractorSearchSelfReviews(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/search/self/reviews?"
	url := route
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//putting wrong csrf to generate error
	editRequest(request, gCookiesContractor, "")

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusInternalServerError {
		log.Info(rec.Body)
		t.Errorf("expected status StatusForbidden, got %v", rec.Result().Status)
	}
}

func TestValidContractorSearchSignatoryProfile(t *testing.T) {
	var signatory database.Signatory
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/search/fall/profile?"
	url := route + "username=testcreatesignatory"
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesContractor, gCsrfContractor)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &signatory) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &signatory))
	}
	if signatory.Username != "testcreatesignatory" || signatory.PhoneNumber != "1000000000" || signatory.Email != "testcreatesignatory@gmail.com" {
		log.Info(rec.Body)
		t.Errorf("wrong data contained in json")
	}
}
func TestWrongContractorSearchSignatoryProfile(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/lungh/search/fall/profile?"

	urls = append(urls,
		//non existing user
		route + "username=testdoesntexist",
		//no username
		route + "username=",
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesContractor, gCsrfContractor)

		gRouter.ServeHTTP(rec, request)
		if (i == 0 && rec.Result().StatusCode != http.StatusForbidden) && rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}

func TestValidContractorSearchSignatoryReviews(t *testing.T) {
	var reviews []database.ContractorReview
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/lungh/search/fall/reviews?"
	url := route + "username=testcreatesignatory&page=0&limit=10"
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesContractor, gCsrfContractor)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &reviews) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &reviews))
	}
	if reviews[0].Body != "review from a lungh" || reviews[0].Rating != 4.2 {
		log.Info(rec.Body)
		t.Errorf("wrong data contained in json")
	}
}
func TestWrongContractorSearchSignatoryReviews(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/lungh/search/fall/reviews?"

	urls = append(urls,
		//non existing user
		route + "username=testdoesntexist",
		//no username
		route + "username=",
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesContractor, gCsrfContractor)

		gRouter.ServeHTTP(rec, request)
		if (i == 0 && rec.Result().StatusCode != http.StatusForbidden) && rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}

func TestValidSignatorySearchSelfProfile(t *testing.T) {
	var signatory database.Signatory
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/search/self/profile?"
	url := route
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &signatory) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &signatory))
	}
	if signatory.Username != "testcreatesignatory" || signatory.PhoneNumber != "1000000000" || signatory.Email != "testcreatesignatory@gmail.com" {
		log.Info(rec.Body)
		t.Errorf("wrong data contained in json")
	}
}
func TestWrongSignatorySearchSelfProfile(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/search/self/profile?"
	url := route
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//putting wrong csrf to generate error
	editRequest(request, gCookiesSignatory, "")

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusInternalServerError {
		log.Info(rec.Body)
		t.Errorf("expected status StatusForbidden, got %v", rec.Result().Status)
	}
}

func TestValidSignatorySearchSelfReviews(t *testing.T) {
	var reviews []database.ContractorReview
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/search/self/reviews?"
	url := route
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &reviews) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &reviews))
	}
	if reviews[0].Body != "review from a lungh" || reviews[0].Rating != 4.2 {
		log.Info(rec.Body)
		t.Errorf("wrong data contained in json")
	}
}
func TestWrongSignatorySearchSelfReviews(t *testing.T) {
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/search/self/reviews?"
	url := route
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//putting wrong token to generate error
	editRequest(request, gCookiesSignatory, "")

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusInternalServerError {
		log.Info(rec.Body)
		t.Errorf("expected status StatusForbidden, got %v", rec.Result().Status)
	}
}

func TestValidSignatorySearchContractorProfile(t *testing.T) {
	var contractor database.Contractor
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/search/lungh/profile?"
	url := route + "username=testcreatecontractor"
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &contractor) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &contractor))
	}
	if contractor.Username != "testcreatecontractor" || contractor.PhoneNumber != "1000000000" || contractor.Email != "testcreatecontractor@gmail.com" {
		log.Info(rec.Body)
		t.Errorf("wrong data contained in json")
	}
}
func TestWrongSignatorySearchContractorProfile(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/fall/search/lungh/profile?"

	urls = append(urls,
		//non existing user
		route + "username=testdoesntexist",
		//no username
		route + "username=",
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesSignatory, gCsrfSignatory)

		gRouter.ServeHTTP(rec, request)
		if (i == 0 && rec.Result().StatusCode != http.StatusForbidden) && rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}

func TestValidSignatorySearchContractorReviews(t *testing.T) {
	var reviews []database.SignatoryReview
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/search/lungh/reviews?"
	url := route + "username=testcreatecontractor&page=0&limit=10"
	request, err := http.NewRequest("GET",  url, nil)
	if err != nil { t.Errorf("expected status OK, got %v", err) }

	//this header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)
	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}
	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &reviews) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &reviews))
	}
	if reviews[0].Body != "review from a fall" || reviews[0].Rating != 2.4 {
		log.Info(rec.Body)
		t.Errorf("wrong data contained in json")
	}
}
func TestWrongSignatorySearchContractorReviews(t *testing.T) {
	var urls []string
	route := "http://localhost:8001/api/fall/search/lungh/reviews?"

	urls = append(urls,
		//non existing user
		route + "username=testdoesntexist",
		//no username
		route + "username=",
	)

	for i, u := range urls {
		rec := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", u, nil)

		//this header refer to the user
		editRequest(request, gCookiesSignatory, gCsrfSignatory)

		gRouter.ServeHTTP(rec, request)
		if (i == 0 && rec.Result().StatusCode != http.StatusForbidden) && rec.Result().StatusCode != http.StatusBadRequest {
			log.Info(rec.Body)
			t.Errorf("expected StatusBadRequest, got %v, iteration: " + strconv.FormatInt(int64(i), 10), rec.Result().Status)
			return
		}
	}
}

func TestValidSignatorySearchContractorContract(t *testing.T) {
	var contracts []database.Contract
	rec := httptest.NewRecorder()
	route := "http://localhost:8001/api/fall/search/lungh/contracts?"
	url := route + "username=testcreatecontractor"
	request, _ := http.NewRequest("GET",  url, nil)

	//This header refer to the user
	editRequest(request, gCookiesSignatory, gCsrfSignatory)

	gRouter.ServeHTTP(rec, request)

	if rec.Result().StatusCode != http.StatusOK {
		log.Info(rec.Body)
		t.Errorf("expected status OK, got %v", rec.Result().Status)
	}

	//tests on the data received
	if json.Unmarshal(rec.Body.Bytes(), &contracts) != nil {
		log.Info(rec.Body)
		t.Errorf("%v", json.Unmarshal(rec.Body.Bytes(), &contracts))
	}
	if contracts[0].Title != "[URGENT] vend kamas & dofus pourpre mp prix" || contracts[0].Body != "create contract \"test\"" {
		log.Info(rec.Body)
		t.Errorf("wrong data contained in json")
	}
}
