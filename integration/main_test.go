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

func TestHttpRequestGetAll(t *testing.T) {
	response, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if err != nil && eType == 1 {
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
	response, err := http.Post("http://142.93.98.36:8080/item", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Do request Error!: %s", err)
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
	url := "http://142.93.98.36:8080/item/" + bson.ObjectId(itemGet.ID).Hex()
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
		fmt.Printf("Do request Error!: %s", err)
	} else {
		defer response.Body.Close()
	}

	responseGetAfter, eType, err := GetAll()
	if err != nil && eType == 0 {
		t.Fatal("End point does not responde!", err.Error())
	} else if err != nil && eType == 1 {
		t.Fatal("Json decode error!:", err)
	}
	itemUpdate := responseGetAfter[len(responseGetAfter)-1]
	if itemGet.Name != itemUpdate.Name {
		t.Fatal(fmt.Printf("Update Fail! Before Update: %v, After Update: %v \n", itemGet, itemUpdate))
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
	url := "http://142.93.98.36:8080/item/" + bson.ObjectId(itemGet.ID).Hex()
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
		fmt.Printf("Do request Error!: %s", err)
	} else {
		defer response.Body.Close()
	}

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
	response, err := http.Get("http://142.93.98.36:8080/item")
	if err != nil {
		return nil, 0, err
	} else {
		defer response.Body.Close()
		var item []data.Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			return item, 1, err
		} else {
			return item, 2, err
		}
	}
}
