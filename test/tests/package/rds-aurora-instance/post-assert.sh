#!/usr/bin/env bash
set -aeuo pipefail

SED=$(which gsed || which sed)
parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
KUTTL=$2

AURORA="$(kubectl get dbinstance.rds.aws.crossplane.io -o=jsonpath='{.items[?(@.metadata.labels.crossplane\.io/claim-name=="example-aurora-mysql-instance")].metadata.name}')"

$SED \
  -e "s/AURORA/${AURORA}/g" \
  ${parent_path}/assert.yaml > ${parent_path}/zz-generated-assert.yaml

$2 assert ${parent_path}/zz-generated-assert.yaml
