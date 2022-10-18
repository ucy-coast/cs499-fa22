//go:build memdb

package profile

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type DatabaseSession struct {
	profiles map[string]*Hotel
}

type Hotel struct {
	Id          string
	Name        string
	PhoneNumber string
	Description string
	Address     *Address
}

type Address struct {
	StreetNumber string
	StreetName   string
	City         string
	State        string
	Country      string
	PostalCode   string
	Lat          float32
	Lon          float32
}

func NewDatabaseSession() *DatabaseSession {
	return &DatabaseSession{}
}

func (db *DatabaseSession) LoadDataFromJsonFile(profilesJsonPath string) {
	var hotels []*Hotel

	// load hotel profiles from json file
	file, err := os.Open(profilesJsonPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	if err := json.Unmarshal(byteValue, &hotels); err != nil {
		log.Fatalf("Failed to load json: %v", err)
	}

	// add hotel profiles to index
	profiles := make(map[string]*Hotel)
	for _, hotel := range hotels {
		profiles[hotel.Id] = hotel
	}

	db.profiles = profiles
}

// GetProfiles returns hotel profiles for requested IDs
func (db *DatabaseSession) GetProfiles(hotelIds []string) ([]*Hotel, error) {
	hotels := make([]*Hotel, 0)

	for _, id := range hotelIds {
		hotels = append(hotels, db.getProfile(id))
	}

	return hotels, nil
}

func (db *DatabaseSession) getProfile(id string) *Hotel {
	return db.profiles[id]
}
