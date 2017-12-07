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
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const RegEncrypted = `^\{[0-9A-Za-z]+/[0-9A-Za-z]+/[0-9A-Za-z]+\}.`
const retrieveEncryptKeyPath = ""

var RegexpEncrypted = regexp.MustCompile(RegEncrypted)

const (
	EncryptAes = "AES"
)

type KmsCipherManager struct {
	// org-level key map {organization: key}
	key map[string][]byte
	// org-level AesCipher map {organization: AesCipher}
	aes    map[string]*cipher.AesCipher
	mutex  *sync.RWMutex
	client *http.Client
}

func (c *KmsCipherManager) retrieveKey(org string) (key []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, retrieveEncryptKeyPath, nil)
	if err != nil {
		return
	}
	res, err := c.client.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("retrieve encryption key failed for org [%v] with status: %v", org, res.Status)
		return
	}
	defer res.Body.Close()
	key64, err := ioutil.ReadAll(res.Body)
	key, err = base64.StdEncoding.DecodeString(string(key64))
	return
}

func (c *KmsCipherManager) getAesCipher(org string) (*cipher.AesCipher, error) {
	// if exists
	c.mutex.RLock()
	if a := c.aes[org]; a != nil {
		c.mutex.RUnlock()
		return a, nil
	}
	c.mutex.RUnlock()
	// if not exists
	key, err := c.retrieveKey(org)
	if err != nil {
		log.Errorf("getAesCipher error for org [%v] when retrieveKey: %v", org, err)
		return nil, err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.key[org] = key
	a, err := cipher.CreateAesCipher(key)
	if err != nil {
		log.Errorf("getAesCipher error for org [%v] when CreateAesCipher: %v", org, err)
		return nil, err
	}
	c.aes[org] = a
	return a, nil
}

// If input is encrypted, it decodes the input with base64,
// and then decrypt it. Otherwise, original input is returned.
// An encrypted input should be ciphertext prepended with algorithm. An unencrypted input can have any other format.
// An example of encrypted input is "{AES/ECB/PKCS5Padding}2jX3V3dQ5xB9C9Zl9sqyo8pmkvVP10rkEVPVhmnLHw4=".
func (c *KmsCipherManager) TryDecryptBase64(input string, org string) (output string, err error) {
	if !IsEncrypted(input) {
		output = input
		return
	}

	text, mode, padding, err := GetCiphertext(input)
	if err != nil {
		log.Errorf("Get ciphertext of [%v] failed: [%v], considered as unencrypted!", input, err)
		return
	}
	bytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Errorf("Decode base64 of [%v] failed: [%v], considered as unencrypted!", text, err)
		return
	}
	aes, err := c.getAesCipher(org)
	if err != nil {
		return
	}
	plaintext, err := aes.Decrypt(bytes, mode, padding)
	if err != nil {
		log.Errorf("Decrypt of [%v] failed: [%v], considered as unencrypted!", bytes, err)
		return
	}
	output = string(plaintext)
	return
}

// It encrypts the input, and then encodes the ciphertext with base64.
// The returned string is the base64 encoding of the encrypted input, prepended with algorithm.
// An example output is "{AES/ECB/PKCS5Padding}2jX3V3dQ5xB9C9Zl9sqyo8pmkvVP10rkEVPVhmnLHw4="
func (c *KmsCipherManager) EncryptBase64(input string, org string, mode cipher.Mode, padding cipher.Padding) (output string, err error) {
	aes, err := c.getAesCipher(org)
	if err != nil {
		return
	}
	ciphertext, err := aes.Encrypt([]byte(input), mode, padding)
	if err != nil {
		return
	}
	output = fmt.Sprintf("{%s/%s/%s}%s", EncryptAes, mode, padding, base64.StdEncoding.EncodeToString(ciphertext))
	return
}

// TODO: make sure this regex has no false positive for all possible inputs
func IsEncrypted(input string) (encrypted bool) {
	return RegexpEncrypted.Match([]byte(input))
}

func GetCiphertext(input string) (ciphertext string, mode cipher.Mode, padding cipher.Padding, err error) {
	l := strings.SplitN(input, "}", 2)
	if len(l) != 2 {
		err = fmt.Errorf("invalid input for GetCiphertext: %v", input)
		return
	}
	ciphertext = l[1]
	l = strings.Split(strings.TrimLeft(l[0], "{"), "/")
	if len(l) != 3 {
		err = fmt.Errorf("invalid input for GetCiphertext: %v", input)
		return
	}
	// encryption algorithm
	if strings.ToUpper(l[0]) != EncryptAes {
		err = fmt.Errorf("unsupported algorithm for GetCiphertext: %v", l[0])
		return
	}
	// mode
	mode = cipher.Mode(strings.ToUpper(l[1]))
	// padding
	padding = cipher.Padding(strings.ToUpper(l[2]))
	return
}
