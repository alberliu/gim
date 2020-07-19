package grpclib

import (
	"context"
	"errors"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
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
	return base.NewBalancerBuilder(Name, &addrPickerBuilder{}, base.Config{HealthCheck: true})
}

func (*addrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	subConns := make(map[string]balancer.SubConn, len(info.ReadySCs))
	for k, sc := range info.ReadySCs {
		subConns[sc.Address.Addr] = k
	}
	return &addrPicker{
		subConns: subConns,
	}
}

type addrPicker struct {
	subConns map[string]balancer.SubConn
}

func (p *addrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	pr := balancer.PickResult{}

	address := info.Ctx.Value(addrKey).(string)
	sc, ok := p.subConns[address]
	if !ok {
		return pr, ErrNoSubConnSelect
	}
	pr.SubConn = sc
	return pr, nil
}
