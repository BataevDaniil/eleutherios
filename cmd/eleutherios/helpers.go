package eleutherios

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func sh(ctx context.Context, name string, args ...string) (string, error) {
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	return string(out), err
}

func printHdr(s string) { fmt.Printf("\n=== %s ===\n", s) }

func indent(s, prefix string) string {
	if s == "" {
		return prefix + "(нет данных)"
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func filter(s, keywords string) string {
	words := strings.Split(keywords, "|")
	var result []string
	for _, line := range strings.Split(s, "\n") {
		for _, w := range words {
			if strings.Contains(line, w) {
				result = append(result, line)
				break
			}
		}
	}
	return strings.Join(result, "\n")
}
