package apidVerifyApiKey

import (
	"database/sql"
	"encoding/json"
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
	. "github.com/30x/apidApigeeSync" // for direct access to Payload types
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var _ = Describe("api", func() {

	var tmpDir string
	var db *sql.DB
	var server *httptest.Server

	BeforeSuite(func() {
		apid.Initialize(factory.DefaultServicesFactory())

		config := apid.Config()

		var err error
		tmpDir, err = ioutil.TempDir("", "api_test")
		Expect(err).NotTo(HaveOccurred())

		config.Set("data_path", tmpDir)

		// init() will create the tables
		apid.InitializePlugins()

		db, err = apid.Data().DB()
		Expect(err).NotTo(HaveOccurred())
		insertTestData(db)

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == apiPath {
				handleRequest(w, req)
			}
		}))
	})

	AfterSuite(func() {
		apid.Events().Close()
		server.Close()
		os.RemoveAll(tmpDir)
	})

	Context("verifyAPIKey() directly", func() {

		It("should reject a bad key", func() {
			rsp, err := verifyAPIKey("credential_x", "/test", "Env_0", "Org_0", "verify")
			Expect(err).ShouldNot(HaveOccurred())

			var respj kmsResponseFail
			json.Unmarshal(rsp, &respj)
			Expect(respj.Type).Should(Equal("ErrorResult"))
			Expect(respj.ErrInfo.ErrorCode).Should(Equal("REQ_ENTRY_NOT_FOUND"))

		})

		It("should reject a key once it's deleted", func() {
			pd0 := &DataPayload{
				EntityIdentifier: "credential_0",
			}
			res := deleteCredential(*pd0, db, "Org_0")
			Expect(res).Should(BeTrue())

			var respj kmsResponseFail
			rsp, err := verifyAPIKey("credential_0", "/test", "Env_0", "Org_0", "verify")
			Expect(err).ShouldNot(HaveOccurred())

			json.Unmarshal(rsp, &respj)
			Expect(respj.Type).Should(Equal("ErrorResult"))
			Expect(respj.ErrInfo.ErrorCode).Should(Equal("REQ_ENTRY_NOT_FOUND"))
		})

		It("should successfully verify good keys", func() {
			for i := 1; i < 10; i++ {
				resulti := strconv.FormatInt(int64(i), 10)
				rsp, err := verifyAPIKey("credential_"+resulti, "/test", "Env_0", "Org_0", "verify")
				Expect(err).ShouldNot(HaveOccurred())

				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Key).Should(Equal("credential_" + resulti))
			}
		})
	})

	Context("access via API", func() {

		It("should reject a bad key", func() {

			uri, err := url.Parse(server.URL)
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

			uri, err := url.Parse(server.URL)
			uri.Path = apiPath

			v := url.Values{}
			v.Add("organization", "Org_0")
			v.Add("key", "credential_1")
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
			Expect(respj.RspInfo.Key).Should(Equal("credential_1"))
		})
	})
})

func insertTestData(db *sql.DB) {

	for i := 0; i < 10; i++ {
		result := strconv.FormatInt(int64(i), 10)
		pd0 := &DataPayload{
			PldCont: Payload{
				AppName:      "Api_product_" + result,
				Resources:    []string{"/**", "/test"},
				Environments: []string{"Env_0", "Env_1"},
			},
		}

		res := insertAPIproduct(*pd0, db, "Org_0")
		Expect(res).Should(BeTrue())
	}

	for i := 0; i < 10; i++ {
		result := strconv.FormatInt(int64(i), 10)

		pd1 := &DataPayload{
			EntityIdentifier: "developer_id_" + result,
			PldCont: Payload{
				Email:     "person_0@apigee.com",
				Status:    "Active",
				UserName:  "user_0",
				FirstName: "user_first_name0",
				LastName:  "user_last_name0",
			},
		}

		res := insertCreateDeveloper(*pd1, db, "Org_0")
		Expect(res).Should(BeTrue())
	}

	var j, k int
	for i := 0; i < 10; i++ {
		resulti := strconv.FormatInt(int64(i), 10)
		for j = k; j < 10+k; j++ {
			resultj := strconv.FormatInt(int64(j), 10)
			pd2 := &DataPayload{
				EntityIdentifier: "application_id_" + resultj,
				PldCont: Payload{
					Email:       "person_0@apigee.com",
					Status:      "Approved",
					AppName:     "application_id_" + resultj,
					DeveloperId: "developer_id_" + resulti,
					CallbackUrl: "call_back_url_0",
				},
			}

			res := insertCreateApplication(*pd2, db, "Org_0")
			Expect(res).Should(BeTrue())
		}
		k = j
	}

	j = 0
	k = 0
	for i := 0; i < 10; i++ {
		resulti := strconv.FormatInt(int64(i), 10)
		for j = k; j < 10+k; j++ {
			resultj := strconv.FormatInt(int64(j), 10)
			pd3 := &DataPayload{
				EntityIdentifier: "credential_" + resultj,
				PldCont: Payload{
					AppId:          "application_id_" + resulti,
					Status:         "Approved",
					ConsumerSecret: "consumer_secret_0",
					IssuedAt:       349583485,
					ApiProducts:    []Apip{{ApiProduct: "Api_product_0", Status: "Approved"}},
				},
			}

			res := insertCreateCredential(*pd3, db, "Org_0")
			Expect(res).Should(BeTrue())
		}
		k = j
	}
}
