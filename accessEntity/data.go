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
	"fmt"
	"github.com/apid/apidApiMetadata/common"
	"strings"
)

const (
	sql_select_api_product = `SELECT * FROM kms_api_product AS ap `
	sql_select_tenant_org  = ` (SELECT o.tenant_id FROM kms_organization AS o WHERE o.name=?)`
)

type DbManager struct {
	common.DbManager
}

func (d *DbManager) GetApiProductNames(id string, idType string) ([]string, error) {
	var query string
	switch idType {
	case TypeConsumerKey:
		query = selectApiProductsById(
			selectAppCredentialMapperByConsumerKey(
				"?",
				"apiprdt_id",
			),
			"name",
		)
	case TypeApp:
		query = selectApiProductsById(
			selectAppCredentialMapperByAppId(
				"?",
				"apiprdt_id",
			),
			"name",
		)
	default:
		return nil, fmt.Errorf("unsupported idType")
	}

	rows, err := d.GetDb().Query(query, id)
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
		"?",
		"name",
	)
	name := sql.NullString{}
	err := d.GetDb().QueryRow(query, comId).Scan(&name)
	if err != nil || !name.Valid {
		return "", err
	}
	return name.String, nil
}

func (d *DbManager) GetDevEmailByDevId(devId string) (string, error) {
	query := selectDeveloperById(
		"?",
		"email",
	)
	email := sql.NullString{}
	err := d.GetDb().QueryRow(query, devId).Scan(&email)
	if err != nil || !email.Valid {
		return "", err
	}
	return email.String, nil
}

func (d *DbManager) GetComNames(id string, idType string) ([]string, error) {
	var query string
	switch idType {
	case TypeDeveloper:
		query = selectCompanyByComId(
			selectCompanyDeveloperByDevId(
				"?",
				"company_id",
			),
			"name",
		)
	case TypeCompany:
		query = selectCompanyByComId(
			"?",
			"name",
		)
	default:
		return nil, fmt.Errorf("unsupported idType")
	}

	rows, err := d.GetDb().Query(query, id)
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

func (d *DbManager) GetAppNames(id string, t string) ([]string, error) {
	var query string
	switch t {
	case TypeDeveloper:
		query = selectAppByDevId(
			"?",
			"name",
		)
	case TypeCompany:
		query = selectAppByComId(
			"?",
			"name",
		)
	default:
		return nil, fmt.Errorf("app type not supported")
	}
	rows, err := d.GetDb().Query(query, id)
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
			"?",
			"status",
		)
	case AppTypeCompany:
		query = selectCompanyByComId(
			"?",
			"status",
		)
	}
	status := sql.NullString{}
	err := d.GetDb().QueryRow(query, id).Scan(&status)
	if err != nil || !status.Valid {
		return "", err
	}

	return status.String, nil
}

func (d *DbManager) GetApiProducts(org, priKey, priVal, secKey, secVal string) (apiProducts []common.ApiProduct, err error) {
	if priKey == IdentifierAppId {
		apiProducts, err = d.getApiProductsByAppId(priVal, org)
		if err != nil {
			return
		}
	} else if priKey == IdentifierApiProductName {
		apiProducts, err = d.getApiProductsByName(priVal, org)
		if err != nil {
			return
		}
	} else if priKey == IdentifierAppName {
		switch secKey {
		case IdentifierDeveloperEmail:
			apiProducts, err = d.getApiProductsByAppName(priVal, secVal, "", "", org)
		case IdentifierDeveloperId:
			apiProducts, err = d.getApiProductsByAppName(priVal, "", secVal, "", org)
		case IdentifierCompanyName:
			apiProducts, err = d.getApiProductsByAppName(priVal, "", "", secVal, org)
		case IdentifierApiResource:
			fallthrough
		case "":
			apiProducts, err = d.getApiProductsByAppName(priVal, "", "", "", org)
		}
		if err != nil {
			return
		}
	} else if priKey == IdentifierConsumerKey {
		apiProducts, err = d.getApiProductsByConsumerKey(priVal, org)
		if err != nil {
			return
		}
	}

	if secKey == IdentifierApiResource {
		apiProducts = filterApiProductsByResource(apiProducts, secVal)
	}
	return
}

func (d *DbManager) GetApps(org, priKey, priVal, secKey, secVal string) (apps []common.App, err error) {
	switch priKey {
	case IdentifierAppId:
		return d.getAppByAppId(priVal, org)
	case IdentifierAppName:
		switch secKey {
		case IdentifierDeveloperEmail:
			return d.getAppByAppName(priVal, secVal, "", "", org)
		case IdentifierDeveloperId:
			return d.getAppByAppName(priVal, "", secVal, "", org)
		case IdentifierCompanyName:
			return d.getAppByAppName(priVal, "", "", secVal, org)
		case "":
			return d.getAppByAppName(priVal, "", "", "", org)
		}
	case IdentifierConsumerKey:
		return d.getAppByConsumerKey(priVal, org)
	}
	return
}

func (d *DbManager) GetCompanies(org, priKey, priVal, secKey, secVal string) (companies []common.Company, err error) {
	switch priKey {
	case IdentifierAppId:
		return d.getCompanyByAppId(priVal, org)
	case IdentifierCompanyName:
		return d.getCompanyByName(priVal, org)
	case IdentifierConsumerKey:
		return d.getCompanyByConsumerKey(priVal, org)
	}
	return
}

func (d *DbManager) GetCompanyDevelopers(org, priKey, priVal, secKey, secVal string) (companyDevelopers []common.CompanyDeveloper, err error) {
	if priKey == IdentifierCompanyName {
		return d.getCompanyDeveloperByComName(priVal, org)
	}
	return
}

func (d *DbManager) GetAppCredentials(org, priKey, priVal, secKey, secVal string) (appCredentials []common.AppCredential, err error) {

	switch priKey {
	case IdentifierConsumerKey:
		return d.getAppCredentialByConsumerKey(priVal, org)
	case IdentifierAppId:
		return d.getAppCredentialByAppId(priVal, org)
	}
	return
}

func (d *DbManager) GetDevelopers(org, priKey, priVal, secKey, secVal string) (developers []common.Developer, err error) {
	switch priKey {
	case IdentifierAppId:
		return d.getDeveloperByAppId(priVal, org)
	case IdentifierDeveloperEmail:
		return d.getDeveloperByEmail(priVal, org)
	case IdentifierConsumerKey:
		return d.getDeveloperByConsumerKey(priVal, org)
	case IdentifierDeveloperId:
		return d.getDeveloperById(priVal, org)
	}
	return
}

func (d *DbManager) getApiProductsByName(apiProdName string, org string) (apiProducts []common.ApiProduct, err error) {
	err = d.GetDb().QueryStructs(&apiProducts,
		sql_select_api_product+
			`WHERE ap.name = ? AND ap.tenant_id IN `+
			sql_select_tenant_org,
		apiProdName,
		org,
	)
	return
}

func (d *DbManager) getApiProductsByAppId(appId string, org string) (apiProducts []common.ApiProduct, err error) {
	cols := []string{"*"}
	query := selectApiProductsById(
		selectAppCredentialMapperByAppId(
			"?",
			"apiprdt_id",
		),
		cols...,
	) + " AND ap.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getApiProductsByAppId: %v", query)
	err = d.GetDb().QueryStructs(&apiProducts, query, appId, org)
	return
}

func (d *DbManager) getApiProductsByConsumerKey(consumerKey string, org string) (apiProducts []common.ApiProduct, err error) {
	cols := []string{"*"}
	query := selectApiProductsById(
		selectAppCredentialMapperByConsumerKey(
			"?",
			"apiprdt_id",
		),
		cols...,
	) + " AND ap.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getApiProductsByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&apiProducts, query, consumerKey, org)
	return
}

func (d *DbManager) getApiProductsByAppName(appName, devEmail, devId, comName, org string) (apiProducts []common.ApiProduct, err error) {
	cols := []string{"*"}
	var appQuery string
	args := []interface{}{appName}
	switch {
	case devEmail != "":
		appQuery = selectAppByNameAndDeveloperId(
			"?",
			selectDeveloperByEmail(
				"?",
				"id",
			),
			"id",
		)
		args = append(args, devEmail)
	case devId != "":
		appQuery = selectAppByNameAndDeveloperId(
			"?",
			"?",
			"id",
		)
		args = append(args, devId)
	case comName != "":
		appQuery = selectAppByNameAndCompanyId(
			"?",
			selectCompanyByName(
				"?",
				"id",
			),
			"id",
		)
		args = append(args, comName)
	default:
		appQuery = selectAppByName(
			"?",
			"id",
		)
	}

	args = append(args, org)

	query := selectApiProductsById(
		selectAppCredentialMapperByAppId(
			appQuery,
			"apiprdt_id",
		),
		cols...,
	) + " AND ap.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getApiProductsByAppName: %v", query)
	err = d.GetDb().QueryStructs(&apiProducts, query, args...)
	return
}

func (d *DbManager) getAppByAppId(id, org string) (apps []common.App, err error) {
	cols := []string{"*"}
	query := selectAppById(
		"?",
		cols...,
	) + " AND a.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getAppByAppId: %v \n %v", query, id)
	err = d.GetDb().QueryStructs(&apps, query, id, org)
	return
}

func (d *DbManager) getAppByAppName(appName, devEmail, devId, comName, org string) (apps []common.App, err error) {
	cols := []string{"*"}
	var query string
	args := []interface{}{appName}
	switch {
	case devEmail != "":
		query = selectAppByNameAndDeveloperId(
			"?",
			selectDeveloperByEmail(
				"?",
				"id",
			),
			cols...,
		)
		args = append(args, devEmail)
	case devId != "":
		query = selectAppByNameAndDeveloperId(
			"?",
			"?",
			cols...,
		)
		args = append(args, devId)
	case comName != "":
		query = selectAppByNameAndCompanyId(
			"?",
			selectCompanyByName(
				"?",
				"id",
			),
			cols...,
		)
		args = append(args, comName)
	default:
		query = selectAppByName(
			"?",
			cols...,
		)
	}

	args = append(args, org)

	query += " AND a.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getAppByAppName: %v", query)
	err = d.GetDb().QueryStructs(&apps, query, args...)
	return
}

func (d *DbManager) getAppByConsumerKey(consumerKey, org string) (apps []common.App, err error) {
	cols := []string{"*"}
	query := selectAppById(
		selectAppCredentialMapperByConsumerKey(
			"?",
			"app_id",
		),
		cols...,
	) + " AND a.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getAppByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&apps, query, consumerKey, org)
	return
}

func (d *DbManager) getAppCredentialByConsumerKey(consumerKey, org string) (appCredentials []common.AppCredential, err error) {
	cols := []string{"*"}
	query := selectAppCredentialByConsumerKey(
		"?",
		cols...,
	) + " AND ac.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getAppCredentialByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&appCredentials, query, consumerKey, org)
	return
}

func (d *DbManager) getAppCredentialByAppId(appId, org string) (appCredentials []common.AppCredential, err error) {
	cols := []string{"*"}
	query := selectAppCredentialByConsumerKey(
		selectAppCredentialMapperByAppId(
			"?",
			"appcred_id",
		),
		cols...,
	) + " AND ac.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getAppCredentialByAppId: %v", query)
	err = d.GetDb().QueryStructs(&appCredentials, query, appId, org)
	return
}

func (d *DbManager) getCompanyByAppId(appId, org string) (companies []common.Company, err error) {
	cols := []string{"*"}
	query := selectCompanyByComId(
		selectAppById(
			"?",
			"company_id",
		),
		cols...,
	) + " AND com.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getCompanyByAppId: %v", query)
	err = d.GetDb().QueryStructs(&companies, query, appId, org)
	return
}

func (d *DbManager) getCompanyByName(name, org string) (companies []common.Company, err error) {
	cols := []string{"*"}
	query := selectCompanyByName(
		"?",
		cols...,
	) + " AND com.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getCompanyByName: %v", query)
	err = d.GetDb().QueryStructs(&companies, query, name, org)
	return
}

func (d *DbManager) getCompanyByConsumerKey(consumerKey, org string) (companies []common.Company, err error) {
	cols := []string{"*"}
	query := selectCompanyByComId(
		selectAppById(
			selectAppCredentialMapperByConsumerKey(
				"?",
				"app_id",
			),
			"company_id",
		),
		cols...,
	) + " AND com.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getCompanyByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&companies, query, consumerKey, org)
	return
}

func (d *DbManager) getCompanyDeveloperByComName(comName, org string) (companyDevelopers []common.CompanyDeveloper, err error) {
	cols := []string{"*"}
	query := selectCompanyDeveloperByComId(
		selectCompanyByName(
			"?",
			"id",
		),
		cols...,
	) + " AND cd.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getCompanyDeveloperByComName: %v", query)
	err = d.GetDb().QueryStructs(&companyDevelopers, query, comName, org)
	return
}

func (d *DbManager) getDeveloperByAppId(appId, org string) (developers []common.Developer, err error) {
	cols := []string{"*"}
	query := selectDeveloperById(
		selectAppById(
			"?",
			"developer_id",
		),
		cols...,
	) + " AND dev.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getDeveloperByAppId: %v", query)
	err = d.GetDb().QueryStructs(&developers, query, appId, org)
	return
}

func (d *DbManager) getDeveloperByConsumerKey(consumerKey, org string) (developers []common.Developer, err error) {
	cols := []string{"*"}
	query := selectDeveloperById(
		selectAppById(
			selectAppCredentialMapperByConsumerKey(
				"?",
				"app_id",
			),
			"developer_id",
		),
		cols...,
	) + " AND dev.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getDeveloperByConsumerKey: %v", query)
	err = d.GetDb().QueryStructs(&developers, query, consumerKey, org)
	return
}

func (d *DbManager) getDeveloperByEmail(email, org string) (developers []common.Developer, err error) {
	cols := []string{"*"}
	query := selectDeveloperByEmail(
		"?",
		cols...,
	) + " AND dev.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getDeveloperByEmail: %v", query)
	err = d.GetDb().QueryStructs(&developers, query, email, org)
	return
}

func (d *DbManager) getDeveloperById(id, org string) (developers []common.Developer, err error) {
	cols := []string{"*"}
	query := selectDeveloperById(
		"?",
		cols...,
	) + " AND dev.tenant_id IN " + sql_select_tenant_org
	//log.Debugf("getDeveloperById: %v", query)
	err = d.GetDb().QueryStructs(&developers, query, id, org)
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
	//log.Debugf("Before filter: %v", apiProducts)
	var prods []common.ApiProduct
	for _, prod := range apiProducts {
		resources := common.JsonToStringArray(prod.ApiResources)
		if Contains(resources, resource) {
			prods = append(prods, prod)
		}
	}
	//log.Debugf("After filter: %v", prods)
	return prods
}

func Contains(sl []string, str string) bool {
	for _, s := range sl {
		if s == str {
			return true
		}
	}
	return false
}
