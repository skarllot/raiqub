/*
 * Copyright (C) 2015 Fabrício Godoy <skarllot@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
 */

package http

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestPreflightHeaders(t *testing.T) {
	var conf = struct {
		method    string
		origin    string
		headers   string
		reqmethod string
	}{
		"OPTIONS",
		"http://localhost",
		"Content-Type",
		"POST",
	}

	cors := NewCORSHandler()
	routes := Routes{
		Route{
			"Test1",
			"GET",
			"/test",
			false,
			nil,
		},
		Route{
			"Test2",
			"POST",
			"/test",
			true,
			nil,
		},
	}
	preflight := cors.CreatePreflight(routes)

	if len(preflight) != 1 {
		t.Fatal("Should generate only one preflight route")
	}

	ts := httptest.NewServer(preflight[0].ActionFunc)
	defer ts.Close()

	client := http.Client{}
	req, err := http.NewRequest(conf.method, ts.URL, nil)
	HttpHeader_Origin().
		SetValue(conf.origin).SetWriter(req.Header)
	HttpHeader_AccessControlRequestHeaders().
		SetValue(conf.headers).SetWriter(req.Header)
	HttpHeader_AccessControlRequestMethod().
		SetValue(conf.reqmethod).SetWriter(req.Header)

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error trying to call HTTP %s: %v", conf.method, err)
	}

	var header *HttpHeader

	header = HttpHeader_AccessControlAllowOrigin()
	if header.GetReader(res.Header).Value == "" {
		t.Errorf("The header %s was not found", header.Name)
	}
	header = HttpHeader_Origin()
	if header.GetReader(res.Header).Value != "" {
		t.Errorf("The header %s should not be found", header.Name)
	}
	header = HttpHeader_AccessControlAllowHeaders()
	if header.GetReader(res.Header).Value == "" {
		t.Error("The header %s was not found", header.Name)
	}
	header = HttpHeader_AccessControlAllowMethods()
	if !strings.Contains(header.GetReader(res.Header).Value, conf.reqmethod) {
		t.Errorf("The header %s doesn't allow '%s' HTTP method",
			header.Name, conf.reqmethod)
	}
	header = HttpHeader_AccessControlAllowCredentials()
	if b, err := strconv.ParseBool(
		header.GetReader(res.Header).Value); err != nil || !b {
		t.Errorf("The header %s should be '%s'",
			header.Name, strconv.FormatBool(true))
	}
}
