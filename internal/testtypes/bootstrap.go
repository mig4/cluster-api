/*
Copyright 2021 The Kubernetes Authors.

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

package testtypes

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// BootstrapGroupVersion is group version used for bootstrap objects.
	BootstrapGroupVersion = schema.GroupVersion{Group: "bootstrap.cluster.x-k8s.io", Version: "v1alpha4"}

	// GenericBootstrapConfigCRD is a generic boostrap CRD.
	GenericBootstrapConfigCRD = generateCRD(BootstrapGroupVersion.WithKind("GenericBootstrapConfig"))

	// GenericBootstrapConfigTemplateCRD is a generic boostrap template CRD.
	GenericBootstrapConfigTemplateCRD = generateCRD(BootstrapGroupVersion.WithKind("GenericBootstrapConfigTemplate"))
)
