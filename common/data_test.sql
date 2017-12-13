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
COMMIT;
