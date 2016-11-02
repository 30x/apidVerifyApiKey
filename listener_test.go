package apidVerifyApiKey

import (
	"encoding/json"
	"github.com/30x/apid"
	. "github.com/30x/apidApigeeSync" // for direct access to Payload types
	"github.com/30x/transicator/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("listener", func() {

	It("should store data from ApigeeSync in the database", func(done Done) {

		var event = common.ChangeList{}
		/* API Product */
		srvItems := common.Row{}
		scv := &common.ColumnVal{
			Value: "ch_api_product_0",
			Type:  1,
		}
		srvItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "{}",
			Type:  1,
		}
		srvItems["api_resources"] = scv

		scv = &common.ColumnVal{
			Value: "{Env_0, Env_1}",
			Type:  1,
		}
		srvItems["environments"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_0",
			Type:  1,
		}
		srvItems["tenant_id"] = scv

		scv = &common.ColumnVal{
			Value: "test_org0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		/* DEVELOPER */
		devItems := common.Row{}
		scv = &common.ColumnVal{
			Value: "ch_developer_id_0",
			Type:  1,
		}
		devItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "Active",
			Type:  1,
		}
		devItems["status"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_0",
			Type:  1,
		}
		devItems["tenant_id"] = scv

		scv = &common.ColumnVal{
			Value: "test_org0",
			Type:  1,
		}
		devItems["_apid_scope"] = scv

		/* APP */
		appItems := common.Row{}
		scv = &common.ColumnVal{
			Value: "ch_application_id_0",
			Type:  1,
		}
		appItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "ch_developer_id_0",
			Type:  1,
		}
		appItems["developer_id"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_0",
			Type:  1,
		}
		appItems["tenant_id"] = scv

		scv = &common.ColumnVal{
			Value: "Approved",
			Type:  1,
		}
		appItems["status"] = scv

		scv = &common.ColumnVal{
			Value: "test_org0",
			Type:  1,
		}
		appItems["_apid_scope"] = scv

		/* CRED */
		credItems := common.Row{}
		scv = &common.ColumnVal{
			Value: "ch_app_credential_0",
			Type:  1,
		}
		credItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "ch_application_id_0",
			Type:  1,
		}
		credItems["app_id"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_0",
			Type:  1,
		}
		credItems["tenant_id"] = scv

		scv = &common.ColumnVal{
			Value: "Approved",
			Type:  1,
		}
		credItems["status"] = scv

		scv = &common.ColumnVal{
			Value: "test_org0",
			Type:  1,
		}
		credItems["_apid_scope"] = scv

		/* APP_CRED_APIPRD_MAPPER */
		mpItems := common.Row{}
		scv = &common.ColumnVal{
			Value: "ch_api_product_0",
			Type:  1,
		}
		mpItems["apiprdt_id"] = scv

		scv = &common.ColumnVal{
			Value: "ch_application_id_0",
			Type:  1,
		}
		mpItems["app_id"] = scv

		scv = &common.ColumnVal{
			Value: "ch_app_credential_0",
			Type:  1,
		}
		mpItems["appcred_id"] = scv

		scv = &common.ColumnVal{
			Value: "Approved",
			Type:  1,
		}
		mpItems["status"] = scv

		scv = &common.ColumnVal{
			Value: "test_org0",
			Type:  1,
		}
		mpItems["_apid_scope"] = scv

		event.Changes = []common.Change{
			{
				Table:     "public.api_product",
				NewRow:    srvItems,
				Operation: 1,
			},
			{
				Table:     "public.developer",
				NewRow:    devItems,
				Operation: 1,
			},
			{
				Table:     "public.app",
				NewRow:    appItems,
				Operation: 1,
			},
			{
				Table:     "public.app_credential",
				NewRow:    credItems,
				Operation: 1,
			},
			{
				Table:     "public.app_credential_apiproduct_mapper",
				NewRow:    mpItems,
				Operation: 1,
			},
		}

		h := &test_handler{
			"checkDatabase",
			func(e apid.Event) {

				// ignore the first event, let standard listener process it
				changeSet := e.(*common.ChangeList)
				if len(changeSet.Changes) > 0 {
					return
				}
				processChange(changeSet)
				rsp, err := verifyAPIKey("ch_app_credential_0", "/test", "Env_0", "test_org0", "verify")
				Expect(err).ShouldNot(HaveOccurred())

				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Key).Should(Equal("ch_app_credential_0"))

				close(done)
			},
		}

		apid.Events().Listen(ApigeeSyncEventSelector, h)
		apid.Events().Emit(ApigeeSyncEventSelector, &event)               // for standard listener
		apid.Events().Emit(ApigeeSyncEventSelector, &common.ChangeList{}) // for test listener
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
