//go:build memdb

package profile

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	pb "github.com/ucy-coast/hotel-app/internal/profile/proto"
)

type DatabaseSession struct {
	profiles map[string]*pb.Hotel
}

func NewDatabaseSession() *DatabaseSession {
	return &DatabaseSession{}
}

func (db *DatabaseSession) LoadDataFromJsonFile(profilesJsonPath string) {
	var hotels []*pb.Hotel

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
	profiles := make(map[string]*pb.Hotel)
	for _, hotel := range hotels {
		profiles[hotel.Id] = hotel
	}

	db.profiles = profiles
}

// GetProfiles returns hotel profiles for requested IDs
func (db *DatabaseSession) GetProfiles(hotelIds []string) ([]*pb.Hotel, error) {
	hotels := make([]*pb.Hotel, 0)

	for _, id := range hotelIds {
		hotels = append(hotels, db.getProfile(id))
	}

	return hotels, nil
}

func (db *DatabaseSession) getProfile(id string) *pb.Hotel {
	return db.profiles[id]
}
