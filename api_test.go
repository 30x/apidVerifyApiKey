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
				row := generateTestApiProduct(i)
				res := insertAPIproducts([]common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}
			// developers
			for i := 0; i < 10; i++ {
				row := generateTestDeveloper(i)
				res := insertDevelopers([]common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			// application
			var j, k int
			for i := 0; i < 10; i++ {
				for j = k; j < 10 + k; j++ {
					row := generateTestApp(j, i)
					res := insertApplications([]common.Row{row}, txn)
					Expect(res).Should(BeTrue())
				}
				k = j
			}
			// app credentials
			for i := 0; i < 10; i++ {
				row := generateTestAppCreds(i)
				res := insert("APP_CREDENTIAL", []common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}
			// api product mapper
			for i := 0; i < 10; i++ {
				row := generateTestApiProductMapper(i)
				res := insert("APP_CREDENTIAL_APIPRODUCT_MAPPER", []common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			// Following are data for company
			// api products
			for i := 100; i < 110; i++ {
				row := generateTestApiProduct(i)
				res := insertAPIproducts([]common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			// companies
			for i := 100; i < 110; i++ {
				row := generateTestCompany(i)
				res := insertCompanies([]common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			// company developers
			for i := 100; i < 110; i++ {
				row := generateTestCompanyDeveloper(i)
				res := insertCompanyDevelopers([]common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			// application
			k = 100
			for i := 100; i < 110; i++ {
				for j = k; j < 100 + k; j++ {
					row := generateTestAppCompany(j, i)
					res := insertApplications([]common.Row{row}, txn)
					Expect(res).Should(BeTrue())
				}
				k = j
			}
			// app credentials
			for i := 100; i < 110; i++ {
				row := generateTestAppCreds(i)
				res := insertCredentials([]common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}
			// api product mapper
			for i := 100; i < 110; i++ {
				row := generateTestApiProductMapper(i)
				res := insertAPIProductMappers([]common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			txn.Commit()
			var count int64
			db.QueryRow("select count(*) from edgex_data_scope").Scan(&count)
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

		It("Positive DB test for Delete operations", func() {
			db := getDB()
			txn, err := db.Begin()
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 10; i++ {
				row := generateTestApiProductMapper(i)
				res := deleteAPIproductMapper(row, txn)
				Expect(res).Should(BeTrue())
			}

			for i := 0; i < 10; i++ {
				row := generateTestAppCreds(i)
				res := delete("APP_CREDENTIAL", []common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}
			for i := 0; i < 100; i++ {
				row := generateTestApp(i, 999) //TODO we use j in above insertions
				res := delete("APP", []common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			for i := 0; i < 10; i++ {
				row := generateTestDeveloper(i)
				res := delete("DEVELOPER", []common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			for i := 0; i < 10; i++ {
				row := generateTestApiProduct(i)
				res := delete("API_PRODUCT", []common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			for i := 100; i < 110; i++ {
				row := generateTestCompanyDeveloper(i)
				res := delete("COMPANY_DEVELOPER", []common.Row{row}, txn)
				Expect(res).Should(BeTrue())
			}

			txn.Commit()
		})

		It("Negative cases for DB Deletes on KMS tables", func() {
			db := getDB()
			txn, err := db.Begin()
			Expect(err).ShouldNot(HaveOccurred())

			row := generateTestApiProductMapper(999)

			res := delete("APP_CREDENTIAL_APIPRODUCT_MAPPER", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = delete("API_PRODUCT", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = delete("APP_CREDENTIAL", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = delete("DEVELOPER", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = delete("APP", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = delete("COMPANY", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = delete("COMPANY_DEVELOPER", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			txn.Rollback()

		})

		It("Negative cases for DB Inserts/updates on KMS tables", func() {

			db := getDB()
			txn, err := db.Begin()
			Expect(err).ShouldNot(HaveOccurred())

			row := generateTestApiProduct(999)
			row["id"] = nil
			res := insert("API_PRODUCT", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = insert( "APP", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = insert( "APP_CREDENTIAL", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = insert("APP_CREDENTIAL_APIPRODUCT_MAPPER", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = insert("COMPANY", []common.Row{row}, txn)
			Expect(res).Should(BeFalse())

			res = insert("COMPANY_DEVELOPER",[]common.Row{row}, txn)
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
