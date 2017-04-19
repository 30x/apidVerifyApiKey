package apidVerifyApiKey

import (
	"github.com/30x/apid-core"
	"github.com/apigee-labs/transicator/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("listener", func() {

	Context("KMS create/updates verification via changes for Developer", func() {

		handler := handler{}

		It("should set DB to appropriate version", func() {

			saveDb := getDB()

			s := &common.Snapshot{
				SnapshotInfo: "test_snapshot",
				Tables:       []common.Table{},
			}

			handler.Handle(s)

			expectedDB, err := data.DBVersion(s.SnapshotInfo)
			Expect(err).NotTo(HaveOccurred())

			Expect(getDB() == expectedDB).Should(BeTrue())

			//restore the db to the valid one
			setDB(saveDb)
		})

	})
})

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
