package apidVerifyApiKey

import (
	"encoding/json"
	"github.com/30x/apid"
	"github.com/apigee-labs/transicator/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("listener", func() {

	It("should store data from ApigeeSync in the database", func(done Done) {

		var event = common.ChangeList{}
		var event2 = common.ChangeList{}

		/* API Product */
		srvItems := common.Row{
			"id": {
				Value: "ch_api_product_0",
				Type:  1,
			},
			"apid_resources": {
				Value: "{}",
				Type:  1,
			},
			"environments": {
				Value: "{Env_0, Env_1}",
				Type:  1,
			},
			"tenant_id": {
				Value: "tenant_id_0",
				Type:  1,
			},
			"_apid_scope": {
				Value: "test_org0",
				Type:  1,
			},
		}

		/* DEVELOPER */
		devItems := common.Row{
			"id": {
				Value: "ch_developer_id_0",
				Type:  1,
			},
			"status": {
				Value: "Active",
				Type:  1,
			},
			"tenant_id": {
				Value: "tenant_id_0",
				Type:  1,
			},
			"_apid_scope": {
				Value: "test_org0",
				Type:  1,
			},
		}

		/* APP */
		appItems := common.Row{
			"id": {
				Value: "ch_application_id_0",
				Type:  1,
			},
			"developer_id": {
				Value: "ch_developer_id_0",
				Type:  1,
			},
			"status": {
				Value: "Approved",
				Type:  1,
			},
			"tenant_id": {
				Value: "tenant_id_0",
				Type:  1,
			},
			"_apid_scope": {
				Value: "test_org0",
				Type:  1,
			},
		}

		/* CRED */
		credItems := common.Row{
			"id": {
				Value: "ch_app_credential_0",
				Type:  1,
			},
			"app_id": {
				Value: "ch_application_id_0",
				Type:  1,
			},
			"tenant_id": {
				Value: "tenant_id_0",
				Type:  1,
			},
			"status": {
				Value: "Approved",
				Type:  1,
			},
			"_apid_scope": {
				Value: "test_org0",
				Type:  1,
			},
		}

		/* APP_CRED_APIPRD_MAPPER */
		mpItems := common.Row{
			"apiprdt_id": {
				Value: "ch_api_product_0",
				Type:  1,
			},
			"app_id": {
				Value: "ch_application_id_0",
				Type:  1,
			},
			"appcred_id": {
				Value: "ch_app_credential_0",
				Type:  1,
			},
			"status": {
				Value: "Approved",
				Type:  1,
			},
			"_apid_scope": {
				Value: "test_org0",
				Type:  1,
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

		event2.Changes = []common.Change{
			{
				Table:     "kms.api_product",
				OldRow:    srvItems,
				Operation: 3,
			},
			{
				Table:     "kms.developer",
				OldRow:    devItems,
				Operation: 3,
			},
			{
				Table:     "kms.app",
				OldRow:    appItems,
				Operation: 3,
			},
			{
				Table:     "kms.app_credential",
				OldRow:    credItems,
				Operation: 3,
			},
			{
				Table:     "kms.app_credential_apiproduct_mapper",
				OldRow:    mpItems,
				Operation: 3,
			},
		}
		h := &test_handler{
			"checkDatabase post Insertion",
			func(e apid.Event) {
				defer GinkgoRecover()

				// ignore the first event, let standard listener process it
				changeSet := e.(*common.ChangeList)
				if len(changeSet.Changes) > 0 {
					return
				}
				v := url.Values{
					"key": []string{"ch_app_credential_0"},
					"uriPath": []string{"/test"},
					"environment": []string{"Env_0"},
					"organization": []string{"test_org0"},
					"action": []string{"verify"},
				}
				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())
				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Key).Should(Equal("ch_app_credential_0"))
				close(done)
			},
		}

		apid.Events().Listen("ApigeeSync", h)
		apid.Events().Emit("ApigeeSync", &event)
		apid.Events().Emit("ApigeeSync", &event2)
		apid.Events().Emit("ApigeeSync", &event)
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
