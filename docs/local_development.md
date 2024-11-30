# Local Development Guide

There are a few steps you'll need to follow in order to test any local code changes end-to-end.

## Pre-Requisites

Please follow the [Local Environment Guide](local_environment.md) before following the steps below.

## Build New Image

Once you've made your changes, you'll need to build new Docker images to pickup your code changes:

```bash
make docker_build
```

## Import Images

Your new local images will also need to be imported into k3d for them to be used:

```bash
make docker_import
```

## Upgrade Controller

Update the [values.yaml](../charts/quibbble-controller/values.yaml) file, changing all `pullPolicy` values from `Always` to `Never`. This ensures your local images are used instead of the official ones hosted on Docker Hub.

Re-run the Helm command to pickup your new images.

```bash
helm upgrade --install quibbble-controller charts/quibbble-controller \
    --values charts/quibbble-controller/values.yaml \
    --namespace quibbble --create-namespace
```
