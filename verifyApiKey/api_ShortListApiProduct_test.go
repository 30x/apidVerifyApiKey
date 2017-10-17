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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("verifyApiKey_shortListApiProduct", func() {

	Context("shortListApiProduct tests", func() {
		It("single-product-happy-path", func() {
			td := shortListApiProductTestDataStruct{
				req:            VerifyApiKeyRequest{EnvironmentName: "test", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
				dbData:         []ApiProductDetails{ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}}},
				expectedResult: "api-product-1",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})
		It("multi-product-custom-resource-happy-path", func() {
			td := shortListApiProductTestDataStruct{

				req: VerifyApiKeyRequest{EnvironmentName: "test", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
				dbData: []ApiProductDetails{
					ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
					ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
				},
				expectedResult: "api-product-2",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})
		It("multi-product-only-one-matches-env", func() {
			td := shortListApiProductTestDataStruct{

				req: VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
				dbData: []ApiProductDetails{
					ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
					ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
					ApiProductDetails{Id: "api-product-3", Environments: []string{"test", "prod", "stage"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
				},
				expectedResult: "api-product-3",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})
		It("multi-product-match-with-no-env", func() {
			td := shortListApiProductTestDataStruct{

				req: VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
				dbData: []ApiProductDetails{
					ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
					ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
					ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
				},
				expectedResult: "api-product-3",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})
		It("multi-product-match-env", func() {
			td := shortListApiProductTestDataStruct{

				req: VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
				dbData: []ApiProductDetails{
					ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
					ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod", "stage"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/this-is-my-path"}},
					ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
				},

				expectedResult: "api-product-2",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})

		It("multi-product-match-empty-res-env-proxy", func() {
			td := shortListApiProductTestDataStruct{

				req: VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
				dbData: []ApiProductDetails{
					ApiProductDetails{Id: "api-product-1"},
					ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod", "stage"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/this-is-my-path"}},
					ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
				},

				expectedResult: "api-product-1",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})

		It("multi-product-match-empty-res-env-proxy-second-indexed", func() {
			td := shortListApiProductTestDataStruct{

				req: VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
				dbData: []ApiProductDetails{
					ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
					ApiProductDetails{Id: "api-product-2"},
					ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
				},
				expectedResult: "api-product-2",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})
		It("multi-product-with-no-resource-match", func() {
			td := shortListApiProductTestDataStruct{

				req: VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-proxy", UriPath: "/this-is-my-path"},
				dbData: []ApiProductDetails{
					ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
					ApiProductDetails{Id: "api-product-2"},
					ApiProductDetails{Id: "api-product-3", Environments: []string{}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/b/**"}},
				},
				expectedResult: "api-product-2",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})
		It("multi-product-non-existent-proxy", func() {
			td := shortListApiProductTestDataStruct{

				req: VerifyApiKeyRequest{EnvironmentName: "stage", ApiProxyName: "test-non-exisitent-proxy", UriPath: "/this-is-my-path"},
				dbData: []ApiProductDetails{
					ApiProductDetails{Id: "api-product-1", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/a/**"}},
					ApiProductDetails{Id: "api-product-2", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/**"}},
					ApiProductDetails{Id: "api-product-3", Environments: []string{"test", "prod"}, Apiproxies: []string{"test-proxy"}, Resources: []string{"/b/**"}},
				},
				expectedResult: "api-product-2",
			}

			actual := shortListApiProduct(td.dbData, td.req)
			Expect(actual.Id).Should(Equal(td.expectedResult))
		})
	})

})

type shortListApiProductTestDataStruct struct {
	req            VerifyApiKeyRequest
	dbData         []ApiProductDetails
	expectedResult string
}
