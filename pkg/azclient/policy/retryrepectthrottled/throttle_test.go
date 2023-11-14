/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package retryrepectthrottled_test

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/policy/retryrepectthrottled"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/utils"
)

var _ = Describe("Throttle", func() {
	Describe("Throttle", func() {
		It("should respect retry-after", func() {
			once := sync.Once{}
			throttlePolicy := &retryrepectthrottled.ThrottlingPolicy{}
			pipeline := runtime.NewPipeline("testmodule", "v0.1.0", runtime.PipelineOptions{}, &policy.ClientOptions{
				PerCallPolicies: []policy.Policy{
					throttlePolicy,
					utils.FuncPolicyWrapper(
						func(req *policy.Request) (*http.Response, error) {
							resp := &http.Response{
								StatusCode: http.StatusOK,
								Body:       http.NoBody,
							}
							once.Do(func() {
								resp = &http.Response{
									StatusCode: http.StatusTooManyRequests,
									Body:       http.NoBody,
									Header: http.Header{
										"Retry-After": []string{"10"},
									},
								}
							})
							return resp, nil
						},
					),
				},
			})
			req, err := runtime.NewRequest(context.Background(), http.MethodPut, "http://localhost:8080")
			Expect(err).NotTo(HaveOccurred())
			err = req.SetBody(streaming.NopCloser(strings.NewReader(`{"etag":"etag"}`)), "application/json")
			Expect(err).NotTo(HaveOccurred())
			_, err = pipeline.Do(req)
			Expect(err).To(HaveOccurred())
			req, err = runtime.NewRequest(context.Background(), http.MethodPut, "http://localhost:8080")
			Expect(err).NotTo(HaveOccurred())
			_, err = pipeline.Do(req)
			Expect(err).To(HaveOccurred())
			throttlePolicy.RetryAfterWriter = time.Now().Add(-time.Second * 10)
			req, err = runtime.NewRequest(context.Background(), http.MethodPut, "http://localhost:8080")
			Expect(err).NotTo(HaveOccurred())
			err = req.SetBody(streaming.NopCloser(strings.NewReader(`{"etag":"etag"}`)), "application/json")
			Expect(err).NotTo(HaveOccurred())
			_, err = pipeline.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
