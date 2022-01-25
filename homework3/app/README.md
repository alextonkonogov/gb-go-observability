# gb-go-observability

# Homework3
Идем в директорию homework3/app

```console
cd homework3/app
```

Запускаем сервисы в контейнерах через docker compose

```bat
docker compose up
```

Затем отдельно собираем приложение (почему см. ниже)
```bat
go build main.go
```
И запускаем его
```bat
./main
```

Проверяем
```console
curl -s  http://localhost:8080/
```

ВОПРОС!
Если приложение будет в докер контейнере, то jaeger его просто не видит. Убил на эту проблему кучу времени, но так и не нашел решение.
Пример в методичке почему-то тоже обходит стороной эту тему и приложение собирается и запускается отдельно.
Решение отсюда не помогает: https://stackoverflow.com/questions/50173643/tracing-with-jaeger-doesnt-work-with-docker-compose 
