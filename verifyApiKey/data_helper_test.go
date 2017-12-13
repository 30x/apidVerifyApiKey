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

package verifyApiKey

import (
	"github.com/apid/apid-core"
	. "github.com/onsi/gomega"
)

//initialize DB for tests
func setupApikeyDeveloperTestDb(db apid.DB) {
	tx, err := db.Begin()
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_organization (id text,name text,display_name text,type text,tenant_id text,customer_id text,description text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_organization" VALUES('85629786-37c5-4e8c-bb45-208f3360d005','apigee-mcrosrvc-client0001','apigee-mcrosrvc-client0001','trial','bc811169','2277ba6c-8991-4a38-a5fc-12d8d36e5812','','2017-07-03 19:21:09.388+00:00','defaultUser','2017-07-05 16:24:35.413+00:00','rajanish@apigee.com','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE kms_developer (id text,tenant_id text,username text,first_name text,last_name text,password text,email text,status text,encrypted_password text,salt text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_developer" VALUES('209ffd18-37e9-4a67-9e30-a5c40a534b6c','bc811169','wilson','Woodre','Wilson','','developer@apigee.com','ACTIVE','','','2017-08-08 17:24:09.008+00:00','defaultUser','2017-08-08 17:24:09.008+00:00','defaultUser','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_company (id text,tenant_id text,name text,display_name text,status text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_company" VALUES('7834c683-9453-4389-b816-34ca24dfccd9','bc811169','DevCompany','East India Company','ACTIVE','2017-08-05 19:54:12.359+00:00','defaultUser','2017-08-05 19:54:12.359+00:00','defaultUser','bc811169');`)

	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_app (id text,tenant_id text,name text,display_name text,access_type text,callback_url text,status text,app_family text,company_id text,developer_id text,parent_id text,type text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_app" VALUES('d371f05a-7c04-430c-b12d-26cf4e4d5d65','bc811169','DeveloperApp','','READ','www.apple.com','APPROVED','default','','209ffd18-37e9-4a67-9e30-a5c40a534b6c','209ffd18-37e9-4a67-9e30-a5c40a534b6c','DEVELOPER','2017-08-07 17:00:54.25+00:00','defaultUser','2017-08-07 17:09:08.259+00:00','defaultUser','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE kms_api_product (id text,tenant_id text,name text,display_name text,description text,api_resources text,approval_type text,scopes text,proxies text,environments text,quota text,quota_time_unit text,quota_interval integer,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_api_product" VALUES('24987a63-edb9-4d6b-9334-87e1d70df8e3','bc811169','KeyProduct4','Sandbox Diamond','','{/zoho,/twitter,/nike}','AUTO','{READ,WRITE}','{DevApplication,KeysApplication}','{test}','','',NULL,'2017-08-08 02:53:32.726+00:00','defaultUser','2017-08-08 02:53:32.726+00:00','defaultUser','bc811169')`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_app_credential (id text,tenant_id text,consumer_secret text,app_id text,method_type text,status text,issued_at blob,expires_at blob,app_status text,scopes text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO kms_app_credential VALUES('63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0','bc811169','Ui8dcyGW3lA04YdX','d371f05a-7c04-430c-b12d-26cf4e4d5d65','','APPROVED','2017-08-07 17:00:54.258+00:00','','','{DELETE}','2017-08-07 17:00:54.258+00:00','-NA-','2017-08-07 17:06:06.242+00:00','-NA-','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_app_credential_apiproduct_mapper (tenant_id text,appcred_id text,app_id text,apiprdt_id text,status text,_change_selector text, primary key (tenant_id,appcred_id,app_id,apiprdt_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_app_credential_apiproduct_mapper" VALUES('bc811169','63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0','d371f05a-7c04-430c-b12d-26cf4e4d5d65','24987a63-edb9-4d6b-9334-87e1d70df8e3','APPROVED','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_attributes (tenant_id text,entity_id text,cust_id text,org_id text,dev_id text,comp_id text,apiprdt_id text,app_id text,appcred_id text,name text,type text,value text,_change_selector text, primary key (tenant_id,entity_id,name,type));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_attributes" VALUES('bc811169','d371f05a-7c04-430c-b12d-26cf4e4d5d65','','','','','','d371f05a-7c04-430c-b12d-26cf4e4d5d65','','Company','APP','Apple','bc811169');`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_attributes" VALUES('bc811169','7834c683-9453-4389-b816-34ca24dfccd9','','','','7834c683-9453-4389-b816-34ca24dfccd9','','','','country','COMPANY','england','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`INSERT INTO "kms_attributes" VALUES('bc811169','63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0','','','','','','','63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0','Device','APP_CREDENTIAL','ios','bc811169');`)
	Expect(err).NotTo(HaveOccurred())
	Expect(tx.Commit()).NotTo(HaveOccurred())

}

//initialize DB for tests
func setupApikeyCompanyTestDb(db apid.DB) {
	tx, err := db.Begin()
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_organization (id text,name text,display_name text,type text,tenant_id text,customer_id text,description text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_organization" VALUES('85629786-37c5-4e8c-bb45-208f3360d005','apigee-mcrosrvc-client0001','apigee-mcrosrvc-client0001','trial','bc811169','2277ba6c-8991-4a38-a5fc-12d8d36e5812','','2017-07-03 19:21:09.388+00:00','defaultUser','2017-07-05 16:24:35.413+00:00','rajanish@apigee.com','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE kms_developer (id text,tenant_id text,username text,first_name text,last_name text,password text,email text,status text,encrypted_password text,salt text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_company (id text,tenant_id text,name text,display_name text,status text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_company" VALUES('7834c683-9453-4389-b816-34ca24dfccd9','bc811169','DevCompany','East India Company','ACTIVE','2017-08-05 19:54:12.359+00:00','defaultUser','2017-08-05 19:54:12.359+00:00','defaultUser','bc811169');`)

	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_app (id text,tenant_id text,name text,display_name text,access_type text,callback_url text,status text,app_family text,company_id text,developer_id text,parent_id text,type text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_app" VALUES('d371f05a-7c04-430c-b12d-26cf4e4d5d65','bc811169','CompApp2','','READ','www.apple.com','APPROVED','default','7834c683-9453-4389-b816-34ca24dfccd9','','7834c683-9453-4389-b816-34ca24dfccd9','COMPANY','2017-08-07 17:00:54.25+00:00','defaultUser','2017-08-07 17:09:08.259+00:00','defaultUser','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE kms_api_product (id text,tenant_id text,name text,display_name text,description text,api_resources text,approval_type text,scopes text,proxies text,environments text,quota text,quota_time_unit text,quota_interval integer,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_api_product" VALUES('24987a63-edb9-4d6b-9334-87e1d70df8e3','bc811169','KeyProduct4','Sandbox Diamond','','{/zoho,/twitter,/nike}','AUTO','{READ,WRITE}','{DevApplication,KeysApplication}','{test}','','',NULL,'2017-08-08 02:53:32.726+00:00','defaultUser','2017-08-08 02:53:32.726+00:00','defaultUser','bc811169')`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_app_credential (id text,tenant_id text,consumer_secret text,app_id text,method_type text,status text,issued_at blob,expires_at blob,app_status text,scopes text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO kms_app_credential VALUES('63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0','bc811169','Ui8dcyGW3lA04YdX','d371f05a-7c04-430c-b12d-26cf4e4d5d65','','APPROVED','2017-08-07 17:00:54.258+00:00','','','{DELETE}','2017-08-07 17:00:54.258+00:00','-NA-','2017-08-07 17:06:06.242+00:00','-NA-','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_app_credential_apiproduct_mapper (tenant_id text,appcred_id text,app_id text,apiprdt_id text,status text,_change_selector text, primary key (tenant_id,appcred_id,app_id,apiprdt_id));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_app_credential_apiproduct_mapper" VALUES('bc811169','63tHSNLKJkcc6GENVWGT1Zw5gek7kVJ0','d371f05a-7c04-430c-b12d-26cf4e4d5d65','24987a63-edb9-4d6b-9334-87e1d70df8e3','APPROVED','bc811169');`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kms_attributes (tenant_id text,entity_id text,cust_id text,org_id text,dev_id text,comp_id text,apiprdt_id text,app_id text,appcred_id text,name text,type text,value text,_change_selector text, primary key (tenant_id,entity_id,name,type));`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_attributes" VALUES('bc811169','d371f05a-7c04-430c-b12d-26cf4e4d5d65','','','','','','d371f05a-7c04-430c-b12d-26cf4e4d5d65','','Company','APP','Apple','bc811169');`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`INSERT INTO "kms_attributes" VALUES('bc811169','7834c683-9453-4389-b816-34ca24dfccd9','','','','7834c683-9453-4389-b816-34ca24dfccd9','','','','country','COMPANY','england','bc811169');`)
	Expect(err).NotTo(HaveOccurred())
	Expect(tx.Commit()).NotTo(HaveOccurred())
}
