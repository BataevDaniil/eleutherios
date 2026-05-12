package eleutherios

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/dns"
	"github.com/BataevDaniil/eleutherios/internal/ipset"
	"github.com/BataevDaniil/eleutherios/internal/iptables"
	"github.com/BataevDaniil/eleutherios/internal/wg"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Показать текущий статус eleutherios",
	RunE: func(cmd *cobra.Command, args []string) error {
		printHdr("WireGuard")
		fmt.Println(indent(wg.Status(), "  "))

		printHdr("Интерфейсы")
		printAddrForRelevantInterfaces()

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

var statusIfaceRe = regexp.MustCompile(`^(br[0-9]+|nwg[0-9]+)$`)

func printAddrForRelevantInterfaces() {
	out, err := sh("ip", "-o", "link", "show")
	if err != nil {
		fmt.Printf("  $ ip -o link show\n")
		fmt.Printf("  ОШИБКА: %v\n", err)
		fmt.Println(indent(strings.TrimSpace(out), "    "))
		return
	}

	ifaces := relevantInterfaces(out)
	if len(ifaces) == 0 {
		fmt.Println("  brN и VPN-интерфейсы не найдены")
		return
	}

	for _, iface := range ifaces {
		printCmd("ip", "addr", "show", iface)
	}
}

func relevantInterfaces(ipLinkOut string) []string {
	var ifaces []string
	seen := map[string]bool{}
	for _, line := range splitLines(ipLinkOut) {
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 2 {
			continue
		}
		name := strings.TrimSpace(parts[1])
		if i := strings.IndexByte(name, '@'); i >= 0 {
			name = name[:i]
		}
		if name == "" || seen[name] || !statusIfaceRe.MatchString(name) {
			continue
		}
		seen[name] = true
		ifaces = append(ifaces, name)
	}
	return ifaces
}

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
