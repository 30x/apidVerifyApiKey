package apidVerifyApiKey

import (
	"encoding/json"
	"github.com/30x/apid"
	. "github.com/30x/apidApigeeSync" // for direct access to Payload types
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("listener", func() {

	It("should store data from ApigeeSync in the database", func(done Done) {

		var event = common.ChangeList{}
		rowitemp := common.Row{}
		scv := &common.ColumnVal{
			Value: "api_product_0",
			Type:  1,
		}
		rowitem["id"] = scv
		event.Changes = []Change{
			{
				Table:  "api_product",
				NewRow: rowitemp,
			},
		}

		h := &test_handler{
			"checkDatabase",
			func(e apid.Event) {

				// ignore the first event, let standard listener process it
				changeSet := e.(*ChangeSet)
				if len(changeSet.Changes) > 0 {
					return
				}

				rsp, err := verifyAPIKey("credential_sync", "/test", "Env_0", "test_org", "verify")
				Expect(err).ShouldNot(HaveOccurred())

				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Key).Should(Equal("credential_sync"))

				close(done)
			},
		}

		apid.Events().Listen(ApigeeSyncEventSelector, h)
		apid.Events().Emit(ApigeeSyncEventSelector, &event)       // for standard listener
		apid.Events().Emit(ApigeeSyncEventSelector, &ChangeSet{}) // for test listener
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
