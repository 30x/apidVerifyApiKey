package common

import (
	"database/sql"
	"github.com/apid/apid-core"
	"strings"
	"sync"
)

type DbManager struct {
	Data      apid.DataService
	Db        apid.DB
	DbMux     sync.RWMutex
	dbVersion string
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
