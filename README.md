# Subglutee Project - Notification Service

worker consume message จาก RabbitMQ แล้วส่งอีเมลผ่าน SMTP, log ผลลง Mongo

### Message contract

producer (เช่น sgt-scheduler) publish JSON ไปคิว `email_notifications`:

```json
{ "to": "user@example.com", "subject": "...", "body": "..." }
```

service นี้ไม่รู้จัก schema ของ subscription/user เลย - format เนื้อหาเมลเป็นหน้าที่ producer

### Structure

```
consumer/       RabbitMQ consumer, decode message -> เรียก usecase
usecases/       ผูก mailer + repository, ตัดสิน status
mailer/         SMTP adapter (usecase กำหนด interface เอง)
repositories/   Mongo adapter เก็บ log (usecase กำหนด interface เอง)
models/         NotificationLog struct
dtos/           queue message struct
routes/         GET /health
config/         env loader + mongo/rabbitmq connection
```

flow: `consumer -> usecases -> mailer (ส่งเมล) + repositories (log ผล)`

### Prerequisite

- Go 1.26
- Docker + Docker Compose
- `make` - `winget install ezwinports.make`
- tools:

```terminal
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install github.com/evilmartians/lefthook@latest
```

### Setup

```terminal
git clone https://github.com/polar-bear-cu/sgt-noti-service.git
cd sgt-noti-service
lefthook install
cp .env.example .env
go mod download
make compose-up
make run
```

### Useful Commands

Check `Makefile`

### Dev tools

- MailHog UI: `localhost:8025` - ดูเมลที่ส่งออก
- RabbitMQ management UI: `localhost:15672` (guest/guest) - ดูคิว
- mongo-express: `localhost:8089` - ดู log
