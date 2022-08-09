TODO: Convert the entire design doc from word doc to md file

# Global Resource Service in Regionless Cloud Platform (Arktos 2.0)

## Background and Context 

 

Demanding of reducing of control plane cost from the cloud providers are getting higher in recent years. To simply scale-up of each cluster to a larger number of compute nodes can help these demands however it will be limited by its control plane scalability and capacities. 

[Arktos](https://github.com/CentaurusInfra/arktos) helps improve the cluster scalability to a new level. It separates the resource management and the application management into different clusters, and it proves that the scalability, in terms of number of nodes in one cluster, number of nodes the scheduler can use for its schedulers.

Facebook’s global infra management with RAS, Twine and Shard Manager [2][3][4] goes one step further to dynamically allocate and adjust the tenant resource allocation to best utilize the compute resources in a region. It improves resource utilization efficiencies. However, Facebook’s infra is mainly designed for its internal applications with dedicated, unique hardware in its data centers. The solution cannot be utilized by generic cloud providers for various customer applications and various hardware types in cloud provider’s data centers. 

Cloud providers offer services around globe with multiple regions to support the customers in different GEO locations. They also ensure data locality policy and compliance policies for their customers. Each region, due to their GEO location, the availability of nature resources and cost from the power companies can be different. For example, the power capacity in Beijing region in China is relatively limited and the cost is higher than it in the remote regions such as Ulan or Guizhou regions. However, it is common for those remote regions to have lower number of application deployment due to various reasons.

Cloud providers often offer incentives for applications deployed in the remote regions and try to lead/guide/allure the customers to deploy their applications in those remote regions. With the increased resource utilization, it reduces provider’s operation cost to maintain those remote regions. However, for some customers, even they were provided with the information how to choose the desired regions for their application deployments, they will have to add the logic to their deployment flows, which increases their deployment complexity [1].

The demand of utilizing and balancing the compute resources in different regions becomes more and more from cloud providers as well as the local governments for better utilizing the nature resources cross regions. For example, China has the computation infrastructure goals in its 13th 5-year-plan to guide cloud providers to leverage the nature resources and less expensive power supply in the western regions in China. The goal is to move the heavily loaded computation from the most populated, high power cost eastern regions to the western regions in China (东数西算)。 

This new trend in policy and in cost-saving of cloud computing raise another level of requirements in addition to improve the control plane efficiencies. A cloud provider must provide a cross-region deployment scheduling logic that places the customer applications to the desired regions based on the application’s desire to satisfy cost, locality, compliance, and performance of the applications. 

Sky computing, was recently raised again by UC Berkeley RISE lab (now SKY lab) with a goal to address those requirements. It solves the common issues such as “provider lock-in”, “not able utilize resources from specific cloud providers even they are with lower pricing model”.  However, this approach is much restricted due to the limited access to each cloud provider’s internal resource management APIs. Currently, its vision is mainly focused on big data related workflow scenarios with a resource broker model to leverage resources from different cloud providers.  

A more promising approach is to provide an infrastructure inside cloud providers to abstract out the region concept and provides a unified, yet flexible mechanism so that customers can deploy their applications/services to their desired regions, without even aware of the region concepts. 

We envision this as a [Regionless Cloud Platform](../Regionless-Platform.png). It can consist of the layers described below conceptually: 


* A new application development model, workload abstraction and APIs
* A global tenant, application management layer that hides the region or GEO location from customers. It provides a set of policies/hints that customers can leverage for their application/services’ placements.
* A global scheduling layer that is capable of handling large scale customers and their deployments
* A global resource management layer that abstracts out the regional compute resources and the scheduler layer will leverage this layer for the deployment placements.
* An infrastructure in each region that provides physical machine managements, network infrastructure etc.


We envision that the network technology and hardware will continue to improve in the coming years and the network latency among regions will continue to improve as well. The aforementioned- global scope services are feasible with the much-improved network latencies, especially with the dedicated networks for the region control planes. 

The regionless cloud platform is not designed to address the “provider lock-in”.  However, as described in the future works, it has the potential to address the “provider lock-in” once cloud providers agree on set of public APIs to render their compute resources in some of their regions.

This document describes the design details of the resource management layer. It first introduces the design requirements, goals, and use cases; then it describes the design of the global resource management service and the key workflows in detail, followed by the evaluation criteria and results. Lastly, it shares some thoughts of future works. 
