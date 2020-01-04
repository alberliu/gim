package grpclib

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

func init() {
	RegisterResolver()
}

func RegisterResolver() {
	resolver.Register(NewAddrsBuilder())
}

type addrsBuilder struct {
}

func NewAddrsBuilder() resolver.Builder {
	return &addrsBuilder{}
}

func (b *addrsBuilder) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	ips := strings.Split(target.Endpoint, ",")

	state := resolver.State{
		Addresses: getAddrs(ips),
	}
	clientConn.UpdateState(state)
	return &addrsResolver{
		addrs:      ips,
		clientConn: clientConn,
	}, nil
}

func (b *addrsBuilder) Scheme() string {
	return "addrs"
}

type addrsResolver struct {
	addrs      []string
	clientConn resolver.ClientConn
}

func (r *addrsResolver) ResolveNow(opt resolver.ResolveNowOption) {
	state := resolver.State{
		Addresses: getAddrs(r.addrs),
	}
	r.clientConn.UpdateState(state)
}

func (r *addrsResolver) Close() {
}

func getAddrs(ips []string) []resolver.Address {
	addresses := make([]resolver.Address, len(ips))
	for i := range ips {
		addresses[i].Addr = ips[i]
		addresses[i].Type = resolver.Backend
	}
	return addresses
}
