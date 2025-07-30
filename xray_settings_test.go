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
