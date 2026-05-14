package network

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var ipRegex = regexp.MustCompile(`inet ([0-9.]+)/`)

// GetNetIP возвращает IP интерфейса выбранной сети
// для Keenetic API не тянем — используем ip addr
func GetNetIP(ctx context.Context, name string) (string, error) {
	if name == "" {
		name = "br0"
	}

	out, err := exec.CommandContext(ctx, "ip", "addr", "show", name).CombinedOutput()
	if err != nil {
		// fallback: ищем любой bridge
		return findBridgeIP(ctx)
	}

	matches := ipRegex.FindStringSubmatch(string(out))
	if len(matches) >= 2 {
		return matches[1], nil
	}

	return findBridgeIP(ctx)
}

// ListNetworks возвращает список доступных сетей (br0, br1...)
func ListNetworks(ctx context.Context) []string {
	var nets []string
	out, err := exec.CommandContext(ctx, "ip", "addr").CombinedOutput()
	if err != nil {
		return []string{"br0"}
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "br") && strings.Contains(line, ": <") {
			parts := strings.SplitN(line, ":", 3)
			if len(parts) >= 2 {
				nets = append(nets, strings.TrimSpace(parts[1]))
			}
		}
	}
	if len(nets) == 0 {
		return []string{"br0"}
	}
	return nets
}

func ResolveIface(ctx context.Context, name string) string {
	if name == "" || name == "Home" || name == "br0" {
		return "br0"
	}
	for _, b := range Bridges(ctx) {
		if name == b.Name || name == b.LinuxName {
			return b.LinuxName
		}
	}
	return name
}

// ValidateIface проверяет, что сеть существует. Возвращает linux-имя или ошибку
// со списком доступных bridge-сетей.
func ValidateIface(ctx context.Context, name string) (string, error) {
	if name == "" {
		name = "br0"
	}
	resolved := ResolveIface(ctx, name)
	if _, err := exec.CommandContext(ctx, "ip", "link", "show", resolved).CombinedOutput(); err == nil {
		return resolved, nil
	}
	available := availableNetworks(ctx)
	if len(available) == 0 {
		return "", fmt.Errorf("сеть %q не найдена и список bridge пуст", name)
	}
	return "", fmt.Errorf("сеть %q не найдена. Доступные: %s", name, strings.Join(available, ", "))
}

func availableNetworks(ctx context.Context) []string {
	bridges := Bridges(ctx)
	if len(bridges) > 0 {
		out := make([]string, 0, len(bridges))
		for _, b := range bridges {
			out = append(out, b.Name+" ("+b.LinuxName+")")
		}
		return out
	}
	return ListNetworks(ctx)
}

func findBridgeIP(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "ip", "addr", "show", "br0").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("не найден br0: %w", err)
	}
	matches := ipRegex.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return "", fmt.Errorf("не найден IP на br0")
	}
	return matches[1], nil
}
