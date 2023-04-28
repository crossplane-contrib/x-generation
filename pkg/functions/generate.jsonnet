local k8s = import 'functions.jsonnet';

local s = {
  config: std.parseJson(std.extVar('config')),
  crd: std.parseJson(std.extVar('crd')),
  data: std.parseJson(std.extVar('data')),
  tagList: std.parseJson(std.extVar('tagList')),
  tagType: std.extVar('tagType'),
  tagProperty: std.extVar('tagProperty'),
  commonTags: std.parseJson(std.extVar('commonTags')),
  labelList: std.parseJson(std.extVar('labelList')),
  commonLabels: std.parseJson(std.extVar('commonLabels')),
  globalLabels: std.parseJson(std.extVar('globalLabels')),
  compositionIdentifier: std.extVar('compositionIdentifier'),
  readinessChecks: std.extVar('readinessChecks'),
};

local plural = k8s.NameToPlural(s.config);
local fqdn = k8s.FQDN(plural, s.config.group);
local resourceFqdn = k8s.FQDN(s.crd.names.kind, s.crd.group);
local version = k8s.GetVersion(s.crd, s.config.provider.crd.version);

local uidFieldPath = k8s.GetUIDFieldPath(s.config);
local uidFieldName = 'uid';

local definitionSpec = k8s.GenerateSchema(
  version.schema.openAPIV3Schema.properties.spec,
  s.config,
  ['spec'],
);

local definitionStatus = k8s.GenerateSchema(
  version.schema.openAPIV3Schema.properties.status,
  s.config,
  ['status'],
);

{
  definition: {
    apiVersion: 'apiextensions.crossplane.io/v1',
    kind: 'CompositeResourceDefinition',
    metadata: {
      name: "composite"+fqdn,
    },
    spec: {
      claimNames: {
        kind: s.config.name,
        plural: plural,
      },
      [if std.objectHas(s.config, "connectionSecretKeys") then "connectionSecretKeys"]:
        s.config.connectionSecretKeys,
      group: s.config.group,
      names: {
        kind: "Composite"+s.config.name,
        plural: "composite"+plural,
        categories: k8s.GenerateCategories(s.config.group),
      },
      versions: [
        {
          name: s.config.version,
          referenceable: version.storage,
          served: version.served,
          schema: {
            openAPIV3Schema: {
              properties: {
                spec: definitionSpec,
                status:
                  definitionStatus
                  {
                    properties+: {
                      [uidFieldName]: {
                        description: 'The unique ID of this %s resource reported by the provider' % [s.config.name],
                        type: 'string',
                      },
                      observed: {
                        description: 'Freeform field containing information about the observed status.',
                        type: 'object',
                        "x-kubernetes-preserve-unknown-fields": true,
                      },
                    },
                  },
              },
            },
          },
          additionalPrinterColumns: k8s.FilterPrinterColumns(version.additionalPrinterColumns),
        },
      ],
      // defaultCompositionRef: {
      //   name: k8s.GetDefaultComposition(s.config.compositions),
      // },
    },
  },
} + {
  ['composition-' + composition.name]: {
    apiVersion: 'apiextensions.crossplane.io/v1',
    kind: 'Composition',
    metadata: {
      name: "composite" + composition.name + "." + s.config.group,
      labels: k8s.GenerateLabels(s.compositionIdentifier,composition.provider),
    },
    spec: {
      local spec = self,
      [if std.objectHas(s.config, "connectionSecretKeys") then "writeConnectionSecretsToNamespace"]:
        'crossplane-system',
      compositeTypeRef: {
        apiVersion: s.config.group + '/' + s.config.version,
        kind: "Composite"+s.config.name,
      },
      patchSets: [
        {
          name: 'Name',
          patches: [{
            type: 'FromCompositeFieldPath',
            fromFieldPath: 'metadata.labels[crossplane.io/claim-name]',
            toFieldPath: if std.objectHas(s.config, 'patchExternalName') && s.config.patchExternalName == false then 'metadata.name' else 'metadata.annotations[crossplane.io/external-name]',
          }],
        },
        {
          name: 'Common',
          patches: k8s.GenLabelsPatch(s.globalLabels)
        },
        {
          name: 'Parameters',
          patches: k8s.GenOptionalPatchFrom(
            k8s.GeneratePatchPaths(
              definitionSpec.properties,
              s.config,
              ['spec']
            ),
            s.config
          ),
        },
        {
          name: 'Labels',
          patches: k8s.GenLabelsPatch(s.labelList)
        }
      ] + k8s.GenTagsPatch(s.tagType, s.tagList, s.tagProperty),
      resources: [
        {
          local resource = self,
          name: s.crd.spec.names.kind,
          base: {
            apiVersion: s.crd.spec.group + '/' + s.config.provider.crd.version,
            kind: resource.name,
            metadata: k8s.GenCommonLabels(s.commonLabels),
            spec: {
              providerConfigRef: {
                name: 'default',
              },
              [if std.objectHas(s.config, "connectionSecretKeys") then "writeConnectionSecretToRef"]:
                {
                  namespace: 'crossplane-system'
                },
              forProvider: k8s.GenTagKeys(s.tagType, s.tagProperty, s.tagList, s.commonTags)
            },
          } + k8s.SetDefaults(s.config),
          patches: [
            {
              type: 'PatchSet',
              patchSetName: ps.name,
            }
            for ps in spec.patchSets
          ] + k8s.GenOptionalPatchTo(
            k8s.GeneratePatchPaths(
              definitionStatus.properties,
              s.config,
              ['status']
            )
          )+ k8s.GenPatch(
              'ToCompositeFieldPath',
              uidFieldPath,
              'status.%s' % [uidFieldName],
              'fromFieldPath',
              'toFieldPath',
              'Optional'
          )+ k8s.GenPatch(
              'ToCompositeFieldPath',
              'status.conditions',
              'status.observed.conditions',
              'fromFieldPath',
              'toFieldPath',
              'Optional'
          )+
          (if std.objectHas(s.config, "connectionSecretKeys") then           
            k8s.GenSecretPatch(
                'FromCompositeFieldPath',
                'metadata.uid',
                'spec.writeConnectionSecretToRef.name',
                'fromFieldPath',
                'toFieldPath',
                'Optional'
            )else []),
          [if s.readinessChecks == "false" then "readinessChecks"]: [{type:"None"}],
          [if std.objectHas(s.config, "connectionSecretKeys") then "connectionDetails"]:
            [
              {
                fromConnectionSecretKey: keys,
              },
              for keys in s.config.connectionSecretKeys
            ],
        },
      ],
    },
  }
  for composition in s.config.compositions
}
