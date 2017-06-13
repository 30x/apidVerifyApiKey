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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var _ = Describe("api", func() {

	Context("DB Inserts/Deletes verification", func() {

		It("should reject a bad key", func() {
			v := url.Values{
				"key":       []string{"credential_x"},
				"uriPath":   []string{"/test"},
				"scopeuuid": []string{"ABCDE"},
				"action":    []string{"verify"},
			}
			rsp, err := verifyAPIKey(v)
			Expect(err).ShouldNot(HaveOccurred())

			var respj kmsResponseFail
			json.Unmarshal(rsp, &respj)
			Expect(respj.Type).Should(Equal("ErrorResult"))
			Expect(respj.ErrInfo.ErrorCode).Should(Equal("REQ_ENTRY_NOT_FOUND"))

		})

		It("should successfully verify good Developer keys", func() {
			for i := 1; i < 10; i++ {
				resulti := strconv.FormatInt(int64(i), 10)
				v := url.Values{
					"key":       []string{"app_credential_" + resulti},
					"uriPath":   []string{"/test"},
					"scopeuuid": []string{"ABCDE"},
					"action":    []string{"verify"},
				}
				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())

				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Type).Should(Equal("developer"))
				Expect(respj.RspInfo.Key).Should(Equal("app_credential_" + resulti))
			}
		})

		It("should successfully verify good Company keys", func() {
			for i := 100; i < 110; i++ {
				resulti := strconv.FormatInt(int64(i), 10)
				v := url.Values{
					"key":       []string{"app_credential_" + resulti},
					"uriPath":   []string{"/test"},
					"scopeuuid": []string{"ABCDE"},
					"action":    []string{"verify"},
				}
				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())

				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Type).Should(Equal("company"))
				Expect(respj.RspInfo.Key).Should(Equal("app_credential_" + resulti))
			}
		})

		It("should reject a bad key", func() {

			uri, err := url.Parse(testServer.URL)
			uri.Path = apiPath

			v := url.Values{}
			v.Add("key", "credential_x")
			v.Add("scopeuuid", "ABCDE")
			v.Add("uriPath", "/test")
			v.Add("action", "verify")

			client := &http.Client{}
			req, err := http.NewRequest("POST", uri.String(), strings.NewReader(v.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

			res, err := client.Do(req)
			defer res.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			var respj kmsResponseFail
			body, err := ioutil.ReadAll(res.Body)
			Expect(err).ShouldNot(HaveOccurred())
			json.Unmarshal(body, &respj)
			Expect(respj.Type).Should(Equal("ErrorResult"))
			Expect(respj.ErrInfo.ErrorCode).Should(Equal("REQ_ENTRY_NOT_FOUND"))
		})

		It("should report error for no scopes", func() {
			v := url.Values{
				"key":       []string{"credential_x"},
				"uriPath":   []string{"/test"},
				"scopeuuid": []string{"ABCDE"},
				"action":    []string{"verify"},
			}

			clearDataScopeTable(getDB())
			rsp, err := verifyAPIKey(v)
			Expect(err).ShouldNot(HaveOccurred())

			var respj kmsResponseFail
			json.Unmarshal(rsp, &respj)
			Expect(respj.Type).Should(Equal("ErrorResult"))
			Expect(respj.ErrInfo.ErrorCode).Should(Equal("ENV_VALIDATION_FAILED"))

		})

		It("should report error for invalid requests", func() {
			v := url.Values{
				"key":       []string{"credential_x"},
				"uriPath":   []string{"/test"},
				"scopeuuid": []string{"ABCDE"},
				"action":    []string{"verify"},
			}

			fields := []string{"key", "uriPath", "scopeuuid", "action"}
			for _, field := range fields {
				tmp := v.Get(field)
				v.Del(field)

				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())
				var respj kmsResponseFail
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("ErrorResult"))
				Expect(respj.ErrInfo.ErrorCode).Should(Equal("INCORRECT_USER_INPUT"))

				v.Set(field, tmp)
			}
		})
	})
})

func clearDataScopeTable(db apid.DB) {
	txn, _ := db.Begin()
	txn.Exec("DELETE FROM EDGEX_DATA_SCOPE")
	log.Info("clear EDGEX_DATA_SCOPE for test")
	txn.Commit()
}
