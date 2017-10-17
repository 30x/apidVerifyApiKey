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

// TODO: end to end IT tests
// 1. happy path for developer
// 2. happy path for company
// 3. error case for developer / company
// 4. input request validation error case
// 5. key not found case

import (
	"encoding/json"
	"errors"
	"github.com/apid/apid-core"
	"github.com/apid/apid-core/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var (
	testServer *httptest.Server
)

var _ = Describe("end to end tests", func() {
	var dataTestTempDir string
	var dbMan *DbManager

	var _ = BeforeEach(func() {
		var err error
		dataTestTempDir, err = ioutil.TempDir(testTempDirBase, "api_test_sqlite3")
		serviceFactoryForTest := factory.DefaultServicesFactory()
		apid.Initialize(serviceFactoryForTest)
		config := apid.Config()
		config.Set("data_path", testTempDir)
		config.Set("log_level", "DEBUG")
		serviceFactoryForTest.Config().Set("local_storage_path", dataTestTempDir)

		Expect(err).NotTo(HaveOccurred())

		dbMan = &DbManager{
			Data:  serviceFactoryForTest.Data(),
			DbMux: sync.RWMutex{},
		}
		dbMan.SetDbVersion(dataTestTempDir)

		apiMan := ApiManager{
			DbMan:             dbMan,
			VerifiersEndpoint: ApiPath,
		}

		testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == ApiPath {
				apiMan.HandleRequest(w, req)
			}
		}))

	})

	Context("veriifyApiKey Api test ", func() {
		It("should return validation error for missing input fields", func() {
			var respObj ErrorResponse
			reqInput := VerifyApiKeyRequest{
				Key: "test",
			}
			jsonBody, _ := json.Marshal(reqInput)

			responseBody, err := performTestOperation(string(jsonBody), 400)
			Expect(err).ShouldNot(HaveOccurred())

			json.Unmarshal(responseBody, &respObj)
			Expect(respObj.ResponseMessage).Should(Equal("Bad_REQUEST"))
			Expect(respObj.ResponseCode).Should(Equal("Missing mandatory fields in the request : action organizationName uriPath"))
		})
		It("should return validation error for inavlid key", func() {
			var respObj ErrorResponse
			reqInput := VerifyApiKeyRequest{
				Key:              "invalid-key",
				Action:           "verify",
				OrganizationName: "apigee-mcrosrvc-client0001",
				EnvironmentName:  "test",
				ApiProxyName:     "DevApplication",
				UriPath:          "/zoho",

				ValidateAgainstApiProxiesAndEnvs: true,
			}
			jsonBody, _ := json.Marshal(reqInput)

			responseBody, err := performTestOperation(string(jsonBody), 200)
			Expect(err).ShouldNot(HaveOccurred())

			json.Unmarshal(responseBody, &respObj)
			Expect(respObj.ResponseMessage).Should(Equal("API Key verify failed for (invalid-key, apigee-mcrosrvc-client0001)"))
			Expect(respObj.ResponseCode).Should(Equal("oauth.v2.InvalidApiKey"))
		})
		It("should return validation error for inavlid env", func() {
			setupApikeyDeveloperTestDb(dbMan.Db)
			var respObj ErrorResponse
			reqInput := VerifyApiKeyRequest{
				Key:              "63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0",
				Action:           "verify",
				OrganizationName: "apigee-mcrosrvc-client0001",
				EnvironmentName:  "prod",
				ApiProxyName:     "DevApplication",
				UriPath:          "/zoho",

				ValidateAgainstApiProxiesAndEnvs: true,
			}
			jsonBody, _ := json.Marshal(reqInput)

			responseBody, err := performTestOperation(string(jsonBody), 200)
			Expect(err).ShouldNot(HaveOccurred())

			json.Unmarshal(responseBody, &respObj)
			Expect(respObj.ResponseMessage).Should(Equal("ENV Validation Failed (test vs prod)"))
			Expect(respObj.ResponseCode).Should(Equal("oauth.v2.InvalidApiKeyForGivenResource"))
		})
		It("should return validation error for inavlid resource", func() {
			setupApikeyDeveloperTestDb(dbMan.Db)
			var respObj ErrorResponse
			reqInput := VerifyApiKeyRequest{
				Key:              "63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0",
				Action:           "verify",
				OrganizationName: "apigee-mcrosrvc-client0001",
				EnvironmentName:  "test",
				ApiProxyName:     "DevApplication",
				UriPath:          "/google",

				ValidateAgainstApiProxiesAndEnvs: true,
			}
			jsonBody, _ := json.Marshal(reqInput)

			responseBody, err := performTestOperation(string(jsonBody), 200)
			Expect(err).ShouldNot(HaveOccurred())

			json.Unmarshal(responseBody, &respObj)
			Expect(respObj.ResponseMessage).Should(Equal("Path Validation Failed. Product not resolved"))
			Expect(respObj.ResponseCode).Should(Equal("oauth.v2.InvalidApiKeyForGivenResource"))
		})
		It("should return validation error for inavlid proxies", func() {
			setupApikeyDeveloperTestDb(dbMan.Db)
			var respObj ErrorResponse
			reqInput := VerifyApiKeyRequest{
				Key:              "63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0",
				Action:           "verify",
				OrganizationName: "apigee-mcrosrvc-client0001",
				EnvironmentName:  "test",
				ApiProxyName:     "Invalid-proxy",
				UriPath:          "/zoho",

				ValidateAgainstApiProxiesAndEnvs: true,
			}
			jsonBody, _ := json.Marshal(reqInput)

			responseBody, err := performTestOperation(string(jsonBody), 200)
			Expect(err).ShouldNot(HaveOccurred())

			json.Unmarshal(responseBody, &respObj)
			Expect(respObj.ResponseMessage).Should(Equal("Proxy Validation Failed (DevApplication, KeysApplication vs Invalid-proxy)"))
			Expect(respObj.ResponseCode).Should(Equal("oauth.v2.InvalidApiKeyForGivenResource"))
		})
		It("should peform verify api key for developer happy path", func() {
			setupApikeyDeveloperTestDb(dbMan.Db)
			var respObj VerifyApiKeySuccessResponse

			reqInput := VerifyApiKeyRequest{
				Action:           "verify",
				OrganizationName: "apigee-mcrosrvc-client0001",
				Key:              "63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0",
				EnvironmentName:  "test",
				ApiProxyName:     "DevApplication",
				UriPath:          "/zoho",

				ValidateAgainstApiProxiesAndEnvs: true,
			}
			jsonBody, _ := json.Marshal(reqInput)

			responseBody, err := performTestOperation(string(jsonBody), 200)
			Expect(err).ShouldNot(HaveOccurred())

			json.Unmarshal(responseBody, &respObj)
			Expect(respObj.Developer.Id).Should(Equal("209ffd18-37e9-4a67-9e30-a5c40a534b6c"))
			Expect(respObj.Developer.FirstName).Should(Equal("Woodre"))
			Expect(respObj.Developer.CreatedAt).Should(Equal("2017-08-08 17:24:09.008+00:00"))
			Expect(respObj.Developer.LastmodifiedAt).Should(Equal("2017-08-08 17:24:09.008+00:00"))
			Expect(respObj.Developer.CreatedBy).Should(Equal("defaultUser"))
			Expect(respObj.Developer.LastmodifiedBy).Should(Equal("defaultUser"))
			Expect(len(respObj.Developer.Attributes)).Should(Equal(0))
			Expect(respObj.Developer.Company).Should(Equal(""))
			Expect(respObj.Developer.Status).Should(Equal("ACTIVE"))
			Expect(respObj.Developer.UserName).Should(Equal("wilson"))
			Expect(respObj.Developer.Email).Should(Equal("developer@apigee.com"))
			Expect(respObj.Developer.LastName).Should(Equal("Wilson"))
			Expect(len(respObj.Developer.Apps)).Should(Equal(0))

			Expect(respObj.ClientId.ClientId).Should(Equal("63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0"))
			Expect(respObj.ClientId.Status).Should(Equal("APPROVED"))
			Expect(respObj.ClientId.Attributes[0].Name).Should(Equal("Device"))
			Expect(respObj.ClientId.Attributes[0].Value).Should(Equal("ios"))
			Expect(respObj.ClientId.ClientSecret).Should(Equal("Ui8dcyGW3lA04YdX"))
			Expect(respObj.ClientId.RedirectURIs[0]).Should(Equal("www.apple.com"))

			Expect(respObj.Company.Id).Should(Equal(""))

			Expect(respObj.App.Id).Should(Equal("d371f05a-7c04-430c-b12d-26cf4e4d5d65"))

			Expect(respObj.ApiProduct.Id).Should(Equal("24987a63-edb9-4d6b-9334-87e1d70df8e3"))

			Expect(respObj.Environment).Should(Equal("test"))

		})

		It("should peform verify api key for company happy path", func() {
			setupApikeyCompanyTestDb(dbMan.Db)
			var respObj VerifyApiKeySuccessResponse

			reqInput := VerifyApiKeyRequest{
				Action:           "verify",
				OrganizationName: "apigee-mcrosrvc-client0001",
				Key:              "63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0",
				EnvironmentName:  "test",
				ApiProxyName:     "DevApplication",
				UriPath:          "/zoho",

				ValidateAgainstApiProxiesAndEnvs: true,
			}
			jsonBody, _ := json.Marshal(reqInput)

			responseBody, err := performTestOperation(string(jsonBody), 200)
			Expect(err).ShouldNot(HaveOccurred())

			json.Unmarshal(responseBody, &respObj)
			Expect(respObj.Developer.Id).Should(Equal(""))

			Expect(respObj.Company.Id).Should(Equal("7834c683-9453-4389-b816-34ca24dfccd9"))
			Expect(respObj.Company.Name).Should(Equal("DevCompany"))
			Expect(respObj.Company.CreatedAt).Should(Equal("2017-08-05 19:54:12.359+00:00"))
			Expect(respObj.Company.LastmodifiedAt).Should(Equal("2017-08-05 19:54:12.359+00:00"))
			Expect(respObj.Company.CreatedBy).Should(Equal("defaultUser"))
			Expect(respObj.Company.LastmodifiedBy).Should(Equal("defaultUser"))
			Expect(len(respObj.Company.Attributes)).Should(Equal(1))
			Expect(respObj.Company.Attributes[0].Name).Should(Equal("country"))
			Expect(respObj.Company.Attributes[0].Value).Should(Equal("england"))
			Expect(respObj.Company.DisplayName).Should(Equal("East India Company"))
			Expect(respObj.Company.Status).Should(Equal("ACTIVE"))
			Expect(len(respObj.Developer.Apps)).Should(Equal(0))

			Expect(respObj.ClientId.ClientId).Should(Equal("63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0"))
			Expect(respObj.ClientId.Status).Should(Equal("APPROVED"))
			Expect(len(respObj.ClientId.Attributes)).Should(Equal(0))
			Expect(respObj.ClientId.ClientSecret).Should(Equal("Ui8dcyGW3lA04YdX"))
			Expect(respObj.ClientId.RedirectURIs[0]).Should(Equal("www.apple.com"))

			Expect(respObj.Company.Id).Should(Equal("7834c683-9453-4389-b816-34ca24dfccd9"))

			Expect(respObj.App.Id).Should(Equal("d371f05a-7c04-430c-b12d-26cf4e4d5d65"))

			Expect(respObj.ApiProduct.Id).Should(Equal("24987a63-edb9-4d6b-9334-87e1d70df8e3"))

			Expect(respObj.Environment).Should(Equal("test"))
		})

	})
})

func performTestOperation(jsonBody string, expectedResponseCode int) ([]byte, error) {
	uri, err := url.Parse(testServer.URL)
	uri.Path = ApiPath
	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", uri.String(), strings.NewReader(string(jsonBody)))
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := client.Do(httpReq)
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != expectedResponseCode {
		err = errors.New("expected response status code does not match. Expected : " + strconv.Itoa(expectedResponseCode) + " ,actual : " + strconv.Itoa(res.StatusCode))
	}

	return responseBody, err
}
