package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Item - Our struct for all items
type Item struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Value       string        `bson:"value" json:"value"`
	Description string        `bson:"description" json:"description"`
}

type Config struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	Collection string `json:"collection"`
}

type Items []Item

var db *mgo.Database

var COLLECTION string

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!Service is working")
}

// GET list of items
func AllitemsEndPoint(w http.ResponseWriter, r *http.Request) {
	items, err := FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, items)
}

// GET a item by its ID
func FinditemsEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	item, err := FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	respondWithJson(w, http.StatusOK, item)
}

// POST a new item
func CreateItemsEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	item.ID = bson.NewObjectId()
	if err := Insert(item); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, item)
}

// PUT update an existing item
func UpdateItemsEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := Update(item); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing item
func DeleteItemsEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := Delete(item); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Establish a connection to database
func Connect(connectionUrl string) {
	session, err := mgo.Dial(connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(COLLECTION)
}

// Find list of Itemss
func FindAll() ([]Item, error) {
	var Items []Item
	err := db.C(COLLECTION).Find(bson.M{}).All(&Items)
	return Items, err
}

// Find a Items by its id
func FindById(id string) (Item, error) {
	var Items Item
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&Items)
	return Items, err
}

// Insert a Items into database
func Insert(Items Item) error {
	err := db.C(COLLECTION).Insert(&Items)
	return err
}

// Delete an existing Items
func Delete(Items Item) error {
	err := db.C(COLLECTION).Remove(&Items)
	return err
}

// Update an existing Items
func Update(Items Item) error {
	err := db.C(COLLECTION).Update(bson.M{"_id": Items.ID}, &Items)
	return err
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func LoadConfiguration() {
	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	var url = config.Host + ":" + config.Port
	COLLECTION = config.Collection
	Connect(url)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/item", AllitemsEndPoint).Methods("GET")
	myRouter.HandleFunc("/item", CreateItemsEndPoint).Methods("POST")
	myRouter.HandleFunc("/item", UpdateItemsEndPoint).Methods("PUT")
	myRouter.HandleFunc("/item", DeleteItemsEndPoint).Methods("DELETE")
	myRouter.HandleFunc("/item/{id}", FinditemsEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	// LoadConfiguration()
	handleRequests()
}
