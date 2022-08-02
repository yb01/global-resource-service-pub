# Release Summary

This is the first release of the Global Resource Service, one of the corner stones for the Regionless Cloud Platform.

### The 0.1.0 release includes the following components:

* Global Resource Service API server, that supports REST APIs for client registration, List Assigned nodes and Watch for node changes.
* Performant distributor, event queues and cache to support large scale of node changes
* Data Aggregator that collects nodes and node changes from each region
* Client development SDK that provides APIs for building scheduler or other clients to the Global Resource Service
* A Region manager simulator that provides region level, multiple Resource Provider simulation of data changes
* A simulated scheduler with the cache layer
* A test infrastructure to automate service deployment, cross region test setup, test execution and result collection



### Key Features:

* Client registration, Node List and Watch APIs
* Distributor algorithms to support multiple region and resource clusters
* Distributor algorithms for efficient, balanced node resource distribution to schedulers
* Scalability: Scale up to 1m nodes cross multiple regions, with up to 40 schedulers
* Performance: End to end latency just 300ms for normal node failures cases (Daily change pattern) and within 1.3 seconds for disaster scenarios(RP outage pattern)
* Abstraction of node resources, aka, logical min. record for node resource
* Abstraction of resource version, aka, Composite RV ( or CRV ) from nodes from different and global origins. 
* Cross region data change simulation of both "Daily" and "RP outage" test scenarios
* Automatic test environment setup, test execution and result collection routines.


P.S.  [Performance test results](../resource-management/docs/performance/730-2022-perf-results.md)
