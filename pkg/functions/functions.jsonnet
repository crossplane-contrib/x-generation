{
  local defaultUIDFieldPath = 'metadata.annotations["crossplane.io/external-name"]',

  local defaultIgnores = [
    'status.conditions',
    'spec.writeConnectionSecretToRef',
    'spec.forProvider.tags',
    'spec.forProvider.tagSpecifications',
    'spec.forProvider.tagging',
    'spec.providerConfigRef.default',
    'spec.providerRef',
    'spec.publishConnectionDetailsTo.configRef.default'
  ],

  NameToPlural(config):: (

    if std.objectHas(config, "plural") then
      std.asciiLower(config.plural)
    else
      local lname = std.asciiLower(config.name);
      local length = std.length(lname);
      local last = std.substr(lname, length - 1, 1);
      local replace = if last == 'y' then 'ie' else last;
      std.substr(lname,0,length -1) + replace + 's'
  ),
  FQDN(name, group):: (
    '%s.%s' % [std.asciiLower(name), group]
  ),
  GenerateCategories(group):: (
    ['crossplane', 'composition', std.split(group, '.')[0]]
  ),
  GenerateLabels(compositionIdentifier, provider):: (
    {
      [compositionIdentifier+'/provider']: provider,
    }
  ),
  local labelize(fqdn) = (
    "metadata.labels['%s']" % [fqdn]
  ),
  local genExternalGenericLabel(fields) = (
    [labelize(field) for field in fields]
  ),
  local genPatch(type, fieldFrom, fieldTo, srcFieldName, dstFieldName, policy) = (
    {
      [srcFieldName]: fieldFrom,
      [dstFieldName]: fieldTo,
      policy: {
        fromFieldPath: policy,
      },
      type: type,
    }
  ),
  GenPatch(type, fieldFrom, fieldTo, srcFieldName, dstFieldName, policy):: (
    [genPatch(type, fieldFrom, fieldTo, srcFieldName, dstFieldName, policy)]
  ),
  local genSecretPatch(type, fieldFrom, fieldTo, srcFieldName, dstFieldName, policy) = (
    {
      [srcFieldName]: fieldFrom,
      [dstFieldName]: fieldTo,
      policy: {
        fromFieldPath: policy,
      },
      type: type,
      transforms: [
        {
          type: "string",
          string: {
            fmt: '%s-secret'
          },
        },
      ],
    }
  ),
  GenSecretPatch(type, fieldFrom, fieldTo, srcFieldName, dstFieldName, policy):: (
    [genSecretPatch(type, fieldFrom, fieldTo, srcFieldName, dstFieldName, policy)]
  ),
   local genOptionalPatchFromConfig(fields, config) = (
      local on = overrideNamesByClaim(config);
   [
      local fieldTo = if field in on then on[field].managedPath else field;
      genPatch('FromCompositeFieldPath', field, fieldTo, 'fromFieldPath', 'toFieldPath', 'Optional')
      for field in std.filter(function(f) (!(f in on && "overrideSettings" in on[f])), fields)
      ] + std.flattenArrays([ on[field].overrideSettings.patches for field in std.filter(function(f) (f in on && "overrideSettings" in on[f] && "managedPath" in on[f]), fields) ])
  ),
  local genOptionalPatchFrom(fields) = (
    [
      genPatch('FromCompositeFieldPath', field, field, 'fromFieldPath', 'toFieldPath', 'Optional')
      for field in fields
    ]
  ),
  GenOptionalPatchFrom(fields, config):: (
    genOptionalPatchFromConfig(fields, config)
  ),
  GenOptionalPatchTo(fields):: (
    [
      genPatch('ToCompositeFieldPath', field, field, 'toFieldPath', 'fromFieldPath', 'Optional')
      for field in fields
    ]
  ),
  GetDefaultComposition(compositions):: (
    local default = [c.name for c in compositions if 'default' in c && c.default];
    assert std.length(default) == 1 : 'Could not find a default composition. One composition must have default: true!';
    default[0]
  ),
  GetVersion(crd, version):: (
    local fv = [v for v in crd.spec.versions if 'name' in v && v.name == version];
    assert std.length(fv) == 1 : 'Could not find CRD with version %s' % version;
    fv[0]
  ),
  local overrides(config) = {
    [o.path]: o.override
    for o in config.overrideFields
    if 'override' in o
  },
  local last(s) = (
    local list = splitPath(s);
    local len = std.length(list);
    {newName: list[len -1 ]}
  ),
  local getOverrides(o) = {
    override:
      if "overrideSettings" in o && "property" in o.overrideSettings  then
       {}+ o.overrideSettings.property
       else if "description" in o then
       {description: o.description}
      else {}
  },
   local getPath(s) = (
      local list = splitPath(s);
      local len = std.length(list);
      local path = std.join(".", list[0:len-1]);
      path
    ),
  local filteredNewEntries = function(config) (std.filter(function(c) (!std.objectHas(c, 'managedPath')), config.overrideFieldsInClaim)),
  local filteredOverrideEntries = function(config) (std.filter(function(c) (std.objectHas(c, 'managedPath')), config.overrideFieldsInClaim)),
  local newProps(config) = {
    local filterForPath = function(path, config) [
        last(p.claimPath) + getOverrides(p)
        for p in std.filter(function(c) (getPath(c.claimPath) == path), filteredNewEntries(config))
    ],
    [getPath(o.claimPath)]: filterForPath(getPath(o.claimPath), config)
    for o in filteredNewEntries(config)
    //["a"]: std.toString(filteredNewEntries(config))
    // local newMap = {};
    // local filMap = (map, path, e) (
    //  if std.objectHas(path) then
    //   map[path] = map[path] + [last(e.claimPath) + getOverrides(e)]
    // else
    //   map[path]: [last(e.claimPath) + getOverrides(e)]
    // )
    //   local path = getPath(o.claimPath)
    //   filMap(newMap,path,o)
    // for o in filteredNewEntries,
    // newMap
  },
  local overrideNames(config) = {
    // local last(s) = {
    //   local list = std.split(s, "."),
    //   local len = std.length(list),

    //   {newName: list[len -1]}
    // };
    [o.managedPath]: o + last(o.claimPath)+
    { override:
      if "overrideSettings" in o && "property" in o.overrideSettings  then
       {}+ o.overrideSettings.property
       else if "description" in o then

       {description: o.description}

      else {}
    }
    for o in filteredOverrideEntries(config)
  },
  local overrideNamesByClaim(config) = {
    [o.claimPath]: o
    for o in config.overrideFieldsInClaim
  },
  local ignores(config) = [
    o.path
    for o in config.overrideFields
    if 'ignore' in o && o.ignore
  ] + defaultIgnores,
  local values(config) = {
    [o.path]: o.value
    for o in config.overrideFields
    if 'value' in o
  },
  local joinPath(path, item) = (
    std.join('.', path + [item])
  ),
  local splitPath(path) = (
    std.split(path, '.')
  ),
  local updatePath(path, name) = (
    if name != 'properties' then
      path + [name]
    else
      path
  ),
  local recurseWithPath(init, obj, path, foldFunc, filterFunc, valueFunc, parent) = (
    if std.isObject(obj) then
      std.foldl(foldFunc, std.filterMap(
        function(n) (
          filterFunc(obj, path, n)
        ),
        function(n) (
          valueFunc(obj, path, n, recurseWithPath, parent)
        ),
        std.objectFields(obj)
      ), init)
    else
      valueFunc(obj, path, path[std.length(path) - 1], recurseWithPath, parent)
  ),
  GenerateSchema(obj, config, path):: (
    local ignorePaths = ignores(config);
    local overridePaths = overrides(config);

    local foldFunc = function(r, n) r + n;

    local filterFunc = function(object, fullPath, name) (
      !std.member(ignorePaths, joinPath(fullPath, name))
    );

    local overrideNamesPath = overrideNames(config);

    local newProperties = newProps(config);

    local ignorPathForRequired = function(ignorePath, reqired) (

      local aux = function(arr , index )  (
        local elem = arr[index];
        if index == std.length(arr) - 2 then
          elem
        else
          elem + "." + aux(arr, index + 1)
      );
      aux(ignorePath,0)
    );

    local valueFunc = function(object, fullPath, name, recurse, parent) (

      local stringPath = std.join(".", fullPath);
      if !std.isObject(object) then
        if name == "required" then
        // filter ignored values from required array
        local ignorePath = ignorPathForRequired(fullPath,object);
        std.map(function(e) (
          local jp = joinPath(fullPath[0:std.length(fullPath) -1 ], e);
          if jp in overrideNamesPath then
            overrideNamesPath[jp].newName
          else
           e
        ),
        std.filterMap(
            function(n) (
              !std.member(ignorePaths, ignorePath+"."+n)
            ),
            function(n) (
              n
            ),
            object,
          ))
        else
          object
      else
        local jp = joinPath(fullPath, name);
        if jp in overrideNamesPath then
          { [overrideNamesPath[jp].newName]+: object[name] +  overrideNamesPath[jp].override}
        else if jp in overridePaths then
          { [name]+: object[name] + overridePaths[jp] }
        else
          { [name]+: recurse(
            {},
            object[name],
            updatePath(fullPath, name),
            foldFunc,
            filterFunc,
            valueFunc,
            name
          )
          }
          +
          if std.objectHas(newProperties, stringPath)  && parent == "properties" then
          {


          //"newProp": std.toString(newProperties[fullPath]) + fullPath
          [p.newName]: p.override
          for p in newProperties[stringPath]
        // [a]: a for a in newProperties
          } else {}

  // {
  //   [prop.newName]: prop.override
  // }
  //  for prop in newProperties[fullPath]
    );
    recurseWithPath({}, obj, path, foldFunc, filterFunc, valueFunc, "")
  ),
  GeneratePatchPaths(obj, config, path):: (
    local filterFunc = function(object, fullPath, name) (
      std.isObject(object[name])
    );
    local foldFunc = function(r, n) (
      if std.isArray(n) then
        r + std.prune(n)
      else
        r + [n]
    );
    local valueFunc = function(object, fullPath, name, recurse, parent) (
      local o = object[name];
      if 'type' in o && o.type == 'object' && !('additionalProperties' in object[name]) || (name == 'properties' && parent != 'properties') then
        recurse([], object[name], updatePath(fullPath, name), foldFunc, filterFunc, valueFunc, name)
      else
        if (name != 'default' || parent == 'properties') then
          joinPath(fullPath, name)
    );
    recurseWithPath([], obj, path, foldFunc, filterFunc, valueFunc,"")
  ),
  SetDefaults(config):: (
    local defaultValues = values(config);
    std.foldl(function(a, b) a + b, std.map(function(key) (
      local sp = splitPath(key);
      local sl = std.length(sp) - 1;
      std.foldr(function(r, n) (
        { [r]+: n }
      ), sp[0:sl], { [sp[sl]]: defaultValues[key] })
    ), std.objectFields(defaultValues)), {})
  ),
  FilterPrinterColumns(columns):: (
    std.filter(function(c) !std.startsWith(c.jsonPath, '.status.conditions'), columns)
  ),
  GetUIDFieldPath(config):: (
    if 'uidFieldPath' in config then
      config.uidFieldPath
    else
      defaultUIDFieldPath
  ),
  GenTagKeys(tagType, tagProperty, tags, commonTags):: (
    local tagProp = if tagProperty == "tag" then "tags" else if tagProperty == "tagSet" then "tagSet";
    local generatedTags = {
      [if tagType == "keyValueArray" then tagProp]: [{
        key: tag
      } for tag in tags ]
      +
      [{
        key: tag,
        value: commonTags[tag],
      } for tag in std.objectFields(commonTags) ],
      [if tagType == "tagKeyValueArray" then tagProp]: [{
        tagKey: tag
      } for tag in tags ]
      +
      [{
        tagKey: tag,
        tagValue: commonTags[tag],
      } for tag in std.objectFields(commonTags) ],
      [if tagType == "stringObject" && std.length(commonTags) > 0 then "tags"]: {
        [tag]: commonTags[tag],
      for tag in std.objectFields(commonTags) }
    };
    if tagProperty == "tag" then generatedTags
    else if tagProperty == "tagSet" then {
      tagging: generatedTags
    }
    else {}
  ),
  GenTagsPatch(tagType, tags, tagProperty):: (
  local tagProp = if tagProperty == "tag" then "tags" else if tagProperty == "tagSet" then "tagging.tagSet";
  if  tagType != "" then [
    {
      name: "Tags",
      patches: if  tagType == "keyValueArray" then [
        genPatch('FromCompositeFieldPath', "metadata.labels["+tags[f]+"]", "spec.forProvider."+tagProp+"["+f+"].value", 'fromFieldPath', 'toFieldPath', "Required")
        for f in std.range(0, std.length(tags)-1)
      ] else if  tagType == "tagKeyValueArray" then [
        genPatch('FromCompositeFieldPath', "metadata.labels["+tags[f]+"]", "spec.forProvider."+tagProp+"["+f+"].tagValue", 'fromFieldPath', 'toFieldPath', "Required")
        for f in std.range(0, std.length(tags)-1)
      ] else if  tagType == "stringObject" then [
        genPatch('FromCompositeFieldPath', "metadata.labels["+tag+"]", "spec.forProvider."+tagProp+"["+tag+"]", 'fromFieldPath', 'toFieldPath', 'Optional')
        for tag in tags
      ]
    }
   ] else []
  ),
  GenLabelsPatch(labelList):: (
    genOptionalPatchFrom(genExternalGenericLabel(labelList))
  ),
  GenCommonLabels(commonLabels):: (
    {
    [if std.length(commonLabels) > 0 then "labels"]: {
        [label]: commonLabels[label] for label in std.objectFields(commonLabels)
      }
    }
  ),
}
