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
	"github.com/apid/apidVerifyApiKey/common"
	"os"
	"testing"
)

const testTempDirBase = "./tmp/"

var (
	testTempDir string
)

var _ = BeforeSuite(func() {
	_ = os.MkdirAll(testTempDirBase, os.ModePerm)
	s := factory.DefaultServicesFactory()
	apid.Initialize(s)
	SetApidServices(s, apid.Log())
	common.SetApidServices(s, s.Log())
})

var _ = AfterSuite(func() {
	apid.Events().Close()
	os.RemoveAll(testTempDirBase)
})

func TestVerifyAPIKey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VerifyAPIKey Suite")
}
