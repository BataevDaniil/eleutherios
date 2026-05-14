package eleutherios

import (
	"context"
	"fmt"
	"time"

	"github.com/BataevDaniil/eleutherios/internal/boot"
	iptables2 "github.com/BataevDaniil/eleutherios/internal/iptables"
	"github.com/BataevDaniil/eleutherios/internal/logging"
	"github.com/spf13/cobra"
)

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

func Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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

func init() {
	rootCmd.PersistentFlags().String("log-file", "", "путь к файлу логов")
	rootCmd.Flags().BoolVar(&iptablesHook, "iptables-hook", false, "восстановить iptables из NDM hook")
	rootCmd.Flags().BoolVar(&fsHook, "fs-hook", false, "создать ipset из NDM fs hook")
	rootCmd.Flags().StringVar(&hookNet, "net", "br0", "сеть для iptables hook")
	rootCmd.Flags().MarkHidden("iptables-hook")
	rootCmd.Flags().MarkHidden("fs-hook")
	rootCmd.Flags().MarkHidden("net")
}
