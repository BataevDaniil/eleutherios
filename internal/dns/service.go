package dns

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/logging"
)

func ensureRunning(ctx context.Context) (bool, error) {
	if readPID(PIDFile) != "" || pidOf(ctx, "dnsmasq") != "" {
		return false, nil
	}
	if _, err := os.Stat(InitFile); err != nil {
		return false, fmt.Errorf("%s не найден: установите dnsmasq-full", InitFile)
	}
	logging.Logger().Info("dnsmasq не запущен, стартуем", "component", "dnsmasq")
	out, err := exec.CommandContext(ctx, InitFile, "start").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("запуск dnsmasq: %w (%s)", err, out)
	}
	if readPID(PIDFile) == "" && pidOf(ctx, "dnsmasq") == "" {
		return false, fmt.Errorf("dnsmasq не смог запуститься (%s)", out)
	}
	return true, nil
}

func restart(ctx context.Context) error {
	out, err := exec.CommandContext(ctx, InitFile, "restart").CombinedOutput()
	if err != nil {
		testOut, _ := exec.CommandContext(ctx, "dnsmasq", "--test").CombinedOutput()
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

func pidOf(ctx context.Context, name string) string {
	out, err := exec.CommandContext(ctx, "pidof", name).CombinedOutput()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
