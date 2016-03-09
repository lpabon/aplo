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
	"encoding/json"
	"github.com/heketi/utils"
	"net/http"
	"os"
)

func (a *App) VolumeCreate(w http.ResponseWriter, r *http.Request) {
	var msg VolumeCreateRequest
	err := utils.GetJsonFromRequest(r, &msg)
	if err != nil {
		http.Error(w, "request unable to be parsed", 422)
		return
	}

	// Check the message has devices
	if msg.Size < 1 {
		http.Error(w, "Invalid volume size", http.StatusBadRequest)
		return
	}

	// Create directory on GlusterFS volume
	os.Mkdir(a.conf.GlusterMountDir+"/"+msg.Name, 0755)

	// Setup quota

	// Send response
	var resp VolumeCreateResponse
	resp.Size = msg.Size
	resp.Name = msg.Name
	resp.Gid = resp.Gid
	resp.Mount.GlusterFS.MountPoint = a.conf.GlusterHost + ":" +
		a.conf.GlusterVolume + "/" + resp.Name
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func (a *App) VolumeExpand(w http.ResponseWriter, r *http.Request) {
}
func (a *App) VolumeDelete(w http.ResponseWriter, r *http.Request) {
}
func (a *App) VolumeList(w http.ResponseWriter, r *http.Request) {
}
