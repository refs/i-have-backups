package http

import (
	"context"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"go.opentelemetry.io/otel"

	"github.com/refs/tpg/grpc/proto"
	"google.golang.org/grpc"
)

func CountHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer(r.URL.Path)
	cc, err := grpc.Dial("localhost:8877", grpc.WithInsecure(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		panic(err)
	}

	ctx, span := tracer.Start(context.Background(), "CountHandler")
	defer span.End()

	c := proto.NewCountServiceClient(cc)
	res, err := c.Add(ctx, &proto.AddRequest{
		CounterName: "default", // TODO read from query parameters
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(strconv.Itoa(int(res.Total))))
}
