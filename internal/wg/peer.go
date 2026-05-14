package wg

import (
	"context"
	"fmt"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/keenetic"
	"github.com/BataevDaniil/eleutherios/internal/logging"
)

func getEntwareName(ctx context.Context, cliName string, wg *keenetic.Interface) (string, error) {
	if wg.InterfaceName != "" && ifaceExists(ctx, wg.InterfaceName) {
		return wg.InterfaceName, nil
	}
	if wg.Address == "" {
		return "", fmt.Errorf("нет адреса для %s", cliName)
	}
	out, err := execCmd(ctx, "ip", "-4", "addr")
	if err != nil {
		logging.Logger().Warn("ip -4 addr не удался", "component", "wg", "iface", cliName, "error", err)
	}
	name := findIfaceByIP(out, wg.Address)
	if name != "" {
		return name, nil
	}
	return "", fmt.Errorf("не найден linux-интерфейс с IP %s", wg.Address)
}

func ifaceExists(ctx context.Context, name string) bool {
	out, err := execCmd(ctx, "ip", "link", "show", name)
	return err == nil && out != ""
}

func findIfaceByIP(out, ip string) string {
	current := ""
	for _, line := range strings.Split(out, "\n") {
		if name := extractIfaceName(line); name != "" {
			current = name
		}
		if current != "" && containsWord(line, ip) {
			return current
		}
	}
	return ""
}

func extractIfaceName(line string) string {
	c1 := -1
	for i := 0; i < len(line); i++ {
		if line[i] == ':' {
			if c1 < 0 {
				c1 = i
			} else {
				return strings.TrimSpace(line[c1+1 : i])
			}
		}
	}
	return ""
}

func containsWord(s, word string) bool {
	for i := 0; i <= len(s)-len(word); i++ {
		if s[i:i+len(word)] == word {
			before := byte(0)
			if i > 0 {
				before = s[i-1]
			}
			after := byte(0)
			if i+len(word) < len(s) {
				after = s[i+len(word)]
			}
			if !isIPChar(before) && !isIPChar(after) {
				return true
			}
		}
	}
	return false
}

func isIPChar(c byte) bool { return (c >= '0' && c <= '9') || c == '.' }
