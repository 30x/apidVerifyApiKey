package apidVerifyApiKey

import (
	"database/sql"
	"encoding/json"
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
	"github.com/30x/transicator/common"
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
		/*
			It("should reject a key once it's deleted", func() {
				pd0 := &DataPayload{
					EntityIdentifier: "app_credential_0",
				}
				res := deleteCredential(*pd0, db)
				Expect(res).Should(BeTrue())

				var respj kmsResponseFail
				rsp, err := verifyAPIKey("app_credential_0", "/test", "Env_0", "Org_0", "verify")
				Expect(err).ShouldNot(HaveOccurred())

				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("ErrorResult"))
				Expect(respj.ErrInfo.ErrorCode).Should(Equal("REQ_ENTRY_NOT_FOUND"))
			})
		*/
		It("should successfully verify good keys", func() {
			for i := 1; i < 10; i++ {
				resulti := strconv.FormatInt(int64(i), 10)
				rsp, err := verifyAPIKey("app_credential_"+resulti, "/test", "Env_0", "Org_0", "verify")
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

func insertTestData(db *sql.DB) {

	for i := 0; i < 10; i++ {
		var rows []common.Row
		srvItems := common.Row{}
		result := strconv.FormatInt(int64(i), 10)

		scv := &common.ColumnVal{
			Value: "api_product_" + result,
			Type:  1,
		}
		srvItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "{/**, /test}",
			Type:  1,
		}
		srvItems["api_resources"] = scv

		scv = &common.ColumnVal{
			Value: "{Env_0, Env_1}",
			Type:  1,
		}
		srvItems["environments"] = scv

		scv = &common.ColumnVal{
			Value: "Org_0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_xxxx",
			Type:  1,
		}
		srvItems["tenant_id"] = scv
		rows = append(rows, srvItems)
		res := insertAPIproducts(rows, db)
		Expect(res).Should(BeTrue())
	}

	for i := 0; i < 10; i++ {
		var rows []common.Row
		srvItems := common.Row{}
		result := strconv.FormatInt(int64(i), 10)

		scv := &common.ColumnVal{
			Value: "developer_id_" + result,
			Type:  1,
		}
		srvItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "test@apigee.com",
			Type:  1,
		}
		srvItems["email"] = scv

		scv = &common.ColumnVal{
			Value: "Active",
			Type:  1,
		}
		srvItems["status"] = scv

		scv = &common.ColumnVal{
			Value: "Apigee",
			Type:  1,
		}
		srvItems["firstName"] = scv

		scv = &common.ColumnVal{
			Value: "Google",
			Type:  1,
		}
		srvItems["lastName"] = scv

		scv = &common.ColumnVal{
			Value: "Org_0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_xxxx",
			Type:  1,
		}
		srvItems["tenant_id"] = scv

		rows = append(rows, srvItems)
		res := insertDevelopers(rows, db)
		Expect(res).Should(BeTrue())
	}

	var j, k int
	for i := 0; i < 10; i++ {
		resulti := strconv.FormatInt(int64(i), 10)
		for j = k; j < 10+k; j++ {
			var rows []common.Row

			srvItems := common.Row{}
			resultj := strconv.FormatInt(int64(j), 10)

			scv := &common.ColumnVal{
				Value: "application_id_" + resultj,
				Type:  1,
			}
			srvItems["id"] = scv

			scv = &common.ColumnVal{
				Value: "developer_id_" + resulti,
				Type:  1,
			}
			srvItems["developer_id"] = scv

			scv = &common.ColumnVal{
				Value: "approved",
				Type:  1,
			}
			srvItems["status"] = scv

			scv = &common.ColumnVal{
				Value: "http://apigee.com",
				Type:  1,
			}
			srvItems["callback_url"] = scv

			scv = &common.ColumnVal{
				Value: "Org_0",
				Type:  1,
			}
			srvItems["_apid_scope"] = scv

			scv = &common.ColumnVal{
				Value: "tenant_id_xxxx",
				Type:  1,
			}
			srvItems["tenant_id"] = scv
			rows = append(rows, srvItems)
			res := insertApplications(rows, db)
			Expect(res).Should(BeTrue())
		}
		k = j
	}

	for i := 0; i < 10; i++ {
		var rows []common.Row
		srvItems := common.Row{}
		result := strconv.FormatInt(int64(i), 10)

		scv := &common.ColumnVal{
			Value: "app_credential_" + result,
			Type:  1,
		}
		srvItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "application_id_" + result,
			Type:  1,
		}
		srvItems["app_id"] = scv

		scv = &common.ColumnVal{
			Value: "approved",
			Type:  1,
		}
		srvItems["status"] = scv

		scv = &common.ColumnVal{
			Value: "Org_0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_xxxx",
			Type:  1,
		}
		srvItems["tenant_id"] = scv
		rows = append(rows, srvItems)
		res := insertCredentials(rows, db)
		Expect(res).Should(BeTrue())
	}

	for i := 0; i < 10; i++ {
		var rows []common.Row
		srvItems := common.Row{}
		result := strconv.FormatInt(int64(i), 10)

		scv := &common.ColumnVal{
			Value: "api_product_" + result,
			Type:  1,
		}
		srvItems["apiprdt_id"] = scv

		scv = &common.ColumnVal{
			Value: "application_id_" + result,
			Type:  1,
		}
		srvItems["app_id"] = scv

		scv = &common.ColumnVal{
			Value: "app_credential_" + result,
			Type:  1,
		}
		srvItems["appcred_id"] = scv
		scv = &common.ColumnVal{
			Value: "approved",
			Type:  1,
		}
		srvItems["status"] = scv
		scv = &common.ColumnVal{
			Value: "Org_0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_xxxx",
			Type:  1,
		}
		srvItems["tenant_id"] = scv
		rows = append(rows, srvItems)
		res := insertAPIProductMappers(rows, db)
		Expect(res).Should(BeTrue())
	}

}
