-- Copyright 2017 Google Inc.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--      http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE kms_attributes (tenant_id text,entity_id text,cust_id text,org_id text,dev_id text,comp_id text,apiprdt_id text,app_id text,appcred_id text,name text,type text,value text,_change_selector text, primary key (tenant_id,entity_id,name,type));
INSERT INTO "kms_attributes" VALUES('bc811169','50321842-d6ee-4e92-91b9-37234a7920c1','','','','','50321842-d6ee-4e92-91b9-37234a7920c1','','','RateLimit','APIPRODUCT','RX100','bc811169');
INSERT INTO "kms_attributes" VALUES('bc811169','85629786-37c5-4e8c-bb45-208f3360d005','','85629786-37c5-4e8c-bb45-208f3360d005','','','','','','features.isEdgexEnabled','ORGANIZATION','true','bc811169');
INSERT INTO "kms_attributes" VALUES('bc811169','85629786-37c5-4e8c-bb45-208f3360d005','','85629786-37c5-4e8c-bb45-208f3360d005','','','','','','features.isCpsEnabled','ORGANIZATION','true','bc811169');
INSERT INTO "kms_attributes" VALUES('bc811169','50321842-d6ee-4e92-91b9-37234a7920c1','','','','','50321842-d6ee-4e92-91b9-37234a7920c1','','','developer.quota.limit','APIPRODUCT','100','bc811169');
INSERT INTO "kms_attributes" VALUES('bc811169','50321842-d6ee-4e92-91b9-37234a7920c1','','','','','50321842-d6ee-4e92-91b9-37234a7920c1','','','developer.quota.interval','APIPRODUCT','10','bc811169');
INSERT INTO "kms_attributes" VALUES('bc811169','50321842-d6ee-4e92-91b9-37234a7920c1','','','','','50321842-d6ee-4e92-91b9-37234a7920c1','','','developer.quota.timeunit','APIPRODUCT','minute','bc811169');
INSERT INTO "kms_attributes" VALUES('bc811169','50321842-d6ee-4e92-91b9-37234a7920c1','','','','','50321842-d6ee-4e92-91b9-37234a7920c1','','','Threshold','APIPRODUCT','TX100','bc811169');
INSERT INTO "kms_attributes" VALUES('bc811169','40753e12-a50a-429d-9121-e571eb4e43a9','','','','','40753e12-a50a-429d-9121-e571eb4e43a9','','','access','APIPRODUCT','public','bc811169');
INSERT INTO "kms_attributes" VALUES('bc811169','2d373ed6-e38f-453b-bb34-6d731d9c4815','','','','','','2d373ed6-e38f-453b-bb34-6d731d9c4815','','DisplayName','APP','demo-app','bc811169');
CREATE TABLE kms_organization (id text,name text,display_name text,type text,tenant_id text,customer_id text,description text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));
INSERT INTO "kms_organization" VALUES('e2cc4caf-40d6-4ecb-8149-ed32d04184b2','apid-haoming','apid-haoming','paid','515211e9','94cd5075-7f33-4afb-9545-a53a254277a1','','2017-08-16 22:16:06.544+00:00','foobar@google.com','2017-08-16 22:29:23.046+00:00','foobar@google.com','515211e9');
CREATE TABLE edgex_data_scope (id text,apid_cluster_id text,scope text,org text,env text,created blob,created_by text,updated blob,updated_by text,_change_selector text,org_scope text,env_scope text, primary key (id));
INSERT INTO "edgex_data_scope" VALUES('cc066263-6355-416d-9d59-7f3135d64953','543230f1-8c41-4bf5-94a3-f10c104ff5d4','155211e9','apid-haoming','test','2017-08-27 22:53:33.859+00:00','foobar@google.com','2017-08-27 22:53:33.859+00:00','foobar@google.com','543230f1-8c41-4bf5-94a3-f10c104ff5d4','12344caf-40d6-4ecb-8149-ed32d04184b2','1234203e-ba88-4cd5-967d-4caa88f64909');
INSERT INTO "edgex_data_scope" VALUES('08c81eeb-57ec-43fe-8fed-cdff5494406f','123430f1-8c41-4bf5-94a3-f10c104ff5d4','165211e9','apid-test','prod','2017-08-29 02:39:34.093+00:00','foobar@google.com','2017-08-29 02:39:34.093+00:00','foobar@google.com','123430f1-8c41-4bf5-94a3-f10c104ff5d4','43214caf-40d6-4ecb-8149-ed32d04184b2','43211cae-f2a6-4663-9f36-eb17d76e6c32');
CREATE TABLE kms_company_developer (tenant_id text,company_id text,developer_id text,roles text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (tenant_id,company_id,developer_id));
CREATE TABLE kms_developer (id text,tenant_id text,username text,first_name text,last_name text,password text,email text,status text,encrypted_password text,salt text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));
CREATE TABLE kms_api_product (id text,tenant_id text,name text,display_name text,description text,api_resources text,approval_type text,scopes text,proxies text,environments text,quota text,quota_time_unit text,quota_interval integer,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));
CREATE TABLE kms_app (id text,tenant_id text,name text,display_name text,access_type text,callback_url text,status text,app_family text,company_id text,developer_id text,parent_id text,type text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));
CREATE TABLE kms_app_credential_apiproduct_mapper (tenant_id text,appcred_id text,app_id text,apiprdt_id text,status text,_change_selector text, primary key (tenant_id,appcred_id,app_id,apiprdt_id));
CREATE TABLE kms_company (id text,tenant_id text,name text,display_name text,status text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));
CREATE TABLE kms_app_credential (id text,tenant_id text,consumer_secret text,app_id text,method_type text,status text,issued_at blob,expires_at blob,app_status text,scopes text,created_at blob,created_by text,updated_at blob,updated_by text,_change_selector text, primary key (id,tenant_id));
COMMIT;
