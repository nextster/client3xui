package client3xui

import (
	"encoding/json"
)

type RealityPublicSettings struct {
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
	ServerName  string `json:"serverName"`
	SpiderX     string `json:"spiderX"`
}
type RealitySettings struct {
	Show        bool                  `json:"show"`
	Xver        int                   `json:"xver"`
	Dest        string                `json:"dest"`
	ServerNames []string              `json:"serverNames"`
	PrivateKey  string                `json:"privateKey"`
	MinClient   string                `json:"minClient"`
	MaxClient   string                `json:"maxClient"`
	MaxTimediff int                   `json:"maxTimediff"`
	ShortIds    []string              `json:"shortIds"`
	Settings    RealityPublicSettings `json:"settings"`
}

type TcpSettings struct {
	AcceptProxyProtocol bool          `json:"acceptProxyProtocol"`
	Header              HeaderSetting `json:"header"`
}

type TcpStreamSettings struct {
	Network       string      `json:"network"`
	Security      string      `json:"security"`
	ExternalProxy []string    `json:"externalProxy"`
	TcpSettings   TcpSettings `json:"tcpSettings"`
	// RealitySettings is only included in JSON if Security is "reality"
	RealitySettings *RealitySettings `json:"realitySettings,omitempty"`
}

// func (s *TcpStreamSettings) SetRealitySettings(rs RealitySettings) {
// 	s.Security = "reality"
// 	s.RealitySettings = &rs
// }

// func (s *TcpStreamSettings) GetRealitySettings() *RealitySettings {
// 	if s.Security != "reality" {
// 		return nil
// 	}
// 	return s.RealitySettings
// }

func (s *TcpStreamSettings) UnmarshalJSON(data []byte) error {
	type Alias TcpStreamSettings
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if s.Security != "reality" {
		s.RealitySettings = nil
	}
	return nil
}

// func (s *TcpStreamSettings) MarshalJSON() ([]byte, error) {
// 	type Alias TcpStreamSettings
// 	return json.Marshal(&struct {
// 		*Alias
// 	}{
// 		Alias: (*Alias)(s),
// 	})
// }

type VlessSettings struct {
	Clients    []InboundClient `json:"clients"`
	Decryption string          `json:"decryption"`
	Fallbacks  []string        `json:"fallbacks"`
}

type VmessSettings struct {
	Clients []InboundClient `json:"clients"`
}

type FallbackOptions struct {
	Name string `json:"name"`
	Alpn string `json:"alpn"`
	Path string `json:"path"`
	Dest string `json:"dest"`
	Xver int    `json:"xver"`
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

type InboundBaseSettings struct {
	Up, Down, Total, Remark, Enable, ExpiryTime, Listen, Port, Protocol string
}

type HeaderSetting struct {
	Type string `json:"type"`
}

type SniffingSettings struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
	MetadataOnly bool     `json:"metadataOnly"`
	RouteOnly    bool     `json:"routeOnly"`
}

type QuicSettings struct {
	Security string        `json:"security"`
	Key      string        `json:"key"`
	Header   HeaderSetting `json:"header"`
}

type QuicStreamSettings struct {
	Network       string       `json:"network"`
	Security      string       `json:"security"`
	ExternalProxy []string     `json:"externalProxy"`
	QuicSettings  QuicSettings `json:"quicSettings"`
}
