package apidVerifyApiKey

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/30x/apid-core"
	"github.com/30x/apid-core/data"
	"github.com/30x/apid-core/factory"
	"io"
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

	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == apiPath {
			handleRequest(w, req)
		}
	}))
})

var _ = BeforeEach(func() {
	dbPath := data.DBPath("common/" + "test")
	_ = os.MkdirAll(dbPath[0:len(dbPath)-7], 0700)
	dst, err := os.Create(dbPath)
	Expect(err).NotTo(HaveOccurred())

	src, err := os.Open("./mockdb.sqlite3")
	Expect(err).NotTo(HaveOccurred())

	defer src.Close()

	log.Info("Copying mockdb.sqlite3 to " + dbPath)
	_, err = io.Copy(dst, src)
	Expect(err).NotTo(HaveOccurred())

	dst.Close()

	db, err := apid.Data().DBVersion("test")

	Expect(err).NotTo(HaveOccurred())
	setDB(db)
	createTables(db)
	createApidClusterTables(db)
	addScopes(db)
})

var _ = AfterSuite(func() {
	apid.Events().Close()
	if testServer != nil {
		testServer.Close()
	}
	//for debugging, comment this out, check the db paths in log output
	//os.RemoveAll(testTempDir)
})

func TestVerifyAPIKey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VerifyAPIKey Suite")
}
