package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func createAddress(w http.ResponseWriter, r *http.Request) {
	address := Address{}
	responseError := Error{}
	addressResponse := AddressResponse{}
	switch r.Method {
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &address)
		addressResponse.Data = append(addressResponse.Data, address)
	default:
		responseError.Code = 403
		responseError.Message = fmt.Sprintf("%s is not allowed on this endpoint", r.Method)
		addressResponse.Error = responseError
	}
	handleResponse(w, addressResponse)
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
