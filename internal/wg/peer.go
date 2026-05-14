package wg

import (
	"context"
	"fmt"
)

func getEntwareName(ctx context.Context, cliName string, wg *ifaceRecord) (string, error) {
	if wg.InterfaceName != "" && ifaceExists(ctx, wg.InterfaceName) {
		return wg.InterfaceName, nil
	}
	if wg.Address == "" {
		return "", fmt.Errorf("нет адреса для %s", cliName)
	}
	out, _ := execCmd(ctx, "ip", "-4", "addr")
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
	for _, line := range splitLines(out) {
		if extractIfaceName(line) != "" {
			current = extractIfaceName(line)
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
				return trimSpace(line[c1+1 : i])
			}
		}
	}
	return ""
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
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

func trimSpace(s string) string {
	for len(s) > 0 && s[0] == ' ' {
		s = s[1:]
	}
	for len(s) > 0 && s[len(s)-1] == ' ' {
		s = s[:len(s)-1]
	}
	return s
}
