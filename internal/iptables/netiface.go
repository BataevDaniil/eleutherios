package iptables

import (
	"context"
	"os/exec"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/network"
)

func NetIface(ctx context.Context, name string) (string, error) {
	name = network.ResolveIface(ctx, name)
	out, err := exec.CommandContext(ctx, "ip", "addr").CombinedOutput()
	if err != nil {
		return name, nil
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 1 && strings.TrimSuffix(fields[1], ":") == name {
			return name, nil
		}
	}
	return name, nil
}
