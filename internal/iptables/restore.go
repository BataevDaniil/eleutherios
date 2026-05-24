package iptables

import (
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
	if err := ensureChain(ctx, "nat", ChainDNS); err != nil {
		return err
	}
	if err := apply(ctx, natRules()); err != nil {
		return err
	}
	return ensurePreroutingJump(ctx, "nat", iface, ChainDNS, true)
}

func restoreMangle(ctx context.Context, iface string) error {
	if err := ipset.CreateSets(ctx); err != nil {
		return err
	}
	if err := ensureChain(ctx, "mangle", ChainMark); err != nil {
		return err
	}
	if err := apply(ctx, mangleRules()); err != nil {
		return err
	}
	return ensurePreroutingJump(ctx, "mangle", iface, ChainMark, false)
}

func ensureChain(ctx context.Context, table, chain string) error {
	if exec.CommandContext(ctx, "iptables", "-t", table, "-F", chain).Run() == nil {
		return nil
	}
	out, err := exec.CommandContext(ctx, "iptables", "-t", table, "-N", chain).CombinedOutput()
	if err != nil {
		return fmt.Errorf("iptables -t %s -N %s: %w (%s)", table, chain, err, out)
	}
	return nil
}

func apply(ctx context.Context, rules [][]string) error {
	for _, args := range rules {
		out, err := exec.CommandContext(ctx, args[0], args[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %w (%s)", args, err, out)
		}
	}
	return nil
}
