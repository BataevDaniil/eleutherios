package iptables

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

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
		out, err := exec.CommandContext(ctx, "iptables", "-w", "-t", "nat", "-N", ChainDNS).CombinedOutput()
		if err != nil {
			return fmt.Errorf("iptables -t nat -N %s: %w (%s)", ChainDNS, err, out)
		}
		if err := apply(ctx, natRules()); err != nil {
			return err
		}
	}
	return ensurePreroutingJump(ctx, "nat", iface, ChainDNS, true)
}

func restoreMangle(ctx context.Context, iface string) error {
	if err := ipset.CreateSets(ctx); err != nil {
		return err
	}
	if !chainExists(ctx, "mangle", ChainMark) {
		out, err := exec.CommandContext(ctx, "iptables", "-w", "-t", "mangle", "-N", ChainMark).CombinedOutput()
		if err != nil {
			return fmt.Errorf("iptables -t mangle -N %s: %w (%s)", ChainMark, err, out)
		}
		if err := apply(ctx, mangleRules()); err != nil {
			return err
		}
	}
	return ensurePreroutingJump(ctx, "mangle", iface, ChainMark, false)
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

func apply(ctx context.Context, rules [][]string) error {
	for _, args := range rules {
		withWait := append([]string{args[0], "-w"}, args[1:]...)
		out, err := exec.CommandContext(ctx, withWait[0], withWait[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %w (%s)", args, err, out)
		}
	}
	return nil
}
