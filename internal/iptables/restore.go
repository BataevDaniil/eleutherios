package iptables

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/ipset"
	"github.com/BataevDaniil/eleutherios/internal/logging"
)

func RunHook(ctx context.Context, iface string) error {
	logger := logging.Logger()
	hookType := os.Getenv("type")
	table := os.Getenv("table")
	logger.Info("iptables hook запущен", "component", "hook", "hook", "iptables", "iface", iface, "type", hookType, "table", table)
	if hookType != "iptables" {
		logger.Info("iptables hook пропущен", "component", "hook", "hook", "iptables", "reason", "unsupported type", "type", hookType)
		return nil
	}
	switch table {
	case "nat":
		if err := restoreNat(ctx, iface); err != nil {
			return err
		}
	case "mangle":
		if err := restoreMangle(ctx, iface); err != nil {
			return err
		}
	default:
		logger.Info("iptables hook пропущен", "component", "hook", "hook", "iptables", "reason", "unsupported table", "table", table)
		return nil
	}
	logger.Info("iptables hook выполнен", "component", "hook", "hook", "iptables", "iface", iface, "table", table)
	return nil
}

func restoreNat(ctx context.Context, iface string) error {
	if !chainExists(ctx, "nat", ChainDNS) {
		if err := applyViaRestore(ctx, "nat", ChainDNS, natRules()); err != nil {
			return err
		}
	}
	if err := ensurePreroutingJump(ctx, "nat", iface, ChainDNS, true); err != nil {
		// Цепочка исчезла между chainExists и ensurePreroutingJump (TOCTOU-гонка с NDM).
		if err2 := applyViaRestore(ctx, "nat", ChainDNS, natRules()); err2 != nil {
			return err2
		}
		return ensurePreroutingJump(ctx, "nat", iface, ChainDNS, true)
	}
	return nil
}

func restoreMangle(ctx context.Context, iface string) error {
	if err := ipset.CreateSets(ctx); err != nil {
		return err
	}
	if !chainExists(ctx, "mangle", ChainMark) {
		if err := applyViaRestore(ctx, "mangle", ChainMark, mangleRules()); err != nil {
			return err
		}
	}
	if err := ensurePreroutingJump(ctx, "mangle", iface, ChainMark, false); err != nil {
		if err2 := applyViaRestore(ctx, "mangle", ChainMark, mangleRules()); err2 != nil {
			return err2
		}
		return ensurePreroutingJump(ctx, "mangle", iface, ChainMark, false)
	}
	return nil
}

// chainExists проверяет наличие цепочки через iptables-save (read-only, без side effects).
// Аналог ip4__chain__is_exist в kvas.
func chainExists(ctx context.Context, table, chain string) bool {
	out, err := exec.CommandContext(ctx, "iptables-save", "-t", table).Output()
	if err != nil {
		return false
	}
	return bytes.Contains(out, []byte(":"+chain+" "))
}

// applyViaRestore применяет цепочку и все правила одним атомарным вызовом iptables-restore --noflush.
// В отличие от N последовательных вызовов iptables, здесь NDM не может удалить цепочку между правилами.
func applyViaRestore(ctx context.Context, table, chain string, rules [][]string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "*%s\n", table)
	fmt.Fprintf(&b, ":%s - [0:0]\n", chain)
	fmt.Fprintf(&b, "-F %s\n", chain)
	for _, rule := range rules {
		// Правила хранятся как ["iptables", "-t", "TABLE", "-A", "CHAIN", ...]
		// iptables-restore принимает начиная с "-A": "-A CHAIN ..."
		fmt.Fprintf(&b, "%s\n", strings.Join(rule[3:], " "))
	}
	b.WriteString("COMMIT\n")

	cmd := exec.CommandContext(ctx, "iptables-restore", "--noflush")
	cmd.Stdin = strings.NewReader(b.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("iptables-restore -t %s: %w (%s)", table, err, out)
	}
	return nil
}
