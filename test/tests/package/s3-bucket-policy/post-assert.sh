#!/usr/bin/env bash
set -aeuo pipefail

SED=$(which gsed || which sed)
parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
KUTTL=$2

POLICY="$(kubectl get bucketpolicy.s3.aws.crossplane.io -o=jsonpath='{.items[?(@.metadata.labels.crossplane\.io/claim-name=="example-bucket-policy")].metadata.name}')"

$SED \
  -e "s/POLICY/${POLICY}/g" \
  ${parent_path}/assert.yaml > ${parent_path}/zz-generated-assert.yaml

$2 assert ${parent_path}/zz-generated-assert.yaml
