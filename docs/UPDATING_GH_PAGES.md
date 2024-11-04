# Updating the `gh-pages` Branch for Scality COSI Driver Helm Chart

This guide provides detailed steps to update the `gh-pages` branch with a new version of the Scality COSI Driver Helm chart, enabling it to be accessible via GitHub Pages as a Helm repository.

## Prerequisites

- Ensure you have a new version of the chart ready in the main branch (e.g., updated `Chart.yaml`).
- Make sure you are familiar with creating Helm package files and updating `index.yaml`.

## Step-by-Step Instructions

### Step 1: Package the New Chart Version

Switch to your main branch where the Helm chart files are located and package the chart:

```bash
# Ensure you are on the main branch
git checkout main

# Package the Helm chart
helm package helm/scality-cosi-driver
```

This command will create a `.tgz` file (e.g., `scality-cosi-driver-0.1.0-beta-PW.8.tgz`) in the current directory.

### Step 2: Switch to the `gh-pages` Branch

Switch to the `gh-pages` branch, where the Helm chart repository files are hosted:

```bash
git checkout gh-pages
```

### Step 3: Move the Packaged Chart to the `gh-pages` Branch

Move the `.tgz` file generated in the main branch to the `gh-pages` branch directory:

```bash
# Move the packaged chart file
mv ../scality-cosi-driver-0.1.0-beta-PW.8.tgz .
```

### Step 4: Update `index.yaml`

Update the `index.yaml` file to include the new chart version. Run the following command in the `gh-pages` branch:

```bash
helm repo index . --url https://scality.github.io/cosi/
```

This command will:

- Update `index.yaml` with the new `.tgz` file.
- Ensure the Helm repository is up-to-date with the latest chart version.

### Step 5: Commit and Push Changes to `gh-pages`

Add, commit, and push the updated `index.yaml` and `.tgz` file to the `gh-pages` branch:

```bash
# Add the new chart package and index.yaml to staging
git add scality-cosi-driver-0.1.0-beta-PW.8.tgz index.yaml

# Commit the changes
git commit -m "Update Helm chart to version 0.1.0-beta-PW.8"

# Push the changes to the remote gh-pages branch
git push origin gh-pages
```

### Step 6: Verify the Update

Once the changes are pushed, verify that the new version is available in your Helm repository.

1. **Refresh the Helm repository** on your local machine:

   ```bash
   helm repo update
   ```

2. **Search for the updated chart version**:

   ```bash
   helm search repo scality-cosi-driver
   ```

3. **Install the new chart version**:

   ```bash
   helm install my-release scality-cosi-driver/scality-cosi-driver --version 0.1.0-beta-PW.8
   ```

This process ensures that the latest chart version is published to GitHub Pages and accessible as a Helm repository.
