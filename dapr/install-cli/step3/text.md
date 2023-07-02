### Install Dapr Runtime

Check the dapr components are working as docker contrainers.

```
docker ps
```{{exec}}


The output will be like below.

```
CONTAINER ID   IMAGE                COMMAND                  CREATED          STATUS                    PORTS                                                 NAMES
1caa60d9aaf9   daprio/dapr:1.11.1   "./placement"            10 minutes ago   Up 10 minutes             0.0.0.0:50005->50005/tcp, :::50005->50005/tcp         dapr_placement
be9d23c6f1b5   openzipkin/zipkin    "start-zipkin"           10 minutes ago   Up 10 minutes (healthy)   9410/tcp, 0.0.0.0:9411->9411/tcp, :::9411->9411/tcp   dapr_zipkin
f7c6f16aa1ea   redis:6              "docker-entrypoint.s…"   10 minutes ago   Up 10 minutes             0.0.0.0:6379->6379/tcp, :::6379->6379/tcp             dapr_redis
```