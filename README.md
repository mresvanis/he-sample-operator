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

Install the Hardware Enablement Sample Operator:

```shell
kubectl apply -k https://github.com/mresvanis/he-sample-operator/config/default
```
