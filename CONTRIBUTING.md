# Contributing to Containerum project

Containerum project welcomes all contributions. 
Before submitting a contribution, please read this document carefully.

Containerum consists of several components written in GO. To contribute, please go to the component repository and sumbit all changes there:

* [**api-gateway**](https://github.com/containerum/gateway) provides routing for Containerum components
* [**user-manager**](https://github.com/containerum/user-manager) is a service for managing users, groups, credentials, blacklists for Containerum
* [**resource**](https://github.com/containerum/resource) manages Kubernetes namespace objects: deployments, ingresses, etc.
* [**permissions**](https://github.com/containerum/permissions) manage user access to enable teamwork
* [**kube-api**](https://github.com/containerum/kube-api) is a set of API for communication between Containerum and K8s
* [**auth**](https://github.com/containerum/auth) handles user authorization and token management
* [**mail**](https://github.com/containerum/mail) is a mail server and newsletter template manager
* [**ui**](https://github.com/containerum/ui) is Web User Interface for Containerum
* [**chkit**](https://github.com/containerum/chkit) is CLI for Containerum

Note: [containerum/containerum](https://github.com/containerum/containerum) contains the **helm charts** for Containerum platform. The **source code** for each component is available in the project repositories above.

## Contribution guidelines

To submit a contribution, please follow these guidelines:

1. Go to the repository with the component you'd like to work on.

2. Create an issue to inform us about the problem you are trying to solve and fork the `master` branch.

3. Make changes and push them to your forked repository.

4. Once done, create a pull request to merge your fork with the `develop` branch of the original repo. Don't forget to describe the changes you are about to commit.

5. We will review the proposed changes and if everythings is alright, we will add them to the next release.

Thanks,
Containerum Team
