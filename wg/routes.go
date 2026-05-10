package wg

import (
	"fmt"
	"regexp"
	"strings"
)

var ipRegex = regexp.MustCompile(`inet ([0-9.]+)/`)

// AddRoutes настраивает маршруты: отдельная таблица для WG
func AddRoutes(name string) error {
	// получаем IP интерфейса WG
	out, err := execCmd("ip", "addr", "show", name)
	if err != nil {
		return fmt.Errorf("ip addr show %s: %w", name, err)
	}

	matches := ipRegex.FindStringSubmatch(out)
	if len(matches) < 2 {
		return fmt.Errorf("не найден IP для интерфейса %s", name)
	}
	wgIP := matches[1]

	// маршрут по умолчанию через WG в отдельной таблице
	cmds := [][]string{
		{"ip", "route", "replace", "default", "dev", name, "src", wgIP, "table", fmt.Sprint(RouteTableID)},
		{"ip", "rule", "add", "fwmark", MarkNum + "/" + MarkNum, "table", fmt.Sprint(RouteTableID), "priority", RulePriority},
	}

	for _, args := range cmds {
		out, err := execCmd(args[0], args[1:]...)
		if err != nil && !strings.Contains(out, "File exists") {
			return fmt.Errorf("%s: %w (%s)", args[0], err, out)
		}
	}
	_, _ = execCmd("ip", "route", "flush", "cache")
	return nil
}
