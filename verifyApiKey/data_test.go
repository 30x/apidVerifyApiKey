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
package verifyApiKey

import (
	"github.com/apid/apid-core"
	"github.com/apid/apid-core/factory"
	"github.com/apid/apidVerifyApiKey/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"sync"
)

var _ = Describe("DataTest", func() {

	Context("query Db to get api key details", func() {
		var dataTestTempDir string
		var dbMan *DbManager
		var _ = BeforeEach(func() {
			var err error
			dataTestTempDir, err = ioutil.TempDir(testTempDirBase, "sqlite3")
			Expect(err).NotTo(HaveOccurred())
			s := factory.DefaultServicesFactory()
			apid.Initialize(s)
			config := apid.Config()
			config.Set("local_storage_path", dataTestTempDir)
			common.SetApidServices(s, s.Log())

			dbMan = &DbManager{
				DbManager: common.DbManager{
					Data:  s.Data(),
					DbMux: sync.RWMutex{},
				},
			}
			dbMan.SetDbVersion(dataTestTempDir)

		})

		It("should get company getApiKeyDetails for happy path", func() {
			setupApikeyCompanyTestDb(dbMan.Db)

			dataWrapper := VerifyApiKeyRequestResponseDataWrapper{
				verifyApiKeyRequest: VerifyApiKeyRequest{
					OrganizationName: "apigee-mcrosrvc-client0001",
					Key:              "63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0",
				},
			}
			err := dbMan.getApiKeyDetails(&dataWrapper)
			Expect(err).NotTo(HaveOccurred())

			Expect(dataWrapper.ctype).Should(BeEquivalentTo("company"))
			Expect(dataWrapper.tenant_id).Should(BeEquivalentTo("bc811169"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.ClientId.Status).Should(BeEquivalentTo("APPROVED"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientSecret).Should(BeEquivalentTo("Ui8dcyGW3lA04YdX"))

			Expect(dataWrapper.tempDeveloperDetails.Id).Should(BeEquivalentTo("7834c683-9453-4389-b816-34ca24dfccd9"))
			Expect(dataWrapper.tempDeveloperDetails.UserName).Should(BeEquivalentTo("East India Company"))
			Expect(dataWrapper.tempDeveloperDetails.FirstName).Should(BeEquivalentTo("DevCompany"))
			Expect(dataWrapper.tempDeveloperDetails.LastName).Should(BeEquivalentTo(""))
			Expect(dataWrapper.tempDeveloperDetails.Email).Should(BeEquivalentTo(""))
			Expect(dataWrapper.tempDeveloperDetails.Status).Should(BeEquivalentTo("ACTIVE"))
			Expect(dataWrapper.tempDeveloperDetails.CreatedAt).Should(BeEquivalentTo("2017-08-05 19:54:12.359+00:00"))
			Expect(dataWrapper.tempDeveloperDetails.CreatedBy).Should(BeEquivalentTo("defaultUser"))
			Expect(dataWrapper.tempDeveloperDetails.LastmodifiedAt).Should(BeEquivalentTo("2017-08-05 19:54:12.359+00:00"))
			Expect(dataWrapper.tempDeveloperDetails.LastmodifiedBy).Should(BeEquivalentTo("defaultUser"))

			Expect(dataWrapper.verifyApiKeySuccessResponse.App.Id).Should(BeEquivalentTo("d371f05a-7c04-430c-b12d-26cf4e4d5d65"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.Name).Should(BeEquivalentTo("CompApp2"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.AccessType).Should(BeEquivalentTo("READ"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl).Should(BeEquivalentTo("www.apple.com"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.DisplayName).Should(BeEquivalentTo(""))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.Status).Should(BeEquivalentTo("APPROVED"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.AppFamily).Should(BeEquivalentTo("default"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.Company).Should(BeEquivalentTo("7834c683-9453-4389-b816-34ca24dfccd9"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.CreatedAt).Should(BeEquivalentTo("2017-08-07 17:00:54.25+00:00"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.CreatedBy).Should(BeEquivalentTo("defaultUser"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.LastmodifiedAt).Should(BeEquivalentTo("2017-08-07 17:09:08.259+00:00"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.LastmodifiedBy).Should(BeEquivalentTo("defaultUser"))

		})

		It("should get developer ApiKeyDetails - happy path", func() {
			setupApikeyDeveloperTestDb(dbMan.Db)

			dataWrapper := VerifyApiKeyRequestResponseDataWrapper{
				verifyApiKeyRequest: VerifyApiKeyRequest{
					OrganizationName: "apigee-mcrosrvc-client0001",
					Key:              "63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0",
				},
			}
			err := dbMan.getApiKeyDetails(&dataWrapper)
			Expect(err).NotTo(HaveOccurred())

			Expect(dataWrapper.ctype).Should(BeEquivalentTo("developer"))
			Expect(dataWrapper.tenant_id).Should(BeEquivalentTo("bc811169"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.ClientId.Status).Should(BeEquivalentTo("APPROVED"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientSecret).Should(BeEquivalentTo("Ui8dcyGW3lA04YdX"))

			Expect(dataWrapper.tempDeveloperDetails.Id).Should(BeEquivalentTo("209ffd18-37e9-4a67-9e30-a5c40a534b6c"))
			Expect(dataWrapper.tempDeveloperDetails.UserName).Should(BeEquivalentTo("wilson"))
			Expect(dataWrapper.tempDeveloperDetails.FirstName).Should(BeEquivalentTo("Woodre"))
			Expect(dataWrapper.tempDeveloperDetails.LastName).Should(BeEquivalentTo("Wilson"))
			Expect(dataWrapper.tempDeveloperDetails.Email).Should(BeEquivalentTo("developer@apigee.com"))
			Expect(dataWrapper.tempDeveloperDetails.Status).Should(BeEquivalentTo("ACTIVE"))
			Expect(dataWrapper.tempDeveloperDetails.CreatedAt).Should(BeEquivalentTo("2017-08-08 17:24:09.008+00:00"))
			Expect(dataWrapper.tempDeveloperDetails.CreatedBy).Should(BeEquivalentTo("defaultUser"))
			Expect(dataWrapper.tempDeveloperDetails.LastmodifiedAt).Should(BeEquivalentTo("2017-08-08 17:24:09.008+00:00"))
			Expect(dataWrapper.tempDeveloperDetails.LastmodifiedBy).Should(BeEquivalentTo("defaultUser"))

			Expect(dataWrapper.verifyApiKeySuccessResponse.App.Id).Should(BeEquivalentTo("d371f05a-7c04-430c-b12d-26cf4e4d5d65"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.Name).Should(BeEquivalentTo("DeveloperApp"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.AccessType).Should(BeEquivalentTo("READ"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl).Should(BeEquivalentTo("www.apple.com"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.DisplayName).Should(BeEquivalentTo(""))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.Status).Should(BeEquivalentTo("APPROVED"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.AppFamily).Should(BeEquivalentTo("default"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.Company).Should(BeEquivalentTo(""))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.CreatedAt).Should(BeEquivalentTo("2017-08-07 17:00:54.25+00:00"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.CreatedBy).Should(BeEquivalentTo("defaultUser"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.LastmodifiedAt).Should(BeEquivalentTo("2017-08-07 17:09:08.259+00:00"))
			Expect(dataWrapper.verifyApiKeySuccessResponse.App.LastmodifiedBy).Should(BeEquivalentTo("defaultUser"))

		})

		It("should throw error when apikey not found", func() {

			setupApikeyCompanyTestDb(dbMan.Db)
			dataWrapper := VerifyApiKeyRequestResponseDataWrapper{
				verifyApiKeyRequest: VerifyApiKeyRequest{
					OrganizationName: "apigee-mcrosrvc-client0001",
					Key:              "invalid-Jkcc6GENVWGT1Zw5gek7kVJ0",
				},
			}
			err := dbMan.getApiKeyDetails(&dataWrapper)
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).Should(BeEquivalentTo("InvalidApiKey"))
		})

		It("should get api products ", func() {

			setupApikeyCompanyTestDb(dbMan.Db)

			apiProducts := dbMan.getApiProductsForApiKey("63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0", "bc811169")
			Expect(len(apiProducts)).Should(BeEquivalentTo(1))

			Expect(apiProducts[0].Id).Should(BeEquivalentTo("24987a63-edb9-4d6b-9334-87e1d70df8e3"))
			Expect(apiProducts[0].Name).Should(BeEquivalentTo("KeyProduct4"))
			Expect(apiProducts[0].DisplayName).Should(BeEquivalentTo("Sandbox Diamond"))
			Expect(apiProducts[0].Status).Should(BeEquivalentTo(""))
			Expect(apiProducts[0].QuotaTimeunit).Should(BeEquivalentTo(""))
			Expect(apiProducts[0].QuotaInterval).Should(BeEquivalentTo(0))
			Expect(apiProducts[0].QuotaLimit).Should(BeEquivalentTo(""))

			Expect(apiProducts[0].Resources).Should(BeEquivalentTo([]string{"/zoho", "/twitter", "/nike"}))
			Expect(apiProducts[0].Apiproxies).Should(BeEquivalentTo([]string{"DevApplication", "KeysApplication"}))
			Expect(apiProducts[0].Environments).Should(BeEquivalentTo([]string{"test"}))
			Expect(apiProducts[0].Company).Should(BeEquivalentTo(""))
			Expect(len(apiProducts[0].Attributes)).Should(BeEquivalentTo(0))

			Expect(apiProducts[0].CreatedBy).Should(BeEquivalentTo("defaultUser"))
			Expect(apiProducts[0].CreatedAt).Should(BeEquivalentTo("2017-08-08 02:53:32.726+00:00"))
			Expect(apiProducts[0].LastmodifiedBy).Should(BeEquivalentTo("defaultUser"))
			Expect(apiProducts[0].LastmodifiedAt).Should(BeEquivalentTo("2017-08-08 02:53:32.726+00:00"))

		})

		It("should return empty array when no api products found", func() {

			setupApikeyCompanyTestDb(dbMan.Db)
			apiProducts := dbMan.getApiProductsForApiKey("invalid-LKJkcc6GENVWGT1Zw5gek7kVJ0", "bc811169")
			Expect(len(apiProducts)).Should(BeEquivalentTo(0))

		})

		It("should get kms attributes", func() {

			setupKmsAttributesdata(dbMan.Db)
			attributes := dbMan.GetKmsAttributes("bc811169", "40753e12-a50a-429d-9121-e571eb4e43a9", "85629786-37c5-4e8c-bb45-208f3360d005", "50321842-d6ee-4e92-91b9-37234a7920c1", "test-invalid")
			Expect(len(attributes)).Should(BeEquivalentTo(3))
			Expect(len(attributes["40753e12-a50a-429d-9121-e571eb4e43a9"])).Should(BeEquivalentTo(1))
			Expect(len(attributes["85629786-37c5-4e8c-bb45-208f3360d005"])).Should(BeEquivalentTo(2))
			Expect(len(attributes["50321842-d6ee-4e92-91b9-37234a7920c1"])).Should(BeEquivalentTo(5))
			Expect(len(attributes["test-invalid"])).Should(BeEquivalentTo(0))

		})

	})
})
