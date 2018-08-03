# Overview

This project is used to generate the data for
https://medium.com/dm03514-tech-blog/sre-debugging-simple-memory-leaks-in-go-e0a9e6d63d4d


## Starting Prometheus/Grafana/Node Exporter

- install docker-compose

```
$ pip install docker-compose
```

- Start prometheus/grafana/node exporter

```
$ make start-prom

cd stack \
&& docker-compose down \
&& docker-compose up
Removing network stack_default
WARNING: The Docker Engine you're using is running in swarm mode.

Compose does not use swarm mode to deploy services to multiple nodes in a swarm. All containers will be scheduled on the current node.

To deploy your application across the swarm, use `docker stack deploy`.

Creating network "stack_default" with the default driver
Creating stack_exporter_1 ... done
Creating stack_prom_1     ... done
Creating stack_grafana_1  ... done
Attaching to stack_exporter_1, stack_prom_1, stack_grafana_1
exporter_1  | time="2018-08-03T22:37:38Z" level=info msg="Starting node_exporter (version=0.16.0, branch=HEAD, revision=d42bd70f4363dced6b77d8fc311ea57b63387e4f)" source="node_exporter.go:82"
exporter_1  | time="2018-08-03T22:37:38Z" level=info msg="Build context (go=go1.9.6, user=root@a67a9bc13a69, date=20180515-15:52:42)" source="node_exporter.go:83"
exporter_1  | time="2018-08-03T22:37:38Z" level=info msg="Enabled collectors:" source="node_exporter.go:90"
exporter_1  | time="2018-08-03T22:37:38Z" level=info msg=" - arp" source="node_exporter.go:97"
exporter_1  | time="2018-08-03T22:37:38Z" level=info msg=" - bcache" source="node_exporter.go:97"
...
grafana_1   | t=2018-08-03T22:37:40+0000 lvl=warn msg="[Deprecated] The folder property is deprecated. Please use path instead." logger=provisioning.dashboard type=file name=default
grafana_1   | t=2018-08-03T22:37:40+0000 lvl=info msg="Initializing Stream Manager"
grafana_1   | t=2018-08-03T22:37:40+0000 lvl=info msg="HTTP Server Listen" logger=http.server address=0.0.0.0:3000 protocol=http subUrl= socket=
```

- Node Exporter data should be available locally

http://localhost:3000/d/yvumWBFmk/mem-leak-service?refresh=5s&orgId=1

## Starting the Mem Leak Service

```
$ go run main.go
Starting Server on port :8080
```

This service will leak memory on every http requests.  Prometheus is registered
to scrape it, and its data will show up on the grafana dashboard linked to above.

## Apply a steady load to the service in order to generate leak data

```
$ echo "GET http://localhost:8080/work" | vegeta attack -rate 2000 > /dev/null
```
