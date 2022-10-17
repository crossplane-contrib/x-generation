#!/usr/bin/env bash
set -aeuo pipefail

SED=$(which gsed || which sed)
parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
KUTTL=$2

GENERATED="$(kubectl get compositepolicy --no-headers -o custom-columns=:metadata.name | head -n1)"

$SED \
  -e "s/xxxxx/${GENERATED}/g" \
  ${parent_path}/assert.yaml > ${parent_path}/zz-generated-assert.yaml

$2 assert ${parent_path}/zz-generated-assert.yaml
