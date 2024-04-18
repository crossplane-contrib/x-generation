package types

import (
	p "github.com/crossplane-contrib/function-patch-and-transform/input/v1beta1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type OverrideField struct {
	Path     string      `yaml:"path" json:"path"`
	Value    interface{} `yaml:"value,omitempty" json:"value,omitempty"`
	Override interface{} `yaml:"override,omitempty" json:"override,omitempty"`
	Ignore   bool        `yaml:"ignore" json:"ignore"`
}

type Composition struct {
	Name     string `yaml:"name" json:"name"`
	Provider string `yaml:"provider" json:"provider"`
	Default  bool   `yaml:"default" json:"default"`
}

type GeneratorConfig struct {
	CompositionIdentifier   string               `yaml:"compositionIdentifier" json:"compositionIdentifier"`
	Provider                GlobalProviderConfig `yaml:"provider" json:"provider"`
	Tags                    TagConfig            `yaml:"tags,omitempty" json:"tags,omitempty"`
	Labels                  LabelConfig          `yaml:"labels,omitempty" json:"labels,omitempty"`
	UsePipeline             *bool                `yaml:"usePipeline,omitempty" json:"usePipeline,omitempty"`
	ExpandCompositionName   *bool                `yaml:"expandCompositionName,omitempty" json:"expandCompositionName,omitempty"`
	AdditionalPipelineSteps []PipelineStep       `yaml:"additionalPipelineSteps,omitempty" json:"additionalPipelineSteps,omitempty"`
	AutoReadyFunction       *AutoReadyFunction   `yaml:"autoReadyFunction,omitempty" json:"autoReadyFunction,omitempty"`
}

type AutoReadyFunction struct {
	Generate *bool   `yaml:"generate,omitempty" json:"generate,omitempty"`
	Name     *string `yaml:"name,omitempty" json:"name,omitempty"`
}

type PipelineFunction struct {
	Name string `yaml:"name" json:"name"`
}

type PipelineStep struct {
	Step        string                 `yaml:"step" json:"step"`
	FunctionRef PipelineFunction       `yaml:"functionRef" json:"functionRef"`
	Condition   *string                `yaml:"condition,omitempty" json:"condition,omitempty"`
	Input       map[string]interface{} `yaml:"input" json:"input"`
	Before      bool                   `yaml:"before" json:"before"`
}

type TagConfig struct {
	FromLabels []string          `yaml:"fromLabels,omitempty" json:"fromLabels,omitempty"`
	Common     map[string]string `yaml:"common,omitempty" json:"common,omitempty"`
}
type LabelConfig struct {
	FromCRD []string          `yaml:"fromCRD,omitempty" json:"fromCRD,omitempty"`
	Common  map[string]string `yaml:"common,omitempty" json:"common,omitempty"`
}

type GlobalHandlingType string

type GlobalHandlingTags struct {
	FromLabels GlobalHandlingType `yaml:"fromLabels,omitempty" json:"fromLabels,omitempty"`
	Common     GlobalHandlingType `yaml:"common,omitempty" json:"common,omitempty"`
}

type GlobalHandlingLabels struct {
	FromCRD GlobalHandlingType `yaml:"fromCRD,omitempty" json:"fromCRD,omitempty"`
	Common  GlobalHandlingType `yaml:"common,omitempty" json:"common,omitempty"`
}

type LocalTagConfig struct {
	TagConfig
	GlobalHandling GlobalHandlingTags `yaml:"globalHandling,omitempty" json:"globalHandling,omitempty"`
}
type LocalLabelConfig struct {
	LabelConfig
	GlobalHandling GlobalHandlingLabels `yaml:"globalHandling,omitempty" json:"globalHandling,omitempty"`
}

type CrdConfig struct {
	File    string `yaml:"file" json:"file"`
	Version string `yaml:"version" json:"version"`
}

type GlobalProviderConfig struct {
	Name    string  `yaml:"name" json:"name"`
	Version string  `yaml:"version" json:"version"`
	BaseURL *string `yaml:"baseURL,omitempty" json:"baseURL,omitempty"`
}
type ProviderConfig struct {
	GlobalProviderConfig
	CRD CrdConfig `yaml:"crd" json:"crd"`
}

type OverrideFieldInClaim struct {
	ClaimPath        string            `yaml:"claimPath" json:"claimPath"`
	ManagedPath      *string           `yaml:"managedPath,omitempty" json:"managedPath,omitempty"`
	OverrideSettings *OverrideSettings `yaml:"overrideSettings,omitempty" json:"overrideSettings,omitempty"`
	Description      *string           `yaml:"description,omitempty" json:"description,omitempty"`
}

type OverrideSettings struct {
	Property *v1.JSONSchemaProps `yaml:"property,omitempty" json:"property,omitempty"`
	Patches  []p.PatchSetPatch   `yaml:"patches" json:"patches"`
	Enum     []*EnumValue        `yaml:"enum" json:"enum"`
	NewEnum  []v1.JSON           `yaml:"newEnum" json:"newEnum"`
}

type EnumValueType string

const EnumValueTypeAdd EnumValueType = "add"
const EnumValueTypeRemove EnumValueType = "remove"
const EnumValueTypeMapTo EnumValueType = "map"

type EnumValue struct {
	Value v1.JSON       `yaml:"value" json:"value" protobuf:"bytes,20,rep,name=value"`
	Type  EnumValueType `yaml:"type" json:"type"`
	MapTo *v1.JSON      `yaml:"mapTo" json:"mapTo"`
}
