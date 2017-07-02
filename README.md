# Minio Broker

This is an implementation of a Service Broker that uses Helm to provision
instances of [Minio](https://kubeapps.com/charts/stable/minio). This is a
**proof-of-concept** for the [Kubernetes Service
Catalog](https://github.com/kubernetes-incubator/service-catalog), and should not
be used in production. Thanks to the [mariadb broker repo](https://github.com/prydonius/mariadb-broker).

## Prerequisites

1. Kubernetes cluster
2. [Helm 2.x](https://github.com/kubernetes/helm)
3. [Service Catalog API](https://github.com/kubernetes-incubator/service-catalog) - follow the [walkthrough](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/walkthrough.md)

## Installing the Broker

The Minio Service Broker can be installed using the Helm chart in this
repository.

```
$ git clone https://github.com/yqf3139/minio-broker.git
$ cd minio-broker
$ helm install --name minio-broker --namespace minio-broker charts/minio-broker
```

To register the Broker with the Service Catalog, create the Broker object:

```
$ kubectl --context service-catalog create -f examples/minio-broker.yaml
```

If the Broker was successfully registered, the `minio` ServiceClass will now
be available in the catalog:

```
$ kubectl --context service-catalog get serviceclasses
NAME      KIND
minio   ServiceClass.v1alpha1.servicecatalog.k8s.io
```

## Usage

### Create the Instance object

```
$ kubectl --context service-catalog create -f examples/minio-instance.yaml
```

This will result in the installation of a new Minio chart:

```
$ helm list
NAME                                  	REVISION	UPDATED                 	STATUS  	CHART               	NAMESPACE
i-3e0e9973-a072-49ba-8308-19568e7f4669	1       	Sat May 13 17:28:35 2017	DEPLOYED	minio-0.6.1       	3e0e9973-a072-49ba-8308-19568e7f4669
```

### Create a Binding to fetch credentials

```
$ kubectl --context service-catalog create -f examples/minio-binding.yaml
```

A secret called `minio-instance-credentials` will be created containing the
connection details for this Minio instance.

```
$ kubectl get secret minio-instance-credentials -o yaml
```
