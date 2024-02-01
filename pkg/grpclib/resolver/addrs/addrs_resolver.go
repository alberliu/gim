package addrs

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

// 实现多个IP地址解析，比如，addrs:///127.0.0.1:50000,127.0.0.1:50001
func init() {
	RegisterResolver()
}

func RegisterResolver() {
	resolver.Register(NewAddrsBuilder())
}

type addrsBuilder struct{}

func NewAddrsBuilder() resolver.Builder {
	return &addrsBuilder{}
}

func (b *addrsBuilder) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ips := strings.Split(target.Endpoint, ",")

	state := resolver.State{
		Addresses: getAddrs(ips),
	}
	_ = clientConn.UpdateState(state)
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

func (r *addrsResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	state := resolver.State{
		Addresses: getAddrs(r.addrs),
	}
	_ = r.clientConn.UpdateState(state)
}

func (r *addrsResolver) Close() {}

func getAddrs(ips []string) []resolver.Address {
	addresses := make([]resolver.Address, len(ips))
	for i := range ips {
		addresses[i].Addr = ips[i]
	}
	return addresses
}
