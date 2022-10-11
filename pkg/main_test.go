package main

import (
	"encoding/json"
	"reflect"
	"testing"

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
			want1:   "tag",
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
			want1:   "tag",
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
			want1:   "tagSet",
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
			want1: "tag",
		},
		{
			name: "Should find tagKeyValueArray",
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
			want:  "tagKeyValueArray",
			want1: "tag",
		},
		{
			name: "Should find stringObject",
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
			want:  "stringObject",
			want1: "tag",
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
			want:  "",
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
			want:  "",
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
			want:  "",
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
					Tags: LocalTagConfig{
						TagConfig: TagConfig{
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
					Tags: LocalTagConfig{
						TagConfig: TagConfig{
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
					Labels: LocalLabelConfig{
						LabelConfig: LabelConfig{
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
					Labels: LocalLabelConfig{
						LabelConfig: LabelConfig{
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
		OverrideFields       []OverrideField
		Compositions         []Composition
		Tags                 LocalTagConfig
		Labels               LocalLabelConfig
		Provider             ProviderConfig
		crdSource            string
		configPath           string
		tagType              string
		tagProperty          string
	}
	type args struct {
		generatorConfig *GeneratorConfig
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
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
						FromLabels: []string{
							"commonA",
							"commonB",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
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
				Labels: LocalLabelConfig{
					LabelConfig: LabelConfig{
						FromCRD: []string{
							"fromCrdA",
							"fromCrdB",
							"fromCrdC",
						},
					},
				},
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
						FromLabels: []string{
							"fromCrdA",
							"fromCrdB",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{},
			},
			wantErr: false,
		},
		{
			name: "Should validate config in globalLabels",
			fields: fields{
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
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
				generatorConfig: &GeneratorConfig{},
			},
			wantErr: false,
		},
		{
			name: "Should validate config in all places",
			fields: fields{
				Labels: LocalLabelConfig{
					LabelConfig: LabelConfig{
						FromCRD: []string{
							"fromCrdA",
							"fromCrdB",
							"fromCrdC",
						},
					},
				},
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
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
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
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
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
						FromLabels: []string{
							"commonA",
							"commonX",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
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
				Labels: LocalLabelConfig{
					LabelConfig: LabelConfig{
						FromCRD: []string{
							"fromCrdA",
							"fromCrdB",
							"fromCrdC",
						},
					},
				},
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
						FromLabels: []string{
							"fromCrdA",
							"fromCrdX",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{},
			},
			wantErr: true,
		},
		{
			name: "Should have errors from globalLabels",
			fields: fields{
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
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
				generatorConfig: &GeneratorConfig{},
			},
			wantErr: true,
		},
		{
			name: "Should have errors from all",
			fields: fields{
				Labels: LocalLabelConfig{
					LabelConfig: LabelConfig{
						FromCRD: []string{
							"fromCrdA",
							"fromCrdB",
							"fromCrdC",
						},
					},
				},
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
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
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
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
				tagType:              tt.fields.tagType,
				tagProperty:          tt.fields.tagProperty,
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
		OverrideFields       []OverrideField
		Compositions         []Composition
		Tags                 LocalTagConfig
		Labels               LocalLabelConfig
		Provider             ProviderConfig
		crdSource            string
		configPath           string
		tagType              string
		tagProperty          string
	}
	type args struct {
		generatorConfig *GeneratorConfig
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: replaceGlobal,
					},
					LabelConfig: LabelConfig{
						FromCRD: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: replaceGlobal,
					},
					LabelConfig: LabelConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: appendGlobal,
					},
					LabelConfig: LabelConfig{
						FromCRD: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: appendGlobal,
					},
					LabelConfig: LabelConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: replaceGlobal,
					},
					LabelConfig: LabelConfig{
						FromCRD: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: replaceGlobal,
					},
					LabelConfig: LabelConfig{
						FromCRD: []string{},
					},
				},
			},
		},
		{
			name: "Should append labels empty no GlobalHandling",
			fields: fields{
				Labels: LocalLabelConfig{

					LabelConfig: LabelConfig{
						FromCRD: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{

					LabelConfig: LabelConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: appendGlobal,
					},
					LabelConfig: LabelConfig{
						FromCRD: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						FromCRD: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: appendGlobal,
					},
					LabelConfig: LabelConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						Common: replaceGlobal,
					},
					LabelConfig: LabelConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdLC": "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdGC": "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						Common: replaceGlobal,
					},
					LabelConfig: LabelConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						Common: appendGlobal,
					},
					LabelConfig: LabelConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						Common: appendGlobal,
					},
					LabelConfig: LabelConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						Common: replaceGlobal,
					},
					LabelConfig: LabelConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						Common: replaceGlobal,
					},
					LabelConfig: LabelConfig{
						Common: map[string]string{},
					},
				},
			},
		},
		{
			name: "Should append labels empty no GlobalHandling",
			fields: fields{
				Labels: LocalLabelConfig{
					LabelConfig: LabelConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{

					LabelConfig: LabelConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						Common: appendGlobal,
					},
					LabelConfig: LabelConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						Common: appendGlobal,
					},
					LabelConfig: LabelConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: replaceGlobal,
					},
					TagConfig: TagConfig{
						FromLabels: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{

				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: replaceGlobal,
					},
					TagConfig: TagConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: appendGlobal,
					},
					TagConfig: TagConfig{
						FromLabels: []string{
							"fromCrdLA",
							"fromCrdLB",
							"fromCrdLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: appendGlobal,
					},
					TagConfig: TagConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: replaceGlobal,
					},
					TagConfig: TagConfig{
						FromLabels: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: replaceGlobal,
					},
					TagConfig: TagConfig{
						FromLabels: []string{},
					},
				},
			},
		},
		{
			name: "Should append tags empty no GlobalHandling",
			fields: fields{
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
						FromLabels: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: appendGlobal,
					},
					TagConfig: TagConfig{
						FromLabels: []string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						FromLabels: []string{
							"fromCrdGA",
							"fromCrdGB",
							"fromCrdGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: appendGlobal,
					},
					TagConfig: TagConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						Common: replaceGlobal,
					},
					TagConfig: TagConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdLC": "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdGC": "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						Common: replaceGlobal,
					},
					TagConfig: TagConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						Common: appendGlobal,
					},
					TagConfig: TagConfig{
						Common: map[string]string{
							"fromCrdLA": "valueLA",
							"fromCrdLB": "valueLB",
							"fromCrdC":  "valueLC",
						},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						Common: appendGlobal,
					},
					TagConfig: TagConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						Common: replaceGlobal,
					},
					TagConfig: TagConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						Common: replaceGlobal,
					},
					TagConfig: TagConfig{
						Common: map[string]string{},
					},
				},
			},
		},
		{
			name: "Should append tags empty no GlobalHandling",
			fields: fields{
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					TagConfig: TagConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						Common: appendGlobal,
					},
					TagConfig: TagConfig{
						Common: map[string]string{},
					},
				},
			},
			args: args{
				generatorConfig: &GeneratorConfig{
					Tags: TagConfig{
						Common: map[string]string{
							"fromCrdGA": "valueGA",
							"fromCrdGB": "valueGB",
							"fromCrdC":  "valueGC",
						},
					},
				},
			},
			want: fields{
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						Common: appendGlobal,
					},
					TagConfig: TagConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: appendGlobal,
						Common:  appendGlobal,
					},
					LabelConfig: LabelConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: appendGlobal,
						Common:     appendGlobal,
					},
					TagConfig: TagConfig{
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
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
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
					Tags: TagConfig{
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
				Labels: LocalLabelConfig{
					GlobalHandling: GlobalHandlingLabels{
						FromCRD: appendGlobal,
						Common:  appendGlobal,
					},
					LabelConfig: LabelConfig{
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
				Tags: LocalTagConfig{
					GlobalHandling: GlobalHandlingTags{
						FromLabels: appendGlobal,
						Common:     appendGlobal,
					},
					TagConfig: TagConfig{
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
				tagType:              tt.fields.tagType,
				tagProperty:          tt.fields.tagProperty,
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
			if tt.want.tagType != g.tagType {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.tagType, tt.want.tagType)
			}
			if tt.want.tagProperty != g.tagProperty {
				t.Errorf("TestGenerator_UpdateConfig() got = %v, want %v", g.tagProperty, tt.want.tagProperty)
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
		generatorConfig *GeneratorConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should be valid common",
			args: args{
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						Common: map[string]string{
							"commonA": "valueCA",
							"commonB": "valueCB",
							"commonC": "valueCB",
						},
					},
					Tags: TagConfig{
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
				generatorConfig: &GeneratorConfig{

					Tags: TagConfig{
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
				generatorConfig: &GeneratorConfig{
					Labels: LabelConfig{
						Common: map[string]string{
							"commonA": "valueCA",
							"commonB": "valueCB",
							"commonC": "valueCB",
						},
					},
					Tags: TagConfig{
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
