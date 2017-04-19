package apidVerifyApiKey

import (
	"github.com/30x/apid-core"
	"github.com/apigee-labs/transicator/common"
)

type handler struct {
}

func (h *handler) String() string {
	return "verifyAPIKey"
}

func (h *handler) Handle(e apid.Event) {

	snapData, ok := e.(*common.Snapshot)
	if ok {
		processSnapshot(snapData)
	}
	return
}

func processSnapshot(snapshot *common.Snapshot) {

	log.Debugf("Snapshot received. Switching to DB version: %s", snapshot.SnapshotInfo)

	db, err := data.DBVersion(snapshot.SnapshotInfo)
	if err != nil {
		log.Panicf("Unable to access database: %v", err)
	}

	setDB(db)
	return
}