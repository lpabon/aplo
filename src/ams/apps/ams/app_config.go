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
	"io"
)

type AmsConfig struct {
	GlusterPodName  string `json:"gluster_pod_name"`
	GlusterHost     string `json:"gluster_hostname"`
	GlusterVolume   string `json:"gluster_volume"`
	GlusterMountDir string `json:"gluster_mountdir"`
}

type ConfigFile struct {
	Ams AmsConfig `json:"ams"`
}

func loadConfiguration(configIo io.Reader) *AmsConfig {
	configParser := json.NewDecoder(configIo)

	var config ConfigFile
	if err := configParser.Decode(&config); err != nil {
		logger.LogError("Unable to parse config file: %v\n",
			err.Error())
		return nil
	}

	if config.Ams.GlusterPodName == "" {
		logger.LogError("Configuration file is missing gluster pod name")
		return nil
	}
	if config.Ams.GlusterHost == "" {
		logger.LogError("Configuration file is missing gluster hostname")
		return nil
	}
	if config.Ams.GlusterVolume == "" {
		logger.LogError("Configuration file is missing gluster volume name")
		return nil
	}

	if config.Ams.GlusterMountDir == "" {
		config.Ams.GlusterMountDir = "/mnt/gluster"
	}

	return &config.Ams
}
