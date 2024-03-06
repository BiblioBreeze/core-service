# core-service

Implements core API for BiblioBreeze.

## Run

### Setup environment

```shell
docker-compose up -d

export DB_DSN=postgresql://postgres@127.0.0.1:5432/postgres
```

### Migrate

```shell
make migrate
```

### Run

```shell
make run
```