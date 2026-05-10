package iptables

import (
	"os/exec"
	"strings"
)

func deleteJumps(table, chain, target string) {
	out, err := exec.Command("iptables-save", "-t", table).CombinedOutput()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 || fields[0] != "-A" || fields[1] != chain {
			continue
		}
		if !hasJump(fields, target) {
			continue
		}
		args := append([]string{"-t", table, "-D", chain}, fields[2:]...)
		_ = exec.Command("iptables", args...).Run()
	}
}

func hasJump(fields []string, target string) bool {
	for i := 0; i < len(fields)-1; i++ {
		if fields[i] == "-j" && fields[i+1] == target {
			return true
		}
	}
	return false
}
