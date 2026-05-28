package dns

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/logging"
)

// Cleanup возвращает dnsmasq в исходное состояние: восстанавливает оригинальный
// конфиг из backup'а (или rollback-конфиг если backup не сохранён) и перезапускает
// dnsmasq. DNS на роутере продолжает работать так же, как до eleutherios start.
func Cleanup(ctx context.Context) error {
	if err := os.Remove(ConfFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("удаление %s: %w", ConfFile, err)
	}
	if _, err := os.Stat(InitFile); err != nil {
		logging.Logger().Info("dnsmasq не установлен, пропускаем", "component", "dnsmasq")
		return nil
	}
	if err := restoreBaseConfig(); err != nil {
		return err
	}
	if err := restart(ctx); err != nil {
		return err
	}
	logging.Logger().Info("dnsmasq перезапущен с оригинальной конфигурацией", "component", "dnsmasq", "config", BaseConfFile)
	return nil
}

func restoreBaseConfig() error {
	backup, err := os.ReadFile(BackupFile)
	if err == nil {
		if wErr := os.WriteFile(BaseConfFile, backup, 0644); wErr != nil {
			return fmt.Errorf("восстановление %s из backup: %w", BaseConfFile, wErr)
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("чтение backup %s: %w", BackupFile, err)
	}
	if wErr := os.WriteFile(BaseConfFile, []byte(rollbackConfig), 0644); wErr != nil {
		return fmt.Errorf("запись rollback %s: %w", BaseConfFile, wErr)
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
