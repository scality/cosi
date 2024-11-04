# Scality COSI Driver

The Scality COSI (Container Object Storage Interface) Driver Helm chart allows you to deploy and manage the Scality COSI Driver in Kubernetes. This driver is designed to integrate with the Kubernetes ecosystem, enabling seamless object storage provisioning.

## Table of Contents

- [Scality COSI Driver](#scality-cosi-driver)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Installation Guide](#installation-guide)
  - [Updating the `gh-pages` Branch](#updating-the-gh-pages-branch)
  - [Contributing](#contributing)
  - [License](#license)

## Introduction

The Scality COSI Driver provides a Kubernetes-native interface for managing object storage. Using this Helm chart, you can deploy the driver and manage its configuration across your Kubernetes clusters.

## Installation Guide

To install the Scality COSI Driver Helm chart and manage your deployment, please refer to the [Installation Guide](docs/INSTALLATION.md).

This guide covers:

- Adding the Helm repository
- Installing and customizing the Helm chart
- Uninstalling the Helm chart
- Cleaning up the Helm repository

## Updating the `gh-pages` Branch

If you're a maintainer and need to update the `gh-pages` branch with a new version of the Helm chart, refer to the [Updating the `gh-pages` Branch Guide](docs/UPDATING_GH_PAGES.md).

This guide includes:

- Steps to package a new Helm chart version
- Instructions for updating `index.yaml`
- Commands to push the updated files to `gh-pages`

## Contributing

We welcome contributions! If you have ideas for improvements or find issues with the chart, please open a pull request or submit an issue.

## License

This project is licensed under the [Apache License 2.0](LICENSE).
