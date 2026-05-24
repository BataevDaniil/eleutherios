package eleutherios

import (
	"github.com/BataevDaniil/eleutherios/internal/boot"
	"github.com/BataevDaniil/eleutherios/internal/dns"
	"github.com/BataevDaniil/eleutherios/internal/ipset"
	"github.com/BataevDaniil/eleutherios/internal/iptables"
	"github.com/BataevDaniil/eleutherios/internal/logging"
	"github.com/BataevDaniil/eleutherios/internal/wg"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:         "stop",
	Short:       "Остановить обход, вернуть всё как было",
	Annotations: map[string]string{AnnotationRequiresRoot: "true"},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := logging.Logger()

		// Порядок важен: iptables должен сняться раньше, чем мы трогаем dnsmasq.
		// Иначе пока dnsmasq рестартится на 53, активный DNAT всё ещё гонит
		// :53 → 127.0.0.1:9753, где уже никто не слушает — клиенты теряют DNS.
		logger.Info("Убираем iptables", "step", "1/5", "component", "iptables")
		iptables.Cleanup(ctx)

		logger.Info("Убираем маршруты WireGuard", "step", "2/5", "component", "wireguard")
		wg.Down(ctx)

		logger.Info("Убираем dnsmasq", "step", "3/5", "component", "dnsmasq")
		if err := dns.Cleanup(ctx); err != nil {
			logger.Warn("dnsmasq cleanup провалился", "error", err)
		}

		logger.Info("Убираем ipset", "step", "4/5", "component", "ipset")
		ipset.DestroySets(ctx)

		logger.Info("Убираем автозапуск", "step", "5/5", "component", "boot")
		boot.Remove()

		logger.Info("Обход остановлен")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
