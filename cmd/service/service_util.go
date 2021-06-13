package service

import (
	. "aerobisoft.com/platform/pkg/common"
	abendpoint "aerobisoft.com/platform/pkg/endpoint"
	abservice "aerobisoft.com/platform/pkg/service"
	abtransporthttp "aerobisoft.com/platform/pkg/transport/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/kit/transport/http"
	opentracinggo "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
)

func defaultGRPCOptions(logger log.Logger, tracer opentracinggo.Tracer, zipkinTracer *stdzipkin.Tracer) map[string][]grpc.ServerOption {
	options := map[string][]grpc.ServerOption{
		EndpointNames.Greeting: {
			grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			grpc.ServerBefore(opentracing.GRPCToContext(tracer, EndpointNames.Greeting, logger)),
		},
	}

	if zipkinTracer != nil {
		for key, value := range options {
			zt := value[:] // make a copy of the server options first
			zt = append(value, zipkin.GRPCServerTrace(zipkinTracer, zipkin.Name(key)))
			options[key] = zt
		}
	}

	return options
}

func defaultHttpOptions(logger log.Logger, tracer opentracinggo.Tracer, zipkinTracer *stdzipkin.Tracer) map[string][]http.ServerOption {
	options := map[string][]http.ServerOption{
		EndpointNames.Health: {
			http.ServerErrorEncoder(abtransporthttp.EncodeError),
			http.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			http.ServerBefore(opentracing.HTTPToContext(tracer, EndpointNames.Health, logger))},
		EndpointNames.Greeting: {
			http.ServerErrorEncoder(abtransporthttp.EncodeError),
			http.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			http.ServerBefore(opentracing.HTTPToContext(tracer, EndpointNames.Greeting, logger))},
	}

	if zipkinTracer != nil {
		for key, value := range options {
			zt := value[:] // make a copy of the server options first
			zt = append(value, zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name(key)))
			options[key] = zt
		}
	}
	return options
}

func addDefaultEndpointMiddleware(logger log.Logger, duration *prometheus.Summary) map[string][]endpoint.Middleware {
	mw := map[string][]endpoint.Middleware{}
	for epName := range EndpointNamesMap {
		mw[epName] = []endpoint.Middleware {
			abendpoint.LoggingMiddleware(log.With(logger, "method", epName)),
			abendpoint.InstrumentingMiddleware(duration.With("method", epName)),
		}
	}

	return mw
}

func addDefaultServiceMiddleware(logger log.Logger, counter metrics.Counter, latency metrics.Histogram) []abservice.Middleware {
	mw := []abservice.Middleware{
		abservice.LoggingMiddleware(logger),
		abservice.InstrumentingMiddleware(counter, latency),
	}
	return mw
}

func addEndpointMiddlewareToAllMethods(mw map[string][]endpoint.Middleware, m endpoint.Middleware) {
	methods := []string{"Greeting"}
	for _, v := range methods {
		mw[v] = append(mw[v], m)
	}
}
