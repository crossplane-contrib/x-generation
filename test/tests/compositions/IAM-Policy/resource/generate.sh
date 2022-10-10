#!/bin/sh
SED=$(which gsed || which sed)

sleep 30
GENERATED="$(kubectl get compositepolicy --no-headers -o custom-columns=:metadata.name | head -n1)"

$SED \
  -e "s/xxxxx/${GENERATED}/g" \
  ./resource/00-assert.yaml > ./resource/zz-generated-assert.yaml
