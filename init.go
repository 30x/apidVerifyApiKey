// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apidVerifyApiKey

import (
	"github.com/apid/apid-core"
	"github.com/apid/apidVerifyApiKey/verifyApiKey"
	"sync"
)

var (
	services apid.Services
	log      apid.LogService
)

func init() {
	apid.RegisterPlugin(initPlugin, pluginData)
}

func initPlugin(s apid.Services) (apid.PluginData, error) {
	services = s
	log = services.Log().ForModule("apidApiMetadata")
	verifyApiKey.SetApidServices(services, log)
	log.Debug("start init")

	log = services.Log()
	verifyDbMan := &verifyApiKey.DbManager{
		Data:  services.Data(),
		DbMux: sync.RWMutex{},
	}
	verifyApiMan := &verifyApiKey.ApiManager{
		DbMan:             verifyDbMan,
		VerifiersEndpoint: verifyApiKey.ApiPath,
	}

	syncHandler := apigeeSyncHandler{
		dbMans:  []DbManagerInterface{verifyDbMan},
		apiMans: []ApiManagerInterface{verifyApiMan},
	}

	syncHandler.initListener(services)

	log.Debug("end init")

	return pluginData, nil
}
