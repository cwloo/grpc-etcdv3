package getcdv3

import (
	grpc_resolver "google.golang.org/grpc/resolver"
)

// <summary>
// Resolver
// <summary>
type Resolver struct {
	grpc_resolver.Resolver
	target string
}

func newResolver(target string) *Resolver {
	return &Resolver{
		target: target,
	}
}

// ResolveNow
func (s *Resolver) ResolveNow(rn grpc_resolver.ResolveNowOptions) {
}

// Close
func (s *Resolver) Close() {
	// logs.Errorf("%v", s.target)
}
