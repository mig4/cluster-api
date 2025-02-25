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

package controllers

import (
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNewMergePatchHelper(t *testing.T) {
	tests := []struct {
		name           string
		original       *unstructured.Unstructured // current
		modified       *unstructured.Unstructured // desired
		wantHasChanges bool
		wantPatch      []byte
	}{
		// Field both in original and in modified --> align to modified

		{
			name: "Field both in original and in modified, no-op when equal",
			original: &unstructured.Unstructured{ // current
				Object: map[string]interface{}{
					"A": "A",
				},
			},
			modified: &unstructured.Unstructured{ // desired
				Object: map[string]interface{}{
					"A": "A",
				},
			},
			wantHasChanges: false,
			wantPatch:      []byte("{}"),
		},
		{
			name: "Field both in original and in modified, align to modified when different",
			original: &unstructured.Unstructured{ // current
				Object: map[string]interface{}{
					"A": "A-changed",
				},
			},
			modified: &unstructured.Unstructured{ // desired
				Object: map[string]interface{}{
					"A": "A",
				},
			},
			wantHasChanges: true,
			wantPatch:      []byte("{\"A\":\"A\"}"),
		},
		{
			name: "Nested field both in original and in modified, no-op when equal",
			original: &unstructured.Unstructured{ // current
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"A": "A",
							},
						},
					},
				},
			},
			modified: &unstructured.Unstructured{ // desired
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"A": "A",
							},
						},
					},
				},
			},
			wantHasChanges: false,
			wantPatch:      []byte("{}"),
		},
		{
			name: "Nested field both in original and in modified, align to modified when different",
			original: &unstructured.Unstructured{ // current
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"A": "A-Changed",
							},
						},
					},
				},
			},
			modified: &unstructured.Unstructured{ // desired
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"A": "A",
							},
						},
					},
				},
			},
			wantHasChanges: true,
			wantPatch:      []byte("{\"spec\":{\"template\":{\"spec\":{\"A\":\"A\"}}}}"),
		},
		{
			name: "Value of type map, enforces entries from modified, preserve entries only in original",
			original: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"map": map[string]string{
						"A": "A-changed",
						"B": "B",
						// C missing
					},
				},
			},
			modified: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"map": map[string]string{
						"A": "A",
						// B missing
						"C": "C",
					},
				},
			},
			wantHasChanges: true,
			wantPatch:      []byte("{\"map\":{\"A\":\"A\",\"C\":\"C\"}}"),
		},
		{
			name: "Value of type Array or Slice, align to modified",
			original: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"slice": []string{
						"D",
						"C",
						"B",
					},
				},
			},
			modified: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"slice": []string{
						"A",
						"B",
						"C",
					},
				},
			},
			wantHasChanges: true,
			wantPatch:      []byte("{\"slice\":[\"A\",\"B\",\"C\"]}"),
		},

		// Field only in modified (not existing in original) --> align to modified

		{
			name: "Field only in modified, align to modified",
			original: &unstructured.Unstructured{ // current
				Object: map[string]interface{}{},
			},
			modified: &unstructured.Unstructured{ // desired
				Object: map[string]interface{}{
					"A": "A",
				},
			},
			wantHasChanges: true,
			wantPatch:      []byte("{\"A\":\"A\"}"),
		},
		{
			name: "Nested field only in modified, align to modified",
			original: &unstructured.Unstructured{ // current
				Object: map[string]interface{}{},
			},
			modified: &unstructured.Unstructured{ // desired
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"A": "A",
							},
						},
					},
				},
			},
			wantHasChanges: true,
			wantPatch:      []byte("{\"spec\":{\"template\":{\"spec\":{\"A\":\"A\"}}}}"),
		},

		// Field only in original (not existing in modified) --> preserve original

		{
			name: "Field only in original, align to modified",
			original: &unstructured.Unstructured{ // current
				Object: map[string]interface{}{
					"A": "A",
				},
			},
			modified: &unstructured.Unstructured{ // desired
				Object: map[string]interface{}{},
			},
			wantHasChanges: false,
			wantPatch:      []byte("{}"),
		},
		{
			name: "Nested field only in original, align to modified",
			original: &unstructured.Unstructured{ // current
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"A": "A",
							},
						},
					},
				},
			},
			modified: &unstructured.Unstructured{ // desired
				Object: map[string]interface{}{},
			},
			wantHasChanges: false,
			wantPatch:      []byte("{}"),
		},

		// More tests
		{
			name: "No changes",
			original: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"A": "A",
						"B": "B",
						"C": "C", // C only in modified
					},
				},
			},
			modified: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"A": "A",
						"B": "B",
					},
				},
			},
			wantHasChanges: false,
			wantPatch:      []byte("{}"),
		},
		{
			name: "Many changes",
			original: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"A": "A",
						// B missing
						"C": "C", // C only in modified
					},
				},
			},
			modified: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"A": "A",
						"B": "B",
					},
				},
			},
			wantHasChanges: true,
			wantPatch:      []byte("{\"spec\":{\"B\":\"B\"}}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			patch, err := newMergePatchHelper(tt.original, tt.modified, nil)
			g.Expect(err).ToNot(HaveOccurred())

			g.Expect(patch.HasChanges()).To(Equal(tt.wantHasChanges))
			g.Expect(patch.data).To(Equal(tt.wantPatch))
		})
	}
}
