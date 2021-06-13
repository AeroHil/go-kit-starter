package service

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"aerobisoft.com/platform/pb"
	abendpoint "aerobisoft.com/platform/pkg/endpoint"
	abservice "aerobisoft.com/platform/pkg/service"
	abgrpc "aerobisoft.com/platform/pkg/transport/grpc"
	abhttp "aerobisoft.com/platform/pkg/transport/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/oklog/oklog/pkg/group"
	prometheus1 "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func createService(endpoints abendpoint.Endpoints) (g *group.Group) {
	g = &group.Group{}
	if grpcAddr != "" {
		initGRPCHandler(endpoints, g)
	}
	if httpAddr != "" {
		initHttpHandler(endpoints, g)
	}
	return g
}

func initGRPCHandler(endpoints abendpoint.Endpoints, g *group.Group) {
	options := defaultGRPCOptions(logger, tracer, zipkinTracer)
	// Add your GRPC options here

	grpcServer := abgrpc.NewGRPCServer(endpoints, options)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", grpcAddr)
		baseServer := grpc.NewServer()
		pb.RegisterGreeterServer(baseServer, grpcServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})
}

func initHttpHandler(endpoints abendpoint.Endpoints, g *group.Group) {
	options := defaultHttpOptions(logger, tracer, zipkinTracer)
	// Add your http options here

	httpHandler := abhttp.NewHTTPHandler(endpoints, options)
	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "HTTP", "addr", httpAddr)
		return http.Serve(httpListener, httpHandler)
	}, func(error) {
		httpListener.Close()
	})

}

func getServiceMiddleware(logger log.Logger) []abservice.Middleware {
	fieldKeys := []string{"method", "error"}
	requestCount := prometheus.NewCounterFrom(prometheus1.CounterOpts{
		Namespace: "aerobisoft",
		Subsystem: environment,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := prometheus.NewSummaryFrom(prometheus1.SummaryOpts{
		Namespace: "aerobisoft",
		Subsystem: environment,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	mw := addDefaultServiceMiddleware(logger, requestCount, requestLatency)
	// Append your middleware here

	return mw
}

func getEndpointMiddleware(logger log.Logger) map[string][]endpoint.Middleware {
	duration := prometheus.NewSummaryFrom(prometheus1.SummaryOpts{
		Help:      "Request duration in seconds.",
		Name:      "request_duration_seconds",
		Namespace: "aerobisoft",
		Subsystem: environment,
	}, []string{"method", "success"})

	return addDefaultEndpointMiddleware(logger, duration)
}

func initMetricsEndpoint(g *group.Group) {
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	debugListener, err := net.Listen("tcp", debugAddr)
	if err != nil {
		logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "debug/HTTP", "addr", debugAddr)
		return http.Serve(debugListener, http.DefaultServeMux)
	}, func(error) {
		debugListener.Close()
	})
}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}
