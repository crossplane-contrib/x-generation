#!/bin/sh
SED=$(which gsed || which sed)

sleep 90

INSTANCE="$(kubectl get dbinstance.rds.aws.crossplane.io -o=jsonpath='{.items[?(@.metadata.labels.crossplane\.io/claim-name=="example-dbinstance")].metadata.name}')"
AURORA="$(kubectl get dbinstance.rds.aws.crossplane.io -o=jsonpath='{.items[?(@.metadata.labels.crossplane\.io/claim-name=="example-aurora-mysql-instance")].metadata.name}')"

$SED \
  -e "s/INSTANCE/${INSTANCE}/g" \
  -e "s/AURORA/${AURORA}/g" \
  ./resource/00-assert.yaml > ./resource/zz-generated-assert.yaml
