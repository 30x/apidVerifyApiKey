// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apidVerifyApiKey

import (
	"database/sql"
	"github.com/30x/apid-core"
	"strconv"
)

func convertSuffix(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func generateTestApiProduct(suffix int, txn *sql.Tx) {

	s, err := txn.Prepare("INSERT INTO kms_api_product (id, api_resources, environments, tenant_id, _change_selector) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug : " + err.Error())
	}
	s.Exec("api_product_"+convertSuffix(suffix), "{/**, /test}", "{Env_0, Env_1}",
		"tenant_id_xxxx", "Org_0")
}

func generateTestDeveloper(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_developer (id, status, email, first_name, last_name, tenant_id, _change_selector)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug : " + err.Error())
	}
	s.Exec("developer_id_"+convertSuffix(suffix), "Active", "test@apigee.com", "Apigee", "Google", "tenant_id_xxxx", "Org_0")
}

func generateTestCompany(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_company (id, status, name, display_name, tenant_id, _change_selector)" +
		"VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("company_id_"+convertSuffix(suffix), "Active", "Apigee Corporation", "Apigee", "tenant_id_xxxx", "Org_0")
}

func generateTestCompanyDeveloper(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_company_developer (developer_id, tenant_id, _change_selector, company_id)" +
		"VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("developer_id_"+convertSuffix(suffix), "tenant_id_0", "test_org0", "company_id_"+convertSuffix(suffix))
}

func generateTestApp(suffix1, suffix2 int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_app (id, developer_id, status, tenant_id, callback_url, _change_selector, parent_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("application_id_"+convertSuffix(suffix1), "developer_id_"+convertSuffix(suffix2), "Approved", "tenant_id_xxxx",
		"http://apigee.com", "Org_0", "developer_id_"+convertSuffix(suffix2))

}

func generateTestAppCompany(suffix1, suffix2 int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_app (id, company_id, status, tenant_id, callback_url, _change_selector, parent_id)" +
		" VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("application_id_"+convertSuffix(suffix1), "company_id_"+convertSuffix(suffix2), "Approved", "tenant_id_xxxx",
		"http://apigee.com", "Org_0", "company_id_"+convertSuffix(suffix2))

}

func generateTestAppCreds(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_app_credential (id, app_id, status, tenant_id, _change_selector) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("app_credential_"+convertSuffix(suffix), "application_id_"+convertSuffix(suffix), "Approved",
		"tenant_id_xxxx", "Org_0")
}

func generateTestApiProductMapper(suffix int, txn *sql.Tx) {
	s, err := txn.Prepare("INSERT INTO kms_app_credential_apiproduct_mapper (apiprdt_id, status, app_id, appcred_id, tenant_id, _change_selector) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Panicf("This is a bug: " + err.Error())
	}
	s.Exec("api_product_"+convertSuffix(suffix), "Approved", "application_id_"+convertSuffix(suffix),
		"app_credential_"+convertSuffix(suffix), "tenant_id_xxxx", "Org_0")
}

func createTables(db apid.DB) {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS kms_api_product (
    id text,
    tenant_id text,
    name text,
    display_name text,
    description text,
    api_resources text[],
    approval_type text,
    _change_selector text,
    proxies text[],
    environments text[],
    quota text,
    quota_time_unit text,
    quota_interval int,
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    PRIMARY KEY (tenant_id, id));
CREATE TABLE IF NOT EXISTS kms_developer (
    id text,
    tenant_id text,
    username text,
    first_name text,
    last_name text,
    password text,
    email text,
    status text,
    encrypted_password text,
    salt text,
    _change_selector text,
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    PRIMARY KEY (tenant_id, id)
);
CREATE TABLE IF NOT EXISTS kms_company (
    id text,
    tenant_id text,
    name text,
    display_name text,
    status text,
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    _change_selector text,
    PRIMARY KEY (tenant_id, id)
);
CREATE TABLE IF NOT EXISTS kms_company_developer (
     tenant_id text,
     company_id text,
     developer_id text,
    roles text[],
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    _change_selector text,
    PRIMARY KEY (tenant_id, company_id,developer_id)
);
CREATE TABLE IF NOT EXISTS kms_app (
    id text,
    tenant_id text,
    name text,
    display_name text,
    access_type text,
    callback_url text,
    status text,
    app_family text,
    company_id text,
    parent_id text,
    developer_id text,
    type int,
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    _change_selector text,
    PRIMARY KEY (tenant_id, id)
);
CREATE TABLE IF NOT EXISTS kms_app_credential (
    id text,
    tenant_id text,
    consumer_secret text,
    app_id text,
    method_type text,
    status text,
    issued_at int64,
    expires_at int64,
    app_status text,
    _change_selector text,
    PRIMARY KEY (tenant_id, id)
);
CREATE TABLE IF NOT EXISTS kms_app_credential_apiproduct_mapper (
    tenant_id text,
    appcred_id text,
    app_id text,
    apiprdt_id text,
    _change_selector text,
    status text,
    PRIMARY KEY (appcred_id, app_id, apiprdt_id,tenant_id)
);
CREATE INDEX IF NOT EXISTS company_id ON kms_company (id);
CREATE INDEX IF NOT EXISTS developer_id ON kms_developer (id);
CREATE INDEX IF NOT EXISTS api_product_id ON kms_api_product (id);
CREATE INDEX IF NOT EXISTS app_id ON kms_app (id);
`)
	if err != nil {
		log.Panic("Unable to initialize DB", err)
	}
}

func createApidClusterTables(db apid.DB) {
	_, err := db.Exec(`
CREATE TABLE edgex_apid_cluster (
    id text,
    instance_id text,
    name text,
    description text,
    umbrella_org_app_name text,
    created int64,
    created_by text,
    updated int64,
    updated_by text,
    _change_selector text,
    snapshotInfo text,
    lastSequence text,
    PRIMARY KEY (id)
);
CREATE TABLE edgex_data_scope (
    id text,
    apid_cluster_id text,
    scope text,
    org text,
    env text,
    created int64,
    created_by text,
    updated int64,
    updated_by text,
    _change_selector text,
    PRIMARY KEY (id)
);
`)
	if err != nil {
		log.Panic("Unable to initialize DB", err)
	}
}

func addScopes(db apid.DB) {
	txn, _ := db.Begin()
	txn.Exec("INSERT INTO EDGEX_DATA_SCOPE (id, _change_selector, apid_cluster_id, scope, org, env) "+
		"VALUES"+
		"($1,$2,$3,$4,$5,$6)",
		"ABCDE",
		"some_cluster_id",
		"some_cluster_id",
		"tenant_id_xxxx",
		"test_org0",
		"Env_0",
	)
	txn.Exec("INSERT INTO EDGEX_DATA_SCOPE (id, _change_selector, apid_cluster_id, scope, org, env) "+
		"VALUES"+
		"($1,$2,$3,$4,$5,$6)",
		"XYZ",
		"test_org0",
		"somecluster_id",
		"tenant_id_0",
		"test_org0",
		"Env_0",
	)
	log.Info("Inserted EDGEX_DATA_SCOPE for test")
	txn.Commit()
}
