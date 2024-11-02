#!/bin/bash
set -e

# Define log file for debugging
LOG_FILE=".github/e2e_tests/artifacts/logs/kind_cluster_logs/cosi_driver/cleanup_debug.log"
mkdir -p "$(dirname "$LOG_FILE")"  # Ensure the log directory exists

# Error handling function
error_handler() {
  echo "An error occurred during the COSI cleanup. Check the log file for details." | tee -a "$LOG_FILE"
  echo "Failed command: $BASH_COMMAND" | tee -a "$LOG_FILE"
  exit 1
}

# Trap errors and call the error handler
trap 'error_handler' ERR

# Log command execution to the log file for debugging
log_and_run() {
  echo "Running: $*" | tee -a "$LOG_FILE"
  "$@" | tee -a "$LOG_FILE"
}

# Step 1: Remove COSI driver and namespace
log_and_run echo "Removing COSI driver manifests and namespace..."
log_and_run kubectl delete -k . || echo "COSI driver manifests not found." | tee -a "$LOG_FILE"
log_and_run kubectl delete namespace scality-object-storage || echo "Namespace scality-object-storage not found." | tee -a "$LOG_FILE"

# Step 2: Verify namespace deletion
log_and_run echo "Verifying namespace deletion..."
if kubectl get namespace scality-object-storage &>/dev/null; then
  echo "Warning: Namespace scality-object-storage was not deleted." | tee -a "$LOG_FILE"
  exit 1
fi

# Step 3: Delete COSI CRDs
log_and_run echo "Deleting COSI CRDs..."
log_and_run kubectl delete -k github.com/kubernetes-sigs/container-object-storage-interface-api || echo "COSI API CRDs not found." | tee -a "$LOG_FILE"
log_and_run kubectl delete -k github.com/kubernetes-sigs/container-object-storage-interface-controller || echo "COSI Controller CRDs not found." | tee -a "$LOG_FILE"

# Step 4: Verify COSI CRDs deletion
log_and_run echo "Verifying COSI CRDs deletion..."
if kubectl get crd | grep 'container-object-storage-interface' &>/dev/null; then
  echo "Warning: Some COSI CRDs were not deleted." | tee -a "$LOG_FILE"
  exit 1
fi

log_and_run echo "COSI cleanup completed successfully."
