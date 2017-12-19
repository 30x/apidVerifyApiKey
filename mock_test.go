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
package apidApiMetadata

import (
	"github.com/apid/apid-core/cipher"
	"github.com/apid/apidApiMetadata/common"
)

type DummyDbMan struct {
	version string
}

func (d *DummyDbMan) GetOrgs() (orgs []string, err error) {
	return
}

func (d *DummyDbMan) SetDbVersion(v string) {
	d.version = v
}
func (d *DummyDbMan) GetDbVersion() string {
	return d.version
}

func (d *DummyDbMan) GetKmsAttributes(tenantId string, entities ...string) map[string][]common.Attribute {
	return nil
}

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
