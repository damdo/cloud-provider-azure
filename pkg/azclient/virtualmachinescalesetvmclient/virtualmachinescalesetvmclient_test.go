// /*
// Copyright The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

// Code generated by client-gen. DO NOT EDIT.
package virtualmachinescalesetvmclient

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var beforeAllFunc func(context.Context)
var afterAllFunc func(context.Context)
var additionalTestCases func()

var _ = Describe("VirtualMachineScaleSetVMsClient", Ordered, func() {

	if beforeAllFunc != nil {
		BeforeAll(beforeAllFunc)
	}

	if additionalTestCases != nil {
		additionalTestCases()
	}

	When("get requests are raised", func() {
		It("should not return error", func(ctx context.Context) {
			newResource, err := realClient.Get(ctx, resourceGroupName, parentResourceName, resourceName)
			Expect(err).NotTo(HaveOccurred())
			Expect(newResource).NotTo(BeNil())
		})
	})
	When("invalid get requests are raised", func() {
		It("should return 404 error", func(ctx context.Context) {
			newResource, err := realClient.Get(ctx, resourceGroupName, parentResourceName, resourceName+"notfound")
			Expect(err).To(HaveOccurred())
			Expect(newResource).To(BeNil())
		})
	})

	When("deletion requests are raised", func() {
		It("should not return error", func(ctx context.Context) {
			err = realClient.Delete(ctx, resourceGroupName, parentResourceName, resourceName)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	if afterAllFunc != nil {
		AfterAll(afterAllFunc)
	}
})
