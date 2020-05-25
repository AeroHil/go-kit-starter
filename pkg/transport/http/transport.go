package http

import (
	"context"
	"encoding/json"
	"net/http"

	abendpoint "aerobisoft.com/platform/pkg/endpoint"

	"github.com/gorilla/mux"
)

// DecodeHTTPHealthRequest method.
func DecodeHTTPHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return abendpoint.HealthRequest{}, nil
}

// DecodeHTTPGreetingRequest method for multiple query strings of the same name.
//func DecodeHTTPGreetingRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	vars := r.URL.Query()
//	names, exists := vars["name"]
//	if !exists || len(names) != 1 {
//		return nil, ErrBadRouting
//	}
//	req := abendpoint.GreetingRequest{Name: names[0]}
//	return req, nil
//}

// DecodeHTTPGreetingRequest method.
func DecodeHTTPGreetingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, ErrBadRouting
	}
	return abendpoint.GreetingRequest{Name: name}, nil
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(abendpoint.Failure); ok && f.Failed() != nil {
		EncodeError(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}
