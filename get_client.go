package client3xui

import (
	"context"
	"fmt"
	"net/http"
)

type GetClientResponse struct {
	Success bool       `json:"success"`
	Msg     string     `json:"msg"`
	Obj     ClientStat `json:"obj"`
}

func (c *Client) GetClientByEmail(ctx context.Context, email string) (*ClientStat, error) {
	resp := &GetClientResponse{}

	err := c.Do(ctx, http.MethodGet, "/panel/api/inbounds/getClientTraffics/"+email, nil, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf(resp.Msg)
	}
	return &resp.Obj, nil
}
