package iptables

import (
	"context"
	"os/exec"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/logging"
)

func deleteJumps(ctx context.Context, table, chain, target string) {
	out, err := exec.CommandContext(ctx, "iptables-save", "-t", table).CombinedOutput()
	if err != nil {
		logging.Logger().Warn("iptables-save не удался", "component", "iptables", "table", table, "error", err, "output", string(out))
		return
	}
	for _, args := range parseDeleteCmds(string(out), table, chain, target) {
		if delOut, err := exec.CommandContext(ctx, "iptables", args...).CombinedOutput(); err != nil {
			logging.Logger().Warn("iptables -D не удался", "component", "iptables", "args", args, "error", err, "output", string(delOut))
		}
	}
}

// parseDeleteCmds выбирает из вывода `iptables-save -t <table>` строки `-A <chain> ... -j <target>`
// и превращает их в аргументы для `iptables -t <table> -D <chain> ...`.
// Чистая функция — тестируется без exec.
func parseDeleteCmds(saveOutput, table, chain, target string) [][]string {
	var out [][]string
	for _, line := range strings.Split(saveOutput, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 || fields[0] != "-A" || fields[1] != chain {
			continue
		}
		if !hasJump(fields, target) {
			continue
		}
		args := append([]string{"-t", table, "-D", chain}, fields[2:]...)
		out = append(out, args)
	}
	return out
}

func hasJump(fields []string, target string) bool {
	for i := 0; i < len(fields)-1; i++ {
		if fields[i] == "-j" && fields[i+1] == target {
			return true
		}
	}
	return false
}
