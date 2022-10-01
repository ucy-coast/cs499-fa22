package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ucy-coast/hotel-app/internal/profile"
	"github.com/ucy-coast/hotel-app/internal/search"

)

// Frontend implements frontend service
type Frontend struct {
	port          int
	addr          string
	search  *search.Search
	profile  *profile.Profile
}

// NewFrontend returns a new server
func NewFrontend(a string, p int, s *search.Search, pr *profile.Profile) *Frontend {
	return &Frontend{
		addr:    a,
		port:    p,
		search:  s,
		profile: pr,
	}
}

// Run the server
func (s *Frontend) Run() error {
	// TODO: Implement me
	// HINT: Follow the instructions in Lab 05: Getting Started with Go Web Apps
}

func (s *Frontend) searchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
	// HINT: Follow the instructions in Lab 05: Getting Started with Go Web Apps
}

// return a geoJSON response that allows google map to plot points directly on map
// https://developers.google.com/maps/documentation/javascript/datalayer#sample_geojson
func geoJSONResponse(hs []*profile.Hotel) map[string]interface{} {
	fs := []interface{}{}

	for _, h := range hs {
		fs = append(fs, map[string]interface{}{
			"type": "Feature",
			"id":   h.Id,
			"properties": map[string]string{
				"name":         h.Name,
				"phone_number": h.PhoneNumber,
			},
			"geometry": map[string]interface{}{
				"type": "Point",
				"coordinates": []float32{
					h.Address.Lon,
					h.Address.Lat,
				},
			},
		})
	}

	return map[string]interface{}{
		"type":     "FeatureCollection",
		"features": fs,
	}
}
