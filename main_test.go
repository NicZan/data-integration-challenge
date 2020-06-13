package main

import (
	"encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"fmt"
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"io"
	"context"
	"time"
	"strings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"
)

func Router() *mux.Router {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:PttmzyPtF0ULT8gF@cluster0-ma8hp.mongodb.net/challenge?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)

    router := mux.NewRouter()
	router.HandleFunc("/company", CreateCompanyEndpoint).Methods("POST")
	router.HandleFunc("/company", UpdateCompanyEndpoint).Methods("PUT")
	router.HandleFunc("/company/{name}", GetCompanyEndpoint).Methods("GET")
    return router
}


//Should import csv q1_catalog.csv

func TestCreateCompany(t *testing.T) {

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("q1_catalog.csv")
	defer file.Close()
	part1,
			errFile1 := writer.CreateFormFile("file",filepath.Base("q1_catalog.csv"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 !=nil {
			
		fmt.Println(errFile1)
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(payload)

	request, _ := http.NewRequest("POST", "/company", payload)
	request.Header.Set("Content-Type", writer.FormDataContentType())
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
    assert.Equal(t, 200, response.Code, "OK response is expected")
}


//Should verify if as imported 
func TestGetCompanyBefore(t *testing.T) {

	// Build our expected body
	body := gin.H{
	   "name": "tola sales group",
	   "zip": "78229",
	   "website": "",
   }

   payload := strings.NewReader("")
   request, _ := http.NewRequest("GET", "/company/tola",payload)
   response := httptest.NewRecorder()
   Router().ServeHTTP(response, request)
   assert.Equal(t, 200, response.Code, "OK response is expected")

	var responseJson map[string]string
	err := json.Unmarshal([]byte(response.Body.String()), &responseJson)
	assert.Nil(t, err)

	valueName, exists1 := responseJson["name"]
	assert.True(t, exists1)
	assert.Equal(t, body["name"], valueName)
	valueZip, exists2 := responseJson["zip"]
	assert.True(t, exists2)
	assert.Equal(t, body["zip"], valueZip)
	valueWebsite, exists3 := responseJson["website"]
	assert.False(t, exists3)
	assert.Equal(t, body["website"], valueWebsite)
}

//Should update website value from csv q2_clientData.csv
func TestUpdateCompany(t *testing.T) {

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("q2_clientData.csv")
	defer file.Close()
	part1,
			errFile1 := writer.CreateFormFile("file",filepath.Base("q2_clientData.csv"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 !=nil {
			
		fmt.Println(errFile1)
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(payload)

	request, _ := http.NewRequest("PUT", "/company", payload)
	request.Header.Set("Content-Type", writer.FormDataContentType())
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
    assert.Equal(t, 200, response.Code, "OK response is expected")
}


//Should verify if as updated website value 
func TestGetCompanyAfter(t *testing.T) {

	 // Build our expected body
	 body := gin.H{
		"name": "tola sales group",
		"zip": "78229",
		"website": "http://repsources.com",
	}

	payload := strings.NewReader("")
	request, _ := http.NewRequest("GET", "/company/tola%20sales%20group",payload)
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

	 var responseJson map[string]string
	 err := json.Unmarshal([]byte(response.Body.String()), &responseJson)
	 assert.Nil(t, err)

	 valueName, exists1 := responseJson["name"]
	 assert.True(t, exists1)
	 assert.Equal(t, body["name"], valueName)
	 valueZip, exists2 := responseJson["zip"]
	 assert.True(t, exists2)
	 assert.Equal(t, body["zip"], valueZip)
	 valueWebsite, exists3 := responseJson["website"]
	 assert.True(t, exists3)
	 assert.Equal(t, body["website"], valueWebsite)
}