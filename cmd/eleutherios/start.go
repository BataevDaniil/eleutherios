package eleutherios

import (
	"fmt"

	"github.com/BataevDaniil/eleutherios/internal/boot"
	"github.com/BataevDaniil/eleutherios/internal/dns"
	"github.com/BataevDaniil/eleutherios/internal/ipset"
	"github.com/BataevDaniil/eleutherios/internal/iptables"
	"github.com/BataevDaniil/eleutherios/internal/logging"
	wg2 "github.com/BataevDaniil/eleutherios/internal/wg"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Запустить обход: *.ru → ISP, всё остальное → WireGuard",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		wgCli, _ := cmd.Flags().GetString("wg")
		netName, _ := cmd.Flags().GetString("net")
		logger := logging.Logger()

		logger.Info("Поднимаем WireGuard", "step", "1/6", "component", "wireguard")
		entName, err := wg2.Up(ctx, wgCli)
		if err != nil {
			return fmt.Errorf("wg up: %w", err)
		}

		logger.Info("Создаём ipset", "step", "2/6", "component", "ipset")
		if err := ipset.CreateSets(ctx); err != nil {
			return fmt.Errorf("ipset: %w", err)
		}

		logger.Info("Настраиваем dnsmasq", "step", "3/6", "component", "dnsmasq")
		if err := dns.Configure(ctx); err != nil {
			return fmt.Errorf("dnsmasq: %w", err)
		}

		logger.Info("Настраиваем iptables", "step", "4/6", "component", "iptables")
		if err := iptables.Setup(ctx, netName, entName); err != nil {
			return fmt.Errorf("iptables: %w", err)
		}

		logger.Info("Добавляем маршруты WireGuard", "step", "5/6", "component", "wireguard")
		if err := wg2.AddRoutes(ctx, entName); err != nil {
			return fmt.Errorf("wg routes: %w", err)
		}

		logger.Info("Настраиваем автозапуск", "step", "6/6", "component", "boot")
		if err := boot.Install(wgCli, netName); err != nil {
			return fmt.Errorf("boot: %w", err)
		}

		logger.Info("ГОТОВО: *.ru → ISP, остальное → WireGuard")
		return nil
	},
}

func init() {
	startCmd.Flags().String("wg", "", "Keenetic-имя WireGuard (напр. Wireguard0)")
	startCmd.Flags().String("net", "br0", "Сеть для обхода (br0, Guest и т.д.)")
	startCmd.MarkFlagRequired("wg")
	rootCmd.AddCommand(startCmd)
}
