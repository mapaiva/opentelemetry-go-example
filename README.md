# opentelemetry-go-example

Service with metrics and traces examples using Prometheus and Elastic APM.

## Components

- Elasticsearch
- Elastic APM Server
- Kibana
- [opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib)

## Running

```sh
> docker-compose up -d
> go run cmd/server/main.go
> curl localhost:2021/users/mapaiva
```

## See tracing

Open Kibana at http://localhost:5601/app/apm and configure the APM Agent.
