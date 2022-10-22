package dialer

import (
	"log"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Dial returns a load balanced grpc client conn with tracing interceptor
func Dial(address string, tracer opentracing.Tracer) (*grpc.ClientConn, error) {

	dialopts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Timeout:             120 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
	}

	log.Printf("Connect to %s", address)
	return grpc.Dial(address, dialopts...)
}
