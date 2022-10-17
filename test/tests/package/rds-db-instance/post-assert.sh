#!/usr/bin/env bash
set -aeuo pipefail

SED=$(which gsed || which sed)
parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
KUTTL=$2

INSTANCE="$(kubectl get dbinstance.rds.aws.crossplane.io -o=jsonpath='{.items[?(@.metadata.labels.crossplane\.io/claim-name=="example-dbinstance")].metadata.name}')"

$SED \
  -e "s/INSTANCE/${INSTANCE}/g" \
  ${parent_path}/assert.yaml > ${parent_path}/zz-generated-assert.yaml

$2 assert ${parent_path}/zz-generated-assert.yaml
