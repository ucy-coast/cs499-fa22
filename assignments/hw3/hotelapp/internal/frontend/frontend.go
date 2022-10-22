package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"

	profile "github.com/ucy-coast/hotel-app/internal/profile/proto"
	search "github.com/ucy-coast/hotel-app/internal/search/proto"
	"github.com/ucy-coast/hotel-app/pkg/dialer"
)

// Frontend implements frontend service
type Frontend struct {
	port          int
	addr          string
	profileAddr   string
	searchAddr    string
	profileClient profile.ProfileClient
	searchClient  search.SearchClient
	tracer        opentracing.Tracer
}

// NewFrontend returns a new server
func NewFrontend(a string, p int, profaddr string, searaddr string, t opentracing.Tracer) *Frontend {
	return &Frontend{
		addr:        a,
		port:        p,
		profileAddr: profaddr,
		searchAddr:  searaddr,
		tracer:      t,
	}
}

// Run the server
func (s *Frontend) Run() error {
	// init grpc clients
	if err := s.initProfileClient(); err != nil {
		return err
	}
	if err := s.initSearchClient(); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("internal/frontend/static")))
	mux.Handle("/hotels", http.HandlerFunc(s.searchHandler))

	log.Printf("Start Frontend server. Addr: %s:%d\n", s.addr, s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
}

func (s *Frontend) initProfileClient() error {
	// TODO: Implement me	
}

func (s *Frontend) initSearchClient() error {
	conn, err := dialer.Dial(s.searchAddr, s.tracer)
	if err != nil {
		return fmt.Errorf("did not connect to search service: %v", err)
	}
	s.searchClient = search.NewSearchClient(conn)
	return nil
}

func (s *Frontend) searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := r.Context()

	// in/out dates from query params
	inDate, outDate := r.URL.Query().Get("inDate"), r.URL.Query().Get("outDate")
	if inDate == "" || outDate == "" {
		http.Error(w, "Please specify inDate/outDate params", http.StatusBadRequest)
		return
	}

	// lan/lon from query params
	sLat, sLon := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	if sLat == "" || sLon == "" {
		http.Error(w, "Please specify location params", http.StatusBadRequest)
		return
	}

	Lat, _ := strconv.ParseFloat(sLat, 32)
	lat := float32(Lat)
	Lon, _ := strconv.ParseFloat(sLon, 32)
	lon := float32(Lon)

	log.Infof("searchHandler [lat: %v, lon: %v, inDate: %v, outDate: %v]", lat, lon, inDate, outDate)
	// search for best hotels
	searchResp, err := s.searchClient.Nearby(ctx, &search.NearbyRequest{
		Lat:     lat,
		Lon:     lon,
		InDate:  inDate,
		OutDate: outDate,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// grab locale from query params or default to en
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en"
	}

	// hotel profiles
	// TODO: Implement me	

	json.NewEncoder(w).Encode(geoJSONResponse(profileResp.Hotels))
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
