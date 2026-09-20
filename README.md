# Subglutee Project - Notification Service

worker consume message จาก RabbitMQ แล้วส่งอีเมลผ่าน SMTP, log ผลลง Mongo

### Message contract

producer (เช่น sgt-scheduler) publish JSON ไป queue `email_notifications`:

```json
{
    "event_id": "billing:22222222-2222-4222-8222-222222222222:2026-09-23:3d",
    "user_id": "11111111-1111-4111-8111-111111111111",
    "subscription_id": "22222222-2222-4222-8222-222222222222",
    "to": "test@example.com",
    "subject": "Upcoming charge: Netflix",
    "body": "Your Netflix subscription renews in 3 days."
}
```

IDs และวันที่ข้างต้นมีไว้สำหรับทดสอบเท่านั้น user/subscription ไม่จำเป็นสำหรับการทดสอบ email worker

| Field             | JSON type | Meaning                                                                         |
| ----------------- | --------- | ------------------------------------------------------------------------------- |
| `event_id`        | string    | ID ของ reminder หนึ่งครั้ง ใช้ ID เดิมเมื่อ retry หรือ publish reminder เดิมซ้ำ |
| `user_id`         | string    | ID เจ้าของ notification สำหรับ in-app notification ในขั้นต่อไป                  |
| `subscription_id` | string    | ID subscription ที่เกี่ยวข้อง                                                   |
| `to`              | string    | Email ผู้รับ                                                                    |
| `subject`         | string    | หัวข้อ email ที่ producer จัดเตรียมแล้ว                                         |
| `body`            | string    | เนื้อหา plain-text email ที่ producer จัดเตรียมแล้ว                             |

service นี้ไม่รู้จัก schema ของ subscription/user เลย - format เนื้อหาเมลเป็นหน้าที่ producer

### Ownership and compatibility

- Producer เป็นผู้ตัดสินว่าเตือนเมื่อใด และ format เนื้อหา email
- Notification Service รับ IDs เป็น metadata; ไม่จำเป็นต้องรู้ schema ทั้งหมดหรือ query database ของ user/subscription
- `event_id` ต้องเหมือนเดิมสำหรับ reminder เดิม แม้ scheduler scan ใหม่หรือ retry; reminder ของรอบบิลใหม่ต้องมี ID ใหม่
- `event_id` เป็นเพียงตัวอย่าง ไม่ได้บังคับในโค้ดและต้องตกลงกับ producer ก่อนใช้งาน
- ช่วงเปลี่ยนผ่าน payload เก่ายัง decode ได้ โดย fields ใหม่จะเป็น empty string
- DTO มีเพียง JSON tags ไม่ได้ตรวจ required fields, UUID หรือ email โดยอัตโนมัติ; validation เป็นงานถัดไป
- Consumer ปัจจุบันยังส่งเพียง `To`, `Subject`, `Body` ไป usecase ดังนั้น fields ใหม่ยังไม่ถูกบันทึกหรือใช้ deduplicate

### Structure

```
consumer/       RabbitMQ consumer, decode message -> เรียก usecase
usecases/       ผูก mailer + repository, ตัดสิน status
mailer/         SMTP adapter (usecase กำหนด interface เอง)
repositories/   Mongo adapter เก็บ log (usecase กำหนด interface เอง)
models/         NotificationLog struct
dtos/           EmailMessage: queue message struct + proposed metadata
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
lefthook install
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

### Run alternatively (container)

```terminal
make image
make container
```

### Useful Commands

Check `Makefile`

### Dev tools

- MailHog UI: `localhost:8025` - ดูเมลที่ส่งออก
- RabbitMQ management UI: `localhost:15672` (guest/guest) - ดูคิว
- mongo-express: `localhost:8089` - ดู log

### Email flow test

1. Restart Go service หลังเปลี่ยน DTO: `Ctrl+C` แล้ว `go run .`
2. เปิด `http://localhost:8088/health` และตรวจ `status: ok`
3. เปิด `RabbitMQ Management` → Queues and Streams → `email_notifications` → Publish message
4. Publish extended JSON ด้านบน
5. เปิด `MailHog` ตรวจ subject/body และผู้รับ
6. เปิด `Mongo Express` → database `noti` → collection `email_logs` ตรวจ record ที่มี `status: sent`

การทดสอบนี้ยืนยันว่า email flow ยังทำงานเมื่อมี fields ใหม่เท่านั้น ไม่ได้ยืนยันว่า metadata ถูกบันทึกหรือป้องกัน duplicate แล้ว.
ถ้าต้องการตรวจ mapping ของ DTO โดยตรง ให้ตรวจ `json.Unmarshal` ด้วย debugger หรือ unit test.
การ publish event เดิมสองครั้งใน implementation ปัจจุบันยังอาจส่ง email สองครั้ง.
