package service

import (
	"context"
	"errors"
)

// Service describes your service.
type API interface {
	Health(ctx context.Context) bool
	Greeting(ctx context.Context, name string) (string, error)
}

type Service struct{}

// Health implementation of the Service.
func (s Service) Health(ctx context.Context) bool {
	return true
}

// Greeting implementation of the Service.
func (s Service) Greeting(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", errors.New("name cannot be empty")
	}
	str := "GO-KIT Hello " + name
	return str, nil
}

// NewService returns a naive, stateless implementation of Service.
func NewService() API {
	return Service{}
}

// New returns a Service with all of the expected middleware wired in.
func New(middleware []Middleware) API {
	var svc API = NewService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
