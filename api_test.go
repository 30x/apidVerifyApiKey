package apidVerifyApiKey

import (
	"encoding/json"
	"github.com/apigee-labs/transicator/common"
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

		It("Positive DB test for Insert operations", func() {
			db := getDB()
			txn, err := db.Begin()
			Expect(err).ShouldNot(HaveOccurred())
			// api products
			for i := 0; i < 10; i++ {
				var rows []common.Row
				result := strconv.FormatInt(int64(i), 10)
				row := common.Row{
					"id": {
						Value: "api_product_" + result,
					},
					"api_resources": {
						Value: "{/**, /test}",
					},
					"environments": {
						Value: "{Env_0, Env_1}",
					},
					"tenant_id": {
						Value: "tenant_id_xxxx",
					},
					"_change_selector": {
						Value: "Org_0",
					},
				}
				rows = append(rows, row)
				res := insertAPIproducts(rows, txn)
				Expect(res).Should(BeTrue())
			}
			// developers
			for i := 0; i < 10; i++ {
				var rows []common.Row
				result := strconv.FormatInt(int64(i), 10)
				row := common.Row{
					"id": {
						Value: "developer_id_" + result,
					},
					"status": {
						Value: "Active",
					},
					"email": {
						Value: "test@apigee.com",
					},
					"first_name": {
						Value: "Apigee",
					},
					"last_name": {
						Value: "Google",
					},
					"tenant_id": {
						Value: "tenant_id_xxxx",
					},
					"_change_selector": {
						Value: "Org_0",
					},
				}

				rows = append(rows, row)
				res := insertDevelopers(rows, txn)
				Expect(res).Should(BeTrue())
			}
			// companies
			for i := 0; i < 10; i++ {
				var rows []common.Row
				result := strconv.FormatInt(int64(i), 10)
				row := common.Row{
					"id": {
						Value: "company_id_" + result,
					},
					"status": {
						Value: "Active",
					},
					"name": {
						Value: "Apigee Corporation",
					},
					"display_name": {
						Value: "Apigee",
					},
					"tenant_id": {
						Value: "tenant_id_xxxx",
					},
					"_change_selector": {
						Value: "Org_0",
					},
				}

				rows = append(rows, row)
				res := insertCompanies(rows, txn)
				Expect(res).Should(BeTrue())
			}

			// company developers
			for i := 0; i < 10; i++ {
				var rows []common.Row
				result := strconv.FormatInt(int64(i), 10)
				row := common.Row{
					"developer_id": {
						Value: "developer_id_" + result,
					},
					"tenant_id": {
						Value: "tenant_id_0",
					},
					"_change_selector": {
						Value: "test_org0",
					},
					"company_id": {
						Value: "company_id_" + result,
					},
				}

				rows = append(rows, row)
				res := insertCompanyDevelopers(rows, txn)
				Expect(res).Should(BeTrue())
			}

			// application
			var j, k int
			for i := 0; i < 10; i++ {
				resulti := strconv.FormatInt(int64(i), 10)
				for j = k; j < 10+k; j++ {
					var rows []common.Row
					resultj := strconv.FormatInt(int64(j), 10)
					row := common.Row{
						"id": {
							Value: "application_id_" + resultj,
						},
						"developer_id": {
							Value: "developer_id_" + resulti,
						},
						"status": {
							Value: "Approved",
						},
						"tenant_id": {
							Value: "tenant_id_xxxx",
						},
						"callback_url": {
							Value: "http://apigee.com",
						},
						"_change_selector": {
							Value: "Org_0",
						},
					}

					rows = append(rows, row)
					res := insertApplications(rows, txn)
					Expect(res).Should(BeTrue())
				}
				k = j
			}
			// app credentials
			for i := 0; i < 10; i++ {
				var rows []common.Row
				result := strconv.FormatInt(int64(i), 10)
				row := common.Row{
					"id": {
						Value: "app_credential_" + result,
					},
					"app_id": {
						Value: "application_id_" + result,
					},
					"status": {
						Value: "Approved",
					},
					"tenant_id": {
						Value: "tenant_id_xxxx",
					},
					"callback_url": {
						Value: "http://apigee.com",
					},
					"_change_selector": {
						Value: "Org_0",
					},
				}

				rows = append(rows, row)
				res := insertCredentials(rows, txn)
				Expect(res).Should(BeTrue())
			}
			// api product mapper
			for i := 0; i < 10; i++ {
				var rows []common.Row
				result := strconv.FormatInt(int64(i), 10)
				var s = "hello";
				row := common.Row{
					"apiprdt_id": {
						Value: "api_product_" + result,
					},
					"status": {
						Value: "Approved",
					},
					"app_id": {
						Value: "application_id_" + result,
					},
					"appcred_id": {
						Value: "app_credential_" + result,
					},
					"tenant_id": {
						Value: "tenant_id_xxxx",
					},
					"_change_selector": {
						Value: "Org_0",
					},
				}

				rows = append(rows, row)
				res := insertAPIProductMappers(rows, txn)
				Expect(res).Should(BeTrue())
			}

			txn.Commit()
			var count int64
			db.QueryRow("select count(*) from data_scope").Scan(&count)
			log.Info("Found ", count)

		})

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

		It("should successfully verify good keys", func() {
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
				Expect(respj.RspInfo.Key).Should(Equal("app_credential_" + resulti))
			}
		})

		It("Positive DB test for Delete operations", func() {
			db := getDB()
			txn, err := db.Begin()
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 10; i++ {
				result := strconv.FormatInt(int64(i), 10)
				row := common.Row{
					"apiprdt_id": {
						Value: "api_product_" + result,
					},
					"status": {
						Value: "Approved",
					},
					"app_id": {
						Value: "application_id_" + result,
					},
					"appcred_id": {
						Value: "app_credential_" + result,
					},
					"tenant_id": {
						Value: "tenant_id_xxxx",
					},
					"_change_selector": {
						Value: "Org_0",
					},
				}

				res := deleteAPIproductMapper(row, txn)
				Expect(res).Should(BeTrue())
			}

			for i := 0; i < 10; i++ {
				srvItems := common.Row{}
				result := strconv.FormatInt(int64(i), 10)

				scv := &common.ColumnVal{
					Value: "app_credential_" + result,
					Type:  3,
				}
				srvItems["id"] = scv

				scv = &common.ColumnVal{
					Value: "Org_0",
					Type:  3,
				}
				srvItems["_change_selector"] = scv

				res := deleteObject("APP_CREDENTIAL", srvItems, txn)
				Expect(res).Should(BeTrue())
			}
			for i := 0; i < 100; i++ {

				srvItems := common.Row{}
				resultj := strconv.FormatInt(int64(i), 10)

				scv := &common.ColumnVal{
					Value: "application_id_" + resultj,
					Type:  1,
				}
				srvItems["id"] = scv

				scv = &common.ColumnVal{
					Value: "Org_0",
					Type:  1,
				}
				srvItems["_change_selector"] = scv

				res := deleteObject("APP", srvItems, txn)
				Expect(res).Should(BeTrue())
			}

			for i := 0; i < 10; i++ {
				srvItems := common.Row{}
				result := strconv.FormatInt(int64(i), 10)

				scv := &common.ColumnVal{
					Value: "developer_id_" + result,
					Type:  1,
				}
				srvItems["id"] = scv

				scv = &common.ColumnVal{
					Value: "Org_0",
					Type:  1,
				}
				srvItems["_change_selector"] = scv

				res := deleteObject("DEVELOPER", srvItems, txn)
				Expect(res).Should(BeTrue())
			}

			for i := 0; i < 10; i++ {
				srvItems := common.Row{}
				result := strconv.FormatInt(int64(i), 10)

				scv := &common.ColumnVal{
					Value: "api_product_" + result,
					Type:  1,
				}
				srvItems["id"] = scv

				scv = &common.ColumnVal{
					Value: "Org_0",
					Type:  1,
				}
				srvItems["_change_selector"] = scv

				res := deleteObject("API_PRODUCT", srvItems, txn)
				Expect(res).Should(BeTrue())
			}

			txn.Commit()

		})

		It("Negative cases for DB Deletes on KMS tables", func() {
			db := getDB()
			txn, err := db.Begin()
			Expect(err).ShouldNot(HaveOccurred())

			srvItems := common.Row{}
			result := "DEADBEEF"

			scv := &common.ColumnVal{
				Value: "api_product_" + result,
				Type:  3,
			}
			srvItems["apiprdt_id"] = scv

			scv = &common.ColumnVal{
				Value: "application_id_" + result,
				Type:  3,
			}
			srvItems["app_id"] = scv

			scv = &common.ColumnVal{
				Value: "app_credential_" + result,
				Type:  3,
			}
			srvItems["appcred_id"] = scv

			scv = &common.ColumnVal{
				Value: "Org_0",
				Type:  3,
			}
			srvItems["_change_selector"] = scv

			res := deleteAPIproductMapper(srvItems, txn)
			Expect(res).Should(BeFalse())

			res = deleteObject("API_PRODUCT", srvItems, txn)
			Expect(res).Should(BeFalse())

			res = deleteObject("APP_CREDENTIAL", srvItems, txn)
			Expect(res).Should(BeFalse())

			res = deleteObject("DEVELOPER", srvItems, txn)
			Expect(res).Should(BeFalse())

			res = deleteObject("APP", srvItems, txn)
			Expect(res).Should(BeFalse())

			txn.Rollback()

		})
		It("Negative cases for DB Inserts/updates on KMS tables", func() {

			db := getDB()
			txn, err := db.Begin()
			Expect(err).ShouldNot(HaveOccurred())

			var rows []common.Row
			srvItems := common.Row{}
			result := "NOPRODID_BADCASE"
			scv := &common.ColumnVal{
				Value: "foobar_" + result,
				Type:  1,
			}
			srvItems[result] = scv

			scv = &common.ColumnVal{
				Value: "{/**, /test}",
				Type:  1,
			}
			srvItems["api_resources"] = scv

			scv = &common.ColumnVal{
				Value: "{Env_1, Env_0}",
				Type:  1,
			}
			srvItems["environments"] = scv

			scv = &common.ColumnVal{
				Value: "Org_0",
				Type:  1,
			}
			srvItems["_change_selector"] = scv

			scv = &common.ColumnVal{
				Value: "tenant_id_xxxx",
				Type:  1,
			}
			srvItems["tenant_id"] = scv

			rows = append(rows, srvItems)
			res := insertAPIproducts(rows, txn)
			Expect(res).Should(BeFalse())

			res = insertApplications(rows, txn)
			Expect(res).Should(BeFalse())

			res = insertCredentials(rows, txn)
			Expect(res).Should(BeFalse())

			res = insertAPIProductMappers(rows, txn)
			Expect(res).Should(BeFalse())

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
