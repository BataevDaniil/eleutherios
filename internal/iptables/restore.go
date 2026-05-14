package iptables

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/BataevDaniil/eleutherios/internal/ipset"
)

func RunHook(ctx context.Context, iface string) error {
	if os.Getenv("type") != "iptables" {
		return nil
	}
	switch os.Getenv("table") {
	case "nat":
		return restoreNat(ctx, iface)
	case "mangle":
		return restoreMangle(ctx, iface)
	default:
		return nil
	}
}

func restoreNat(ctx context.Context, iface string) error {
	_ = exec.CommandContext(ctx, "iptables", "-t", "nat", "-N", ChainDNS).Run()
	_ = exec.CommandContext(ctx, "iptables", "-t", "nat", "-F", ChainDNS).Run()
	if err := apply(ctx, natRules()); err != nil {
		return err
	}
	return ensurePreroutingJump(ctx, "nat", iface, ChainDNS, true)
}

func restoreMangle(ctx context.Context, iface string) error {
	if err := ipset.CreateSets(ctx); err != nil {
		return err
	}
	_ = exec.CommandContext(ctx, "iptables", "-t", "mangle", "-N", ChainMark).Run()
	_ = exec.CommandContext(ctx, "iptables", "-t", "mangle", "-F", ChainMark).Run()
	if err := apply(ctx, mangleRules()); err != nil {
		return err
	}
	return ensurePreroutingJump(ctx, "mangle", iface, ChainMark, false)
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
