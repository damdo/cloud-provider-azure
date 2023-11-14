/*
Copyright 2019 The Kubernetes Authors.

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

package utils

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	acr "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"

	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/utils/pointer"
)

// CreateContainerRegistry creates a test acr
func (tc *AzureTestClient) CreateContainerRegistry() (acr.Registry, error) {
	acrClient := tc.createACRClient()
	rgName, location := tc.GetResourceGroup(), tc.GetLocation()
	acrName := "e2eacr" + string(uuid.NewUUID())[0:4]
	template := acr.Registry{
		SKU: &acr.SKU{
			Name: to.Ptr(acr.SKUNameBasic),
			Tier: to.Ptr(acr.SKUTierBasic),
		},
		Name:     pointer.String(acrName),
		Location: pointer.String(location),
	}

	Logf("Creating acr %s in resource group %s.", acrName, rgName)
	registry, err := acrClient.Create(context.Background(), rgName, acrName, template)
	if err != nil {
		return acr.Registry{}, err
	}

	return *registry, nil
}

// DeleteContainerRegistry deletes an existing acr
func (tc *AzureTestClient) DeleteContainerRegistry(registryName string) (err error) {
	acrClient := tc.createACRClient()
	rgName := tc.GetResourceGroup()

	Logf("Deleting acr %s in resource group %s.", registryName, rgName)
	err = acrClient.Delete(context.Background(), rgName, registryName)
	if err != nil {
		return fmt.Errorf("failed to delete acr %s in resource group %s with error: %w", registryName, rgName, err)
	}
	return nil
}

func AZACRLogin() (err error) {
	authConfig, _, _, err := azureAuthConfigFromTestProfile()
	if err != nil {
		return err
	}

	Logf("Attempting az login with azure cred.")
	//nolint:gosec // G204 ignore this!
	cmd := exec.Command("az", "login", "--service-principal",
		"--username", authConfig.AADClientID,
		"--password", authConfig.AADClientSecret,
		"--tenant", authConfig.TenantID)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("az failed to login with error: %w", err)
	}
	Logf("az login success.")

	cmd = exec.Command("az", "account", "show")
	var output []byte
	if output, err = cmd.Output(); err != nil {
		return fmt.Errorf("az failed to account show with output: %s\n error: %w", string(output), err)
	}
	Logf("az account show success.")
	return nil
}

// PushImageToACR pull an image from Docker Hub and push
// it to the given azure container registry
func (tc *AzureTestClient) PushImageToACR(registryName, image string) (string, error) {
	acrClient := tc.createACRClient()
	rgName := tc.GetResourceGroup()

	err := acrClient.ImportImage(context.Background(), rgName, registryName, acr.ImportImageParameters{
		Source: &acr.ImportSource{
			RegistryURI: pointer.String("docker.io"),
			SourceImage: pointer.String("library/" + image + ":latest"),
		},
		TargetTags: []*string{
			to.Ptr(image + ":latest"),
		},
		Mode: to.Ptr(acr.ImportModeNoForce),
	})
	if err != nil {
		return "", fmt.Errorf("failed to import image %s from docker hub to acr %s with error: %w", image, registryName, err)
	}
	return "latest", nil
}

// AZACRCacheCreate enables acr cache for a image.
func AZACRCacheCreate(acrName, ruleName, imageURL, imageName string) (err error) {
	Logf("Attempting az acr cache create for image URL %q.", imageURL)
	cmd := exec.Command("az", "acr", "cache", "create",
		"-r", acrName,
		"-n", ruleName,
		"-s", imageURL,
		"-t", imageName)
	if output, err := cmd.Output(); err != nil {
		return fmt.Errorf("az acr cache create failed with output %s\n error: %w", string(output), err)
	}
	Logf("az acr cache create success.")
	return nil
}
