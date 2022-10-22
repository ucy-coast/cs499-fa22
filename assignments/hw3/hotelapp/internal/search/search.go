package search

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	geo "github.com/ucy-coast/hotel-app/internal/geo/proto"
	rate "github.com/ucy-coast/hotel-app/internal/rate/proto"
	pb "github.com/ucy-coast/hotel-app/internal/search/proto"
	"github.com/ucy-coast/hotel-app/pkg/dialer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Search implements the search service
type Search struct {
	port       int
	addr       string
	geoAddr    string
	rateAddr   string
	geoClient  geo.GeoClient
	rateClient rate.RateClient
	tracer     opentracing.Tracer
}

// NewSearch returns a new server
func NewSearch(a string, p int, geoaddr string, rateaddr string, t opentracing.Tracer) *Search {
	return &Search{
		addr:     a,
		port:     p,
		geoAddr:  geoaddr,
		rateAddr: rateaddr,
		tracer:   t,
	}
}

// Run starts the server
func (s *Search) Run() error {
	func (s *Search) Run() error {
		if s.port == 0 {
			return fmt.Errorf("server port must be set")
		}
	
		opts := []grpc.ServerOption{
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Timeout: 120 * time.Second,
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				PermitWithoutStream: true,
			}),
			grpc.UnaryInterceptor(
				otgrpc.OpenTracingServerInterceptor(s.tracer),
			),
		}
	
		srv := grpc.NewServer(opts...)
		pb.RegisterSearchServer(srv, s)
	
		// Register reflection service on gRPC server.
		reflection.Register(srv)
	
		// init grpc clients
		if err := s.initGeoClient(); err != nil {
			return err
		}
		if err := s.initRateClient(); err != nil {
			return err
		}
	
		// listener
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	
		log.Printf("Start Search server. Addr: %s:%d\n", s.addr, s.port)
		return srv.Serve(lis)
	}
}

func (s *Search) initGeoClient() error {
	// TODO: Implement me
}

func (s *Search) initRateClient() error {
	// TODO: Implement me
}

// Nearby returns ids of nearby hotels ordered by ranking algo
func (s *Search) Nearby(ctx context.Context, req *pb.NearbyRequest) (*pb.SearchResult, error) {
	// TODO: Implement me
	// HINT: Reuse the implementation from the monolithic implementation 
	// HINT: and modify as needed.
}
