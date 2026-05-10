package cmd

import (
	"fmt"

	"github.com/BataevDaniil/eleutherios/internal/boot"
	iptables2 "github.com/BataevDaniil/eleutherios/internal/iptables"
	"github.com/spf13/cobra"
)

var iptablesHook bool
var fsHook bool
var hookNet string

var rootCmd = &cobra.Command{
	Use:   "eleutherios",
	Short: "WireGuard split-tunnel: *.ru → ISP, остальное → WG",
	RunE: func(cmd *cobra.Command, args []string) error {
		if fsHook {
			if err := boot.RunFSHook(); err != nil {
				return fmt.Errorf("fs hook: %w", err)
			}
			return nil
		}
		if !iptablesHook {
			return cmd.Help()
		}
		iface, err := iptables2.NetIface(hookNet)
		if err != nil {
			return err
		}
		if err := iptables2.RunHook(iface); err != nil {
			return fmt.Errorf("iptables hook: %w", err)
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVar(&iptablesHook, "iptables-hook", false, "восстановить iptables из NDM hook")
	rootCmd.Flags().BoolVar(&fsHook, "fs-hook", false, "создать ipset из NDM fs hook")
	rootCmd.Flags().StringVar(&hookNet, "net", "br0", "сеть для iptables hook")
	rootCmd.Flags().MarkHidden("iptables-hook")
	rootCmd.Flags().MarkHidden("fs-hook")
	rootCmd.Flags().MarkHidden("net")
}
