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
	"testing"
)

var (
	testTempDir string
	testServer  *httptest.Server
)

var _ = BeforeSuite(func() {
	apid.Initialize(factory.DefaultServicesFactory())

	config := apid.Config()

	var err error
	testTempDir, err = ioutil.TempDir("", "api_test")
	Expect(err).NotTo(HaveOccurred())

	config.Set("data_path", testTempDir)

	apid.InitializePlugins("")

	db, err := apid.Data().DB()
	Expect(err).NotTo(HaveOccurred())
	setDB(db)
	createTables(db)
	createApidClusterTables(db)
	addScopes(db)
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == apiPath {
			handleRequest(w, req)
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
		for j = k; j < 10 + k; j++ {
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
		for j = k; j < 100 + k; j++ {
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
	db.QueryRow("select count(*) from data_scope").Scan(&count)
	log.Info("Found ", count)
}
