package wg

import (
	"context"
	"fmt"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/keenetic"
	"github.com/BataevDaniil/eleutherios/internal/logging"
)

const RouteTableID = 1001
const MarkNum = "0xd1000"
const RulePriority = "1778"

// Up поднимает WireGuard Keenetic через REST API. Возвращает linux-имя интерфейса.
func Up(ctx context.Context, cliName string) (string, error) {
	requestedName := cliName
	ifaces, err := getInterfaces(ctx)
	if err != nil {
		return "", fmt.Errorf("запрос к API: %w", err)
	}

	wg, cliName := findWireguard(ifaces, cliName)
	if wg == nil {
		available := listWireguardNames(ifaces)
		if len(available) == 0 {
			return "", fmt.Errorf("WireGuard %q не найден. В панели Keenetic нет ни одного WireGuard-интерфейса", requestedName)
		}
		return "", fmt.Errorf("WireGuard %q не найден. Доступные: %s", requestedName, strings.Join(available, ", "))
	}

	if wg.State == "down" {
		_, err := keenetic.Post(ctx, keenetic.APIBase+"/interface/"+cliName, `{"up":"true"}`)
		if err != nil {
			return "", fmt.Errorf("поднять %s: %w", cliName, err)
		}
	}

	entName, err := getEntwareName(ctx, cliName, wg)
	if err != nil {
		return "", fmt.Errorf("linux-интерфейс для %s: %w", cliName, err)
	}

	logging.Logger().Info("WireGuard поднят", "component", "wireguard", "description", wg.Description, "cli_name", cliName, "interface", entName)
	return entName, nil
}

// Down чистит маршруты
func Down(ctx context.Context) {
	cmds := [][]string{
		{"ip", "route", "flush", "table", fmt.Sprint(RouteTableID)},
		{"ip", "rule", "del", "fwmark", MarkNum + "/" + MarkNum, "table", fmt.Sprint(RouteTableID), "priority", RulePriority},
		{"ip", "route", "flush", "cache"},
	}
	for _, args := range cmds {
		out, err := execCmd(ctx, args[0], args[1:]...)
		if err != nil && !strings.Contains(out, "No such file") && !strings.Contains(out, "No such process") {
			logging.Logger().Warn("wg down команда провалилась", "component", "wireguard", "cmd", strings.Join(args, " "), "error", err, "output", strings.TrimSpace(out))
		}
	}
	logging.Logger().Info("WireGuard маршруты убраны", "component", "wireguard")
}

// Status возвращает список WireGuard-интерфейсов
func Status(ctx context.Context) string {
	ifaces, err := getInterfaces(ctx)
	if err != nil {
		return "ошибка API: " + err.Error()
	}
	out := ""
	for _, wg := range ifaces {
		if wg.Type != "Wireguard" {
			continue
		}
		state := "ОТКЛЮЧЕН"
		if wg.State == "up" {
			state = "ПОДКЛЮЧЕН"
		}
		out += fmt.Sprintf("%s (%s): %s [%s]\n", wg.Description, wg.ID, state, wg.Address)
	}
	if out == "" {
		return "WireGuard не найден"
	}
	return out
}

func listWireguardNames(ifaces []keenetic.Interface) []string {
	var out []string
	for _, w := range ifaces {
		if w.Type != "Wireguard" {
			continue
		}
		name := w.ID
		if w.Description != "" {
			name = w.Description + " (" + w.ID + ")"
		}
		out = append(out, name)
	}
	return out
}

func List(ctx context.Context) string {
	ifaces, err := getInterfaces(ctx)
	if err != nil {
		return "ошибка API: " + err.Error()
	}
	out := ""
	for _, wg := range ifaces {
		if wg.Type != "Wireguard" {
			continue
		}
		out += fmt.Sprintf("%s\t%s\t%s\t%s\n", wg.ID, wg.Description, wg.State, wg.InterfaceName)
	}
	if out == "" {
		return "WireGuard не найден"
	}
	return out
}
