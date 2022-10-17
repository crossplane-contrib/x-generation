#!/usr/bin/env bash
set -aeuo pipefail

echo "Running setup.sh"
echo "Waiting until configuration package is healthy/installed..."
${KUBECTL} wait configuration.pkg x-generation --for=condition=Healthy --timeout 5m
${KUBECTL} wait configuration.pkg x-generation --for=condition=Installed --timeout 5m
${KUBECTL} wait "provider.pkg.crossplane.io/crossplane-contrib-provider-aws" --for=condition=healthy --timeout=300s
${KUBECTL} wait --for=condition=established --timeout=300s crd/providerconfigs.aws.crossplane.io

echo "Creating a secret default provider config"
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: aws-creds
  namespace: upbound-system
type: Opaque
stringData:
  key: nocreds
EOF

echo "Creating a default provider config"
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: aws.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: 123456789101-example
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: upbound-system
      name: aws-creds
      key: key
EOF

${KUBECTL} create ns assert-db
${KUBECTL} create ns assert-bucket
${KUBECTL} create ns assert-bucket-policy
${KUBECTL} create ns assert-kms
${KUBECTL} create ns assert-kms-alias
${KUBECTL} create ns assert-instance-profile
${KUBECTL} create ns assert-iam-policy
${KUBECTL} create ns assert-role-policy-attachment
${KUBECTL} create ns assert-role
