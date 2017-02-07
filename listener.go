package apidVerifyApiKey

import (
	"database/sql"
	"github.com/30x/apid"
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
	} else {
		changeSet, ok := e.(*common.ChangeList)
		if ok {
			processChange(changeSet)
		} else {
			log.Errorf("Received Invalid event. Ignoring. %v", e)
		}
	}
	return
}

func processSnapshot(snapshot *common.Snapshot) {

	log.Debugf("Snapshot received. Switching to DB version: %s", snapshot.SnapshotInfo)

	db, err := data.DBVersion(snapshot.SnapshotInfo)
	if err != nil {
		log.Panicf("Unable to access database: %v", err)
	}

	createTables(db)

	if len(snapshot.Tables) > 0 {
		txn, err := db.Begin()
		if err != nil {
			log.Panicf("Unable to create transaction: %v", err)
			return
		}

		/*
		 * Iterate the tables, and insert the rows,
		 * Commit them in bulk.
		 */
		ok := true
		for _, payload := range snapshot.Tables {
			switch payload.Name {
			case "kms.developer":
				ok = insertDevelopers(payload.Rows, txn)
			case "kms.app":
				ok = insertApplications(payload.Rows, txn)
			case "kms.app_credential":
				ok = insertCredentials(payload.Rows, txn)
			case "kms.api_product":
				ok = insertAPIproducts(payload.Rows, txn)
			case "kms.app_credential_apiproduct_mapper":
				ok = insertAPIProductMappers(payload.Rows, txn)
			}
			if !ok {
				log.Error("Error encountered in Downloading Snapshot for VerifyApiKey")
				txn.Rollback()
				return
			}
		}
		log.Debug("Downloading Snapshot for VerifyApiKey complete")
		txn.Commit()
	}

	setDB(db)
	return
}

/*
 * Performs bulk insert of credentials
 */
func insertCredentials(rows []common.Row, txn *sql.Tx) bool {

	var scope, id, appId, consumerSecret, appstatus, status, tenantId string
	var issuedAt int64

	prep, err := txn.Prepare("INSERT INTO APP_CREDENTIAL (_change_selector, id, app_id, consumer_secret, app_status, status, issued_at, tenant_id)VALUES($1,$2,$3,$4,$5,$6,$7,$8);")
	if err != nil {
		log.Error("INSERT Cred Failed: ", err)
		return false
	}
	defer prep.Close()
	for _, ele := range rows {
		ele.Get("_change_selector", &scope)
		ele.Get("id", &id)
		ele.Get("app_id", &appId)
		ele.Get("consumer_secret", &consumerSecret)
		ele.Get("app_status", &appstatus)
		ele.Get("status", &status)
		ele.Get("issued_at", &issuedAt)
		ele.Get("tenant_id", &tenantId)

		/* Mandatory params check */
		if id == "" || scope == "" || tenantId == "" {
			log.Error("INSERT APP_CREDENTIAL: i/p args missing")
			return false
		}
		_, err = txn.Stmt(prep).Exec(
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
			log.Debug("INSERT CRED Success: (", id, ", ", scope, ")")
		}
	}
	return true
}

/*
 * Performs Bulk insert of Applications
 */
func insertApplications(rows []common.Row, txn *sql.Tx) bool {

	var scope, EntityIdentifier, DeveloperId, CallbackUrl, Status, AppName, AppFamily, tenantId, CreatedBy, LastModifiedBy string
	var CreatedAt, LastModifiedAt int64

	prep, err := txn.Prepare("INSERT INTO APP (_change_selector, id, developer_id,callback_url,status, name, app_family, created_at, created_by,updated_at, updated_by,tenant_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);")
	if err != nil {
		log.Error("INSERT APP Failed: ", err)
		return false
	}

	defer prep.Close()
	for _, ele := range rows {

		ele.Get("_change_selector", &scope)
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

		/* Mandatory params check */
		if EntityIdentifier == "" || scope == "" || tenantId == "" {
			log.Error("INSERT APP: i/p args missing")
			return false
		}
		_, err = txn.Stmt(prep).Exec(
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
			log.Debug("INSERT APP Success: (", EntityIdentifier, ", ", tenantId, ")")
		}
	}
	return true

}

/*
 * Performs bulk insert of Developers
 */
func insertDevelopers(rows []common.Row, txn *sql.Tx) bool {

	var scope, EntityIdentifier, Email, Status, UserName, FirstName, LastName, tenantId, CreatedBy, LastModifiedBy, Username string
	var CreatedAt, LastModifiedAt int64

	prep, err := txn.Prepare("INSERT INTO DEVELOPER (_change_selector,email,id,tenant_id,status,username,first_name,last_name,created_at,created_by,updated_at,updated_by) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);")
	if err != nil {
		log.Error("INSERT DEVELOPER Failed: ", err)
		return false
	}

	defer prep.Close()
	for _, ele := range rows {

		ele.Get("_change_selector", &scope)
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

		/* Mandatory params check */
		if EntityIdentifier == "" || scope == "" || tenantId == "" {
			log.Error("INSERT DEVELOPER: i/p args missing")
			return false
		}
		_, err = txn.Stmt(prep).Exec(
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
			log.Debug("INSERT DEVELOPER Success: (", EntityIdentifier, ", ", scope, ")")
		}
	}
	return true
}

/*
 * Performs Bulk insert of Company Developers
 */
func insertCompanyDevelopers(rows []common.Row, txn *sql.Tx) bool {
	var scope, CompanyId, DeveloperId, tenantId, CreatedBy, LastModifiedBy string
	var CreatedAt, LastModifiedAt int64

	prep, err := txn.Prepare("INSERT INTO COMPANY_DEVELOPER (_change_selector,company_id,tenant_id,developer_id,created_at,created_by,updated_at,updated_by) VALUES($1,$2,$3,$4,$5,$6,$7,$8);")
	if err != nil {
		log.Error("INSERT COMPANY_DEVELOPER Failed: ", err)
		return false
	}
	defer prep.Close()
	for _, ele := range rows {

		ele.Get("_change_selector", &scope)
		ele.Get("company_id", &CompanyId)
		ele.Get("tenant_id", &tenantId)
		ele.Get("developer_id", &DeveloperId)
		ele.Get("created_at", &CreatedAt)
		ele.Get("created_by", &CreatedBy)
		ele.Get("updated_at", &LastModifiedAt)
		ele.Get("updated_by", &LastModifiedBy)

		/* Mandatory params check */
		if scope == "" || tenantId == "" || CompanyId == "" || DeveloperId == ""{
			log.Error("INSERT COMPANY_DEVELOPER: i/p args missing")
			return false
		}
		_, err = txn.Stmt(prep).Exec(
			scope,
			CompanyId,
			tenantId,
			DeveloperId,
			CreatedAt,
			CreatedBy,
			LastModifiedAt,
			LastModifiedBy)

		if err != nil {
			log.Error("INSERT COMPANY_DEVELOPER Failed: (", DeveloperId, ", ", CompanyId, ", ", scope, ")", err)
			return false
		} else {
			log.Debug("INSERT COMPANY_DEVELOPER Success: (", DeveloperId, ", ", CompanyId, ", ", scope, ")")
		}
	}
	return true
}

/*
 * Performs Bulk insert of Companies
 */
func insertCompanies(rows []common.Row, txn *sql.Tx) bool {
	var scope, EntityIdentifier, Name, DisplayName, Status, tenantId, CreatedBy, LastModifiedBy string
	var CreatedAt, LastModifiedAt int64

	prep, err := txn.Prepare("INSERT INTO COMPANY (_change_selector,id,tenant_id,status,name,display_name,created_at,created_by,updated_at,updated_by) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);")
	if err != nil {
		log.Error("INSERT COMPANY Failed: ", err)
		return false
	}
	defer prep.Close()
	for _, ele := range rows {

		ele.Get("_change_selector", &scope)
		ele.Get("id", &EntityIdentifier)
		ele.Get("tenant_id", &tenantId)
		ele.Get("status", &Status)
		ele.Get("name", &Name)
		ele.Get("display_name", &DisplayName)
		ele.Get("created_at", &CreatedAt)
		ele.Get("created_by", &CreatedBy)
		ele.Get("updated_at", &LastModifiedAt)
		ele.Get("updated_by", &LastModifiedBy)

		/* Mandatory params check */
		if EntityIdentifier == "" || scope == "" || tenantId == "" {
			log.Error("INSERT COMPANY: i/p args missing")
			return false
		}
		_, err = txn.Stmt(prep).Exec(
			scope,
			EntityIdentifier,
			tenantId,
			Status,
			Name,
			DisplayName,
			CreatedAt,
			CreatedBy,
			LastModifiedAt,
			LastModifiedBy)

		if err != nil {
			log.Error("INSERT COMPANY Failed: (", EntityIdentifier, ", ", scope, ")", err)
			return false
		} else {
			log.Debug("INSERT COMPANY Success: (", EntityIdentifier, ", ", scope, ")")
		}
	}
	return true
}

/*
 * Performs Bulk insert of API products
 */
func insertAPIproducts(rows []common.Row, txn *sql.Tx) bool {

	var scope, apiProduct, res, env, tenantId string

	prep, err := txn.Prepare("INSERT INTO API_PRODUCT (id, api_resources, environments, tenant_id,_change_selector) VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		log.Error("INSERT API_PRODUCT Failed: ", err)
		return false
	}

	defer prep.Close()
	for _, ele := range rows {

		ele.Get("_change_selector", &scope)
		ele.Get("id", &apiProduct)
		ele.Get("api_resources", &res)
		ele.Get("environments", &env)
		ele.Get("tenant_id", &tenantId)

		/* Mandatory params check */
		if apiProduct == "" || scope == "" || tenantId == "" {
			log.Error("INSERT API_PRODUCT: i/p args missing")
			return false
		}
		_, err = txn.Stmt(prep).Exec(
			apiProduct,
			res,
			env,
			tenantId,
			scope)

		if err != nil {
			log.Error("INSERT API_PRODUCT Failed: (", apiProduct, ", ", tenantId, ")", err)
			return false
		} else {
			log.Debug("INSERT API_PRODUCT Success: (", apiProduct, ", ", tenantId, ")")
		}
	}
	return true
}

/*
 * Performs a bulk insert of all APP_CREDENTIAL_APIPRODUCT_MAPPER rows
 */
func insertAPIProductMappers(rows []common.Row, txn *sql.Tx) bool {

	var ApiProduct, AppId, EntityIdentifier, tenantId, Scope, Status string

	prep, err := txn.Prepare("INSERT INTO APP_CREDENTIAL_APIPRODUCT_MAPPER(apiprdt_id, app_id, appcred_id, tenant_id, _change_selector, status) VALUES($1,$2,$3,$4,$5,$6);")
	if err != nil {
		log.Error("INSERT APP_CREDENTIAL_APIPRODUCT_MAPPER Failed: ", err)
		return false
	}

	defer prep.Close()
	for _, ele := range rows {

		ele.Get("apiprdt_id", &ApiProduct)
		ele.Get("app_id", &AppId)
		ele.Get("appcred_id", &EntityIdentifier)
		ele.Get("tenant_id", &tenantId)
		ele.Get("_change_selector", &Scope)
		ele.Get("status", &Status)

		/* Mandatory params check */
		if ApiProduct == "" || AppId == "" || EntityIdentifier == "" || tenantId == "" || Scope == "" {
			log.Error("INSERT APP_CREDENTIAL_APIPRODUCT_MAPPER : i/p args missing")
			return false
		}

		/*
		 * If the credentials has been successfully inserted, insert the
		 * mapping entries associated with the credential
		 */

		_, err = txn.Stmt(prep).Exec(
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
			log.Debug("INSERT APP_CREDENTIAL_APIPRODUCT_MAPPER Success: (",
				ApiProduct, ", ",
				AppId, ", ",
				EntityIdentifier, ", ",
				tenantId, ", ",
				Scope, ", ",
				Status,
				")")
		}
	}
	return true
}

func processChange(changes *common.ChangeList) {

	db := getDB()

	txn, err := db.Begin()
	if err != nil {
		log.Error("Unable to create transaction")
		return
	}
	defer txn.Rollback()

	var rows []common.Row
	ok := true

	log.Debugf("apigeeSyncEvent: %d changes", len(changes.Changes))
	for _, payload := range changes.Changes {
		rows = nil
		switch payload.Table {
		case "kms.developer":
			switch payload.Operation {
			case common.Insert:
				rows = append(rows, payload.NewRow)
				ok = insertDevelopers(rows, txn)

			case common.Update:
				ok = deleteObject("DEVELOPER", payload.OldRow, txn)
				rows = append(rows, payload.NewRow)
				ok = insertDevelopers(rows, txn)

			case common.Delete:
				ok = deleteObject("DEVELOPER", payload.OldRow, txn)
			}
		case "kms.app":
			switch payload.Operation {
			case common.Insert:
				rows = append(rows, payload.NewRow)
				ok = insertApplications(rows, txn)

			case common.Update:
				ok = deleteObject("APP", payload.OldRow, txn)
				rows = append(rows, payload.NewRow)
				ok = insertApplications(rows, txn)

			case common.Delete:
				ok = deleteObject("APP", payload.OldRow, txn)
			}
		case "kms.company":
			switch payload.Operation {
			case common.Insert:
				rows = append(rows, payload.NewRow)
				ok = insertCompanies(rows, txn)

			case common.Update:
				ok = deleteObject("COMPANY", payload.OldRow, txn)
				rows = append(rows, payload.NewRow)
				ok = insertCompanies(rows, txn)

			case common.Delete:
				ok = deleteObject("COMPANY", payload.OldRow, txn)
			}
		case "kms.company_developer":
			switch payload.Operation {
			case common.Insert:
				rows = append(rows, payload.NewRow)
				ok = insertCompanyDevelopers(rows, txn)

			case common.Update:
				ok = deleteCompanyDeveloper(payload.OldRow, txn)
				rows = append(rows, payload.NewRow)
				ok = insertCompanyDevelopers(rows, txn)

			case common.Delete:
				ok = deleteCompanyDeveloper(payload.OldRow, txn)
			}
		case "kms.app_credential":
			switch payload.Operation {
			case common.Insert:
				rows = append(rows, payload.NewRow)
				ok = insertCredentials(rows, txn)

			case common.Update:
				ok = deleteObject("APP_CREDENTIAL", payload.OldRow, txn)
				rows = append(rows, payload.NewRow)
				ok = insertCredentials(rows, txn)

			case common.Delete:
				ok = deleteObject("APP_CREDENTIAL", payload.OldRow, txn)
			}
		case "kms.api_product":
			switch payload.Operation {
			case common.Insert:
				rows = append(rows, payload.NewRow)
				ok = insertAPIproducts(rows, txn)

			case common.Update:
				ok = deleteObject("API_PRODUCT", payload.OldRow, txn)
				rows = append(rows, payload.NewRow)
				ok = insertAPIproducts(rows, txn)

			case common.Delete:
				ok = deleteObject("API_PRODUCT", payload.OldRow, txn)
			}

		case "kms.app_credential_apiproduct_mapper":
			switch payload.Operation {
			case common.Insert:
				rows = append(rows, payload.NewRow)
				ok = insertAPIProductMappers(rows, txn)

			case common.Update:
				ok = deleteAPIproductMapper(payload.OldRow, txn)
				rows = append(rows, payload.NewRow)
				ok = insertAPIProductMappers(rows, txn)

			case common.Delete:
				ok = deleteAPIproductMapper(payload.OldRow, txn)
			}
		}
		if !ok {
			log.Error("Sql Operation error. Operation rollbacked")
			txn.Rollback()
			return
		}
	}
	txn.Commit()
	return
}

/*
 * DELETE OBJECT as passed in the input
 */
func deleteObject(object string, ele common.Row, txn *sql.Tx) bool {

	var scope, objid string
	ssql := "DELETE FROM " + object + " WHERE id = $1 AND _change_selector = $2"
	prep, err := txn.Prepare(ssql)
	if err != nil {
		log.Error("DELETE ", object, " Failed: ", err)
		return false
	}
	defer prep.Close()
	ele.Get("_change_selector", &scope)
	ele.Get("id", &objid)

	res, err := txn.Stmt(prep).Exec(objid, scope)
	if err == nil {
		affect, err := res.RowsAffected()
		if err == nil && affect != 0 {
			log.Debugf("DELETE %s (%s, %s) success.", object, objid, scope)
			return true
		}
	}
	log.Errorf("DELETE %s (%s, %s) failed.", object, objid, scope)
	return false

}

/*
 * DELETE  APIPRDT MAPPER
 */
func deleteAPIproductMapper(ele common.Row, txn *sql.Tx) bool {
	var ApiProduct, AppId, EntityIdentifier, apid_scope string

	prep, err := txn.Prepare("DELETE FROM APP_CREDENTIAL_APIPRODUCT_MAPPER WHERE apiprdt_id=$1 AND app_id=$2 AND appcred_id=$3 AND _change_selector=$4;")
	if err != nil {
		log.Error("DELETE APP_CREDENTIAL_APIPRODUCT_MAPPER Failed: ", err)
		return false
	}

	defer prep.Close()

	ele.Get("apiprdt_id", &ApiProduct)
	ele.Get("app_id", &AppId)
	ele.Get("appcred_id", &EntityIdentifier)
	ele.Get("_change_selector", &apid_scope)

	res, err := txn.Stmt(prep).Exec(ApiProduct, AppId, EntityIdentifier, apid_scope)
	if err == nil {
		affect, err := res.RowsAffected()
		if err == nil && affect != 0 {
			log.Debugf("DELETE APP_CREDENTIAL_APIPRODUCT_MAPPER (%s, %s, %s, %s) success.", ApiProduct, AppId, EntityIdentifier, apid_scope)
			return true
		}
	}
	log.Errorf("DELETE APP_CREDENTIAL_APIPRODUCT_MAPPER (%s, %s, %s, %s) failed.", ApiProduct, AppId, EntityIdentifier, apid_scope, err)
	return false
}

func deleteCompanyDeveloper(ele common.Row, txn *sql.Tx) bool {
	prep, err := txn.Prepare(`
	DELETE FROM COMPANY_DEVELOPER
	WHERE tenant_id=$1 AND company_id=$2 AND developer_id=$3`)
	if err != nil {
		log.Errorf("DELETE COMPANY_DEVELOPER Failed: %v", err)
		return false
	}
	defer prep.Close()

	var tenantId, companyId, developerId string
	ele.Get("tenant_id", &tenantId)
	ele.Get("company_id", &companyId)
	ele.Get("developer_id", &developerId)

	res, err := txn.Stmt(prep).Exec(tenantId, companyId, developerId)
	if err == nil {
		affect, err := res.RowsAffected()
		if err == nil && affect != 0 {
			log.Debugf("DELETE COMPANY_DEVELOPER (%s, %s, %s) success.", tenantId, companyId, developerId)
			return true
		}
	}
	log.Errorf("DELETE COMPANY_DEVELOPER (%s, %s, %s) failed: %v", tenantId, companyId, developerId, err)
	return false
}
