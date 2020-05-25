package grpc

import (
	"context"

	. "aerobisoft.com/platform/pkg/common"
	abendpoint "aerobisoft.com/platform/pkg/endpoint"
	"aerobisoft.com/platform/pb"

	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	greeting grpc.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC GreeterServer.
func NewGRPCServer(endpoints abendpoint.Endpoints, options map[string][]grpc.ServerOption) pb.GreeterServer {
	return &grpcServer{
		greeting: makeGreetingHandler(endpoints, options[EndpointNames.Greeting]),
	}
}

// makeGreetingHandler creates the handler logic
func makeGreetingHandler(endpoints abendpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(
		endpoints.GreetingEndpoint,
		decodeGreetingRequest,
		encodeGreetingResponse,
		options...
	)
}

// Greeting implementation of the method of the GreeterService interface.
func (g *grpcServer) Greeting(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	_, rep, err := g.greeting.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GreetingResponse), nil
}
