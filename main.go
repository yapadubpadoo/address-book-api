package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// Address ...
type Address struct {
	Name  string
	Phone string
}

type Error struct {
	Code    int
	Message string
}

type AddressResponse struct {
	Data  []Address
	Error Error
}

var dbSession = getDBSession()

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API is up")
}

func handleResponse(w http.ResponseWriter, adr AddressResponse) {
	addressResponseJSON, _ := json.Marshal(adr)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(addressResponseJSON))
}

func validateJSONRequest(r *http.Request) (bool, string) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return false, "Request body must be JSON"
	}
	return true, ""
}

func createAddress(w http.ResponseWriter, r *http.Request) {
	address := Address{}
	responseError := Error{}
	addressResponse := AddressResponse{}
	valid, err := validateJSONRequest(r)
	if valid {
		switch r.Method {
		case "POST":
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Println("Going to decode JSON")
			error := json.Unmarshal(body, &address)
			if error != nil {
				responseError.Message = fmt.Sprintf("%+v", error)
			} else {
				c := dbSession.DB("address_book").C("people")
				err := c.Insert(address)
				if err != nil {
					log.Fatal(err)
				}
			}
			addressResponse.Data = append(addressResponse.Data, address)
		default:
			responseError.Code = 403
			responseError.Message = fmt.Sprintf("%s is not allowed on this endpoint", r.Method)
		}
	} else {
		responseError.Code = -1
		responseError.Message = fmt.Sprintf("%s", err)
	}
	addressResponse.Error = responseError
	handleResponse(w, addressResponse)
}

func getDBSession() *mgo.Session {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	return session
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/healthcheck", healthcheck)

	// # /create => create record => POST
	http.HandleFunc("/create", createAddress)

	// # /record => get all records => GET
	// # /record/id => get record with specificy id => GET
	// # /record/id/update => update record with specify id => PUT
	// # /record/id => delete specify record => DELETE
	// http.HandleFunc("/record", healthcheck)

	// http.HandleFunc("/record/:id", healthcheck)

	// http.HandleFunc("/record/:id", healthcheck)

	http.ListenAndServe(":8080", nil)
}
