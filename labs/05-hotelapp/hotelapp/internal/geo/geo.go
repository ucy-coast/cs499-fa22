package geo

import (
	"github.com/hailocab/go-geoindex"
)

// Geo implements the geo subsystem
type Geo struct {
	geoidx *geoindex.ClusteringIndex
}

// The latitude and longitude of the current location.
type Request struct {
	Lat float32
	Lon float32
}

type Result struct {
	HotelIds []string
}

const (
	maxSearchRadius  = 10
	maxSearchResults = 5
)

// NewGeo returns a new Geo subsystem
func NewGeo(db *DatabaseSession) *Geo {
	return &Geo{
		geoidx: db.newGeoIndex(),
	}
}

// Nearby returns all hotels within a given distance.
func (s *Geo) Nearby(req *Request) (*Result, error) {
	var (
		points = s.getNearbyPoints(float64(req.Lat), float64(req.Lon))
		res    = &Result{}
	)
	for _, p := range points {
		res.HotelIds = append(res.HotelIds, p.Id())
	}

	return res, nil
}

func (s *Geo) getNearbyPoints(lat, lon float64) []geoindex.Point {
	center := &geoindex.GeoPoint{
		Pid:  "",
		Plat: lat,
		Plon: lon,
	}

	return s.geoidx.KNearest(
		center,
		maxSearchResults,
		geoindex.Km(maxSearchRadius), func(p geoindex.Point) bool {
			return true
		},
	)
}

// Implement Point interface
func (p *point) Lat() float64 { return p.Plat }
func (p *point) Lon() float64 { return p.Plon }
func (p *point) Id() string   { return p.Pid }
