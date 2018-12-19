package data

import (
	"fmt"
	"os"
	"time"

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

type Items []Item

var db *mgo.Database

//DB string
var DB string

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
		fmt.Println("Mongodb connection error!:", err.Error())
	}
	db = session.DB(DB)
}

//LoadConfiguration Parse the configuration file 'config.toml', and establish a connection to DB
func LoadConfiguration() {
	var url = os.Getenv("HOST_ENV")
	DB = os.Getenv("DATABASE_ENV")
	Connect(url)
}
