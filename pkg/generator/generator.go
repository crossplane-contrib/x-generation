package generator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	p "github.com/crossplane-contrib/function-patch-and-transform/input/v1beta1"
	t "github.com/crossplane-contrib/x-generation/pkg/types"
	c "github.com/crossplane/crossplane/apis/apiextensions/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type XGenerator struct {
	Group                     string                      `yaml:"group" json:"group"`
	Name                      string                      `yaml:"name" json:"name"`
	Plural                    *string                     `yaml:"plural,omitempty" json:"plural,omitempty"`
	PatchExternalName         *bool                       `yaml:"patchExternalName,omitempty" json:"patchExternalName,omitempty"`
	PatchlName                *bool                       `yaml:"patchName,omitempty" json:"patchName,omitempty"`
	ConnectionSecretKeys      *[]string                   `yaml:"connectionSecretKeys,omitempty" json:"connectionSecretKeys,omitempty"`
	Compositions              []t.Composition             `yaml:"compositions" json:"compositions"`
	Version                   string                      `yaml:"version" json:"version"`
	Crd                       v1.CustomResourceDefinition `yaml:"crd" json:"crd"`
	Provider                  t.ProviderConfig            `yaml:"provider" json:"provider"`
	OverrideFields            []t.OverrideField           `yaml:"overrideFields" json:"overrideFields"`
	OverrideFieldsInClaim     []t.OverrideFieldInClaim    `yaml:"overrideFieldsInClaim" json:"overrideFieldsInClaim"`
	Labels                    t.LocalLabelConfig          `yaml:"labels,omitempty" json:"labels,omitempty"`
	ReadinessChecks           *bool                       `yaml:"readinessChecks,omitempty" json:"readinessChecks,omitempty"`
	ResourceName              *string                     `yaml:"resourceName,omitempty" json:"resourceName,omitempty"`
	UIDFieldPath              *string                     `yaml:"uidFieldPath,omitempty" json:"uidFieldPath,omitempty"`
	ExpandCompositionName     *bool                       `yaml:"expandCompositionName,omitempty" json:"expandCompositionName,omitempty"`
	AdditionalPipelineSteps   []t.PipelineStep            `yaml:"additionalPipelineSteps,omitempty" json:"additionalPipelineSteps,omitempty"`
	TagType                   *string                     `yaml:"tagType,omitempty" json:"tagType,omitempty"`
	TagProperty               *string                     `yaml:"tagProperty,omitempty" json:"tagProperty,omitempty"`
	AutoReadyFunction         *t.AutoReadyFunction        `yaml:"autoReadyFunction,omitempty" json:"autoReadyFunction,omitempty"`
	PatchAndTransfromFunction *string                     `yaml:"patchAndTransfromFunction,omitempty" json:"patchAndTransfromFunction,omitempty"`

	GlobalLabels             []string
	GeneratorConfig          t.GeneratorConfig
	xrdSchema                *v1.JSONSchemaProps
	overrideFieldDefinitions []*OverrideFieldDefinition
}

type OverrideFieldDefinition struct {
	ClaimPath     string
	ManagedPath   string
	Schema        *v1.JSONSchemaProps
	Required      bool
	Replacement   bool
	PathSegments  []pathSegment
	Patches       []p.PatchSetPatch
	OriginalEnum  []v1.JSON
	Overwrites    *t.OverrideFieldInClaim
	IgnoreInClaim bool
}

type NamedComposition struct {
	Name        string
	Composition c.Composition
}

func (g *XGenerator) GenerateXRD() (*c.CompositeResourceDefinition, error) {

	plural := g.nameToPlural()
	defaultCompositionName, _ := g.getDefaultCompositionName()
	version, _ := g.getVersion()
	status, err := g.generateSchema("status")
	if err != nil {
		return nil, err
	}
	status.Properties["observed"] = v1.JSONSchemaProps{
		Description:            "Freeform field containing information about the observed status.",
		Type:                   "object",
		XPreserveUnknownFields: pointer(true),
	}
	status.Properties["uid"] = v1.JSONSchemaProps{
		Description: fmt.Sprintf("The unique ID of this %s resource reported by the provider", g.Name),
		Type:        "string",
	}
	g.overrideFieldDefinitions = mapOverwrittenFields(g.OverrideFieldsInClaim)
	specSchema, err := g.generateSchema("spec")
	if err != nil {
		return nil, err
	}
	g.xrdSchema = specSchema
	xrd := c.CompositeResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apiextensions.crossplane.io/v1",
			Kind:       "CompositeResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "composite" + g.fqdn(),
		},
		Spec: c.CompositeResourceDefinitionSpec{
			ClaimNames: &v1.CustomResourceDefinitionNames{
				Kind:   g.Name,
				Plural: plural,
			},
			DefaultCompositionRef: &c.CompositionReference{
				Name: *defaultCompositionName,
			},
			Group: g.Group,
			Names: v1.CustomResourceDefinitionNames{
				Kind:       "Composite" + g.Name,
				Plural:     "composite" + plural,
				Categories: g.generateCategories(),
			},
			Versions: []c.CompositeResourceDefinitionVersion{
				{
					Name:          g.Version,
					Referenceable: true,
					Served:        true,
					Schema: &c.CompositeResourceValidation{
						OpenAPIV3Schema: runtime.RawExtension{
							Object: &unstructured.Unstructured{
								Object: map[string]interface{}{
									"properties": map[string]interface{}{
										"spec":   g.xrdSchema,
										"status": status,
									},
								},
							},
						},
					},
					AdditionalPrinterColumns: filterCustomResourceColumnDefinitions(version.AdditionalPrinterColumns),
				},
			},
		},
	}

	//
	// g.generateSchema()
	//

	if g.ConnectionSecretKeys != nil {
		xrd.Spec.ConnectionSecretKeys = *g.ConnectionSecretKeys
	}
	xrd.Status = c.CompositeResourceDefinitionStatus{}

	return &xrd, nil
}

func (g *XGenerator) GenerateComposition() ([]NamedComposition, error) {
	compositions := []NamedComposition{}

	for _, comp := range g.Compositions {

		rName := g.Crd.Spec.Names.Kind
		if g.ResourceName != nil {
			rName = *g.ResourceName
		}
		resource := p.ComposedTemplate{
			Name: rName,
			Base: &runtime.RawExtension{
				Raw: g.generateBase(comp),
			},
		}
		if g.ReadinessChecks != nil && !*g.ReadinessChecks {
			resource.ReadinessChecks = []p.ReadinessCheck{{
				Type: p.ReadinessCheckTypeNone,
			},
			}
		}
		name := comp.Name
		if g.ExpandCompositionName != nil && *g.ExpandCompositionName {
			name = "composite" + name + "." + g.Group
		}
		composition := c.Composition{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apiextensions.crossplane.io/v1",
				Kind:       "Composition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					g.GeneratorConfig.CompositionIdentifier + "/provider": comp.Provider,
				},
			},
			Spec: c.CompositionSpec{
				CompositeTypeRef: c.TypeReference{
					APIVersion: g.Group + "/" + g.Version,
					Kind:       "Composite" + g.Name,
				},
				Mode: pointer(c.CompositionModePipeline),
			},
		}

		patchSets := []p.PatchSet{}

		if g.PatchlName == nil || *g.PatchlName {
			var toFieldPath string
			if g.PatchExternalName != nil && !*g.PatchExternalName {
				toFieldPath = "metadata.name"
			} else {
				toFieldPath = "metadata.annotations[crossplane.io/external-name]"
			}
			patchSets = append(patchSets, p.PatchSet{
				Name: "Name",
				Patches: []p.PatchSetPatch{
					{
						Type: p.PatchTypeFromCompositeFieldPath,
						Patch: p.Patch{
							FromFieldPath: pointer("metadata.labels[crossplane.io/claim-name]"),
							ToFieldPath:   &toFieldPath,
						},
					},
				},
			})
		}
		patchSets = append(patchSets, p.PatchSet{
			Name: "External-Name",
			Patches: []p.PatchSetPatch{
				{
					Patch: p.Patch{
						FromFieldPath: pointer("metadata.annotations[crossplane.io/external-name]"),
						ToFieldPath:   pointer("metadata.annotations[crossplane.io/external-name]"),
						Policy: &p.PatchPolicy{
							FromFieldPath: pointer(p.FromFieldPathPolicyOptional),
						},
					},
					Type: p.PatchTypeFromCompositeFieldPath,
				},
			},
		})

		patchSets = append(patchSets, generateLabelPatchset("Common", g.GlobalLabels))

		if g.xrdSchema == nil {
			specSchema, err := g.generateSchema("spec")
			if err != nil {
				return nil, err
			}
			g.xrdSchema = specSchema
		}

		statusSchema, err := g.generateSchema("status")
		if err != nil {
			return nil, err
		}

		xrdStatusSchema := statusSchema
		patchSets = append(patchSets, p.PatchSet{
			Name:    "Parameters",
			Patches: g.generateSortedPropertyPatchesFor(*g.xrdSchema, "spec", p.PatchTypeFromCompositeFieldPath),
		})
		patchSets = append(patchSets, p.PatchSet{
			Name:    "Status",
			Patches: g.generateSortedPropertyPatchesFor(*xrdStatusSchema, "status", p.PatchTypeToCompositeFieldPath),
		})

		labelPatchset := generateLabelPatchset("Labels", g.Labels.FromCRD)

		if len(labelPatchset.Patches) > 0 {

			patchSets = append(patchSets, labelPatchset)
		}

		// composition.Spec.PatchSets = patchSets

		for _, ps := range patchSets {
			resource.Patches = append(resource.Patches, p.ComposedPatch{
				Type:         p.PatchTypePatchSet,
				PatchSetName: pointer(ps.Name),
			})
		}

		resource.Patches = append(resource.Patches, p.ComposedPatch{
			Patch: p.Patch{
				FromFieldPath: pointer(g.getUidFieldPath()),
				ToFieldPath:   pointer("status.uid"),
				Policy: &p.PatchPolicy{
					FromFieldPath: pointer(p.FromFieldPathPolicyOptional),
				},
			},
			Type: p.PatchTypeToCompositeFieldPath,
		})
		resource.Patches = append(resource.Patches, p.ComposedPatch{
			Patch: p.Patch{
				FromFieldPath: pointer("status.conditions"),
				ToFieldPath:   pointer("status.observed.conditions"),
				Policy: &p.PatchPolicy{
					FromFieldPath: pointer(p.FromFieldPathPolicyOptional),
				},
			},
			Type: p.PatchTypeToCompositeFieldPath,
		})

		if g.ConnectionSecretKeys != nil {
			composition.Spec.WriteConnectionSecretsToNamespace = pointer("crossplane-system")
			resource.Patches = append(resource.Patches, p.ComposedPatch{
				Patch: p.Patch{
					FromFieldPath: pointer("metadata.uid"),
					ToFieldPath:   pointer("spec.writeConnectionSecretToRef.name"),
					Policy: &p.PatchPolicy{
						FromFieldPath: pointer(p.FromFieldPathPolicyOptional),
					},
					Transforms: []p.Transform{
						{
							Type: p.TransformTypeString,
							String: &p.StringTransform{
								Format: pointer("%s-secret"),
								Type:   p.StringTransformTypeFormat,
							},
						},
					},
				},
				Type: p.PatchTypeFromCompositeFieldPath,
			})
			for _, k := range *g.ConnectionSecretKeys {

				resource.ConnectionDetails = append(resource.ConnectionDetails, p.ConnectionDetail{
					Name:                    k,
					Type:                    p.ConnectionDetailTypeFromConnectionSecretKey,
					FromConnectionSecretKey: pointer(k),
				})

			}
		}
		// composition.Spec.Resources = []c.ComposedTemplate{
		// 	resource,
		// }

		name = "function-patch-and-transform"
		if g.PatchAndTransfromFunction != nil {
			name = *g.PatchAndTransfromFunction
		}
		patchAndTransform := c.PipelineStep{
			Step: "patch-and-transform",
			FunctionRef: c.FunctionReference{
				Name: name,
			},
		}

		patchAndTransformResource := p.Resources{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "pt.fn.crossplane.io/v1beta1",
				Kind:       "Resources",
			},
			PatchSets: patchSets,
			Resources: []p.ComposedTemplate{resource},
		}

		patchAndTransformRaw := map[string]interface{}{
			"apiVersion": patchAndTransformResource.APIVersion,
			"kind":       patchAndTransformResource.Kind,
			"patchSets":  patchAndTransformResource.PatchSets,
			"resources":  patchAndTransformResource.Resources,
		}
		raw, err := json.Marshal(patchAndTransformRaw)
		if err != nil {
			return nil, err
		}
		patchAndTransform.Input = &runtime.RawExtension{
			Raw: raw,
		}
		composition.Spec.Pipeline = append(composition.Spec.Pipeline, patchAndTransform)

		if g.AdditionalPipelineSteps != nil {
			startSteps := []c.PipelineStep{}
			for _, s := range g.AdditionalPipelineSteps {
				step, err := g.generateAdditonalPipelineStep(s)
				if err != nil {
					return nil, err
				}
				if step != nil {
					if s.Before {
						startSteps = append(startSteps, *step)
					} else {
						composition.Spec.Pipeline = append(composition.Spec.Pipeline, *step)
					}
				}
			}
			composition.Spec.Pipeline = append(startSteps, composition.Spec.Pipeline...)
			if g.AutoReadyFunction == nil || g.AutoReadyFunction.Generate == nil || *g.AutoReadyFunction.Generate {
				functionName := "function-auto-ready"
				if g.AutoReadyFunction != nil && g.AutoReadyFunction.Name != nil {
					functionName = *g.AutoReadyFunction.Name
				}
				composition.Spec.Pipeline = append(composition.Spec.Pipeline, c.PipelineStep{
					Step: "automatically-detect-readiness",
					FunctionRef: c.FunctionReference{
						Name: functionName,
					},
				})
			}
		}
		compositions = append(compositions, NamedComposition{
			Name:        comp.Name,
			Composition: composition,
		})
	}
	return compositions, nil
}

func (g *XGenerator) updateKubernetesValidation(schema *v1.JSONSchemaProps) {

	kubernetesValidations := schema.XValidations

	replaceMap := map[string]string{}
	replaceMessageMap := map[string]string{}

	for _, override := range g.OverrideFieldsInClaim {
		if override.ManagedPath != nil {
			// var updatedClaimPath, updatedManagedPath string
			updatedClaimPath := strings.Replace(override.ClaimPath, "spec", "self", 1)
			updatedManagedPath := strings.Replace(*override.ManagedPath, "spec", "self", 1)
			replaceMap[updatedManagedPath] = updatedClaimPath
			replaceMessageMap[*override.ManagedPath] = override.ClaimPath
		}
	}
	validationMapArray := []v1.ValidationRule{}
	for _, validation := range kubernetesValidations {
		rule := validation.Rule
		message := validation.Message
		for old, new := range replaceMap {
			rule = strings.Replace(rule, old, new, -1)
		}
		for old, new := range replaceMessageMap {
			message = strings.Replace(message, old, new, -1)
		}

		validation.Rule = rule
		validation.Message = message
		validationMapArray = append(validationMapArray, validation)
	}
	schema.XValidations = validationMapArray
}

func (g *XGenerator) generatePropertyPatchesFor(schema v1.JSONSchemaProps, path string, patchType p.PatchType) []p.PatchSetPatch {
	patches := []p.PatchSetPatch{}
	if schema.Type == "object" && schema.AdditionalProperties == nil {
		for key, prop := range schema.Properties {
			patches = append(patches, g.generatePropertyPatchesFor(prop, path+"."+key, patchType)...)
		}
	} else {
		definition := getOverwriteDefinition(g.overrideFieldDefinitions, path, CLAIMPATH)
		var toFieldPath string
		if definition != nil {
			toFieldPath = definition.ManagedPath
		} else {
			toFieldPath = path
		}
		definitionPatches := getPatchesFromDefinition(definition, patchType)
		if len(definitionPatches) > 0 {
			patches = append(patches, definitionPatches...)
		} else {
			patches = append(patches, p.PatchSetPatch{
				Patch: p.Patch{
					FromFieldPath: pointer(path),
					ToFieldPath:   pointer(toFieldPath),
					Policy: &p.PatchPolicy{
						FromFieldPath: pointer(p.FromFieldPathPolicyOptional),
					},
				},
				Type: patchType,
			})
		}
	}

	return patches
}

func (g *XGenerator) generatePropertyPatchesForIgnoredProperties(patchType p.PatchType) []p.PatchSetPatch {
	patches := []p.PatchSetPatch{}
	definitions := getOverwriteDefinitionForIgnoredFileds(g.overrideFieldDefinitions)
	for _, d := range definitions {
		definitionPatches := getPatchesFromDefinition(d, patchType)

		if len(definitionPatches) > 0 {
			patches = append(patches, definitionPatches...)
		}
	}
	return patches
}
func (g *XGenerator) generateSortedPropertyPatchesFor(schema v1.JSONSchemaProps, path string, patchType p.PatchType) []p.PatchSetPatch {
	patches := g.generatePropertyPatchesFor(schema, path, patchType)
	if patchType == p.PatchTypeFromCompositeFieldPath {
		patches = append(patches, g.generatePropertyPatchesForIgnoredProperties(patchType)...)
	}
	sort.Slice(patches, func(i, j int) bool {
		return *patches[i].FromFieldPath < *patches[j].FromFieldPath
	})
	return patches
}

func fieldsTo(path string, fields []string) []string {
	values := []string{}
	for _, f := range fields {
		values = append(values, fmt.Sprintf("%s['%s']", path, f))
	}
	return values
}

func generateOptionalToFromPatches(paths []string, patchType p.PatchType) []p.PatchSetPatch {
	optional := make([]p.FromFieldPathPolicy, len(paths))
	for i := 0; i < len(paths); i++ {
		optional[i] = p.FromFieldPathPolicyOptional
	}
	patches, _ := generatePatches(paths, paths, optional, patchType)

	return patches
}

func generatePatches(fromFields, toFields []string, policies []p.FromFieldPathPolicy, patchType p.PatchType) ([]p.PatchSetPatch, error) {

	patches := []p.PatchSetPatch{}
	if len(fromFields) != len(toFields) || len(fromFields) != len(policies) {
		return patches, errors.New("unequal length of parameters")
	}

	for i := 0; i < len(fromFields); i++ {
		patches = append(patches, p.PatchSetPatch{
			Patch: p.Patch{
				FromFieldPath: &fromFields[i],
				ToFieldPath:   &toFields[i],
				Policy: &p.PatchPolicy{
					FromFieldPath: &policies[i],
				},
			},
			Type: patchType,
		})
	}
	return patches, nil
}

func generateLabelPatchset(name string, fields []string) p.PatchSet {

	labelFields := fieldsTo("metadata.labels", fields)
	patches := generateOptionalToFromPatches(labelFields, p.PatchTypeFromCompositeFieldPath)

	return p.PatchSet{
		Name:    name,
		Patches: patches,
	}
}

func (g *XGenerator) nameToPlural() string {
	if g.Plural != nil {
		return strings.ToLower(*g.Plural)
	}

	lname := strings.ToLower(g.Name)
	last := string(lname[len(lname)-1])
	if last == "y" {
		lname = string(lname[:len(lname)-1]) + "ie"
	}
	lname = lname + "s"
	return lname
}

func (g *XGenerator) fqdn() string {

	plural := g.nameToPlural()
	return fmt.Sprintf("%s.%s", strings.ToLower(plural), g.Group)
}

func (g *XGenerator) getDefaultCompositionName() (*string, error) {

	for _, c := range g.Compositions {
		if c.Default {
			if g.ExpandCompositionName != nil && *g.ExpandCompositionName {
				return pointer("composite" + c.Name + "." + g.Group), nil
			}
			return &c.Name, nil
		}
	}
	return nil, errors.New("could not find a default composition - exactly one composition must have default: true")
}

func (g *XGenerator) generateCategories() []string {
	return []string{
		"crossplane",
		"composition",
		strings.Split(g.Group, ".")[0],
	}
}

func (g *XGenerator) getVersion() (*v1.CustomResourceDefinitionVersion, error) {
	for _, v := range g.Crd.Spec.Versions {
		if v.Name == g.Provider.CRD.Version {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("could not find CRD with version %s", g.Provider.Version)
}

func (g *XGenerator) generateSchema(prop string) (*v1.JSONSchemaProps, error) {
	version, _ := g.getVersion()
	a := version.Schema.OpenAPIV3Schema.Properties[prop]
	b := g.generateSchemaFor(a, prop)
	err := g.overwrittenFields(b, prop)
	if err != nil {
		return nil, err
	}
	g.updateKubernetesValidation(b)
	return b, nil
}

func (g *XGenerator) generateSchemaFor(schema v1.JSONSchemaProps, path string) *v1.JSONSchemaProps {
	switch schema.Type {
	case "object":
		result := g.generateSchemaForObject(schema, path)
		return result
	}
	return nil
}

func (g *XGenerator) generateSchemaForObject(schema v1.JSONSchemaProps, path string) *v1.JSONSchemaProps {
	result := schema.DeepCopy()
	for key, value := range schema.Properties {
		currentPath := path + "." + key
		if !listIncludes(g.getIgnored(), currentPath) {
			overwrite := getOverwriteDefinition(g.overrideFieldDefinitions, currentPath, MANAGEDPATH)
			propertySchema := g.generateSchemaFor(value, currentPath)
			if propertySchema != nil {
				result.Properties[key] = *propertySchema
			}
			if overwrite != nil && !overwrite.IgnoreInClaim {
				if overwrite.Schema == nil {
					overwrite.Schema = pointer(result.Properties[key])
				}
				delete(result.Properties, key)
				if listIncludes(result.Required, key) {
					overwrite.Required = true
					result.Required = filterList(result.Required, key)
				}
			}
		} else {
			delete(result.Properties, key)
			result.Required = filterList(result.Required, key)
		}
	}
	if listIncludes(g.getIgnored(), path+"."+"default") {
		result.Default = nil
	}
	return result
}

func filterList(list []string, value string) []string {
	filterd := []string{}
	for _, e := range list {
		if e != value {
			filterd = append(filterd, e)
		}
	}
	return filterd
}

func listIncludes(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func pointer[K any](input K) *K {
	return &input
}

func filterCustomResourceColumnDefinitions(list []v1.CustomResourceColumnDefinition) []v1.CustomResourceColumnDefinition {
	result := []v1.CustomResourceColumnDefinition{}
	for _, c := range list {
		if !strings.HasPrefix(c.JSONPath, ".status.conditions") {
			result = append(result, c)
		}
	}
	return result
}

func (g *XGenerator) getIgnored() []string {
	defaultIgnored := []string{
		"status.conditions",
		"spec.writeConnectionSecretToRef",
		"spec.forProvider.tags",
		"spec.forProvider.tagSpecifications",
		"spec.forProvider.tagging",
		"spec.providerConfigRef.default",
		"spec.providerRef",
		"spec.publishConnectionDetailsTo.configRef.default",
	}
	for _, o := range g.OverrideFields {
		if o.Ignore {
			defaultIgnored = append(defaultIgnored, o.Path)
		}
	}
	for _, o := range g.OverrideFieldsInClaim {
		if o.Ignore {
			defaultIgnored = append(defaultIgnored, o.ClaimPath)
		}
	}
	return defaultIgnored
}

func (g *XGenerator) generateBase(comp t.Composition) []byte {

	version, _ := g.getVersion()
	spec := version.Schema.OpenAPIV3Schema.Properties["spec"]

	baseSpec := map[string]interface{}{}

	if _, ok := spec.Properties["providerConfigRef"]; ok {

		baseSpec["providerConfigRef"] = map[string]interface{}{
			"name": "default",
		}
	}

	commonLabels := map[string]string{}

	for key, value := range g.Labels.Common {
		commonLabels[key] = value
	}
	base := map[string]interface{}{
		"apiVersion": g.Crd.Spec.Group + "/" + g.Provider.CRD.Version,
		"kind":       &g.Crd.Spec.Names.Kind,
		"metadata":   map[string]interface{}{},
		"spec":       baseSpec,
	}

	if len(commonLabels) > 0 {
		base["metadata"].(map[string]interface{})["labels"] = commonLabels
	}

	if g.ConnectionSecretKeys != nil {
		base["spec"].(map[string]interface{})["writeConnectionSecretToRef"] = map[string]string{
			"namespace": "crossplane-system",
		}
	}

	base = applyOverrideFields(base, g.OverrideFields)

	object, err := json.Marshal(base)
	if err != nil {
		fmt.Printf("unable to marshal base: %v\n", err)
	}

	return object
}

func applyOverrideFields(base map[string]interface{}, overrideFields []t.OverrideField) map[string]interface{} {
	for _, overwite := range overrideFields {
		if overwite.Value != nil {
			path := splitPath(overwite.Path)
			var current interface{}
			current = base
			pathLength := len(path)

			for i := 0; i < pathLength-1; i++ {
				segment := path[i]
				property := path[i].path
				if segment.pathType == "object" {
					if current.(map[string]interface{})[property] == nil {
						current.(map[string]interface{})[property] = map[string]interface{}{}
					}

					current = current.(map[string]interface{})[property].(map[string]interface{})
				} else if segment.pathType == "array" {
					if current.(map[string]interface{})[property] == nil {
						current.(map[string]interface{})[property] = []map[string]interface{}{}
					}

					var b interface{}
					b = current.(map[string]interface{})[property].([]map[string]interface{})
					currentSize := len(b.([]map[string]interface{}))
					wantedSize := segment.arrayPosition + 1
					if currentSize < wantedSize {
						sizeToGrow := wantedSize - currentSize
						b = slices.Grow(b.([]map[string]interface{}), sizeToGrow)
						b = b.([]map[string]interface{})[:cap(b.([]map[string]interface{}))]
						b.([]map[string]interface{})[segment.arrayPosition] = map[string]interface{}{}
					}
					current.(map[string]interface{})[property] = b
					current = b.([]map[string]interface{})[segment.arrayPosition]
				}
			}
			segment := path[pathLength-1]
			if segment.pathType == "object" {
				(current).(map[string]interface{})[path[pathLength-1].path] = overwite.Value
			}
			if segment.pathType == "array" {
				property := path[pathLength-1].path

				if (current.(map[string]interface{}))[property] == nil {
					(current.(map[string]interface{}))[property] = []interface{}{}
				}

				var b interface{}
				b = (current.(map[string]interface{}))[property].([]interface{})
				currentSize := len(b.([]interface{}))
				wantedSize := segment.arrayPosition + 1
				if currentSize < wantedSize {
					sizeToGrow := wantedSize - currentSize
					b = slices.Grow(b.([]interface{}), sizeToGrow)
					b = b.([]interface{})[:cap(b.([]interface{}))]
					(b.([]interface{})[segment.arrayPosition]) = overwite.Value
				}
				current.(map[string]interface{})[property] = b

			}
		}
	}
	return base
}

type pathSegment struct {
	path          string
	pathType      string
	arrayPosition int
}

func splitPath(path string) []pathSegment {
	inString := false
	result := []pathSegment{}
	current := ""
	escaped := false
	for _, r := range path {
		switch r {
		case '"':
			inString = !inString
			escaped = false

		case '\\':
			escaped = true
		case '.':
			if current != "" {
				if !inString && !escaped {
					segment := pathSegment{
						path:     current,
						pathType: "object",
					}
					result = append(result, segment)
					current = ""
				} else {
					current += string(r)
				}
			}
		case '[':
			if !inString && !escaped {
				segment := pathSegment{
					path:     current,
					pathType: "object",
				}
				result = append(result, segment)
				current = ""
			} else {
				current += string(r)
			}
		case ']':
			if !inString && !escaped {
				lastSegemnt := result[len(result)-1]
				arrayIndex, err := strconv.Atoi(current)
				if err == nil {
					lastSegemnt.pathType = "array"
					lastSegemnt.arrayPosition = arrayIndex
					result[len(result)-1] = lastSegemnt
				} else {
					segment := pathSegment{
						path:     current,
						pathType: "object",
					}
					result = append(result, segment)
				}
				current = ""

			} else {
				current += string(r)
			}
		default:
			current += string(r)
			escaped = false

		}
	}
	if current != "" {
		segment := pathSegment{
			path:     current,
			pathType: "object",
		}
		result = append(result, segment)
	}
	return result
}

func (g *XGenerator) getUidFieldPath() string {
	if g.UIDFieldPath != nil {
		return *g.UIDFieldPath
	}
	return "metadata.annotations[\"crossplane.io/external-name\"]"
}

func (g *XGenerator) generateAdditonalPipelineStep(s t.PipelineStep) (*c.PipelineStep, error) {
	rawInput, _ := json.Marshal(s.Input)
	rawInput = bytes.ReplaceAll(rawInput, []byte("{tagProperty}"), []byte(*g.TagProperty))
	rawInput = bytes.ReplaceAll(rawInput, []byte("{tagType}"), []byte(*g.TagType))
	render := true
	if s.Condition != nil {
		var err error
		data := ConditonData{
			TagProperty: *g.TagProperty,
			TagType:     *g.TagType,
		}
		render, err = EvaluateCondition(s.Condition, data)
		if err != nil {
			return nil, err
		}
	}
	if render {
		return &c.PipelineStep{
			Step: s.Step,
			FunctionRef: c.FunctionReference{
				Name: s.FunctionRef.Name,
			},
			Input: &runtime.RawExtension{
				Raw: rawInput,
			},
		}, nil
	}

	return nil, nil
}

func (g *XGenerator) overwrittenFields(schema *v1.JSONSchemaProps, path string) error {
	for _, o := range g.overrideFieldDefinitions {
		if o.PathSegments[0].path == path && !o.IgnoreInClaim {
			err := overwrittenFields(schema, path, o, 1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func overwrittenFields(schema *v1.JSONSchemaProps, path string, definition *OverrideFieldDefinition, level int) error {
	if len(definition.PathSegments)-1 > level {
		if schema.Type == "object" {
			pathSegment := definition.PathSegments[level].path
			prop, ok := schema.Properties[pathSegment]
			if !ok {
				schema.Properties[pathSegment] = v1.JSONSchemaProps{
					Type:       "object",
					Properties: map[string]v1.JSONSchemaProps{},
				}
				prop = schema.Properties[pathSegment]
			}
			err := overwrittenFields(&prop, path+"."+pathSegment, definition, level+1)
			if err != nil {
				return err
			}
		}
	} else {
		pathSegment := definition.PathSegments[level].path
		if definition.Schema == nil {
			return fmt.Errorf("schema must be given for new property: %s", definition.ClaimPath)
		}
		err := handleEnumFor(definition.Schema, definition)
		if err != nil {
			return err
		}
		schema.Properties[pathSegment] = *definition.Schema
	}
	return nil
}

type definitionProperty string

const CLAIMPATH definitionProperty = "claimPath"
const MANAGEDPATH definitionProperty = "managedPath"

func getOverwriteDefinition(list []*OverrideFieldDefinition, path string, prop definitionProperty) *OverrideFieldDefinition {
	for _, d := range list {
		var oPath string
		if prop == CLAIMPATH {
			oPath = d.ClaimPath
		} else {
			oPath = d.ManagedPath
		}
		if oPath == path {
			return d
		}
	}
	return nil
}

func getOverwriteDefinitionForIgnoredFileds(list []*OverrideFieldDefinition) []*OverrideFieldDefinition {
	definitions := []*OverrideFieldDefinition{}
	for _, d := range list {
		if d.IgnoreInClaim {
			definitions = append(definitions, d)
		}
	}
	return definitions
}

func mapOverwrittenFields(fields []t.OverrideFieldInClaim) []*OverrideFieldDefinition {
	overrideFieldDefinitions := []*OverrideFieldDefinition{}
	for _, o := range fields {
		definition := &OverrideFieldDefinition{

			ClaimPath: o.ClaimPath,

			PathSegments: splitPath(o.ClaimPath),
			Overwrites:   &o,
		}
		if o.Ignore {
			definition.IgnoreInClaim = true
		}
		if o.ManagedPath != nil {
			definition.ManagedPath = *o.ManagedPath
			definition.Replacement = true
		} else {
			definition.ManagedPath = o.ClaimPath
			definition.Replacement = false
			if o.OverrideSettings != nil {
				definition.Schema = o.OverrideSettings.Property
			}
		}
		overrideFieldDefinitions = append(overrideFieldDefinitions, definition)
	}
	return overrideFieldDefinitions
}

func handleEnumFor(schema *v1.JSONSchemaProps, definition *OverrideFieldDefinition) error {
	if definition.Overwrites.OverrideSettings != nil {
		if definition.Overwrites.OverrideSettings.NewEnum != nil {
			if schema.Enum == nil {
				schema.Enum = definition.Overwrites.OverrideSettings.NewEnum
				return nil
			} else {
				return errors.New("cannot set new enum to existing enum. Use enum property to overwrite existing enum")
			}
		}
		if definition.Overwrites.OverrideSettings.Enum != nil {
			if schema.Enum == nil {
				return errors.New("cannot overwirite enum if non existing. Use newEnum property to create new enum")
			}
			definition.OriginalEnum = schema.Enum
			schema.Enum = handleExistingEnum(schema.Enum, definition.Overwrites.OverrideSettings.Enum)
		}
	}
	return nil
}

func handleExistingEnum(existing []v1.JSON, enumValues []*t.EnumValue) []v1.JSON {
	newEnum := []v1.JSON{}
	for _, e := range existing {
		enumValue := getMatchingEnumValue(e, enumValues)
		if enumValue == nil || enumValue.Type != t.EnumValueTypeRemove {
			newEnum = append(newEnum, e)
		}
	}
	for _, e := range enumValues {
		if e.Type == t.EnumValueTypeAdd {
			newEnum = append(newEnum, e.Value)
		}
	}
	return newEnum
}

func getMatchingEnumValue(value v1.JSON, enumValues []*t.EnumValue) *t.EnumValue {
	for _, e := range enumValues {
		if e.Value.String() == value.String() {
			return e
		}
	}
	return nil
}

func getPatchesFromDefinition(definition *OverrideFieldDefinition, patchType p.PatchType) []p.PatchSetPatch {
	patches := []p.PatchSetPatch{}
	if definition != nil {
		if definition.Overwrites != nil && definition.Overwrites.OverrideSettings != nil {
			if definition.Overwrites.OverrideSettings.Patches != nil {
				patches = append(patches, definition.Overwrites.OverrideSettings.Patches...)

			} else if definition.OriginalEnum != nil {
				transformPairs := map[string]v1.JSON{}

				for _, e := range definition.OriginalEnum {
					newEnum := getMatchingEnumValue(e, definition.Overwrites.OverrideSettings.Enum)
					if newEnum == nil {
						transformPairs[jsonToString(e.Raw)] = e
					} else if newEnum.Type == t.EnumValueTypeMapTo {
						transformPairs[jsonToString(e.Raw)] = *newEnum.MapTo
					}
				}
				for _, e := range definition.Overwrites.OverrideSettings.Enum {
					if e.Type == t.EnumValueTypeAdd {
						transformPairs[jsonToString(e.Value.Raw)] = *e.MapTo
					}
				}
				patch := p.PatchSetPatch{
					Patch: p.Patch{
						FromFieldPath: pointer(definition.ClaimPath),
						ToFieldPath:   pointer(definition.ManagedPath),
						Policy: &p.PatchPolicy{
							FromFieldPath: pointer(p.FromFieldPathPolicyOptional),
						},
						Transforms: []p.Transform{
							{
								Type: p.TransformTypeMap,
								Map: &p.MapTransform{
									Pairs: transformPairs,
								},
							},
						},
					},
					Type: patchType,
				}
				patches = append(patches, patch)
			}
		} else if definition.Replacement {
			patches = append(patches, p.PatchSetPatch{
				Patch: p.Patch{
					FromFieldPath: pointer(definition.ClaimPath),
					ToFieldPath:   pointer(definition.ManagedPath),
					Policy: &p.PatchPolicy{
						FromFieldPath: pointer(p.FromFieldPathPolicyOptional),
					},
				},
				Type: patchType,
			})
		}
	}
	return patches
}

func jsonToString(json []byte) string {
	data := string(json[:])
	if data[0] == '"' && data[len(data)-1] == '"' {
		return data[1 : len(data)-1]
	}
	return data
}
