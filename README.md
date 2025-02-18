**SIMPLE GOLANG EWALLET**

## 1.Example .env
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
```

## 2.Example database diagram

![img.png](img.png)

## 3.Example database data

 - User Account
![img_1.png](img_1.png)
 - Transaction
![img_2.png](img_2.png)
 - Provider
![img_4.png](img_4.png)
 - Provider Setting
![img_3.png](img_3.png)
 - Balance Movement
![img_5.png](img_5.png)

## 4.API Doc
```
- POST http://127.0.0.1:10001/v1/auth/login-pin
    BODY
        {
        "phone": "62111111",
        "pin": "123456"
        }
        
- POST http://127.0.0.1:10001/v1/auth/verify-pin
    HEADER
        Authorization Bearer {{TOKEN}}
    BODY
        {
        "pin": "123456",
        "type": "WITHDRAW"
        }
        
- POST http://127.0.0.1:10001/v1/transaction/topup
    HEADER
        Authorization Bearer {{TOKEN}}
    BODY
        {
        "va_number": "111111",
        "amount": 50000
        }
        
- POST http://127.0.0.1:10001/v1/transaction/transfer
    HEADER
        Authorization Bearer {{TOKEN}}
    BODY
        {
        "to_account_id": "300ee889-a156-4e46-bcbe-ba7c8c3b12d8",
        "amount": 9999999999
        }
        
- POST http://127.0.0.1:10001/v1/transaction/withdraw
    HEADER
        Authorization Bearer {{TOKEN}}
    BODY
        {
        "provider_id": "bd76b150-7e6e-467e-b671-82f11a6273c5",
        "amount": 50000
        }
        
- GET http://127.0.0.1:10001/v1/user-account/get-info
    HEADER
        Authorization Bearer {{TOKEN}}
        
- GET http://127.0.0.1:10001/v1/user-account/get-balance-history
    HEADER
        Authorization Bearer {{TOKEN}}
```