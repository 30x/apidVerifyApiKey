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
	"github.com/apid/apidApiMetadata/common"
	"net/http"
)

type ApiManagerInterface interface {
	common.ApiManagerInterface
	HandleRequest(w http.ResponseWriter, r *http.Request)
}

type DbManagerInterface interface {
	common.DbManagerInterface
	GetApiProducts(org, priKey, priVal, secKey, secVal string) (apiProducts []common.ApiProduct, err error)
	GetApps(org, priKey, priVal, secKey, secVal string) (apps []common.App, err error)
	GetCompanies(org, priKey, priVal, secKey, secVal string) (companies []common.Company, err error)
	GetCompanyDevelopers(org, priKey, priVal, secKey, secVal string) (companyDevelopers []common.CompanyDeveloper, err error)
	GetAppCredentials(org, priKey, priVal, secKey, secVal string) (appCredentials []common.AppCredential, err error)
	GetDevelopers(org, priKey, priVal, secKey, secVal string) (developers []common.Developer, err error)
	// utils
	GetApiProductNames(id string, idType string) ([]string, error)
	GetAppNames(id string, idType string) ([]string, error)
	GetComNames(id string, idType string) ([]string, error)
	GetDevEmailByDevId(devId string, org string) (string, error)
	GetStatus(id, t string) (string, error)
}
