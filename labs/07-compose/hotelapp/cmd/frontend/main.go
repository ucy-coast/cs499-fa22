package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/ucy-coast/hotel-app/internal/frontend"
	"github.com/ucy-coast/hotel-app/pkg/tracing"
)

var (
	port        = flag.Int("port", 8080, "The service port")
	addr        = flag.String("addr", "0.0.0.0", "Address of the service")
	jaegeraddr  = flag.String("jaeger", "jaeger:6831", "Jaeger address")
	profileaddr = flag.String("profileaddr", "profile:8081", "Address of the profile service")
	searchaddr  = flag.String("searchaddr", "search:8082", "Address of the search service")
)

func main() {
	flag.Parse()

	tracer, err := tracing.NewTracer("frontend", *jaegeraddr)
	if err != nil {
		log.Panicf("Got error while initializing jaeger agent: %v", err)
	}

	srv := frontend.NewFrontend(*addr, *port, *profileaddr, *searchaddr, tracer)

	if err := srv.Run(); err != nil {
		log.Fatalf("run main error: %v", err)
	}
}
