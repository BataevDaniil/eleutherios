package eleutherios

import (
	"fmt"
	"os/exec"
)

func sh(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return string(out), err
}

func printHdr(s string) { fmt.Printf("\n=== %s ===\n", s) }

func indent(s, prefix string) string {
	if s == "" {
		return prefix + "(нет данных)"
	}
	result := ""
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			result += prefix + s[start:i] + "\n"
			start = i + 1
		}
	}
	if start < len(s) {
		result += prefix + s[start:]
	}
	return result
}

func filter(s, keywords string) string {
	result := ""
	for _, line := range splitLines(s) {
		for _, w := range splitBy(keywords, '|') {
			if contains(line, w) {
				result += line + "\n"
				break
			}
		}
	}
	return result
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

func splitBy(s string, sep byte) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
