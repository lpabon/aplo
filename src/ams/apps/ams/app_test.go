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
	"bytes"
	"github.com/heketi/tests"
	"testing"
)

func TestAppBadConfigData(t *testing.T) {
	data := []byte(`{ bad json }`)
	app := NewApp(bytes.NewBuffer(data))
	tests.Assert(t, app == nil)

	data = []byte(`{}`)
	app = NewApp(bytes.NewReader(data))
	tests.Assert(t, app == nil)

	data = []byte(`{
		"ams" : {}
		}`)
	app = NewApp(bytes.NewReader(data))
	tests.Assert(t, app == nil)
}

func TestAppConfig(t *testing.T) {

	data := []byte(`{
		"ams" : {
			"gluster_pod_name" : "gluster-pod",
			"gluster_hostname" : "gluster-host",
			"gluster_volume" : "gluster-volume"
		}
	}`)

	app := NewApp(bytes.NewReader(data))
	tests.Assert(t, app != nil)
	tests.Assert(t, app.conf.GlusterPodName == "gluster-pod")
	tests.Assert(t, app.conf.GlusterHost == "gluster-host")
	tests.Assert(t, app.conf.GlusterVolume == "gluster-volume")
}
