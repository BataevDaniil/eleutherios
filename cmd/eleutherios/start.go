package eleutherios

import (
	"fmt"

	"github.com/BataevDaniil/eleutherios/internal/boot"
	"github.com/BataevDaniil/eleutherios/internal/dns"
	"github.com/BataevDaniil/eleutherios/internal/ipset"
	"github.com/BataevDaniil/eleutherios/internal/iptables"
	wg2 "github.com/BataevDaniil/eleutherios/internal/wg"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Запустить обход: *.ru → ISP, всё остальное → WireGuard",
	RunE: func(cmd *cobra.Command, args []string) error {
		wgCli, _ := cmd.Flags().GetString("wg")
		netName, _ := cmd.Flags().GetString("net")

		fmt.Println("[1/6] Поднимаем WireGuard...")
		entName, err := wg2.Up(wgCli)
		if err != nil {
			return fmt.Errorf("wg up: %w", err)
		}

		fmt.Println("[2/6] Создаём ipset...")
		if err := ipset.CreateSets(); err != nil {
			return fmt.Errorf("ipset: %w", err)
		}

		fmt.Println("[3/6] Настраиваем dnsmasq...")
		if err := dns.Configure(); err != nil {
			return fmt.Errorf("dnsmasq: %w", err)
		}

		fmt.Println("[4/6] Настраиваем iptables...")
		if err := iptables.Setup(netName, entName); err != nil {
			return fmt.Errorf("iptables: %w", err)
		}

		fmt.Println("[5/6] Добавляем маршруты WireGuard...")
		if err := wg2.AddRoutes(entName); err != nil {
			return fmt.Errorf("wg routes: %w", err)
		}

		fmt.Println("[6/6] Настраиваем автозапуск...")
		if err := boot.Install(wgCli, netName); err != nil {
			return fmt.Errorf("boot: %w", err)
		}

		fmt.Println("ГОТОВО: *.ru → ISP, остальное → WireGuard")
		return nil
	},
}

func init() {
	startCmd.Flags().String("wg", "", "Keenetic-имя WireGuard (напр. Wireguard0)")
	startCmd.Flags().String("net", "br0", "Сеть для обхода (br0, Guest и т.д.)")
	startCmd.MarkFlagRequired("wg")
	rootCmd.AddCommand(startCmd)
}
