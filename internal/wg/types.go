package wg

import (
	"context"
	"strings"

	"github.com/BataevDaniil/eleutherios/internal/keenetic"
)

func getInterfaces(ctx context.Context) ([]keenetic.Interface, error) {
	return keenetic.ShowInterfaces(ctx)
}

func findWireguard(ifaces []keenetic.Interface, name string) (*keenetic.Interface, string) {
	name = strings.TrimSpace(name)
	for i := range ifaces {
		if ifaces[i].Type == "Wireguard" && (ifaces[i].ID == name || ifaces[i].Description == name) {
			return &ifaces[i], ifaces[i].ID
		}
	}
	if name != "" {
		return nil, ""
	}
	for i := range ifaces {
		if ifaces[i].Type == "Wireguard" {
			return &ifaces[i], ifaces[i].ID
		}
	}
	return nil, ""
}
