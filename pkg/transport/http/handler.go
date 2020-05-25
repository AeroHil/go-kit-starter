package http

import (
	"errors"
	"net/http"

	. "aerobisoft.com/platform/pkg/common"
	abendpoint "aerobisoft.com/platform/pkg/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints abendpoint.Endpoints, options map[string][]httptransport.ServerOption) http.Handler {
	m := mux.NewRouter()

	healthHandler := httptransport.NewServer(
		endpoints.HealthEndpoint,
		DecodeHTTPHealthRequest,
		EncodeHTTPGenericResponse,
		options[EndpointNames.Health]...,
	)

	greetingHandler := httptransport.NewServer(
		endpoints.GreetingEndpoint,
		DecodeHTTPGreetingRequest,
		EncodeHTTPGenericResponse,
		options[EndpointNames.Greeting]...,
	)

	// GET /health         retrieves service heath information
	// GET /greeting?name  retrieves greeting

	m.Methods("GET").Path("/health").Handler(healthHandler)
	m.Methods("GET").Path("/greeting/{name}").Handler(greetingHandler)

	return m
}
