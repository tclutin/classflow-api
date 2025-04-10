# classflow-api
## 📖Описание
ClassFlow API — это сервис для управления группами, расписанием занятий и пользователями. Администраторы могут создавать группы и загружать/удалять в них расписание. Пользователи же могут вступать и выходить из групп, смотреть расписание на разные недели и дни, а также настраивать уведомления о предстоящих парах.

Данный сервис был разработан в рамках университетского задания и имеет дополнительный компоненты: [сервис оповещений](https://github.com/tclutin/classflow-notification-service), [телеграм бот](https://github.com/Clonsaldafon/CsuScheduleBot)

## ⚙️Установка
1️⃣ Изменить название файла
```bash
.env.example -> .env
```
2️⃣ Поменять конфигурацию в зависимосит от своих потребностей
```bash
ENVIRONMENT=dev #can be prod

HTTP_HOST=app
HTTP_PORT=8080

POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DB=schedule_backend
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres

JWT_SECRET=asdokghi090qw902109
JWT_EXPIRE=1h
```
3️⃣ Запустить сервис
```bash
docker-compose up --build
```

## 📚 Документация
Если ENVIRONMENT=dev, то документация и спецификация будут доступны [тут](http://localhost:8080/swagger/index.html)
