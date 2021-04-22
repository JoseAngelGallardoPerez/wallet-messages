**Dependencies**

 - Golang V1.10
 - Docker
 - MySql

**Build microservice**
Run the following command from microservice directory:

    GOOS=linux go build

**Containerize microservice**

Build an image:

    sudo docker build -t repository/name:version directory

 - **repository** - dockerhub or AWS ECS Registry repository name
 - **name**  - image name
 - **version** - image version
 - **directory** - build directory

Container can be started running:

    sudo docker run -e ENV_VAR1=value1 -e ENV_VAR2=value2 -i -t -p host_port:container_port repository/name:version

 - **-e** - flag to pass environment variable to the container
 - **ENV_VAR1=value1** and **ENV_VAR2=value2** - examples of environment variables
 - **repository** - dockerhub or AWS ECS Registry repository name
 - **name**  - image name
 - **version** - image version

Following environment variables are required for Velmie Wallet Message Service:

 - VELMIE_WALLET_MESSAGES_DB_PASS=passsord
 - VELMIE_WALLET_MESSAGES_DB_USER=user
 - VELMIE_WALLET_MESSAGES_DB_NAME=schema
 - VELMIE_WALLET_MESSAGES_DB_PORT=port
 - VELMIE_WALLET_MESSAGES_DB_HOST=host
 - VELMIE_WALLET_MESSAGES_DB_DRIV=mysql
 - VELMIE_WALLET_MESSAGES_DB_IS_DEBUG_MODE=false
 - VELMIE_WALLET_MESSAGES_CORS_HEADERS=*
 - VELMIE_WALLET_MESSAGES_CORS_ORIGINS=*
 - VELMIE_WALLET_MESSAGES_CORS_METHODS=GET,POST,PUT,DELETE,OPTIONS

**Migrate schema**

[[Golang migration tool]](https://github.com/golang-migrate/migrate)

    migrate -database mysql://username:password@protocol(host)/dbname -path database/migrations

**Push Image to repository**

    sudo docker push repository/name:version

 - **repository** - dockerhub or AWS ECS Registry repository name
 - **name**  - image name
 - **version** - image version

**Deploy containerized  microservice to AWS ECS**

To be completed ...

## Wallet Messages Helm chart configuration

For usage examples and tips see [this article](https://velmie.atlassian.net/wiki/spaces/WAL/pages/52004603/Wallet-+Helm+charts+getting+started).

The following table lists the configurable parameters of the wallet-messages chart, and their default values.

| Parameter                      | Description                                                                                                                      | Default                                 |
|--------------------------------|----------------------------------------------------------------------------------------------------------------------------------|:---------------------------------------:|
| service.type                   | The type of a service e.g. ClusterIp, NodePort, LoadBalancer                                                                     | ClusterIp                               |
| service.ports.public           | Application public API port.                                                                                                     | 10308                                   |
| service.ports.rpc              | Application RPC port.                                                                                                            | 12308                                   |
| service.ports.unsafeExposeRPC  | Forces to expose RPC port even if service.type other than ClusterIp                                                              | false                                   |
| service.selectors              | List of additional selectors                                                                                                     | {}                                      |
| containerPorts                 | List of ports that should be exposed on application container but in the service object.                                         | []                                      |
| containerLivenessProbe.enabled | Determines whether liveness probe should be performed on a pod.                                                                  |                                         |
| containerLivenessProbe.failureThreshold | Number of requests that should be failed in order to treat container unhealthy                                          | 5                                       |
| containerLivenessProbe.periodSeconds | Number of seconds between check requests.                                                                                  | 15                                      |
| appApiPathPrefix               | API prefix path. Used with internal health check functionality.                                                                  | messages                              |
| mysqlAdmin.user                | Privileged database user name. Used in order to create DB schema and user. Required if hooks.dbInit.enabled=true.                |                                         |
| mysqlAdmin.password            | Privileged database user password.                                                                                               |                                         |
| hooks.dbInit.enabled           | Enabled database init job.                                                                                                       | false                                   |
| hooks.dbInit.createSchema      | Determines whether to create database schema. Depends on hooks.dbInit.enabled                                                    | true                                    |
| hooks.dbInit.createUser        | Determines whether to create database user that will be restricted to only use specified database schema.                        | true                                    |
| hooks.dbMigration.enabled      | Determines whether to run database migrations.                                                                                   |                                         |
| ingress.enabled                | Determines whether to create ingress resource for the service.                                                                   | true                                    |
| ingress.annotations            | List of additional annotations for the ingress.                                                                                  | {"kubernetes.io/ingress.class": "nginx"}|
| ingress.tls.enabled            | Determines whether TLS (https) connection should be set.                                                                         | false                                   |
| ingress.tls.host               | Host name that is covered by a certificate. This value is required if ingress.tls.enabled=true.                                  |                                         |
| ingress.tls.secretName         | [Kubernetes secret](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls) name where TLS certificate is stored.  |                                         |
| appEnv.corsMethods             | Access-Control-Allow-Methods header that will be returned by the application.                                                    | GET,POST,PUT,OPTIONS                    |
| appEnv.corsOrigins             | Access-Control-Allow-Origin header that will be returned by the application.                                                     | *                                       |
| appEnv.corsHeaders             | Access-Control-Allow-Headers header that will be returned by the application.                                                    | *                                       |
| appEnv.dbHost                  | Database host to which application will be connected                                                                             | mysql                                   |
| appEnv.dbPort                  | Application database port.                                                                                                       | 3306                                    |
| appEnv.dbUser                  | Application database user.                                                                                                       |                                         |
| appEnv.dbName                  | Application database name.                                                                                                       |                                         |
| appEnv.dbDebugMode             | Whether database queries should be logged. Debugging mode.                                                                       | false                                   |
| image.repository               | What docker image to deploy.                                                                                                     | 360021420270.dkr.ecr.eu-central-1.amazonaws.com/velmie/wallet-currencies |
| image.pullPolicy               | What image pull policy to use.                                                                                                   | IfNotPresent                             |
| image.tag                      | What docker image tag to use.                                                                                                    | {Chart.yaml - appVersion}                |
| image.dbMigrationRepository    | What docker image to run in order to execute database migrations. By default the value if image.repository + "-db-migration"     | {image.tag}-db-migration                 |
| image.dbMigrationTag           | What docker image tag should be used for the db migration image.                                                                 | Same as image.tag                        |
| imagePullSecrets               | List of secrets which contain credentials to private docker repositories.                                                        | []                                       |
| nameOverride                   | Override this chart name.                                                                                                        | wallet-messages                        |
| fullnameOverride               | Override this chart full name. By default it is composed from release name and the chart name.                                   | {releaseName}-{chartName}                |
| serviceAccount.create          | Whether Kubernetes service account resource should be created.                                                                   | false                                    |
| serviceAccount.annotations     | Annotations to add to the service account                                                                                        | {}                                       |
| serviceAccount.name            | The name of the service account to use. If not set and create is true, a name is generated using the fullname template.          | See description                          |
| podAnnotations                 | Kubernetes pod annotations.                                                                                                      | {}                                       |
| securityContext                | A security context defines privilege and access control settings for a Pod or Container. [See details](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) | {} |
| resources                      | Limit Pod computing resources. [See details](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/)             | {}                                       |
| autoscaling.enabled            | Determines whether autoscaling functionality is enabled.                                                                         | false                                    |
| autoscaling.minReplicas        | [See details](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)                            | 1                                        |
| autoscaling.maxReplicas        | [See details](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)                            | 5                                        |
| autoscaling.targetCPUUtilizationPercentage | [See details](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)                | 80                                       |
| nodeSelector                   | [See details](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector)                             | {}                                       |
| tolerations                    | [See details](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/)                                     | []                                       |
| affinity                       | [See details](https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes-using-node-affinity/)                   | {}                                       |

## Run the project with Tilt

[Tilt](https://tilt.dev/) automates all the steps from a code change to a new process: watching files, building container images, and bringing your environment up-to-date.

[Install Tilt](https://docs.tilt.dev/install.html)

See [this article](https://velmie.atlassian.net/wiki/spaces/WAL/pages/56001240/Running+core+services+with+Tilt) which explains how to work with Tilt regarding this project.
