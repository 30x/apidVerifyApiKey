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
	"github.com/apid/apidApiMetadata/common"
	tran "github.com/apigee-labs/transicator/common"
)

const (
	APIGEE_SYNC_EVENT = "ApigeeSync"
)

type apigeeSyncHandler struct {
	dbMans    []common.DbManagerInterface
	apiMans   []common.ApiManagerInterface
	cipherMan common.CipherManagerInterface
}

func (h *apigeeSyncHandler) initListener(services apid.Services) {
	services.Events().Listen(APIGEE_SYNC_EVENT, h)
}

func (h *apigeeSyncHandler) String() string {
	return "verifyAPIKey"
}

func (h *apigeeSyncHandler) processSnapshot(snapshot *tran.Snapshot) {
	log.Debugf("Snapshot received. Switching to DB version: %s", snapshot.SnapshotInfo)
	// set db version for all packages
	for _, dbMan := range h.dbMans {
		dbMan.SetDbVersion(snapshot.SnapshotInfo)
	}
	// retrieve encryption keys
	orgs, err := h.dbMans[0].GetOrgs()
	if err != nil {
		log.Panicf("Failed to get orgs: %v", err)
	}
	h.cipherMan.AddOrgs(orgs)
	// idempotent init api for all packages
	for _, apiMan := range h.apiMans {
		apiMan.InitAPI()
	}
	log.Debug("Snapshot processed")
}

func (h *apigeeSyncHandler) Handle(e apid.Event) {

	if snapData, ok := e.(*tran.Snapshot); ok {
		h.processSnapshot(snapData)
	} else { //TODO handle changelist and retrieve key for new orgs
		log.Debugf("Received event. No action required for verifyApiKey plugin. Ignoring. %v", e)
	}
}
