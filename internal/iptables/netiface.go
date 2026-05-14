package iptables

import (
	"context"

	"github.com/BataevDaniil/eleutherios/internal/network"
)

func NetIface(ctx context.Context, name string) (string, error) {
	return network.ResolveIface(ctx, name), nil
}
