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
	"github.com/apid/apid-core/cipher"
	"github.com/apid/apidApiMetadata/common"
)

type DummyCipherMan struct {
}

func (c *DummyCipherMan) AddOrgs(orgs []string) {
}

func (d *DummyCipherMan) TryDecryptBase64(input string, org string) (string, error) {
	return input, nil
}

func (d *DummyCipherMan) EncryptBase64(input string, org string, mode cipher.Mode, padding cipher.Padding) (string, error) {
	return input, nil
}

type DummyDbMan struct {
	apiProducts       []common.ApiProduct
	apps              []common.App
	companies         []common.Company
	companyDevelopers []common.CompanyDeveloper
	appCredentials    []common.AppCredential
	developers        []common.Developer
	apiProductNames   []string
	appNames          []string
	comNames          []string
	email             string
	status            string
	attrs             map[string][]common.Attribute
	err               error
}

func (d *DummyDbMan) GetOrgs() (orgs []string, err error) {
	return
}

func (d *DummyDbMan) SetDbVersion(string) {

}
func (d *DummyDbMan) GetDbVersion() string {
	return ""
}

func (d *DummyDbMan) GetKmsAttributes(tenantId string, entities ...string) map[string][]common.Attribute {
	return d.attrs
}

func (d *DummyDbMan) GetApiProducts(org, priKey, priVal, secKey, secVal string) (apiProducts []common.ApiProduct, err error) {
	return d.apiProducts, d.err
}

func (d *DummyDbMan) GetApps(org, priKey, priVal, secKey, secVal string) (apps []common.App, err error) {
	return d.apps, d.err
}

func (d *DummyDbMan) GetCompanies(org, priKey, priVal, secKey, secVal string) (companies []common.Company, err error) {
	return d.companies, d.err
}

func (d *DummyDbMan) GetCompanyDevelopers(org, priKey, priVal, secKey, secVal string) (companyDevelopers []common.CompanyDeveloper, err error) {
	return d.companyDevelopers, d.err
}

func (d *DummyDbMan) GetAppCredentials(org, priKey, priVal, secKey, secVal string) (appCredentials []common.AppCredential, err error) {
	return d.appCredentials, d.err
}

func (d *DummyDbMan) GetDevelopers(org, priKey, priVal, secKey, secVal string) (developers []common.Developer, err error) {
	return d.developers, d.err
}

func (d *DummyDbMan) GetApiProductNames(id string, idType string) ([]string, error) {
	return d.apiProductNames, d.err
}

func (d *DummyDbMan) GetAppNames(id string, idType string) ([]string, error) {
	return d.appNames, d.err
}

func (d *DummyDbMan) GetComNames(id string, idType string) ([]string, error) {
	return d.comNames, d.err
}

func (d *DummyDbMan) GetDevEmailByDevId(devId string) (string, error) {
	return d.email, d.err
}

func (d *DummyDbMan) GetStatus(id, t string) (string, error) {
	return d.status, d.err
}
