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
package loadbalancerclient

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armnetwork "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3"
	. "github.com/onsi/gomega"
)

var pipClient *armnetwork.PublicIPAddressesClient
var pipResource *armnetwork.PublicIPAddress

func init() {
	additionalTestCases = func() {
	}

	beforeAllFunc = func(ctx context.Context) {
		pipClient, err = armnetwork.NewPublicIPAddressesClient(recorder.SubscriptionID(), recorder.TokenCredential(), &arm.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: recorder.HTTPClient(),
			},
		})

		poller, err := pipClient.BeginCreateOrUpdate(ctx, resourceGroupName, "pip1", armnetwork.PublicIPAddress{
			Location: to.Ptr(location),
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		resp, err := poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 1 * time.Second})
		Expect(err).NotTo(HaveOccurred())
		pipResource = &resp.PublicIPAddress
		newResource = &armnetwork.LoadBalancer{
			Location: to.Ptr(location),
			Properties: &armnetwork.LoadBalancerPropertiesFormat{
				FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
					{
						Name: to.Ptr("frontendConfig1"),
						Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: pipResource,
						},
					},
				},
			},
		}
	}
	afterAllFunc = func(ctx context.Context) {
		poller, err := pipClient.BeginDelete(ctx, resourceGroupName, "pip1", nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 10 * time.Second})
		Expect(err).NotTo(HaveOccurred())
	}
}
