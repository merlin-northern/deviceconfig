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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/mendersoftware/go-lib-micro/rest.utils"

	"github.com/mendersoftware/deviceconfig/app"
	"github.com/mendersoftware/deviceconfig/model"
	"github.com/mendersoftware/deviceconfig/store"
)

type ManagementAPI struct {
	App app.App
}

func NewManagementAPI(app app.App) *ManagementAPI {
	return &ManagementAPI{
		App: app,
	}
}

func (api *ManagementAPI) ConfigurationSet(c *gin.Context) {
	var configuration model.Configuration

	ctx := c.Request.Context()
	devID, err := uuid.Parse(c.Param("device_id"))
	if err != nil {
		rest.RenderError(c,
			http.StatusBadRequest,
			errors.Wrap(err, "correctly formatted device id is needed"),
		)
		return
	}

	err = c.ShouldBindJSON(&configuration)
	if err != nil {
		rest.RenderError(c,
			http.StatusBadRequest,
			errors.Wrap(err, "malformed request body"),
		)
		return
	}

	err = api.App.ConfigurationSet(ctx, devID, configuration)
	if err != nil {
		switch cause := errors.Cause(err); cause {
		case store.ErrDeviceAlreadyExists:
			rest.RenderError(c, http.StatusConflict, cause)
		default:
			c.Error(err) //nolint:errcheck
			rest.RenderError(c,
				http.StatusInternalServerError,
				errors.New(http.StatusText(http.StatusInternalServerError)),
			)
		}
		return
	}
	c.Status(http.StatusCreated)
}

func (api *ManagementAPI) ConfigurationGet(c *gin.Context) {
	ctx := c.Request.Context()

	devID, err := uuid.Parse(c.Param("device_id"))
	if err != nil {
		rest.RenderError(c,
			http.StatusBadRequest,
			errors.Wrap(err, "correctly formatted device id is needed"),
		)
		return
	}

	device, err := api.App.GetDevice(ctx, devID)
	if err != nil {
		rest.RenderError(c,
			http.StatusNotFound,
			errors.New(http.StatusText(http.StatusNotFound)),
		)
		return
	}

	c.JSON(http.StatusOK, device)
}
