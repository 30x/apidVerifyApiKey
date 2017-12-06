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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/apid/apid-core"
	"github.com/apid/apid-core/factory"
	"github.com/apid/apidApiMetadata/common"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var testTempDirBase string

func initSetup(s apid.Services) (apid.PluginData, error) {
	SetApidServices(s, apid.Log())
	common.SetApidServices(s, s.Log())
	return common.PluginData, nil
}

var _ = BeforeSuite(func() {
	apid.RegisterPlugin(initSetup, common.PluginData)
	var err error
	testTempDirBase, err = ioutil.TempDir("", "verify_apikey_")
	Expect(err).Should(Succeed())
	apid.Initialize(factory.DefaultServicesFactory())
	apid.InitializePlugins("0.0.0")
	go services.API().Listen()
	time.Sleep(time.Second)
}, 2)

var _ = AfterSuite(func() {
	apid.Events().Close()
	Expect(os.RemoveAll(testTempDirBase)).Should(Succeed())
})

func TestVerifyAPIKey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AccessEntity Suite")
}
