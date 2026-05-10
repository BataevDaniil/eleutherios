package iptables

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/BataevDaniil/eleutherios/internal/ipset"
)

func RunHook(iface string) error {
	if os.Getenv("type") != "iptables" {
		return nil
	}
	switch os.Getenv("table") {
	case "nat":
		return restoreNat(iface)
	case "mangle":
		return restoreMangle(iface)
	default:
		return nil
	}
}

func restoreNat(iface string) error {
	_ = exec.Command("iptables", "-t", "nat", "-N", ChainDNS).Run()
	_ = exec.Command("iptables", "-t", "nat", "-F", ChainDNS).Run()
	if err := apply(natRules()); err != nil {
		return err
	}
	return ensurePreroutingJump("nat", iface, ChainDNS, true)
}

func restoreMangle(iface string) error {
	if err := ipset.CreateSets(); err != nil {
		return err
	}
	_ = exec.Command("iptables", "-t", "mangle", "-N", ChainMark).Run()
	_ = exec.Command("iptables", "-t", "mangle", "-F", ChainMark).Run()
	if err := apply(mangleRules()); err != nil {
		return err
	}
	return ensurePreroutingJump("mangle", iface, ChainMark, false)
}

func apply(rules [][]string) error {
	for _, args := range rules {
		out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %w (%s)", args, err, out)
		}
	}
	return nil
}
