package service

import (
	"context"
	"errors"
)

// Service describes your service.
type API interface {
	Greeting(ctx context.Context, name string) (string, error)
	Health(ctx context.Context) bool
}

type Service struct{}

// Greeting implementation of the Service.
func (s Service) Greeting(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", errors.New("name cannot be empty")
	}
	str := "GO-KIT Hello " + name
	return str, nil
}

// Health implementation of the Service.
func (s Service) Health(ctx context.Context) bool {
	return true
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
