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
	"encoding/json"
	"github.com/apid/apidApiMetadata/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	apiTestUrl = "http://127.0.0.1:9000"
)

var _ = Describe("API Tests", func() {
	var apiMan *ApiManager
	testCount := 0
	client := &http.Client{}
	dbMan := &DummyDbMan{}
	var attrs []common.Attribute
	var testId string
	clientGet := func(path string, pars map[string][]string) (int, []byte) {
		uri, err := url.Parse(apiTestUrl + path)
		Expect(err).Should(Succeed())
		query := url.Values(pars)
		uri.RawQuery = query.Encode()
		httpReq, err := http.NewRequest("GET", uri.String(), nil)
		Expect(err).Should(Succeed())
		res, err := client.Do(httpReq)
		Expect(err).Should(Succeed())
		defer res.Body.Close()
		responseBody, err := ioutil.ReadAll(res.Body)
		Expect(err).Should(Succeed())
		return res.StatusCode, responseBody
	}

	BeforeEach(func() {
		testCount++
		testId = "test-" + strconv.Itoa(testCount)
		apiMan = &ApiManager{
			DbMan:            dbMan,
			AccessEntityPath: AccessEntityPath + strconv.Itoa(testCount),
			apiInitialized:   false,
		}
		attrs = setAttrs(dbMan, testId)
		apiMan.InitAPI()
		time.Sleep(100 * time.Millisecond)
	})

	It("ApiProduct", func() {
		testProd := []common.ApiProduct{
			{
				Id:            testId,
				Name:          "apstest",
				DisplayName:   "apstest",
				Description:   "",
				ApiResources:  "{/**}",
				ApprovalType:  "AUTO",
				Scopes:        `{foo,bar}`,
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
			},
		}

		expected := ApiProductSuccessResponse{
			ApiProduct: &ApiProductDetails{
				ApiProxies:     []string{"aps", "perfBenchmark"},
				ApiResources:   []string{"/**"},
				ApprovalType:   testProd[0].ApprovalType,
				Attributes:     attrs,
				CreatedAt:      testProd[0].CreatedAt,
				CreatedBy:      testProd[0].CreatedBy,
				Description:    testProd[0].Description,
				DisplayName:    testProd[0].DisplayName,
				Environments:   []string{"prod", "test"},
				ID:             testProd[0].Id,
				LastModifiedAt: testProd[0].UpdatedAt,
				LastModifiedBy: testProd[0].UpdatedBy,
				Name:           testProd[0].Name,
				QuotaInterval:  testProd[0].QuotaInterval,
				QuotaLimit:     10000000,
				QuotaTimeUnit:  testProd[0].QuotaTimeUnit,
				Scopes:         []string{"foo", "bar"},
			},
			Organization:             "test-org",
			PrimaryIdentifierType:    IdentifierAppName,
			PrimaryIdentifierValue:   "test-app",
			SecondaryIdentifierType:  IdentifierDeveloperId,
			SecondaryIdentifierValue: "test-dev",
		}

		testData := [][]common.ApiProduct{
			testProd,
			nil,
			testProd,
			testProd,
		}

		testPars := []map[string][]string{
			// positive
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppName:      {"test-app"},
				IdentifierDeveloperId:  {"test-dev"},
			},
			// negative
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppName:      {"test-app"},
				IdentifierDeveloperId:  {"test-dev"},
			},
			{
				IdentifierAppName:     {"test-app"},
				IdentifierDeveloperId: {"test-dev"},
			},
			{
				IdentifierDeveloperId: {"test-dev"},
			},
		}

		results := [][]interface{}{
			{http.StatusOK, expected},
			{http.StatusNotFound, nil},
			{http.StatusBadRequest, nil},
			{http.StatusBadRequest, nil},
		}

		for i, data := range testData {
			dbMan.apiProducts = data
			code, body := clientGet(apiMan.AccessEntityPath+EndpointApiProduct, testPars[i])
			Expect(code).Should(Equal(results[i][0]))
			if results[i][1] != nil {
				var res ApiProductSuccessResponse
				Expect(json.Unmarshal(body, &res)).Should(Succeed())
				Expect(res).Should(Equal(results[i][1]))
			}
		}

	})

	It("Apps", func() {
		testApp := []common.App{
			{
				Id:          testId,
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
			},
		}

		testComApp := []common.App{
			{
				Id:          testId,
				TenantId:    "515211e9",
				Name:        "apstest",
				DisplayName: "apstest",
				AccessType:  "READ",
				CallbackUrl: "https://www.google.com",
				Status:      "APPROVED",
				AppFamily:   "default",
				CompanyId:   "a94f75e2-69b0-44af-8776-155df7c7d22e",
				DeveloperId: "",
				ParentId:    "a94f75e2-69b0-44af-8776-155df7c7d22e",
				Type:        "COMPANY",
				CreatedAt:   "2017-08-18 22:13:18.325+00:00",
				CreatedBy:   "haoming@apid.git",
				UpdatedAt:   "2017-08-18 22:13:18.325+00:00",
				UpdatedBy:   "haoming@apid.git",
			},
		}

		testProductNames := []string{"foo", "bar"}
		testStatus := "test-status"
		testCreds := []common.AppCredential{
			{
				Id:             testId,
				TenantId:       "515211e9",
				ConsumerSecret: "secret1",
				AppId:          testId,
				MethodType:     "GET",
				Status:         "APPROVED",
				IssuedAt:       "2017-08-18 22:13:18.35+00:00",
				ExpiresAt:      "2018-08-18 22:13:18.35+00:00",
				AppStatus:      "APPROVED",
				Scopes:         "{foo,bar}",
				CreatedAt:      "2017-08-18 22:13:18.35+00:00",
				CreatedBy:      "-NA-",
				UpdatedAt:      "2017-08-18 22:13:18.352+00:00",
				UpdatedBy:      "-NA-",
			},
		}
		expected := AppSuccessResponse{
			App: &AppDetails{
				AccessType:  testApp[0].AccessType,
				ApiProducts: testProductNames,
				AppCredentials: []*CredentialDetails{
					{
						ApiProductReferences: testProductNames,
						AppID:                testCreds[0].AppId,
						AppStatus:            testApp[0].Status,
						Attributes:           attrs,
						ConsumerKey:          testCreds[0].Id,
						ConsumerSecret:       testCreds[0].ConsumerSecret,
						ExpiresAt:            testCreds[0].ExpiresAt,
						IssuedAt:             testCreds[0].IssuedAt,
						MethodType:           testCreds[0].MethodType,
						Scopes:               []string{"foo", "bar"},
						Status:               testCreds[0].Status,
					},
				},
				AppFamily:       testApp[0].AppFamily,
				AppParentID:     testApp[0].ParentId,
				AppParentStatus: testStatus,
				AppType:         testApp[0].Type,
				Attributes:      attrs,
				CallbackUrl:     testApp[0].CallbackUrl,
				CreatedAt:       testApp[0].CreatedAt,
				CreatedBy:       testApp[0].CreatedBy,
				DisplayName:     testApp[0].DisplayName,
				Id:              testApp[0].Id,
				LastModifiedAt:  testApp[0].UpdatedAt,
				LastModifiedBy:  testApp[0].UpdatedBy,
				Name:            testApp[0].Name,
				Status:          testApp[0].Status,
			},
			Organization:             "test-org",
			PrimaryIdentifierType:    IdentifierAppName,
			PrimaryIdentifierValue:   "test-app",
			SecondaryIdentifierType:  IdentifierDeveloperId,
			SecondaryIdentifierValue: "test-dev",
		}

		expectedComApp := AppSuccessResponse{
			App: &AppDetails{
				AccessType:  testComApp[0].AccessType,
				ApiProducts: testProductNames,
				AppCredentials: []*CredentialDetails{
					{
						ApiProductReferences: testProductNames,
						AppID:                testCreds[0].AppId,
						AppStatus:            testComApp[0].Status,
						Attributes:           attrs,
						ConsumerKey:          testCreds[0].Id,
						ConsumerSecret:       testCreds[0].ConsumerSecret,
						ExpiresAt:            testCreds[0].ExpiresAt,
						IssuedAt:             testCreds[0].IssuedAt,
						MethodType:           testCreds[0].MethodType,
						Scopes:               []string{"foo", "bar"},
						Status:               testCreds[0].Status,
					},
				},
				AppFamily:       testComApp[0].AppFamily,
				AppParentID:     "testcompanyhflxv",
				AppParentStatus: testStatus,
				AppType:         testComApp[0].Type,
				Attributes:      attrs,
				CallbackUrl:     testComApp[0].CallbackUrl,
				CreatedAt:       testComApp[0].CreatedAt,
				CreatedBy:       testComApp[0].CreatedBy,
				DisplayName:     testComApp[0].DisplayName,
				Id:              testComApp[0].Id,
				LastModifiedAt:  testComApp[0].UpdatedAt,
				LastModifiedBy:  testComApp[0].UpdatedBy,
				Name:            testComApp[0].Name,
				Status:          testComApp[0].Status,
			},
			Organization:             "test-org",
			PrimaryIdentifierType:    IdentifierAppName,
			PrimaryIdentifierValue:   "test-app",
			SecondaryIdentifierType:  IdentifierDeveloperId,
			SecondaryIdentifierValue: "test-dev",
		}

		testData := [][]common.App{
			testApp,
			testComApp,
			nil,
			testApp,
			testApp,
		}

		testPars := []map[string][]string{
			// positive
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppName:      {"test-app"},
				IdentifierDeveloperId:  {"test-dev"},
			},
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppName:      {"test-app"},
				IdentifierDeveloperId:  {"test-dev"},
			},
			// negative
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppName:      {"test-app"},
				IdentifierDeveloperId:  {"test-dev"},
			},
			{
				IdentifierAppName:     {"test-app"},
				IdentifierDeveloperId: {"test-dev"},
			},
			{
				IdentifierDeveloperId: {"test-dev"},
			},
		}

		results := [][]interface{}{
			{http.StatusOK, expected},
			{http.StatusOK, expectedComApp},
			{http.StatusNotFound, nil},
			{http.StatusBadRequest, nil},
			{http.StatusBadRequest, nil},
		}

		dbMan.apiProductNames = testProductNames
		dbMan.status = testStatus
		dbMan.appCredentials = testCreds
		for i, data := range testData {
			dbMan.apps = data
			dbMan.comNames = []string{"testcompanyhflxv"}
			code, body := clientGet(apiMan.AccessEntityPath+EndpointApp, testPars[i])
			Expect(code).Should(Equal(results[i][0]))
			if results[i][1] != nil {
				var res AppSuccessResponse
				Expect(json.Unmarshal(body, &res)).Should(Succeed())
				Expect(res).Should(Equal(results[i][1]))
			}
		}

	})

	It("Company", func() {
		testCom := []common.Company{
			{
				Id:          testId,
				TenantId:    "515211e9",
				Name:        "testcompanyhflxv",
				DisplayName: "testcompanyhflxv",
				Status:      "ACTIVE",
				CreatedAt:   "2017-11-02 16:00:16.287+00:00",
				CreatedBy:   "haoming@apid.git",
				UpdatedAt:   "2017-11-02 16:00:16.287+00:00",
				UpdatedBy:   "haoming@apid.git",
			},
		}

		testAppNames := []string{"foo", "bar"}

		expected := CompanySuccessResponse{
			Company: &CompanyDetails{
				Apps:           testAppNames,
				Attributes:     attrs,
				CreatedAt:      testCom[0].CreatedAt,
				CreatedBy:      testCom[0].CreatedBy,
				DisplayName:    testCom[0].DisplayName,
				ID:             testCom[0].Id,
				LastModifiedAt: testCom[0].UpdatedAt,
				LastModifiedBy: testCom[0].UpdatedBy,
				Name:           testCom[0].Name,
				Status:         testCom[0].Status,
			},
			Organization:           "test-org",
			PrimaryIdentifierType:  IdentifierAppId,
			PrimaryIdentifierValue: "test-app",
		}

		testData := [][]common.Company{
			testCom,
			nil,
			testCom,
			testCom,
		}

		testPars := []map[string][]string{
			// positive
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppId:        {"test-app"},
			},
			// negative
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppId:        {"test-app"},
			},
			{
				IdentifierAppId: {"test-app"},
			},
			{
				IdentifierDeveloperId: {"test-dev"},
			},
		}

		results := [][]interface{}{
			{http.StatusOK, expected},
			{http.StatusNotFound, nil},
			{http.StatusBadRequest, nil},
			{http.StatusBadRequest, nil},
		}

		dbMan.appNames = testAppNames

		for i, data := range testData {
			dbMan.companies = data
			code, body := clientGet(apiMan.AccessEntityPath+EndpointCompany, testPars[i])
			Expect(code).Should(Equal(results[i][0]))
			if results[i][1] != nil {
				var res CompanySuccessResponse
				Expect(json.Unmarshal(body, &res)).Should(Succeed())
				Expect(res).Should(Equal(results[i][1]))
			}
		}

	})

	It("Developer", func() {
		testDev := []common.Developer{
			{
				Id:                testId,
				TenantId:          "515211e9",
				UserName:          "haoming",
				FirstName:         "haoming",
				LastName:          "zhang",
				Password:          "111",
				Email:             "bar@google.com",
				Status:            "ACTIVE",
				EncryptedPassword: "222",
				Salt:              "333",
				CreatedAt:         "2017-08-16 22:39:46.669+00:00",
				CreatedBy:         "foo@google.com",
				UpdatedAt:         "2017-08-16 22:39:46.669+00:00",
				UpdatedBy:         "foo@google.com",
			},
		}

		testAppNames := []string{"foo", "bar"}
		testComNames := []string{"foo", "bar"}

		expected := DeveloperSuccessResponse{
			Developer: &DeveloperDetails{
				Apps:           testAppNames,
				Attributes:     attrs,
				Companies:      testComNames,
				CreatedAt:      testDev[0].CreatedAt,
				CreatedBy:      testDev[0].CreatedBy,
				Email:          testDev[0].Email,
				FirstName:      testDev[0].FirstName,
				ID:             testDev[0].Id,
				LastModifiedAt: testDev[0].UpdatedAt,
				LastModifiedBy: testDev[0].UpdatedBy,
				LastName:       testDev[0].LastName,
				Password:       testDev[0].Password,
				Status:         testDev[0].Status,
				UserName:       testDev[0].UserName,
			},
			Organization:           "test-org",
			PrimaryIdentifierType:  IdentifierAppId,
			PrimaryIdentifierValue: "test-app",
		}

		testData := [][]common.Developer{
			testDev,
			nil,
			testDev,
			testDev,
		}

		testPars := []map[string][]string{
			// positive
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppId:        {"test-app"},
			},
			// negative
			{
				IdentifierOrganization: {"test-org"},
				IdentifierAppId:        {"test-app"},
			},
			{
				IdentifierAppId: {"test-app"},
			},
			{
				IdentifierDeveloperId: {"test-dev"},
			},
		}

		results := [][]interface{}{
			{http.StatusOK, expected},
			{http.StatusNotFound, nil},
			{http.StatusBadRequest, nil},
			{http.StatusBadRequest, nil},
		}

		dbMan.appNames = testAppNames
		dbMan.comNames = testComNames

		for i, data := range testData {
			dbMan.developers = data
			code, body := clientGet(apiMan.AccessEntityPath+EndpointDeveloper, testPars[i])
			Expect(code).Should(Equal(results[i][0]))
			if results[i][1] != nil {
				var res DeveloperSuccessResponse
				Expect(json.Unmarshal(body, &res)).Should(Succeed())
				Expect(res).Should(Equal(results[i][1]))
			}
		}

	})

	It("AppCredential", func() {
		testAppCred := []common.AppCredential{
			{
				Id:             testId,
				TenantId:       "515211e9",
				ConsumerSecret: "secret1",
				AppId:          testId,
				MethodType:     "GET",
				Status:         "APPROVED",
				IssuedAt:       "2017-08-18 22:13:18.35+00:00",
				ExpiresAt:      "2018-08-18 22:13:18.35+00:00",
				AppStatus:      "APPROVED",
				Scopes:         "{foo,bar}",
				CreatedAt:      "2017-08-18 22:13:18.35+00:00",
				CreatedBy:      "-NA-",
				UpdatedAt:      "2017-08-18 22:13:18.352+00:00",
				UpdatedBy:      "-NA-",
			},
		}
		testApp := []common.App{
			{
				Id:          testId,
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
			},
		}

		testProductNames := []string{"foo", "bar"}
		testStatus := "test-status"

		expected := AppCredentialSuccessResponse{
			AppCredential: &AppCredentialDetails{
				AppID:       testAppCred[0].AppId,
				AppName:     testApp[0].Name,
				Attributes:  attrs,
				ConsumerKey: testAppCred[0].Id,
				ConsumerKeyStatus: &ConsumerKeyStatusDetails{
					AppCredential: &CredentialDetails{
						ApiProductReferences: testProductNames,
						AppID:                testAppCred[0].AppId,
						AppStatus:            testApp[0].Status,
						Attributes:           attrs,
						ConsumerKey:          testAppCred[0].Id,
						ConsumerSecret:       testAppCred[0].ConsumerSecret,
						ExpiresAt:            testAppCred[0].ExpiresAt,
						IssuedAt:             testAppCred[0].IssuedAt,
						MethodType:           testAppCred[0].MethodType,
						Scopes:               []string{"foo", "bar"},
						Status:               testAppCred[0].Status,
					},
					AppID:           testAppCred[0].AppId,
					AppName:         testApp[0].Name,
					AppStatus:       testApp[0].Status,
					AppType:         testApp[0].Type,
					DeveloperID:     testApp[0].DeveloperId,
					DeveloperStatus: testStatus,
					IsValidKey:      true,
				},
				ConsumerSecret: testAppCred[0].ConsumerSecret,
				DeveloperID:    testApp[0].DeveloperId,
				RedirectUris:   []string{testApp[0].CallbackUrl},
				Scopes:         []string{"foo", "bar"},
				Status:         testAppCred[0].Status,
			},
			Organization:           "test-org",
			PrimaryIdentifierType:  IdentifierConsumerKey,
			PrimaryIdentifierValue: "test-key",
		}

		testData := [][]common.AppCredential{
			testAppCred,
			nil,
			testAppCred,
			testAppCred,
		}

		testPars := []map[string][]string{
			// positive
			{
				IdentifierOrganization: {"test-org"},
				IdentifierConsumerKey:  {"test-key"},
			},
			// negative
			{
				IdentifierOrganization: {"test-org"},
				IdentifierConsumerKey:  {"test-key"},
			},
			{
				IdentifierConsumerKey: {"test-key"},
			},
			{
				IdentifierDeveloperId: {"test-dev"},
			},
		}

		results := [][]interface{}{
			{http.StatusOK, expected},
			{http.StatusNotFound, nil},
			{http.StatusBadRequest, nil},
			{http.StatusBadRequest, nil},
		}

		dbMan.apiProductNames = testProductNames
		dbMan.status = testStatus
		dbMan.apps = testApp
		for i, data := range testData {
			dbMan.appCredentials = data
			code, body := clientGet(apiMan.AccessEntityPath+EndpointAppCredentials, testPars[i])
			Expect(code).Should(Equal(results[i][0]))
			if results[i][1] != nil {
				var res AppCredentialSuccessResponse
				Expect(json.Unmarshal(body, &res)).Should(Succeed())
				Expect(res).Should(Equal(results[i][1]))
			}
		}
	})

	It("CompanyDeveloper", func() {
		testDev := []common.CompanyDeveloper{
			{
				TenantId:    "515211e9",
				CompanyId:   "a94f75e2-69b0-44af-8776-155df7c7d22e",
				DeveloperId: "590f33bf-f05c-48c1-bb93-183759bd9ee1",
				Roles:       "{foo,bar}",
				CreatedAt:   "2017-11-02 16:00:16.287+00:00",
				CreatedBy:   "haoming@apid.git",
				UpdatedAt:   "2017-11-02 16:00:16.287+00:00",
				UpdatedBy:   "haoming@apid.git",
			},
		}

		testComNames := []string{"foo"}
		testEmail := "haoming@apid.git"

		expected := CompanyDevelopersSuccessResponse{
			CompanyDevelopers: []*CompanyDeveloperDetails{
				{
					CompanyName:    testComNames[0],
					CreatedAt:      testDev[0].CreatedAt,
					CreatedBy:      testDev[0].CreatedBy,
					DeveloperEmail: testEmail,
					LastModifiedAt: testDev[0].UpdatedAt,
					LastModifiedBy: testDev[0].UpdatedBy,
					Roles:          []string{"foo", "bar"},
				},
			},
			Organization:           "test-org",
			PrimaryIdentifierType:  IdentifierCompanyName,
			PrimaryIdentifierValue: "test-com",
		}

		testData := [][]common.CompanyDeveloper{
			testDev,
			nil,
			testDev,
			testDev,
		}

		testPars := []map[string][]string{
			// positive
			{
				IdentifierOrganization: {"test-org"},
				IdentifierCompanyName:  {"test-com"},
			},
			// negative
			{
				IdentifierOrganization: {"test-org"},
				IdentifierCompanyName:  {"test-com"},
			},
			{
				IdentifierCompanyName: {"test-com"},
			},
			{
				IdentifierDeveloperId: {"test-dev"},
			},
		}

		results := [][]interface{}{
			{http.StatusOK, expected},
			{http.StatusNotFound, nil},
			{http.StatusBadRequest, nil},
			{http.StatusBadRequest, nil},
		}

		dbMan.comNames = testComNames
		dbMan.email = testEmail

		for i, data := range testData {
			dbMan.companyDevelopers = data
			code, body := clientGet(apiMan.AccessEntityPath+EndpointCompanyDeveloper, testPars[i])
			Expect(code).Should(Equal(results[i][0]))
			if results[i][1] != nil {
				var res CompanyDevelopersSuccessResponse
				Expect(json.Unmarshal(body, &res)).Should(Succeed())
				Expect(res).Should(Equal(results[i][1]))
			}
		}

	})
})

func setAttrs(dbMan *DummyDbMan, id string) []common.Attribute {
	dbMan.attrs = map[string][]common.Attribute{
		id: {
			{
				Name:  "foo",
				Value: "bar",
			},
			{
				Name:  "bar",
				Value: "foo",
			},
		},
	}
	return dbMan.attrs[id]
}
