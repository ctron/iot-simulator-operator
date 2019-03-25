# Building a Cluster Service Version (CSV) for the Operator Framework

This guide is intended to guide an Operator author to package a version of their Operator to run with the [Operator Lifecycle Manager](https://github.com/operator-framework/operator-lifecycle-manager). This will be a manual method that will walk through each section of the file, what it’s used for and how to populate it.

## What is a Cluster Service Version (CSV)?

A CSV is the metadata that accompanies your Operator container image. It can be used to populate user interfaces with info like your logo/description/version and it is also a source of technical information needed to run the Operator, like the RBAC rules it requires and which Custom Resources it manages or depends on.

The Lifecycle Manager will parse this and do all of the hard work to wire up the correct Roles and Role Bindings, ensure that the Operator is started (or updated) within the desired namespace and check for various other requirements, all without the end users having to do anything.

You can read about the [full architecture in more detail](architecture.md#what-is-a-clusterserviceversion).

## CSV Metadata

The object has the normal Kubernetes metadata. Since the CSV pertains to the specific version, the naming scheme is the name of the Operator + the semantic version number, eg `mongodboperator.v0.3`.

The namespace is used when a CSV will remain private to a specific namespace. Only users of that namespace will be able to view or instantiate the Operator. If you plan on distributing your Operator to many namespaces or clusters, you may want to explore bundling it into a [Catalog](architecture.md#catalog-registry-design).

The namespace listed in the CSV within a catalog is actually a placeholder, so it is common to simply list `placeholder`. Otherwise, loading a CSV directly into a namespace requires that namespace, of course.

```yaml
apiVersion: operators.coreos.com/v1alpha1
kind: ClusterServiceVersion
metadata:
  name: mongodboperator.v0.3
  namespace: placeholder
```

## Your Custom Resource Definitions
There are two types of CRDs that your Operator may use, ones that are “owned” by it and ones that it depends on, which are “required”.
### Owned CRDs

The CRDs owned by your Operator are the most important part of your CSV. This establishes the link between your Operator and the required RBAC rules, dependency management and other under-the-hood Kubernetes concepts.

It’s common for your Operator to use multiple CRDs to link together concepts, such as top-level database configuration in one object and a representation of replica sets in another. List out each one in the CSV file.

**DisplayName**: A human readable version of your CRD name, eg. “MongoDB Standalone”

**Description**: A short description of how this CRD is used by the Operator or a description of the functionality provided by the CRD.

**Group**: The API group that this CRD belongs to, eg. database.example.com

**Kind**: The machine readable name of your CRD

**Name**: The full name of your CRD

The next two sections require more explanation. 

**Resources**:
Your CRDs will own one or more types of Kubernetes objects. These are listed in the resources section to inform your end-users of the objects they might need to troubleshoot or how to connect to the application, such as the Service or Ingress rule that exposes a database.

It’s recommended to only list out the objects that are important to a human, not an exhaustive list of everything you orchestrate. For example, ConfigMaps that store internal state that shouldn’t be modified by a user shouldn’t appear here.

**SpecDescriptors**:

These are a way to hint UIs with certain inputs or outputs of your Operator that are most important to an end users. If your CRD contains the name of a Secret or ConfigMap that the user must provide, you can specify that here. These items will be linked to and highlighted in compatible UIs. The description allows for you to explain how this is used by the Operator.

todo(?): list out the options

Below is an example of a MongoDB “standalone” CRD that requires some user input in the form of a Secret and ConfigMap, and orchestrates Services, StatefulSets, Pods and ConfigMaps.

```yaml
      - displayName: MongoDB Standalone
        group: mongodb.com
        kind: MongoDbStandalone
        name: mongodbstandalones.mongodb.com
        resources:
          - kind: Service
            name: ''
            version: v1
          - kind: StatefulSet
            name: ''
            version: v1beta2
          - kind: Pod
            name: ''
            version: v1
          - kind: ConfigMap
            name: ''
            version: v1
        specDescriptors:
          - description: Credentials for Ops Manager or Cloud Manager.
            displayName: Credentials
            path: credentials
            x-descriptors:
              - 'urn:alm:descriptor:com.tectonic.ui:selector:core:v1:Secret'
          - description: Project this deployment belongs to.
            displayName: Project
            path: project
            x-descriptors:
              - 'urn:alm:descriptor:com.tectonic.ui:selector:core:v1:ConfigMap'
          - description: MongoDB version to be installed.
            displayName: Version
            path: version
            x-descriptors:
              - 'urn:alm:descriptor:com.tectonic.ui:label'
        version: v1
        description: >-
          MongoDB Deployment consisting of only one host. No replication of
          data.
```

### Required CRDs

Relying on other “required” CRDs is completely optional and only exists to reduce the scope of individual Operators and provide a way to compose multiple Operators together to solve an end-to-end use case. An example of this is an Operator that might set up an application and install an etcd cluster (from an etcd Operator) to use for distributed locking and a Postgres database (from a Postgres Operator) for data storage.

The Lifecycle Manager will check against the available CRDs and Operators in the cluster to fulfill these requirements. If suitable versions are found, the Operators will be started within the desired namespace and a Service Account created for each Operator to create/watch/modify the Kubernetes resources required.

**Name**: The full name of the CRD you require

**Version**: The version of that object API

**Kind**: The Kubernetes object kind

**DisplayName**: A human readable version of the CRD

**Description**: A summary of how the component fits in your larger architecture

```yaml
    required:
    - name: etcdclusters.etcd.database.coreos.com
      version: v1beta2
      kind: EtcdCluster
      displayName: etcd Cluster
      description: Represents a cluster of etcd nodes.
```
## CRD Templates
Users of your Operator will need to be aware of which options are required vs optional. You can provide templates for each of our CRDs with a minimum set of configuration as an annotation named `alm-examples`. Compatible UIs will pre-enter this template for users to further customize.

The annotation consists of a list of the `kind`, eg. the CRD name, and the corresponding `metadata` and `spec` of the Kubernetes object. Here’s a full example that provides templates for `EtcdCluster`, `EtcdBackup` and `EtcdRestore`:

```yaml
metadata:
  annotations:
    alm-examples: >-
      [{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdCluster","metadata":{"name":"example","namespace":"default"},"spec":{"size":3,"version":"3.2.13"}},{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdRestore","metadata":{"name":"example-etcd-cluster"},"spec":{"etcdCluster":{"name":"example-etcd-cluster"},"backupStorageType":"S3","s3":{"path":"<full-s3-path>","awsSecret":"<aws-secret>"}}},{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdBackup","metadata":{"name":"example-etcd-cluster-backup"},"spec":{"etcdEndpoints":["<etcd-cluster-endpoints>"],"storageType":"S3","s3":{"path":"<full-s3-path>","awsSecret":"<aws-secret>"}}}]

```
## Operator Metadata
The metadata section contains general metadata around the name, version and other info that aids users in discovery of your Operator.

**DisplayName**: Human readable name that describes your Operator and the CRDs that it implements

**Keywords**: A list of categories that your Operator falls into. Used for filtering within compatible UIs.

**Provider**: The name of the publishing entity behind the Operator

**Maturity**: The stability of the Operator, eg. alpha, beta or stable

**Version**: The semanic version of the Operator. This value should be incremented each time a new Operator image is published.

**Icon**: a base64 encoded image of the Operator logo or the logo of the publisher. The `base64data` parameter contains the data and the `mediatype` specifies the type of image, eg. `image/png` or `image/svg`.

**Links**: A list of relevant links for the Operator. Common links include documentation, how-to guides, blog posts, and the company homepage.

**Maintainers**: A list of names and email addresses of the maintainers of the Operator code. This can be a list of individuals or a shared email alias, eg. support@example.com.

**Description**: A markdown blob that describes the Operator. Important information to include: features, limitations and common use-cases for the Operator. If your Operator manages different types of installs, eg. standalone vs clustered, it is useful to give an overview of how each differs from each other, or which ones are supported for production use.

## Operator Install
The install block is how the Lifecycle Manager will instantiate the Operator on the cluster. There are two subsections within install: one to describe the `deployment` that will be started within the desired namespace and one that describes the Role `permissions` required to successfully run the Operator.

Ensure that the `serviceAccountName` used in the `deployment` spec matches one of the Roles described under `permissions`.

Multiple Roles should be described to reduce the scope of any actions needed containers that the Operator may run on the cluster. For example, if you have a component that generates a TLS Secret upon start up, a Role that allows `create` but not `list` on Secrets is more secure than using a single all-powerful Service Account.

Here’s a full example:

```yaml
  install:
    spec:
      deployments:
        - name: example-operator
          spec:
            replicas: 1
            selector:
              matchLabels:
                k8s-app: example-operator
            template:
              metadata:
                labels:
                  k8s-app: example-operator
              spec:
                containers:
                    image: 'quay.io/example/example-operator:v0.0.1'
                    imagePullPolicy: Always
                    name: example-operator
                    resources:
                      limits:
                        cpu: 200m
                        memory: 100Mi
                      requests:
                        cpu: 100m
                        memory: 50Mi
                imagePullSecrets:
                  - name: ''
                nodeSelector:
                  beta.kubernetes.io/os: linux
                serviceAccountName: example-operator
      permissions:
        - serviceAccountName: example-operator
          rules:
            - apiGroups:
                - ''
              resources:
                - configmaps
                - secrets
                - services
              verbs:
                - get
                - list
                - create
                - update
                - delete
            - apiGroups:
                - apps
              resources:
                - statefulsets
              verbs:
                - '*'
            - apiGroups:
                - apiextensions.k8s.io
              resources:
                - customresourcedefinitions
              verbs:
                - get
                - list
                - watch
                - create
                - delete
            - apiGroups:
                - mongodb.com
              resources:
                - '*'
              verbs:
                - '*'
        - serviceAccountName: example-operator-list
          rules:
            - apiGroups:
                - ''
              resources:
                - services
              verbs:
                - get
                - list
    strategy: deployment
```

## Full Examples

Several [complete examples of CSV files](https://github.com/operator-framework/operator-lifecycle-manager/tree/master/deploy/chart/catalog_resources/rh-operators) are stored in Github.
