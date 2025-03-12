# Ozon URL Shortener

Этот сервис позволяет сокращать длинные URL и получать их оригинальные версии по коротким ссылкам

## API Эндпоинты

### 1. Создание короткой ссылки
`POST /url`
Принимает оригинальный URL через параметр `url` в теле запроса и возвращает сокращённый URL.

**Пример запроса:**
```json
{
  "url": "https://example.com"
}
```

### 2. Получение оригинального URL
`GET /url?alias=<short_link>` Принимает сокращённый URL через параметр `?alias=` и возвращает оригинальный.

## Запуск проекта

Требуется **Docker**, **Task**, **golang-migrate** и **mockery**.

### 1️⃣ Установите `Task`
```sh
go install github.com/go-task/task/v3/cmd/task@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### 2️⃣ Установите `golang-migrate`
```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### 3️⃣ Установите `mockery`
```sh
go install github.com/vektra/mockery/v2@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### 4️⃣ Скопируйте и настройте `.env` файл
```sh
cp deploy/.env.sample deploy/.env
```
Измените пароли в `deploy/.env`

### 5️⃣ Выбор хранилища
Тип хранилища задаётся в `configs/local.yaml` параметром `db_type`:
- **`postgres`** — хранение ссылок в PostgreSQL
- **`inmemory`** — хранение в памяти (данные теряются при перезапуске)

Пример:
```yaml
db_type: postgres
```

### 6️⃣ Запустите сервис
```sh
task start
```
Теперь сервис доступен по адресу [`http://localhost:8080`](http://localhost:8080).

### 7️⃣ Остановите сервис
```sh
task stop
```

## Генерация моков

- **Сгенерировать моки:** `task gen-mocks`

## Тестирование

- **Юнит-тесты:** `task unit-test`

## Работа с базой

- **Применить миграции:** `task migrate-up`
- **Откатить миграцию:** `task migrate-down`
- **Создать новую миграцию:** `task migrate-create name=<название>`