package iptables

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/BataevDaniil/eleutherios/ipset"
	"github.com/BataevDaniil/eleutherios/wg"
)

const (
	ChainDNS  = "ELEUTHERIOS_DNS"
	ChainMark = "ELEUTHERIOS_MARK"
	DNSPort   = "9753"
)

func Setup(netName, wgName string) error {
	iface, err := NetIface(netName)
	if err != nil {
		return err
	}
	Cleanup()
	if err := restoreNat(iface); err != nil {
		return err
	}
	if err := restoreMangle(iface); err != nil {
		return err
	}
	if err := InstallHook(iface); err != nil {
		return err
	}
	fmt.Printf("  iptables: %s DNS->%s, НЕ *.ru -> %s\n", iface, DNSPort, wgName)
	return nil
}

func natRules() [][]string {
	return [][]string{
		{"iptables", "-t", "nat", "-A", ChainDNS, "-p", "udp", "--dport", "53", "-j", "DNAT", "--to-destination", "127.0.0.1:" + DNSPort},
		{"iptables", "-t", "nat", "-A", ChainDNS, "-p", "tcp", "--dport", "53", "-j", "DNAT", "--to-destination", "127.0.0.1:" + DNSPort},
	}
}

func mangleRules() [][]string {
	return [][]string{
		{"iptables", "-t", "mangle", "-A", ChainMark, "-m", "set", "--match-set", ipset.SetExcluded, "dst", "-j", "RETURN"},
		{"iptables", "-t", "mangle", "-A", ChainMark, "-m", "set", "--match-set", ipset.SetRU, "dst", "-j", "RETURN"},
		{"iptables", "-t", "mangle", "-A", ChainMark, "-p", "udp", "--dport", "53", "-j", "RETURN"},
		{"iptables", "-t", "mangle", "-A", ChainMark, "-p", "tcp", "--dport", "53", "-j", "RETURN"},
		{"iptables", "-t", "mangle", "-A", ChainMark, "-j", "CONNMARK", "--restore-mark"},
		{"iptables", "-t", "mangle", "-A", ChainMark, "-m", "mark", "--mark", wg.MarkNum, "-j", "RETURN"},
		{"iptables", "-t", "mangle", "-A", ChainMark, "-m", "conntrack", "!", "--ctstate", "NEW", "-j", "RETURN"},
		{"iptables", "-t", "mangle", "-A", ChainMark, "-j", "MARK", "--set-mark", wg.MarkNum},
		{"iptables", "-t", "mangle", "-A", ChainMark, "-j", "CONNMARK", "--save-mark"},
	}
}

func ensurePreroutingJump(table, iface, target string, first bool) error {
	check := []string{"-t", table, "-C", "PREROUTING", "-i", iface, "-j", target}
	if exec.Command("iptables", check...).Run() == nil {
		return nil
	}
	args := []string{"-t", table, "-I", "PREROUTING"}
	if first {
		args = append(args, "1")
	}
	args = append(args, "-i", iface, "-j", target)
	out, err := exec.Command("iptables", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("iptables %s: %w (%s)", strings.Join(args, " "), err, out)
	}
	return nil
}

func Cleanup() {
	RemoveHook()
	deleteJumps("nat", "PREROUTING", ChainDNS)
	deleteJumps("mangle", "PREROUTING", ChainMark)
	_ = exec.Command("iptables", "-t", "nat", "-F", ChainDNS).Run()
	_ = exec.Command("iptables", "-t", "nat", "-X", ChainDNS).Run()
	_ = exec.Command("iptables", "-t", "mangle", "-F", ChainMark).Run()
	_ = exec.Command("iptables", "-t", "mangle", "-X", ChainMark).Run()
	fmt.Println("  iptables очищены")
}
