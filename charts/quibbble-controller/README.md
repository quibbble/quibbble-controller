# Quibbble Controller

This chart install the Quibbble controller onto a K8s cluster. This controller currently requires a load balancer such as NGINX to sit in front of it to allow for dynamic routing of games. Meaning that if a game with key `tictactoe` and id `example` is created, then NGINX will handle the routing to that game over path `/tictactoe/example`.

## Setup

#### Install NGINX

```
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```

#### Install Quibbble Controller

```
helm install quibbble-controller charts/quibbble-controller/ \
    --values charts/quibbble-controller/values.yaml \
    --namespace quibbble --create-namespace
```
