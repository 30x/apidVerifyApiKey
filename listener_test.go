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

		var event = ChangeSet{}
		event.Changes = []ChangePayload{
			{
				Data: DataPayload{
					EntityType: "apiproduct",
					Operation:  "create",
					PldCont: Payload{
						Organization: "test_org",
						AppName:      "Api_product_sync",
						Resources:    []string{"/**", "/test"},
						Environments: []string{"Env_0", "Env_1"},
					},
				},
			},
			{
				Data: DataPayload{
					EntityType:       "developer",
					Operation:        "create",
					EntityIdentifier: "developer_id_sync",
					PldCont: Payload{
						Organization: "test_org",
						Email:        "person_sync@apigee.com",
						Status:       "Active",
						UserName:     "user_sync",
						FirstName:    "user_first_name_sync",
						LastName:     "user_last_name_sync",
					},
				},
			},
			{
				Data: DataPayload{
					EntityType:       "app",
					Operation:        "create",
					EntityIdentifier: "application_id_sync",
					PldCont: Payload{
						Organization: "test_org",
						Email:        "person_sync@apigee.com",
						Status:       "Approved",
						AppName:      "application_id_sync",
						DeveloperId:  "developer_id_sync",
						CallbackUrl:  "call_back_url",
					},
				},
			},
			{
				Data: DataPayload{
					EntityType:       "credential",
					Operation:        "create",
					EntityIdentifier: "credential_sync",
					PldCont: Payload{
						Organization:   "test_org",
						AppId:          "application_id_sync",
						Status:         "Approved",
						ConsumerSecret: "consumer_secret_sync",
						IssuedAt:       349583485,
						ApiProducts: []Apip{
							{
								ApiProduct: "Api_product_sync",
								Status:     "Approved",
							},
						},
					},
				},
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
