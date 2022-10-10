#!/bin/sh
SED=$(which gsed || which sed)

sleep 30
GENERATED="$(kubectl get compositerole --no-headers -o custom-columns=:metadata.name | head -n1)"
CROSSPLANE_NAME="$(kubectl get roles.iam.aws.crossplane.io --no-headers -o custom-columns=:metadata.name | head -n1)"

$SED \
  -e "s/xxxxx/${GENERATED}/g" \
  -e "s/CROSSPLANE_NAME/${CROSSPLANE_NAME}/g" \
  ./resource/00-assert.yaml > ./resource/zz-generated-assert.yaml