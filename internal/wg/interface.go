package wg

import "fmt"

const apiBase = "http://127.0.0.1:79/rci"
const RouteTableID = 1001
const MarkNum = "0xd1000"
const RulePriority = "1778"

// Up поднимает WireGuard Keenetic через REST API. Возвращает linux-имя интерфейса.
func Up(cliName string) (string, error) {
	requestedName := cliName
	ifaces, err := getInterfaces()
	if err != nil {
		return "", fmt.Errorf("запрос к API: %w", err)
	}

	wg, cliName := findWireguard(ifaces, cliName)
	if wg == nil {
		return "", fmt.Errorf("WireGuard %q не найден. Проверьте имя интерфейса в панели Keenetic", requestedName)
	}

	if wg.State == "down" {
		_, err := httpPost(apiBase+"/interface/"+cliName, `{"up":"true"}`)
		if err != nil {
			return "", fmt.Errorf("поднять %s: %w", cliName, err)
		}
	}

	entName, err := getEntwareName(cliName, wg)
	if err != nil {
		return "", fmt.Errorf("linux-интерфейс для %s: %w", cliName, err)
	}

	fmt.Printf("  WireGuard: %s (%s) → %s\n", wg.Description, cliName, entName)
	return entName, nil
}

// Down чистит маршруты
func Down() error {
	execCmd("ip", "route", "flush", "table", fmt.Sprint(RouteTableID))
	execCmd("ip", "rule", "del", "fwmark", MarkNum+"/"+MarkNum, "table", fmt.Sprint(RouteTableID), "priority", RulePriority)
	execCmd("ip", "route", "flush", "cache")
	return nil
}

// Status возвращает список WireGuard-интерфейсов
func Status() string {
	ifaces, err := getInterfaces()
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

func List() string {
	ifaces, err := getInterfaces()
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
