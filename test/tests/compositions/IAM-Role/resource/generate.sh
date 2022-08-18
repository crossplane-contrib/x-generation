#!/bin/sh
SED=$(which gsed || which sed)

$SED \
 -e "s/xxxxx/xxxxx/g" \
 ./resource/00-assert.yaml > ./resource/zz-generated-assert.yaml
