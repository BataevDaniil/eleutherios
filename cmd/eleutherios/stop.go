package eleutherios

import (
	"fmt"

	"github.com/BataevDaniil/eleutherios/internal/boot"
	"github.com/BataevDaniil/eleutherios/internal/dns"
	"github.com/BataevDaniil/eleutherios/internal/ipset"
	"github.com/BataevDaniil/eleutherios/internal/iptables"
	"github.com/BataevDaniil/eleutherios/internal/wg"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Остановить обход, вернуть всё как было",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		fmt.Println("[1/5] Убираем dnsmasq...")
		if err := dns.Cleanup(ctx); err != nil {
			return fmt.Errorf("dnsmasq cleanup: %w", err)
		}

		fmt.Println("[2/5] Убираем iptables...")
		iptables.Cleanup(ctx)

		fmt.Println("[3/5] Убираем ipset...")
		ipset.DestroySets(ctx)

		fmt.Println("[4/5] Убираем маршруты WireGuard...")
		wg.Down(ctx)

		fmt.Println("[5/5] Убираем автозапуск...")
		boot.Remove()

		fmt.Println("Обход остановлен.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
