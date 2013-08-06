package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/bmizerany/pat"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	baseUrl = "http://dockerbox:9200"
)

type httpHandlerFunc func(w http.ResponseWriter, req *http.Request)
type jsonHttpBuilderFunc func(req *http.Request) interface{}

func JsonBuilder(req *http.Request) interface{} {
	jsonMap := make(map[string]interface{})
	name := req.URL.Query().Get(":name")
	jsonMap["message"] = "hello " + name
	return jsonMap
}

func JsonServer(builderFunc jsonHttpBuilderFunc) (hanlderFunc httpHandlerFunc) {
	hanlderFunc = func(w http.ResponseWriter, req *http.Request) {
		enc := json.NewEncoder(w)
		jsonMap := builderFunc(req)
		w.Header().Set("Content-Type", "application/json")
		if err := enc.Encode(&jsonMap); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func PutCard(rw http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get(":username")
	project := req.URL.Query().Get(":project")
	cardnumber := req.URL.Query().Get(":cardnumber")

	body, error := ioutil.ReadAll(req.Body)
	if error != nil {
		log.Println("Error reading the bytes from the body")
		return
	}

	jsonBody, error := simplejson.NewJson(body)
	if error != nil {
		log.Println("Error marshalling json document %s", req.Body)
		return
	}

	httpClient := &http.Client{}
	url := fmt.Sprintf("%s/%s/%s_card/%s", baseUrl, username, project, cardnumber)

	jsonBytes, error := jsonBody.Encode()
	if error != nil {
		log.Printf("Error getting json bytes %s with message %s\n", url, error)
		return
	}
	log.Printf("BODY %s/%s/card/%s: %s\n", username, project, cardnumber, jsonBytes)
	clientReq, _ := http.NewRequest("PUT", url, bytes.NewReader(jsonBytes))
	response, error := httpClient.Do(clientReq)
	if error != nil {
		log.Println("Error PUT %s with message %s", url, error)
		return
	}
	responseBody, _ := ioutil.ReadAll(req.Body)
	log.Printf("status %s responseBody %s: %s\n", response.Status, url, responseBody)
	response.Body.Close()

}

func main() {
	m := pat.New()
	m.Get("/api/hello/:name", http.HandlerFunc(JsonServer(JsonBuilder)))
	m.Put("/api/:username/:project/card/:cardnumber", http.HandlerFunc(PutCard))

	http.Handle("/", m)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
