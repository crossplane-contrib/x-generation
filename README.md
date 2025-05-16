# X-GENERATION

generate compositions from crossplane provider crds

## configure

The generation of crds can be configured in two places, either in the global configuration file, or in the local generation file for each composition.
### global configuration
In the local configuration, the provider used and labels and tags that should be used for all generated crds can be configured in a global configuration. The default name for this file is `generator-config.yaml`, the filename can be overwritten using the `--configFile` flag. 

| Property              | Type              | Description |
|-----------------------|-------------------|-------------|
| compositionIdentifier | string            | Defines the refix used for the provider label of the composition |
| provider              | object            | Object used to configure the provider used for the generation |
| provider.baseURL      | string            | The url globaly used to retrieve the crds needed for generating the compositions, three placeholders are provided during the generation of compositions: The name of the provider, the version of the provider and the crd file name|
| provider.name         | string            | The name of the provider |
| provider.version      | string            | The version of the provider |
| labels                | object            | Configure the labels and label patches for each crd |
| labels.fromCRD        | array of strings  | For each entry `e` a patch that copies the value of the `metadata.labels[e]` field from the CompositeResourceDefinition to the same field of the resource |
| labels.common         | object of strings | For each property of the object a label with the given value is created in the resource |
| tags                  | object            | Configure the tags and tag patches for each crd |
| tags.fromLabels       | array of strings  | For each entry `e` a patch that copies the value of the `metadata.labels[e]` field to a tag with the same name and value is created
| tags.common           | object of strings | For each property of the object a tag with the given value is created in the resource |
| usePipeline           | boolean | if true, x-generation generates compositions in pipeline mode, additional pipelinestepts can be added using `additionalPipelineSteps` |
| additionalPipelineSteps           | array of objects | add additional pipeline steps when in pipeline mode, see section using pipelelines  |


The values in `tags.fromLabels` must exist in `lables.fromCRD` otherwise no values that can be patched to the resources exist.

The creation of tags depends on the underlying resource crd, the generator can distinguish between crds without tags at all, crds with objects of strings, arrays of key-value pairs and arrays of tagKey-tagValue pairs. If tags reside inside forProvider.tagging.tagSet, this property is used instead of forProvider.tags.

 #### example
 ```yaml
provider:
  baseURL: https://raw.githubusercontent.com/crossplane-contrib/%s/%s/package/crds/%s
  name: provider-aws
  version: v0.32.0
labels:
  fromCRD:
    - controlling.example.cloud/cost-reference
    - controlling.example.cloud/owner
    - controlling.example.cloud/product
    - tags.example.cloud/account
    - tags.example.cloud/environment
    - tags.example.cloud/protection-requirement
    - tags.example.cloud/repourl
    - tags.example.cloud/zone
  common:
    commonLabelA: commonLabelAValue
    commonLabelB: commonLabelBValue
tags:
  fromLabels:
    - tags.example.cloud/account
    - tags.example.cloud/environment
    - tags.example.cloud/protection-requirement
    - tags.example.cloud/repourl
    - tags.example.cloud/zone
  common:
    commonTagA: comonTagAValue
    commonTagB: comonTagBValue
 ```
# local configuration

The local configuration is placed in the subfolder of the composition to be created. The name of the file defaults to `generate.yaml`. The name of the file can be changed using the `inputName`- flag. Settings in the local configuration overwirte settings in the global configuration.

| Property                       | Type                  | Description |
|--------------------------------|-----------------------|-------------|
| group                          | string                | The group that should be used for the composition |
| name                           | string                | The name that should be used for the composition |
| version                        | string                | The version that should be used for the composition |
| provider                       | object                | Object used to configure the provider used for the generation |
| provider.baseURL               | string                | The url used to retrieve the crd needed for generating the composition, three placeholders are provided during the generation of compositions: The name of the provider, the version of the provider and the crd file name|
| provider.name                  | string                | The name of the provider |
| provider.version               | string                | The version of the provider |
| provider.crd                   | object                | Object used to configure the crd used for the generation |
| provider.crd.file              | object                | The name of the crd file used for generating the composition |
| provider.crd.version           | object                | The version of the object in the crd file used for generating the composition |
| ignore                         | boolean               | If true, no composition is created for this configuration |
| labels                         | object                | Configure the labels and label patches for each crd |
| labels.fromCRD                 | array of strings      | For each entry `e` a patch that copies the value of the `metadata.labels[e]` field from the CompositeResourceDefinition to the same field of the resource |
| labels.common                  | object of strings     | For each property of the object a label with the given value is created in the resource |
| labels.globalHandling.fromCRD  | "append" or "replace" | If append, the labels in labels.fromCRD are appended to the labels in the global configuration labels.fromCRD, otherwise those will be replaced |
| labels.globalHandling.common   | "append" or "replace" | If append, the labels in labels.common are appended to the labels in the global configuration labels.common, otherwise those will be replaced |
| tags                           | object                | Configure the tags and tag patches for each crd |
| tags.fromLabels                | array of strings      | For each entry `e` a patch that copies the value of the `metadata.labels[e]` field to a tag with the same name and value is created
| tags.common                    | object of strings     | For each property of the object a tag with the given value is created in the resource |
| tags.globalHandling.fromLabels | "append" or "replace" | If append, the tags in tags.fromLabels are appended to the tags in the global configuration tags.fromLabels, otherwise those will be replaced |
| tags.globalHandling.common     | "append" or "replace" | If append, the tags in labels.common are appended to the tasg in the global configuration tags.common, otherwise those will be replaced |
| overrideFieldsInClaim          | object                | This optional property can be used to override the names in the composite and the claim or add properties. See description below |
| patchName                      | boolean               | If set to false, the name of the object will not be patched, otherwise`patchExternalName` decides if the name of the claim will be patched to `metadata.name` or `metadata.annotations[crossplane.io/external-name]` |
| patchExternalName              | boolean               | Decides if if the name of the claim will be patched to `metadata.name` or `metadata.annotations[crossplane.io/external-name]`. Not applied if `patchName` is false |
| defaultCompositeDeletePolicy   | string                | This optional property can be used to set the defaultCompositeDeletePolicy on the xrd, possible values Foreground or Background |


## overrideFieldsInClaim
The overrideFieldsInClaim property can be used to change the name of a property in the claim and the composite or to add properties in the claim and composite. This can for example be helpfull if one wants to change the provider of the managed resource without changing the crds for the claim and the composite. OverrideFieldsInClaim has the following properties:

| Property                  | Type        | Description |
|---------------------------|-------------|-------------|
| claimPath                 | string      | The path of the property in the claim and the composite |
| managedPath               | string      | The path of the property in the managed resource. Currently only the name of the property is allowed to change between the claimPath and the managedPath |
| description               | string      | An optional description to override the description of the property from the managed resource |
| overrideSettings          | object      | This allows to override the definition of the new property and the patches applied in the composition for it |
| overrideSettings.property | interface{} | The definition of the property |
| overrideSettings.patches  | []Patch     | A list of pathces that will be placed inside the composition for this property |

To simply rename a property, only claimPath and managedPath is needed. The definition and description is taken from the property in the managed resource, a patch is applied to patch from the new property in the composite to the old name in the managed resource. E.g.:

```yaml
overrideFieldsInClaim:
  - claimPath: spec.forProvider.assumeRolePolicyDocument # user sees assumeRolePolicyDocument
    managedPath: spec.forProvider.assumeRolePolicy # managed resource has assumeRolePolicy
```
leads to

```yaml
## definition.yaml
...
forProvider:
  properties:
    assumeRolePolicyDocument:
      description: Policy that grants an entity permission to assume
        the role.
      type: string
  ...
  required:
    - assumeRolePolicyDocument
  ...
## compsition.yaml
...
patches:
  - fromFieldPath: spec.forProvider.assumeRolePolicyDocument
    policy:
      fromFieldPath: Optional
    toFieldPath: spec.forProvider.assumeRolePolicy
    type: FromCompositeFieldPath
...
```
 If `description` is added, this overrides the decription form the managed resource:
```yaml
overrideFieldsInClaim:
  - claimPath: spec.forProvider.assumeRolePolicyDocument # user sees assumeRolePolicyDocument
    managedPath: spec.forProvider.assumeRolePolicy # managed resource has assumeRolePolicy
    description: The policy document that grants an entity permission to assume
        the role.
```
leads to

```yaml
## definition.yaml
...
forProvider:
  properties:
    assumeRolePolicyDocument:
      description: The policy document that grants an entity permission to assume
        the role.
      type: string
  ...
  required:
    - assumeRolePolicyDocument
  ...
## compsition.yaml
...
patches:
  - fromFieldPath: spec.forProvider.assumeRolePolicyDocument
    policy:
      fromFieldPath: Optional
    toFieldPath: spec.forProvider.assumeRolePolicy
    type: FromCompositeFieldPath
...
```

If not only the name and the description of a property should change, `overrideSettings` with `property` and `patches` can be used:

```yaml
overrideFieldsInClaim:
- claimPath: spec.forProvider.newProp
    managedPath: spec.forProvider.forceDetachPolicies
    overrideSettings:
      property:
        description: New Property
        type: array
        items:
          type: boolean
      patches:
        - fromFieldPath: spec.forProvider.newProp[0]
          policy:
            fromFieldPath: Optional
          toFieldPath: spec.forProvider.forceDetachPolicies
          type: FromCompositeFieldPath
        - fromFieldPath: spec.forProvider.newProp[1]
          policy:
            fromFieldPath: Optional
          toFieldPath: spec.forProvider.forceDetachPolicies
          type: FromCompositeFieldPath
```
leads to
```yaml
## definition.yaml
...
forProvider:
  properties:
    newProp:
      description: New Property
      items:
        type: boolean
      type: array
  ...
## compsition.yaml
...
patches:
  - fromFieldPath: spec.forProvider.newProp[0]
    policy:
      fromFieldPath: Optional
    toFieldPath: spec.forProvider.forceDetachPolicies
    type: FromCompositeFieldPath
  - fromFieldPath: spec.forProvider.newProp[1]
    policy:
      fromFieldPath: Optional
    toFieldPath: spec.forProvider.forceDetachPolicies
    type: FromCompositeFieldPath
...
```

If a property is no longer used by the managed resource but should be keept in the composite and the claim to not break existing resources, 

```yaml
overrideFieldsInClaim:
  - claimPath: spec.forProvider.namedProp
    overrideSettings:
      property:
        description: Deprecated
        type: string
```
leads to
```yaml
## definition.yaml
...
forProvider:
  namedProp:
    description: Deprecated
    type: string
...
## compsition.yaml
...
patches: []
...
```

## Licensing

x-generation is under the Apache 2.0 license.

> The material in `pkg/functions` in x-generation is partially derived from "[crossplane-composition-generator](https://github.com/benagricola/crossplane-composition-generator)"

| Property                       | Function              | Repository  |
|--------------------------------|-----------------------|-------------|
| pkg/functions/                 | generation functions  | [crossplane-composition-generator](https://github.com/benagricola/crossplane-composition-generator) |
| build &                        | submodule for build   | [upbound/build](https://github.com/upbound/build) |
| make e2e                       | kuttl-tests           | [upbound/uptest](https://github.com/upbound/uptest)|


## using pipelelines
When setting the property `usePipeline` to true, x-generation generates composition with `mode: pipeline`. By default the standard function-patch-and-transform function is used to generate the resource and the corresponding patches. Using the prooperty `patchAndTransfromFunction` the function name can be overwritten. The attribute `autoReadyFunction` can be used to configure if and witch autoReady function will be used in the pipeline. The default is:
```yaml
autoReadyFunction:
  generate: true
  name: function-auto-ready
```

Using the property `additionalPipelineSteps` one can configure x-generation to add additional pipeline steps before or after the default patch-and-transform step. Additional pipeline steps need a `step` and a `functionRef.name`. If you want to add the step before the patch-and-transform part, set `bevore` to true, by default the step will be appended. To add a step conditionally, use `condition`, here you can add a CEL string to determine if the composition will include the step. At the moment, the only properties that can be used in the CEL string are tagProperty and tagType, both will be determined by x-generation. To configure the iput of the pipeline step, use the `input` field. Inside the input the placeholders `{tagType}` and `{tagProperty}` will be replaced by the values determined by x-generation.

A example of `additionalPipelineSteps` could look like:

```yaml
usePipeline: true
additionalPipelineSteps: 
  - step: labels
    functionRef:
      name: function-add-labels
    input: 
      apiVersion: labels.fn.crossplane.io/v1beta1
      kind: Input
      labels:
        exclude:
          - "helm.*"
          - "kustomize.*"
          - "crossplane.*"
      annotations:
        ignore: true
  - step: tag
    functionRef:
      name: function-add-tags
    condition: "tagType != 'noTag'"
    input: 
      apiVersion: tags.fn.crossplane.io/v1beta1
      kind: Input
      tagsFrom: metadata.labels
      ignoreTags:
        - kustomize.*
        - crossplane.*
      tags:
        - type: "{tagType}"
          path: "{tagProperty}"
```
