package cmd

import (
	"fmt"

	"github.com/BataevDaniil/eleutherios/network"
	"github.com/BataevDaniil/eleutherios/wg"
	"github.com/spf13/cobra"
)

var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Команды для VPN-интерфейсов",
}

var vpnLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Показать WireGuard интерфейсы Keenetic",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(wg.List())
	},
}

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Команды для сетей",
}

var netLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Показать доступные bridge-сети",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(network.List())
	},
}

func init() {
	vpnCmd.AddCommand(vpnLsCmd)
	netCmd.AddCommand(netLsCmd)
	rootCmd.AddCommand(vpnCmd, netCmd)
}
