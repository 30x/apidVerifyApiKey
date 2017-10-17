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

const sql_GET_API_KEY_DETAILS_SQL = `
			SELECT
				COALESCE("developer","") as ctype,
				COALESCE(c.tenant_id,""),

				COALESCE(c.status,""),
				COALESCE(c.consumer_secret,""),

				COALESCE(ad.id,"") as dev_id,
				COALESCE(ad.username,"") as dev_username,
				COALESCE(ad.first_name,"") as dev_first_name,
				COALESCE(ad.last_name,"") as dev_last_name,
				COALESCE(ad.email,"") as dev_email,
				COALESCE(ad.status,"") as dev_status,
				COALESCE(ad.created_at,"") as dev_created_at,
				COALESCE(ad.created_by,"") as dev_created_by,
				COALESCE(ad.updated_at,"") as dev_updated_at,
				COALESCE(ad.updated_by,"") as dev_updated_by,

				COALESCE(a.id,"") as app_id,
				COALESCE(a.name,"") as app_name,
				COALESCE(a.access_type,"") as app_access_type,
				COALESCE(a.callback_url,"") as app_callback_url,
				COALESCE(a.display_name,"") as app_display_name,
				COALESCE(a.status,"") as app_status,
				COALESCE(a.app_family,"") as app_app_family,
				COALESCE(a.company_id,"") as app_company_id,
				COALESCE(a.created_at,"") as app_created_at,
				COALESCE(a.created_by,"") as app_created_by,
				COALESCE(a.updated_at,"") as app_updated_at,
				COALESCE(a.updated_by,"") as app_updated_by

			FROM
				KMS_APP_CREDENTIAL AS c
				INNER JOIN KMS_APP AS a
					ON c.app_id = a.id
				INNER JOIN KMS_DEVELOPER AS ad
					ON ad.id = a.developer_id
				INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
					ON mp.appcred_id = c.id
				INNER JOIN KMS_ORGANIZATION AS o
					ON o.tenant_id = c.tenant_id
			WHERE 	(
				mp.app_id = a.id
				AND mp.appcred_id = c.id
				AND c.id = $1
				AND o.name = $2)
			UNION ALL
			SELECT
				COALESCE("company","") as ctype,
				COALESCE(c.tenant_id,""),

				COALESCE(c.status,""),
				COALESCE(c.consumer_secret,""),

				COALESCE(ad.id,"") as dev_id,
				COALESCE(ad.display_name,"") as dev_username,
				COALESCE(ad.name,"") as dev_first_name,
				COALESCE("","") as dev_last_name,
				COALESCE("","") as dev_email,
				COALESCE(ad.status,"") as dev_status,
				COALESCE(ad.created_at,"") as dev_created_at,
				COALESCE(ad.created_by,"") as dev_created_by,
				COALESCE(ad.updated_at,"") as dev_updated_at,
				COALESCE(ad.updated_by,"") as dev_updated_by,

				COALESCE(a.id,"") as app_id,
				COALESCE(a.name,"") as app_name,
				COALESCE(a.access_type,"") as app_access_type,
				COALESCE(a.callback_url,"") as app_callback_url,
				COALESCE(a.display_name,"") as app_display_name,
				COALESCE(a.status,"") as app_status,
				COALESCE(a.app_family,"") as app_app_family,
				COALESCE(a.company_id,"") as app_company_id,
				COALESCE(a.created_at,"") as app_created_at,
				COALESCE(a.created_by,"") as app_created_by,
				COALESCE(a.updated_at,"") as app_updated_at,
				COALESCE(a.updated_by,"") as app_updated_by

			FROM
				KMS_APP_CREDENTIAL AS c
				INNER JOIN KMS_APP AS a
					ON c.app_id = a.id
				INNER JOIN KMS_COMPANY AS ad
					ON ad.id = a.company_id
				INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
					ON mp.appcred_id = c.id
				INNER JOIN KMS_ORGANIZATION AS o
					ON o.tenant_id = c.tenant_id
			WHERE   (
				mp.app_id = a.id
				AND mp.appcred_id = c.id
				AND c.id = $1
				AND o.name = $2)
		;`

const sql_GET_API_PRODUCTS_FOR_KEY_SQL = `
			SELECT
				COALESCE(ap.id,"") as prod_id,
				COALESCE(ap.name,"") as prod_name,
				COALESCE(ap.display_name,"") as prod_display_name,
				COALESCE(ap.quota,"") as prod_quota,
				COALESCE(ap.quota_interval, 0) as prod_quota_interval,
				COALESCE(ap.quota_time_unit,"") as prod_quota_time_unit,
				COALESCE(ap.created_at,"") as prod_created_at,
				COALESCE(ap.created_by,"") as prod_created_by,
				COALESCE(ap.updated_at,"") as prod_updated_at,
				COALESCE(ap.updated_by,"") as prod_updated_by,
				COALESCE(ap.proxies,"") as prod_proxies,
				COALESCE(ap.environments,"") as prod_environments,
				COALESCE(ap.api_resources,"") as prod_resources
			FROM
				KMS_APP_CREDENTIAL AS c
				INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
					ON mp.appcred_id = c.id
				INNER JOIN KMS_API_PRODUCT as ap
					ON ap.id = mp.apiprdt_id
			WHERE 	(
				mp.apiprdt_id = ap.id
				AND mp.appcred_id = c.id
				AND mp.status = 'APPROVED'
				AND c.id = $1
				AND ap.tenant_id = $2
				)
		;`

const sql_GET_KMS_ATTRIBUTES_FOR_TENANT = "select entity_id, name, value from kms_attributes where tenant_id = $1"
