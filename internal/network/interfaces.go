package network

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var ipRegex = regexp.MustCompile(`inet ([0-9.]+)/`)

// GetNetIP возвращает IP интерфейса выбранной сети
// для Keenetic API не тянем — используем ip addr
func GetNetIP(name string) (string, error) {
	if name == "" {
		name = "br0"
	}

	out, err := exec.Command("ip", "addr", "show", name).CombinedOutput()
	if err != nil {
		// fallback: ищем любой bridge
		return findBridgeIP()
	}

	matches := ipRegex.FindStringSubmatch(string(out))
	if len(matches) >= 2 {
		return matches[1], nil
	}

	return findBridgeIP()
}

// ListNetworks возвращает список доступных сетей (br0, br1...)
func ListNetworks() []string {
	var nets []string
	out, err := exec.Command("ip", "addr").CombinedOutput()
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

func ResolveIface(name string) string {
	if name == "" || name == "Home" || name == "br0" {
		return "br0"
	}
	for _, b := range Bridges() {
		if name == b.Name || name == b.LinuxName {
			return b.LinuxName
		}
	}
	return name
}

func findBridgeIP() (string, error) {
	out, err := exec.Command("ip", "addr", "show", "br0").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("не найден br0: %w", err)
	}
	matches := ipRegex.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return "", fmt.Errorf("не найден IP на br0")
	}
	return matches[1], nil
}
