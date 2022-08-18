#!/bin/sh
SED=$(which gsed || which sed)

sleep 30

BUCKET="$(kubectl get bucket.s3.aws.crossplane.io -o=jsonpath='{.items[?(@.metadata.labels.crossplane\.io/claim-name=="example-bucket")].metadata.name}')"

$SED \
  -e "s/BUCKET/${BUCKET}/g" \
  ./resource/00-assert.yaml > ./resource/zz-generated-assert.yaml
