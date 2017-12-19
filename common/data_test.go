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
package common

import (
	"github.com/apid/apid-core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"reflect"
	"sort"
	"sync"
)

const fileDataTest = "data_test.sql"

var _ = Describe("DataTest", func() {

	Context("DB", func() {
		var dataTestTempDir string
		var testDbMan *DbManager
		BeforeEach(func() {
			var err error
			dataTestTempDir, err = ioutil.TempDir(testTempDirBase, "sqlite3")
			Expect(err).NotTo(HaveOccurred())
			services.Config().Set("local_storage_path", dataTestTempDir)

			testDbMan = &DbManager{
				Data:  services.Data(),
				DbMux: sync.RWMutex{},
			}
			testDbMan.SetDbVersion(dataTestTempDir)
			Expect(testDbMan.GetDbVersion()).Should(Equal(dataTestTempDir))
			setupTestDb(testDbMan.GetDb())
		})

		It("should get kms attributes", func() {
			attributes := testDbMan.GetKmsAttributes("bc811169", "40753e12-a50a-429d-9121-e571eb4e43a9", "85629786-37c5-4e8c-bb45-208f3360d005", "50321842-d6ee-4e92-91b9-37234a7920c1", "test-invalid")
			Expect(len(attributes)).Should(BeEquivalentTo(3))
			Expect(len(attributes["40753e12-a50a-429d-9121-e571eb4e43a9"])).Should(BeEquivalentTo(1))
			Expect(len(attributes["85629786-37c5-4e8c-bb45-208f3360d005"])).Should(BeEquivalentTo(2))
			Expect(len(attributes["50321842-d6ee-4e92-91b9-37234a7920c1"])).Should(BeEquivalentTo(5))
			Expect(len(attributes["test-invalid"])).Should(BeEquivalentTo(0))
		})

		It("Should get all orgs", func() {
			orgs, err := testDbMan.GetOrgs()
			Expect(err).Should(Succeed())
			sort.Strings(orgs)
			Expect(orgs).Should(Equal([]string{"apid-haoming", "apid-test"}))
		})

	})

	Context("Validate common.JsonToStringArray", func() {

		It("should transform simple valid json", func() {
			array := JsonToStringArray("[\"test-1\", \"test-2\"]")
			Expect(reflect.DeepEqual(array, []string{"test-1", "test-2"})).Should(BeTrue())
		})
		It("should transform simple single valid json", func() {
			array := JsonToStringArray("[\"test-1\"]")
			Expect(reflect.DeepEqual(array, []string{"test-1"})).Should(BeTrue())
		})
		It("should transform simple fake json", func() {
			s := JsonToStringArray("{test-1,test-2}")
			Expect(reflect.DeepEqual(s, []string{"test-1", "test-2"})).Should(BeTrue())
		})
		It("should transform simple single valued fake json", func() {
			s := JsonToStringArray("{test-1}")
			Expect(reflect.DeepEqual(s, []string{"test-1"})).Should(BeTrue())
		})
		It("space between fields considered as valid char", func() {
			s := JsonToStringArray("{test-1, test-2}")
			Expect(reflect.DeepEqual(s, []string{"test-1", " test-2"})).Should(BeTrue())
		})
		It("remove only last braces", func() {
			s := JsonToStringArray("{test-1,test-2}}")
			Expect(reflect.DeepEqual(s, []string{"test-1", "test-2}"})).Should(BeTrue())
		})

	})
})

func setupTestDb(db apid.DB) {
	bytes, err := ioutil.ReadFile(fileDataTest)
	Expect(err).Should(Succeed())
	query := string(bytes)
	_, err = db.Exec(query)
	Expect(err).Should(Succeed())
}
