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
)

// Ugly function signature due to a limitation in Go, this function cannot be a method of *Client.
func AddInbound[T VlessSettings | VmessSettings, K TcpStreamSettings | QuicStreamSettings](ctx context.Context, c *Client, inOpt InboundBaseSettings, protoOpt T, streamOpt K, sniffOpt SniffingSettings) (*ApiResponse, error) {
	form := url.Values{}

	protoSettings, err := json.Marshal(protoOpt)
	if err != nil {
		return nil, err
	}
	form.Add("settings", string(protoSettings))

	streamSettings, err := json.Marshal(streamOpt)
	if err != nil {
		return nil, err
	}
	form.Add("streamSettings", string(streamSettings))

	sniffingSettings, err := json.Marshal(sniffOpt)
	if err != nil {
		return nil, err
	}
	form.Add("sniffing", string(sniffingSettings))

	form.Add("up", inOpt.Up)
	form.Add("down", inOpt.Down)
	form.Add("total", inOpt.Total)
	form.Add("remark", inOpt.Remark)
	form.Add("enable", inOpt.Enable)
	form.Add("expiryTime", inOpt.ExpiryTime)
	form.Add("listen", inOpt.Listen)
	form.Add("port", inOpt.Port)
	form.Add("protocol", inOpt.Protocol)

	resp := &ApiResponse{}
	err = c.DoForm(ctx, http.MethodPost, "/panel/inbound/add", form, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return resp, fmt.Errorf(resp.Msg)
	}
	return resp, err
}
