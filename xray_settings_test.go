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
	"encoding/json"
	"testing"
)

func TestXraySettingsMarshaling(t *testing.T) {
	// Test the Xray settings structure
	settings := &XraySettings{
		Log: &XrayLog{
			Access:      "none",
			DNSLog:      false,
			Error:       "",
			LogLevel:    "warning",
			MaskAddress: "",
		},
		API: &XrayAPI{
			Tag:      "api",
			Services: []string{"HandlerService", "LoggerService", "StatsService"},
		},
		Inbounds: []XrayInbound{
			{
				Tag:      "api",
				Listen:   "127.0.0.1",
				Port:     62789,
				Protocol: "dokodemo-door",
				Settings: map[string]interface{}{
					"address": "127.0.0.1",
				},
			},
		},
		Outbounds: []XrayOutbound{
			{
				Tag:      "direct",
				Protocol: "freedom",
				Settings: map[string]interface{}{
					"domainStrategy": "UseIP",
					"redirect":       "",
					"noises":         []interface{}{},
				},
			},
		},
		Stats: map[string]interface{}{},
	}

	// Test marshaling
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Failed to marshal XraySettings: %v", err)
	}

	// Test unmarshaling
	var unmarshaled XraySettings
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal XraySettings: %v", err)
	}

	// Basic validation
	if unmarshaled.Log.LogLevel != "warning" {
		t.Errorf("Expected log level 'warning', got '%s'", unmarshaled.Log.LogLevel)
	}
	if len(unmarshaled.API.Services) != 3 {
		t.Errorf("Expected 3 API services, got %d", len(unmarshaled.API.Services))
	}
}

func TestXraySettingsWrapperParsing(t *testing.T) {
	// Test parsing the wrapper response format
	jsonStr := `{
		"xraySetting": {
			"log": {
				"access": "none",
				"dnsLog": false,
				"error": "",
				"loglevel": "warning",
				"maskAddress": ""
			},
			"api": {
				"tag": "api",
				"services": ["HandlerService", "LoggerService", "StatsService"]
			},
			"stats": {}
		},
		"inboundTags": ["inbound-14509", "inbound-44488"]
	}`

	var wrapper XraySettingsWrapper
	err := json.Unmarshal([]byte(jsonStr), &wrapper)
	if err != nil {
		t.Fatalf("Failed to unmarshal XraySettingsWrapper: %v", err)
	}

	if wrapper.XraySetting == nil {
		t.Fatal("XraySetting is nil")
	}
	if len(wrapper.InboundTags) != 2 {
		t.Errorf("Expected 2 inbound tags, got %d", len(wrapper.InboundTags))
	}
}

func TestParseVlessURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		check   func(t *testing.T, outbound *XrayOutbound)
	}{
		{
			name: "Valid VLESS URL with Reality",
			url:  "vless://8e72473d-3c52-4153-b5ba-3b06035d0ad1@89.169.53.31:36989?type=tcp&security=reality&pbk=QpIeLuq1OYR1dSWituaXb0c8h4iZtkFPIjKxLKiyC3o&fp=random&sni=rt.com&sid=82c54a0dbca8&spx=/#test-server",
			check: func(t *testing.T, outbound *XrayOutbound) {
				if outbound.Tag != "test-server" {
					t.Errorf("Expected tag 'test-server', got '%s'", outbound.Tag)
				}
				if outbound.Protocol != "vless" {
					t.Errorf("Expected protocol 'vless', got '%s'", outbound.Protocol)
				}

				// Check vnext settings
				vnext := outbound.Settings["vnext"].([]map[string]interface{})
				if len(vnext) != 1 {
					t.Fatalf("Expected 1 vnext entry, got %d", len(vnext))
				}
				server := vnext[0]
				if server["address"] != "89.169.53.31" {
					t.Errorf("Expected address '89.169.53.31', got '%v'", server["address"])
				}
				if server["port"] != 36989 {
					t.Errorf("Expected port 36989, got '%v'", server["port"])
				}

				// Check stream settings
				if outbound.StreamSettings["network"] != "tcp" {
					t.Errorf("Expected network 'tcp', got '%v'", outbound.StreamSettings["network"])
				}
				if outbound.StreamSettings["security"] != "reality" {
					t.Errorf("Expected security 'reality', got '%v'", outbound.StreamSettings["security"])
				}

				// Check reality settings
				reality := outbound.StreamSettings["realitySettings"].(map[string]interface{})
				if reality["publicKey"] != "QpIeLuq1OYR1dSWituaXb0c8h4iZtkFPIjKxLKiyC3o" {
					t.Errorf("Unexpected public key: %v", reality["publicKey"])
				}
				if reality["serverName"] != "rt.com" {
					t.Errorf("Expected serverName 'rt.com', got '%v'", reality["serverName"])
				}
			},
		},
		{
			name: "Simple VLESS URL",
			url:  "vless://uuid-here@example.com:443",
			check: func(t *testing.T, outbound *XrayOutbound) {
				vnext := outbound.Settings["vnext"].([]map[string]interface{})
				server := vnext[0]
				if server["address"] != "example.com" {
					t.Errorf("Expected address 'example.com', got '%v'", server["address"])
				}
				if server["port"] != 443 {
					t.Errorf("Expected port 443, got '%v'", server["port"])
				}
			},
		},
		{
			name:    "Invalid scheme",
			url:     "vmess://uuid@example.com:443",
			wantErr: true,
		},
		{
			name:    "Missing UUID",
			url:     "vless://@example.com:443",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outbound, err := ParseVlessURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVlessURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, outbound)
			}
		})
	}
}

func TestCreateOutbounds(t *testing.T) {
	// Test CreateFreedomOutbound
	freedom := CreateFreedomOutbound("direct", "UseIP")
	if freedom.Tag != "direct" {
		t.Errorf("Expected tag 'direct', got '%s'", freedom.Tag)
	}
	if freedom.Protocol != "freedom" {
		t.Errorf("Expected protocol 'freedom', got '%s'", freedom.Protocol)
	}
	if freedom.Settings["domainStrategy"] != "UseIP" {
		t.Errorf("Expected domainStrategy 'UseIP', got '%v'", freedom.Settings["domainStrategy"])
	}

	// Test CreateBlackholeOutbound
	blackhole := CreateBlackholeOutbound("blocked")
	if blackhole.Tag != "blocked" {
		t.Errorf("Expected tag 'blocked', got '%s'", blackhole.Tag)
	}
	if blackhole.Protocol != "blackhole" {
		t.Errorf("Expected protocol 'blackhole', got '%s'", blackhole.Protocol)
	}

	// Test CreateVlessOutbound
	reality := &XrayRealityOutboundSettings{
		PublicKey:   "test-key",
		Fingerprint: "random",
		ServerName:  "example.com",
		ShortID:     "abc123",
		SpiderX:     "/",
	}
	vless := CreateVlessOutbound("test-vless", "1.2.3.4", 443, "uuid-here", "xtls-rprx-vision", "reality", reality)
	if vless.Tag != "test-vless" {
		t.Errorf("Expected tag 'test-vless', got '%s'", vless.Tag)
	}
	if vless.Protocol != "vless" {
		t.Errorf("Expected protocol 'vless', got '%s'", vless.Protocol)
	}

	// Check reality settings were applied
	if vless.StreamSettings["security"] != "reality" {
		t.Errorf("Expected security 'reality', got '%v'", vless.StreamSettings["security"])
	}
	realitySettings := vless.StreamSettings["realitySettings"].(map[string]interface{})
	if realitySettings["publicKey"] != "test-key" {
		t.Errorf("Expected publicKey 'test-key', got '%v'", realitySettings["publicKey"])
	}
}
