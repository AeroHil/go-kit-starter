package grpc

import (
	"context"

	"aerobisoft.com/platform/pb"
	abendpoint "aerobisoft.com/platform/pkg/endpoint"
)

// decodeGreeterRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Greeter request.
func decodeGreetingRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GreetingRequest)
	return abendpoint.GreetingRequest{Name: req.Name}, nil
}

// encodeGreetingResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGreetingResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(abendpoint.GreetingResponse)
	return &pb.GreetingResponse{Greeting: resp.Greeting}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
