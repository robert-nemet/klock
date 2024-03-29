# Klock

## Motive

Clean and easy way to protect the k8s resource from deletion or updates. It could be achieved with RBAC but
it is not explicit enough. This does not exclude RBAC, an operator which set the lock on the resource still needs
rights to access it.

## Goals

* Mandatory locking (no one can update/delete a resource) as long the lock exists
* Exclusive locking (only one can update/delete a resource) as long the lock exists

## Design

Lock CRD which resign in a namespace, protects any resource which matches at least one label. Protection is only for
specified operations.

Avoid deadlock by ignoring lock on a lock.

Which resources are tracked by the lock is defined in `ValidatingWebhookConfiguration`.

### Mandatory locking

```yaml
apiVersion: klock.rnemet.dev/v1
kind: Lock
metadata:
  name: lock-red-blue
  namespace: yellow
spec:
  operations:
    - UPDATE
    - DELETE
  matcher:
    aura: red
```

In above example `Lock` is mandatory lock. In the namespace `yellow` is will prevent any UPDATE/DELETE operation
on every resource that have at least one label as defined in `matcher` section. So, if there is a Pod with label
`aura: red` in the same namespace as the lock it can be deleted as long the lock exists.

#### Mandatory locking with expressions

If you need to lock multiple resources with same label name but different values you can use expressions like:

```yaml
aura: red|blue|green
```

Which means protect from update and delete all resources if they have label named `aura` which value is `red` or
`blue` or `green`.

In case label named `aura` can have values: red, blue, green and black, above expression can be:

```yaml
aura: ^black
```

Which can be read as protect all resources with label named `aura` which value is __not__ `black`. 

Other examples:

```yaml
...
aura: ^(red|blue)
...
aura: black|^blue
```

Supported operations are &(AND) , |(OR) and ^(NOT).

See test: /tests/e2e/lock-expression

#### Mandatory locking matching multiple labels

Locks like this:

```yaml
apiVersion: klock.rnemet.dev/v1
kind: Lock
metadata:
  name: lock-protect-sun
spec:
  operations:
    - UPDATE
    - DELETE
  matcher:
    aura: red
    element: ^wind
```

Before version v3.0.0 would lock any resource that have label `aura: red`. In version v3.0.0 for resource to be locked it has to match all expressions under matcher.
Looking in above example a Pod with labels `aura: red` and `element: earth` will be locked, while a Pod with labels `aura: red` and `element: wind` will not be lock.

See test: /tests/e2e/lock-expressions-two-labels

### Exclusive locking

```yaml
apiVersion: klock.rnemet.dev/v1
kind: Lock
metadata:
  name: lockred
spec:
  operations:
    - UPDATE
    - DELETE
  matcher:
    aura: red
  exclusive:
    name: johny
```

Same as in previous example except if operation requester is the `johny`. The `johny` can delete a pod.

### Field spec

Field `spec` contains `Lock` specification. It has two sub resources:

* `operations` is a list containing operations to watch for.
* `matcher` is a set of labels to match when looking for locked resources. For operation to be denied ony one needs to be matched.
* `exclusive` define requester params when setting exclusive lock. One parameter is enough. If both are set then both
  needs to match: 
  * `name` requester name
  * `uid` requester uid

### How it works

Based on `ValidatingWebhookConfiguration` when ever one of operations is invoked:

* check if in target `namespace` exists any `Lock`
* for every found `Lock` which contains target operation:
  * try to find matching `Lock` and target object

A `lock` __will not__ lock other `lock`, even if it set by configuration.

## Configuration

Cpnfiguration is autogenerated by [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder). Initially only `Pods`, `Deployments` and,
`Secretcs` can be locked. This can be changed by modifying `ValidatingWebhookConfiguration` named `klock-validating-webhook-configuration`.
As well, `UPDATE` and `DELETE` operations are set there as operations to watch for. If needed user can modify that as well. It is not recommended
to watch for all possible resources.

### Default webhook configuration

```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  annotations:
    cert-manager.io/inject-ca-from: klock-system/klock-serving-cert
  name: klock-validating-webhook-configuration
webhooks:
- admissionReviewVersions:
  - v1
  clientConfig:
    service:
      name: klock-webhook-service
      namespace: klock-system
      path: /validate-all
  failurePolicy: Fail
  name: klocks.rnemet.dev
  rules:
  - apiGroups:
    - '*'
    apiVersions:
    - '*'
    operations:
    - DELETE
    - UPDATE
    resources:
    - pods
    - deployments
    - secrets
  sideEffects: None
```

If you need some other resources or operations modify webhook configuration. Operations used in webhook configuration are used in `Lock`.

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
make docker-build IMG=IMG=<some-registry>/klock:tag
```

Load image into local cluster: `kind load docker-image IMG=<some-registry>/klock:tag`

### Install

#### For development

* Kustomize: `make deploy IMG=IMG=<some-registry>/klock:tag`
  
To your cluster:

```sh
kubectl apply -f install/klock-1.0.0.yaml
```

#### Helm

```sh
$ helm repo add rnemet https://rnemet.dev/helm-charts
$ helm install klock rnemet/klock
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
