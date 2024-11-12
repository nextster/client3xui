package client3xui

import (
	"encoding/json"
)

type RealitySettings struct {
	Show        bool     `json:"show"`
	Xver        int      `json:"xver"`
	Dest        string   `json:"dest"`
	ServerNames []string `json:"serverNames"`
	PrivateKey  string   `json:"privateKey"`
	MinClient   string   `json:"minClient"`
	MaxClient   string   `json:"maxClient"`
	MaxTimediff int      `json:"maxTimediff"`
	ShortIds    []string `json:"shortIds"`
	Settings    struct {
		PublicKey   string `json:"publicKey"`
		Fingerprint string `json:"fingerprint"`
		ServerName  string `json:"serverName"`
		SpiderX     string `json:"spiderX"`
	} `json:"settings"`
}

type TcpHeader struct {
	Type string `json:"type"`
}

type TcpSettings struct {
	AcceptProxyProtocol bool      `json:"acceptProxyProtocol"`
	Header              TcpHeader `json:"header"`
}

type InboundStreamSetting struct {
	Network       string        `json:"network"`
	Security      string        `json:"security"`
	ExternalProxy []interface{} `json:"externalProxy"`
	TcpSettings   TcpSettings   `json:"tcpSettings"`
	// RealitySettings is only included in JSON if Security is "reality"
	RealitySettings *RealitySettings `json:"realitySettings,omitempty"`
}

// func (iss *InboundStreamSetting) SetRealitySettings(rs RealitySettings) {
// 	if iss.Security == "reality" {
// 		iss.RealitySettings = &rs
// 	} else {
// 		iss.RealitySettings = nil
// 	}
// }

func (iss *InboundStreamSetting) UnmarshalJSON(data []byte) error {
	type Alias InboundStreamSetting
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(iss),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if iss.Security != "reality" {
		iss.RealitySettings = nil
	}
	return nil
}

type InboundSettings struct {
	Clients    []InboundClient `json:"clients"`
	Decryption string          `json:"decryption"`
	Fallbacks  []string        `json:"fallbacks"`
}

type InboundClient struct {
	Email      string `json:"email"`
	Enable     bool   `json:"enable"`
	ExpiryTime int    `json:"expiryTime"`
	Flow       string `json:"flow,omitempty"`
	ID         string `json:"id"`
	LimitIp    int    `json:"limitIp"`
	Reset      int    `json:"reset"`
	SubId      string `json:"subId,omitempty"`
	// TgId       string `json:"tgId,omitempty"`
	TotalGB int `json:"totalGB"`
}

type Fallback struct {
	// Define the fields for Fallback if any
}
