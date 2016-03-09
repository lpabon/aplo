//
// Copyright (c) 2016 The ams Authors
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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gluster/aplo/src/ams/apps/ams"
	"github.com/gorilla/mux"
	"github.com/heketi/heketi/apps"
	//"github.com/heketi/heketi/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Port        string `json:"port"`
	AuthEnabled bool   `json:"use_auth"`
	//	JwtConfig   middleware.JwtAuthConfig `json:"jwt"`
}

var (
	AMS_VERSION = "(dev)"
	configfile  string
	showVersion bool
)

func init() {
	flag.StringVar(&configfile, "config", "", "Configuration file")
	flag.BoolVar(&showVersion, "version", false, "Show version")
}

func printVersion() {
	fmt.Printf("AMS %v\n", AMS_VERSION)
}

func main() {
	flag.Parse()
	printVersion()

	// Quit here if all we needed to do was show version
	if showVersion {
		return
	}

	// Check configuration file was given
	if configfile == "" {
		fmt.Fprintln(os.Stderr, "Please provide configuration file")
		os.Exit(1)
	}

	// Read configuration
	fp, err := os.Open(configfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open config file %v: %v\n",
			configfile,
			err.Error())
		os.Exit(1)
	}
	defer fp.Close()

	configParser := json.NewDecoder(fp)
	var options Config
	if err = configParser.Decode(&options); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to parse %v: %v\n",
			configfile,
			err.Error())
		os.Exit(1)
	}

	// Go to the beginning of the file when we pass it
	// to the application
	fp.Seek(0, os.SEEK_SET)

	// Setup a new GlusterFS application
	var app apps.Application
	amsApp := ams.NewApp(fp)
	if amsApp == nil {
		fmt.Fprintln(os.Stderr, "ERROR: Unable to start application")
		os.Exit(1)
	}
	app = amsApp

	// Create a router and do not allow any routes
	// unless defined.
	router := mux.NewRouter().StrictSlash(true)
	err = app.SetRoutes(router)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: Unable to create http server endpoints")
		os.Exit(1)
	}

	// Use negroni to add middleware.  Here we add two
	// middlewares: Recovery and Logger, which come with
	// Negroni
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())

	// Load authorization JWT middleware
	/*
		if options.AuthEnabled {
			jwtauth := middleware.NewJwtAuth(&options.JwtConfig)
			if jwtauth == nil {
				fmt.Fprintln(os.Stderr, "ERROR: Missing JWT information in config file")
				os.Exit(1)
			}

			// Add Token parser
			n.Use(jwtauth)

			fmt.Println("Authorization loaded")
		}
	*/

	// Setup routes
	n.UseHandler(router)

	// Shutdown on CTRL-C signal
	// For a better cleanup, we should shutdown the server and
	signalch := make(chan os.Signal, 1)
	signal.Notify(signalch, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to know if the server was unable to start
	done := make(chan bool)
	go func() {
		// Start the server.
		err = http.ListenAndServe(":"+options.Port, n)
		if err != nil {
			fmt.Printf("ERROR: HTTP Server error: %v\n", err)
		}
		done <- true
	}()

	// Block here for signals and errors from the HTTP server
	select {
	case <-signalch:
	case <-done:
	}
	fmt.Printf("Shutting down...\n")

	// Shutdown the application
	// :TODO: Need to shutdown the server
	app.Close()

}
