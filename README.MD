# classflow-api

## Установка

1.Склонируйте репозиторий
```bash
git clone https://github.com/tclutin/classflow-api
```
2.Перейдите в директорию 
```bash
cd classflow-api
```
3.Измените название файла .env.example на .env, затем настройте его по мере необходимости

4.Запустите контейнеры
```bash
docker-compose up
```

5.Также примените миграции для контейнера с приложением
```bash
docker ps
```
```bash
docker exec <name of container with app> goose -dir migrations postgres "postgresql://postgres:postgres@localhost:5432/schedule_backend?sslmode=disable" up 
```

