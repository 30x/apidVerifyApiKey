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
	"encoding/json"
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type performValidationsTestDataStruct struct {
	testDesc                           string
	dataWrapper                        VerifyApiKeyRequestResponseDataWrapper
	expectedResult                     string
	expectedWhenValidateProxyEnvIsTrue string
}

var _ = Describe("performValidationsTest", func() {

	apid.Initialize(factory.DefaultServicesFactory())
	log = factory.DefaultServicesFactory().Log()
	a := apiManager{}

	Context("performValidationsTest tests", func() {
		It("happy-path", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "",
				expectedWhenValidateProxyEnvIsTrue: "",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "ACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{
							Id:           "test-api-product",
							Resources:    []string{"/**"},
							Apiproxies:   []string{"test-proxy-name"},
							Environments: []string{"test-env-name"},
							Status:       "APPROVED",
						},
						App: AppDetails{
							Status: "APPROVED",
						},
						ClientId: ClientIdDetails{
							Status: "APPROVED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})
		It("Inactive Developer", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "{\"response_code\":\"keymanagement.service.DeveloperStatusNotActive\",\"response_message\":\"API Key verify failed for (test-key, test-org)\"}",
				expectedWhenValidateProxyEnvIsTrue: "{\"response_code\":\"keymanagement.service.DeveloperStatusNotActive\",\"response_message\":\"API Key verify failed for (test-key, test-org)\"}",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "INACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{
							Id:           "test-api-product",
							Resources:    []string{"/**"},
							Apiproxies:   []string{"test-proxy-name"},
							Environments: []string{"test-env-name"},
							Status:       "APPROVED",
						},
						App: AppDetails{
							Status: "APPROVED",
						},
						ClientId: ClientIdDetails{
							Status: "APPROVED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})
		It("Revoked Client Id", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "{\"response_code\":\"oauth.v2.ApiKeyNotApproved\",\"response_message\":\"API Key verify failed for (test-key, test-org)\"}",
				expectedWhenValidateProxyEnvIsTrue: "{\"response_code\":\"oauth.v2.ApiKeyNotApproved\",\"response_message\":\"API Key verify failed for (test-key, test-org)\"}",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "ACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{
							Id:           "test-api-product",
							Resources:    []string{"/**"},
							Apiproxies:   []string{"test-proxy-name"},
							Environments: []string{"test-env-name"},
							Status:       "APPROVED",
						},
						App: AppDetails{
							Status: "APPROVED",
						},
						ClientId: ClientIdDetails{
							Status: "REVOKED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})
		It("Revoked App", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "{\"response_code\":\"keymanagement.service.invalid_client-app_not_approved\",\"response_message\":\"API Key verify failed for (test-key, test-org)\"}",
				expectedWhenValidateProxyEnvIsTrue: "{\"response_code\":\"keymanagement.service.invalid_client-app_not_approved\",\"response_message\":\"API Key verify failed for (test-key, test-org)\"}",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "ACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{
							Id:           "test-api-product",
							Resources:    []string{"/**"},
							Apiproxies:   []string{"test-proxy-name"},
							Environments: []string{"test-env-name"},
							Status:       "APPROVED",
						},
						App: AppDetails{
							Status: "REVOKED",
						},
						ClientId: ClientIdDetails{
							Status: "APPROVED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})
		It("Company Inactive", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "{\"response_code\":\"keymanagement.service.CompanyStatusNotActive\",\"response_message\":\"API Key verify failed for (test-key, test-org)\"}",
				expectedWhenValidateProxyEnvIsTrue: "{\"response_code\":\"keymanagement.service.CompanyStatusNotActive\",\"response_message\":\"API Key verify failed for (test-key, test-org)\"}",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					ctype: "company",
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "INACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{
							Id:           "test-api-product",
							Resources:    []string{"/**"},
							Apiproxies:   []string{"test-proxy-name"},
							Environments: []string{"test-env-name"},
							Status:       "APPROVED",
						},
						App: AppDetails{
							Status: "APPROVED",
						},
						ClientId: ClientIdDetails{
							Status: "APPROVED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})
		It("Product not resolved", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "{\"response_code\":\"oauth.v2.InvalidApiKeyForGivenResource\",\"response_message\":\"Path Validation Failed. Product not resolved\"}",
				expectedWhenValidateProxyEnvIsTrue: "{\"response_code\":\"oauth.v2.InvalidApiKeyForGivenResource\",\"response_message\":\"Path Validation Failed. Product not resolved\"}",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "ACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{},
						App: AppDetails{
							Status: "APPROVED",
						},
						ClientId: ClientIdDetails{
							Status: "APPROVED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})
		It("resources not configured in db", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "",
				expectedWhenValidateProxyEnvIsTrue: "",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "ACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{
							Id:           "test-api-product",
							Resources:    []string{},
							Apiproxies:   []string{"test-proxy-name"},
							Environments: []string{"test-env-name"},
							Status:       "APPROVED",
						},
						App: AppDetails{
							Status: "APPROVED",
						},
						ClientId: ClientIdDetails{
							Status: "APPROVED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})
		It("proxies not configured in db", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "",
				expectedWhenValidateProxyEnvIsTrue: "",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "ACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{
							Id:           "test-api-product",
							Resources:    []string{"/test"},
							Apiproxies:   []string{},
							Environments: []string{"test-env-name"},
							Status:       "APPROVED",
						},
						App: AppDetails{
							Status: "APPROVED",
						},
						ClientId: ClientIdDetails{
							Status: "APPROVED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})
		It("environments not configured in db", func() {
			td := performValidationsTestDataStruct{
				expectedResult:                     "",
				expectedWhenValidateProxyEnvIsTrue: "",
				dataWrapper: VerifyApiKeyRequestResponseDataWrapper{
					verifyApiKeyRequest: VerifyApiKeyRequest{
						Key:              "test-key",
						OrganizationName: "test-org",
						UriPath:          "/test",
						ApiProxyName:     "test-proxy-name",
						EnvironmentName:  "test-env-name",
					},
					tempDeveloperDetails: DeveloperDetails{
						Status: "ACTIVE",
					},
					verifyApiKeySuccessResponse: VerifyApiKeySuccessResponse{
						ApiProduct: ApiProductDetails{
							Id:           "test-api-product",
							Resources:    []string{"/test"},
							Apiproxies:   []string{"test-proxy-name"},
							Environments: []string{},
							Status:       "APPROVED",
						},
						App: AppDetails{
							Status: "APPROVED",
						},
						ClientId: ClientIdDetails{
							Status: "APPROVED",
						},
					},
				},
			}
			actualObject := a.performValidations(td.dataWrapper)
			var actual string
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))

			td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
			actualObject = a.performValidations(td.dataWrapper)
			if actualObject != nil {
				a, _ := json.Marshal(&actualObject)
				actual = string(a)
			} else {
				actual = ""
			}
			Expect(actual).Should(Equal(td.expectedResult))
		})

	})
})
