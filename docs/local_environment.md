# Local Environment Guide

These docs provide step by step instructions for how to get the Quibbble controller up and running on a local [K8s](https://kubernetes.io) cluster running on your machine.

## Pre-Requisites

Please install the following tools:
- [Docker](https;//docker.com)
- [Kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl)

## Create Local Cluster

There are a number of local cluster solutions out there; [Minikube](https://minikube.sigs.k8s.io), [Kind](https://kind.sigs.k8s.io), and  [K3d](https://k3d.io) to name a few. Quibbble currently runs on [K3s](https://k3s.io) in production and given that [K3d](https://k3d.io) is just a wrapper on top of [K3s](https://k3s.io) we recommend using [K3d](https://k3d.io) for local testing.

### Install K3d

Please follow the installation instructions on [k3d.io](https://k3d.io) to install it onto your specific OS. 

## Create Cluster

This will start a single node cluster on your local machine and expose port `80` to allow you to interact with the cluster over http. Please make sure this `quibbble-controller` repo is your root directory before running the `make` command.

```bash
make create_cluster
```

## Install NGINX

```bash
helm upgrade --install ingress-nginx ingress-nginx \
    --repo https://kubernetes.github.io/ingress-nginx \
    --namespace ingress-nginx --create-namespace
```

## Create Admin Password

Create a k8s secret that holds your admin password.

```bash
kubectl create secret generic quibbble-controller \
    --namespace=quibbble \
    --from-literal=admin-password=<YOUR_PASSWORD_HERE>
```

## Install Quibbble Controller

This will install the Quibbble Controller onto your cluster. You should not need to make any changes but take [values.yaml](../charts/quibbble-controller/values.yaml) file beforehand and change any values as desired. 

```bash
helm upgrade --install quibbble-controller charts/quibbble-controller \
    --values charts/quibbble-controller/values.yaml \
    --namespace quibbble --create-namespace
```

## Check Setup

You should now be able to curl and hit the Quibbble Controller. 

```bash
curl http://127.0.0.1/health
```
