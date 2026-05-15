package dns

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/fsutil"
	"github.com/BataevDaniil/eleutherios/internal/logging"
)

const (
	ConfFile     = "/opt/etc/dnsmasq.d/eleutherios.dnsmasq"
	BaseConfFile = "/opt/etc/dnsmasq.conf"
	BackupFile   = "/opt/etc/dnsmasq.conf.backup"
	InitFile     = "/opt/etc/init.d/S56dnsmasq"
	PIDFile      = "/var/run/opt-dnsmasq.pid"
	Port         = "9753"
)

//go:embed dnsmasq.conf
var baseConfigTemplate string

//go:embed overlay.dnsmasq
var overlayConfig string

func Configure(ctx context.Context) error {
	if err := os.MkdirAll("/opt/etc/dnsmasq.d", 0755); err != nil {
		return fmt.Errorf("создание dnsmasq.d: %w", err)
	}
	if err := fsutil.WriteAtomic(ConfFile, []byte(renderOverlay()), 0644); err != nil {
		return fmt.Errorf("запись %s: %w", ConfFile, err)
	}
	if err := ensureBaseConfig(); err != nil {
		return err
	}
	if _, err := ensureRunning(ctx); err != nil {
		return err
	}
	if err := restart(ctx); err != nil {
		return err
	}
	logging.Logger().Info("dnsmasq настроен", "component", "dnsmasq", "domain", "*.ru", "ipset", "ELEUTHERIOS_RU", "port", Port, "config", ConfFile)
	return nil
}

// renderOverlay возвращает содержимое /opt/etc/dnsmasq.d/eleutherios.dnsmasq.
// Эта политика — "домены *.ru попадают в ipset ELEUTHERIOS_RU при резолве".
// Список доменов задаётся в overlay.dnsmasq.
func renderOverlay() string {
	return strings.TrimRight(overlayConfig, "\n") + "\n"
}

// renderBaseConfig подставляет порт в встроенный шаблон dnsmasq.conf
// и нормализует завершающие переводы строк.
func renderBaseConfig(port string) string {
	text := strings.ReplaceAll(baseConfigTemplate, "@PORT", port)
	return strings.TrimRight(text, "\n") + "\n"
}

func ensureBaseConfig() error {
	if err := backupBaseConfig(); err != nil {
		return err
	}
	if err := fsutil.WriteAtomic(BaseConfFile, []byte(renderBaseConfig(Port)), 0644); err != nil {
		return fmt.Errorf("запись %s: %w", BaseConfFile, err)
	}
	return nil
}

func backupBaseConfig() error {
	if _, err := os.Stat(BackupFile); err == nil {
		return nil
	}
	data, err := os.ReadFile(BaseConfFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("чтение %s для backup: %w", BaseConfFile, err)
	}
	if err := fsutil.WriteAtomic(BackupFile, data, 0644); err != nil {
		return fmt.Errorf("запись backup %s: %w", BackupFile, err)
	}
	return nil
}
