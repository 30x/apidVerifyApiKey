package apidVerifyApiKey

import (
	"encoding/json"
	"fmt"
	"github.com/30x/apid-core"
	"github.com/apigee-labs/transicator/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
)

var _ = Describe("listener", func() {

	Context("KMS create/updates verification via changes for Developer", func() {
		It("Create KMS tables via changes, and Verify via verifyApiKey", func(done Done) {
			server := mockKMSserver()
			var event = common.ChangeList{}
			closed := 0
			/* API Product */
			srvItems := common.Row{
				"id": {
					Value: "ch_api_product_2",
				},
				"api_resources": {
					Value: "{}",
				},
				"environments": {
					Value: "{Env_0, Env_1}",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* DEVELOPER */
			devItems := common.Row{
				"id": {
					Value: "ch_developer_id_2",
				},
				"status": {
					Value: "Active",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* APP */
			appItems := common.Row{
				"id": {
					Value: "ch_application_id_2",
				},
				"developer_id": {
					Value: "ch_developer_id_2",
				},
				"status": {
					Value: "Approved",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
				"parent_id": {
					Value: "ch_developer_id_2",
				},
			}

			/* CRED */
			credItems := common.Row{
				"id": {
					Value: "ch_app_credential_2",
				},
				"app_id": {
					Value: "ch_application_id_2",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"status": {
					Value: "Approved",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* APP_CRED_APIPRD_MAPPER */
			mpItems := common.Row{
				"apiprdt_id": {
					Value: "ch_api_product_2",
				},
				"app_id": {
					Value: "ch_application_id_2",
				},
				"appcred_id": {
					Value: "ch_app_credential_2",
				},
				"status": {
					Value: "Approved",
				},
				"_change_selector": {
					Value: "test_org0",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
			}

			event.Changes = []common.Change{
				{
					Table:     "kms.api_product",
					NewRow:    srvItems,
					Operation: 1,
				},
				{
					Table:     "kms.developer",
					NewRow:    devItems,
					Operation: 1,
				},
				{
					Table:     "kms.app",
					NewRow:    appItems,
					Operation: 1,
				},
				{
					Table:     "kms.app_credential",
					NewRow:    credItems,
					Operation: 1,
				},
				{
					Table:     "kms.app_credential_apiproduct_mapper",
					NewRow:    mpItems,
					Operation: 1,
				},
			}

			h := &test_handler{
				"checkDatabase post Insertion",
				func(e apid.Event) {
					defer GinkgoRecover()

					// ignore the first event, let standard listener process it
					changeSet := e.(*common.ChangeList)
					if len(changeSet.Changes) > 0 || closed == 1 {
						return
					}

					rsp, err := http.PostForm(
						fmt.Sprintf("%s/verifiers/apikey",
							server.URL),
						url.Values{"key": {"ch_app_credential_2"},
							"uriPath":   {"/test"},
							"scopeuuid": {"XYZ"},
							"action":    {"verify"}})

					Expect(err).Should(Succeed())
					defer rsp.Body.Close()
					body, readErr := ioutil.ReadAll(rsp.Body)
					Expect(readErr).Should(Succeed())
					var respj kmsResponseSuccess
					json.Unmarshal(body, &respj)
					Expect(respj.Type).Should(Equal("APIKeyContext"))
					Expect(rsp.StatusCode).To(Equal(http.StatusOK))
					Expect(respj.RspInfo.Key).Should(Equal("ch_app_credential_2"))
					Expect(respj.RspInfo.Type).Should(Equal("developer"))
					dataValue := rsp.Header.Get("Content-Type")
					Expect(dataValue).To(Equal("application/json"))

					rsp, err = http.PostForm(
						fmt.Sprintf("%s/verifiers/apikey",
							server.URL),
						url.Values{"key": {"ch_app_credential_2"},
							"uriPath":   {"/test"},
							"scopeuuid": {"badscope"},
							"action":    {"verify"}})

					Expect(err).Should(Succeed())
					defer rsp.Body.Close()
					body, readErr = ioutil.ReadAll(rsp.Body)
					Expect(readErr).Should(Succeed())
					var respe kmsResponseFail
					json.Unmarshal(body, &respe)
					Expect(rsp.StatusCode).To(Equal(http.StatusOK))
					dataValue = rsp.Header.Get("Content-Type")
					Expect(dataValue).To(Equal("application/json"))
					Expect(respe.Type).Should(Equal("ErrorResult"))
					Expect(respe.ErrInfo.ErrorCode).Should(Equal("ENV_VALIDATION_FAILED"))

					closed = 1
					close(done)
				},
			}

			apid.Events().Listen("ApigeeSync", h)
			apid.Events().Emit("ApigeeSync", &event)
			apid.Events().Emit("ApigeeSync", &common.ChangeList{})
		})
	})

	Context("KMS create/updates verification via changes for Company", func() {
		It("Create KMS tables via changes, and Verify via verifyApiKey", func(done Done) {
			server := mockKMSserver()
			var event = common.ChangeList{}
			closed := 0
			/* API Product */
			srvItems := common.Row{
				"id": {
					Value: "ch_api_product_5",
				},
				"api_resources": {
					Value: "{}",
				},
				"environments": {
					Value: "{Env_0, Env_1}",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* COMPANY */
			companyItems := common.Row{
				"id": {
					Value: "ch_company_id_5",
				},
				"status": {
					Value: "Active",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
				"name": {
					Value: "test_company_name0",
				},
				"display_name": {
					Value: "test_company_display_name0",
				},
			}
			/* COMPANY_DEVELOPER */
			companyDeveloperItems := common.Row{
				"id": {
					Value: "ch_developer_id_5",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
				"company_id": {
					Value: "ch_company_id_5",
				},
				"developer_id": {
					Value: "ch_developer_id_5",
				},
			}

			/* APP */
			appItems := common.Row{
				"id": {
					Value: "ch_application_id_5",
				},
				"company_id": {
					Value: "ch_company_id_5",
				},
				"status": {
					Value: "Approved",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
				"parent_id": {
					Value: "ch_company_id_5",
				},
			}

			/* CRED */
			credItems := common.Row{
				"id": {
					Value: "ch_app_credential_5",
				},
				"app_id": {
					Value: "ch_application_id_5",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"status": {
					Value: "Approved",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* APP_CRED_APIPRD_MAPPER */
			mpItems := common.Row{
				"apiprdt_id": {
					Value: "ch_api_product_5",
				},
				"app_id": {
					Value: "ch_application_id_5",
				},
				"appcred_id": {
					Value: "ch_app_credential_5",
				},
				"status": {
					Value: "Approved",
				},
				"_change_selector": {
					Value: "test_org0",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
			}

			event.Changes = []common.Change{
				{
					Table:     "kms.api_product",
					NewRow:    srvItems,
					Operation: 1,
				},

				{
					Table:     "kms.app",
					NewRow:    appItems,
					Operation: 1,
				},
				{
					Table:     "kms.app_credential",
					NewRow:    credItems,
					Operation: 1,
				},
				{
					Table:     "kms.app_credential_apiproduct_mapper",
					NewRow:    mpItems,
					Operation: 1,
				},
				{
					Table:     "kms.company",
					NewRow:    companyItems,
					Operation: 1,
				},
				{
					Table:     "kms.company_developer",
					NewRow:    companyDeveloperItems,
					Operation: 1,
				},
			}

			h := &test_handler{
				"checkDatabase post Insertion",
				func(e apid.Event) {
					defer GinkgoRecover()

					// ignore the first event, let standard listener process it
					changeSet := e.(*common.ChangeList)
					if len(changeSet.Changes) > 0 || closed == 1 {
						return
					}

					rsp, err := http.PostForm(
						fmt.Sprintf("%s/verifiers/apikey",
							server.URL),
						url.Values{"key": {"ch_app_credential_5"},
							"uriPath":   {"/test"},
							"scopeuuid": {"XYZ"},
							"action":    {"verify"}})

					Expect(err).Should(Succeed())
					defer rsp.Body.Close()
					body, readErr := ioutil.ReadAll(rsp.Body)
					Expect(readErr).Should(Succeed())
					var respj kmsResponseSuccess
					json.Unmarshal(body, &respj)
					Expect(rsp.StatusCode).To(Equal(http.StatusOK))
					Expect(respj.RspInfo.Type).Should(Equal("company"))
					Expect(respj.Type).Should(Equal("APIKeyContext"))
					Expect(respj.RspInfo.Key).Should(Equal("ch_app_credential_5"))
					dataValue := rsp.Header.Get("Content-Type")
					Expect(dataValue).To(Equal("application/json"))

					rsp, err = http.PostForm(
						fmt.Sprintf("%s/verifiers/apikey",
							server.URL),
						url.Values{"key": {"ch_app_credential_5"},
							"uriPath":   {"/test"},
							"scopeuuid": {"badscope"},
							"action":    {"verify"}})

					Expect(err).Should(Succeed())
					defer rsp.Body.Close()
					body, readErr = ioutil.ReadAll(rsp.Body)
					Expect(readErr).Should(Succeed())
					var respe kmsResponseFail
					json.Unmarshal(body, &respe)
					Expect(rsp.StatusCode).To(Equal(http.StatusOK))
					dataValue = rsp.Header.Get("Content-Type")
					Expect(dataValue).To(Equal("application/json"))
					Expect(respe.Type).Should(Equal("ErrorResult"))
					Expect(respe.ErrInfo.ErrorCode).Should(Equal("ENV_VALIDATION_FAILED"))

					closed = 1
					close(done)
				},
			}

			apid.Events().Listen("ApigeeSync", h)
			apid.Events().Emit("ApigeeSync", &event)
			apid.Events().Emit("ApigeeSync", &common.ChangeList{})
		})
	})

	It("Modify tables in KMS tables, and verify via verifyApiKey for Developer", func(done Done) {
		closed := 0
		var event = common.ChangeList{}
		var event2 = common.ChangeList{}

		/* Orig data */
		/* API Product */
		srvItemsOld := common.Row{
			"id": {
				Value: "ch_api_product_0",
			},
			"api_resources": {
				Value: "{}",
			},
			"environments": {
				Value: "{Env_0, Env_1}",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* DEVELOPER */
		devItemsOld := common.Row{
			"id": {
				Value: "ch_developer_id_0",
			},
			"status": {
				Value: "Active",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP */
		appItemsOld := common.Row{
			"id": {
				Value: "ch_application_id_0",
			},
			"developer_id": {
				Value: "ch_developer_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"parent_id": {
				Value: "ch_developer_id_0",
			},
		}

		/* CRED */
		credItemsOld := common.Row{
			"id": {
				Value: "ch_app_credential_0",
			},
			"app_id": {
				Value: "ch_application_id_0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP_CRED_APIPRD_MAPPER */
		mpItemsOld := common.Row{
			"apiprdt_id": {
				Value: "ch_api_product_0",
			},
			"app_id": {
				Value: "ch_application_id_0",
			},
			"appcred_id": {
				Value: "ch_app_credential_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
		}

		/* New to be replaced data */
		/* API PRODUCT */
		srvItemsNew := common.Row{
			"id": {
				Value: "ch_api_product_1",
			},
			"api_resources": {
				Value: "{}",
			},
			"environments": {
				Value: "{Env_0, Env_1}",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* DEVELOPER */
		devItemsNew := common.Row{
			"id": {
				Value: "ch_developer_id_1",
			},
			"status": {
				Value: "Active",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP */
		appItemsNew := common.Row{
			"id": {
				Value: "ch_application_id_1",
			},
			"developer_id": {
				Value: "ch_developer_id_1",
			},
			"status": {
				Value: "Approved",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"parent_id": {
				Value: "ch_developer_id_1",
			},
		}

		/* CRED */
		credItemsNew := common.Row{
			"id": {
				Value: "ch_app_credential_1",
			},
			"app_id": {
				Value: "ch_application_id_1",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP_CRED_APIPRD_MAPPER */
		mpItemsNew := common.Row{
			"apiprdt_id": {
				Value: "ch_api_product_1",
			},
			"app_id": {
				Value: "ch_application_id_1",
			},
			"appcred_id": {
				Value: "ch_app_credential_1",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
		}

		event.Changes = []common.Change{
			{
				Table:     "kms.api_product",
				NewRow:    srvItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.developer",
				NewRow:    devItemsOld,
				Operation: 1,
			},

			{
				Table:     "kms.app",
				NewRow:    appItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.app_credential",
				NewRow:    credItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.app_credential_apiproduct_mapper",
				NewRow:    mpItemsOld,
				Operation: 1,
			},
		}

		event2.Changes = []common.Change{
			{
				Table:     "kms.api_product",
				OldRow:    srvItemsOld,
				NewRow:    srvItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.developer",
				OldRow:    devItemsOld,
				NewRow:    devItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app",
				OldRow:    appItemsOld,
				NewRow:    appItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app_credential",
				OldRow:    credItemsOld,
				NewRow:    credItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app_credential_apiproduct_mapper",
				OldRow:    mpItemsOld,
				NewRow:    mpItemsNew,
				Operation: 2,
			},
		}

		h := &test_handler{
			"checkDatabase post Insertion",
			func(e apid.Event) {
				defer GinkgoRecover()

				// ignore the first event, let standard listener process it
				changeSet := e.(*common.ChangeList)
				if len(changeSet.Changes) > 0 || closed == 1 {
					return
				}
				v := url.Values{
					"key":       []string{"ch_app_credential_1"},
					"uriPath":   []string{"/test"},
					"scopeuuid": []string{"XYZ"},
					"action":    []string{"verify"},
				}
				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())
				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Key).Should(Equal("ch_app_credential_1"))
				Expect(respj.RspInfo.Type).Should(Equal("developer"))
				closed = 1
				close(done)
			},
		}

		apid.Events().Listen("ApigeeSync", h)
		apid.Events().Emit("ApigeeSync", &event)
		apid.Events().Emit("ApigeeSync", &event2)
		apid.Events().Emit("ApigeeSync", &common.ChangeList{})
	})

	It("Modify tables in KMS tables, and verify via verifyApiKey for Company", func(done Done) {
		closed := 0
		var event = common.ChangeList{}
		var event2 = common.ChangeList{}

		/* Orig data */
		/* API Product */
		srvItemsOld := common.Row{
			"id": {
				Value: "ch_api_product_0",
			},
			"api_resources": {
				Value: "{}",
			},
			"environments": {
				Value: "{Env_0, Env_1}",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* DEVELOPER */
		devItemsOld := common.Row{
			"id": {
				Value: "ch_developer_id_0",
			},
			"status": {
				Value: "Active",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP */
		appItemsOld := common.Row{
			"id": {
				Value: "ch_application_id_0",
			},
			"company_id": {
				Value: "ch_company_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"parent_id": {
				Value: "ch_company_id_0",
			},
		}

		/* COMPANY */
		companyItemsOld := common.Row{
			"id": {
				Value: "ch_company_id_0",
			},
			"status": {
				Value: "Active",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"name": {
				Value: "test_company_name0",
			},
			"display_name": {
				Value: "test_company_display_name0",
			},
		}

		/* COMPANY_DEVELOPER */
		companyDeveloperItemsOld := common.Row{
			"id": {
				Value: "ch_developer_id_0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"company_id": {
				Value: "ch_company_id_0",
			},
			"developer_id": {
				Value: "ch_developer_id_0",
			},
		}

		/* CRED */
		credItemsOld := common.Row{
			"id": {
				Value: "ch_app_credential_0",
			},
			"app_id": {
				Value: "ch_application_id_0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP_CRED_APIPRD_MAPPER */
		mpItemsOld := common.Row{
			"apiprdt_id": {
				Value: "ch_api_product_0",
			},
			"app_id": {
				Value: "ch_application_id_0",
			},
			"appcred_id": {
				Value: "ch_app_credential_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
		}

		/* New to be replaced data */
		/* API PRODUCT */
		srvItemsNew := common.Row{
			"id": {
				Value: "ch_api_product_1",
			},
			"api_resources": {
				Value: "{}",
			},
			"environments": {
				Value: "{Env_0, Env_1}",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* DEVELOPER */
		devItemsNew := common.Row{
			"id": {
				Value: "ch_developer_id_1",
			},
			"status": {
				Value: "Active",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* COMPANY */
		companyItemsNew := common.Row{
			"id": {
				Value: "ch_company_id_1",
			},
			"status": {
				Value: "Active",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"name": {
				Value: "test_company_name0",
			},
			"display_name": {
				Value: "test_company_display_name0",
			},
		}

		/* COMPANY_DEVELOPER */
		companyDeveloperItemsNew := common.Row{
			"id": {
				Value: "ch_developer_id_1",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"company_id": {
				Value: "ch_company_id_1",
			},
			"developer_id": {
				Value: "ch_developer_id_1",
			},
		}

		/* APP */
		appItemsNew := common.Row{
			"id": {
				Value: "ch_application_id_1",
			},
			"company_id": {
				Value: "ch_company_id_1",
			},
			"status": {
				Value: "Approved",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"parent_id": {
				Value: "ch_company_id_1",
			},
		}

		/* CRED */
		credItemsNew := common.Row{
			"id": {
				Value: "ch_app_credential_1",
			},
			"app_id": {
				Value: "ch_application_id_1",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP_CRED_APIPRD_MAPPER */
		mpItemsNew := common.Row{
			"apiprdt_id": {
				Value: "ch_api_product_1",
			},
			"app_id": {
				Value: "ch_application_id_1",
			},
			"appcred_id": {
				Value: "ch_app_credential_1",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
		}

		event.Changes = []common.Change{
			{
				Table:     "kms.api_product",
				NewRow:    srvItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.developer",
				NewRow:    devItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.company",
				NewRow:    companyItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.company_developer",
				NewRow:    companyDeveloperItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.app",
				NewRow:    appItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.app_credential",
				NewRow:    credItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.app_credential_apiproduct_mapper",
				NewRow:    mpItemsOld,
				Operation: 1,
			},
		}

		event2.Changes = []common.Change{
			{
				Table:     "kms.api_product",
				OldRow:    srvItemsOld,
				NewRow:    srvItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.developer",
				OldRow:    devItemsOld,
				NewRow:    devItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.company",
				OldRow:    companyItemsOld,
				NewRow:    companyItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.company_developer",
				OldRow:    companyDeveloperItemsOld,
				NewRow:    companyDeveloperItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app",
				OldRow:    appItemsOld,
				NewRow:    appItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app_credential",
				OldRow:    credItemsOld,
				NewRow:    credItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app_credential_apiproduct_mapper",
				OldRow:    mpItemsOld,
				NewRow:    mpItemsNew,
				Operation: 2,
			},
		}

		h := &test_handler{
			"checkDatabase post Insertion",
			func(e apid.Event) {
				defer GinkgoRecover()

				// ignore the first event, let standard listener process it
				changeSet := e.(*common.ChangeList)
				if len(changeSet.Changes) > 0 || closed == 1 {
					return
				}
				v := url.Values{
					"key":       []string{"ch_app_credential_1"},
					"uriPath":   []string{"/test"},
					"scopeuuid": []string{"XYZ"},
					"action":    []string{"verify"},
				}
				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())
				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Key).Should(Equal("ch_app_credential_1"))

				closed = 1
				close(done)
			},
		}

		apid.Events().Listen("ApigeeSync", h)
		apid.Events().Emit("ApigeeSync", &event)
		apid.Events().Emit("ApigeeSync", &event2)
		apid.Events().Emit("ApigeeSync", &common.ChangeList{})
	})

})

type test_handler struct {
	description string
	f           func(event apid.Event)
}

func (t *test_handler) String() string {
	return t.description
}

func (t *test_handler) Handle(event apid.Event) {
	t.f(event)
}

func addScopes(db apid.DB) {
	txn, _ := db.Begin()
	txn.Exec("INSERT INTO DATA_SCOPE (id, _change_selector, apid_cluster_id, scope, org, env) "+
		"VALUES"+
		"($1,$2,$3,$4,$5,$6)",
		"ABCDE",
		"some_cluster_id",
		"some_cluster_id",
		"tenant_id_xxxx",
		"test_org0",
		"Env_0",
	)
	txn.Exec("INSERT INTO DATA_SCOPE (id, _change_selector, apid_cluster_id, scope, org, env) "+
		"VALUES"+
		"($1,$2,$3,$4,$5,$6)",
		"XYZ",
		"test_org0",
		"somecluster_id",
		"tenant_id_0",
		"test_org0",
		"Env_0",
	)
	log.Info("Inserted DATA_SCOPE for test")
	txn.Commit()
}

func mockKMSserver() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r)
	}))
}
