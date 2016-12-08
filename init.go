package apidVerifyApiKey

import (
	"github.com/30x/apid"
	"sync"
)

const (
	apiPath = "/verifiers/apikey"
)

var (
	log    apid.LogService
	data   apid.DataService
	events apid.EventsService
	unsafeDB apid.DB
	dbMux    sync.RWMutex
)

func getDB() apid.DB {
	dbMux.RLock()
	db := unsafeDB
	dbMux.RUnlock()
	return db
}

func setDB(db apid.DB) {
	dbMux.Lock()
	unsafeDB = db
	dbMux.Unlock()
}

func init() {
	apid.RegisterPlugin(initPlugin)
}

func initPlugin(services apid.Services) (apid.PluginData, error) {
	log = services.Log().ForModule("apidVerifyAPIKey")
	log.Debug("start init")

	data = services.Data()
	events = services.Events()

	services.API().HandleFunc(apiPath, handleRequest).Methods("POST")

	events.Listen("ApigeeSync", &handler{})
	log.Debug("end init")

	return pluginData, nil
}

func createTables(db apid.DB) {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS api_product (
    id text,
    tenant_id text,
    name text,
    display_name text,
    description text,
    api_resources text[],
    approval_type text,
    _apid_scope text,
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
CREATE TABLE IF NOT EXISTS developer (
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
    _apid_scope text,
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    PRIMARY KEY (tenant_id, id)
);
CREATE TABLE IF NOT EXISTS company (
    id text,
    tenant_id text,
    name text,
    display_name text,
    status text,
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    _apid_scope text,
    PRIMARY KEY (tenant_id, id)
);
CREATE TABLE IF NOT EXISTS company_developer (
     tenant_id text,
     company_id text,
     developer_id text,
    roles text[],
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    _apid_scope text,
    PRIMARY KEY (tenant_id, company_id,developer_id)
);
CREATE TABLE IF NOT EXISTS app (
    id text,
    tenant_id text,
    name text,
    display_name text,
    access_type text,
    callback_url text,
    status text,
    app_family text,
    company_id text,
    developer_id text,
    type int,
    created_at int64,
    created_by text,
    updated_at int64,
    updated_by text,
    _apid_scope text,
    PRIMARY KEY (tenant_id, id)
);
CREATE TABLE IF NOT EXISTS app_credential (
    id text,
    tenant_id text,
    consumer_secret text,
    app_id text,
    method_type text,
    status text,
    issued_at int64,
    expires_at int64,
    app_status text,
    _apid_scope text,
    PRIMARY KEY (tenant_id, id)
);
CREATE TABLE IF NOT EXISTS app_credential_apiproduct_mapper (
    tenant_id text,
    appcred_id text,
    app_id text,
    apiprdt_id text,
    _apid_scope text,
    status text,
    PRIMARY KEY (appcred_id, app_id, apiprdt_id,tenant_id)
);
`)
	if err != nil {
		log.Panic("Unable to initialize DB", err)
	}
}
