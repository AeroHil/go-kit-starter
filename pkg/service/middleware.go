package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Middleware describes a service middleware.
type Middleware func(API) API

type loggingMiddleware struct {
	logger log.Logger
	next API
}

// LoggingMiddleware takes a logger as a dependency
// and returns a API Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next API) API {
		return &loggingMiddleware{logger, next}
	}

}

func (m loggingMiddleware) Health(ctx context.Context) (healthy bool) {
	defer func(begin time.Time) {
		m.logger.Log(
			"method", "Health",
			"healthy", healthy,
			"took", time.Since(begin),
		)
	}(time.Now())

	return m.next.Health(ctx)
}

func (m loggingMiddleware) Greeting(ctx context.Context, name string) (greeting string, err error) {
	defer func(begin time.Time) {
		m.logger.Log(
			"method", "Greeting",
			"name", name,
			"greeting", greeting,
			"took", time.Since(begin),
		)
	}(time.Now())

	return m.next.Greeting(ctx, name)
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next API
}

// InstrumentingMiddleware takes a logger as a dependency
// and returns a API Middleware.
func InstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) Middleware {
	return func(next API) API {
		return &instrumentingMiddleware{requestCount, requestLatency, next}
	}
}

func (i instrumentingMiddleware) Health(ctx context.Context) (healthy bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Health", "error", "false"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Health(ctx)
}

func (i instrumentingMiddleware) Greeting(ctx context.Context, name string) (string, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Greeting", "error", "false"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Greeting(ctx, name)
}
