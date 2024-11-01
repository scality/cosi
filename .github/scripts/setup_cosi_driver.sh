#!/bin/bash
set -e

# Error handling function
error_handler() {
  echo "An error occurred during the COSI setup. Exiting."
  exit 1
}

# Trap errors and call the error handler
trap 'error_handler' ERR

# Step 1: Install COSI CRDs
echo "Installing COSI CRDs..."
kubectl create -k github.com/kubernetes-sigs/container-object-storage-interface-api
kubectl create -k github.com/kubernetes-sigs/container-object-storage-interface-controller

# Step 2: Verify COSI Controller Pod Status
echo "Verifying COSI Controller Pod status..."
kubectl wait --namespace default --for=condition=ready pod -l app.kubernetes.io/name=container-object-storage-interface-controller --timeout=10s
kubectl get pods --namespace default

# Step 3: Build COSI driver Docker image
echo "Building COSI driver image..."
docker build -t ghcr.io/scality/cosi:latest .

# Step 4: Load COSI driver image into KIND cluster
echo "Loading COSI driver image into KIND cluster..."
kind load docker-image ghcr.io/scality/cosi:latest --name object-storage-cluster

# Step 5: Run COSI driver
echo "Applying COSI driver manifests..."
kubectl apply -k .

# Step 6: Verify COSI driver Pod Status
echo "Verifying COSI driver Pod status..."
kubectl wait --namespace scality-object-storage --for=condition=ready pod --selector=app.kubernetes.io/name=scality-cosi-driver --timeout=20s
kubectl get pods -n scality-object-storage

echo "COSI setup completed successfully."
