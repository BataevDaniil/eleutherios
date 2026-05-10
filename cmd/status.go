package cmd

import (
	"fmt"
	"strings"

	"github.com/BataevDaniil/eleutherios/dns"
	"github.com/BataevDaniil/eleutherios/ipset"
	"github.com/BataevDaniil/eleutherios/iptables"
	"github.com/BataevDaniil/eleutherios/wg"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Показать текущий статус eleutherios",
	RunE: func(cmd *cobra.Command, args []string) error {
		printHdr("WireGuard")
		fmt.Println(indent(wg.Status(), "  "))

		printHdr("Интерфейсы")
		printCmd("ip", "addr", "show")

		printHdr(fmt.Sprintf("Маршруты WG (таблица %d)", wg.RouteTableID))
		printCmd("ip", "route", "show", "table", fmt.Sprint(wg.RouteTableID))

		printHdr("Правила ip rule")
		printFilteredCmd("0xd1000", "ip", "rule")

		printHdr("iptables: nat")
		printCmd("iptables", "-t", "nat", "-S", iptables.ChainDNS)
		printFilteredCmd(iptables.ChainDNS, "iptables", "-t", "nat", "-S", "PREROUTING")

		printHdr("iptables: mangle")
		printCmd("iptables", "-t", "mangle", "-S", iptables.ChainMark)
		printFilteredCmd(iptables.ChainMark, "iptables", "-t", "mangle", "-S", "PREROUTING")

		printHdr("ipset: " + ipset.SetRU)
		printCmd("ipset", "list", ipset.SetRU)

		printHdr("ipset: " + ipset.SetExcluded)
		printCmd("ipset", "list", ipset.SetExcluded)

		printHdr("dnsmasq")
		fmt.Println(indent(dns.Status(), "  "))

		return nil
	},
}

func init() { rootCmd.AddCommand(statusCmd) }

func printCmd(name string, args ...string) {
	out, err := sh(name, args...)
	fmt.Printf("  $ %s\n", strings.Join(append([]string{name}, args...), " "))
	if err != nil {
		fmt.Printf("  ОШИБКА: %v\n", err)
	} else {
		fmt.Println("  OK")
	}
	fmt.Println(indent(strings.TrimSpace(out), "    "))
}

func printFilteredCmd(keywords string, name string, args ...string) {
	out, err := sh(name, args...)
	fmt.Printf("  $ %s\n", strings.Join(append([]string{name}, args...), " "))
	if err != nil {
		fmt.Printf("  ОШИБКА: %v\n", err)
		fmt.Println(indent(strings.TrimSpace(out), "    "))
		return
	}
	fmt.Println("  OK")
	fmt.Println(indent(strings.TrimSpace(filter(out, keywords)), "    "))
}
