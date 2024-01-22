# Overview

ArgoCD Tag Plugin is an innovative tool designed to dynamically fetch the latest or a specific tag of a Docker image. This tool is particularly useful for continuous deployment workflows using ArgoCD, allowing for automated updates of image tags in Kubernetes deployments.

# Features

- **Dynamic Tag Fetching**: Automatically fetches the latest or specified tag of an image from a Docker registry.
- **Tag Constraints and Selectors**: Supports specifying constraints (like 'latest') and selectors (like prefixes or suffixes) for image tags.
- **Integration with ArgoCD**: Seamlessly integrates with ArgoCD, enabling automatic image updates during the deployment process.

# How it Works

1. **Configuration in values.yaml:**
   Users specify the image tag in their values.yaml file using a special token format. For example:

```yaml
kind: Deployment
---
image: <account_id>.dkr.ecr.<region>.amazonaws.com/the/path:<image-tag-plugin:latest#prefix(commit-)>
```

2. **Commit to Git:**
   This file is then committed to a Git repository.
3. **Configure ArgoCD to Use the Plugin:**
   Users must configure [ArgoCD to use this plugin](https://argo-cd.readthedocs.io/en/stable/operator-manual/config-management-plugins/).
4. **Automatic Tag Detection and Fetching:**
   Whenever a new 'latest' image is pushed to the Docker registry, the ArgoCD refresh process detects it and fetches the relevant commit hash for deployment.

# Prerequisites

- A working ArgoCD setup.
- Docker images tagged with both 'latest' and specific commit hashes in your Docker registry.

# Token Structure

The token used in the **values.yaml** file has a specific structure:

- `<image-tag-plugin:` This is the identifier for the plugin to recognize its tokens.
- `latest` This represents a tag constraint. It specifies which tag to look for when searching available tags.
- `prefix(commit-)` This is a tag selector. After finding a list of tags using the tag constraint, it selects which tag to use for the image in the target yaml file. Options include prefix(...), suffix(...), or contains(...).

# Installation and Configuration

## Install the Plugin as a Sidecar Plugin

Begin by installing the ArgoCD Tag Plugin as a sidecar plugin. Follow the instructions detailed in the [Argo CD documentation](https://argo-cd.readthedocs.io/en/stable/operator-manual/config-management-plugins/#place-the-plugin-configuration-file-in-the-sidecar) to correctly place the plugin configuration file in the sidecar.

## Create an Argo Plugin Configuration

Next, create an Argo plugin configuration as per the guidelines found in the [Argo CD documentation](https://argo-cd.readthedocs.io/en/stable/operator-manual/config-management-plugins/#write-the-plugin-configuration-file). This step involves writing the plugin configuration file necessary for the ArgoCD Tag Plugin to function correctly within your ArgoCD setup.

## Configure the generate Phase

During the generate phase of your ArgoCD workflow, you can run the ArgoCD Tag Plugin using the following command:

```sh
helm template $ARGOCD_APP_NAME -n $ARGOCD_APP_NAMESPACE --include-crds -f values.yaml . |
        argocd-tag-plugin generate -
```

This command integrates the plugin with Helm, allowing it to process the templates and apply the dynamic tag fetching feature of the ArgoCD Tag Plugin.

## Final Steps

After completing these installation and configuration steps, the ArgoCD Tag Plugin should be fully integrated into your ArgoCD environment, ready to dynamically manage your Docker image tags according to your specified constraints and selectors.

Always ensure that your configurations are in line with your deployment requirements and thoroughly test the setup in a controlled environment before applying it to your production workflows.

# License

Licensed under the MIT License.
