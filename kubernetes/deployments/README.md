What is included in this blog:
- What is a Kubernetes Deployment
- How to create a Deployment
- The Relationship between ReplicaSets and Deployments
- How to update a Deployment
- Use Case of Deployments
 
## prerequisites
I recommend you have basic knowledge of Kubernetes Pods before reading this blog. You can check [this doc](https://kubernetes.io/docs/concepts/workloads/pods/pod/ "Pods in Kubernetes") for details about Kubernetes Pods.

## What Is A Deployment

A Deployment is a Kubernetes object designed to manage stateless applications. It scales up a set of Pods to the desired number you describe in a config file (a JSON or a yaml file). These Pods are replicas of each other: They run the same containers and have exactly the same workflow. Kubernetes also provides a simple way to manage a Deployment Pods, including scaling up/down Pods and doing rolling updates.


## How to Create A Deployment

The following is an example of a Deployment configuration: It scales up three Pods and each Pod runs a Nginx container.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
 name: nginx-deployment-demo
spec:
 selector:
   matchLabels:
     app: nginx
     env: demo
 replicas: 3
 strategy:
   rollingUpdate:
     maxUnavailable: 0
 template:
   metadata:
     labels:
       app: nginx
       env: demo
   spec:
     affinity:
       nodeAffinity:
         requiredDuringSchedulingIgnoredDuringExecution:
           nodeSelectorTerms:
           - matchExpressions:
             - key: failure-domain.beta.kubernetes.io/zone
               operator: In
               values:
               - us-central1-a
               - us-central1-b
               - us-central1-c
       podAntiAffinity:
         preferredDuringSchedulingIgnoredDuringExecution:
         - weight: 100
           podAffinityTerm:
             labelSelector:
               matchExpressions:
               - key: app
                 operator: In
                 values:
                 - nginx
             topologyKey: kubernetes.io/hostname
     containers:
     - name: nginx
       image: nginx:1.15.3
       ports:
       - containerPort: 80
```

### Required Fields

The fields `apiVersion`, `kind` and `metadata` are required for a Deployment as they define the metadata. `metadata.name` is used to name the deployment and `metadata.labels` is used to add some labels to the deployment.

The `spec` field is also required as it describes the specification of the deployment and `spec.template` defines the pods managed by the Deployment. **`spec.template` is the boundary between the definitoin of the Deployment and the definition of the Deployment Pods.**


### Pod Selector

The `spec.selector` is used for the Deployment to find which pods to manage. In this example, the Deployment uses `app: nginx && env: demo` defined in `sepc.selector.matchLabels` to find the pods that have labels `{app: nginx, env: demo}` (defined in `spec.template.metadata.labels`). `sepc.selector.matchLabels` is a map of key-value pairs and requirements are ANDed.

Instead of using `sepc.selector.matchLabels`, you can use `sepc.selector.matchExpressions` to define more sophisticated match roles. You can check [this doc](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#labelselector-v1-meta "LabelSelector in Kubernetes") for more details about the usage of `sepc.selector.matchExpressions`.

**As you can see, a Deployment relies on pod labels and pod selector to find its pods. Therefore, it is recommended to put some unique pod labels for a Deployment. Otherwise, Deployment A may end up with managing the pods belonging to Deployment B.**

### Replica

The `spec.replica` field specifies the desired number of Pods for the Deployment. **It is Highly recommended to run at least two replicas for any Deployment in Production.** This is because having at least two replicas at the beginning can help you keep your Deployments stateless, as the problem can be easily detected when you are trying to introduce "stateful stuff" to a Deployment with at least two replicas. For example, you will quickly realize the problem when you are trying to add a cron job to a two-replicas Deployment to process some data in a daily base: the data will be processed twice a day as all replicas have to execute this cron job, which may cause some unexpected behaviour. In addition, a singleton pod may cause some downtime in some cases. For example, a single-replica deployment will not be available for a moment when the single Pod is triggered to restart for whatever reason.

### Rolling Update Strategies

The `spec.strategy` field defines the strategy for replacing old pods with new ones when a rolling update occurs. The `spec.strategy.type` can be `Recreate` or `RollingUpdate`. The default value is `RollingUpdate`.

In general, it is not recommended use `Recreate` in Production based on the consideration of availability. This is because `Recreate` will introduce downtime when a rolling update occurs: All the existing pods need to be terminated before new ones are created when `spec.strategy.type` is `Recreate`.

You can use `maxUnavailable` and `maxSurge` to control the update process when you set `spec.strategy.type == RollingUpdate`. The `maxUnavailable` field sets the maximum number of Pods that can be unavailable during an update process, while the `maxSurge` specifies the maximum number of Pods that can be created over the desired number of Pods. The default value is 25% for these two fields. Moreover, they cannot be set 0 at the same time, as this stops the Deployment from performing the rolling update. **One suggestion is to set `maxUnavailable` 0 as this is the most effective way to prevent your old pods from being terminated while there are some problems spinning up new pods.**


### Affinity for Pods
The `affinity` field inside the `spec.template.spec` allows you to specify which nodes should run in the Deployment's Pods. **As shown in the following picture, the ideal scenario of running a Deployment is running multiple replicas in different nodes in different zones, and avoid running multiple replicas in the same node.**

![The Ideal Scenario of Running A Deployment](https://raw.githubusercontent.com/azhuox/blogs/master/kubernetes/deployments/k8s-ideal-scenario-of-running-deployment.png "The Ideal Scenario of Running A Deployment")

You can utilize `spec.template.spec.affinity` to achieve this goal. Kubernetes provides `nodeAffinity` for you to constrain which nodes to run your Pods based on node labels. It also provides `podAffinity` and `podAntiAffinity` for you to specify inter-pod affinity. The official explanation of `podAffinity` and `podAntiAffinity` is "Inter-pod affinity and anti-affinity allow you to constrain which nodes your pod is eligible to be scheduled based on labels on pods that are already running on the node rather than based on labels on nodes." You can check [this doc](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/ "Assigning Pods to Nodes") for more details about node/pod affinity.


## ReplicaSets vs Deployments

The following picture demonstrates the relationship between ReplicaSets and Deployments

![The Relationship Between ReplicaSets and Deployments](https://raw.githubusercontent.com/azhuox/blogs/master/kubernetes/deployments/k8s-deploys-vs-replicasets.png "The Relationship Between ReplicaSets and Deployments")

ReplicaSet is the next generation Replication Controller designed to replace the old
[ReplicationController](https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/ "Kubernetes ReplicationController"). **ReplicaSet ensures that a specific number of pod replicas are running at a given time, based on `replica number` you define in a ReplicaSet spec.** Although it provides an easy way to replicates Pods, it lacks the ability to rolling update pods.

Deployments are built on the top of ReplicaSets. A Deployment essentially is a set of ReplicaSets. Kubernetes rolls out a new ReplicaSet for a Deployment when a rolling updte occurs. It runs the desired number of pods for the new ReplicaSet and smoothly terminates pods in the old ReplicaSet. **In other words, a Deployment performs the rolling update by replacing current ReplicaSet with a new one.** You can check [this doc](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#updating-a-deployment "Updating a Deployment") for more details about rolling update or rolling back Deployments.


## Use Case

### Recommended Use Cases

**Deployments are suitable for running "stateless" applications. This because Pods inside a Deployment do not have sticky identification and they execute exactly the same workflow.** 

The following picture shows the typical use case of Deployments. You can see that a three-replicas Deployment is used to run the user micro-service. The replicas share the storage and serve the same APIs.

![A Use Case of Deployments](https://raw.githubusercontent.com/azhuox/blogs/master/kubernetes/deployments/k8s-deploy-user-uservice.png "A Use Case of Deployments")

### Something to avoid

**One important principle that you should follow is: DO NOT add anything stateful to any replica of any Deployment. You should keep this in mind and review this principle whenever you make a change to a Deployment.** Some typical mistakes that I have made are: 1. Stick shared data to each Pod inside a Deployment; 2. Run cron jobs on a Deployment that has only a single Pod.

The following picture demonstrates the first case. From the picture you can see that Nginx default caching system uses local disk to store cache, and each Nginx replica maintains its own cache. This causes two problems: 1. A Nginx replica will lose its cache whenever it restarts; 2. the whole cashing system is low efficiency as a page request needs to be served in all the replicas in order to get itself "fully" cached. The root cause was that I stored page cache to each Nginx Pod, while the cache is supposed to be stored in a place where it can be shared among all the Nginx Pods. 

![A "Statful" System in Deployments](https://raw.githubusercontent.com/azhuox/blogs/master/kubernetes/deployments/k8s-deploy-ngx-cache-system.png "A "Statful" System in Deployments")

Be careful when you need to run cron jobs in a Deployment for several reasons. Firstly, running a time-consuming cron job in a Pod is not safe as it can be aborted and cannot be resumed from any Pod disaster. Secondly, as all the pods within a Deployment execute the same code, having a cron job in a deployment means all the Pods have to run this cron job at a given frequency, which may duplicate the data that you do not want. It is ok to do that if the cron job is tightly related to each Pod, like a daily cron to clean up each Pod's /tmp directory. Otherwise, running a cron job inside a Deployment may become a headache one day, when you need to increase the replica to two (This is why I recommed running at least two replicas for any Deployment).

The following picture shows the wrong usage of cron jobs in Deployments. A user micro-service, a Deployment with a single pod, intends to utilize a cron job to sends the daily digest to all the subscribed users. It loops through all the users and sends today's daily digest to each user. It may work well when the system does not have too many users. For example, it may have two hundred thousand users at the beginning and it may just take three hours to finish the cron job. However things will start to break when the user number grows to two million and it has to spend thirty hours running the cron job, provided everything goes well. The worst part is that increasing the number of replicas won't help because they have no way to partition the work amongst themselves.  **One way to avoid the problems like this, when using a Deployment, is to increase the replica to at least two at the beginning to force you to make the Deployment "stateless"**

![Wrong Usage of Cron Jobs in Deployments](https://raw.githubusercontent.com/azhuox/blogs/master/kubernetes/deployments/k8s-deploy-cron-jobs.png "Wrong Usage of Cron Jobs in Deployments")

That's it. Thank you for reading this blog.

## Reference
- [Pods in Kubernetes](https://kubernetes.io/docs/concepts/workloads/pods/pod/ "Pods in Kubernetes")
- [LabelSelector in Kubernetes](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#labelselector-v1-meta "LabelSelector in Kubernetes")
- [Assigning Pods to Nodes](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/ "Assigning Pods to Nodes")
- [Kubernetes ReplicationController](https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/ "Kubernetes ReplicationController")
- [Kubernetes Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#updating-a-deployment "Updating a Deployment")