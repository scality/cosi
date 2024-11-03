#!/bin/bash
set -e

NAMESPACE="scality-object-storage"

echo "Verifying Helm installation..."

# Check that Helm release exists
if ! helm status scality-cosi-driver -n $NAMESPACE; then
  echo "Helm release scality-cosi-driver not found in namespace $NAMESPACE"
  exit 1
fi

# Check that all pods are running
if ! kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=scality-cosi-driver -n $NAMESPACE --timeout=120s; then
  echo "One or more pods failed to start within the expected time"
  exit 1
fi

echo "Helm installation verified successfully."
