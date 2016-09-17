package apidVerifyApiKey

import (
	"database/sql"
	"encoding/json"
	"github.com/30x/apid"
	"github.com/30x/apidApigeeSync"
)

type handler struct {
}

func (h *handler) String() string {
	return "verifyAPIKey"
}

// todo: The following was basically just copied from old APID - needs review.

func (h *handler) Handle(e apid.Event) {
	changeSet, ok := e.(*apidApigeeSync.ChangeSet)
	if !ok {
		log.Errorf("Received non-ChangeSet event. This shouldn't happen!")
		return
	}

	log.Debugf("apigeeSyncEvent: %d changes", len(changeSet.Changes))

	db, err := data.DB()
	if err != nil {
		panic("help me!") // todo: handle
	}

	for _, payload := range changeSet.Changes {

		org := payload.Data.PldCont.Organization

		switch payload.Data.EntityType {
		case "developer":
			switch payload.Data.Operation {
			case "create":
				insertCreateDeveloper(payload.Data, db, org)
			}

		case "app":
			switch payload.Data.Operation {
			case "create":
				insertCreateApplication(payload.Data, db, org)
			}

		case "credential":
			switch payload.Data.Operation {
			case "create":
				insertCreateCredential(payload.Data, db, org)

			case "delete":
				deleteCredential(payload.Data, db, org)
			}

		case "apiproduct":
			switch payload.Data.Operation {
			case "create":
				insertAPIproduct(payload.Data, db, org)
			}
		}

	}
}

/*
 * INSERT INTO APP_CREDENTIAL op
 */
func insertCreateCredential(ele apidApigeeSync.DataPayload, db *sql.DB, org string) bool {

	txn, _ := db.Begin()
	isPass := true
	_, err := txn.Exec("INSERT INTO APP_CREDENTIAL (org, id, app_id, cons_secret, status, issued_at)VALUES(?,?,?,?,?,?);",
		org,
		ele.EntityIdentifier,
		ele.PldCont.AppId,
		ele.PldCont.ConsumerSecret,
		ele.PldCont.Status,
		ele.PldCont.IssuedAt)

	if err != nil {
		isPass = false
		log.Error("INSERT CRED Failed: ", ele.EntityIdentifier, org, ")", err)
		goto OT
	} else {
		log.Info("INSERT CRED Success: (", ele.EntityIdentifier, org, ")")
	}

	/*
	 * If the credentials has been successfully inserted, insert the
	 * mapping entries associated with the credential
	 */

	for _, elem := range ele.PldCont.ApiProducts {

		_, err = txn.Exec("INSERT INTO APP_AND_API_PRODUCT_MAPPER (org, api_prdt_id, app_id, app_cred_id, api_prdt_status) VALUES(?,?,?,?,?);",
			org,
			elem.ApiProduct,
			ele.PldCont.AppId,
			ele.EntityIdentifier,
			elem.Status)

		if err != nil {
			isPass = false
			log.Error("INSERT APP_AND_API_PRODUCT_MAPPER Failed: (",
				org,
				elem.ApiProduct,
				ele.PldCont.AppId,
				ele.EntityIdentifier,
				")",
				err)
			break
		} else {
			log.Info("INSERT APP_AND_API_PRODUCT_MAPPER Success: (",
				org,
				elem.ApiProduct,
				ele.PldCont.AppId,
				ele.EntityIdentifier,
				")")
		}
	}
OT:
	if isPass == true {
		txn.Commit()
	} else {
		txn.Rollback()
	}
	return isPass

}

/*
 * DELETE CRED
 */
func deleteCredential(ele apidApigeeSync.DataPayload, db *sql.DB, org string) bool {

	txn, _ := db.Begin()

	_, err := txn.Exec("DELETE FROM APP_CREDENTIAL WHERE org=? AND id=?;", org, ele.EntityIdentifier)

	if err != nil {
		log.Error("DELETE CRED Failed: (", ele.EntityIdentifier, org, ")", err)
		txn.Rollback()
		return false
	} else {
		log.Info("DELETE CRED Success: (", ele.EntityIdentifier, org, ")")
		txn.Commit()
		return true
	}

}

/*
 * Helper function to convert string slice in to JSON format
 */
func convertSlicetoStringFormat(inpslice []string) string {

	bytes, _ := json.Marshal(inpslice)
	return string(bytes)
}

/*
 * INSERT INTO API product op
 */
func insertAPIproduct(ele apidApigeeSync.DataPayload, db *sql.DB, org string) bool {

	txn, _ := db.Begin()
	restr := convertSlicetoStringFormat(ele.PldCont.Resources)
	envstr := convertSlicetoStringFormat(ele.PldCont.Environments)

	_, err := txn.Exec("INSERT INTO API_PRODUCT (org, id, res_names, env) VALUES(?,?,?,?)",
		org,
		ele.PldCont.AppName,
		restr,
		envstr)

	if err != nil {
		log.Error("INSERT API_PRODUCT Failed: (", ele.PldCont.AppName, org, ")", err)
		txn.Rollback()
		return false
	} else {
		log.Info("INSERT API_PRODUCT Success: (", ele.PldCont.AppName, org, ")")
		txn.Commit()
		return true
	}

}

/*
 * INSERT INTO APP op
 */
func insertCreateApplication(ele apidApigeeSync.DataPayload, db *sql.DB, org string) bool {

	txn, _ := db.Begin()

	_, err := txn.Exec("INSERT INTO APP (org, id, dev_id,cback_url,status, name, app_family, created_at, created_by,updated_at, updated_by) VALUES(?,?,?,?,?,?,?,?,?,?,?);",
		org,
		ele.EntityIdentifier,
		ele.PldCont.DeveloperId,
		ele.PldCont.CallbackUrl,
		ele.PldCont.Status,
		ele.PldCont.AppName,
		ele.PldCont.AppFamily,
		ele.PldCont.CreatedAt,
		ele.PldCont.CreatedBy,
		ele.PldCont.LastModifiedAt,
		ele.PldCont.LastModifiedBy)

	if err != nil {
		log.Error("INSERT APP Failed: (", ele.EntityIdentifier, org, ")", err)
		txn.Rollback()
		return false
	} else {
		log.Info("INSERT APP Success: (", ele.EntityIdentifier, org, ")")
		txn.Commit()
		return true
	}

}

/*
 * INSERT INTO DEVELOPER op
 */
func insertCreateDeveloper(ele apidApigeeSync.DataPayload, db *sql.DB, org string) bool {

	txn, _ := db.Begin()

	_, err := txn.Exec("INSERT INTO DEVELOPER (org, email, id, sts, username, firstname, lastname, created_at,created_by, updated_at, updated_by) VALUES(?,?,?,?,?,?,?,?,?,?,?);",
		org,
		ele.PldCont.Email,
		ele.EntityIdentifier,
		ele.PldCont.Status,
		ele.PldCont.UserName,
		ele.PldCont.FirstName,
		ele.PldCont.LastName,
		ele.PldCont.CreatedAt,
		ele.PldCont.CreatedBy,
		ele.PldCont.LastModifiedAt,
		ele.PldCont.LastModifiedBy)

	if err != nil {
		log.Error("INSERT DEVELOPER Failed: (", ele.PldCont.UserName, org, ")", err)
		txn.Rollback()
		return false
	} else {
		log.Info("INSERT DEVELOPER Success: (", ele.PldCont.UserName, org, ")")
		txn.Commit()
		return true
	}

}
