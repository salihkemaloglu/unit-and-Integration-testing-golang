package main

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Find list of Itemss
func FindAll() ([]Item, error) {
	var Items []Item
	err := db.C("unit").Find(bson.M{}).All(&Items)
	return Items, err
}

// Find a Items by its id
func FindById(id string) (Item, error) {
	var Items Item
	err := db.C("unit").FindId(bson.ObjectIdHex(id)).One(&Items)
	return Items, err
}

// Insert a Items into database
func Insert(Items Item) error {
	err := db.C("unit").Insert(&Items)
	return err
}

// Delete an existing Items
func Delete(Items Item) error {
	err := db.C("unit").Remove(&Items)
	return err
}

// Update an existing Items
func Update(Items Item) error {
	err := db.C("unit").Update(bson.M{"_id": Items.ID}, &Items)
	return err
}

//Connect Establish a connection to database
func Connect(connectionUrl string) {
	info := &mgo.DialInfo{
		Addrs:    []string{connectionUrl},
		Timeout:  5 * time.Second,
		Database: DB,
		Username: "",
		Password: "",
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		fmt.Println(err.Error())
	}
	db = session.DB(DB)
}

//LoadConfiguration Parse the configuration file 'config.toml', and establish a connection to DB
func LoadConfiguration() {
	// var url = os.Getenv("HOST_ENV")
	// DB = os.Getenv("DATABASE_ENV")
	var url = "localhost:27017"
	DB = "UnitGoo"
	Connect(url)
}
