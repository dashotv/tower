package plex

import (
	"fmt"

	"github.com/pkg/errors"
)

func (p *Client) GetClients() (map[string]any, error) {
	clients := map[string]any{}
	resp, err := p._server().SetResult(clients).
		SetHeaders(p.Headers).Get("/clients")
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get clients: %s", resp.Status())
	}
	fmt.Printf("clients: %s\n", resp.String())
	return clients, nil
}

func (p *Client) GetDevices() (map[string]any, error) {
	devices := map[string]any{}
	resp, err := p._server().SetResult(devices).
		SetHeaders(p.Headers).Get("/devices")
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get devices: %s", resp.Status())
	}
	fmt.Printf("devices: %s\n", resp.String())
	return devices, nil
}
