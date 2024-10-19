# Go Authentication API с PostgreSQL и Keycloak

Это проект API с аутентификацией, использующий Go, PostgreSQL и Keycloak. Проект полностью контейнеризован с помощью Docker.

## Что используется:

- **Go**: Для создания API.
- **PostgreSQL**: Для хранения данных (например, стран).
- **Keycloak**: Для управления пользователями и аутентификацией.
- **Docker**: Для запуска всех сервисов через Docker Compose.

## Шаги для запуска проекта:

### 1. Требования:

- **Docker** и **Docker Compose** должны быть установлены на твоем компьютере.

### 2. Настройка:

1. Склонируй репозиторий проекта:

   ```bash
   git clone https://github.com/rtmelsov/auth-jwt.git
Перейди в директорию проекта:

```bash
cd auth-jwt
````
Создай файл .env в корне проекта и добавь туда переменные окружения для базы данных и Keycloak. Пример:

.env
```gitignore
KEY_CLOAK_CERT_URL=http://localhost:8081/realms/go-realm/protocol/openid-connect/certs
TOKEN_URL= http://localhost:8081/realms/go-realm/protocol/openid-connect/token
CLIENT_ID=go-client
CLIENT_SECRET=...

DATABASE_URL=postgres://myuser:mypassword@localhost:5432/myappdb?sslmode=disable
```

### 3. Запуск проекта:
   Собери и запусти контейнеры с помощью Docker Compose:

bash
Copy code
docker-compose up --build
Проверь, что все сервисы запущены. Ты можешь посмотреть работающие контейнеры с помощью:

bash
Copy code
docker-compose ps
### 4. Доступные API эндпоинты:
   GET /countries: Получить список всех стран.
   GET /countries/
   : Получить информацию о конкретной стране по id. 
   
### 5. Остановка проекта:
   Чтобы остановить и удалить все контейнеры:

```bash
docker-compose down
```
Примечание:
Убедись, что ты настроил Keycloak и добавил нужные client_id и client_secret в Keycloak для работы с JWT.
Все данные сохраняются в PostgreSQL, который также работает в Docker-контейнере.
