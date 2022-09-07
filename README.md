![lint](https://github.com/mresvanis/he-sample-operator/actions/workflows/lint.yaml/badge.svg)
![tests](https://github.com/mresvanis/he-sample-operator/actions/workflows/test.yaml/badge.svg)
[![codecov](https://codecov.io/gh/mresvanis/he-sample-operator/branch/main/graph/badge.svg?token=EMH9QLP6NR)](https://codecov.io/gh/mresvanis/he-sample-operator)
[![go report](https://goreportcard.com/badge/github.com/mresvanis/he-sample-operator)](https://goreportcard.com/report/github.com/mresvanis/he-sample-operator)

# Hardware Enablement Sample Operator

The Hardware Enablement Sample Operator is an example of a Kubernetes
[Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/), which uses the
[Kernel Module Management Operator](https://github.com/kubernetes-sigs/kernel-module-management) to
manage [out-of-tree](https://www.kernel.org/doc/Documentation/kbuild/modules.txt) kernel modules in
[Kubernetes](https://kubernetes.io/).

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get
a local cluster for testing, or run against a remote cluster. **Note:** Your controller will
automatically use the current context in your kubeconfig file (i.e. whatever cluster
`kubectl cluster-info` shows).

### Running on the cluster

1. Install Instances of Custom Resources:

```shell
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```shell
make docker-build docker-push IMG=<some-registry>/he-sample-operator:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```shell
make deploy IMG=<some-registry>/he-sample-operator:tag
```

### Uninstall CRDs

To delete the CRDs from the cluster:

```shell
make uninstall
```

### Undeploy controller

UnDeploy the controller to the cluster:

```shell
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources untile the desired state
is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```shell
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```shell
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```shell
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

