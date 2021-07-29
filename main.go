package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/gorilla/mux"
	tgrpc "github.com/refs/tpg/grpc"
	"github.com/refs/tpg/grpc/proto"
	thttp "github.com/refs/tpg/http"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

func main() {
	close := make(chan os.Signal, 1)

	setupTracing()
	go startgRPC()
	go startHTTP()

	signal.Notify(close, os.Interrupt)
	<-close
}

func startgRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8877))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterCountServiceServer(grpcServer, tgrpc.NewService())
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func startHTTP() {
	r := mux.NewRouter()
	r.HandleFunc("/count", thttp.CountHandler)
	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		panic(err)
	}
}

func setupTracing() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		panic(err)
	}

	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)
}

func tracerProvider(url string) (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in an Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Counter"),
			attribute.String("environment", "development"),
		)),
	)
	return tp, nil
}
