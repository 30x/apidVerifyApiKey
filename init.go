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
	row := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table';")
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
	_, err := db.Exec("CREATE TABLE COMPANY (org varchar(255), id varchar(255), PRIMARY KEY (id, org));CREATE TABLE DEVELOPER (org varchar(255), email varchar(255), id varchar(255), sts varchar(255), username varchar(255), firstname varchar(255), lastname varchar(255), apigee_scope varchar(255), enc_password varchar(255), salt varchar(255), created_at integer, created_by varchar(255), updated_at integer, updated_by varchar(255), PRIMARY KEY (id, org));CREATE TABLE APP (org varchar(255), id varchar(255), dev_id varchar(255) null, cmp_id varchar(255) null, display_name varchar(255), apigee_scope varchar(255), type varchar(255), access_type varchar(255), cback_url varchar(255), status varchar(255), name varchar(255), app_family varchar(255), created_at integer, created_by varchar(255), updated_at integer, updated_by varchar(255), PRIMARY KEY (id, org), FOREIGN KEY (dev_id, org) references DEVELOPER (id, org) ON DELETE CASCADE);CREATE TABLE APP_CREDENTIAL (org varchar(255), id varchar(255), app_id varchar(255), cons_secret varchar(255), method_type varchar(255), status varchar(255), issued_at integer, expire_at integer, created_at integer, created_by varchar(255), updated_at integer, updated_by varchar(255), PRIMARY KEY (id, org), FOREIGN KEY (app_id, org) references app (id, org) ON DELETE CASCADE);CREATE TABLE API_PRODUCT (org varchar(255), id varchar(255), res_names varchar(255), env varchar(255), PRIMARY KEY (id, org));CREATE TABLE COMPANY_DEVELOPER (org varchar(255), dev_id varchar(255), id varchar(255), cmpny_id varchar(255), PRIMARY KEY (id, org), FOREIGN KEY (cmpny_id) references company(id) ON DELETE CASCADE, FOREIGN KEY (dev_id, org) references DEVELOPER(id, org) ON DELETE CASCADE);CREATE TABLE APP_AND_API_PRODUCT_MAPPER (org varchar(255), api_prdt_id varchar(255), app_id varchar(255), app_cred_id varchar(255), api_prdt_status varchar(255), PRIMARY KEY (org, api_prdt_id, app_id, app_cred_id), FOREIGN KEY (api_prdt_id, org) references api_product(id, org) ON DELETE CASCADE, FOREIGN KEY (app_cred_id, org) references app_credential(id, org) ON DELETE CASCADE, FOREIGN KEY (app_id, org) references app(id, org) ON DELETE CASCADE);")
	if err != nil {
		log.Panic("Unable to initialize DB", err)
	}
}
