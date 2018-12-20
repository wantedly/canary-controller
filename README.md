# Canary Controller

Canary Controller is a Custom Resource and Controller for managing Canary server.

## Concept

The Canary Release is an approach to release a new application more safely. In short, the Canary Release consists of the following steps:

1. Deploy new application. (called "Canary server")
1. Send some of trafics to the Canary server. Typicatlly we'd join a Canary server for the Load Balancer and Load Balancer routes traffic using the DNS round robin. In Kubernetes, `Service` assumes this role
1. Evaluate a Canary server by comparing with other server deployed old application (called Baseline server) using some metrics, for example RPS, resource usage and error rate.
1. Roll out the new application to all servers OR rollback to the old version.

Canary Controller supports the second phase, managing a Canary server and routing. Canary Controller creates a `Deployment` of a Canary server based on the existing ones. Additionally, Canary Controller put some metadata to a generated Deployment for identifying whether an application is a Canary.

* Attach `canary: true` to labels
  * Kubernetes manages resources using labels like selector of Service
* Give a name the host of Pod to `canary`
* Add `CANARY_ENABLED: 1` to environment variables of all containers

## Installation

This installation is dependent by [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and [kustomize](https://github.com/kubernetes-sigs/kustomize/blob/master/docs/INSTALL.md). If your environment is OSX then you can install these using the following commands:

```
$ brew install kubernetes-cli
$ brew install kustomize
```

You can deploy to your cluster using `make` command. 

```
$ make deploy
```

Canary Controller is installed in `canary-controller-manager` namespace.

## Usage

Apply a very simple manifest if you want to deploy a Canary server: 

```yaml
apiVersion: canary.k8s.wantedly.com/v1beta1
kind: Canary
metadata:
  name: foo-canary
  namespace: default
spec:
  targetDeploymentName: "foo"
```

Canary Controller looks up existing Deployment named "foo" in the "default" namespace and creates a Deployment based on "foo":

```
$ kubectl get deploy -n default
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
foo       1         1         1            1           1m

$ kubectl apply -f config/sample/canary_v1beta1_canary.yaml
canary.canary.k8s.wantedly.com "canary-sample" created

$ kubectl get deploy -n default
NAME         DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
foo          1         1         1            1           1m
foo-canary   1         1         1            1           9s
```

And the Deployment created by Canary Controller will cleanup when deletes the manifest:

```
$ kubectl delete -f config/sample/canary_v1beta1_canary.yaml
canary.canary.k8s.wantedly.com "canary-sample" deleted

$ kubectl get deploy
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
foo       1         1         1            1           1m
```
