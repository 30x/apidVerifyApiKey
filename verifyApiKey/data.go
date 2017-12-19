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
	"errors"
	"github.com/apid/apidApiMetadata/common"
)

type DbManagerInterface interface {
	common.DbManagerInterface
	getApiKeyDetails(dataWrapper *VerifyApiKeyRequestResponseDataWrapper) error
}

type DbManager struct {
	common.DbManager
}

func (dbc *DbManager) getApiKeyDetails(dataWrapper *VerifyApiKeyRequestResponseDataWrapper) error {

	db := dbc.Db

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

	secret, err := dbc.CipherManager.TryDecryptBase64(dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientSecret,
		dataWrapper.verifyApiKeyRequest.OrganizationName)
	if err != nil {
		return err
	}
	dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientSecret = secret

	if dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl != "" {
		dataWrapper.verifyApiKeySuccessResponse.ClientId.RedirectURIs = []string{dataWrapper.verifyApiKeySuccessResponse.App.CallbackUrl}
	}

	dataWrapper.apiProducts = dbc.getApiProductsForApiKey(dataWrapper.verifyApiKeyRequest.Key, dataWrapper.tenant_id)

	log.Debug("dataWrapper : ", dataWrapper)

	return err
}

func (dbc *DbManager) getApiProductsForApiKey(key, tenantId string) []ApiProductDetails {

	db := dbc.Db
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
		apiProductDetais.Apiproxies = common.JsonToStringArray(proxies)
		apiProductDetais.Environments = common.JsonToStringArray(environments)
		apiProductDetais.Resources = common.JsonToStringArray(resources)

		allProducts = append(allProducts, apiProductDetais)
	}

	log.Debug("Api products retrieved for key : [%s] , tenantId : [%s] is ", key, tenantId, allProducts)

	return allProducts
}
