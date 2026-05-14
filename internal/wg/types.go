package wg

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ifaceRecord — запись из /rci/show/interface
type ifaceRecord struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	State          string `json:"state"`
	Description    string `json:"description"`
	InterfaceName  string `json:"interface-name"`
	Address        string `json:"address"`
	DefaultGateway bool   `json:"defaultgw"`
}

func getInterfaces(ctx context.Context) ([]ifaceRecord, error) {
	resp, err := httpGet(ctx, apiBase+"/show/interface")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Keenetic RCI возвращает как map[id]record, а не массив
	var raw map[string]ifaceRecord
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("разбор JSON: %w", err)
	}

	ifaces := make([]ifaceRecord, 0, len(raw))
	for _, v := range raw {
		ifaces = append(ifaces, v)
	}
	return ifaces, nil
}

func findWireguard(ifaces []ifaceRecord, name string) (*ifaceRecord, string) {
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
