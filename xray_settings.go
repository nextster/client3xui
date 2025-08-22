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
	"net/url"
	"strconv"
	"strings"
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

// OutboundSettings represents different types of outbound settings
type OutboundSettings interface{}

// FreedomSettings represents freedom protocol settings
type FreedomSettings struct {
	DomainStrategy string   `json:"domainStrategy,omitempty"`
	Redirect       string   `json:"redirect,omitempty"`
	Noises         []string `json:"noises,omitempty"`
}

// BlackholeSettings represents blackhole protocol settings
type BlackholeSettings struct{}

// XrayVlessSettings represents VLESS protocol settings for outbound
type XrayVlessSettings struct {
	Vnext []VlessServerConfig `json:"vnext"`
}

// VlessServerConfig represents VLESS server configuration
type VlessServerConfig struct {
	Address string      `json:"address"`
	Port    int         `json:"port"`
	Users   []VlessUser `json:"users"`
}

// VlessUser represents a VLESS user
type VlessUser struct {
	ID         string `json:"id"`
	Flow       string `json:"flow,omitempty"`
	Encryption string `json:"encryption"`
}

// StreamSettings represents stream settings
type StreamSettings struct {
	Network         string                       `json:"network"`
	Security        string                       `json:"security"`
	RealitySettings *XrayRealityOutboundSettings `json:"realitySettings,omitempty"`
	TcpSettings     *XrayTcpSettings             `json:"tcpSettings,omitempty"`
}

// XrayRealityOutboundSettings represents reality security settings for outbound
type XrayRealityOutboundSettings struct {
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
	ServerName  string `json:"serverName"`
	ShortID     string `json:"shortId"`
	SpiderX     string `json:"spiderX"`
}

// XrayTcpSettings represents TCP transport settings for outbound
type XrayTcpSettings struct {
	Header *HeaderSettings `json:"header,omitempty"`
}

// HeaderSettings represents header settings
type HeaderSettings struct {
	Type string `json:"type"`
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
	form := url.Values{}
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}
	form.Add("xraySetting", string(settingsJSON))
	err = c.DoForm(ctx, http.MethodPost, "/panel/xray/update", form, resp)
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

// ParseVlessURL parses a VLESS URL string into an XrayOutbound
func ParseVlessURL(vlessURL string) (*XrayOutbound, error) {
	// Parse the URL
	u, err := url.Parse(vlessURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Check protocol
	if u.Scheme != "vless" {
		return nil, fmt.Errorf("invalid scheme: expected 'vless', got '%s'", u.Scheme)
	}

	// Extract UUID (user info)
	uuid := u.User.Username()
	if uuid == "" {
		return nil, fmt.Errorf("missing UUID in VLESS URL")
	}

	// Extract host and port
	host := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "443" // Default port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	// Parse query parameters
	params := u.Query()

	// Create outbound structure
	outbound := &XrayOutbound{
		Tag:      u.Fragment, // Use fragment as tag
		Protocol: "vless",
		Settings: map[string]interface{}{
			"vnext": []map[string]interface{}{
				{
					"address": host,
					"port":    port,
					"users": []map[string]interface{}{
						{
							"id":         uuid,
							"flow":       params.Get("flow"),
							"encryption": "none",
						},
					},
				},
			},
		},
	}

	// Build stream settings
	streamSettings := map[string]interface{}{
		"network": params.Get("type"),
	}

	// Handle security settings
	security := params.Get("security")
	if security != "" {
		streamSettings["security"] = security

		switch security {
		case "reality":
			realitySettings := map[string]interface{}{
				"publicKey":   params.Get("pbk"),
				"fingerprint": params.Get("fp"),
				"serverName":  params.Get("sni"),
				"shortId":     params.Get("sid"),
				"spiderX":     params.Get("spx"),
			}
			streamSettings["realitySettings"] = realitySettings
		case "tls":
			tlsSettings := map[string]interface{}{
				"serverName": params.Get("sni"),
			}
			if alpn := params.Get("alpn"); alpn != "" {
				tlsSettings["alpn"] = strings.Split(alpn, ",")
			}
			streamSettings["tlsSettings"] = tlsSettings
		}
	}

	// Handle TCP settings
	if params.Get("type") == "tcp" {
		headerType := params.Get("headerType")
		if headerType == "" {
			headerType = "none"
		}
		tcpSettings := map[string]interface{}{
			"header": map[string]interface{}{
				"type": headerType,
			},
		}
		if headerType == "http" {
			if host := params.Get("host"); host != "" {
				tcpSettings["header"].(map[string]interface{})["request"] = map[string]interface{}{
					"headers": map[string]interface{}{
						"Host": strings.Split(host, ","),
					},
				}
			}
		}
		streamSettings["tcpSettings"] = tcpSettings
	}

	outbound.StreamSettings = streamSettings

	return outbound, nil
}

// CreateFreedomOutbound creates a freedom outbound
func CreateFreedomOutbound(tag string, domainStrategy string) *XrayOutbound {
	return &XrayOutbound{
		Tag:      tag,
		Protocol: "freedom",
		Settings: map[string]interface{}{
			"domainStrategy": domainStrategy,
			"redirect":       "",
			"noises":         []interface{}{},
		},
	}
}

// CreateBlackholeOutbound creates a blackhole outbound
func CreateBlackholeOutbound(tag string) *XrayOutbound {
	return &XrayOutbound{
		Tag:      tag,
		Protocol: "blackhole",
		Settings: map[string]interface{}{},
	}
}

// CreateVlessOutbound creates a VLESS outbound with detailed settings
func CreateVlessOutbound(tag string, address string, port int, uuid string, flow string, security string, realitySettings *XrayRealityOutboundSettings) *XrayOutbound {
	outbound := &XrayOutbound{
		Tag:      tag,
		Protocol: "vless",
		Settings: map[string]interface{}{
			"vnext": []map[string]interface{}{
				{
					"address": address,
					"port":    port,
					"users": []map[string]interface{}{
						{
							"id":         uuid,
							"flow":       flow,
							"encryption": "none",
						},
					},
				},
			},
		},
		StreamSettings: map[string]interface{}{
			"network":  "tcp",
			"security": security,
			"tcpSettings": map[string]interface{}{
				"header": map[string]interface{}{
					"type": "none",
				},
			},
		},
	}

	// Add reality settings if provided
	if security == "reality" && realitySettings != nil {
		outbound.StreamSettings["realitySettings"] = map[string]interface{}{
			"publicKey":   realitySettings.PublicKey,
			"fingerprint": realitySettings.Fingerprint,
			"serverName":  realitySettings.ServerName,
			"shortId":     realitySettings.ShortID,
			"spiderX":     realitySettings.SpiderX,
		}
	}

	return outbound
}
