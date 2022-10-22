package rate

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	pb "github.com/ucy-coast/hotel-app/internal/rate/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Rate implements the rate service
type Rate struct {
	port      int
	addr      string
	dbsession *DatabaseSession
	tracer    opentracing.Tracer
}

// NewRate returns a new server
func NewRate(a string, p int, db *DatabaseSession, tr opentracing.Tracer) *Rate {
	return &Rate{
		addr:      a,
		port:      p,
		dbsession: db,
		tracer:    tr,
	}
}

// Run starts the server
func (s *Rate) Run() error {
	// TODO: Implement me
}

func inTimeSpan(start, end, check time.Time) bool {
	return (check.Equal(start) || check.After(start)) && (check.Equal(end) || check.Before(end))
}

// GetRates gets rates for hotels for specific date range.
func (s *Rate) GetRates(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	// TODO: Implement me
	// HINT: Reuse the implementation from the monolithic implementation 
	// HINT: and modify as needed.
}

type RatePlans []*pb.RatePlan

func (r RatePlans) Len() int {
	return len(r)
}

func (r RatePlans) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RatePlans) Less(i, j int) bool {
	return r[i].RoomType.TotalRate > r[j].RoomType.TotalRate
}
