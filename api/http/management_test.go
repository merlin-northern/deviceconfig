// Copyright 2021 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package http

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	mapp "github.com/mendersoftware/deviceconfig/app/mocks"
	"github.com/mendersoftware/deviceconfig/model"
	"github.com/mendersoftware/go-lib-micro/rest.utils"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationSet(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name string

		TenantID string
		Request  *http.Request

		App    *mapp.App
		Error  *rest.Error
		Status int
	}{
		{
			Name: "ok",

			Request: func() *http.Request {
				body, _ := json.Marshal(map[string]interface{}{
					"expected": []map[string]interface{}{
						{
							"key":   "key0",
							"value": "value0",
						},
					},
				})

				repl := strings.NewReplacer(
					":device_id", uuid.NewSHA1(
						uuid.NameSpaceDNS, []byte("mender.io"),
					).String(),
				)
				req, _ := http.NewRequest("PUT",
					"http://localhost"+URIManagement+repl.Replace(URIConfiguration),
					bytes.NewReader(body),
				)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibWVuZGVyLnBsYW4iOiJlbnRlcnByaXNlIn0.s27fi93Qik81WyBmDB5APE0DfGko7Pq8BImbp33-gy4")
				return req
			}(),

			App: func() *mapp.App {
				app := new(mapp.App)
				app.On("ConfigurationSet",
					contextMatcher,
					mock.AnythingOfType("uuid.UUID"),
					model.Configuration{
						Expected: []model.Attribute{
							{
								Key:   "key0",
								Value: "value0",
							},
						},
						Actual: nil,
					},
				).Return(nil)
				return app
			}(),
			Status: http.StatusCreated,
		},

		{
			Name: "error no auth",

			Request: func() *http.Request {
				body, _ := json.Marshal(map[string]interface{}{
					"expected": []map[string]interface{}{
						{
							"key":   "key0",
							"value": "value0",
						},
					},
				})

				repl := strings.NewReplacer(
					":device_id", uuid.NewSHA1(
						uuid.NameSpaceDNS, []byte("mender.io"),
					).String(),
				)
				req, _ := http.NewRequest("PUT",
					"http://localhost"+URIManagement+repl.Replace(URIConfiguration),
					bytes.NewReader(body),
				)
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),

			App: func() *mapp.App {
				app := new(mapp.App)
				return app
			}(),
			Status: http.StatusUnauthorized,
		},

		{
			Name: "error bad token format",

			Request: func() *http.Request {
				body, _ := json.Marshal(map[string]interface{}{
					"expected": []map[string]interface{}{
						{
							"key":   "key0",
							"value": "value0",
						},
					},
				})

				repl := strings.NewReplacer(
					":device_id", uuid.NewSHA1(
						uuid.NameSpaceDNS, []byte("mender.io"),
					).String(),
				)
				req, _ := http.NewRequest("PUT",
					"http://localhost"+URIManagement+repl.Replace(URIConfiguration),
					bytes.NewReader(body),
				)
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),

			App: func() *mapp.App {
				app := new(mapp.App)
				return app
			}(),
			Status: http.StatusUnauthorized,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			defer tc.App.AssertExpectations(t)
			router := NewRouter(tc.App)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, tc.Request)
			assert.Equal(t, tc.Status, w.Code)
			if tc.Error != nil {
				b, _ := json.Marshal(tc.Error)
				assert.JSONEq(t, string(b), string(w.Body.Bytes()))
			}
		})
	}
}

func TestConfigurationGet(t *testing.T) {
	t.Parallel()

	device := model.Device{
		ID: uuid.New(),
		DesiredAttributes: []model.Attribute{
			{
				Key:   "key0",
				Value: "value0",
			},
			{
				Key:   "key2",
				Value: "value2",
			},
		},
		CurrentAttributes: []model.Attribute{
			{
				Key:   "key1",
				Value: "value1",
			},
			{
				Key:   "key3",
				Value: "value3",
			},
		},
		UpdatedTS: time.Now(),
		ReportTS:  time.Now(),
	}

	testCases := []struct {
		Name string

		TenantID string
		Request  *http.Request

		App    *mapp.App
		Error  *rest.Error
		Status int
	}{
		{
			Name: "ok",

			Request: func() *http.Request {
				repl := strings.NewReplacer(
					":device_id", uuid.NewSHA1(
						uuid.NameSpaceDNS, []byte("mender.io"),
					).String(),
				)
				req, _ := http.NewRequest("GET",
					"http://localhost"+URIManagement+repl.Replace(URIConfiguration),
					nil,
				)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibWVuZGVyLnBsYW4iOiJlbnRlcnByaXNlIn0.s27fi93Qik81WyBmDB5APE0DfGko7Pq8BImbp33-gy4")
				return req
			}(),

			App: func() *mapp.App {
				app := new(mapp.App)
				app.On("GetDevice",
					contextMatcher,
					mock.AnythingOfType("uuid.UUID"),
				).Return(device, nil)
				return app
			}(),
			Status: http.StatusOK,
		},

		{
			Name: "error no auth",

			Request: func() *http.Request {
				repl := strings.NewReplacer(
					":device_id", uuid.NewSHA1(
						uuid.NameSpaceDNS, []byte("mender.io"),
					).String(),
				)
				req, _ := http.NewRequest("GET",
					"http://localhost"+URIManagement+repl.Replace(URIConfiguration),
					nil,
				)
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),

			App: func() *mapp.App {
				app := new(mapp.App)
				return app
			}(),
			Status: http.StatusUnauthorized,
		},

		{
			Name: "error bad token format",

			Request: func() *http.Request {
				repl := strings.NewReplacer(
					":device_id", uuid.NewSHA1(
						uuid.NameSpaceDNS, []byte("mender.io"),
					).String(),
				)
				req, _ := http.NewRequest("GET",
					"http://localhost"+URIManagement+repl.Replace(URIConfiguration),
					nil,
				)
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),

			App: func() *mapp.App {
				app := new(mapp.App)
				return app
			}(),
			Status: http.StatusUnauthorized,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			defer tc.App.AssertExpectations(t)
			router := NewRouter(tc.App)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, tc.Request)
			assert.Equal(t, tc.Status, w.Code)
			if w.Code == http.StatusOK {
				var d model.Device
				json.Unmarshal(w.Body.Bytes(), &d)
				t.Logf("got: %+v", d)
				d.UpdatedTS = time.Unix(1, 0)
				device.UpdatedTS = time.Unix(1, 0)
				d.ReportTS = time.Unix(1, 0)
				device.ReportTS = time.Unix(1, 0)
				assert.Equal(t, d, device)
			}
			if tc.Error != nil {
				b, _ := json.Marshal(tc.Error)
				assert.JSONEq(t, string(b), string(w.Body.Bytes()))
			}
		})
	}
}
