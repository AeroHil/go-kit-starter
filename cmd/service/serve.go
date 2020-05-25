package service

import (
	"os"

	abendpoint "aerobisoft.com/platform/pkg/endpoint"
	abservice "aerobisoft.com/platform/pkg/service"

	"github.com/go-kit/kit/log"
	opentracinggo "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/urfave/cli/v2"
)

var tracer opentracinggo.Tracer
var zipkinTracer *zipkin.Tracer
var logger log.Logger

func Run(c *cli.Context) error {
	// Create a single logger, which we'll use and give to other components.
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	logger.Log("environment", environment)

	// Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency
	if zipkinAddr != "" {
		logger.Log("tracer", "Zipkin", "type", "OpenTracing", "URL", zipkinAddr)

		var (
			err         error
			hostPort    = "localhost:80"
			serviceName = "go-kit-service"
			reporter    = zipkinhttp.NewReporter(zipkinAddr)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP))
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}

		tracer = zipkinot.Wrap(zipkinTracer)
	} else {
		logger.Log("tracer", "none")
		tracer = opentracinggo.GlobalTracer()
	}

	svc := abservice.New(getServiceMiddleware(logger))
	eps := abendpoint.New(svc, getEndpointMiddleware(logger))
	g := createService(eps)
	initMetricsEndpoint(g)
	initCancelInterrupt(g)
	logger.Log("exit", g.Run())

	return nil
}
