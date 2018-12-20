package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/salihkemaloglu/UnitAndIntegrationTesting-Golang/operations"
)

var baseUrl = "http://api:8080/item"

func TestHttpRequestGetAll(t *testing.T) {
	response, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if eType == 1 {
		t.Fatal("Page not found!")
	} else if err != nil && eType == 2 {
		t.Fatal("Json decode error!:", err)
	} else {
		if len(response) != 3 {
			t.Fatal("Expected value: 3 Received value:", len(response))
		}
	}
}
func TestHttpRequestInsert(t *testing.T) {
	responseGetBefore, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if err != nil && eType == 1 {
		t.Fatal("Json decode error!:", err)
	}
	item := data.Item{
		Name:        "hey",
		Value:       "val",
		Description: "desc",
	}
	bytesRepresentation, err := json.Marshal(item)
	if err != nil {
		fmt.Printf("Json decode error!: %s", err)
	}
	response, err := http.Post(baseUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Do request Error!: %s", err)
	} else if response.StatusCode == 404 {
		t.Fatal("Page not found!")
	} else {
		defer response.Body.Close()
		var item data.Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			fmt.Printf("Json decode error!: %s", err)
		}
	}

	responseGetAfter, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if err != nil && eType == 1 {
		t.Fatal("Json decode error!:", err)
	}
	responseGetBeforeCount := len(responseGetBefore)
	responseGetBeforeCount++
	if responseGetBeforeCount != len(responseGetAfter) {
		t.Fatal(fmt.Printf("Insert Fail! Before Insert: %v, After Insert: %v \n", responseGetBeforeCount, len(responseGetAfter)))
	}
}

func TestHttpRequestUpdate(t *testing.T) {
	responseGetBefore, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if err != nil && eType == 1 {
		t.Fatal("Json decode error!:", err)
	}
	itemGet := responseGetBefore[len(responseGetBefore)-1]
	itemGet.Name = "UpdateName"
	itemGet.Value = "UpdateValue"
	itemGet.Description = "UpdateDesc"
	url := baseUrl + "/" + bson.ObjectId(itemGet.ID).Hex()
	bytesRepresentation, err := json.Marshal(itemGet)
	if err != nil {
		fmt.Printf("Json decode error!: %s", err)
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Create request Error!: %s", err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal("Do request Error!:", err)
	} else if response.StatusCode == 404 {
		t.Fatal("Page not found error!:", err)
	}
	defer response.Body.Close()
	responseGetAfter, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if err != nil && eType == 1 {
		t.Fatal("Json decode error!:", err)
	}
	itemUpdate := responseGetAfter[len(responseGetAfter)-1]
	if itemGet.Name != itemUpdate.Name {
		t.Fatal(fmt.Printf("Update Fail! Before Update: %v, After Update: %v \n", itemGet.Name, itemUpdate.Name))
	}

}

func TestHttpRequestDelete(t *testing.T) {
	responseGetBefore, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if err != nil && eType == 1 {
		t.Fatal("Json decode error!:", err)
	}
	itemGet := responseGetBefore[0]
	url := baseUrl + "/" + bson.ObjectId(itemGet.ID).Hex()
	bytesRepresentation, err := json.Marshal(itemGet)
	if err != nil {
		fmt.Printf("Json decode error!: %s", err)
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Create request Error!: %s", err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal("Do request Error!:", err)
	} else if response.StatusCode == 404 {
		t.Fatal("Page not found error!:", err)
	}
	defer response.Body.Close()
	responseGetAfter, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if err != nil && eType == 1 {
		t.Fatal("Json decode error!:", err)
	}
	responseGetBeforeCount := len(responseGetBefore)
	responseGetBeforeCount--
	if responseGetBeforeCount != len(responseGetAfter) {
		t.Fatal(fmt.Printf("Delete Fail! Before Delete: %v, After Delete: %v \n", responseGetBeforeCount, len(responseGetAfter)))
	}
}

func GetAll() ([]data.Item, int, error) {
	response, err := http.Get(baseUrl)
	if err != nil {
		return nil, 0, err
	} else if response.StatusCode == 404 {
		return nil, 1, err
	} else {
		defer response.Body.Close()
		var item []data.Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			return item, 2, err
		} else {
			return item, 3, err
		}
	}
}
