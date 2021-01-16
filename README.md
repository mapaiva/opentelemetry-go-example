# Opentel

Tests with open telemetry

## Running

```sh
> docker-compose up -d
> go run cmd/server/main.go
```

## See tracing metrics

Open Kibana at http://localhost:5601/app/apm and configure APM Agent.
