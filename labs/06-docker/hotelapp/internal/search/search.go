package search

import (
	"log"

	"github.com/ucy-coast/hotel-app/internal/geo"
	"github.com/ucy-coast/hotel-app/internal/rate"
)

// Search implements the search service
type Search struct {
	geo *geo.Geo 
	rate *rate.Rate
}

type SearchResult struct {
	HotelIds []string
}

type NearbyRequest struct {
  Lat float32
  Lon float32
  InDate string
  OutDate string
}

// NewSearch returns a new server
func NewSearch(g *geo.Geo, r *rate.Rate) *Search {
	return &Search{
		geo:  g,
		rate: r,
	}
}

// Nearby returns ids of nearby hotels ordered by ranking algo
func (s *Search) Nearby(req *NearbyRequest) (*SearchResult, error) {
		// find nearby hotels
	nearby, err := s.geo.Nearby(&geo.Request{
		Lat: req.Lat,
		Lon: req.Lon,
	})
	if err != nil {
		log.Fatalf("nearby error: %v", err)
	}

	// find rates for hotels
	rates, err := s.rate.GetRates(&rate.RateRequest{
		HotelIds: nearby.HotelIds,
		InDate:   req.InDate,
		OutDate:  req.OutDate,
	})
	if err != nil {
		log.Fatalf("rates error: %v", err)
	}

	// build the response
	res := new(SearchResult)
	for _, ratePlan := range rates.RatePlans {
		res.HotelIds = append(res.HotelIds, ratePlan.HotelId)
	}

	return res, nil
}
