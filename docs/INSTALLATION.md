# Scality COSI Driver Helm Chart

This README provides instructions to add the Helm repository, install the Scality COSI Driver chart, uninstall it, and clean up resources. These commands are designed to be copy-paste ready.

## Adding the Helm Repository

To add the Scality COSI Driver Helm repository:

```bash
helm repo add scality-cosi-driver https://scality.github.io/cosi/
helm repo update
```

This command will add the `scality-cosi-driver` repository and fetch the latest metadata.

## Installing the Scality COSI Driver Chart

To install the Scality COSI Driver chart from the repository, use the following command:

```bash
helm install my-release scality-cosi-driver/scality-cosi-driver --version 0.1.0-beta-PW.10 --namespace scality-object-storage --create-namespace
```

Replace `my-release` with your desired release name. This command will:

- Install the chart with the specified version (`0.1.0-beta-PW.7`).
- Create the `scality-object-storage` namespace if it doesn’t already exist.

### Customizing the Installation

You can customize values using the `--set` flag or a `values.yaml` file.

For example, to specify a custom image tag:

```bash
helm install my-release scality-cosi-driver/scality-cosi-driver \
  --version 0.1.0-beta-PW.7 \
  --namespace scality-object-storage \
  --create-namespace \
  --set image.tag=latest
```

## Verifying the Installation

To verify that the release is installed and running:

```bash
# Check the Helm release status
helm list -n scality-object-storage

# Check the pods in the namespace
kubectl get pods -n scality-object-storage
```

## Uninstalling the Scality COSI Driver Chart

To uninstall the Helm release and delete all associated resources within the namespace:

```bash
helm uninstall my-release --namespace scality-object-storage
```

Replace `my-release` with the name of your release. This command will remove all resources managed by the Helm release within the specified namespace.

### Optional: Delete the Namespace

If you want to delete the namespace as well, use:

```bash
kubectl delete namespace scality-object-storage
```

> **Note**: Only delete the namespace if you are sure that no other resources are using it, as this will remove all resources within that namespace.

## Removing the Helm Repository

To remove the Scality COSI Driver Helm repository from your local setup:

```bash
helm repo remove scality-cosi-driver
```

You can verify that the repository has been removed by listing your Helm repositories:

```bash
helm repo list
```
