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
	"database/sql"
	"github.com/apid/apidVerifyApiKey/common"
	"strings"
)

const (
	sql_select_api_product = `SELECT * FROM kms_api_product AS ap `
	sql_select_org         = `SELECT * FROM kms_organization AS o WHERE o.tenant_id=$1 LIMIT 1;`
)

type DbManager struct {
	common.DbManager
}

func (d *DbManager) GetOrgName(tenantId string) (string, error) {
	row := d.GetDb().QueryRow(sql_select_org, tenantId)
	org := sql.NullString{}
	if err := row.Scan(&org); err != nil {
		return "", err
	}
	if org.Valid {
		return org.String, nil
	}
	return "", nil
}

func (d *DbManager) GetApiProductNamesByAppId(appId string) ([]string, error) {
	query := selectApiProductsById(
		selectAppCredentialMapperByAppId(
			"'"+appId+"'",
			"apiprdt_id",
		),
		"name",
	)
	rows, err := d.GetDb().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		name := sql.NullString{}
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		if name.Valid {
			names = append(names, name.String)
		}
	}
	return names, nil
}

func (d *DbManager) GetAppNamesByComId(comId string) ([]string, error) {
	query := selectAppByComId(
		"'"+comId+"'",
		"name",
	)
	rows, err := d.GetDb().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		name := sql.NullString{}
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		if name.Valid {
			names = append(names, name.String)
		}
	}
	return names, nil
}

func (d *DbManager) GetComNameByComId(comId string) (string, error) {
	query := selectCompanyByComId(
		"'"+comId+"'",
		"name",
	)
	name := sql.NullString{}
	err := d.GetDb().QueryRow(query).Scan(&name)
	if err != nil || !name.Valid {
		return "", err
	}
	return name.String, nil
}

func (d *DbManager) GetDevEmailByDevId(devId string) (string, error) {
	query := selectDeveloperById(
		"'"+devId+"'",
		"email",
	)
	email := sql.NullString{}
	err := d.GetDb().QueryRow(query).Scan(&email)
	if err != nil || !email.Valid {
		return "", err
	}
	return email.String, nil
}

func (d *DbManager) GetComNamesByDevId(devId string) ([]string, error) {
	query := selectCompanyByComId(
		selectCompanyDeveloperByDevId(
			"'"+devId+"'",
			"company_id",
		),
		"name",
	)
	rows, err := d.GetDb().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		name := sql.NullString{}
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		if name.Valid {
			names = append(names, name.String)
		}
	}
	return names, nil
}

func (d *DbManager) GetAppNamesByDevId(devId string) ([]string, error) {
	query := selectAppByDevId(
		"'"+devId+"'",
		"name",
	)
	rows, err := d.GetDb().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		name := sql.NullString{}
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		if name.Valid {
			names = append(names, name.String)
		}
	}
	return names, nil
}

func (d *DbManager) GetStatus(id, t string) (string, error) {
	var query string
	switch t {
	case AppTypeDeveloper:
		query = selectDeveloperById(
			"'"+id+"'",
			"status",
		)
	case AppTypeCompany:
		query = selectDeveloperById(
			"'"+id+"'",
			"status",
		)
	}
	status := sql.NullString{}
	err := d.GetDb().QueryRow(query).Scan(&status)
	if err != nil || !status.Valid {
		return "", err
	}

	return status.String, nil
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

func (d *DbManager) GetApps(priKey, priVal, secKey, secVal string) (apps []common.App, err error) {
	switch priKey {
	case IdentifierAppId:
		return d.getAppByAppId(priVal)
	case IdentifierAppName:
		switch secKey {
		case IdentifierDeveloperEmail:
			return d.getAppByAppName(priVal, secVal, "", "")
		case IdentifierDeveloperId:
			return d.getAppByAppName(priVal, "", secVal, "")
		case IdentifierCompanyName:
			return d.getAppByAppName(priVal, "", "", secVal)
		case "":
			return d.getAppByAppName(priVal, "", "", "")
		}
	case IdentifierConsumerKey:
		return d.getAppByConsumerKey(priVal)
	}
	return
}

func (d *DbManager) GetCompanies(priKey, priVal, secKey, secVal string) (companies []common.Company, err error) {
	switch priKey {
	case IdentifierAppId:
		return d.getCompanyByAppId(priVal)
	case IdentifierCompanyName:
		return d.getCompanyByName(priVal)
	case IdentifierConsumerKey:
		return d.getCompanyByConsumerKey(priVal)
	}
	return
}

func (d *DbManager) GetCompanyDevelopers(priKey, priVal, secKey, secVal string) (companyDevelopers []common.CompanyDeveloper, err error) {
	if priKey == IdentifierCompanyName {
		return d.getCompanyDeveloperByComName(priVal)
	}
	return
}

func (d *DbManager) GetAppCredentials(priKey, priVal, secKey, secVal string) (appCredentials []common.AppCredential, err error) {
	if priKey == IdentifierConsumerKey {
		return d.getAppCredentialByConsumerKey(priVal)
	}
	return
}

func (d *DbManager) GetDevelopers(priKey, priVal, secKey, secVal string) (developers []common.Developer, err error) {
	switch priKey {
	case IdentifierAppId:
		return d.getDeveloperByAppId(priVal)
	case IdentifierDeveloperEmail:
		return d.getDeveloperByEmail(priVal)
	case IdentifierConsumerKey:
		return d.getDeveloperByConsumerKey(priVal)
	case IdentifierDeveloperId:
		return d.getDeveloperById(priVal)
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
	query := selectApiProductsById(
		selectAppCredentialMapperByAppId(
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
	query := selectApiProductsById(
		selectAppCredentialMapperByConsumerKey(
			"'"+consumerKey+"'",
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

	query := selectApiProductsById(
		selectAppCredentialMapperByAppId(
			appQuery,
			"apiprdt_id",
		),
		cols...,
	)
	log.Debugf("getApiProductsByAppName: %v", query)
	err = d.GetDb().QueryStructs(&apiProducts, query)
	return
}

func (d *DbManager) getAppByAppId(id string) (apps []common.App, err error) {
	cols := []string{"*"}
	query := selectAppById(
		"'"+id+"'",
		cols...,
	)
	log.Debugf("getAppByAppId: %v", query)
	err = d.GetDb().QueryStructs(&apps, query)
	return
}

func (d *DbManager) getAppByAppName(appName, devEmail, devId, comName string) (apps []common.App, err error) {
	cols := []string{"*"}
	var query string
	switch {
	case devEmail != "":
		query = selectAppByNameAndDeveloperId(
			"'"+appName+"'",
			selectDeveloperByEmail(
				"'"+devEmail+"'",
				"id",
			),
			cols...,
		)
	case devId != "":
		query = selectAppByNameAndDeveloperId(
			"'"+appName+"'",
			"'"+devId+"'",
			cols...,
		)
	case comName != "":
		query = selectAppByNameAndCompanyId(
			"'"+appName+"'",
			selectCompanyByName(
				"'"+comName+"'",
				"id",
			),
			cols...,
		)
	default:
		query = selectAppByName(
			"'"+appName+"'",
			cols...,
		)
	}
	log.Debugf("getAppByAppName: %v", query)
	err = d.GetDb().QueryStructs(&apps, query)
	return
}

func (d *DbManager) getAppByConsumerKey(consumerKey string) (apps []common.App, err error) {
	cols := []string{"*"}
	query := selectAppById(
		selectAppCredentialMapperByConsumerKey(
			"'"+consumerKey+"'",
			"app_id",
		),
		cols...,
	)
	log.Debugf("getAppByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&apps, query)
	return
}

func (d *DbManager) getAppCredentialByConsumerKey(consumerKey string) (appCredentials []common.AppCredential, err error) {
	cols := []string{"*"}
	query := selectAppCredentialByConsumerKey(
		"'"+consumerKey+"'",
		cols...,
	)
	log.Debugf("getAppCredentialByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&appCredentials, query)
	return
}

func (d *DbManager) getCompanyByAppId(appId string) (companies []common.Company, err error) {
	cols := []string{"*"}
	query := selectCompanyByComId(
		selectAppById(
			"'"+appId+"'",
			"company_id",
		),
		cols...,
	)
	log.Debugf("getCompanyByAppId: %v", query)
	err = d.GetDb().QueryStructs(&companies, query)
	return
}

func (d *DbManager) getCompanyByName(name string) (companies []common.Company, err error) {
	cols := []string{"*"}
	query := selectCompanyByName(
		"'"+name+"'",
		cols...,
	)
	log.Debugf("getCompanyByName: %v", query)
	err = d.GetDb().QueryStructs(&companies, query)
	return
}

func (d *DbManager) getCompanyByConsumerKey(consumerKey string) (companies []common.Company, err error) {
	cols := []string{"*"}
	query := selectCompanyByComId(
		selectAppById(
			selectAppCredentialMapperByConsumerKey(
				"'"+consumerKey+"'",
				"app_id",
			),
			"company_id",
		),
		cols...,
	)
	log.Debugf("getCompanyByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&companies, query)
	return
}

func (d *DbManager) getCompanyDeveloperByComName(comName string) (companyDevelopers []common.CompanyDeveloper, err error) {
	cols := []string{"*"}
	query := selectCompanyDeveloperByComId(
		selectCompanyByName(
			"'"+comName+"'",
			"id",
		),
		cols...,
	)
	log.Debugf("getCompanyDeveloperByComName: %v", query)
	err = d.GetDb().QueryStructs(&companyDevelopers, query)
	return
}

func (d *DbManager) getDeveloperByAppId(appId string) (developers []common.Developer, err error) {
	cols := []string{"*"}
	query := selectDeveloperById(
		selectAppById(
			"'"+appId+"'",
			"developer_id",
		),
		cols...,
	)
	log.Debugf("getDeveloperByAppId: %v", query)
	err = d.GetDb().QueryStructs(&developers, query)
	return
}

func (d *DbManager) getDeveloperByConsumerKey(consumerKey string) (developers []common.Developer, err error) {
	cols := []string{"*"}
	query := selectDeveloperById(
		selectAppById(
			selectAppCredentialMapperByConsumerKey(
				"'"+consumerKey+"'",
				"app_id",
			),
			"developer_id",
		),
		cols...,
	)
	log.Debugf("getDeveloperByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&developers, query)
	return
}

func (d *DbManager) getDeveloperByEmail(email string) (developers []common.Developer, err error) {
	cols := []string{"*"}
	query := selectDeveloperByEmail(
		"'"+email+"'",
		cols...,
	)
	log.Debugf("getDeveloperByEmail: %v", query)
	err = d.GetDb().QueryStructs(&developers, query)
	return
}

func (d *DbManager) getDeveloperById(id string) (developers []common.Developer, err error) {
	cols := []string{"*"}
	query := selectDeveloperById(
		"'"+id+"'",
		cols...,
	)
	log.Debugf("getDeveloperById: %v", query)
	err = d.GetDb().QueryStructs(&developers, query)
	return
}

func selectApiProductsById(idQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_api_product AS ap WHERE ap.id IN (" +
		idQuery +
		")"

	return query
}

func selectAppCredentialMapperByAppId(idQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app_credential_apiproduct_mapper AS acm WHERE acm.app_id IN (" +
		idQuery +
		")"
	return query
}

func selectAppCredentialMapperByConsumerKey(keyQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app_credential_apiproduct_mapper AS acm WHERE acm.appcred_id IN (" +
		keyQuery +
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

func selectAppById(appIdQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app AS a WHERE a.id IN (" +
		appIdQuery +
		")"
	return query
}

func selectAppByComId(comIdQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app AS a WHERE a.company_id IN (" +
		comIdQuery +
		")"
	return query
}

func selectAppByDevId(devIdQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app AS a WHERE a.developer_id IN (" +
		devIdQuery +
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

func selectDeveloperById(idQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_developer AS dev WHERE dev.id IN (" +
		idQuery +
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

func selectCompanyByComId(comIdQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_company AS com WHERE com.id IN (" +
		comIdQuery +
		")"
	return query
}

func selectCompanyDeveloperByComId(comIdQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_company_developer AS cd WHERE cd.company_id IN (" +
		comIdQuery +
		")"
	return query
}

func selectCompanyDeveloperByDevId(devIdQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_company_developer AS cd WHERE cd.developer_id IN (" +
		devIdQuery +
		")"
	return query
}

func selectAppCredentialByConsumerKey(consumerQuery string, colNames ...string) string {
	query := "SELECT " +
		strings.Join(colNames, ",") +
		" FROM kms_app_credential AS ac WHERE ac.id IN (" +
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
