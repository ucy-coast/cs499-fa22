package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/ucy-coast/hotel-app/internal/search"
	"github.com/ucy-coast/hotel-app/pkg/tracing"
)

var (
	port       = flag.Int("port", 8082, "The service port")
	addr       = flag.String("addr", "0.0.0.0", "Address of the service")
	jaegeraddr = flag.String("jaeger", "jaeger:6831", "Jaeger address")
	geoaddr    = flag.String("geoaddr", "geo:8083", "Address of the geo service")
	rateaddr   = flag.String("rateaddr", "rate:8084", "Address of the rate service")
)

func main() {
	flag.Parse()

	tracer, err := tracing.NewTracer("geo", *jaegeraddr)
	if err != nil {
		log.Panicf("Got error while initializing jaeger agent: %v", err)
	}

	srv := search.NewSearch(*addr, *port, *geoaddr, *rateaddr, tracer)

	if err := srv.Run(); err != nil {
		log.Fatalf("run main error: %v", err)
	}
}
