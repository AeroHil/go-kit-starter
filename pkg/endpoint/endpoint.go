package endpoint

import (
	"context"

	. "aerobisoft.com/platform/pkg/common"
	abservice "aerobisoft.com/platform/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	HealthEndpoint   endpoint.Endpoint
	GreetingEndpoint endpoint.Endpoint
}

// New returns a Endpoints struct that wraps the provided service, and wires in all of the
// expected endpoint middlewares
func New(s abservice.API, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GreetingEndpoint: makeGreetingEndpoint(s),
		HealthEndpoint: makeHealthEndpoint(s),
	}
	for _, m := range mdw[EndpointNames.Health] {
		eps.HealthEndpoint = m(eps.HealthEndpoint)
	}
	for _, m := range mdw[EndpointNames.Greeting] {
		eps.GreetingEndpoint = m(eps.GreetingEndpoint)
	}

	return eps
}

// MakeHealthEndpoint constructs a Health endpoint wrapping the service.
func makeHealthEndpoint(s abservice.API) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		healthy := s.Health(ctx)
		return HealthResponse{Healthy: healthy}, nil
	}
}

// HealthRequest collects the request parameters for the Health method.
type HealthRequest struct{}

// HealthResponse collects the response values for the Health method.
type HealthResponse struct {
	Healthy bool  `json:"healthy,omitempty"`
	Err     error `json:"err,omitempty"`
}

// Failed implements Failer.
func (r HealthResponse) Failed() error { return r.Err }

// MakeGreetingEndpoint constructs a Greeter endpoint wrapping the service.
func makeGreetingEndpoint(s abservice.API) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GreetingRequest)
		greeting, err := s.Greeting(ctx, req.Name)
		if err != nil {
			return GreetingResponse{}, err
		}
		return GreetingResponse{Greeting: greeting}, nil
	}
}

// GreetingRequest collects the request parameters for the Greeting method.
type GreetingRequest struct {
	Name string `json:"name"`
}

// GreetingResponse collects the response values for the Greeting method.
type GreetingResponse struct {
	Greeting string `json:"greeting,omitempty"`
	Err      error  `json:"err,omitempty"`
}

// Failed implements Failer.
func (r GreetingResponse) Failed() error { return r.Err }

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failure, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}
