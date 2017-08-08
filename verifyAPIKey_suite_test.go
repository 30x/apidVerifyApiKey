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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

var (
	testTempDir     string
	testServer      *httptest.Server
	testSyncHandler apigeeSyncHandler
)

var _ = BeforeSuite(func() {
	var err error
	testTempDir, err = ioutil.TempDir("", "api_test")
	s := factory.DefaultServicesFactory()
	apid.Initialize(s)
	config := apid.Config()
	config.Set("data_path", testTempDir)
	config.Set("log_level", "DEBUG")
	log = apid.Log()
	Expect(err).NotTo(HaveOccurred())

	apid.InitializePlugins("")

	db, err := apid.Data().DB()
	Expect(err).NotTo(HaveOccurred())

	dbMan := &dbManager{
		data:  s.Data(),
		dbMux: sync.RWMutex{},
	}
	dbMan.initDb()
	apiMan := apiManager{
		dbMan:             dbMan,
		verifiersEndpoint: apiPath,
	}

	testSyncHandler = apigeeSyncHandler{
		dbMan:  dbMan,
		apiMan: apiMan,
	}

	testSyncHandler.initListener(s)

	createTables(db)
	createApidClusterTables(db)
	addScopes(db)
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == apiPath {
			apiMan.handleRequest(w, req)
		}
	}))

	createTestData(db)
})

var _ = AfterSuite(func() {
	apid.Events().Close()
	if testServer != nil {
		testServer.Close()
	}
	os.RemoveAll(testTempDir)
})

func TestVerifyAPIKey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VerifyAPIKey Suite")
}

func createTestData(db apid.DB) {
	txn, err := db.Begin()
	Expect(err).ShouldNot(HaveOccurred())
	// api products
	for i := 0; i < 10; i++ {
		generateTestApiProduct(i, txn)
	}
	// developers
	for i := 0; i < 10; i++ {
		generateTestDeveloper(i, txn)
	}

	// application
	var j, k int
	for i := 0; i < 10; i++ {
		for j = k; j < 10+k; j++ {
			generateTestApp(j, i, txn)
		}
		k = j
	}
	// app credentials
	for i := 0; i < 10; i++ {
		generateTestAppCreds(i, txn)
	}
	// api product mapper
	for i := 0; i < 10; i++ {
		generateTestApiProductMapper(i, txn)
	}

	// Following are data for company
	// api products
	for i := 100; i < 110; i++ {
		generateTestApiProduct(i, txn)
	}

	// companies
	for i := 100; i < 110; i++ {
		generateTestCompany(i, txn)
	}

	// company developers
	for i := 100; i < 110; i++ {
		generateTestCompanyDeveloper(i, txn)
	}

	// application
	k = 100
	for i := 100; i < 110; i++ {
		for j = k; j < 100+k; j++ {
			generateTestAppCompany(j, i, txn)
		}
		k = j
	}
	// app credentials
	for i := 100; i < 110; i++ {
		generateTestAppCreds(i, txn)
	}
	// api product mapper
	for i := 100; i < 110; i++ {
		generateTestApiProductMapper(i, txn)
	}

	txn.Commit()
	var count int64
	db.QueryRow("select count(*) from EDGEX_DATA_SCOPE").Scan(&count)
	log.Info("Found ", count)
}
