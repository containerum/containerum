[![Build Status](https://travis-ci.org/containerum/containerum.svg?branch=master)](https://travis-ci.org/containerum/containerum) [![HitCount](http://hits.dwyl.com/containerum/containerum.svg)](http://hits.dwyl.com/containerum/containerum) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Coverage](coverage_badge.png)]('https://github.com/jpoles1/gopherbadger)

# Embark to containerum!

![Containerum logo](../logo.svg)

A tool to quickly install [Containerum](https://containerum.com/software/) on a pre-configured Kubernetes.

## Getting Started

  * [Prerequisites](#prerequisites)
  * [Installing](#installing)
  * [Running the tests](#running-the-tests)
  * [Contributing](#contributing)
  * [Useful packages](#useful-packages)
  * [Versioning](#versioning)
  * [Authors](#authors)
  * [License](#license)
  * [Acknowledgments](#acknowledgments)

### Prerequisites

Embark requires installed and configured Kubernetes node. If you need to start a Kubernetes cluster, checkout these articles:
  + [4 ways to bootstrap a Kubernetes cluster](https://medium.com/containerum/4-ways-to-bootstrap-a-kubernetes-cluster-de0d5150a1e4)
  + [How to deploy Kubernetes and Containerum on Digital Ocean](https://medium.com/containerum/how-to-deploy-kubernetes-and-containerum-on-digital-ocean-eca93e6b4d26)
  + [Installing Kubernetes from binaries](https://medium.com/containerum/installing-kubernetes-from-binaries-pt-1-preparing-your-cluster-c229b2a8dca7)

### Installing
Just run 
```bash 
kubect create -f -
```
and then paste 

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: containerum/embark
  labels:
    app: containerum
    name: embark
spec:
  serviceAccountName: containerum-embark
  selector: {}
  template:
    metadata:
      labels:
        app: containerum
        name: embark
    spec:
      containers:
        - name: embark
          image: containerum/embark
          imagePullPolicy: Always
          volumeMounts:
            - name: kube
              mountPath: /etc/kube
      volumes:
        - name: kube
          configMap:
            defaultMode: 420
            name: kube-config
      restartPolicy: Never
```
This will launch Containerum installation job. Then wait until installation is complete:
```bash
kubectl wait --for=condition=complete containerum/embark
```

## Running the tests
Requirements:
  + go standart toolset

__Run the unit tests__
```bash
go test ./...
```

__Run the integration tests__

You need Kubernetes kluster up and running and kubectl config `~/.kube/config`.

**⛔️ !!DANGEROUS!!! ⛔️**

**DUE TO THE FACT THAT TESTS CAN CREATE AND DELETE OBJECTS IN THE KUBERNETES, YOU MUST BE SURE THAT THE CONFIGURATION GIVES ACCESS ONLY TO THE TEST CLUSTER, AND NO PRODUCTION SERVICES WILL BE DESTROYED OR CORRUPTED!**

```bash
go test -tags="IntegrationTests" ./...  
```

## Contributing

We are using [`dep`](https://github.com/golang/dep) ad dependency manager and [`fileb0x`](https://github.com/UnnoTed/fileb0x) as embedded filesystem for static assets (default configs, charts, kubernetes objects, etc.)

All code generation tools are vendored, so you don't need to install them, the only command required is `go generate -v`.

We welcome any help from the open source community. To submit your contributions, fork the project you want to contribute to (e.g. permissions, auth, etc.), commit changes and create a pull request to the develop branch. We will review the changes and include them to the project. Read more about contributing in this [document](../CONTRIBUTING.md).

## Useful packages 
You may want to use some of the components of this utility in your projects.
First of all you should pay attention to:
  + [`cgraph`](pkg/cgraph) oriented task graph builder (used to resolve dependency conflicts)
  + [`object`](pkg/object) generic Kube object model, useful if you want to manipulate structured data and then create custom Kubernetes objects
  + [`spin`](pkg/utils/spin) KISS terminal spinner package
  + [`why`](pkg/utils/why) KISS package for structured printing to STDOUT and io.Writer
    ```go
      why.Print("root item", "first item", "second item", "third item")
      // root item
      //     ╠═ first item
      //     ╠═ second item
      //     ╚═ third item
    ```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/containerum/containerum/tags).

## Authors

* **Paul Petrukhin** - *Initial work* - [ninedraft](https://github.com/ninedraft)

See also the list of [contributors](https://github.com/containerum/containerum/contributors) who participated in this project.

## License

This project is licensed under the Apache License - see the [LICENSE](../LICENSE) file for details

## Acknowledgments

* Thanks to Kubernetes team for providing well documented code as such components as `apimachinery` and `client-go`
* Thanks to Helm gang for excelent quality of their product
