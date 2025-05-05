package picker

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

// AddrPickerName 实现指定地址调用的RPC调用
const AddrPickerName = "addr"

type addrKey struct{}

var ErrNoSubConnSelect = errors.New("no sub conn select")

func init() {
	balancer.Register(newBuilder())
}

func ContextWithAddr(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, addrKey{}, addr)
}

type addrPickerBuilder struct{}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(AddrPickerName, &addrPickerBuilder{}, base.Config{HealthCheck: true})
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
		subConnes: subConns,
	}
}

type addrPicker struct {
	subConnes map[string]balancer.SubConn
}

func (p *addrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	pr := balancer.PickResult{}

	address := info.Ctx.Value(addrKey{}).(string)
	sc, ok := p.subConnes[address]
	if !ok {
		slog.Error("Pick error", "address", address, "subConnes", p.subConnes)
		return pr, ErrNoSubConnSelect
	}
	pr.SubConn = sc
	return pr, nil
}
