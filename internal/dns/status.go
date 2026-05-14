package dns

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Cleanup(ctx context.Context) error {
	if err := os.Remove(ConfFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("удаление %s: %w", ConfFile, err)
	}
	if err := restoreBaseConfig(); err != nil {
		return err
	}
	if readPID(PIDFile) != "" || pidOf(ctx, "dnsmasq") != "" {
		if err := restart(ctx); err != nil {
			return err
		}
	}
	fmt.Printf("  dnsmasq: конфиг %s удалён\n", ConfFile)
	return nil
}

func restoreBaseConfig() error {
	data, err := os.ReadFile(BackupFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("чтение backup %s: %w", BackupFile, err)
	}
	if err := os.WriteFile(BaseConfFile, data, 0644); err != nil {
		return fmt.Errorf("восстановление %s из %s: %w", BaseConfFile, BackupFile, err)
	}
	fmt.Printf("  dnsmasq: восстановлен %s из %s\n", BaseConfFile, BackupFile)
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
