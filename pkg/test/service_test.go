package test

import (
	"context"
	"testing"

	abservice "aerobisoft.com/platform/pkg/service"

	"github.com/stretchr/testify/assert"
)

func TestServiceGreeting(t *testing.T) {
	svc := abservice.NewService()
	greeting, err := svc.Greeting(context.Background(), "Test")
	assert.Nil(t, err,"Error should be nil")
	assert.Equal(t,"GO-KIT Hello Test", greeting,"Greeting should have correct string")
}

func TestServiceHealth(t *testing.T) {
	svc := abservice.NewService()
	health := svc.Health(context.Background())
	assert.NotNil(t, health,"Response should not be nil")
	assert.Equal(t,true, health,"Health should be true")
}