package apidVerifyApiKey

import (
	"strconv"
	"database/sql"
)

func convertSuffix(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func generateTestApiProduct(suffix int, txn *sql.Tx) {

	s, err := txn.Prepare("INSERT INTO kms_api_product (id, api_resources, environments, tenant_id, _change_selector) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("WHAT THE FUCK: " + err.Error())
	}
	s.Exec("api_product_" + convertSuffix(suffix), "{/**, /test}", "{Env_0, Env_1}",
		"tenant_id_xxxx", "Org_0")
}

func generateTestDeveloper(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_developer (id, status, email, first_name, last_name, tenant_id, _change_selector)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug : " + err.Error())
	}
	s.Exec("developer_id_" + convertSuffix(suffix), "Active", "test@apigee.com", "Apigee", "Google", "tenant_id_xxxx", "Org_0")
}

func generateTestCompany(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_company (id, status, name, display_name, tenant_id, _change_selector)" +
		"VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("company_id_" + convertSuffix(suffix), "Active", "Apigee Corporation", "Apigee", "tenant_id_xxxx", "Org_0")
}

func generateTestCompanyDeveloper(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_company_developer (developer_id, tenant_id, _change_selector, company_id)" +
		"VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("developer_id_" + convertSuffix(suffix), "tenant_id_0", "test_org0", "company_id_" + convertSuffix(suffix))
}

func generateTestApp(suffix1, suffix2 int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_app (id, developer_id, status, tenant_id, callback_url, _change_selector, parent_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("application_id_" + convertSuffix(suffix1), "developer_id_" + convertSuffix(suffix2), "Approved", "tenant_id_xxxx",
		"http://apigee.com", "Org_0", "developer_id_" + convertSuffix(suffix2))

}

func generateTestAppCompany(suffix1, suffix2 int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_app (id, company_id, status, tenant_id, callback_url, _change_selector, parent_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("application_id_" + convertSuffix(suffix1), "company_id_" + convertSuffix(suffix2), "Approved", "tenant_id_xxxx",
		"http://apigee.com", "Org_0", "company_id_" + convertSuffix(suffix2))

}

func generateTestAppCreds(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_app_credential (id, app_id, status, tenant_id, _change_selector) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("app_credential_" + convertSuffix(suffix), "application_id_" + convertSuffix(suffix), "Approved",
		"tenant_id_xxxx", "Org_0")
}

func generateTestApiProductMapper(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_app_credential_apiproduct_mapper (apiprdt_id, status, app_id, appcred_id, tenant_id, _change_selector) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("api_product_" + convertSuffix(suffix), "Approved", "application_id_" + convertSuffix(suffix),
		"app_credential_" + convertSuffix(suffix), "tenant_id_xxxx", "Org_0")
}