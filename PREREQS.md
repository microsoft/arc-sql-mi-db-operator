# Prerequisites

A lot of moving parts to an operator with webhooks.

## CertManager

The first thing needed is to install the CertManager or another tool that can provide the required ssl certificate communication.  Here are the articles to explain what is required.

[Cert Manger](https://book.kubebuilder.io/cronjob-tutorial/cert-manager.html)

You can follow [the cert manager documentation](https://docs.cert-manager.io/en/latest/getting-started/install/kubernetes.html) to install it.

## Docker Repository

You will need to have a docker registry that you can push the operator and the sync image to.  You will also need to verify that your `K8S` cluster will have the credentials to pull images from that registry.

## External `Kubectl` Image

To allow seamless `K8S` api access, the Database Operator uses a sidecar image `bitnami/kubectl:latest`.  This container is resposible to proxy the api to the main operator container.  This allows the operator to connect to `K8S` api through localhost.

```bash
kubectl proxy --port=9090 &
```

For more info about why this is needed and other potential options, you can look at the [documentation for accessing K8S api from a pod](https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/).

### Why

The sidecar method is used because we do not have access to the `SQL Managed Instance` code to use the native `go-client` with typed support.