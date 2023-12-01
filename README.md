
[![forthebadge](https://forthebadge.com/images/featured/featured-built-with-love.svg)](https://forthebadge.com)
# TODO-list api
todo list api

Используемые технологии: 
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Fiber (веб фреймворк)
- golang-migrate/migrate (для миграций БД)
- JWT-token (для авторизации)

# Usage

**Скопируйте проект**
```bash
  git clone https://github.com/artemiyKew/todo-list-rest-api.git
```

**Перейдите в каталог проекта**
```bash
  cd todo-list-rest-api
```

**Запустите сервер**
```bash
  make compose
```

## Examples
- [Регистрация](#регистрация)
- [Аутентификация](#аутентификация) 
- [Получение информации о текущем пользователе](#получение-информации-о-текущем-пользователе)
- [Создание работы](#cоздание-работы)
- [Получение списка работ](#получение-списка-работ)
- [Обновление работы](#обновление-работы)
- [Удаление работы](#удаление-работы)

## Регистрация
Регистрация пользователя: 

```bash
    curl -X POST http://localhost:1234/sign-up
        -H "Content-Type: application/json" \
        -d '{
            "email": "test@example.com", 
            "password": "password"
        }'
```
Пример ответа: 
```json
{
    "id":1,
    "email":"test@example.com"
}
```

## Аутентификация
Аутентификация пользователя:
```bash
    curl -X POST http://localhost:1234/sign-in
        -H "Content-Type: application/json" \
        -d '{
            "email": "test@example.com", 
            "password": "password"
        }'
```

Пример ответа: 
```json
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE1NDYzMDIsInN1YiI6MX0.YYByTZhM3HUk0oS-HAbmejxyXX6KqqmGUGH9dr33U3w
```
## Получение информации о текущем пользователе
Получение информации о текущем пользователе:

```bash
    curl -X http://localhost:1234/auth/whoami \
        -H "Token: <your_jwt_token> "
```
Пример ответа: 
```json
{
    "id":1,
    "email":"test@example.com"
}
```

## Создание работы
Создание работы пользователя(достпуно только авторизированным пользователям):
```bash
    curl -X POST  http://localhost:1234/auth/work \
        -H "Content-Type: application/json" \
        -H "Token: <your_jwt_token> \
        -d '{
            "name": "Task 1", 
            "description": "Description for Task 1"
        }' 
```

Пример ответа: 
```json
{
    "id":1,
    "name":"Task 1",
    "user_id":1,
    "description":"Description for Task 1",  
    "created_at":"2023-12-01T19:49:44.496621387Z",
    "exp_at":"2023-12-02T19:49:44.496621387Z"
}
```

## Получение списка работ
Получение списка работ пользователя:
```bash
    curl -X http://localhost:1234/auth/works \
        -H "Token:<your_jwt_token>"
```
Пример ответа: 
```json
{
    [
        {
            "id":1,
            "name":"Task 1",
            "user_id":1,
            "description":"Description for Task 1",
            "created_at":"2023-12-01T19:49:44.496621Z",
            "exp_at":"2023-12-02T19:49:44.496621Z"
        },
        {
            "id":2,
            "name":"Task 1",
            "user_id":1,
            "description":"Description for Task 1",
            "created_at":"2023-12-01T19:53:48.47464Z",
            "exp_at":"2023-12-02T19:53:48.47464Z"
        }
    ]
}
```

## Обновление работы
Обновление работы пользователя:
```bash
    curl -X PUT http://localhost:1234/auth/work/{work_id} \
        -H "Content-Type: application/json" \
        -H "Token:<your_jwt_token>" \
        -d '{
            "name": "Updated Task", 
            "description": "Updated Description"
        }'
```
Пример ответа: 
```json
{
    "id":0,
    "name":"Updated Task",
    "user_id":1,
    "description":"Updated Description",
    "created_at":"0001-01-01T00:00:00Z",
    "exp_at":"0001-01-01T00:00:00Z"
}
```

## Удаление работы
Удаление работы пользователя:
```bash
    curl -X DELETE  http://localhost:1234/auth/work/{work_id} \
        -H "Token:<your_jwt_token>"
```



