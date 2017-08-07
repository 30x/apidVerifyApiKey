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
	"github.com/30x/apid-core"
	"github.com/apigee-labs/transicator/common"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	//"github.com/30x/apid-core/data"
)

var _ = Describe("listener", func() {

	Context("KMS create/updates verification via changes for Developer", func() {

		handler := apigeeSyncHandler{}

		It("should set DB to appropriate version", func() {

			//saveDb := handler.dbMan.getDb()

			s := &common.Snapshot{
				SnapshotInfo: "test_snapshot",
				Tables:       []common.Table{},
			}

			handler.Handle(s)

			//expectedDB, err := handler.dbMan.data.DBVersion(s.SnapshotInfo)
			//Expect(err).NotTo(HaveOccurred())
			//
			//Expect(getDB() == expectedDB).Should(BeTrue())
			//
			////restore the db to the valid one
			//setDB(saveDb)
		})

	})
})

func addScopes(db apid.DB) {
	txn, _ := db.Begin()
	txn.Exec("INSERT INTO EDGEX_DATA_SCOPE (id, _change_selector, apid_cluster_id, scope, org, env) "+
		"VALUES"+
		"($1,$2,$3,$4,$5,$6)",
		"ABCDE",
		"some_cluster_id",
		"some_cluster_id",
		"tenant_id_xxxx",
		"test_org0",
		"Env_0",
	)
	txn.Exec("INSERT INTO EDGEX_DATA_SCOPE (id, _change_selector, apid_cluster_id, scope, org, env) "+
		"VALUES"+
		"($1,$2,$3,$4,$5,$6)",
		"XYZ",
		"test_org0",
		"somecluster_id",
		"tenant_id_0",
		"test_org0",
		"Env_0",
	)
	log.Info("Inserted EDGEX_DATA_SCOPE for test")
	txn.Commit()
}
