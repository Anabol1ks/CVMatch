# CVMatch

CVMatch — backend сервис для парсинга, хранения и поиска резюме, реализованный на Go. Поддерживает регистрацию/авторизацию пользователей, загрузку и структурирование резюме, хранение данных в PostgreSQL, предоставляет REST API с документацией Swagger и покрыт тестами. Парсинг резюме реализован с помощью YandexGPT.

---

## ⚡ Быстрый старт

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/Anabol1ks/CVMatch.git
cd CVMatch
```

### 2. Настройка переменных окружения

- Для **Linux/macOS**:
    ```bash
    cp .env.example .env
    ```
- Для **Windows** (cmd):
    ```cmd
    copy .env.example .env
    ```
- Обязательно отредактируйте файл `.env`:
    - Проверьте и задайте свои значения для всех секретов (`ACCESS_SECRET`, `REFRESH_SECRET`).
    - Укажите параметры подключения к базе данных.
    - Установите параметры интеграции с YandexGPT (`YANDEXGPT_IAM`, `YANDEXGPT_CATALOG_ID`).
    - Проверьте и настройте `BASE_URL`, `ACCESS_EXP`, `REFRESH_EXP` и другие переменные под ваш сценарий.

### 3. Запуск через Docker

- Убедитесь, что установлен Docker и Docker Compose.
- Запустите сервисы:
    ```bash
    docker-compose up -d --build
    ```
- Будут запущены контейнеры:
    - PostgreSQL (порт 5432)
    - Основной сервис (порт 8080)

### 4. Запуск локально (без Docker)

- Установите **Go 1.24+**.
- Убедитесь, что PostgreSQL запущен и параметры в `.env` корректны.
- Запустите приложение:
    ```bash
    make run
    ```
    или вручную:
    ```bash
    go run cmd/main.go
    ```

---

## 🧩 Makefile

Доступные команды:

- `make run` — запуск приложения локально
- `make swag` — генерация Swagger-документации
- `make doc` — запуск docker-compose (БД и сервисы)
- `make test` — запуск всех тестов

---

## 📝 Swagger-документация

- Актуальная спецификация находится в `docs/swagger.yaml` и `docs/swagger.json`.
- Для генерации/обновления используйте:
    ```bash
    make swag
    ```
- Документация API доступна в браузере по адресу:
    [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🗄️ Структура проекта

- **cmd/** — точка входа (main.go)
- **internal/** — бизнес-логика (handlers, service, repository, models и др.)
- **docs/** — OpenAPI/Swagger спецификация
- **testdata/** — тестовые файлы, если есть

---

## 🚦 Тесты

- Покрывают слой репозиториев и сервисов.
- Запуск всех тестов:
    ```bash
    make test
    ```
    или
    ```bash
    go test ./...
    ```

---

## 🔒 Авторизация

- JWT авторизация для всех защищённых эндпоинтов.
- Секреты и параметры токенов настраиваются в `.env` (`ACCESS_SECRET`, `REFRESH_SECRET`, `ACCESS_EXP`, `REFRESH_EXP` и др).

---

## 🤖 Парсинг резюме

- Парсинг резюме реализован с помощью **YandexGPT**.
- Для работы YandexGPT необходимо заполнить параметры `YANDEXGPT_IAM` и `YANDEXGPT_CATALOG_ID` в `.env`.

---

## 🌍 Roadmap (дальнейшее развитие)

- [x] Авторизация и регистрация пользователей
- [x] CRUD для резюме, загрузка и парсинг PDF
- [x] Хранение структуры резюме (скиллы, опыт, образование)
- [x] Документация Swagger
- [x] Покрытие тестами (repository, service)
- [ ] CRUD для вакансий (Job Description)
- [ ] Алгоритм сравнения резюме и вакансий (matching, scoring)
- [ ] Улучшение парсера (поддержка разных форматов)
- [ ] Интеграция с внешними сервисами и AI (YandexGPT)
- [ ] E2E и HTTP тесты для API
- [ ] CI/CD pipeline (GitHub Actions)

---

## 🛠️ Контрибьютинг

Pull requests и идеи приветствуются! Перед пушем убедитесь, что проходят все тесты и линтеры.

---

## 📄 Лицензия

[MIT](./LICENSE)

---

**Контакты:**  
Автор: [Anabol1ks](https://github.com/Anabol1ks)