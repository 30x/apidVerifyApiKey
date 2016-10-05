package apidVerifyApiKey

import (
	"database/sql"
	"github.com/30x/apid"
	"github.com/30x/apidApigeeSync"
)

const (
	apiPath = "/verifyAPIKey"
)

var (
	log    apid.LogService
	data   apid.DataService
	events apid.EventsService
)

func init() {
	apid.RegisterPlugin(initPlugin)
}

func initPlugin(services apid.Services) error {
	log = services.Log().ForModule("apidVerifyAPIKey")
	log.Debug("start init")

	data = services.Data()
	events = services.Events()

	db, err := data.DB()
	if err != nil {
		log.Panic("Unable to access DB", err)
	}

	var count int
	row := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='API_PRODUCT';")
	if err := row.Scan(&count); err != nil {
		log.Panic("Unable to setup database", err)
	}
	if count == 0 {
		createTables(db)
	}

	services.API().HandleFunc(apiPath, handleRequest)

	events.Listen(apidApigeeSync.ApigeeSyncEventSelector, &handler{})
	log.Debug("end init")

	return nil
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`
CREATE TABLE api_product (
    id uuid,
    tenant_id text,
    name text,
    display_name text,
    description text,
    api_resources text[],
    approval_type text,
    scopes text[],
    proxies text[],
    environments text[],
    quota text,
    quota_time_unit text,
    quota_interval int,
    created_at timestamp,
    created_by text,
    updated_at timestamp,
    updated_by text,
    PRIMARY KEY (tenant_id, id));
CREATE TABLE developer (
    id uuid,
    tenant_id text,
    username text,
    first_name text,
    last_name text,
    password text,
    email text,
    status text,
    encrypted_password text,
    salt text,
    created_at timestamp,
    created_by text,
    updated_at timestamp,
    updated_by text,
    PRIMARY KEY (tenant_id, id),
    constraint developer_email_uq unique(tenant_id, email)
);
CREATE TABLE company (
    id uuid,
    tenant_id text,
    name text,
    display_name text,
    status text,
    created_at timestamp,
    created_by text,
    updated_at timestamp,
    updated_by text,
    PRIMARY KEY (tenant_id, id),
    constraint comp_name_uq unique(tenant_id, name)
);
CREATE TABLE company_developer (
     tenant_id text,
     company_id uuid,
     developer_id uuid,
    roles text[],
    created_at timestamp,
    created_by text,
    updated_at timestamp,
    updated_by text,
    PRIMARY KEY (tenant_id, company_id,developer_id),
    FOREIGN KEY (tenant_id,company_id) references company(tenant_id,id),
    FOREIGN KEY (tenant_id,developer_id) references developer(tenant_id,id)
);
CREATE TABLE app (
    id uuid,
    tenant_id text,
    name text,
    display_name text,
    access_type text,
    callback_url text,
    status text,
    app_family text,
    company_id uuid,
    developer_id uuid,
    type app_type,
    created_at timestamp,
    created_by text,
    updated_at timestamp,
    updated_by text,
    PRIMARY KEY (tenant_id, id),
    constraint app_name_uq unique(tenant_id, name),
    FOREIGN KEY (tenant_id,company_id) references company(tenant_id,id),
    FOREIGN KEY (tenant_id,developer_id) references developer(tenant_id,id)
);
CREATE TABLE app_credential (
    id text,
    tenant_id text,
    consumer_secret text,
    app_id uuid,
    method_type text,
    status text,
    issued_at timestamp,
    expires_at timestamp,
    app_status text,
    scopes text[],
    PRIMARY KEY (tenant_id, id),
    FOREIGN KEY (tenant_id,app_id) references app(tenant_id,id)
);
CREATE TABLE app_credential_apiproduct_mapper (
    tenant_id text,
    appcred_id text,
    app_id uuid,
    apiprdt_id uuid,
    status appcred_apiprdt_status,
    PRIMARY KEY (tenant_id,appcred_id,app_id,apiprdt_id),
    FOREIGN KEY (tenant_id,appcred_id) references app_credential(tenant_id,id),
    FOREIGN KEY (tenant_id,app_id) references app(tenant_id,id)
);
CREATE TABLE attributes (
   tenant_id text,
   dev_id uuid,
   comp_id uuid,
   apiprdt_id uuid,
   app_id uuid,
   appcred_id text,
   type entity_type,
   name text ,
   value text,
   PRIMARY KEY (tenant_id,dev_id,comp_id,apiprdt_id,app_id,appcred_id,type,name),
   FOREIGN KEY (tenant_id,appcred_id) references app_credential(tenant_id,id),
   FOREIGN KEY (tenant_id,app_id) references app(tenant_id,id),
   FOREIGN KEY (tenant_id,dev_id) references developer(tenant_id,id),
   FOREIGN KEY (tenant_id,comp_id) references company(tenant_id,id),
   FOREIGN KEY (tenant_id,apiprdt_id) references api_product(tenant_id,id)
);
CREATE TABLE apidconfig (
    id uuid,
    consumer_key text,
    consumer_secret text,
    scope text[],
    app_id uuid,
    created_at timestamp,
    created_by text,
    updated_at timestamp,
    updated_by text,
    PRIMARY KEY(id),
    constraint apidconfig_key_uq unique(consumer_key),
    constraint apidconfig_appid_uq unique(app_id)
);
`)
	if err != nil {
		log.Panic("Unable to initialize DB", err)
	}
}
