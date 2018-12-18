package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/salihkemaloglu/UnitAndIntegrationTesting-Golang/operations"
)

func TestHttpRequestGetAll(t *testing.T) {
	response := GetAll()
	if response != nil {
		if len(response) != 8 {
			t.Fatal("Expected value: 3 Received value:", len(response))
		}
	} else {
		t.Fatal("End point does not responde!")
	}

}
func TestHttpRequestInsert(t *testing.T) {
	responseGetBefore := GetAll()
	if responseGetBefore == nil {
		t.Fatal("End point does not responde!")
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
	response, err := http.Post("http://localhost:8080/item", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
	} else {
		defer response.Body.Close()
		var item data.Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			fmt.Printf("Json decode error!: %s", err)
		}
	}
	responseGetAfter := GetAll()
	if responseGetAfter == nil {
		t.Fatal("End point does not responde!")
	}
	responseGetBeforeCount := len(responseGetBefore)
	responseGetBeforeCount++
	if responseGetBeforeCount != len(responseGetAfter) {
		t.Fatal(fmt.Printf("Insert Fail! Before Insert: %v, After Insert: %v \n", responseGetBeforeCount, len(responseGetAfter)))
	}
}

func TestHttpRequestUpdate(t *testing.T) {
	responseGetBefore := GetAll()
	if responseGetBefore == nil {
		t.Fatal("End point does not responde!")
	}
	itemGet := responseGetBefore[len(responseGetBefore)-1]
	itemGet.Name = "UpdateName"
	itemGet.Value = "UpdateValue"
	itemGet.Description = "UpdateDesc"
	url := "http://localhost:8080/item/" + bson.ObjectId(itemGet.ID).Hex()
	bytesRepresentation, err := json.Marshal(itemGet)
	if err != nil {
		fmt.Printf("Json decode error!: %s", err)
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
	} else {
		defer response.Body.Close()
		// contents, err := ioutil.ReadAll(response.Body)
		// if err != nil {
		// 	fmt.Printf("Json decode error!: %s", err)
		// }
		// fmt.Println("The update result is:", string(contents))
	}

	responseGetAfter := GetAll()
	if responseGetAfter == nil {
		t.Fatal("End point does not responde!")
	}
	itemUpdate := responseGetAfter[len(responseGetAfter)-1]
	if itemGet.Name != itemUpdate.Name {
		t.Fatal(fmt.Printf("Update Fail! Before Update: %v, After Update: %v \n", itemGet, itemUpdate))
	}

}

func TestHttpRequestDelete(t *testing.T) {
	responseGetBefore := GetAll()
	if responseGetBefore == nil {
		t.Fatal("End point does not responde!")
	}
	itemGet := responseGetBefore[0]
	url := "http://localhost:8080/item/" + bson.ObjectId(itemGet.ID).Hex()
	bytesRepresentation, err := json.Marshal(itemGet)
	if err != nil {
		fmt.Printf("Json decode error!: %s", err)
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
		os.Exit(1)
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
	} else {
		defer response.Body.Close()
	}

	responseGetAfter := GetAll()
	if responseGetAfter == nil {
		t.Fatal("End point does not responde!")
	}
	responseGetBeforeCount := len(responseGetBefore)
	responseGetBeforeCount--
	if responseGetBeforeCount != len(responseGetAfter) {
		t.Fatal(fmt.Printf("Delete Fail! Before Delete: %v, After Delete: %v \n", responseGetBeforeCount, len(responseGetAfter)))
	}
}

func GetAll() []data.Item {
	response, err := http.Get("http://localhost:8080/item")
	if err != nil {
		return nil
	} else {
		defer response.Body.Close()
		var item []data.Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			fmt.Printf("Json decode error!: %s", err)
		}
		return item
	}
}
