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
		case "kms.developer", "public.developer":
			switch payload.Operation {
			case 1:
				insertCreateDeveloper(payload.NewRow, db)
			}

		case "kms.app", "public.app":
			switch payload.Operation {
			case 1:
				insertCreateApplication(payload.NewRow, db)
			}

		case "kms.app_credential", "public.app_credential":
			switch payload.Operation {
			case 1:
				insertCreateCredential(payload.NewRow, db)
			}
		case "kms.api_product", "public.api_product":
			switch payload.Operation {
			case 1:
				insertAPIproduct(payload.NewRow, db)
			}

		case "kms.app_credential_apiproduct_mapper", "public.app_credential_apiproduct_mapper":
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
func insertCreateCredential(ele common.Row, db *sql.DB) bool {

	var scope, id, appId, consumerSecret, appstatus, status, tenantId string
	var issuedAt int64

	ele.Get("_apid_scope", &scope)
	ele.Get("id", &id)
	ele.Get("app_id", &appId)
	ele.Get("consumer_secret", &consumerSecret)
	ele.Get("app_status", &appstatus)
	ele.Get("status", &status)
	ele.Get("issued_at", &issuedAt)
	ele.Get("tenant_id", &tenantId)

	stmt, err := db.Prepare("INSERT INTO APP_CREDENTIAL (_apid_scope, id, app_id, consumer_secret, app_status, status, issued_at, tenant_id)VALUES($1,$2,$3,$4,$5,$6,$7,$8);")
	if err != nil {
		log.Error("INSERT CRED Failed: ", id, ", ", scope, ")", err)
		return false
	}
	_, err = stmt.Exec(
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
		return false
	} else {
		log.Info("INSERT CRED Success: (", id, ", ", scope, ")")
		return true
	}

}
func insertApiProductMapper(ele common.Row, db *sql.DB) bool {

	var ApiProduct, AppId, EntityIdentifier, tenantId, Scope, Status string

	ele.Get("apiprdt_id", &ApiProduct)
	ele.Get("app_id", &AppId)
	ele.Get("appcred_id", &EntityIdentifier)
	ele.Get("tenant_id", &tenantId)
	ele.Get("_apid_scope", &Scope)
	ele.Get("status", &Status)

	/*
	 * If the credentials has been successfully inserted, insert the
	 * mapping entries associated with the credential
	 */

	stmt, err := db.Prepare("INSERT INTO APP_CREDENTIAL_APIPRODUCT_MAPPER(apiprdt_id, app_id, appcred_id, tenant_id, _apid_scope, status) VALUES($1,$2,$3,$4,$5,$6);")
	if err != nil {
		log.Error("INSERT APP_CREDENTIAL_APIPRODUCT_MAPPER Failed: ", err)
		return false
	}
	_, err = stmt.Exec(
		ApiProduct,
		AppId,
		EntityIdentifier,
		tenantId,
		Scope,
		Status)

	if err != nil {
		log.Error("INSERT APP_CREDENTIAL_APIPRODUCT_MAPPER Failed: (",
			ApiProduct, ", ",
			AppId, ", ",
			EntityIdentifier, ", ",
			tenantId, ", ",
			Scope, ", ",
			Status,
			")",
			err)
		return false
	} else {
		log.Info("INSERT APP_CREDENTIAL_APIPRODUCT_MAPPER Success: (",
			ApiProduct, ", ",
			AppId, ", ",
			EntityIdentifier, ", ",
			tenantId, ", ",
			Scope, ", ",
			Status,
			")")
		return true
	}
}

/*
 * DELETE CRED
 */
func deleteCredential(ele common.Row, db *sql.DB) bool {
	return true
}

/*
 * INSERT INTO API product op
 */
func insertAPIproduct(ele common.Row, db *sql.DB) bool {

	var scope, apiProduct, res, env, tenantId string

	ele.Get("_apid_scope", &scope)
	ele.Get("id", &apiProduct)
	ele.Get("api_resources", &res)
	ele.Get("environments", &env)
	ele.Get("tenant_id", &tenantId)

	stmt, err := db.Prepare("INSERT INTO API_PRODUCT (id, api_resources, environments, tenant_id,_apid_scope) VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		log.Error("INSERT API_PRODUCT Failed: ", err)
		return false
	}
	_, err = stmt.Exec(
		apiProduct,
		res,
		env,
		tenantId,
		scope)

	if err != nil {
		log.Error("INSERT API_PRODUCT Failed: (", apiProduct, ", ", tenantId, ")", err)
		return false
	} else {
		log.Info("INSERT API_PRODUCT Success: (", apiProduct, ", ", tenantId, ")")
		return true
	}

}

/*
 * INSERT INTO APP op
 */
func insertCreateApplication(ele common.Row, db *sql.DB) bool {

	var scope, EntityIdentifier, DeveloperId, CallbackUrl, Status, AppName, AppFamily, tenantId, CreatedBy, LastModifiedBy string
	var CreatedAt, LastModifiedAt int64

	ele.Get("_apid_scope", &scope)
	ele.Get("id", &EntityIdentifier)
	ele.Get("developer_id", &DeveloperId)
	ele.Get("callback_url", &CallbackUrl)
	ele.Get("status", &Status)
	ele.Get("name", &AppName)
	ele.Get("app_family", &AppFamily)
	ele.Get("created_at", &CreatedAt)
	ele.Get("created_by", &CreatedBy)
	ele.Get("updated_at", &LastModifiedAt)
	ele.Get("updated_by", &LastModifiedBy)
	ele.Get("tenant_id", &tenantId)

	stmt, err := db.Prepare("INSERT INTO APP (_apid_scope, id, developer_id,callback_url,status, name, app_family, created_at, created_by,updated_at, updated_by,tenant_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);")
	if err != nil {
		log.Error("INSERT APP Failed: ", err)
		return false
	}
	_, err = stmt.Exec(
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
		log.Error("INSERT APP Failed: (", EntityIdentifier, ", ", tenantId, ")", err)
		return false
	} else {
		log.Info("INSERT APP Success: (", EntityIdentifier, ", ", tenantId, ")")
		return true
	}

}

/*
 * INSERT INTO DEVELOPER op
 */
func insertCreateDeveloper(ele common.Row, db *sql.DB) bool {
	var scope, EntityIdentifier, Email, Status, UserName, FirstName, LastName, tenantId, CreatedBy, LastModifiedBy, Username string
	var CreatedAt, LastModifiedAt int64

	ele.Get("_apid_scope", &scope)
	ele.Get("email", &Email)
	ele.Get("id", &EntityIdentifier)
	ele.Get("tenant_id", &tenantId)
	ele.Get("status", &Status)
	ele.Get("username", &Username)
	ele.Get("first_name", &FirstName)
	ele.Get("last_name", &LastName)
	ele.Get("created_at", &CreatedAt)
	ele.Get("created_by", &CreatedBy)
	ele.Get("updated_at", &LastModifiedAt)
	ele.Get("updated_by", &LastModifiedBy)

	stmt, err := db.Prepare("INSERT INTO DEVELOPER (_apid_scope,email,id,tenant_id,status,username,first_name,last_name,created_at,created_by,updated_at,updated_by) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);")
	if err != nil {
		log.Error("INSERT DEVELOPER Failed: ", err)
		return false
	}

	_, err = stmt.Exec(
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
		log.Error("INSERT DEVELOPER Failed: (", EntityIdentifier, ", ", scope, ")", err)
		return false
	} else {
		log.Info("INSERT DEVELOPER Success: (", EntityIdentifier, ", ", scope, ")")
		return true
	}
}
