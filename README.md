# klock

## Description

Locking Kubernetes resources to prevent update or/and delete operation. Locking is done by defining `Lock` CRD. The `Lock`
define what operations should be watched and for what operation. Matching resources are done trough labels.

Example:

```yaml
apiVersion: klock.rnemet.dev/v1
kind: Lock
metadata:
  name: lock-sample
  namespace: test
spec:
  operations:
    - UPDATE
    - DELETE
  matcher:
    test: test
```

### Field spec

Field `spec` contains `Lock` specification. It has two sub resources:

* `operations` is a list containing operations to watch for
* `matcher` is a set of labels to match when looking for locked resources

### How it works

Based on `ValidatingWebhookConfiguration` when ever one of operations is invoked:

* check if in target `namespace` exists any `Lock`
* for every found `Lock` which contains target operation:
  * try to find matching `Lock` and target object

## Development

### Testing

Tests are located in folder `./tests/`:

```sh
make ktest
```

### Install tools

Tests are run by [kind](https://kind.sigs.k8s.io/) and [KUTTL](https://kuttl.dev/).

### Prepare local environment

Start `kind` cluster:

```sh
make make-cluster
```

Install `cert-manager`:

```sh
make cert-manager
```

### Build image locally

```sh
make make docker-build IMG=IMG=<some-registry>/klock:tag
```

Load image into local cluster: `kind load docker-image IMG=<some-registry>/klock:tag`

### Install

For development:

* Kustomize: `make deploy IMG=IMG=<some-registry>/klock:tag`
* Helm: `helm install klock klock` from `helm` folder.

To your cluster:

```sh
kubectl apply -f install/klock-0.0.1.yaml
```

### Run tests

```shell
make ktest
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

## License

Copyright 2022 Robert Nemet.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

