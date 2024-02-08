package plex

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientsResponse struct {
	MediaContainer struct {
		Size   int64           `json:"size"`
		Server []*ServerClient `json:"Server"`
	} `json:"MediaContainer"`
}

type ServerClient struct {
	Name                 string `json:"name"`
	Host                 string `json:"host"`
	Address              string `json:"address"`
	Port                 int64  `json:"port"`
	MachineIdentifier    string `json:"machineIdentifier"`
	Version              string `json:"version"`
	Protocol             string `json:"protocol"`
	Product              string `json:"product"`
	DeviceClass          string `json:"deviceClass"`
	ProtocolVersion      string `json:"protocolVersion"`
	ProtocolCapabilities string `json:"protocolCapabilities"`
}

func (p *Client) GetClients() (*ClientsResponse, error) {
	clients := &ClientsResponse{}
	resp, err := p._server().
		SetResult(clients).
		SetHeaders(p.Headers).
		Get("/clients")
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
