package containerregistry

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 0.17.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// Registry is an object that represents a container registry.
type Registry struct {
	autorest.Response   `json:"-"`
	ID                  *string             `json:"id,omitempty"`
	Name                *string             `json:"name,omitempty"`
	Type                *string             `json:"type,omitempty"`
	Location            *string             `json:"location,omitempty"`
	Tags                *map[string]*string `json:"tags,omitempty"`
	*RegistryProperties `json:"properties,omitempty"`
}

// RegistryCredentials is the result of a request to get the administrator
// login credentials for a container registry.
type RegistryCredentials struct {
	autorest.Response `json:"-"`
	Username          *string `json:"username,omitempty"`
	Password          *string `json:"password,omitempty"`
}

// RegistryListResult is the result of a request to list container registries.
type RegistryListResult struct {
	autorest.Response `json:"-"`
	Value             *[]Registry `json:"value,omitempty"`
	NextLink          *string     `json:"nextLink,omitempty"`
}

// RegistryListResultPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client RegistryListResult) RegistryListResultPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// RegistryNameCheckRequest is a request to check whether the container
// registry name is available.
type RegistryNameCheckRequest struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

// RegistryNameStatus is the result of a request to check the availability of
// a container registry name.
type RegistryNameStatus struct {
	autorest.Response `json:"-"`
	NameAvailable     *bool   `json:"nameAvailable,omitempty"`
	Reason            *string `json:"reason,omitempty"`
	Message           *string `json:"message,omitempty"`
}

// RegistryProperties is the properties of a container registry.
type RegistryProperties struct {
	LoginServer      *string                   `json:"loginServer,omitempty"`
	CreationDate     *date.Time                `json:"creationDate,omitempty"`
	AdminUserEnabled *bool                     `json:"adminUserEnabled,omitempty"`
	StorageAccount   *StorageAccountProperties `json:"storageAccount,omitempty"`
}

// RegistryPropertiesUpdateParameters is the parameters for updating the
// properties of a container registry.
type RegistryPropertiesUpdateParameters struct {
	AdminUserEnabled *bool                     `json:"adminUserEnabled,omitempty"`
	StorageAccount   *StorageAccountProperties `json:"storageAccount,omitempty"`
}

// RegistryUpdateParameters is the parameters for updating a container
// registry.
type RegistryUpdateParameters struct {
	Tags                                *map[string]*string `json:"tags,omitempty"`
	*RegistryPropertiesUpdateParameters `json:"properties,omitempty"`
}

// Resource is an Azure resource.
type Resource struct {
	ID       *string             `json:"id,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Type     *string             `json:"type,omitempty"`
	Location *string             `json:"location,omitempty"`
	Tags     *map[string]*string `json:"tags,omitempty"`
}

// StorageAccountProperties is the properties of a storage account for a
// container registry.
type StorageAccountProperties struct {
	Name      *string `json:"name,omitempty"`
	AccessKey *string `json:"accessKey,omitempty"`
}
