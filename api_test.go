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

	Context("verifyAPIKey() directly", func() {

		It("should reject a bad key", func() {
			v := url.Values{
				"key": []string{"credential_x"},
				"uriPath": []string{"/test"},
				"environment": []string{"Env_0"},
				"organization": []string{"Org_0"},
				"action": []string{"verify"},
			}
			rsp, err := verifyAPIKey(v)
			Expect(err).ShouldNot(HaveOccurred())

			var respj kmsResponseFail
			json.Unmarshal(rsp, &respj)
			Expect(respj.Type).Should(Equal("ErrorResult"))
			Expect(respj.ErrInfo.ErrorCode).Should(Equal("REQ_ENTRY_NOT_FOUND"))

		})
		/*
			It("should reject a key once it's deleted", func() {
				pd0 := &DataPayload{
					EntityIdentifier: "app_credential_0",
				}
				res := deleteCredential(*pd0, db)
				Expect(res).Should(BeTrue())

				var respj kmsResponseFail
				v := url.Values{
					"key": []string{"app_credential_0"},
					"path": []string{"/test"},
					"env": []string{"Env_0"},
					"organization": []string{"Org_0"},
					"action": []string{"verify"},
				}
				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())

				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("ErrorResult"))
				Expect(respj.ErrInfo.ErrorCode).Should(Equal("REQ_ENTRY_NOT_FOUND"))
			})
		*/
		It("should successfully verify good keys", func() {
			for i := 1; i < 10; i++ {
				resulti := strconv.FormatInt(int64(i), 10)
				v := url.Values{
					"key": []string{"app_credential_"+resulti},
					"uriPath": []string{"/test"},
					"environment": []string{"Env_0"},
					"organization": []string{"Org_0"},
					"action": []string{"verify"},
				}
				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())

				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Key).Should(Equal("app_credential_" + resulti))
			}
		})
	})

	Context("access via API", func() {

		It("should reject a bad key", func() {

			uri, err := url.Parse(testServer.URL)
			uri.Path = apiPath

			v := url.Values{}
			v.Add("organization", "Org_0")
			v.Add("key", "credential_x")
			v.Add("environment", "Env_0")
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

		It("should successfully verify a good key", func() {

			uri, err := url.Parse(testServer.URL)
			uri.Path = apiPath

			v := url.Values{}
			v.Add("organization", "Org_0")
			v.Add("key", "app_credential_1")
			v.Add("environment", "Env_0")
			v.Add("uriPath", "/test")
			v.Add("action", "verify")

			client := &http.Client{}
			req, err := http.NewRequest("POST", uri.String(), strings.NewReader(v.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

			res, err := client.Do(req)
			defer res.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			var respj kmsResponseSuccess
			body, err := ioutil.ReadAll(res.Body)
			Expect(err).ShouldNot(HaveOccurred())
			json.Unmarshal(body, &respj)
			Expect(respj.Type).Should(Equal("APIKeyContext"))
			Expect(respj.RspInfo.Key).Should(Equal("app_credential_1"))
		})
	})
})
