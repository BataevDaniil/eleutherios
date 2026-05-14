// Package keenetic — клиент к локальному RCI API роутера Keenetic.
// Используется и для WireGuard-интерфейсов, и для bridge-сетей.
package keenetic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIBase — корень RCI на роутере. Keenetic поднимает HTTP-сервис на localhost:79.
const APIBase = "http://127.0.0.1:79/rci"

// Interface — запись из /rci/show/interface. Объединяет поля, нужные
// разным потребителям (wg использует State/DefaultGateway, network — Index).
type Interface struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	State          string `json:"state"`
	Description    string `json:"description"`
	InterfaceName  string `json:"interface-name"`
	Address        string `json:"address"`
	Index          int    `json:"index"`
	DefaultGateway bool   `json:"defaultgw"`
}

// ShowInterfaces возвращает список интерфейсов из /rci/show/interface.
// Порядок не определён — RCI отдаёт map[id]record.
func ShowInterfaces(ctx context.Context) ([]Interface, error) {
	resp, err := Get(ctx, APIBase+"/show/interface")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw map[string]Interface
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("разбор JSON: %w", err)
	}
	out := make([]Interface, 0, len(raw))
	for _, v := range raw {
		out = append(out, v)
	}
	return out, nil
}

// Get делает GET-запрос с контекстом к RCI. Экспортируется, чтобы другие пакеты
// (wg для POST на /interface/<name>) могли использовать общую инфраструктуру.
func Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// Post делает POST с JSON-телом к RCI.
func Post(ctx context.Context, url, body string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, io.NopCloser(strings.NewReader(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}
