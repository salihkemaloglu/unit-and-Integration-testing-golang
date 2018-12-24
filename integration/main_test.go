package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/mgo.v2/bson"

	. "github.com/salihkemaloglu/UnitAndIntegrationTesting-Golang/operations"
)

var baseUrl = "http://api:8080/item"

func TestHttpRequestGetAll(t *testing.T) {
	response := GetAll(t)
	if len(response) != 3 {
		t.Fatal("Expected value: 3 Received value:", len(response))
	}
}
func TestHttpRequestInsert(t *testing.T) {
	responseGetBefore := GetAll(t)
	item := Item{
		Name:        "hey",
		Value:       "val",
		Description: "desc",
	}
	bytesRepresentation, err := json.Marshal(item)
	if err != nil {
		t.Fatal("Json decode error!:", err)
	}
	response, err := http.Post(baseUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal("Do request Error!:", err)
	} else if response.StatusCode != 200 {
		t.Fatal("Server side response not Ok!,Response StatusCode:", response.StatusCode)
	} else {
		defer response.Body.Close()
		var item Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			t.Fatal("Json decode error!:", err)
		}
	}

	responseGetAfter := GetAll(t)
	responseGetBeforeCount := len(responseGetBefore)
	responseGetBeforeCount++
	if responseGetBeforeCount != len(responseGetAfter) {
		t.Fatal(fmt.Printf("Insert Fail! Before Insert: %v, After Insert: %v \n", responseGetBeforeCount, len(responseGetAfter)))
	}
}

func TestHttpRequestUpdate(t *testing.T) {
	responseGetBefore := GetAll(t)
	itemGet := responseGetBefore[len(responseGetBefore)-1]
	itemGet.Name = "UpdateName"
	itemGet.Value = "UpdateValue"
	itemGet.Description = "UpdateDesc"
	url := baseUrl + "/" + bson.ObjectId(itemGet.ID).Hex()
	bytesRepresentation, err := json.Marshal(itemGet)
	if err != nil {
		t.Fatal("Json decode error!:", err)
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal("Create request Error!:", err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal("Do request Error!:", err)
	} else if response.StatusCode != 200 {
		t.Fatal("Server side response not Ok!,Response StatusCode:", response.StatusCode)
	}
	defer response.Body.Close()
	responseGetAfter := GetAll(t)
	itemUpdate := responseGetAfter[len(responseGetAfter)-1]
	if itemGet.Name != itemUpdate.Name {
		t.Fatal(fmt.Printf("Update Fail! Before Update: %v, After Update: %v \n", itemGet.Name, itemUpdate.Name))
	}

}

func TestHttpRequestDelete(t *testing.T) {
	responseGetBefore := GetAll(t)
	itemGet := responseGetBefore[0]
	url := baseUrl + "/" + bson.ObjectId(itemGet.ID).Hex()
	bytesRepresentation, err := json.Marshal(itemGet)
	if err != nil {
		t.Fatal("Json decode error!:", err)
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal("Create request Error!:", err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal("Do request Error!:", err)
	} else if response.StatusCode != 200 {
		t.Fatal("Server side response not Ok!,Response StatusCode:", response.StatusCode)
	}
	defer response.Body.Close()
	responseGetAfter := GetAll(t)
	responseGetBeforeCount := len(responseGetBefore)
	responseGetBeforeCount--
	if responseGetBeforeCount != len(responseGetAfter) {
		t.Fatal(fmt.Printf("Delete Fail! Before Delete: %v, After Delete: %v \n", responseGetBeforeCount, len(responseGetAfter)))
	}
}

func GetAll(t *testing.T) []Item {
	response, err := http.Get(baseUrl)
	if err != nil {
		t.Fatal("End point does not responde!", err.Error())
		return nil
	} else if response.StatusCode != 200 {
		t.Fatal("Server side response not Ok!,Response StatusCode:", response.StatusCode)
		return nil
	} else {
		defer response.Body.Close()
		var item []Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			t.Fatal("Json decode error!:", err)
			return nil
		} else {
			return item
		}
	}
}
