package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/ucy-coast/hotel-app/internal/frontend"
	"github.com/ucy-coast/hotel-app/internal/geo"
	"github.com/ucy-coast/hotel-app/internal/profile"
	"github.com/ucy-coast/hotel-app/internal/rate"
	"github.com/ucy-coast/hotel-app/internal/search"
)

type server interface {
	Run() error
}

var (
	port        = flag.Int("port", 8080, "The service port")
	addr        = flag.String("addr", "0.0.0.0", "Address of the service")
	jsonDataDir = flag.String("jsondata", "data/medium", "Directory containing json data files")
)

func main() {
	flag.Parse()

	var srv server
	var g *geo.Geo
	var r *rate.Rate
	var s *search.Search
	var p *profile.Profile

	// Initialize Database
	pdb := initializeProfileDatabase()
	rdb := initializeRateDatabase()
	gdb := initializeGeoDatabase()

	g = geo.NewGeo(gdb)
	r = rate.NewRate(rdb)
	s = search.NewSearch(g, r)
	p = profile.NewProfile(pdb)

	srv = frontend.NewFrontend(*addr, *port, s, p)

	if err := srv.Run(); err != nil {
		log.Fatalf("run main error: %v", err)
	}
}
