# dagagent

Airflow DAG Agent.

## Build

```shell
% docker build -t dagagent -f docker/Dockerfile .
```

## Configuration

Add the following configuration to [airflow/docker-compose.yaml](https://airflow.apache.org/docs/apache-airflow/stable/docker-compose.yaml).

```yaml
...
  dagagent:
    image: dagagent:latest
    ports:
      - 1323:1323
    environment:
      AIRFLOW__CORE__DAGS_FOLDER: /opt/airflow/dags
      _AIRFLOW_WWW_USER_USERNAME: airflow
      _AIRFLOW_WWW_USER_PASSWORD: airflow
    volumes_from:
      - airflow-webserver:rw
    healthcheck:
      test: ["CMD", "curl", "--fail", "http://airflow:airflow@localhost:1323/health"]
      interval: 10s
      timeout: 10s
      retries: 5
    restart: on-failure
...
```

## API Documentation

Create and start containers, and browser to http://localhost:1323/swagger/index.html.
