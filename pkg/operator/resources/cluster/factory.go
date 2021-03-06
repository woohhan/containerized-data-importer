/*
Copyright 2018 The CDI Authors.

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

package cluster

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"kubevirt.io/containerized-data-importer/pkg/operator/resources/utils"
)

// FactoryArgs contains the required parameters to generate all cluster-scoped resources
type FactoryArgs struct {
	Namespace string
}

type factoryFunc func(*FactoryArgs) []runtime.Object

var factoryFunctions = map[string]factoryFunc{
	"apiserver-rbac":   createAPIServerResources,
	"controller-rbac":  createControllerResources,
	"crd-resources":    createCRDResources,
	"uploadproxy-rbac": createUploadProxyResources,
	"aggregate-roles":  createAggregateClusterRoles,
}

func createCRDResources(args *FactoryArgs) []runtime.Object {
	return []runtime.Object{
		createDataVolumeCRD(),
		createCDIConfigCRD(),
	}
}

// CreateAllResources creates all cluster-wide resources
func CreateAllResources(args *FactoryArgs) ([]runtime.Object, error) {
	var resources []runtime.Object
	for group := range factoryFunctions {
		rs, err := CreateResourceGroup(group, args)
		if err != nil {
			return nil, err
		}
		resources = append(resources, rs...)
	}
	return resources, nil
}

// CreateResourceGroup creates all cluster resources fr a specific group/component
func CreateResourceGroup(group string, args *FactoryArgs) ([]runtime.Object, error) {
	f, ok := factoryFunctions[group]
	if !ok {
		return nil, fmt.Errorf("group %s does not exist", group)
	}
	resources := f(args)
	utils.ValidateGVKs(resources)
	return resources, nil
}

// GetClusterRolePolicyRules returns all cluster PolicyRules
func GetClusterRolePolicyRules() []rbacv1.PolicyRule {
	result := getAPIServerClusterPolicyRules()
	result = append(result, getControllerClusterPolicyRules()...)
	result = append(result, getUploadProxyClusterPolicyRules()...)
	return result
}
