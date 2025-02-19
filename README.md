**SIMPLE GOLANG EWALLET**

## Required installation
- Go
- PostgreSQL
- redis
- NSQ

## Example .env
```
APP_VERSION=v1
DB_HOST=127.0.0.1
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=PASSWORD
DB_PORT=5432
DB_SSL=disable
DB_TIMEZONE=Asia/Jakarta
DB_AUTO_MIGRATE=false
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
HTTP_HOST=127.0.0.1
HTTP_PORT=10001
ACCESS_TOKEN_SECRET=MY_SECRET
ACCESS_TOKEN_EXPIRE_DURATION=48h
REFRESH_TOKEN_SECRET=MY_SECRET
REFRESH_TOKEN_EXPIRE_DURATION=120h
NSQ_PRODUCER_HOST=127.0.0.1:4150
NSQ_CONSUMER_HOST=127.0.0.1:4161
NSQ_N_CONSUMER=2
NSQ_MAX_IN_FLIGHT=20
NSQ_REQUEUE_TIME=10
NSQ_MAX_ATTEMPTS=20
TOPIC_FINISH_TRANSACTION=finish_transaction
CHANNEL_UPDATE_TRANSACTION_STATUS=update_transaction_status
```

## Example database diagram

![img.png](img.png)

## Example request & response

- API Register (Success)
![img.png](img.png)

- API Register (Failed)
![img_1.png](img_1.png)

- API Login (Success)
![img_2.png](img_2.png)

- API Login (Failed)
![img_3.png](img_3.png)

- API Update Profile (Success)
![img_5.png](img_5.png)

- API Update Profile (Failed)
![img_4.png](img_4.png)
