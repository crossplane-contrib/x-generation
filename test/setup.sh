#!/usr/bin/env bash
set -aeuo pipefail

echo "Running setup.sh"
echo "Waiting until configuration package is healthy/installed..."
${KUBECTL} wait configuration.pkg x-generation --for=condition=Healthy --timeout 5m
${KUBECTL} wait configuration.pkg x-generation --for=condition=Installed --timeout 5m
${KUBECTL} wait "provider.pkg.crossplane.io/crossplane-contrib-provider-aws" --for=condition=healthy --timeout=300s
${KUBECTL} wait --for=condition=established --timeout=300s crd/providerconfigs.aws.crossplane.io
sleep 5
${KUBECTL} wait "provider.pkg.crossplane.io/crossplane-contrib-provider-zpa" --for=condition=healthy --timeout=300s
${KUBECTL} wait --for=condition=established --timeout=300s crd/providerconfigs.zpa.crossplane.io

${KUBECTL} wait "provider.pkg.crossplane.io/upbound-release-candidates-provider-aws-iam" --for=condition=healthy --timeout=300s
${KUBECTL} wait --for=condition=established --timeout=300s crd/providerconfigs.aws.upbound.io
sleep 60

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
echo "Creating a secret default provider zpa config"
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: zpa-creds
  namespace: upbound-system
type: Opaque
stringData:
  key: nocreds
EOF

echo "Creating a default provider zpa config"
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: zpa.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: 123456789101-example
spec:
  clientID:
    source: None
  clientSecret:
    source: None
  customerId: testID
  host: zpa.local
EOF

echo "Creating a secret default offical provider awsconfig"
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: offical-aws-creds
  namespace: upbound-system
type: Opaque
stringData:
  key: nocreds
EOF

echo "Creating a default provider config"
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: aws.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: 123456789101-example
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: upbound-system
      name: offical-aws-creds
      key: key
EOF

${KUBECTL} create ns assert-db || true
${KUBECTL} create ns assert-bucket || true
${KUBECTL} create ns assert-bucket-policy || true
${KUBECTL} create ns assert-kms || true
${KUBECTL} create ns assert-kms-alias || true
${KUBECTL} create ns assert-instance-profile || true
${KUBECTL} create ns assert-iam-policy || true
${KUBECTL} create ns assert-role-policy-attachment || true
${KUBECTL} create ns assert-role || true
${KUBECTL} create ns assert-zpa || true
