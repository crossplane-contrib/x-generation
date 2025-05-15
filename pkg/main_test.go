package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	xtype "github.com/crossplane-contrib/x-generation/pkg/types"
	cv1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func Test_tryToGetTags(t *testing.T) {
	type args struct {
		crd     extv1.CustomResourceDefinition
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    *extv1.JSONSchemaProps
		want1   string
		wantErr bool
	}{
		{
			name: "Should have tags",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tags": {
																Type: "array",
																Items: &extv1.JSONSchemaPropsOrArray{
																	Schema: &extv1.JSONSchemaProps{
																		Type: "object",
																		Properties: map[string]extv1.JSONSchemaProps{
																			"key": {
																				Type: "string",
																			},
																			"value": {
																				Type: "string",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want: &extv1.JSONSchemaProps{
				Type: "array",
				Items: &extv1.JSONSchemaPropsOrArray{
					Schema: &extv1.JSONSchemaProps{
						Type: "object",
						Properties: map[string]extv1.JSONSchemaProps{
							"key": {
								Type: "string",
							},
							"value": {
								Type: "string",
							},
						},
					},
				},
			},
			want1:   "spec.forProvider.tags",
			wantErr: false,
		},
		{
			name: "Should have find the right version",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1beta1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tags": {
																Type: "object",
																AdditionalProperties: &extv1.JSONSchemaPropsOrBool{
																	Schema: &extv1.JSONSchemaProps{
																		Type: "string",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tags": {
																Type: "array",
																Items: &extv1.JSONSchemaPropsOrArray{
																	Schema: &extv1.JSONSchemaProps{
																		Type: "object",
																		Properties: map[string]extv1.JSONSchemaProps{
																			"key": {
																				Type: "string",
																			},
																			"value": {
																				Type: "string",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want: &extv1.JSONSchemaProps{
				Type: "array",
				Items: &extv1.JSONSchemaPropsOrArray{
					Schema: &extv1.JSONSchemaProps{
						Type: "object",
						Properties: map[string]extv1.JSONSchemaProps{
							"key": {
								Type: "string",
							},
							"value": {
								Type: "string",
							},
						},
					},
				},
			},
			want1:   "spec.forProvider.tags",
			wantErr: false,
		},
		{
			name: "Should find tagging",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tagging": {
																Properties: map[string]extv1.JSONSchemaProps{
																	"tagSet": {
																		Type: "array",
																		Items: &extv1.JSONSchemaPropsOrArray{
																			Schema: &extv1.JSONSchemaProps{
																				Type: "object",
																				Properties: map[string]extv1.JSONSchemaProps{
																					"keyA": {
																						Type: "string",
																					},
																					"valueA": {
																						Type: "string",
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want: &extv1.JSONSchemaProps{
				Type: "array",
				Items: &extv1.JSONSchemaPropsOrArray{
					Schema: &extv1.JSONSchemaProps{
						Type: "object",
						Properties: map[string]extv1.JSONSchemaProps{
							"keyA": {
								Type: "string",
							},
							"valueA": {
								Type: "string",
							},
						},
					},
				},
			},
			want1:   "spec.forProvider.tagging.tagSet",
			wantErr: false,
		},
		{
			name: "Should find no tags",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tagging": {
																Properties: map[string]extv1.JSONSchemaProps{},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want:    nil,
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tryToGetTags(tt.args.crd, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("tryToGetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tryToGetTags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("tryToGetTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_checkTagType(t *testing.T) {
	type args struct {
		crd     extv1.CustomResourceDefinition
		version string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Should find keyValueArray",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tags": {
																Type: "array",
																Items: &extv1.JSONSchemaPropsOrArray{
																	Schema: &extv1.JSONSchemaProps{
																		Type: "object",
																		Properties: map[string]extv1.JSONSchemaProps{
																			"key": {
																				Type: "string",
																			},
																			"value": {
																				Type: "string",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want:  "keyValueArray",
			want1: "spec.forProvider.tags",
		},
		{
			name: "Should find tagKeyTagValueArray",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tags": {
																Type: "array",
																Items: &extv1.JSONSchemaPropsOrArray{
																	Schema: &extv1.JSONSchemaProps{
																		Type: "object",
																		Properties: map[string]extv1.JSONSchemaProps{
																			"tagKey": {
																				Type: "string",
																			},
																			"tagValue": {
																				Type: "string",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want:  "tagKeyTagValueArray",
			want1: "spec.forProvider.tags",
		},
		{
			name: "Should find tagObject",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tags": {
																Type: "object",
																AdditionalProperties: &extv1.JSONSchemaPropsOrBool{
																	Schema: &extv1.JSONSchemaProps{
																		Type: "string",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want:  "tagObject",
			want1: "spec.forProvider.tags",
		},
		{
			name: "Should find no tags if no tags",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want:  "noTag",
			want1: "",
		},
		{
			name: "Should find no tags if wrong types array",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tags": {
																Type: "array",
																Items: &extv1.JSONSchemaPropsOrArray{
																	Schema: &extv1.JSONSchemaProps{
																		Type: "object",
																		Properties: map[string]extv1.JSONSchemaProps{
																			"key1": {
																				Type: "string",
																			},
																			"value": {
																				Type: "string",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want:  "noTag",
			want1: "",
		},
		{
			name: "Should find no tags if wrong types object",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "v1alpha1",
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Properties: map[string]extv1.JSONSchemaProps{
															"tags": {
																Type: "object",
																AdditionalProperties: &extv1.JSONSchemaPropsOrBool{
																	Schema: &extv1.JSONSchemaProps{
																		Type: "int",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want:  "noTag",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := checkTagType(tt.args.crd, tt.args.version)
			if got != tt.want {
				t.Errorf("checkTagType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("checkTagType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getTagListAsString(t *testing.T) {
	type args struct {
		g *Generator
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should generate list",
			args: args{
				g: &Generator{
					Tags: xtype.LocalTagConfig{
						TagConfig: xtype.TagConfig{
							FromLabels: []string{
								"tagA",
								"tagB",
							},
						},
					},
				},
			},
			want: "[\"tagA\",\"tagB\"]",
		},
		{
			name: "Should generate empty list",
			args: args{
				g: &Generator{},
			},
			want: "[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTagListAsString(tt.args.g); got != tt.want {
				t.Errorf("getTagListAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCommonTagsAsString(t *testing.T) {
	type args struct {
		g *Generator
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should generate object",
			args: args{

				g: &Generator{
					Tags: xtype.LocalTagConfig{
						TagConfig: xtype.TagConfig{
							Common: map[string]string{
								"commonA": "valueA",
								"commonB": "valueB",
							},
						},
					},
				},
			},
			want: "{\"commonA\":\"valueA\",\"commonB\":\"valueB\"}",
		},
		{
			name: "Should generate empty object",
			args: args{
				g: &Generator{},
			},
			want: "{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCommonTagsAsString(tt.args.g); got != tt.want {
				t.Errorf("getCommonTagsAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLabelListAsString(t *testing.T) {
	type args struct {
		g *Generator
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should generate list",
			args: args{
				g: &Generator{
					Labels: xtype.LocalLabelConfig{
						LabelConfig: xtype.LabelConfig{
							FromCRD: []string{
								"labelA",
								"labelB",
							},
						},
					},
				},
			},
			want: "[\"labelA\",\"labelB\"]",
		},
		{
			name: "Should generate empty list",
			args: args{
				g: &Generator{},
			},
			want: "[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLabelListAsString(tt.args.g); got != tt.want {
				t.Errorf("getLabelListAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCommonLabelsString(t *testing.T) {
	type args struct {
		g *Generator
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should generate object",
			args: args{

				g: &Generator{
					Labels: xtype.LocalLabelConfig{
						LabelConfig: xtype.LabelConfig{
							Common: map[string]string{
								"commonA": "valueA",
								"commonB": "valueB",
							},
						},
					},
				},
			},
			want: "{\"commonA\":\"valueA\",\"commonB\":\"valueB\"}",
		},
		{
			name: "Should generate empty object",
			args: args{
				g: &Generator{},
			},
			want: "{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCommonLabelsString(tt.args.g); got != tt.want {
				t.Errorf("getCommonLabelsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendLists(t *testing.T) {
	type args struct {
		a *[]string
		b *[]string
	}
	tests := []struct {
		name string
		args args
		want *[]string
	}{
		{
			name: "Should append independent list",
			args: args{
				a: &[]string{
					"a1",
					"a2",
					"a3",
				},
				b: &[]string{
					"b1",
					"b2",
					"b3",
				},
			},
			want: &[]string{
				"a1",
				"a2",
				"a3",
				"b1",
				"b2",
				"b3",
			},
		},
		{
			name: "Should append overlapping list",
			args: args{
				a: &[]string{
					"a1",
					"a2",
					"a3",
				},
				b: &[]string{
					"a1",
					"a2",
					"b3",
					"b4",
					"b5",
				},
			},
			want: &[]string{
				"a1",
				"a2",
				"a3",
				"b3",
				"b4",
				"b5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendLists(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendLists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendStringMaps(t *testing.T) {
	type args struct {
		a map[string]string
		b map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Should appmend indepentent maps",
			args: args{
				a: map[string]string{
					"keyAA": "valueAA",
					"keyAB": "valueAB",
					"keyAC": "valueAC",
				},
				b: map[string]string{
					"keyBA": "valueBA",
					"keyBB": "valueBB",
					"keyBC": "valueBC",
				},
			},
			want: map[string]string{
				"keyAA": "valueAA",
				"keyAB": "valueAB",
				"keyAC": "valueAC",
				"keyBA": "valueBA",
				"keyBB": "valueBB",
				"keyBC": "valueBC",
			},
		},
		{
			name: "Should appmend overlapping maps",
			args: args{
				a: map[string]string{
					"keyAA": "valueAA",
					"keyAB": "valueAB",
					"keyAC": "valueAC",
				},
				b: map[string]string{
					"keyAA": "valueBA",
					"keyAB": "valueBB",
					"keyBC": "valueBC",
					"keyBD": "valueBD",
					"keyBE": "valueBE",
				},
			},
			want: map[string]string{
				"keyAC": "valueAC",
				"keyAA": "valueBA",
				"keyAB": "valueBB",
				"keyBC": "valueBC",
				"keyBD": "valueBD",
				"keyBE": "valueBE",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendStringMaps(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendStringMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_CheckConfig(t *testing.T) {
	type fields struct {
		Group                string
		Name                 string
		Plural               *string
		Version              string
		ScriptFileName       *string
		ConnectionSecretKeys *[]string
		Ignore               bool
		PatchExternalName    *bool
		UIDFieldPath         *string
		OverrideFields       []xtype.OverrideField
		Compositions         []xtype.Composition
		Tags                 xtype.LocalTagConfig
		Labels               xtype.LocalLabelConfig
		Provider             xtype.ProviderConfig
		crdSource            string
		configPath           string
		tagType              string
		tagProperty          string
	}
	type args struct {
		generatorConfig *xtype.GeneratorConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should validate config in commonLables",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"commonA",
							"commonB",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"commonA": "valueA",
							"commonB": "valueB",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should validate config in fromCrd",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdA",
							"fromCrdB",
							"fromCrdC",
						},
					},
				},
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdA",
							"fromCrdB",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{},
			},
			wantErr: false,
		},
		{
			name: "Should validate config in globalLabels",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"crossplane.io/claim-name",
							"crossplane.io/claim-namespace",
							"crossplane.io/composite",
							"external-name",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{},
			},
			wantErr: false,
		},
		{
			name: "Should validate config in all places",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdA",
							"fromCrdB",
							"fromCrdC",
						},
					},
				},
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdA",
							"fromCrdB",
							"commonA",
							"commonB",
							"crossplane.io/claim-name",
							"crossplane.io/claim-namespace",
							"crossplane.io/composite",
							"external-name",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"commonA": "valueA",
							"commonB": "valueB",
							"commonC": "valueC",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should have errors from commonLables",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"commonA",
							"commonX",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"commonA": "valueA",
							"commonB": "valueB",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should have errors from fromCrd",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdA",
							"fromCrdB",
							"fromCrdC",
						},
					},
				},
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdA",
							"fromCrdX",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{},
			},
			wantErr: true,
		},
		{
			name: "Should have errors from globalLabels",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"crossplane.io/claim-name",
							"crossplane.io/claim-namespace",
							"crossplane.io/composite",
							"external-nameX",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{},
			},
			wantErr: true,
		},
		{
			name: "Should have errors from all",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdA",
							"fromCrdB",
							"fromCrdC",
						},
					},
				},
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdA",
							"fromCrdX",
							"commonA",
							"commonX",
							"crossplane.io/claim-name",
							"crossplane.io/claim-namespace",
							"crossplane.io/composite",
							"external-name",
							"external-nameX",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"commonA": "valueA",
							"commonB": "valueB",
							"commonC": "valueC",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Group:                tt.fields.Group,
				Name:                 tt.fields.Name,
				Plural:               tt.fields.Plural,
				Version:              tt.fields.Version,
				ScriptFileName:       tt.fields.ScriptFileName,
				ConnectionSecretKeys: tt.fields.ConnectionSecretKeys,
				Ignore:               tt.fields.Ignore,
				PatchExternalName:    tt.fields.PatchExternalName,
				UIDFieldPath:         tt.fields.UIDFieldPath,
				OverrideFields:       tt.fields.OverrideFields,
				Compositions:         tt.fields.Compositions,
				Tags:                 tt.fields.Tags,
				Labels:               tt.fields.Labels,
				Provider:             tt.fields.Provider,
				crdSource:            tt.fields.crdSource,
				configPath:           tt.fields.configPath,
				TagType:              &tt.fields.tagType,
				TagProperty:          &tt.fields.tagProperty,
			}
			if err := g.CheckConfig(tt.args.generatorConfig); (err != nil) != tt.wantErr {
				t.Errorf("Generator.CheckConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_UpdateConfig(t *testing.T) {
	type fields struct {
		Group                string
		Name                 string
		Plural               *string
		Version              string
		ScriptFileName       *string
		ConnectionSecretKeys *[]string
		Ignore               bool
		PatchExternalName    *bool
		UIDFieldPath         *string
		OverrideFields       []xtype.OverrideField
		Compositions         []xtype.Composition
		Tags                 xtype.LocalTagConfig
		Labels               xtype.LocalLabelConfig
		Provider             xtype.ProviderConfig
		crdSource            string
		configPath           string
		tagType              string
		tagProperty          string
	}
	type args struct {
		generatorConfig *xtype.GeneratorConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "Should not append labels",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: replaceGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: replaceGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
		},
		{
			name: "Should append labels",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
		},
		{
			name: "Should not append labels empty",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: replaceGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: replaceGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{},
					},
				},
			},
		},
		{
			name: "Should append labels empty no GlobalHandling",
			fields: fields{
				Labels: xtype.LocalLabelConfig{

					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{

					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
		},
		{
			name: "Should append labels empty",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
		},
		{
			name: "Should not append common labels",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						Common: replaceGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdLC": "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdGC": "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						Common: replaceGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdLC": "valueLC",
						},
					},
				},
			},
		},
		{
			name: "Should append common labels",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						Common: appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						Common: appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
						},
					},
				},
			},
		},
		{
			name: "Should not append common labels empty",

			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						Common: replaceGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						Common: replaceGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{},
					},
				},
			},
		},
		{
			name: "Should append labels empty no GlobalHandling",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{

					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
		},
		{
			name: "Should append common labels empty",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						Common: appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						Common: appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
		},
		{
			name: "Should not append tags",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: replaceGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{

				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: replaceGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
		},
		{
			name: "Should append tags",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
		},
		{
			name: "Should not append tags empty",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: replaceGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: replaceGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{},
					},
				},
			},
		},
		{
			name: "Should append tags empty no GlobalHandling",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
		},
		{
			name: "Should append tags empty",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
		},
		{
			name: "Should not append common tags",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						Common: replaceGlobal,
					},
					TagConfig: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdLC": "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdGC": "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						Common: replaceGlobal,
					},
					TagConfig: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdLC": "valueLC",
						},
					},
				},
			},
		},
		{
			name: "Should append common tags",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						Common: appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						Common: appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
						},
					},
				},
			},
		},
		{
			name: "Should not append common tags empty",

			fields: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						Common: replaceGlobal,
					},
					TagConfig: xtype.TagConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						Common: replaceGlobal,
					},
					TagConfig: xtype.TagConfig{
						Common: map[string]string{},
					},
				},
			},
		},
		{
			name: "Should append tags empty no GlobalHandling",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					TagConfig: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
		},
		{
			name: "Should append common labels empty",
			fields: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						Common: appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Tags: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						Common: appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
		},
		{
			name: "Should append all",
			fields: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: appendGlobal,
						Common:  appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
						},
					},
				},
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: appendGlobal,
						Common:     appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
					Tags: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: xtype.LocalLabelConfig{
					GlobalHandling: xtype.GlobalHandlingLabels{
						FromCRD: appendGlobal,
						Common:  appendGlobal,
					},
					LabelConfig: xtype.LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
						},
					},
				},
				Tags: xtype.LocalTagConfig{
					GlobalHandling: xtype.GlobalHandlingTags{
						FromLabels: appendGlobal,
						Common:     appendGlobal,
					},
					TagConfig: xtype.TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Group:                tt.fields.Group,
				Name:                 tt.fields.Name,
				Plural:               tt.fields.Plural,
				Version:              tt.fields.Version,
				ScriptFileName:       tt.fields.ScriptFileName,
				ConnectionSecretKeys: tt.fields.ConnectionSecretKeys,
				Ignore:               tt.fields.Ignore,
				PatchExternalName:    tt.fields.PatchExternalName,
				UIDFieldPath:         tt.fields.UIDFieldPath,
				OverrideFields:       tt.fields.OverrideFields,
				Compositions:         tt.fields.Compositions,
				Tags:                 tt.fields.Tags,
				Labels:               tt.fields.Labels,
				Provider:             tt.fields.Provider,
				crdSource:            tt.fields.crdSource,
				configPath:           tt.fields.configPath,
				TagType:              &tt.fields.tagType,
				TagProperty:          &tt.fields.tagProperty,
			}
			g.UpdateConfig(tt.args.generatorConfig)
			marshaledWantTags, _ := json.Marshal(tt.want.Tags)
			wantTagsString := string(marshaledWantTags)
			marshaledIsTags, _ := json.Marshal(g.Tags)
			isTagsString := string(marshaledIsTags)

			if wantTagsString != isTagsString {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", isTagsString, wantTagsString)
			}
			marshaledWantLabels, _ := json.Marshal(tt.want.Labels)
			wantLabelsString := string(marshaledWantLabels)
			marshaledIsLabels, _ := json.Marshal(g.Labels)
			isLabelsString := string(marshaledIsLabels)

			if wantLabelsString != isLabelsString {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", isLabelsString, wantLabelsString)
			}
			if tt.want.Group != g.Group {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.Group, tt.want.Group)
			}
			if tt.want.Name != g.Name {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.Name, tt.want.Name)
			}
			if tt.want.Plural != g.Plural {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.Plural, tt.want.Plural)
			}
			if tt.want.Version != g.Version {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.Version, tt.want.Version)
			}
			if tt.want.ScriptFileName != g.ScriptFileName {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.ScriptFileName, tt.want.ScriptFileName)
			}
			if tt.want.ConnectionSecretKeys != g.ConnectionSecretKeys {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.ConnectionSecretKeys, tt.want.ConnectionSecretKeys)
			}
			if tt.want.Ignore != g.Ignore {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.Ignore, tt.want.Ignore)
			}
			if tt.want.PatchExternalName != g.PatchExternalName {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.PatchExternalName, tt.want.PatchExternalName)
			}
			if tt.want.Provider != g.Provider {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.Provider, tt.want.Provider)
			}
			if tt.want.crdSource != g.crdSource {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.crdSource, tt.want.crdSource)
			}
			if tt.want.configPath != g.configPath {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.configPath, tt.want.configPath)
			}
			if tt.want.tagType != *g.TagType {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.TagType, tt.want.tagType)
			}
			if tt.want.tagProperty != *g.TagProperty {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.TagProperty, tt.want.tagProperty)
			}
		})
	}
}

func Test_listHas(t *testing.T) {
	type args struct {
		list  *[]string
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should have",
			args: args{
				list: &[]string{
					"A",
					"B",
					"C",
				},
				value: "B",
			},
			want: true,
		},
		{
			name: "Should not have",
			args: args{
				list: &[]string{
					"A",
					"B",
					"C",
				},
				value: "D",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listHas(tt.args.list, tt.args.value); got != tt.want {
				t.Errorf("listHas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkConfig(t *testing.T) {
	type args struct {
		generatorConfig *xtype.GeneratorConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should be valid common",
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"commonA": "valueCA",
							"commonB": "valueCB",
							"commonC": "valueCB",
						},
					},
					Tags: xtype.TagConfig{
						FromLabels: []string{
							"commonA",
							"commonB",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should be valid globalLabels",
			args: args{
				generatorConfig: &xtype.GeneratorConfig{

					Tags: xtype.TagConfig{
						FromLabels: []string{
							"crossplane.io/claim-name",
							"crossplane.io/claim-namespace",
							"crossplane.io/composite",
							"external-name",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should be not be valid",
			args: args{
				generatorConfig: &xtype.GeneratorConfig{
					Labels: xtype.LabelConfig{
						Common: map[string]string{
							"commonA": "valueCA",
							"commonB": "valueCB",
							"commonC": "valueCB",
						},
					},
					Tags: xtype.TagConfig{
						FromLabels: []string{
							"commonA",
							"commonX",
							"external-nameX",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should be valid vor nil ",
			args: args{
				generatorConfig: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkConfig(tt.args.generatorConfig); (err != nil) != tt.wantErr {
				t.Errorf("checkConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tryProperties(t *testing.T) {
	type args struct {
		crd     extv1.CustomResourceDefinition
		version string
	}
	tests := []struct {
		name string
		args args
		want *extv1.JSONSchemaProps
	}{
		{
			name: "Should have default property",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "testv1",
								AdditionalPrinterColumns: []extv1.CustomResourceColumnDefinition{{
									JSONPath: ".metadata.annotations.crossplane.io/external-name",
									Name:     "EXTERNAL-NAME",
									Type:     "string",
								}},
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Description: "For Provider property",
														Properties: map[string]extv1.JSONSchemaProps{
															"default": {
																Description: "This is a property called properties",
																Type:        "array",
																Items: &extv1.JSONSchemaPropsOrArray{
																	Schema: &extv1.JSONSchemaProps{
																		Type:        "object",
																		Description: "This is an item of the property properties",
																		Properties: map[string]extv1.JSONSchemaProps{
																			"arrayProp": {
																				Type: "string",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
											"status": {
												Properties: map[string]extv1.JSONSchemaProps{
													"name": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want: &extv1.JSONSchemaProps{
				Description: "For Provider property",
				Properties: map[string]extv1.JSONSchemaProps{
					"default": {
						Description: "This is a property called properties",
						Type:        "array",
						Items: &extv1.JSONSchemaPropsOrArray{
							Schema: &extv1.JSONSchemaProps{
								Type:        "object",
								Description: "This is an item of the property properties",
								Properties: map[string]extv1.JSONSchemaProps{
									"arrayProp": {
										Type: "string",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Should have properties property",
			args: args{
				crd: extv1.CustomResourceDefinition{
					Spec: extv1.CustomResourceDefinitionSpec{
						Versions: []extv1.CustomResourceDefinitionVersion{
							{
								Name: "testv1",
								AdditionalPrinterColumns: []extv1.CustomResourceColumnDefinition{{
									JSONPath: ".metadata.annotations.crossplane.io/external-name",
									Name:     "EXTERNAL-NAME",
									Type:     "string",
								}},
								Schema: &extv1.CustomResourceValidation{
									OpenAPIV3Schema: &extv1.JSONSchemaProps{
										Properties: map[string]extv1.JSONSchemaProps{
											"spec": {
												Properties: map[string]extv1.JSONSchemaProps{
													"forProvider": {
														Description: "For Provider property",
														Properties: map[string]extv1.JSONSchemaProps{
															"properties": {
																Description: "This is a property called properties",
																Type:        "array",
																Items: &extv1.JSONSchemaPropsOrArray{
																	Schema: &extv1.JSONSchemaProps{
																		Type:        "object",
																		Description: "This is an item of the property properties",
																		Properties: map[string]extv1.JSONSchemaProps{
																			"arrayProp": {
																				Type: "string",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
											"status": {
												Properties: map[string]extv1.JSONSchemaProps{
													"name": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				version: "v1alpha1",
			},
			want: &extv1.JSONSchemaProps{
				Description: "For Provider property",
				Properties: map[string]extv1.JSONSchemaProps{
					"properties": {
						Description: "This is a property called properties",
						Type:        "array",
						Items: &extv1.JSONSchemaPropsOrArray{
							Schema: &extv1.JSONSchemaProps{
								Type:        "object",
								Description: "This is an item of the property properties",
								Properties: map[string]extv1.JSONSchemaProps{
									"arrayProp": {
										Type: "string",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			crdSource, err := json.Marshal(tt.args.crd)
			if err != nil {
				t.Errorf("could not marshal crdSource")
			}
			tempDir, err := os.MkdirTemp("", "g-generation*")
			if err != nil {
				t.Errorf("could not generate tempDir")
			}
			plural := "TestObjects"
			g := Generator{
				Group:       "example.cloud",
				Name:        "TestObject",
				Version:     "testv1",
				crdSource:   string(crdSource),
				configPath:  tempDir,
				TagType:     nil,
				TagProperty: nil,
				Provider: xtype.ProviderConfig{
					CRD: xtype.CrdConfig{
						Version: "testv1",
					},
				},
				Plural:                &plural,
				Compositions:          []xtype.Composition{},
				OverrideFields:        []xtype.OverrideField{},
				OverrideFieldsInClaim: []xtype.OverrideFieldInClaim{},
			}

			gConfig := xtype.GeneratorConfig{
				CompositionIdentifier: "example.cloud",
			}

			cwd, _ := os.Getwd()

			sp := filepath.Join(cwd, "functions")
			g.Exec(&gConfig, sp, "", "")

			path := filepath.Join(tempDir, "definition.yaml")
			y, err := os.ReadFile(path)

			if err != nil {
				t.Errorf("could not load definition.yaml file")
				return
			}

			var newCRD cv1.CompositeResourceDefinition
			err = yaml.Unmarshal(y, &newCRD)

			if err != nil {
				t.Errorf("could not parse definition.yaml file")
				return
			}

			openApiSchema := newCRD.Spec.Versions[0].Schema.OpenAPIV3Schema

			var properties extv1.JSONSchemaProps

			err = yaml.Unmarshal(openApiSchema.Raw, &properties)
			if err != nil {
				t.Errorf("could not unmarshal properties")
				return
			}

			if spec, ok := properties.Properties["spec"]; ok {
				if forProvider, ok := spec.Properties["forProvider"]; ok {
					sourceJSON, err := json.Marshal(forProvider)
					if err != nil {
						t.Error("source cannot be marshaled")
					}
					wantJSON, err := json.Marshal(tt.want)
					if err != nil {
						t.Error("want cannot be marshaled")
					}
					patchJSON, err := jsonpatch.CreateMergePatch(wantJSON, sourceJSON)
					if err != nil {
						t.Error("error generating patch")
					}
					if string(patchJSON) != "{}" {
						t.Errorf("objects not the same: %s", patchJSON)
						return
					}
				}
			}
			err = os.Remove(path)
			if err != nil {
				t.Error("could not delete definition file")
			}
			err = os.Remove(tempDir)
			if err != nil {
				t.Error("could not delete temp directory")
			}
		})
	}
}
func Test_noDefaultPatch(t *testing.T) {
	type args struct {
		crd     extv1.CustomResourceDefinition
		version string
	}
	testData := struct {
		name string
		args args
	}{

		name: "Should have properties property",
		args: args{
			crd: extv1.CustomResourceDefinition{
				Spec: extv1.CustomResourceDefinitionSpec{
					Versions: []extv1.CustomResourceDefinitionVersion{
						{
							Name: "testv1",
							AdditionalPrinterColumns: []extv1.CustomResourceColumnDefinition{{
								JSONPath: ".metadata.annotations.crossplane.io/external-name",
								Name:     "EXTERNAL-NAME",
								Type:     "string",
							}},
							Schema: &extv1.CustomResourceValidation{
								OpenAPIV3Schema: &extv1.JSONSchemaProps{
									Properties: map[string]extv1.JSONSchemaProps{
										"spec": {
											Properties: map[string]extv1.JSONSchemaProps{
												"providerConfigRef": {
													Default: &extv1.JSON{
														Raw: []byte("{\"name\": \"default\"}"),
													},
													Description: "ProviderConfigReference specifies how the provider that will be used to create, observe, update, and delete this managed",
													Type:        "object",
													Properties: map[string]extv1.JSONSchemaProps{
														"name": {
															Type:        "string",
															Description: "Name of the referenced object.",
														},
														"policy": {
															Type:        "string",
															Description: "description: Policies for referencing.",
														},
													},
												},
											},
										},
										"status": {
											Properties: map[string]extv1.JSONSchemaProps{
												"name": {
													Type: "string",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			version: "v1alpha1",
		},
	}
	t.Run(testData.name, func(t *testing.T) {

		crdSource, err := json.Marshal(testData.args.crd)
		if err != nil {
			t.Errorf("could not marshal crdSource")
		}
		tempDir, err := os.MkdirTemp("", "g-generation*")
		if err != nil {
			t.Errorf("could not generate tempDir")
		}
		plural := "TestObjects"
		g := Generator{
			Group:       "example.cloud",
			Name:        "TestObject",
			Version:     "testv1",
			crdSource:   string(crdSource),
			configPath:  tempDir,
			TagType:     nil,
			TagProperty: nil,
			Provider: xtype.ProviderConfig{
				CRD: xtype.CrdConfig{
					Version: "testv1",
				},
			},
			Plural: &plural,
			Compositions: []xtype.Composition{
				{
					Name:     "configuration",
					Provider: "sop",
					Default:  true,
				},
			},
			OverrideFields:        []xtype.OverrideField{},
			OverrideFieldsInClaim: []xtype.OverrideFieldInClaim{},
		}

		gConfig := xtype.GeneratorConfig{
			CompositionIdentifier: "example.cloud",
		}

		cwd, _ := os.Getwd()

		sp := filepath.Join(cwd, "functions")
		g.Exec(&gConfig, sp, "", "")

		path := filepath.Join(tempDir, "composition-configuration.yaml")
		y, err := os.ReadFile(path)

		if err != nil {
			t.Errorf("could not load definition.yaml file")
		}

		var newComposition cv1.Composition
		err = yaml.Unmarshal(y, &newComposition)

		if err != nil {
			t.Errorf("could not parse definition.yaml file")
		}

		patchsets := newComposition.Spec.PatchSets

		var parameterPatchSet *cv1.PatchSet

		for _, patchset := range patchsets {
			if patchset.Name == "Parameters" {
				parameterPatchSet = &patchset
				break
			}
		}

		if parameterPatchSet != nil {
			for _, patch := range parameterPatchSet.Patches {
				if *patch.FromFieldPath == "spec.providerConfigRef.default" {
					t.Error("There should be no patch for a default property")
					return
				}
			}
		}

		err = os.Remove(path)
		if err != nil {
			t.Error("could not delete definition file")
		}
		definitionPath := filepath.Join(tempDir, "definition.yaml")
		err = os.Remove(definitionPath)
		if err != nil {
			t.Error("could not delete composition file")
		}
		err = os.Remove(tempDir)
		if err != nil {
			t.Error("could not delete temp directory")
		}
	})
}

func Test_setDefaultCompositeDeletePolicy_Foreground(t *testing.T) {
	type args struct {
		crd     extv1.CustomResourceDefinition
		version string
	}
	testData := struct {
		name string
		args args
	}{

		name: "Should set default composite delete policy to Foreground",
		args: args{
			crd: extv1.CustomResourceDefinition{
				Spec: extv1.CustomResourceDefinitionSpec{
					Versions: []extv1.CustomResourceDefinitionVersion{
						{
							Name: "testv1",
							AdditionalPrinterColumns: []extv1.CustomResourceColumnDefinition{{
								JSONPath: ".metadata.annotations.crossplane.io/external-name",
								Name:     "EXTERNAL-NAME",
								Type:     "string",
							}},
							Schema: &extv1.CustomResourceValidation{
								OpenAPIV3Schema: &extv1.JSONSchemaProps{
									Properties: map[string]extv1.JSONSchemaProps{
										"spec": {
											Properties: map[string]extv1.JSONSchemaProps{
												"providerConfigRef": {
													Default: &extv1.JSON{
														Raw: []byte("{\"name\": \"default\"}"),
													},
													Description: "ProviderConfigReference specifies how the provider that will be used to create, observe, update, and delete this managed",
													Type:        "object",
													Properties: map[string]extv1.JSONSchemaProps{
														"name": {
															Type:        "string",
															Description: "Name of the referenced object.",
														},
														"policy": {
															Type:        "string",
															Description: "description: Policies for referencing.",
														},
													},
												},
											},
										},
										"status": {
											Properties: map[string]extv1.JSONSchemaProps{
												"name": {
													Type: "string",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			version: "v1alpha1",
		},
	}
	t.Run(testData.name, func(t *testing.T) {

		crdSource, err := json.Marshal(testData.args.crd)
		if err != nil {
			t.Errorf("could not marshal crdSource")
		}
		tempDir, err := os.MkdirTemp("", "g-generation*")
		if err != nil {
			t.Errorf("could not generate tempDir")
		}
		plural := "TestObjects"
		g := Generator{
			Group:       "example.cloud",
			Name:        "TestObject",
			Version:     "testv1",
			crdSource:   string(crdSource),
			configPath:  tempDir,
			TagType:     nil,
			TagProperty: nil,
			Provider: xtype.ProviderConfig{
				CRD: xtype.CrdConfig{
					Version: "testv1",
				},
			},
			Plural: &plural,
			Compositions: []xtype.Composition{
				{
					Name:     "configuration",
					Provider: "sop",
					Default:  true,
				},
			},
			OverrideFields:        []xtype.OverrideField{},
			OverrideFieldsInClaim: []xtype.OverrideFieldInClaim{},
			DefaultCompositeDeletePolicy: func() *string { s := "Foreground"; return &s }(),
		}

		gConfig := xtype.GeneratorConfig{
			CompositionIdentifier: "example.cloud",
		}

		cwd, _ := os.Getwd()

		sp := filepath.Join(cwd, "functions")
		g.Exec(&gConfig, sp, "", "")

		path := filepath.Join(tempDir, "definition.yaml")
		y, err := os.ReadFile(path)

		if err != nil {
			t.Errorf("could not load definition.yaml file")
		}

		var newCRD cv1.CompositeResourceDefinition
			err = yaml.Unmarshal(y, &newCRD)

		if err != nil {
			t.Errorf("could not parse definition.yaml file")
		}

		// assert that the default composite delete policy is set to Foreground
		if newCRD.Spec.DefaultCompositeDeletePolicy != nil {
			if *newCRD.Spec.DefaultCompositeDeletePolicy != "Foreground" {
				t.Errorf("expected default composite delete policy to be Foreground, got %s", *newCRD.Spec.DefaultCompositeDeletePolicy)
			}
		}

		err = os.Remove(path)
		if err != nil {
			t.Error("could not delete definition file")
		}
		definitionPath := filepath.Join(tempDir, "composition-configuration.yaml")
		err = os.Remove(definitionPath)
		if err != nil {
			t.Error("could not delete composition file")
		}
		err = os.Remove(tempDir)
		if err != nil {
			t.Error("could not delete temp directory")
		}
	})
}

func Test_setDefaultCompositeDeletePolicy_Background(t *testing.T) {
	type args struct {
		crd     extv1.CustomResourceDefinition
		version string
	}
	testData := struct {
		name string
		args args
	}{

		name: "Should set default composite delete policy to Background",
		args: args{
			crd: extv1.CustomResourceDefinition{
				Spec: extv1.CustomResourceDefinitionSpec{
					Versions: []extv1.CustomResourceDefinitionVersion{
						{
							Name: "testv1",
							AdditionalPrinterColumns: []extv1.CustomResourceColumnDefinition{{
								JSONPath: ".metadata.annotations.crossplane.io/external-name",
								Name:     "EXTERNAL-NAME",
								Type:     "string",
							}},
							Schema: &extv1.CustomResourceValidation{
								OpenAPIV3Schema: &extv1.JSONSchemaProps{
									Properties: map[string]extv1.JSONSchemaProps{
										"spec": {
											Properties: map[string]extv1.JSONSchemaProps{
												"providerConfigRef": {
													Default: &extv1.JSON{
														Raw: []byte("{\"name\": \"default\"}"),
													},
													Description: "ProviderConfigReference specifies how the provider that will be used to create, observe, update, and delete this managed",
													Type:        "object",
													Properties: map[string]extv1.JSONSchemaProps{
														"name": {
															Type:        "string",
															Description: "Name of the referenced object.",
														},
														"policy": {
															Type:        "string",
															Description: "description: Policies for referencing.",
														},
													},
												},
											},
										},
										"status": {
											Properties: map[string]extv1.JSONSchemaProps{
												"name": {
													Type: "string",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			version: "v1alpha1",
		},
	}
	t.Run(testData.name, func(t *testing.T) {

		crdSource, err := json.Marshal(testData.args.crd)
		if err != nil {
			t.Errorf("could not marshal crdSource")
		}
		tempDir, err := os.MkdirTemp("", "g-generation*")
		if err != nil {
			t.Errorf("could not generate tempDir")
		}
		plural := "TestObjects"
		g := Generator{
			Group:       "example.cloud",
			Name:        "TestObject",
			Version:     "testv1",
			crdSource:   string(crdSource),
			configPath:  tempDir,
			TagType:     nil,
			TagProperty: nil,
			Provider: xtype.ProviderConfig{
				CRD: xtype.CrdConfig{
					Version: "testv1",
				},
			},
			Plural: &plural,
			Compositions: []xtype.Composition{
				{
					Name:     "configuration",
					Provider: "sop",
					Default:  true,
				},
			},
			OverrideFields:        []xtype.OverrideField{},
			OverrideFieldsInClaim: []xtype.OverrideFieldInClaim{},
			DefaultCompositeDeletePolicy: func() *string { s := "Background"; return &s }(),
		}

		gConfig := xtype.GeneratorConfig{
			CompositionIdentifier: "example.cloud",
		}

		cwd, _ := os.Getwd()

		sp := filepath.Join(cwd, "functions")
		g.Exec(&gConfig, sp, "", "")

		path := filepath.Join(tempDir, "definition.yaml")
		y, err := os.ReadFile(path)

		if err != nil {
			t.Errorf("could not load definition.yaml file")
		}

		var newCRD cv1.CompositeResourceDefinition
			err = yaml.Unmarshal(y, &newCRD)

		if err != nil {
			t.Errorf("could not parse definition.yaml file")
		}

		// assert that the default composite delete policy is set to Background
		if newCRD.Spec.DefaultCompositeDeletePolicy != nil {
			if *newCRD.Spec.DefaultCompositeDeletePolicy != "Background" {
				t.Errorf("expected default composite delete policy to be Background, got %s", *newCRD.Spec.DefaultCompositeDeletePolicy)
			}
		}

		err = os.Remove(path)
		if err != nil {
			t.Error("could not delete definition file")
		}
		definitionPath := filepath.Join(tempDir, "composition-configuration.yaml")
		err = os.Remove(definitionPath)
		if err != nil {
			t.Error("could not delete composition file")
		}
		err = os.Remove(tempDir)
		if err != nil {
			t.Error("could not delete temp directory")
		}
	})
}
