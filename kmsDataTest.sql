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

INSERT INTO "kms"."customer" (id,name,display_name,description,created_at,created_by,updated_at,updated_by) VALUES
('f8c9ffd6-a234-4723-bd2a-68379df33ff0' /*not nullable*/,'s' /*not nullable*/,'s','s','2016-12-16 16:38:16.593',
's' /*not nullable*/,'2016-12-16 16:38:16.593','s' /*not nullable*/);


INSERT INTO "kms"."organization" (id,name,tenant_id,customer_id,description,created_at,created_by,updated_at,updated_by) VALUES
('f8c9ffd6-a234-4723-bd2a-68379df33ff1' /*not nullable*/,'5cfc6415' /*not nullable*/,'5cfc6415' /*not nullable*/,'
f8c9ffd6-a234-4723-bd2a-68379df33ff0' /*not nullable*/,'5cfc6415',
'2016-12-16 16:36:49.407','s','2016-12-16 16:36:49.407','s');

INSERT INTO "kms"."developer"
(id,tenant_id,
username,first_name,last_name,password,email,
status,encrypted_password,salt,
created_at,created_by,updated_at,updated_by,_change_selector) VALUES
('f8c9ffd6-a234-4723-bd2a-68379df33fff','5cfc6415','Dev1','Alex' /*not nullable*/,'Khimich' /*not nullable*/,
'Passwd','alexkhimich+edgex@google.com' /*not nullable*/,
'ACTIVE','s','s',now(),'sql',now(),'sql','5cfc6415');

INSERT INTO "kms"."api_product"
(id,tenant_id,name,display_name,description,
api_resources,approval_type,scopes,proxies,environments,
quota,quota_time_unit,quota_interval,
created_at,created_by,updated_at,updated_by,_change_selector) VALUES
('f8c9ffd6-a234-4723-bd2a-68379df33ff2' /*not nullable*/,'5cfc6415' /*not nullable*/,'Product1' /*not nullable*/,'Product1','s',
'{/**}','AUTO','{""}','{iloveapis}','{prod,test}',
's','s',0,
'2016-12-16 16:39:11.596','s' ,'2016-12-16 16:39:11.596','s' /*not nullable*/,'5cfc6415' /*not nullable*/);

INSERT INTO "kms"."app"
(id,tenant_id,name,display_name,
access_type,callback_url,status,app_family,
company_id,developer_id,parent_id,
type,
created_at,created_by,updated_at,updated_by,_change_selector) VALUES
('f8c9ffd6-a234-4723-bd2a-68379df33ff3' /*not nullable*/,'5cfc6415' /*not nullable*/,'App1' /*not nullable*/,'App1',
's','s','APPROVED','default',
NULL,'f8c9ffd6-a234-4723-bd2a-68379df33fff','f8c9ffd6-a234-4723-bd2a-68379df33fff' /*not nullable*/,
'DEVELOPER',
'2016-12-16 16:43:56.839','s' ,'2016-12-16 16:43:56.839','s' /*not nullable*/,'5cfc6415' /*not nullable*/);

INSERT INTO "kms"."app_credential"
(id,tenant_id,consumer_secret,app_id,
method_type,status,
issued_at,expires_at,app_status,scopes,
created_at,created_by,updated_at,updated_by,_change_selector) VALUES
('E1HbT0ecxK1M7CMTsJ4TvAzWLTrF3zB9', '5cfc6415', 'DER234','f8c9ffd6-a234-4723-bd2a-68379df33ff3',
NULL,'APPROVED','2016-12-16 16:46:50.590',NULL,NULL,'{}','2016-12-16 16:46:50.590','s' ,'2016-12-16 16:46:50.590','s' ,'5cfc6415' );



INSERT INTO "kms"."app_credential_apiproduct_mapper"
(tenant_id,appcred_id,app_id,apiprdt_id,
status,_change_selector) VALUES
('5cfc6415','E1HbT0ecxK1M7CMTsJ4TvAzWLTrF3zB9','f8c9ffd6-a234-4723-bd2a-68379df33ff3','f8c9ffd6-a234-4723-bd2a-68379df33ff2'
,'APPROVED','5cfc6415' );

curl -s -d "key=E1HbT0ecxK1M7CMTsJ4TvAzWLTrF3zB9&scopeuuid=9356ba25-38df-4e72-bb4d-f974ce9683d6&uriPath=/&action=verify"
"http://localhost:9090/verifiers/apikey" | python -m json.tool
