package apidVerifyApiKey

import (
	"database/sql"
	"github.com/30x/apid-core"
	"github.com/apigee-labs/transicator/common"
	"strings"
	"fmt"
	"sort"
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
				ok = insert("developer", payload.Rows, txn)
			case "kms.app":
				ok = insert("app", payload.Rows, txn)
			case "kms.app_credential":
				ok = insert("app_credential", payload.Rows, txn)
			case "kms.api_product":
				ok = insert("api_product", payload.Rows, txn)
			case "kms.app_credential_apiproduct_mapper":
				ok = insert("app_credential_apiproduct_mapper", payload.Rows, txn)
			case "kms.company":
				ok = insert("company", payload.Rows, txn)
			case "kms.company_developer":
				ok = insert("company_developer", payload.Rows, txn)
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

func processChange(changes *common.ChangeList) {

	db := getDB()

	txn, err := db.Begin()
	if err != nil {
		log.Error("Unable to create transaction")
		return
	}
	defer txn.Rollback()

	ok := true

	log.Debugf("apigeeSyncEvent: %d changes", len(changes.Changes))
	for _, payload := range changes.Changes {
		var newrows = []common.Row{payload.NewRow}
		var oldrows = []common.Row{payload.OldRow}
		switch payload.Table {
		case "kms.developer":
			switch payload.Operation {
			case common.Insert:
				ok = insert("developer", newrows, txn)
			case common.Update:
				ok = delete("DEVELOPER", oldrows, txn)
				ok = ok && insert("developer", newrows, txn)
			case common.Delete:
				ok = delete("DEVELOPER", oldrows, txn)
			}
		case "kms.app":
			switch payload.Operation {
			case common.Insert:
				ok = insert("app", newrows, txn)
			case common.Update:
				ok = delete("APP", oldrows, txn)
				ok = ok && insert("app", newrows, txn)
			case common.Delete:
				ok = delete("APP", oldrows, txn)
			}
		case "kms.company":
			switch payload.Operation {
			case common.Insert:
				ok = insert("company", newrows, txn)
			case common.Update:
				ok = delete("COMPANY", oldrows, txn)
				ok = ok && insert("company", newrows, txn)
			case common.Delete:
				ok = delete("COMPANY", oldrows, txn)
			}
		case "kms.company_developer":
			switch payload.Operation {
			case common.Insert:
				ok = insert("company_developer", newrows, txn)
			case common.Update:
				ok = delete("company_developer", oldrows, txn)
				ok = ok && insert("company_developer", newrows, txn)
			case common.Delete:
				ok = delete("company_developer", oldrows, txn)
			}
		case "kms.app_credential":
			switch payload.Operation {
			case common.Insert:
				ok = insert("app_credential", newrows, txn)
			case common.Update:
				ok = delete("APP_CREDENTIAL", oldrows, txn)
				ok = ok && insert("app_credential", newrows, txn)
			case common.Delete:
				ok = delete("APP_CREDENTIAL", oldrows, txn)
			}
		case "kms.api_product":
			switch payload.Operation {
			case common.Insert:
				ok = insert("api_product", newrows, txn)
			case common.Update:
				ok = delete("API_PRODUCT", oldrows, txn)
				ok = ok && insert("api_product", newrows, txn)
			case common.Delete:
				ok = delete("API_PRODUCT", oldrows, txn)
			}

		case "kms.app_credential_apiproduct_mapper":
			switch payload.Operation {
			case common.Insert:
				ok = insert("app_credential_apiproduct_mapper", newrows, txn)
			case common.Update:
				ok = delete("app_credential_apiproduct_mapper", oldrows, txn)
				ok = ok && insert("app_credential_apiproduct_mapper", newrows, txn)
			case common.Delete:
				ok = delete("app_credential_apiproduct_mapper", oldrows, txn)
			}
		}
		if !ok {
			log.Error("Sql Operation error. Operation rollbacked")
			return
		}
	}
	txn.Commit()
	return
}

func delete(tableName string, rows[] common.Row, txn *sql.Tx) bool {
	pkeys, err := getPkeysForTable(tableName)
	if (len(pkeys) == 0 || err != nil) {
		log.Errorf("DELETE No primary keys found for table.", tableName)
		return false
	} else {
		sql := buildDeleteSql(tableName, pkeys);
		prep, err := txn.Prepare(sql)
		if err != nil {
			log.Errorf("DELETE Fail to prep statement [%s] error=[%v]", sql, err)
			return false
		}
		defer prep.Close()
		for _, row := range rows {
			var values []interface{};
			for _, pkey := range pkeys {
				var value interface{}
				row.Get(pkey, &value)
				values = append(values, value)
			}
			_, err = prep.Exec(values)

			if err != nil {
				log.Errorf("DELETE Fail [%s] value=[%v] error=[%v]", sql, values, err)
				return false
			} else {
				log.Debugf("DELETE Success [%s] value=[%v]", sql, values)
			}
		}
		return true
	}
}

func update(tableName string, oldRows, newRows []common.Row, txn *sql.Tx) bool {
	pkeys, err := getPkeysForTable(tableName)
	if (len(pkeys) == 0 || err != nil) {
		log.Errorf("UPDATE No primary keys found for table.", tableName)
		return false
	} else {
		if len(oldRows) == 0 || len(newRows) == 0 {
			return false
		}

		var columnNames []string

		//extract columnNames
		for columnName := range oldRows[0] {
			columnNames = append(columnNames, columnName)
		}
		sort.Strings(columnNames)


		//build update statement, use arbitrary row as template
		sql := buildUpdateSql(tableName, newRows[0], pkeys);
		prep, err := txn.Prepare(sql)

		for i, row := range newRows {
			if err != nil {
				log.Errorf("UPDATE Fail to prep statement [%s] error=[%v]", sql, err)
				return false
			}
			defer prep.Close()
			var values []interface{};

			//sort to ensure order parity with sql creation, add values for set clause
			for _, columnName := range columnNames {
				values = append(values, row[columnName])
			}

			//add values for where clause, use PKs of old row
			for _, pk := range pkeys {
				values = append(values, oldRows[i][pk])

			}

			//create prepared statement from existing template statement
			_, err = txn.Stmt(prep).Exec(values)

			if err != nil {
				log.Errorf("UPDATE Fail [%s] value=[%v] error=[%v]", sql, values, err)
				return false
			} else {
				log.Debugf("UPDATE Success [%s] value=[%v]", sql, values)
			}
		}

		return true
	}
}

func buildDeleteSql(tableName string, pkeys []string) string {

	normalizedTableName := normalizeTableName(tableName)
	var clauses []string
	for i, columnName := range pkeys {
		clauses = append(clauses, fmt.Sprint(columnName, "= $", (i + 1)))
	}

	sql := []string{"DELETE FROM ", normalizedTableName, "WHERE", strings.Join(clauses, "AND"), ";"}
	return strings.Join(sql, " ")
}

func getPkeysForTable(tableName string) ([]string, error) {
	db := getDB()
	normalizedTableName := normalizeTableName(tableName)
	sql := "SELECT columnName FROM _transicator_tables WHERE tableName = $1 ORDER BY columnName"
	rows, err := db.Query(sql, normalizedTableName)
	if err != nil {
		log.Errorf("Failed [%s] values=[s%] Error: %v", sql, normalizedTableName, err)
		return nil, err
	}
	var columnNames []string
	defer rows.Close()
	for rows.Next() {
		var value interface{}
		err := rows.Scan(&value)
		if err != nil {
			log.Fatal(err)
		}
		columnNames = append(columnNames, fmt.Sprint(value))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return columnNames, nil;
}

func buildUpdateSql(tableName string, row common.Row, pkeys []string) string {
	if row == nil{
		return ""
	}
	normalizedTableName := strings.Replace(tableName, ".", "_", 0)

	var columns, setPlaceholders, wherePlaceholders []string
	i := 1

	for columnName := range row {
		columns = append(columns, columnName)

	}

	sort.Strings(columns)
	for _, columnName := range columns {
		setPlaceholders = append(setPlaceholders, fmt.Sprintf("%s=$%v", columnName, i))
		i++
	}

	for _, pk := range pkeys {
		wherePlaceholders = append(wherePlaceholders, fmt.Sprintf("%s=$%v", pk, i))
		i++
	}

	sql := "UPDATE " + normalizedTableName + " SET "
	sql = sql + strings.Join(setPlaceholders, ", ")
	sql = sql + " WHERE "
	sql = sql + strings.Join(wherePlaceholders, " AND ")

	return sql
}

func buildInsertSql(tableName string, rows []common.Row) string {
	if len(rows) == 0 {
		return ""
	}
	normalizedTableName := normalizeTableName(tableName)
	row := rows[0]
	var columns, placeholders []string
	i := 1
	for columnName := range row {
		columns = append(columns, columnName)
		placeholders = append(placeholders, fmt.Sprint("$", i))
		i++
	}

	sql := []string{"INSERT INTO", normalizedTableName, "(", strings.Join(columns, ","), ")",
		"VALUES", "(", strings.Join(placeholders, ","), ");"}
	return strings.Join(sql, " ")
}

func normalizeTableName(tableName string) string {
	if (strings.Contains(tableName, ".")) {
		split := strings.Split(tableName, ".")
		return split[len(split) - 1];
	}
	return tableName;
}

func insert(tableName string, rows []common.Row, txn *sql.Tx) bool {

	sql := buildInsertSql(tableName, rows)

	prep, err := txn.Prepare(sql)
	if err != nil {
		log.Errorf("INSERT Fail to prepare statement [%s] error=[%v]", sql, err)
		return false
	}
	defer prep.Close()
	for _, ele := range rows {
		var values []interface{};
		for _, value := range ele {
			values = append(values, value)
		}

		_, err = prep.Exec(values)

		if err != nil {
			log.Errorf("INSERT Fail [%s] value=[%v] error=[%v]", sql, values, err)
			return false
		} else {
			log.Debugf("INSERT Success [%s] value=[%v]", sql, values)
		}
	}
	return true
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
		_, err = prep.Exec(
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

	var scope, EntityIdentifier, DeveloperId, CompanyId, ParentId, CallbackUrl, Status, AppName, AppFamily, tenantId, CreatedBy, LastModifiedBy string
	var CreatedAt, LastModifiedAt int64

	prep, err := txn.Prepare("INSERT INTO APP (_change_selector, id, developer_id, company_id, parent_id, callback_url,status, name, app_family, created_at, created_by,updated_at, updated_by,tenant_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);")
	if err != nil {
		log.Error("INSERT APP Failed: ", err)
		return false
	}

	defer prep.Close()
	for _, ele := range rows {

		ele.Get("_change_selector", &scope)
		ele.Get("id", &EntityIdentifier)
		ele.Get("developer_id", &DeveloperId)
		ele.Get("company_id", &CompanyId)
		ele.Get("parent_id", &ParentId)
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
		_, err = prep.Exec(
			scope,
			EntityIdentifier,
			DeveloperId,
			CompanyId,
			ParentId,
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
		_, err = prep.Exec(
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
		if scope == "" || tenantId == "" || CompanyId == "" || DeveloperId == "" {
			log.Error("INSERT COMPANY_DEVELOPER: i/p args missing")
			return false
		}
		_, err = prep.Exec(
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
		_, err = prep.Exec(
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
		_, err = prep.Exec(
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

		_, err = prep.Exec(
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

	res, err := prep.Exec(ApiProduct, AppId, EntityIdentifier, apid_scope)
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

	res, err := prep.Exec(tenantId, companyId, developerId)
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
