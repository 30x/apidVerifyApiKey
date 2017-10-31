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
package accessEntity

import (
	"github.com/apid/apidVerifyApiKey/common"
	"strings"
)

const (
	sql_select_api_product = `SELECT * FROM kms_api_product AS ap `
	sql_select_api         = `
	SELECT * FROM kms_api_product AS ap WHERE ap.id IN (
		SELECT apiprdt_id FROM kms_app_credential_apiproduct_mapper AS acm WHERE acm.app_id IN (
			SELECT a.id FROM kms_app AS a WHERE a.name IN ('apstest')
		)
	);

	`
)

type DbManager struct {
	common.DbManager
}

func (d *DbManager) GetApiProducts(priKey, priVal, secKey, secVal string) (apiProducts []common.ApiProduct, err error) {
	if priKey == IdentifierAppId {
		apiProducts, err = d.getApiProductsByAppId(priVal)
		if err != nil {
			return
		}
	} else if priKey == IdentifierApiProductName {
		apiProducts, err = d.getApiProductsByName(priVal)
		if err != nil {
			return
		}
	} else if priKey == IdentifierAppName {
		switch secKey {
		case IdentifierDeveloperEmail:
			apiProducts, err = d.getApiProductsByAppName(priVal, secVal, "", "")
		case IdentifierDeveloperId:
			apiProducts, err = d.getApiProductsByAppName(priVal, "", secVal, "")
		case IdentifierCompanyName:
			apiProducts, err = d.getApiProductsByAppName(priVal, "", "", secVal)
		case IdentifierApiResource:
			fallthrough
		case "":
			apiProducts, err = d.getApiProductsByAppName(priVal, "", "", "")
		}
		if err != nil {
			return
		}
	} else if priKey == IdentifierConsumerKey {
		apiProducts, err = d.getApiProductsByConsumerKey(priVal)
		if err != nil {
			return
		}
	}

	if secKey == IdentifierApiResource {
		apiProducts = filterApiProductsByResource(apiProducts, secVal)
	}
	return
}

func (d *DbManager) getApiProductsByName(apiProdName string) (apiProducts []common.ApiProduct, err error) {
	err = d.GetDb().QueryStructs(&apiProducts,
		sql_select_api_product+`WHERE ap.name = $1;`,
		apiProdName,
	)
	return
}

func (d *DbManager) getApiProductsByAppId(appId string) (apiProducts []common.ApiProduct, err error) {
	cols := []string{"*"}
	query := selectApiProductsByIds(
		selectAppCredentialMapperByAppIds(
			"'"+appId+"'",
			"apiprdt_id",
		),
		cols...,
	)
	log.Debugf("getApiProductsByAppId: %v", query)
	err = d.GetDb().QueryStructs(&apiProducts, query)
	return
}

func (d *DbManager) getApiProductsByConsumerKey(consumerKey string) (apiProducts []common.ApiProduct, err error) {
	cols := []string{"*"}
	query := selectApiProductsByIds(
		selectAppCredentialMapperByAppIds(
			selectAppCredentialByConsumerKey(
				"'"+consumerKey+"'",
				"app_id",
			),
			"apiprdt_id",
		),
		cols...,
	)
	log.Debugf("getApiProductsByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&apiProducts, query)
	return
}

func (d *DbManager) getApiProductsByAppName(appName, devEmail, devId, comName string) (apiProducts []common.ApiProduct, err error) {
	cols := []string{"*"}
	var appQuery string
	switch {
	case devEmail != "":
		appQuery = selectAppByNameAndDeveloperId(
			"'"+appName+"'",
			selectDeveloperByEmail(
				"'"+devEmail+"'",
				"id",
			),
			"id",
		)
	case devId != "":
		appQuery = selectAppByNameAndDeveloperId(
			"'"+appName+"'",
			"'"+devId+"'",
			"id",
		)
	case comName != "":
		appQuery = selectAppByNameAndCompanyId(
			"'"+appName+"'",
			selectCompanyByName(
				"'"+comName+"'",
				"id",
			),
			"id",
		)
	default:
		appQuery = selectAppByName(
			"'"+appName+"'",
			"id",
		)
	}

	query := selectApiProductsByIds(
		selectAppCredentialMapperByAppIds(
			appQuery,
			"apiprdt_id",
		),
		cols...,
	)
	log.Debugf("getApiProductsByAppName: %v", query)
	err = d.GetDb().QueryStructs(&apiProducts, query)
	return
}

func selectApiProductsByIds(idQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_api_product AS ap WHERE ap.id IN (" +
		idQuery +
		")"

	return query
}

func selectAppCredentialMapperByAppIds(idQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app_credential_apiproduct_mapper AS acm WHERE acm.app_id IN (" +
		idQuery +
		")"
	return query
}

func selectAppByName(nameQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app AS a WHERE a.name IN (" +
		nameQuery +
		")"
	return query
}

func selectAppByNameAndDeveloperId(nameQuery string, developerIdQuery string, colNames ...string) string {
	query := selectAppByName(nameQuery, colNames...) +
		" AND developer_id IN (" +
		developerIdQuery +
		")"
	return query
}

func selectAppByNameAndCompanyId(nameQuery string, companyIdQuery string, colNames ...string) string {
	query := selectAppByName(nameQuery, colNames...) +
		" AND company_id IN (" +
		companyIdQuery +
		")"
	return query
}

func selectDeveloperByEmail(emailQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_developer AS dev WHERE dev.email IN (" +
		emailQuery +
		")"
	return query
}

func selectCompanyByName(nameQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_company AS com WHERE com.name IN (" +
		nameQuery +
		")"
	return query
}

func selectAppCredentialByConsumerKey(consumerQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app_credential AS ac WHERE ac.consumer_secret IN (" +
		consumerQuery +
		")"
	return query
}

func filterApiProductsByResource(apiProducts []common.ApiProduct, resource string) []common.ApiProduct {
	log.Debugf("Before filter: %v", apiProducts)
	var prods []common.ApiProduct
	for _, prod := range apiProducts {
		resources := common.JsonToStringArray(prod.ApiResources)
		if InSlice(resources, resource) {
			prods = append(prods, prod)
		}
	}
	log.Debugf("After filter: %v", prods)
	return prods
}

func InSlice(sl []string, str string) bool {
	for _, s := range sl {
		if s == str {
			return true
		}
	}
	return false
}
