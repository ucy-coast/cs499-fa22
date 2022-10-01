//go:build memdb

package geo

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/hailocab/go-geoindex"
)

type DatabaseSession struct {
	points []*point
}

// point represents a hotels's geo location on map
type point struct {
	Pid  string  `json:"hotelId"`
	Plat float64 `json:"lat"`
	Plon float64 `json:"lon"`
}

func NewDatabaseSession() *DatabaseSession {
	return &DatabaseSession{}
}

func (db *DatabaseSession) LoadDataFromJsonFile(geoJsonPath string) {
	var points []*point

	// load geo points from json file
	file, err := os.Open(geoJsonPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	if err := json.Unmarshal(byteValue, &points); err != nil {
		log.Fatalf("Failed to load json: %v", err)
	}

	db.points = points
}

// newGeoIndex returns a geo index with points loaded
func (db *DatabaseSession) newGeoIndex() *geoindex.ClusteringIndex {
	// add points to index
	index := geoindex.NewClusteringIndex()
	for _, point := range db.points {
		index.Add(point)
	}

	return index
}
