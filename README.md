# X-GENERATION

generate compositions from crossplane provider crds

## configure

The generation of crds can be configured in two places, either in the global configuration file, or in the local generation file for each composition.
### global configuration
In the local configuration, the provider used and labels and tags that should be used for all generated crds can be configured in a global configuration. The default name for this file is `generator-config.yaml`, the filename can be overwritten using the `--configFile` flag. 

| Property         | Type              | Description |
|------------------|-------------------|-------------|
| provider         | object            | Object used to configure the provider used for the generation |
| provider.baseURL | string            | The url globaly used to retrieve the crds needed for generating the compositions, three placeholders are provided during the generation of compositions: The name of the provider, the version of the provider and the crd file name|
| provider.name    | string            | The name of the provider |
| provider.version | string            | The version of the provider |
| labels           | object            | Configure the labels and label patches for each crd |
| labels.fromCRD   | array of strings  | For each entry `e` a patch that copies the value of the `metadata.labels[e]` field from the CompositeResourceDefinition to the same field of the resource |
| labels.common    | object of strings | For each property of the object a label with the given value is created in the resource |
| tags             | object            | Configure the tags and tag patches for each crd |
| tags.fromLabels  | array of strings  | For each entry `e` a patch that copies the value of the `metadata.labels[e]` field to a tag with the same name and value is created
| tags.common      | object of strings | For each property of the object a tag with the given value is created in the resource |


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

## tests
```
mkdir test
kubectl kuttl test
```

