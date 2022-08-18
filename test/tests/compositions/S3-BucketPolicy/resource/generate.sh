#!/bin/sh
SED=$(which gsed || which sed)

sleep 30

POLICY="$(kubectl get bucketpolicy.s3.aws.crossplane.io -o=jsonpath='{.items[?(@.metadata.labels.crossplane\.io/claim-name=="example-bucket-policy")].metadata.name}')"

$SED \
  -e "s/POLICY/${POLICY}/g" \
  ./resource/00-assert.yaml > ./resource/zz-generated-assert.yaml
