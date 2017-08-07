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
	"testing"
)

var performValidationsTestData = []struct {
	testDesc                           string
	dataWrapper                        VerifyApiKeyRequestResponseDataWrapper
	expectedResult                     string
	expectedWhenValidateProxyEnvIsTrue string
}{
	{
		testDesc:                           "happy-path",
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
	}, {
		testDesc:                           "Inactive Developer",
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
	},
	{
		testDesc:                           "Revoked Client Id",
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
	},
	{
		testDesc:                           "Revoked App",
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
	},
	{
		testDesc:                           "Company Inactive",
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
	},
	{
		testDesc:                           "Product not resolved",
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
	},
}

func TestPerformValidation(t *testing.T) {

	// tODO : what is the right way to get this ?
	apid.Initialize(factory.DefaultServicesFactory())
	log = factory.DefaultServicesFactory().Log()
	a := apiManager{}
	for _, td := range performValidationsTestData {
		actual, err := a.performValidations(td.dataWrapper)
		if string(actual) != td.expectedResult {
			t.Errorf("TestData (%s) ValidateProxyEnv (%t) : expected (%s), actual (%s)", td.testDesc, td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs, td.expectedResult, string(actual))
		}
		if err != nil && err.Error() != "200" {
			t.Error("Expected to return 200 status code")
		}
	}
}

func TestPerformValidationValidateProxyEnv(t *testing.T) {

	// tODO : what is the right way to get this ?
	apid.Initialize(factory.DefaultServicesFactory())
	log = factory.DefaultServicesFactory().Log()
	a := apiManager{}
	for _, td := range performValidationsTestData {
		td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs = true
		actual, err := a.performValidations(td.dataWrapper)
		if string(actual) != td.expectedWhenValidateProxyEnvIsTrue {
			t.Errorf("TestData (%s) ValidateProxyEnv (%t) : expected (%s), actual (%s)", td.testDesc, td.dataWrapper.verifyApiKeyRequest.ValidateAgainstApiProxiesAndEnvs, td.expectedWhenValidateProxyEnvIsTrue, string(actual))
		}
		if err != nil && err.Error() != "200" {
			t.Error("Expected to return 200 status code")
		}
	}
}
