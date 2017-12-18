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
	"encoding/json"
	"encoding/xml"
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
		Context("Retrieve new key with lazy method", func() {
			BeforeEach(func() {
				// set key server
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()
					Expect(r.URL.Path).Should(Equal(retrieveEncryptKeyPath))
					Expect(r.URL.Query().Get(parameterOrganization)).Should(Equal(testOrg))
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

		Context("Retrieve new keys during initialization", func() {

			AfterEach(func() {
				server.Close()
			})

			It("Retrieve Key happy path", func() {
				// set key server
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()
					Expect(r.URL.Path).Should(Equal(retrieveEncryptKeyPath))
					Expect(r.URL.Query().Get(parameterOrganization)).Should(HavePrefix(testOrg))
					Expect(w.Write([]byte(base64.StdEncoding.EncodeToString(key)))).Should(Equal(24))
				}))
				time.Sleep(100 * time.Millisecond)
				testCipherMan = CreateCipherManager(&http.Client{}, server.URL)

				//test 2 orgs
				testOrg1 := testOrg + "_1"
				testCipherMan.AddOrgs([]string{testOrg, testOrg1})
				for {
					time.Sleep(100 * time.Millisecond)
					testCipherMan.mutex.RLock()
					l := len(testCipherMan.aes)
					testCipherMan.mutex.RUnlock()
					if l == 2 {
						//close server to make sure key was retrieved by "AddOrgs"
						server.Close()
						Expect(testCipherMan.EncryptBase64(plaingtext, testOrg, cipher.ModeEcb, cipher.PaddingPKCS5)).
							Should(Equal(cipher64))
						Expect(testCipherMan.EncryptBase64(plaingtext, testOrg1, cipher.ModeEcb, cipher.PaddingPKCS5)).
							Should(Equal(cipher64))
						return
					}
				}
			}, 2)

			It("Retrieve Key should retry for internal server error", func() {
				// set key server
				count := 0
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()
					Expect(r.URL.Path).Should(Equal(retrieveEncryptKeyPath))
					Expect(r.URL.Query().Get(parameterOrganization)).Should(Equal(testOrg))
					count++
					if count == 1 {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if count == 2 {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					Expect(w.Write([]byte(base64.StdEncoding.EncodeToString(key)))).Should(Equal(24))
				}))
				time.Sleep(100 * time.Millisecond)
				testCipherMan = CreateCipherManager(&http.Client{}, server.URL)
				testCipherMan.interval = 100 * time.Millisecond
				//should retry in case of error
				testCipherMan.AddOrgs([]string{testOrg})
				for {
					time.Sleep(100 * time.Millisecond)
					testCipherMan.mutex.RLock()
					aes := testCipherMan.aes[testOrg]
					testCipherMan.mutex.RUnlock()
					if aes != nil {
						//close server to make sure key was retrieved by "AddOrgs"
						server.Close()
						Expect(testCipherMan.EncryptBase64(plaingtext, testOrg, cipher.ModeEcb, cipher.PaddingPKCS5)).
							Should(Equal(cipher64))
						return
					}
				}
			}, 2)

			It("Retrieve Key should stop retrying for JSON organizations.EncryptionKeyDoesNotExist", func() {
				// set key server
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()
					Expect(r.URL.Path).Should(Equal(retrieveEncryptKeyPath))
					Expect(r.URL.Query().Get(parameterOrganization)).Should(Equal(testOrg))

					response := KeyErrorResponse{
						Code:    errorCodeNoKey,
						Message: fmt.Sprintf("Encryption key does not exist for the org [%s].", testOrg),
					}
					bytes, err := json.Marshal(response)
					Expect(err).Should(Succeed())
					w.Header().Set(headerContentType, typeJson)
					w.WriteHeader(http.StatusNotFound)
					Expect(w.Write(bytes)).Should(Equal(len(bytes)))
				}))
				time.Sleep(100 * time.Millisecond)
				testCipherMan = CreateCipherManager(&http.Client{}, server.URL)
				//should stop retrying after one try
				testCipherMan.startRetrieve(testOrg, 100*time.Millisecond, 10*time.Minute)
			}, 2)

			It("Retrieve Key should stop retrying for XML organizations.EncryptionKeyDoesNotExist", func() {
				// set key server
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()
					Expect(r.URL.Path).Should(Equal(retrieveEncryptKeyPath))
					Expect(r.URL.Query().Get(parameterOrganization)).Should(Equal(testOrg))

					response := KeyErrorResponse{
						Code:    errorCodeNoKey,
						Message: fmt.Sprintf("Encryption key does not exist for the org [%s].", testOrg),
					}
					bytes, err := xml.Marshal(response)
					Expect(err).Should(Succeed())
					w.Header().Set(headerContentType, typeXml)
					w.WriteHeader(http.StatusNotFound)
					Expect(w.Write(bytes)).Should(Equal(len(bytes)))
				}))
				time.Sleep(100 * time.Millisecond)
				testCipherMan = CreateCipherManager(&http.Client{}, server.URL)
				//should stop retrying after one try
				testCipherMan.startRetrieve(testOrg, 100*time.Millisecond, 10*time.Minute)
			}, 2)
		})

	})

	Context("IsEncrypted", func() {
		It("IsEncrypted", func() {
			testData := [][]interface{}{
				{"{AES/ECB/PKCS5Padding}foo", true},
				{"AES/ECB/PKCS5Padding}foo", false},
				{"{AES/ECB/PKCS5Paddingfoo", false},
				{"{AES/ECB/}foo", false},
				{"{AES/PKCS5Padding}foo", false},
				{"{AES//PKCS5Padding}foo", false},
				{"foo", false},
			}
			for i := range testData {
				Expect(IsEncrypted(testData[i][0].(string))).Should(Equal(testData[i][1]))
			}
		})
	})
})
