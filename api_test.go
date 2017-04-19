package apidVerifyApiKey

import (
	"encoding/json"
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
	})
})