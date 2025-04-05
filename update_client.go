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
)

// Add client to an inbound.
func (c *Client) DeleteClient(ctx context.Context, inboundId uint, clientUuid string) (*ApiResponse, error) {
	resp := &ApiResponse{}
	inboundIdStr := strconv.FormatUint(uint64(inboundId), 10)
	err := c.Do(ctx, http.MethodPost, "/panel/api/inbounds/"+inboundIdStr+"/delClient/"+clientUuid, nil, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return resp, fmt.Errorf(resp.Msg)
	}
	return resp, err
}

func (c *Client) UpdateClient(ctx context.Context, inboundId uint, client InboundClient) (*ApiResponse, error) {
	resp := &ApiResponse{}
	inboundIdStr := strconv.FormatUint(uint64(inboundId), 10)

	// Create client settings using InboundClient struct
	clientSettings := struct {
		Clients []InboundClient `json:"clients"`
	}{
		Clients: []InboundClient{client},
	}

	// Convert settings to JSON string
	settingsBytes, err := json.Marshal(clientSettings)
	if err != nil {
		return nil, err
	}

	// Create form data
	form := url.Values{}
	form.Add("id", inboundIdStr)
	form.Add("settings", string(settingsBytes))

	err = c.DoForm(ctx, http.MethodPost, "/panel/inbound/updateClient/"+client.ID, form, resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return resp, fmt.Errorf(resp.Msg)
	}
	return resp, nil
}
