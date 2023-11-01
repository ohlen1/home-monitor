# home-monitor
This service, built with Go, collects data about a home's electricity consumption from Tibber's [GraphQL Websocket API](https://developer.tibber.com/docs/guides/getting-started). It requires an active Tibber subscription. Based on the data from Tibber, Prometheus metrics are produced. 

**Disclaimer**: *This is work in progress. I don't commit to anything or give any guarantees whatsoever. You're welcome to use it, but do so at your own risk.*

## Deployment with Docker Compose
Running the [docker-compose.yaml](docker/docker-compose.yaml) will spin up the following components:
* `home-monitor` - The home-monitor service, based on a pre-built image on DockerHub
* `prometheus` - A Prometheus instance
* `grafana` - A Grafana instance, bootstrapped with datasource for Prometheus, for visualiziation of the Tibber metrics

1. Make a copy of [config_example.yaml](config_example.yaml)

2. [.env](docker/.env) sets var `DOCKERDIR="~/tmp/docker"`. Change this if you want to persist the application data for Grafana and Prometheus in some other directory.

3. Start services:
   ```bash
   docker compose -f docker/docker-compose.yaml up prometheus grafana home-monitor -d
   ```

Check logs for all services. 
```bash
docker compose -f docker/docker-compose.yaml logs -f --tail 2000 prometheus grafana home-monitor
```

Run a shell in the home-monitor container:
```bash
docker compose -f docker/docker-compose.yaml exec $1 sh
``````

See this article on how to setup docker-compose on a RaspberryPi:
https://dev.to/elalemanyo/how-to-install-docker-and-docker-compose-on-raspberry-pi-1mo