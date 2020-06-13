package main

import (
	"context"
	"encoding/json"
	"encoding/csv"
	"io"
	"log"
	"fmt"
	"strings"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

var client *mongo.Client

type Company struct {
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Zip  string             `json:"zip,omitempty" bson:"zip,omitempty"`
	Website  string         `json:"website,omitempty" bson:"website,omitempty"`
}

func CreateCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	
	f, _, err := request.FormFile("file")

	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(f)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}	

		word := strings.Split(record[0], ";")
		fmt.Printf("%s       %s        ", word[0], word[1])

		company := Company{word[0], word[1], ""}

		collection := client.Database("challenge").Collection("company")
		
		insertResult, err := collection.InsertOne(context.TODO(), company)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}

	json.NewEncoder(response)
}

func GetCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	name, _ := params["name"]

	filter := bson.D{{"name", primitive.Regex{Pattern: name, Options: ""}}}

	var company Company
	collection := client.Database("challenge").Collection("company")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(company)
}

func UpdateCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	
	f, _, err := request.FormFile("file")

	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(f)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}	

		word := strings.Split(record[0], ";")
		fmt.Printf("%s       %s        ", word[0], word[1])

		filter := bson.D{{"name", word[0]}}

		update := bson.D{
			{"$set", bson.D{
				{"website", word[2]},
			}},
		}

		collection := client.Database("challenge").Collection("company")
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	json.NewEncoder(response)
}

func main() {

	fmt.Println("Starting the application...")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:PttmzyPtF0ULT8gF@cluster0-ma8hp.mongodb.net/challenge?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)

	// Init Router
	router := mux.NewRouter()

	// Route Handlers / Endpoints
	router.HandleFunc("/company", CreateCompanyEndpoint).Methods("POST")
	router.HandleFunc("/company", UpdateCompanyEndpoint).Methods("PUT")
	router.HandleFunc("/company/{name}", GetCompanyEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
