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
	"github.com/apid/apidApiMetadata/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"sync"
)

const (
	fileDataTest = "data_test.sql"
	// SQL Injection
	// If select foo from bar where id in (' + condition + ') is used,
	// the hacked sql would be like "select XXX from XXX where id in ('1') or ('1'=='1');"
	sqlInjectionStmt = "1') or ('1'=='1"
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
					Data:          services.Data(),
					DbMux:         sync.RWMutex{},
					CipherManager: &DummyCipherMan{},
				},
			}
			dbMan.SetDbVersion(dataTestTempDir)
			setupTestDb(dbMan.GetDb())
		})

		Describe("Get structs", func() {
			It("should get apiProducts", func() {
				testData := [][]string{
					//positive tests
					{IdentifierApiProductName, "apstest", "", "", "apid-haoming"},
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", "", "apid-haoming"},
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", IdentifierApiResource, "/**", "apid-haoming"},
					{IdentifierAppName, "apstest", "", "", "apid-haoming"},
					{IdentifierConsumerKey, "abcd", "", "", "apid-haoming"},
					{IdentifierConsumerKey, "abcd", IdentifierApiResource, "/**", "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierDeveloperId, "e41f04e8-9d3f-470a-8bfd-c7939945896c", "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierDeveloperEmail, "bar@google.com", "apid-haoming"},
					{IdentifierAppName, "testappahhis", IdentifierCompanyName, "testcompanyhflxv", "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierApiResource, "/**", "apid-haoming"},
					// negative tests
					{IdentifierApiProductName, "apstest", "", "", "non-existent"},
					{IdentifierApiProductName, "non-existent", "", "", "apid-haoming"},
					{IdentifierAppId, "non-existent", "", "", "apid-haoming"},
					{IdentifierAppId, "non-existent", IdentifierApiResource, "non-existent", "apid-haoming"},
					{IdentifierAppName, "non-existent", "", "", "apid-haoming"},
					{IdentifierConsumerKey, "non-existent", "", "", "apid-haoming"},
					{IdentifierConsumerKey, "non-existent", IdentifierApiResource, "non-existent", "apid-haoming"},
					{IdentifierAppName, "non-existent", IdentifierDeveloperId, "non-existent", "apid-haoming"},
					{IdentifierAppName, "non-existent", IdentifierDeveloperEmail, "non-existent", "apid-haoming"},
					{IdentifierAppName, "non-existent", IdentifierCompanyName, "non-existent", "apid-haoming"},
					{IdentifierAppName, "non-existent", IdentifierApiResource, "non-existent", "apid-haoming"},
					// SQL Injection
					{IdentifierApiProductName, "apstest", "", "", sqlInjectionStmt},
					{IdentifierApiProductName, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierAppId, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierAppName, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierConsumerKey, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierDeveloperId, sqlInjectionStmt, "apid-haoming"},
					{IdentifierAppName, "testappahhis", IdentifierDeveloperEmail, sqlInjectionStmt, "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierCompanyName, sqlInjectionStmt, "apid-haoming"},
				}

				var expectedDevApiProd = common.ApiProduct{
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

				var expectedComApiProd = common.ApiProduct{
					Id:            "fea8a6d5-8d34-477f-ac82-c397eaec06af",
					Name:          "testproductsdljnkpt",
					DisplayName:   "testproductsdljnkpt",
					Description:   "",
					ApiResources:  "{/res1}",
					ApprovalType:  "AUTO",
					Scopes:        `{}`,
					Proxies:       `{}`,
					Environments:  `{test}`,
					Quota:         "",
					QuotaTimeUnit: "",
					QuotaInterval: 0,
					CreatedAt:     "2017-11-02 16:00:15.608+00:00",
					CreatedBy:     "haoming@apid.git",
					UpdatedAt:     "2017-11-02 16:00:18.125+00:00",
					UpdatedBy:     "haoming@apid.git",
					TenantId:      "515211e9",
				}

				results := [][]common.ApiProduct{
					{expectedDevApiProd},
					{expectedDevApiProd},
					{expectedDevApiProd},
					{expectedDevApiProd},
					{expectedDevApiProd},
					{expectedDevApiProd},
					{expectedDevApiProd},
					{expectedDevApiProd},
					{expectedComApiProd},
					{expectedDevApiProd},
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
					priKey, priVal, secKey, secVal, org := data[0], data[1], data[2], data[3], data[4]
					prods, err := dbMan.GetApiProducts(org, priKey, priVal, secKey, secVal)
					Expect(err).Should(Succeed())
					if len(results[i]) > 0 {
						Expect(prods).Should(Equal(results[i]))
					} else {
						Expect(prods).Should(BeZero())
					}
				}
			})

			It("should get apps", func() {
				testData := [][]string{
					//positive tests
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", "", "apid-haoming"},
					{IdentifierAppName, "apstest", "", "", "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierDeveloperId, "e41f04e8-9d3f-470a-8bfd-c7939945896c", "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierDeveloperEmail, "bar@google.com", "apid-haoming"},
					{IdentifierAppName, "testappahhis", IdentifierCompanyName, "testcompanyhflxv", "apid-haoming"},
					{IdentifierConsumerKey, "abcd", "", "", "apid-haoming"},
					// negative tests
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", "", "non-existent"},
					{IdentifierAppId, "non-existent", "", "", "apid-haoming"},
					{IdentifierAppName, "non-existent", "", "", "apid-haoming"},
					{IdentifierAppName, "non-existent", IdentifierDeveloperId, "non-existent", "apid-haoming"},
					{IdentifierAppName, "non-existent", IdentifierDeveloperEmail, "non-existent", "apid-haoming"},
					{IdentifierAppName, "non-existent", IdentifierCompanyName, "non-existent", "apid-haoming"},
					{IdentifierConsumerKey, "non-existent", "", "", "apid-haoming"},
					// SQL Injection
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", "", sqlInjectionStmt},
					{IdentifierAppId, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierAppName, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierDeveloperId, sqlInjectionStmt, "apid-haoming"},
					{IdentifierAppName, "apstest", IdentifierDeveloperEmail, sqlInjectionStmt, "apid-haoming"},
					{IdentifierAppName, "testappahhis", IdentifierCompanyName, sqlInjectionStmt, "apid-haoming"},
					{IdentifierConsumerKey, sqlInjectionStmt, "", "", "apid-haoming"},
				}

				var expectedDevApp = common.App{
					Id:          "408ad853-3fa0-402f-90ee-103de98d71a5",
					TenantId:    "515211e9",
					Name:        "apstest",
					DisplayName: "apstest",
					AccessType:  "READ",
					CallbackUrl: "https://www.google.com",
					Status:      "APPROVED",
					AppFamily:   "default",
					CompanyId:   "",
					DeveloperId: "e41f04e8-9d3f-470a-8bfd-c7939945896c",
					ParentId:    "e41f04e8-9d3f-470a-8bfd-c7939945896c",
					Type:        "DEVELOPER",
					CreatedAt:   "2017-08-18 22:13:18.325+00:00",
					CreatedBy:   "haoming@apid.git",
					UpdatedAt:   "2017-08-18 22:13:18.325+00:00",
					UpdatedBy:   "haoming@apid.git",
				}

				var expectedComApp = common.App{
					Id:          "35608afe-2715-4064-bb4d-3cbb4e82c474",
					TenantId:    "515211e9",
					Name:        "testappahhis",
					DisplayName: "testappahhis",
					AccessType:  "READ",
					CallbackUrl: "",
					Status:      "APPROVED",
					AppFamily:   "default",
					CompanyId:   "a94f75e2-69b0-44af-8776-155df7c7d22e",
					DeveloperId: "",
					ParentId:    "a94f75e2-69b0-44af-8776-155df7c7d22e",
					Type:        "COMPANY",
					CreatedAt:   "2017-11-02 16:00:16.504+00:00",
					CreatedBy:   "haoming@apid.git",
					UpdatedAt:   "2017-11-02 16:00:16.504+00:00",
					UpdatedBy:   "haoming@apid.git",
				}

				results := [][]common.App{
					{expectedDevApp},
					{expectedDevApp},
					{expectedDevApp},
					{expectedDevApp},
					{expectedComApp},
					{expectedDevApp},
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
					nil,
					nil,
					nil,
					nil,
				}

				for i, data := range testData {
					priKey, priVal, secKey, secVal, org := data[0], data[1], data[2], data[3], data[4]
					apps, err := dbMan.GetApps(org, priKey, priVal, secKey, secVal)
					Expect(err).Should(Succeed())
					if len(results[i]) > 0 {
						Expect(apps).Should(Equal(results[i]))
					} else {
						Expect(apps).Should(BeZero())
					}
				}
			})

			It("should get Companies", func() {
				testData := [][]string{
					//positive tests
					{IdentifierAppId, "35608afe-2715-4064-bb4d-3cbb4e82c474", "", "", "apid-haoming"},
					{IdentifierCompanyName, "testcompanyhflxv", "", "", "apid-haoming"},
					{IdentifierConsumerKey, "wxyz", "", "", "apid-haoming"},
					// negative tests
					{IdentifierAppId, "35608afe-2715-4064-bb4d-3cbb4e82c474", "", "", "non-existent"},
					{IdentifierAppId, "non-existent", "", "", "apid-haoming"},
					{IdentifierCompanyName, "non-existent", "", "", "apid-haoming"},
					{IdentifierConsumerKey, "non-existent", "", "", "apid-haoming"},
					// SQL Injection
					{IdentifierAppId, "35608afe-2715-4064-bb4d-3cbb4e82c474", "", "", sqlInjectionStmt},
					{IdentifierAppId, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierCompanyName, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierConsumerKey, sqlInjectionStmt, "", "", "apid-haoming"},
				}

				var expectedCom = common.Company{
					Id:          "a94f75e2-69b0-44af-8776-155df7c7d22e",
					TenantId:    "515211e9",
					Name:        "testcompanyhflxv",
					DisplayName: "testcompanyhflxv",
					Status:      "ACTIVE",
					CreatedAt:   "2017-11-02 16:00:16.287+00:00",
					CreatedBy:   "haoming@apid.git",
					UpdatedAt:   "2017-11-02 16:00:16.287+00:00",
					UpdatedBy:   "haoming@apid.git",
				}

				results := [][]common.Company{
					{expectedCom},
					{expectedCom},
					{expectedCom},
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
					priKey, priVal, secKey, secVal, org := data[0], data[1], data[2], data[3], data[4]
					apps, err := dbMan.GetCompanies(org, priKey, priVal, secKey, secVal)
					Expect(err).Should(Succeed())
					if len(results[i]) > 0 {
						Expect(apps).Should(Equal(results[i]))
					} else {
						Expect(apps).Should(BeZero())
					}
				}
			})

			It("should get developers", func() {
				testData := [][]string{
					//positive tests
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", "", "apid-haoming"},
					{IdentifierConsumerKey, "abcd", "", "", "apid-haoming"},
					{IdentifierDeveloperEmail, "bar@google.com", "", "", "apid-haoming"},
					{IdentifierDeveloperId, "e41f04e8-9d3f-470a-8bfd-c7939945896c", "", "", "apid-haoming"},

					// negative tests
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", "", "non-existent"},
					{IdentifierAppId, "non-existent", "", "", "apid-haoming"},
					{IdentifierConsumerKey, "non-existent", "", "", "apid-haoming"},
					{IdentifierDeveloperEmail, "non-existent", "", "", "apid-haoming"},
					{IdentifierDeveloperId, "non-existent", "", "", "apid-haoming"},
					// SQL Injection
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", "", sqlInjectionStmt},
					{IdentifierAppId, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierConsumerKey, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierDeveloperEmail, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierDeveloperId, sqlInjectionStmt, "", "", "apid-haoming"},
				}

				var expectedDev = common.Developer{
					Id:                "e41f04e8-9d3f-470a-8bfd-c7939945896c",
					TenantId:          "515211e9",
					UserName:          "haoming",
					FirstName:         "haoming",
					LastName:          "zhang",
					Password:          "",
					Email:             "bar@google.com",
					Status:            "ACTIVE",
					EncryptedPassword: "",
					Salt:              "",
					CreatedAt:         "2017-08-16 22:39:46.669+00:00",
					CreatedBy:         "foo@google.com",
					UpdatedAt:         "2017-08-16 22:39:46.669+00:00",
					UpdatedBy:         "foo@google.com",
				}

				results := [][]common.Developer{
					{expectedDev},
					{expectedDev},
					{expectedDev},
					{expectedDev},
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
					priKey, priVal, secKey, secVal, org := data[0], data[1], data[2], data[3], data[4]
					prods, err := dbMan.GetDevelopers(org, priKey, priVal, secKey, secVal)
					Expect(err).Should(Succeed())
					if len(results[i]) > 0 {
						Expect(prods).Should(Equal(results[i]))
					} else {
						Expect(prods).Should(BeZero())
					}
				}
			})

			It("should get appCredentials", func() {
				testData := [][]string{
					// positive tests
					{IdentifierConsumerKey, "abcd", "", "", "apid-haoming"},
					{IdentifierAppId, "408ad853-3fa0-402f-90ee-103de98d71a5", "", "", "apid-haoming"},
					// negative tests
					{IdentifierConsumerKey, "abcd", "", "", "non-existent"},
					{IdentifierConsumerKey, "non-existent", "", "", "apid-haoming"},
					{IdentifierAppId, "non-existent", "", "", "apid-haoming"},
					// SQL Injection
					{IdentifierConsumerKey, "abcd", "", "", sqlInjectionStmt},
					{IdentifierConsumerKey, sqlInjectionStmt, "", "", "apid-haoming"},
					{IdentifierAppId, sqlInjectionStmt, "", "", "apid-haoming"},
				}

				var expectedCred = common.AppCredential{
					Id:             "abcd",
					TenantId:       "515211e9",
					ConsumerSecret: "secret1",
					AppId:          "408ad853-3fa0-402f-90ee-103de98d71a5",
					MethodType:     "",
					Status:         "APPROVED",
					IssuedAt:       "2017-08-18 22:13:18.35+00:00",
					ExpiresAt:      "",
					AppStatus:      "",
					Scopes:         "{}",
					CreatedAt:      "2017-08-18 22:13:18.35+00:00",
					CreatedBy:      "-NA-",
					UpdatedAt:      "2017-08-18 22:13:18.352+00:00",
					UpdatedBy:      "-NA-",
				}

				results := [][]common.AppCredential{
					{expectedCred},
					{expectedCred},
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
				}

				for i, data := range testData {
					priKey, priVal, secKey, secVal, org := data[0], data[1], data[2], data[3], data[4]
					prods, err := dbMan.GetAppCredentials(org, priKey, priVal, secKey, secVal)
					Expect(err).Should(Succeed())
					if len(results[i]) > 0 {
						Expect(prods).Should(Equal(results[i]))
					} else {
						Expect(prods).Should(BeZero())
					}
				}
			})

			It("should get CompanyDevelopers", func() {
				testData := [][]string{
					// positive tests
					{IdentifierCompanyName, "testcompanyhflxv", "", "", "apid-haoming"},
					// negative tests
					{IdentifierCompanyName, "testcompanyhflxv", "", "", "non-existent"},
					{IdentifierCompanyName, "non-existent", "", "", "apid-haoming"},
					// SQL Injection
					{IdentifierCompanyName, "testcompanyhflxv", "", "", sqlInjectionStmt},
					{IdentifierCompanyName, sqlInjectionStmt, "", "", "apid-haoming"},
				}

				var expectedComDev = common.CompanyDeveloper{
					TenantId:    "515211e9",
					CompanyId:   "a94f75e2-69b0-44af-8776-155df7c7d22e",
					DeveloperId: "590f33bf-f05c-48c1-bb93-183759bd9ee1",
					Roles:       "admin",
					CreatedAt:   "2017-11-02 16:00:16.287+00:00",
					CreatedBy:   "haoming@apid.git",
					UpdatedAt:   "2017-11-02 16:00:16.287+00:00",
					UpdatedBy:   "haoming@apid.git",
				}

				results := [][]common.CompanyDeveloper{
					{expectedComDev},
					nil,
					nil,
					nil,
					nil,
				}

				for i, data := range testData {
					priKey, priVal, secKey, secVal, org := data[0], data[1], data[2], data[3], data[4]
					prods, err := dbMan.GetCompanyDevelopers(org, priKey, priVal, secKey, secVal)
					Expect(err).Should(Succeed())
					if len(results[i]) > 0 {
						Expect(prods).Should(Equal(results[i]))
					} else {
						Expect(prods).Should(BeZero())
					}
				}
			})

		})

		Describe("utils", func() {
			It("GetApiProductNamesByConsumerKey", func() {
				data := "abcd"
				expected := []string{"apstest"}
				Expect(dbMan.GetApiProductNames(data, TypeConsumerKey)).Should(Equal(expected))

				data = "408ad853-3fa0-402f-90ee-103de98d71a5"
				expected = []string{"apstest"}
				Expect(dbMan.GetApiProductNames(data, TypeApp)).Should(Equal(expected))
			})

			It("GetAppNames", func() {
				data := "a94f75e2-69b0-44af-8776-155df7c7d22e"
				expected := []string{"testappahhis"}
				Expect(dbMan.GetAppNames(data, TypeCompany)).Should(Equal(expected))

				data = "e41f04e8-9d3f-470a-8bfd-c7939945896c"
				expected = []string{"apstest"}
				Expect(dbMan.GetAppNames(data, TypeDeveloper)).Should(Equal(expected))
			})

			It("GetComNames", func() {
				data := "8ba5b747-5104-4a40-89ca-a0a51798fe34"
				expected := []string{"DevCompany"}
				Expect(dbMan.GetComNames(data, TypeCompany)).Should(Equal(expected))
				data = "590f33bf-f05c-48c1-bb93-183759bd9ee1"
				expected = []string{"testcompanyhflxv"}
				Expect(dbMan.GetComNames(data, TypeDeveloper)).Should(Equal(expected))
			})

			It("GetDevEmailByDevId", func() {
				data := "e41f04e8-9d3f-470a-8bfd-c7939945896c"
				expected := "bar@google.com"
				Expect(dbMan.GetDevEmailByDevId(data, "apid-haoming")).Should(Equal(expected))
			})

			It("GetStatus", func() {
				data := "e41f04e8-9d3f-470a-8bfd-c7939945896c"
				expected := "ACTIVE"
				Expect(dbMan.GetStatus(data, AppTypeDeveloper)).Should(Equal(expected))
				data = "8ba5b747-5104-4a40-89ca-a0a51798fe34"
				expected = "ACTIVE"
				Expect(dbMan.GetStatus(data, AppTypeCompany)).Should(Equal(expected))
			})

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
