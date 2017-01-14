package apidVerifyApiKey

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/30x/apid"
	"github.com/30x/apid/factory"
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

	apid.InitializePlugins()

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
