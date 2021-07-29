package grpc

import (
	"context"
	"fmt"

	"github.com/refs/tpg/grpc/proto"
)

func NewService() S {
	defaultCounters := make(map[string]*counter)

	defaultCounters["default"] = &counter{
		name:   "default",
		latest: 0,
	}

	return S{
		counters: defaultCounters,
	}
}

// S implements the proto.CountServiceServer interface
type S struct {
	counters map[string]*counter
}

// Add increases a named counter by one
func (s S) Add(ctx context.Context, request *proto.AddRequest) (*proto.AddResponse, error) {
	c, ok := s.counters[request.CounterName]
	if !ok {
		return nil, fmt.Errorf("wrong counter")
	}

	c.Increase()

	r := proto.AddResponse{
		Total: c.latest,
	}

	return &r, nil
}
