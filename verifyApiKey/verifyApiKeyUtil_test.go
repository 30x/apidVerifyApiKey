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

package verifyApiKey

import (
	"github.com/apid/apidApiMetadata/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Validate Env", func() {

	It("validation1", func() {
		s := contains([]string{"foo", "bar"}, "foo")
		Expect(s).Should(BeTrue())
	})
	It("validation2", func() {
		s := contains([]string{"foo", "bar"}, "bar")
		Expect(s).Should(BeTrue())
	})
	It("validation3", func() {
		s := contains([]string{"foo", "bar"}, "xxx")
		Expect(s).Should(BeFalse())
	})
	It("validation4", func() {
		s := contains([]string{}, "xxx")
		Expect(s).Should(BeFalse())
	})
})

var _ = Describe("Validate common.JsonToStringArray", func() {

	It("should tranform simple valid json", func() {
		array := common.JsonToStringArray("[\"test-1\", \"test-2\"]")
		Expect(reflect.DeepEqual(array, []string{"test-1", "test-2"})).Should(BeTrue())
	})
	It("should tranform simple single valid json", func() {
		array := common.JsonToStringArray("[\"test-1\"]")
		Expect(reflect.DeepEqual(array, []string{"test-1"})).Should(BeTrue())
	})
	It("should tranform simple fake json", func() {
		s := common.JsonToStringArray("{test-1,test-2}")
		Expect(reflect.DeepEqual(s, []string{"test-1", "test-2"})).Should(BeTrue())
	})
	It("should tranform simple single valued fake json", func() {
		s := common.JsonToStringArray("{test-1}")
		Expect(reflect.DeepEqual(s, []string{"test-1"})).Should(BeTrue())
	})
	It("space between fields considered as valid char", func() {
		s := common.JsonToStringArray("{test-1, test-2}")
		Expect(reflect.DeepEqual(s, []string{"test-1", " test-2"})).Should(BeTrue())
	})
	It("remove only last braces", func() {
		s := common.JsonToStringArray("{test-1,test-2}}")
		Expect(reflect.DeepEqual(s, []string{"test-1", "test-2}"})).Should(BeTrue())
	})

})
