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

package apidApiMetadata

import (
	"github.com/apid/apid-core"
	"github.com/apid/apid-core/util"
	"github.com/apid/apidApiMetadata/accessEntity"
	"github.com/apid/apidApiMetadata/common"
	"github.com/apid/apidApiMetadata/verifyApiKey"
	"net/http"
	"sync"
	"time"
)

const (
	maxIdleConnsPerHost      = 1
	httpTimeout              = 5 * time.Minute
	configBearerToken        = "apigeesync_bearer_token"
	configRetrieveEncKeyBase = "apimetadata_encryption_key_server_base"
)

var (
	services apid.Services
	log      apid.LogService
)

func init() {
	apid.RegisterPlugin(initPlugin, common.PluginData)
}

func initPlugin(s apid.Services) (apid.PluginData, error) {
	services = s
	log = services.Log().ForModule("apidApiMetadata")
	verifyApiKey.SetApidServices(services, log)
	accessEntity.SetApidServices(services, log)
	common.SetApidServices(services, log)
	log.Debug("start init")
	initManagers(services)
	log.Debug("end init")

	return common.PluginData, nil
}

// init http client
func createHttpClient() *http.Client {
	tr := util.Transport(services.Config().GetString(util.ConfigfwdProxyPortURL))
	tr.MaxIdleConnsPerHost = maxIdleConnsPerHost
	client := &http.Client{
		Transport: tr,
		Timeout:   httpTimeout,
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			req.Header.Set("Authorization", "Bearer "+services.Config().GetString(configBearerToken))
			return nil
		},
	}
	return client
}

func initManagers(services apid.Services) *apigeeSyncHandler {

	cipherMan := common.CreateCipherManager(createHttpClient(), services.Config().GetString(configRetrieveEncKeyBase))

	verifyDbMan := &verifyApiKey.DbManager{
		DbManager: common.DbManager{
			Data:          services.Data(),
			DbMux:         sync.RWMutex{},
			CipherManager: cipherMan,
		},
	}
	verifyApiMan := &verifyApiKey.ApiManager{
		DbMan:             verifyDbMan,
		VerifiersEndpoint: verifyApiKey.ApiPath,
	}

	entityDbMan := &accessEntity.DbManager{
		DbManager: common.DbManager{
			Data:          services.Data(),
			DbMux:         sync.RWMutex{},
			CipherManager: cipherMan,
		},
	}

	entityApiMan := &accessEntity.ApiManager{
		DbMan:            entityDbMan,
		AccessEntityPath: accessEntity.AccessEntityPath,
	}

	syncHandler := &apigeeSyncHandler{
		dbMans:    []common.DbManagerInterface{verifyDbMan, entityDbMan},
		apiMans:   []common.ApiManagerInterface{verifyApiMan, entityApiMan},
		cipherMan: cipherMan,
	}
	syncHandler.initListener(services)
	return syncHandler
}
