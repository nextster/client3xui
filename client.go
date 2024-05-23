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

import "net/http"

type Client struct {
	url, subUrl, username, password string
	httpClient                      *http.Client
	sessionCookie                   *http.Cookie
}

func New(c Config) *Client {
	cl := &Client{url: c.Url, subUrl: c.SubUrl, username: c.Username, password: c.Password}
	if c.Client == nil {
		cl.httpClient = http.DefaultClient
	} else {
		cl.httpClient = c.Client
	}
	return cl
}
