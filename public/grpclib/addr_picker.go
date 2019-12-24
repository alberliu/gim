package grpclib

import (
	"context"
	"errors"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const Name = "addr"

const addrKey = "addr"

var ErrNoSubConnSelect = errors.New("no sub conn select")

func init() {
	balancer.Register(newBuilder())
}

func ContextWithAddr(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, addrKey, addr)
}

type addrPickerBuilder struct{}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilderWithConfig(Name, &addrPickerBuilder{}, base.Config{HealthCheck: true})
}

func (*addrPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	subConns := make(map[string]balancer.SubConn, len(readySCs))
	for k, sc := range readySCs {
		subConns[k.Addr] = sc
	}
	return &addrPicker{
		subConns: subConns,
	}
}

type addrPicker struct {
	subConns map[string]balancer.SubConn
}

func (p *addrPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	address := ctx.Value(addrKey).(string)
	sc, ok := p.subConns[address]
	if !ok {
		return nil, nil, ErrNoSubConnSelect
	}
	return sc, nil, nil
}
