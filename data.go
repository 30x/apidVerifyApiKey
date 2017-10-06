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
	"github.com/apid/apid-core"
	"strings"
	"sync"
)

type dbManager struct {
	data      apid.DataService
	db        apid.DB
	dbMux     sync.RWMutex
	dbVersion string
}

func (dbc *dbManager) setDbVersion(version string) {
	db, err := dbc.data.DBVersion(version)
	if err != nil {
		log.Panicf("Unable to access database: %v", err)
	}
	dbc.dbMux.Lock()
	dbc.db = db
	dbc.dbMux.Unlock()
	dbc.dbVersion = version
	// TODO : check if we need to release old db here...
}

func (dbc *dbManager) getDb() apid.DB {
	dbc.dbMux.RLock()
	defer dbc.dbMux.RUnlock()
	return dbc.db
}

func (dbc *dbManager) getDbVersion() string {
	return dbc.dbVersion
}

type dbManagerInterface interface {
	setDbVersion(string)
	getDb() apid.DB
	getDbVersion() string
	getKmsAttributes(tenantId string, entities ...string) map[string][]Attribute
	getApiKeyDetails(dataWrapper *VerifyApiKeyRequestResponseDataWrapper) error
}

func (dbc *dbManager) getKmsAttributes(tenantId string, entities ...string) map[string][]Attribute {

	db := dbc.db
	var attName, attValue, entity_id sql.NullString
	// TODO : is there no other better way to do in caluse???
	sql := sql_GET_KMS_ATTRIBUTES_FOR_TENANT + ` and entity_id in ('` + strings.Join(entities, `','`) + `')`
	mapOfAttributes := make(map[string][]Attribute)
	attributes, err := db.Query(sql, tenantId)
	defer attributes.Close()
	if err != nil {
		log.Error("Error while fetching attributes for tenant id : %s and entityId : %s", tenantId, err)
		return mapOfAttributes
	}
	for attributes.Next() {
		err := attributes.Scan(
			&entity_id,
			&attName,
			&attValue,
		)
		if err != nil {
			log.Error("error fetching attributes for entityid ", entities, err)
			return nil
		}
		if attName.Valid && entity_id.Valid {
			att := Attribute{Name: attName.String, Value: attValue.String}
			mapOfAttributes[entity_id.String] = append(mapOfAttributes[entity_id.String], att)
		} else {
			log.Debugf("Not valid. AttName: %s Entity_id: %s", attName.String, entity_id.String)
		}
	}
	log.Debug("attributes returned for query ", sql, " are ", mapOfAttributes)
	return mapOfAttributes
}

func (dbc dbManager) getApiKeyDetails(dataWrapper *VerifyApiKeyRequestResponseDataWrapper) error {

	db := dbc.db

	err := db.QueryRow(sql_GET_API_KEY_DETAILS_SQL, dataWrapper.verifyApiKeyRequest.Key, dataWrapper.verifyApiKeyRequest.OrganizationName).
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
		)

	if err != nil {
		log.Debug("error fetching verify apikey details ", err)
		return errors.New("InvalidApiKey")
	}

	if dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl != "" {
		dataWrapper.verifyApiKeySuccessResponse.ClientId.RedirectURIs = []string{dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl}
	}

	dataWrapper.apiProducts = dbc.getApiProductsForApiKey(dataWrapper.verifyApiKeyRequest.Key, dataWrapper.tenant_id)

	log.Debug("dataWrapper : ", dataWrapper)

	return err
}

func (dbc dbManager) getApiProductsForApiKey(key, tenantId string) []ApiProductDetails {

	db := dbc.db
	allProducts := []ApiProductDetails{}
	var proxies, environments, resources string

	rows, err := db.Query(sql_GET_API_PRODUCTS_FOR_KEY_SQL, key, tenantId)
	defer rows.Close()
	if err != nil {
		log.Error("error fetching apiProduct details", err)
		return allProducts
	}

	for rows.Next() {
		apiProductDetais := ApiProductDetails{}
		rows.Scan(
			&apiProductDetais.Id,
			&apiProductDetais.Name,
			&apiProductDetais.DisplayName,
			&apiProductDetais.QuotaLimit,
			&apiProductDetais.QuotaInterval,
			&apiProductDetais.QuotaTimeunit,
			&apiProductDetais.CreatedAt,
			&apiProductDetais.CreatedBy,
			&apiProductDetais.LastmodifiedAt,
			&apiProductDetais.LastmodifiedBy,
			&proxies,
			&environments,
			&resources,
		)
		apiProductDetais.Apiproxies = jsonToStringArray(proxies)
		apiProductDetais.Environments = jsonToStringArray(environments)
		apiProductDetais.Resources = jsonToStringArray(resources)

		allProducts = append(allProducts, apiProductDetais)
	}

	log.Debug("Api products retrieved for key : [%s] , tenantId : [%s] is ", key, tenantId, allProducts)

	return allProducts
}
