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
	"github.com/30x/apid-core/factory"
	"github.com/apigee-labs/transicator/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"sync"
)

var _ = Describe("listener", func() {

	var listnerTestSyncHandler apigeeSyncHandler
	var listnerTestTempDir string
	var _ = BeforeEach(func() {
		var err error
		listnerTestTempDir, err = ioutil.TempDir("", "listner_test")
		s := factory.DefaultServicesFactory()
		apid.Initialize(s)
		config := apid.Config()
		config.Set("data_path", listnerTestTempDir)
		Expect(err).NotTo(HaveOccurred())

		apid.InitializePlugins("")

		db, err := apid.Data().DB()
		Expect(err).NotTo(HaveOccurred())

		dbMan := &dbManager{
			data:  s.Data(),
			dbMux: sync.RWMutex{},
			db:    db,
		}

		listnerTestSyncHandler = apigeeSyncHandler{
			dbMan:  dbMan,
			apiMan: apiManager{},
		}

		listnerTestSyncHandler.initListener(s)
	})

	var _ = AfterEach(func() {
		os.RemoveAll(listnerTestTempDir)
	})

	Context("Apigee Sync Event Processing", func() {

		It("should set DB to appropriate version", func() {
			s := &common.Snapshot{
				SnapshotInfo: "test_snapshot",
				Tables:       []common.Table{},
			}
			listnerTestSyncHandler.Handle(s)
			Expect(listnerTestSyncHandler.dbMan.getDbVersion()).Should(BeEquivalentTo(s.SnapshotInfo))

		})

		It("should not change version for chang event", func() {

			version := listnerTestSyncHandler.dbMan.getDbVersion()
			s := &common.Change{
				ChangeSequence: 12321,
				Table:          "",
			}
			testSyncHandler.Handle(s)
			Expect(listnerTestSyncHandler.dbMan.getDbVersion() == version).Should(BeTrue())

		})

	})
})
