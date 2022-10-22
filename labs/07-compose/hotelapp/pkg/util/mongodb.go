package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	log "github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
)

func LoadDataFromJsonFile(mongosession *mgo.Session, database string, collection string, jsonPath string) {
	session := mongosession.Copy()
	defer session.Close()
	c := session.DB(database).C(collection)
	
	c.RemoveAll(nil)
	
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var v []interface{}
	if err := json.Unmarshal(byteValue, &v); err != nil {
		log.Fatalf("Failed to load json: %v", err)
	}

	if err := c.Insert(v...); err != nil {
		log.Fatalf("Failed to insert json: %v", err)
 	}
}
