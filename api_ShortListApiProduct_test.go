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

import "testing"

var shortListApiProductTestData = []struct {
	testDesc                           string
	req                                VerifyApiKeyRequest
	dbData                             []ApiProductDetails
	expectedResult                     string
	expectedWhenValidateProxyEnvIsTrue string
}{
	{
		testDesc:                           "single-product-happy-path",
		req:                                VerifyApiKeyRequest{EnvironmentName: "test", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
		dbData:                             []ApiProductDetails{ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}}},
		expectedResult:                     "api-product-1",
		expectedWhenValidateProxyEnvIsTrue: "api-product-1",
	},
	{
		testDesc: "multi-product-custom-resource-happy-path",
		req:      VerifyApiKeyRequest{EnvironmentName: "test", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
		dbData: []ApiProductDetails{
			ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
			ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
		},
		expectedResult:                     "api-product-2",
		expectedWhenValidateProxyEnvIsTrue: "api-product-2",
	}, {
		testDesc: "multi-product-only-one-matches-env",
		req:      VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
		dbData: []ApiProductDetails{
			ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
			ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
			ApiProductDetails{Id: "api-product-3", Environments: []string{"test", "prod", "stage"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
		},
		expectedResult:                     "api-product-3",
		expectedWhenValidateProxyEnvIsTrue: "api-product-3",
	},
	{
		testDesc: "multi-product-match-with-no-env",
		req:      VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
		dbData: []ApiProductDetails{
			ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
			ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
			ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
		},
		expectedResult:                     "api-product-3",
		expectedWhenValidateProxyEnvIsTrue: "api-product-3",
	},
	{
		testDesc: "multi-product-match-env",
		req:      VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
		dbData: []ApiProductDetails{
			ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
			ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod", "stage"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/this-is-my-path"}},
			ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
		},

		expectedResult:                     "api-product-2",
		expectedWhenValidateProxyEnvIsTrue: "api-product-2",
	},

	{
		testDesc: "multi-product-match-empty-res-env-proxy",
		req:      VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
		dbData: []ApiProductDetails{
			ApiProductDetails{Id: "api-product-1"},
			ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod", "stage"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/this-is-my-path"}},
			ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
		},

		expectedResult:                     "api-product-1",
		expectedWhenValidateProxyEnvIsTrue: "api-product-1",
	},

	{
		testDesc: "multi-product-match-empty-res-env-proxy-second-indexed",
		req:      VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
		dbData: []ApiProductDetails{
			ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
			ApiProductDetails{Id: "api-product-2"},
			ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
		},
		expectedResult:                     "api-product-2",
		expectedWhenValidateProxyEnvIsTrue: "api-product-2",
	},
	{
		testDesc: "multi-product-with-no-resource-match",
		req:      VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
		dbData: []ApiProductDetails{
			ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
			ApiProductDetails{Id: "api-product-2"},
			ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/b/**"}},
		},
		expectedResult:                     "api-product-2",
		expectedWhenValidateProxyEnvIsTrue: "api-product-2",
	},
	{
		testDesc: "multi-product-non-existent-proxy",
		req:      VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-non-exisitent-proxy", UriPath: "/this-is-my-path"},
		dbData: []ApiProductDetails{
			ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
			ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
			ApiProductDetails{Id: "api-product-3", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/b/**"}},
		},
		expectedResult:                     "api-product-2",
		expectedWhenValidateProxyEnvIsTrue: "",
	},
}

func TestShortListApiProduct(t *testing.T) {
	for _, td := range shortListApiProductTestData {
		actual := shortListApiProduct(td.dbData, td.req)
		if actual.Id != td.expectedResult {
			t.Errorf("TestData (%s) ValidateProxyEnv (%t) : expected (%s), actual (%s)", td.testDesc, td.req.ValidateAgainstApiProxiesAndEnvs, td.expectedResult, actual.Id)
		}
	}
}
