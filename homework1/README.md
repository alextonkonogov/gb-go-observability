# gb-go-observability

# Homework1
Идем в директорию homework1

```console
cd homework1
```

Запускаем постгрес в контейнере

```bat
docker run \
    --rm -it \
    -p 5432:5432 \
    --name postgres \
    -e POSTGRES_PASSWORD=password \
    -e PGDATA=/var/lib/postgresql/data \
    postgres:13.1
```

Заливаем туда данные

```console
migrate -database "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable" -path migrations up
```

Поднимаем контейнеры с prometheus и grafana
```console
docker run -d -p 9090:9090 prom/prometheus:v2.25.0
```

```console
docker run -d -p 3000:3000 --name grafana grafana/grafana
```
Запускаем приложение

```console
go run cmd/motivation-keeper/main.go
```

Проверяем работу
```console
curl -s  http://localhost:1234
```
