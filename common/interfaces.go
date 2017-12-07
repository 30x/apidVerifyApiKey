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

import "github.com/apid/apid-core/cipher"

type ApiManagerInterface interface {
	InitAPI()
}

type DbManagerInterface interface {
	SetDbVersion(string)
	GetDbVersion() string
	GetKmsAttributes(tenantId string, entities ...string) map[string][]Attribute
}

type CipherManagerInterface interface {
	// If input is encrypted, it decodes the input with base64,
	// and then decrypt it. Otherwise, original input is returned.
	// An encrypted input should be ciphertext prepended with algorithm. An unencrypted input can have any other format.
	// An example input is "{AES/ECB/PKCS5Padding}2jX3V3dQ5xB9C9Zl9sqyo8pmkvVP10rkEVPVhmnLHw4=".
	TryDecryptBase64(input string, org string) (output string, err error)
	// It encrypts the input, and then encodes the ciphertext with base64.
	// The returned string is the base64 encoding of the encrypted input, prepended with algorithm.
	// An example output is "{AES/ECB/PKCS5Padding}2jX3V3dQ5xB9C9Zl9sqyo8pmkvVP10rkEVPVhmnLHw4="
	EncryptBase64(input string, org string, mode cipher.Mode, padding cipher.Padding) (output string, err error)
}
