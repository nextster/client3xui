/* Copyright 2024 İrem Kuyucu <irem@digilol.net>
 * Copyright 2024 Laurynas Četyrkinas <laurynas@digilol.net>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client3xui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// XraySettings represents the complete Xray configuration
type XraySettings struct {
	Log       *XrayLog               `json:"log,omitempty"`
	API       *XrayAPI               `json:"api,omitempty"`
	Inbounds  []XrayInbound          `json:"inbounds,omitempty"`
	Outbounds []XrayOutbound         `json:"outbounds,omitempty"`
	Policy    *XrayPolicy            `json:"policy,omitempty"`
	Routing   *XrayRouting           `json:"routing,omitempty"`
	Stats     map[string]interface{} `json:"stats,omitempty"`
}

// XrayLog represents log configuration
type XrayLog struct {
	Access      string `json:"access,omitempty"`
	DNSLog      bool   `json:"dnsLog,omitempty"`
	Error       string `json:"error,omitempty"`
	LogLevel    string `json:"loglevel,omitempty"`
	MaskAddress string `json:"maskAddress,omitempty"`
}

// XrayAPI represents API configuration
type XrayAPI struct {
	Tag      string   `json:"tag,omitempty"`
	Services []string `json:"services,omitempty"`
}

// XrayInbound represents inbound configuration
type XrayInbound struct {
	Tag      string                 `json:"tag,omitempty"`
	Listen   string                 `json:"listen,omitempty"`
	Port     int                    `json:"port,omitempty"`
	Protocol string                 `json:"protocol,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// XrayOutbound represents outbound configuration
type XrayOutbound struct {
	Tag            string                 `json:"tag,omitempty"`
	Protocol       string                 `json:"protocol,omitempty"`
	Settings       map[string]interface{} `json:"settings,omitempty"`
	StreamSettings map[string]interface{} `json:"streamSettings,omitempty"`
}

// XrayPolicy represents policy configuration
type XrayPolicy struct {
	Levels map[string]interface{} `json:"levels,omitempty"`
	System map[string]interface{} `json:"system,omitempty"`
}

// XrayRouting represents routing configuration
type XrayRouting struct {
	DomainStrategy string     `json:"domainStrategy,omitempty"`
	Rules          []XrayRule `json:"rules,omitempty"`
}

// XrayRule represents a routing rule
type XrayRule struct {
	Type        string   `json:"type"`
	InboundTag  []string `json:"inboundTag,omitempty"`
	OutboundTag string   `json:"outboundTag,omitempty"`
	IP          []string `json:"ip,omitempty"`
	Protocol    []string `json:"protocol,omitempty"`
}

// XraySettingsResponse represents the response from getting Xray settings
type XraySettingsResponse struct {
	Success bool            `json:"success"`
	Msg     string          `json:"msg"`
	Obj     json.RawMessage `json:"obj"`
}

// XraySettingsWrapper represents the wrapper for xray settings and inbound tags
type XraySettingsWrapper struct {
	XraySetting *XraySettings `json:"xraySetting"`
	InboundTags []string      `json:"inboundTags"`
}

// GetXraySettings retrieves the current Xray settings
func (c *Client) GetXraySettings(ctx context.Context) (*XraySettingsWrapper, error) {
	resp := &XraySettingsResponse{}
	err := c.Do(ctx, http.MethodPost, "/panel/xray/", nil, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("%s", resp.Msg)
	}

	// The response object is a string containing JSON, so we need to unmarshal it twice
	var objStr string
	if err := json.Unmarshal(resp.Obj, &objStr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal obj as string: %w", err)
	}

	// Now unmarshal the actual settings
	var wrapper XraySettingsWrapper
	if err := json.Unmarshal([]byte(objStr), &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return &wrapper, nil
}

// UpdateXraySettings updates the Xray settings
func (c *Client) UpdateXraySettings(ctx context.Context, settings *XraySettings) error {
	resp := &ApiResponse{}
	err := c.Do(ctx, http.MethodPost, "/panel/xray/update", settings, resp)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("%s", resp.Msg)
	}
	return nil
}

// GetXrayResult gets the Xray result (usually empty)
func (c *Client) GetXrayResult(ctx context.Context) (*ApiResponse, error) {
	resp := &ApiResponse{}
	err := c.Do(ctx, http.MethodGet, "/panel/xray/getXrayResult", nil, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return resp, fmt.Errorf("%s", resp.Msg)
	}
	return resp, nil
}

// RestartXrayService restarts the Xray service
func (c *Client) RestartXrayService(ctx context.Context) (*ApiResponse, error) {
	resp := &ApiResponse{}
	err := c.Do(ctx, http.MethodPost, "/server/restartXrayService", nil, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return resp, fmt.Errorf("%s", resp.Msg)
	}
	return resp, nil
}
