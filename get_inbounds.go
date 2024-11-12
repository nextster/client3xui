package client3xui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetInboundsResponse struct {
	Success bool      `json:"success"`
	Msg     string    `json:"msg"`
	Obj     []Inbound `json:"obj"`
}

type Inbound struct {
	ID             int                  `json:"id"`
	Up             int                  `json:"up"`
	Down           int                  `json:"down"`
	Total          int                  `json:"total"`
	Remark         string               `json:"remark"`
	Enable         bool                 `json:"enable"`
	ExpiryTime     int                  `json:"expiryTime"`
	ClientStats    []ClientStat         `json:"clientStats"`
	Listen         string               `json:"listen"`
	Port           int                  `json:"port"`
	Protocol       string               `json:"protocol"`
	Settings       InboundSettings      `json:"settings"`
	StreamSettings InboundStreamSetting `json:"streamSettings"`
	Tag            string               `json:"tag"`
	Sniffing       string               `json:"sniffing"`
}

func (i *Inbound) UnmarshalJSON(data []byte) error {
	type Alias Inbound
	aux := &struct {
		Settings       string `json:"settings"`
		StreamSettings string `json:"streamSettings"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.StreamSettings != "" {
		if err := json.Unmarshal([]byte(aux.StreamSettings), &i.StreamSettings); err != nil {
			return err
		}
	}
	if aux.Settings != "" {
		if err := json.Unmarshal([]byte(aux.Settings), &i.Settings); err != nil {
			return err
		}
	}
	return nil
}

type ClientStat struct {
	ID         int    `json:"id"`
	InboundID  int    `json:"inboundId"`
	Enable     bool   `json:"enable"`
	Email      string `json:"email"`
	Up         int    `json:"up"`
	Down       int    `json:"down"`
	ExpiryTime int    `json:"expiryTime"`
	Total      int    `json:"total"`
	Reset      int    `json:"reset"`
}

func (c *Client) GetInbounds(ctx context.Context) (*GetInboundsResponse, error) {
	resp := &GetInboundsResponse{}
	err := c.Do(ctx, http.MethodPost, "/panel/inbound/list", nil, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return resp, fmt.Errorf(resp.Msg)
	}
	return resp, err
}
