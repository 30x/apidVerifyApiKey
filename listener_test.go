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
	"github.com/apid/apid-core/factory"
	tran "github.com/apigee-labs/transicator/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
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
		listnerTestSyncHandler = initManagers(s)
	})

	var _ = AfterEach(func() {
		os.RemoveAll(listnerTestTempDir)
	})

	Context("Apigee Sync Event Processing", func() {

		It("should set DB to appropriate version", func() {
			s := &tran.Snapshot{
				SnapshotInfo: "test_snapshot",
				Tables:       []tran.Table{},
			}
			listnerTestSyncHandler.Handle(s)
			for _, dbMan := range listnerTestSyncHandler.dbMans {
				Expect(dbMan.GetDbVersion()).Should(BeEquivalentTo(s.SnapshotInfo))
			}

		})

		It("should not change version for chang event", func() {

			version := listnerTestSyncHandler.dbMans[0].GetDbVersion()
			s := &tran.Change{
				ChangeSequence: 12321,
				Table:          "",
			}
			listnerTestSyncHandler.Handle(s)
			for _, dbMan := range listnerTestSyncHandler.dbMans {
				Expect(dbMan.GetDbVersion() == version).Should(BeTrue())
			}

		})

	})
})
