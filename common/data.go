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
package common

import (
	"database/sql"
	"encoding/json"
	"github.com/apid/apid-core"
	"strings"
	"sync"
	"unicode/utf8"
)

type DbManager struct {
	Data          apid.DataService
	Db            apid.DB
	DbMux         sync.RWMutex
	CipherManager CipherManagerInterface
	dbVersion     string
}

const (
	sql_GET_KMS_ATTRIBUTES_FOR_TENANT = `select entity_id, name, value from kms_attributes where tenant_id = $1`
)

var (
	services apid.Services
	log      apid.LogService
)

func SetApidServices(s apid.Services, l apid.LogService) {
	services = s
	log = l
}

func (dbc *DbManager) SetDbVersion(version string) {
	db, err := dbc.Data.DBVersion(version)
	if err != nil {
		log.Panicf("Unable to access database: %v", err)
	}
	dbc.DbMux.Lock()
	dbc.Db = db
	dbc.DbMux.Unlock()
	dbc.dbVersion = version
}

func (dbc *DbManager) GetDb() apid.DB {
	dbc.DbMux.RLock()
	defer dbc.DbMux.RUnlock()
	return dbc.Db
}

func (dbc *DbManager) GetDbVersion() string {
	return dbc.dbVersion
}

func (dbc *DbManager) GetKmsAttributes(tenantId string, entities ...string) map[string][]Attribute {

	db := dbc.Db
	var attName, attValue, entity_id sql.NullString
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
	return mapOfAttributes
}

func (dbc *DbManager) GetOrgs() (orgs []string, err error) {
	db := dbc.GetDb()
	rows, err := db.Query(`SELECT DISTINCT org FROM edgex_data_scope`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmp sql.NullString
		if err = rows.Scan(&tmp); err != nil {
			return nil, err
		}
		if tmp.Valid {
			orgs = append(orgs, tmp.String)
		}
	}
	err = rows.Err()
	return
}

func AddIndexes(version string) error {
	db, err := services.Data().DBVersion(version)
	if err != nil {
		log.Errorf("Unable to access database: %v", err)
		return err
	}
	log.Debugf("adding indexes to sqlite file")
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("AddIndexes: Unable to get DB tx Err: {%v}", err)
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`
	CREATE INDEX IF NOT EXISTS mp_appcred_id on KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER (appcred_id);
	CREATE INDEX IF NOT EXISTS mp_app_id on KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER (app_id);
	CREATE INDEX IF NOT EXISTS mp_apiprdt_id on KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER (apiprdt_id);
	CREATE INDEX IF NOT EXISTS app_name on KMS_APP (name);
	CREATE INDEX IF NOT EXISTS app_company_id on KMS_APP (company_id);
	CREATE INDEX IF NOT EXISTS app_developer_id on KMS_APP (developer_id);
	CREATE INDEX IF NOT EXISTS dev_email on KMS_DEVELOPER (email);
	CREATE INDEX IF NOT EXISTS com_name on KMS_COMPANY (name);
	CREATE INDEX IF NOT EXISTS com_dev_com_id on KMS_COMPANY_DEVELOPER (company_id);
	CREATE INDEX IF NOT EXISTS com_dev_dev_id on KMS_COMPANY_DEVELOPER (developer_id);
	CREATE INDEX IF NOT EXISTS cred_app_id on KMS_APP_CREDENTIAL (app_id);
	CREATE INDEX IF NOT EXISTS org_tenant_id on KMS_ORGANIZATION (tenant_id);
	CREATE INDEX IF NOT EXISTS org_name on KMS_ORGANIZATION (name);
	`)
	if err != nil {
		log.Errorf("AddIndexes: Tx Exec Err: {%v}", err)
		return err
	}
	if err = tx.Commit(); err != nil {
		log.Errorf("Commit error in AddIndexes: %v", err)
	}
	return err
}

func JsonToStringArray(fjson string) []string {
	var array []string
	if err := json.Unmarshal([]byte(fjson), &array); err == nil {
		return array
	}
	s := strings.TrimPrefix(fjson, "{")
	s = strings.TrimSuffix(s, "}")
	if utf8.RuneCountInString(s) > 0 {
		array = strings.Split(s, ",")
	}
	return array
}
