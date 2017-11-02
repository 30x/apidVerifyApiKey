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
package accessEntity

import (
	"github.com/apid/apid-core"
	"github.com/apid/apidVerifyApiKey/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"sync"
)

const (
	fileDataTest = "data_test.sql"
)

var _ = Describe("DataTest", func() {

	Context("query Db to get entities", func() {
		var dataTestTempDir string
		var dbMan *DbManager
		BeforeEach(func() {
			var err error
			dataTestTempDir, err = ioutil.TempDir(testTempDirBase, "sqlite3")
			Expect(err).NotTo(HaveOccurred())
			services.Config().Set("local_storage_path", dataTestTempDir)

			dbMan = &DbManager{
				DbManager: common.DbManager{
					Data:  services.Data(),
					DbMux: sync.RWMutex{},
				},
			}
			dbMan.SetDbVersion(dataTestTempDir)
		})

		It("should get apiProducts", func() {
			setupTestDb(dbMan.GetDb())
			org := "edgex01"
			testData := [][]string{
				//positive tests
				{IdentifierApiProductName, "apstest", "", ""},
				{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", ""},
				{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", IdentifierApiResource, "/**"},
				{IdentifierAppName, "apstest", "", ""},
				{IdentifierConsumerKey, "abcd", "", ""},
				{IdentifierConsumerKey, "abcd", IdentifierApiResource, "/**"},
				{IdentifierAppName, "apstest", IdentifierDeveloperId, "e41f04e8-9d3f-470a-8bfd-c7939945896c"},
				{IdentifierAppName, "apstest", IdentifierDeveloperEmail, "bar@google.com"},
				{IdentifierAppName, "apstest", IdentifierCompanyName, "DevCompany"},
				{IdentifierAppName, "apstest", IdentifierApiResource, "/**"},
				// negative tests
				{IdentifierApiProductName, "non-existent", "", ""},
				{IdentifierAppId, "non-existent", "", ""},
				{IdentifierAppId, "non-existent", IdentifierApiResource, "non-existent"},
				{IdentifierAppName, "non-existent", "", ""},
				{IdentifierConsumerKey, "non-existent", "", ""},
				{IdentifierConsumerKey, "non-existent", IdentifierApiResource, "non-existent"},
				{IdentifierAppName, "non-existent", IdentifierDeveloperId, "non-existent"},
				{IdentifierAppName, "non-existent", IdentifierDeveloperEmail, "non-existent"},
				{IdentifierAppName, "non-existent", IdentifierCompanyName, "non-existent"},
				{IdentifierAppName, "non-existent", IdentifierApiResource, "non-existent"},
			}

			var expectedApiProd = common.ApiProduct{
				Id:            "b7e0970c-4677-4b05-8105-5ea59fdcf4e7",
				Name:          "apstest",
				DisplayName:   "apstest",
				Description:   "",
				ApiResources:  "{/**}",
				ApprovalType:  "AUTO",
				Scopes:        `{""}`,
				Proxies:       `{aps,perfBenchmark}`,
				Environments:  `{prod,test}`,
				Quota:         "10000000",
				QuotaTimeUnit: "MINUTE",
				QuotaInterval: 1,
				CreatedAt:     "2017-08-18 22:12:49.363+00:00",
				CreatedBy:     "haoming@apid.git",
				UpdatedAt:     "2017-08-18 22:26:50.153+00:00",
				UpdatedBy:     "haoming@apid.git",
				TenantId:      "515211e9",
			}

			results := [][]common.ApiProduct{
				{expectedApiProd},
				{expectedApiProd},
				{expectedApiProd},
				{expectedApiProd},
				{expectedApiProd},
				{expectedApiProd},
				{expectedApiProd},
				{expectedApiProd},
				{expectedApiProd},
				{expectedApiProd},
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
			}

			for i, data := range testData {
				priKey, priVal, secKey, secVal := data[0], data[1], data[2], data[3]
				prods, err := dbMan.GetApiProducts(org, priKey, priVal, secKey, secVal)
				Expect(err).Should(Succeed())
				if len(results[i]) > 0 {
					Expect(prods).Should(Equal(results[i]))
				} else {
					Expect(prods).Should(BeZero())
				}
			}
		})
	})

})

func setupTestDb(db apid.DB) {
	bytes, err := ioutil.ReadFile(fileDataTest)
	Expect(err).Should(Succeed())
	query := string(bytes)
	log.Debug(query)
	_, err = db.Exec(query)
	Expect(err).Should(Succeed())
}
