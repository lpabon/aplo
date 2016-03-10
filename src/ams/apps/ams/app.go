//
// Copyright (c) 2015 The heketi Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ams

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/heketi/rest"
	"github.com/heketi/utils"
	"io"
	"net/http"
)

const (
	ASYNC_ROUTE = "/queue"
)

var (
	logger = utils.NewLogger("[ams]", utils.LEVEL_DEBUG)
)

type App struct {
	asyncManager *rest.AsyncHttpManager
	conf         *AmsConfig
}

func NewApp(configIo io.Reader) *App {
	app := &App{}

	// Load configuration file
	app.conf = loadConfiguration(configIo)
	if app.conf == nil {
		return nil
	}

	// Setup asynchronous manager
	app.asyncManager = rest.NewAsyncHttpManager(ASYNC_ROUTE)

	// Show application has loaded
	logger.Info("Ams Application Loaded")

	return app
}

// Register Routes
func (a *App) SetRoutes(router *mux.Router) error {

	routes := rest.Routes{

		// HelloWorld
		rest.Route{
			Name:        "Hello",
			Method:      "GET",
			Pattern:     "/hello",
			HandlerFunc: a.Hello},

		// Asynchronous Manager
		rest.Route{
			Name:        "Async",
			Method:      "GET",
			Pattern:     ASYNC_ROUTE + "/{id:[A-Fa-f0-9]+}",
			HandlerFunc: a.asyncManager.HandlerStatus},

		// Endpoints
		rest.Route{
			Name:        "AmsVolumeCreate",
			Method:      "POST",
			Pattern:     "/aplo/volumes",
			HandlerFunc: a.VolumeCreate},
		rest.Route{
			Name:        "AmsVolumeExpand",
			Method:      "POST",
			Pattern:     "/aplo/volumes/expand",
			HandlerFunc: a.VolumeExpand},
		rest.Route{
			Name:        "AmsVolumeDelete",
			Method:      "DELETE",
			Pattern:     "/aplo/volumes",
			HandlerFunc: a.VolumeDelete},
		rest.Route{
			Name:        "AmsVolumeList",
			Method:      "GET",
			Pattern:     "/aplo/volumes",
			HandlerFunc: a.VolumeList},
	}

	// Register all routes from the App
	for _, route := range routes {

		// Add routes from the table
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

	}

	return nil

}

func (a *App) Close() {
	logger.Info("Closed")
}

func (a *App) Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "HelloWorld from AMS Application")
}

// Middleware function
func (a *App) Auth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Everything is clean
	next(w, r)
}
