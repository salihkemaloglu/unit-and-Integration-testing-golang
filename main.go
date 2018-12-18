package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/salihkemaloglu/UnitAndIntegrationTesting-Golang/operations"
	"goji.io"
	"goji.io/pat"
	"gopkg.in/mgo.v2/bson"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!Service is working")
}

// GET list of items
func GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := data.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, items)
}

// GET a item by its ID
func GetById(w http.ResponseWriter, r *http.Request) {
	params := pat.Param(r, "id")
	if length := len(params); length != 24 {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID length")
		return
	}
	item, err := data.FindById(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	respondWithJson(w, http.StatusOK, item)
}

// POST a new item
func InsertItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var item data.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	item.ID = bson.NewObjectId()
	if err := data.Insert(item); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, item)
}

// PUT update an existing item
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var newItem data.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	params := pat.Param(r, "id")
	if length := len(params); length != 24 {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID length")
		return
	}
	oldItem, err := data.FindById(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	newItem.ID = oldItem.ID
	if err := data.Update(newItem); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing item
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := pat.Param(r, "id")
	var newItem data.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if length := len(params); length != 24 {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID length")
		return
	}
	oldItem, err := data.FindById(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	newItem.ID = oldItem.ID
	if err := data.Delete(newItem); err != nil {
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

func handleRequests() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/"), homePage)
	mux.HandleFunc(pat.Get("/item"), GetAll)
	mux.HandleFunc(pat.Get("/item/:id"), GetById)
	mux.HandleFunc(pat.Post("/item"), InsertItem)
	mux.HandleFunc(pat.Put("/item/:id"), UpdateItem)
	mux.HandleFunc(pat.Delete("/item/:id"), DeleteItem)
	log.Fatal(http.ListenAndServe(":8080", cors.AllowAll().Handler(mux)))
}

func main() {
	data.LoadConfiguration()
	handleRequests()
}
