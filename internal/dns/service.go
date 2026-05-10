package dns

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ensureRunning() (bool, error) {
	if readPID(PIDFile) != "" || pidOf("dnsmasq") != "" {
		return false, nil
	}
	if _, err := os.Stat(InitFile); err != nil {
		return false, fmt.Errorf("%s не найден: установите dnsmasq-full", InitFile)
	}
	fmt.Println("  dnsmasq не запущен, стартуем...")
	out, err := exec.Command(InitFile, "start").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("запуск dnsmasq: %w (%s)", err, out)
	}
	if readPID(PIDFile) == "" && pidOf("dnsmasq") == "" {
		return false, fmt.Errorf("dnsmasq не смог запуститься (%s)", out)
	}
	return true, nil
}

func restart() error {
	out, err := exec.Command(InitFile, "restart").CombinedOutput()
	if err != nil {
		testOut, _ := exec.Command("dnsmasq", "--test").CombinedOutput()
		return fmt.Errorf("перезапуск dnsmasq: %w (%s) test: %s", err, out, testOut)
	}
	return nil
}

func readPID(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func pidOf(name string) string {
	out, err := exec.Command("pidof", name).CombinedOutput()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
