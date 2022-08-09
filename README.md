# Global Resource Service(GRS)


[![LICENSE](https://img.shields.io/badge/license-apache%202.0-green)](https://github.com/CentaurusInfra/global-resource-service/blob/master/LICENSE)


## What is Global Resource Service

Global Resource Service(GRS) is one of the corner stones for the [Regionless Cloud Platform](docs/Regionless-Platform.png). It provides abstraction, collection, aggregation of resources from each region and dynamic resource distribution to schedulers cross all regions globally.

GRS aims to be an open source solution to address key challenges of managing compute resources in the regionless cloud platform, such as:
1. Provide resource abstraction for various hardware types and machine types
2. Aggregate compute resources from all regions globally and distribute the desired set of resources to schedulers globally
3. Manage large scale of resource pool and distribute them with low latency
4. Balance resource allocation for ideal resource utilization globally

## Architecture
![Architecture Diagram](docs/GRS-Architecture.png)

## Key Features

* #### Large Scalability and high performance resource service and resource change propagation pipeline
* #### Provide resource collection and aggregation from multiple regions globally
* #### Heterogenious hardware type and machine type support
* #### Abstraction of compute resources
* #### Distribute/Balance/adjustment of resource allocation for ideal resource utilization

## Build and Run GRS

Please refer to the [dev-env-setup instruction](docs/setup-guide/setup-dev-env.md) for setting up local dev environment.

Please refer to the [setup-grs-test instruction](docs/setup-guide/setup-grs-test.md) for building, running end to end GRS tests.

## Contribute to the GRS project

Please refer to the [CONTRIBUTING guideline](./CONTRIUTING.md) for details and guidelines on how to contribute to the GRS project.

## Community Meetings 

 Pacific Time: **Tuesday, 6:00PM PST (Weekly)**

 Resources: [Meeting Link](https://futurewei.zoom.us/j/92636035970) | [Meeting Summary](https://docs.google.com/document/d/1Cwpp44pQhMZ_MQ4ebralDHCt0AZHqhSkj14kNAzA7lY/edit#)

## Documents and Support

The [design document folder](docs/design-proposals/) contains the detailed design of already implemented features, and also some thoughts for planned features.

To report a problem, please [create an issue](https://github.com/CentaurusInfra/global-resource-service/issues) in the project repo. 

To ask a question, here is [the invitation](https://join.slack.com/t/arktosworkspace/shared_invite/zt-cmak5gjq-rBxX4vX2TGMyNeU~jzAMLQ) to join [Arktos slack channels](http://arktosworkspace.slack.com/). You can also post in the [email group](https://groups.google.com/forum/#!forum/arktos-user), or [create an issue](https://github.com/CentaurusInfra/global-resource-service/issues) of question type in the repo.
