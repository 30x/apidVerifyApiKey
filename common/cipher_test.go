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
	"encoding/base64"
	"fmt"
	"github.com/apid/apid-core/cipher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Cipher Test", func() {
	var testCipherMan *KmsCipherManager
	var testCount int
	var testOrg string
	plaingtext := "aUWQKgAwmaR0p2kY"
	cipher64 := "{AES/ECB/PKCS5Padding}2jX3V3dQ5xB9C9Zl9sqyo8pmkvVP10rkEVPVhmnLHw4="
	key := []byte{2, 122, 212, 83, 150, 164, 180, 4, 148, 242, 65, 189, 3, 188, 76, 247}
	BeforeEach(func() {
		testCount++
		testOrg = fmt.Sprintf("org%d", testCount)
	})

	Context("Encryption/Decryption", func() {
		BeforeEach(func() {
			testCipherMan = CreateCipherManager(nil, "")
			// set key locally
			testCipherMan.key[testOrg] = key
			var err error
			testCipherMan.aes[testOrg], err = cipher.CreateAesCipher(key)
			Expect(err).Should(Succeed())
		})

		It("Encryption", func() {
			Expect(testCipherMan.EncryptBase64(plaingtext, testOrg, cipher.ModeEcb, cipher.PaddingPKCS5)).
				Should(Equal(cipher64))
		})

		It("Decryption", func() {
			Expect(testCipherMan.TryDecryptBase64(cipher64, testOrg)).Should(Equal(plaingtext))
		})

		It("Try to decrypt unencrypted input", func() {
			Expect(testCipherMan.TryDecryptBase64(plaingtext, testOrg)).Should(Equal(plaingtext))
		})
	})

	Context("Retrieve new key", func() {
		var server *httptest.Server
		BeforeEach(func() {
			// set key server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer GinkgoRecover()
				Expect(w.Write([]byte(base64.StdEncoding.EncodeToString(key)))).Should(Equal(24))
			}))
			time.Sleep(100 * time.Millisecond)
			testCipherMan = CreateCipherManager(&http.Client{}, server.URL)
		})

		AfterEach(func() {
			server.Close()
		})

		It("Encryption", func() {
			Expect(testCipherMan.EncryptBase64(plaingtext, testOrg, cipher.ModeEcb, cipher.PaddingPKCS5)).
				Should(Equal(cipher64))
		})

		It("Decryption", func() {
			Expect(testCipherMan.TryDecryptBase64(cipher64, testOrg)).Should(Equal(plaingtext))
		})
	})
})
