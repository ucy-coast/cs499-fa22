//go:build memdb

package rate

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	pb "github.com/ucy-coast/hotel-app/internal/rate/proto"
)

type DatabaseSession struct {
	rateTable map[string]RatePlans
}

func NewDatabaseSession() *DatabaseSession {
	return &DatabaseSession{}
}

func (db *DatabaseSession) LoadDataFromJsonFile(rateJsonPath string) {
	rates := []*pb.RatePlan{}

	// load hotel profiles from json file
	file, err := os.Open(rateJsonPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	if err := json.Unmarshal(byteValue, &rates); err != nil {
		log.Fatalf("Failed to load json: %v", err)
	}

	// add hotel inventory to index
	rateTable := make(map[string]RatePlans)
	for _, ratePlan := range rates {
		rateTable[ratePlan.HotelId] = append(rateTable[ratePlan.HotelId], ratePlan)
	}

	db.rateTable = rateTable
}

// GetRates gets rates for hotels for specific date range.
func (db *DatabaseSession) GetRates(hotelIds []string) (RatePlans, error) {
	ratePlans := make(RatePlans, 0)

	for _, hotelID := range hotelIds {
		if db.rateTable[hotelID] != nil {
			for _, r := range db.rateTable[hotelID] {
				ratePlans = append(ratePlans, r)
			}
		}
	}

	return ratePlans, nil
}
