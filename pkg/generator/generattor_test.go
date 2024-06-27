package generator

import (
	"encoding/json"
	"testing"

	tp "github.com/crossplane-contrib/x-generation/pkg/types"
)

func Test_generateOverrideFields(t *testing.T) {
	tests := []struct {
		name          string
		base          map[string]interface{}
		overridePaths []tp.OverrideField
		want          string
	}{
		{
			name: "Should render path",
			base: map[string]interface{}{},
			overridePaths: []tp.OverrideField{
				{
					Path:  "spec.forProvider.certificateAuthorityConfiguration",
					Value: "testA",
				},
			},
			want: `{"spec": {"forProvider": {"certificateAuthorityConfiguration":"testA"}}}`,
		},
		{
			name: "Should not override existing path",
			base: map[string]interface{}{
				"spec": map[string]interface{}{
					"forProvider": map[string]interface{}{
						"test": "testValue",
					},
				},
			},
			overridePaths: []tp.OverrideField{
				{
					Path:  "spec.forProvider.certificateAuthorityConfiguration",
					Value: "testA",
				},
			},
			want: `{"spec": {"forProvider": {"test":"testValue", "certificateAuthorityConfiguration":"testA"}}}`,
		},
		{
			name: "Should create arrays of strings",
			base: map[string]interface{}{},
			overridePaths: []tp.OverrideField{
				{
					Path:  "spec.forProvider.certificateAuthorityConfiguration[0]",
					Value: "testA",
				},
			},
			want: `{"spec": {"forProvider": { "certificateAuthorityConfiguration": ["testA"]}}}`,
		},
		{
			name: "Should create arrays of arrays",
			base: map[string]interface{}{},
			overridePaths: []tp.OverrideField{
				{
					Path:  "spec.forProvider.certificateAuthorityConfiguration[0].subproperty[0]",
					Value: "testA",
				},
			},
			want: `{"spec": {"forProvider": { "certificateAuthorityConfiguration": [{"subproperty":["testA"]}]}}}`,
		},
		{
			name: "Should create arrays of arrays with properies",
			base: map[string]interface{}{},
			overridePaths: []tp.OverrideField{
				{
					Path:  "spec.forProvider.certificateAuthorityConfiguration[0].subproperty[0].arg",
					Value: "testA",
				},
			},
			want: `{"spec": {"forProvider": { "certificateAuthorityConfiguration": [{"subproperty":[{"arg":"testA"}]}]}}}`,
		},
		{
			name: "Should render path",
			base: map[string]interface{}{},
			overridePaths: []tp.OverrideField{
				{
					Path:  "spec.forProvider[\"certificateAuthorityConfiguration\"]",
					Value: "testA",
				},
			},
			want: `{"spec": {"forProvider": {"certificateAuthorityConfiguration":"testA"}}}`,
		},
		{
			name: "Should render path",
			base: map[string]interface{}{},
			overridePaths: []tp.OverrideField{
				{
					Path:  "spec.forProvider[\"certificateAuthority.Configuration\"]",
					Value: "testA",
				},
			},
			want: `{"spec": {"forProvider": {"certificateAuthority.Configuration":"testA"}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyOverrideFields(tt.base, tt.overridePaths)
			object, err := json.Marshal(got)
			if err != nil {
				t.Errorf("error marshalling object %v", err)
			}

			// unmarshal and marshal wanted value to be independent of formatting and oder of properties
			var wantobject interface{}
			err = json.Unmarshal([]byte(tt.want), &wantobject)
			if err != nil {
				t.Errorf("want is not valid json %v", err)
			}

			wantbyte, err := json.Marshal(wantobject)

			if err != nil {
				t.Errorf("error marshalling object %v", err)
			}

			gotString := string(object)
			wantString := string(wantbyte)
			if gotString != wantString {
				t.Errorf("Expaced object does not match got\n%v, want\n%v", gotString, wantString)
			}
		})
	}
}
