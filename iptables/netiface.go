package iptables

import (
	"os/exec"
	"strings"

	"github.com/BataevDaniil/eleutherios/network"
)

func NetIface(name string) (string, error) {
	name = network.ResolveIface(name)
	out, err := exec.Command("ip", "addr").CombinedOutput()
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
