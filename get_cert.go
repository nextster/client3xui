package client3xui

import (
	"context"
	"fmt"
	"net/http"
)

type GetX25519CertResponse struct {
	Success bool       `json:"success"`
	Msg     string     `json:"msg"`
	Obj     X25519Cert `json:"obj"`
}

type X25519Cert struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func (c *Client) GetX25519Cert(ctx context.Context) (*X25519Cert, error) {
	resp := &GetX25519CertResponse{}

	err := c.Do(ctx, http.MethodPost, "/server/getNewX25519Cert", nil, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf(resp.Msg)
	}
	return &resp.Obj, nil
}
