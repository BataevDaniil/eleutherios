package eleutherios

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/BataevDaniil/eleutherios/internal/boot"
	iptables2 "github.com/BataevDaniil/eleutherios/internal/iptables"
	"github.com/BataevDaniil/eleutherios/internal/logging"
	"github.com/spf13/cobra"
)

// AnnotationRequiresRoot помечает команду как требующую root.
// PersistentPreRunE падает с понятной ошибкой до запуска RunE.
const AnnotationRequiresRoot = "requires-root"

var iptablesHook bool
var fsHook bool
var hookNet string

func SetVersion(v string) {
	rootCmd.Version = v
}

var rootCmd = &cobra.Command{
	Use:   "eleutherios",
	Short: "WireGuard split-tunnel: *.ru → ISP, остальное → WG",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("log-file")
		if path == "" {
			path, _ = cmd.Root().PersistentFlags().GetString("log-file")
		}
		if err := logging.Configure(path); err != nil {
			return fmt.Errorf("log file: %w", err)
		}
		if needsRoot(cmd) && os.Geteuid() != 0 {
			return fmt.Errorf("команда %q требует root (запустите от root или через sudo)", cmd.CommandPath())
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if fsHook {
			logging.Logger().Info("Запуск fs hook", "component", "hook", "hook", "fs")
			if err := boot.RunFSHook(cmd.Context()); err != nil {
				return fmt.Errorf("fs hook: %w", err)
			}
			logging.Logger().Info("fs hook завершён", "component", "hook", "hook", "fs")
			return nil
		}
		if !iptablesHook {
			return cmd.Help()
		}
		logging.Logger().Info("Запуск iptables hook", "component", "hook", "hook", "iptables", "net", hookNet)
		iface, err := iptables2.NetIface(cmd.Context(), hookNet)
		if err != nil {
			return err
		}
		if err := iptables2.RunHook(cmd.Context(), iface); err != nil {
			return fmt.Errorf("iptables hook: %w", err)
		}
		logging.Logger().Info("iptables hook завершён", "component", "hook", "hook", "iptables", "net", hookNet, "iface", iface)
		return nil
	},
}

// defaultTimeout — общий таймаут всей команды. На медленном роутере с холодным
// dnsmasq и большим iptables-save 1 минуты прежнего лимита было тесно.
// Перекрывается флагом --timeout.
const defaultTimeout = 3 * time.Minute

func Execute() error {
	// Парсим persistent-флаги вручную, чтобы вытащить --timeout до создания
	// контекста. cobra потом распарсит их же повторно внутри ExecuteContext.
	_ = rootCmd.ParseFlags(os.Args[1:])
	timeout, _ := rootCmd.PersistentFlags().GetDuration("timeout")
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		logging.Logger().Error("command failed", "error", err)
	}
	if closeErr := logging.Close(); err == nil && closeErr != nil {
		return fmt.Errorf("close log file: %w", closeErr)
	}
	return err
}

// needsRoot решает, нужен ли root для текущей команды.
// Учитывает Annotations и скрытые hook-флаги rootCmd (--iptables-hook/--fs-hook),
// потому что rootCmd сам по себе аннотацию не имеет: hook'и вызываются как `eleutherios --iptables-hook ...`.
func needsRoot(cmd *cobra.Command) bool {
	if cmd.Annotations[AnnotationRequiresRoot] == "true" {
		return true
	}
	// rootCmd сам не имеет родителя; hook-режим запускается как `eleutherios --iptables-hook ...`.
	if cmd.Parent() == nil && (iptablesHook || fsHook) {
		return true
	}
	return false
}

func init() {
	rootCmd.PersistentFlags().String("log-file", "", "путь к файлу логов")
	rootCmd.PersistentFlags().Duration("timeout", defaultTimeout, "общий таймаут команды (напр. 30s, 5m)")
	rootCmd.Flags().BoolVar(&iptablesHook, "iptables-hook", false, "восстановить iptables из NDM hook")
	rootCmd.Flags().BoolVar(&fsHook, "fs-hook", false, "создать ipset из NDM fs hook")
	rootCmd.Flags().StringVar(&hookNet, "net", "br0", "сеть для iptables hook")
	mustHide(rootCmd, "iptables-hook", "fs-hook", "net")
}

// mustHide вызывает MarkHidden и паникует при ошибке — флаг гарантированно объявлен выше.
func mustHide(cmd *cobra.Command, names ...string) {
	for _, n := range names {
		if err := cmd.Flags().MarkHidden(n); err != nil {
			panic(fmt.Sprintf("MarkHidden(%q): %v", n, err))
		}
	}
}
