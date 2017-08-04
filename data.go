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
	"errors"
	"github.com/30x/apid-core"
	"sync"
)

type dbManager struct {
	data  apid.DataService
	db    apid.DB
	dbMux sync.RWMutex
}

func (dbc *dbManager) setDbVersion(version string) {
	db, err := dbc.data.DBVersion(version)
	if err != nil {
		log.Panicf("Unable to access database: %v", err)
	}
	dbc.dbMux.Lock()
	dbc.db = db
	dbc.dbMux.Unlock()
}

func (dbc *dbManager) getDb() apid.DB {
	dbc.dbMux.RLock()
	defer dbc.dbMux.RUnlock()
	return dbc.db
}

func (dbc *dbManager) initDb() error {
	db := dbc.getDb()
	if db == nil {
		return errors.New("DB not initialized")
	}
	return nil
}

type dbManagerInterface interface {
	setDbVersion(string)
	initDb() error
	getDb() apid.DB
	getKmsAttributes(tenantId string, entityId string) []Attribute
	getApiKeyDetails(dataWrapper *VerifyApiKeyRequestResponseDataWrapper) error
}

func (dbc *dbManager) getKmsAttributes(tenantId string, entityId string) []Attribute {

	db := dbc.db
	var attName, attValue sql.NullString
	sql := "select name, value from kms_attributes where tenant_id = $1 and entity_id = $2"
	attributesForQuery := []Attribute{}
	attributes, err := db.Query(sql, tenantId, entityId)
	if err != nil {
		log.Error("Error while fetching attributes for tenant id : %s and entityId : %s", tenantId, entityId, err)
		return attributesForQuery
	}

	for attributes.Next() {
		err := attributes.Scan(
			&attName,
			&attValue,
		)
		if err != nil {
			log.Error("error fetching attributes for entityid ", entityId, err)
		}
		if attName.String != "" {
			att := Attribute{Name: attName.String, Value: attValue.String}
			attributesForQuery = append(attributesForQuery, att)
		}
	}
	log.Debug("attributes returned for query ", sql, " are ", attributesForQuery, tenantId, entityId)
	return attributesForQuery
}

func (dbc dbManager) getApiKeyDetails(dataWrapper *VerifyApiKeyRequestResponseDataWrapper) error {

	db := dbc.db
	var proxies, environments, resources string
	sSql := `
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
				COALESCE(a.updated_by,"") as app_updated_by,

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
				INNER JOIN KMS_APP AS a
					ON c.app_id = a.id
				INNER JOIN KMS_DEVELOPER AS ad
					ON ad.id = a.developer_id
				INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
					ON mp.appcred_id = c.id
				INNER JOIN KMS_API_PRODUCT as ap
					ON ap.id = mp.apiprdt_id
				INNER JOIN KMS_ORGANIZATION AS o
					ON o.tenant_id = c.tenant_id
			WHERE 	(mp.apiprdt_id = ap.id
				AND mp.app_id = a.id
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
				COALESCE("","") as dev_first_name,
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
				COALESCE(a.updated_by,"") as app_updated_by,

				COALESCE(ap.id,"") as prod_id,
				COALESCE(ap.name,"") as prod_name,
				COALESCE(ap.display_name,"") as prod_display_name,
				COALESCE(ap.quota,"") as prod_quota,
				COALESCE(ap.quota_interval,0) as prod_quota_interval,
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
				INNER JOIN KMS_APP AS a
					ON c.app_id = a.id
				INNER JOIN KMS_COMPANY AS ad
					ON ad.id = a.company_id
				INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
					ON mp.appcred_id = c.id
				INNER JOIN KMS_API_PRODUCT as ap
					ON ap.id = mp.apiprdt_id
				INNER JOIN KMS_ORGANIZATION AS o
					ON o.tenant_id = c.tenant_id
			WHERE   (mp.apiprdt_id = ap.id
				AND mp.app_id = a.id
				AND mp.appcred_id = c.id
				AND c.id = $1
				AND o.name = $2)
		;`

	//cid,csecret,did,dusername,dfirstname,dlastname,demail,dstatus,dcreated_at,dcreated_by,dlast_modified_at,dlast_modified_by, aid,aname,aaccesstype,acallbackurl,adisplay_name,astatus,aappfamily, acompany,acreated_at,acreated_by,alast_modified_at,alast_modified_by,pid,pname,pdisplayname,pquota_limit,pqutoainterval,pquotatimeout,pcreated_at,pcreated_by,plast_modified_at,plast_modified_by sql.NullString

	err := db.QueryRow(sSql, dataWrapper.verifyApiKeyRequest.Key, dataWrapper.verifyApiKeyRequest.OrganizationName).
		Scan(
			&dataWrapper.ctype,
			&dataWrapper.tenant_id,
			&dataWrapper.verifyApiKeySuccessResponse.ClientId.Status,
			&dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientSecret,

			&dataWrapper.tempDeveloperDetails.Id,
			&dataWrapper.tempDeveloperDetails.UserName,
			&dataWrapper.tempDeveloperDetails.FirstName,
			&dataWrapper.tempDeveloperDetails.LastName,
			&dataWrapper.tempDeveloperDetails.Email,
			&dataWrapper.tempDeveloperDetails.Status,
			&dataWrapper.tempDeveloperDetails.CreatedAt,
			&dataWrapper.tempDeveloperDetails.CreatedBy,
			&dataWrapper.tempDeveloperDetails.LastmodifiedAt,
			&dataWrapper.tempDeveloperDetails.LastmodifiedBy,

			&dataWrapper.verifyApiKeySuccessResponse.App.Id,
			&dataWrapper.verifyApiKeySuccessResponse.App.Name,
			&dataWrapper.verifyApiKeySuccessResponse.App.AccessType,
			&dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl,
			&dataWrapper.verifyApiKeySuccessResponse.App.DisplayName,
			&dataWrapper.verifyApiKeySuccessResponse.App.Status,
			&dataWrapper.verifyApiKeySuccessResponse.App.AppFamily,
			&dataWrapper.verifyApiKeySuccessResponse.App.Company,
			&dataWrapper.verifyApiKeySuccessResponse.App.CreatedAt,
			&dataWrapper.verifyApiKeySuccessResponse.App.CreatedBy,
			&dataWrapper.verifyApiKeySuccessResponse.App.LastmodifiedAt,
			&dataWrapper.verifyApiKeySuccessResponse.App.LastmodifiedBy,

			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Id,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Name,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.DisplayName,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.QuotaLimit,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.QuotaInterval,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.QuotaTimeunit,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.CreatedAt,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.CreatedBy,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.LastmodifiedAt,
			&dataWrapper.verifyApiKeySuccessResponse.ApiProduct.LastmodifiedBy,
			&proxies,
			&environments,
			&resources,
		)

	if err != nil {
		log.Error("error fetching verify apikey details", err)
		return err
	}

	dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Apiproxies = jsonToStringArray(proxies)
	dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Environments = jsonToStringArray(environments)
	dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Resources = jsonToStringArray(resources)

	if dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl != "" {
		dataWrapper.verifyApiKeySuccessResponse.ClientId.RedirectURIs = []string{dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl}
	}

	log.Debug("dataWrapper : ", dataWrapper)

	return err
}
