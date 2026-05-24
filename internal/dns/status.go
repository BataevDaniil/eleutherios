package dns

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/logging"
)

// Cleanup останавливает Entware dnsmasq и возвращает /opt/etc/dnsmasq.conf
// в безопасный embedded rollback. DNS клиентов после этого обслуживает
// встроенный в Keenetic NDM (ndnproxy + iptables redirect на :53) — то же
// состояние, что было до eleutherios start.
func Cleanup(ctx context.Context) error {
	if err := os.Remove(ConfFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("удаление %s: %w", ConfFile, err)
	}
	if err := os.WriteFile(BaseConfFile, []byte(rollbackConfig), 0644); err != nil {
		return fmt.Errorf("запись %s: %w", BaseConfFile, err)
	}
	if err := stopDnsmasq(ctx); err != nil {
		return err
	}
	logging.Logger().Info("dnsmasq остановлен, DNS обслуживается NDM", "component", "dnsmasq", "config", BaseConfFile)
	return nil
}

func stopDnsmasq(ctx context.Context) error {
	if readPID(PIDFile) == "" && pidOf(ctx, "dnsmasq") == "" {
		return nil
	}
	if _, err := os.Stat(InitFile); err != nil {
		return nil
	}
	out, err := exec.CommandContext(ctx, InitFile, "stop").CombinedOutput()
	if err != nil {
		return fmt.Errorf("остановка dnsmasq: %w (%s)", err, out)
	}
	return nil
}

func Status(ctx context.Context) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("init: %s\n", pathState(InitFile)))
	b.WriteString(fmt.Sprintf("pid-file: %s (%s)\n", PIDFile, readPID(PIDFile)))
	b.WriteString(fmt.Sprintf("pidof dnsmasq: %s\n", valueOr(pidOf(ctx, "dnsmasq"), "нет")))
	out, err := exec.CommandContext(ctx, InitFile, "status").CombinedOutput()
	b.WriteString(fmt.Sprintf("$ %s status: %v\n%s\n", InitFile, err, strings.TrimSpace(string(out))))
	b.WriteString(fmt.Sprintf("config: %s\n", pathState(ConfFile)))
	if data, err := os.ReadFile(ConfFile); err == nil {
		b.WriteString(strings.TrimSpace(string(data)) + "\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func pathState(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "нет (" + err.Error() + ")"
	}
	return fmt.Sprintf("есть, %d байт", info.Size())
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
