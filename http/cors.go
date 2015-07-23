/*
 * Copyright 2015 Fabrício Godoy
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

package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/skarllot/raiqub"
)

const (
	DEFAULT_CORS_PREFLIGHT_METHOD = "OPTIONS"
	DEFAULT_CORS_MAX_AGE          = time.Hour * 24 / time.Second
	DEFAULT_CORS_METHODS          = "OPTIONS, GET, HEAD, POST, PUT, DELETE, TRACE, CONNECT"
	DEFAULT_CORS_ORIGIN           = "*"
)

type CORSHandler struct {
	PredicateOrigin raiqub.PredicateStringFunc
	Headers         []string
	ExposedHeaders  []string
}

func NewCORSHandler() *CORSHandler {
	return &CORSHandler{
		PredicateOrigin: raiqub.TrueForAll,
		Headers: []string{
			"Origin", "X-Requested-With", "Content-Type",
			"Accept", "Authorization",
		},
		ExposedHeaders: make([]string, 0),
	}
}

func (s *CORSHandler) CreateOptionsRoutes(routes Routes) Routes {
	list := make(Routes, 0, len(routes))
	hList := make(map[string]*CORSPreflight, len(routes))
	for _, v := range routes {
		preflight, ok := hList[v.Path]
		if !ok {
			preflight = &CORSPreflight{
				*s,
				make([]string, 0, 1),
				v.MustAuth,
			}
			hList[v.Path] = preflight
		}
		preflight.Methods = append(preflight.Methods, v.Method)
	}

	for k, v := range hList {
		list = append(list, Route{
			Name:       "",
			Method:     DEFAULT_CORS_PREFLIGHT_METHOD,
			Path:       k,
			MustAuth:   v.UseCredentials,
			ActionFunc: v.ServeHTTP,
		})
	}
	return list
}

type CORSPreflight struct {
	CORSHandler
	Methods        []string
	UseCredentials bool
}

func (s *CORSPreflight) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := HttpHeader_Origin().GetReader(r.Header)
	status := http.StatusBadRequest
	defer func() {
		w.WriteHeader(status)
	}()

	if origin.Value != "" {
		if !s.PredicateOrigin(origin.Value) {
			status = http.StatusForbidden
			return
		}

		method := r.Header.Get("Access-Control-Request-Method")
		header := strings.Split(
			r.Header.Get("Access-Control-Request-Headers"),
			", ")

		if !raiqub.StringSlice(s.Methods).Exists(method) {
			return
		}

		if len(s.Headers) == 0 {
			HttpHeader_AccessControlAllowHeaders().
				SetWriter(w.Header())
		} else {
			if len(header) > 0 &&
				!raiqub.StringSlice(s.Headers).ExistsAllIgnoreCase(header) {
				return
			}
			HttpHeader_AccessControlAllowHeaders().
				SetValue(strings.Join(s.Headers, ", ")).
				SetWriter(w.Header())
		}

		HttpHeader_AccessControlAllowMethods().
			SetValue(strings.Join(s.Methods, ", ")).
			SetWriter(w.Header())
		origin.SetWriter(w.Header())
		HttpHeader_AccessControlAllowCredentials().
			SetValue(strconv.FormatBool(s.UseCredentials)).
			SetWriter(w.Header())
		// Optional
		HttpHeader_AccessControlMaxAge().
			SetValue(strconv.Itoa(int(DEFAULT_CORS_MAX_AGE))).
			SetWriter(w.Header())
		status = http.StatusOK
	} else {
		status = http.StatusNotFound
	}
}
