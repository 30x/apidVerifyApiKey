package apidVerifyApiKey

import (
	"database/sql"
	"github.com/30x/apid"
	"github.com/30x/transicator/common"
)

type handler struct {
}

func (h *handler) String() string {
	return "verifyAPIKey"
}

// todo: The following was basically just copied from old APID - needs review.

func (h *handler) Handle(e apid.Event) {

	snapData, ok := e.(*common.Snapshot)
	if ok {
		processSnapshot(snapData)
	} else {
		changeSet, ok := e.(*common.ChangeList)
		if ok {
			processChange(changeSet)
		} else {
			log.Errorf("Received Invalid event. This shouldn't happen!")
		}
	}
	return
}

func processSnapshot(snapshot *common.Snapshot) {

	log.Debugf("Process Snapshot data")

	db, err := data.DB()
	if err != nil {
		panic("Unable to access Sqlite DB")
	}

	for _, payload := range snapshot.Tables {

		switch payload.Name {
		case "developer":
			for _, row := range payload.Rows {
				insertCreateDeveloper(row, db)
			}
		case "app":
			for _, row := range payload.Rows {
				insertCreateApplication(row, db)
			}
		case "app_credential":
			for _, row := range payload.Rows {
				insertCreateCredential(row, db)
			}
		case "api_product":
			for _, row := range payload.Rows {
				insertAPIproduct(row, db)
			}
		case "app_credential_apiproduct_mapper":
			for _, row := range payload.Rows {
				insertApiProductMapper(row, db)
			}

		}
	}
}

func processChange(changes *common.ChangeList) {

	log.Debugf("apigeeSyncEvent: %d changes", len(changes.Changes))

	db, err := data.DB()
	if err != nil {
		panic("Unable to access Sqlite DB")
	}

	for _, payload := range changes.Changes {

		switch payload.Table {
		case "public.developer":
			switch payload.Operation {
			case 1:
				insertCreateDeveloper(payload.NewRow, db)
			}

		case "public.app":
			switch payload.Operation {
			case 1:
				insertCreateApplication(payload.NewRow, db)
			}

		case "public.app_credential":
			switch payload.Operation {
			case 1:
				insertCreateCredential(payload.NewRow, db)
			}
		case "public.api_product":
			switch payload.Operation {
			case 1:
				insertAPIproduct(payload.NewRow, db)
			}

		case "public.app_credential_apiproduct_mapper":
			switch payload.Operation {
			case 1:
				insertApiProductMapper(payload.NewRow, db)
			}

		}
	}
}

/*
 * INSERT INTO APP_CREDENTIAL op
 */
func insertCreateCredential(ele common.Row, db *sql.DB) {

	var scope, id, appId, consumerSecret, appstatus, status, tenantId string
	var issuedAt int64

	txn, _ := db.Begin()
	err := ele.Get("_apid_scope", &scope)
	err = ele.Get("id", &id)
	err = ele.Get("app_id", &appId)
	err = ele.Get("consumer_secret", &consumerSecret)
	err = ele.Get("app_status", &appstatus)
	err = ele.Get("status", &status)
	err = ele.Get("issued_at", &issuedAt)
	err = ele.Get("tenant_id", &tenantId)

	_, err = txn.Exec("INSERT INTO APP_CREDENTIAL (_apid_scope, id, app_id, consumer_secret, app_status, status, issued_at, tenant_id)VALUES(?,?,?,?,?,?,?,?);",
		scope,
		id,
		appId,
		consumerSecret,
		appstatus,
		status,
		issuedAt,
		tenantId)

	if err != nil {
		log.Error("INSERT CRED Failed: ", id, ", ", scope, ")", err)
		txn.Rollback()
	} else {
		log.Info("INSERT CRED Success: (", id, ", ", scope, ")")
		txn.Commit()
	}

}
func insertApiProductMapper(ele common.Row, db *sql.DB) {

	var ApiProduct, AppId, EntityIdentifier, tenantId, Scope, Status string

	txn, _ := db.Begin()
	err := ele.Get("apiprdt_id", &ApiProduct)
	err = ele.Get("app_id", &AppId)
	err = ele.Get("appcred_id", &EntityIdentifier)
	err = ele.Get("tenant_id", &tenantId)
	err = ele.Get("_apid_scope", &Scope)
	err = ele.Get("status", &Status)

	/*
	 * If the credentials has been successfully inserted, insert the
	 * mapping entries associated with the credential
	 */

	_, err = txn.Exec("INSERT INTO APP_CREDENTIAL_APIPRODUCT_MAPPER(apiprdt_id, app_id, appcred_id, tenant_id, _apid_scope, status) VALUES(?,?,?,?,?,?);",
		ApiProduct,
		AppId,
		EntityIdentifier,
		tenantId,
		Scope,
		Status)

	if err != nil {
		log.Error("INSERT APP_CREDENTIAL_APIPRODUCT_MAPPER Failed: (",
			ApiProduct,
			AppId,
			EntityIdentifier,
			tenantId,
			Scope,
			Status,
			")",
			err)
		txn.Rollback()
	} else {
		log.Info("INSERT APP_CREDENTIAL_APIPRODUCT_MAPPER Success: (",
			ApiProduct,
			AppId,
			EntityIdentifier,
			tenantId,
			Scope,
			Status,
			")")
		txn.Commit()
	}
}

/*
 * DELETE CRED
 */
func deleteCredential(ele common.Row, db *sql.DB) {

}

/*
 * INSERT INTO API product op
 */
func insertAPIproduct(ele common.Row, db *sql.DB) {

	var scope, apiProduct, res, env, tenantId string

	txn, _ := db.Begin()
	err := ele.Get("_apid_scope", &scope)
	err = ele.Get("id", &apiProduct)
	err = ele.Get("api_resources", &res)
	err = ele.Get("environments", &env)
	err = ele.Get("tenant_id", &tenantId)

	_, err = txn.Exec("INSERT INTO API_PRODUCT (id, api_resources, environments, tenant_id,_apid_scope) VALUES(?,?,?,?,?)",
		apiProduct,
		res,
		env,
		tenantId,
		scope)

	if err != nil {
		log.Error("INSERT API_PRODUCT Failed: (", apiProduct, tenantId, ")", err)
		txn.Rollback()
	} else {
		log.Info("INSERT API_PRODUCT Success: (", apiProduct, tenantId, ")")
		txn.Commit()
	}

}

/*
 * INSERT INTO APP op
 */
func insertCreateApplication(ele common.Row, db *sql.DB) {

	var scope, EntityIdentifier, DeveloperId, CallbackUrl, Status, AppName, AppFamily, tenantId, CreatedBy, LastModifiedBy string
	var CreatedAt, LastModifiedAt int64
	txn, _ := db.Begin()

	err := ele.Get("_apid_scope", &scope)
	err = ele.Get("id", &EntityIdentifier)
	err = ele.Get("developer_id", &DeveloperId)
	err = ele.Get("callback_url", &CallbackUrl)
	err = ele.Get("status", &Status)
	err = ele.Get("name", &AppName)
	err = ele.Get("app_family", &AppFamily)
	err = ele.Get("created_at", &CreatedAt)
	err = ele.Get("created_by", &CreatedBy)
	err = ele.Get("updated_at", &LastModifiedAt)
	err = ele.Get("updated_by", &LastModifiedBy)
	err = ele.Get("tenant_id", &tenantId)

	_, err = txn.Exec("INSERT INTO APP (_apid_scope, id, developer_id,callback_url,status, name, app_family, created_at, created_by,updated_at, updated_by,tenant_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?);",
		scope,
		EntityIdentifier,
		DeveloperId,
		CallbackUrl,
		Status,
		AppName,
		AppFamily,
		CreatedAt,
		CreatedBy,
		LastModifiedAt,
		LastModifiedBy,
		tenantId)

	if err != nil {
		log.Error("INSERT APP Failed: (", EntityIdentifier, tenantId, ")", err)
		txn.Rollback()
	} else {
		log.Info("INSERT APP Success: (", EntityIdentifier, tenantId, ")")
		txn.Commit()
	}

}

/*
 * INSERT INTO DEVELOPER op
 */
func insertCreateDeveloper(ele common.Row, db *sql.DB) {
	var scope, EntityIdentifier, Email, Status, UserName, FirstName, LastName, tenantId, CreatedBy, LastModifiedBy, Username string
	var CreatedAt, LastModifiedAt int64
	txn, _ := db.Begin()

	err := ele.Get("_apid_scope", &scope)
	err = ele.Get("email", &Email)
	err = ele.Get("id", &EntityIdentifier)
	err = ele.Get("tenant_id", &tenantId)
	err = ele.Get("status", &Status)
	err = ele.Get("username", &Username)
	err = ele.Get("first_name", &FirstName)
	err = ele.Get("last_name", &LastName)
	err = ele.Get("created_at", &CreatedAt)
	err = ele.Get("created_by", &CreatedBy)
	err = ele.Get("updated_at", &LastModifiedAt)
	err = ele.Get("updated_by", &LastModifiedBy)

	_, err = txn.Exec("INSERT INTO DEVELOPER (_apid_scope,email,id,tenant_id,status,username,first_name,last_name,created_at,created_by,updated_at,updated_by) VALUES(?,?,?,?,?,?,?,?,?,?,?,?);",
		scope,
		Email,
		EntityIdentifier,
		tenantId,
		Status,
		UserName,
		FirstName,
		LastName,
		CreatedAt,
		CreatedBy,
		LastModifiedAt,
		LastModifiedBy)

	if err != nil {
		log.Error("INSERT DEVELOPER Failed: (", EntityIdentifier, scope, ")", err)
		txn.Rollback()
	} else {
		log.Info("INSERT DEVELOPER Success: (", EntityIdentifier, scope, ")")
		txn.Commit()
	}
}
