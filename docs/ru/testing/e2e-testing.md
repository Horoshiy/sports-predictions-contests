# Руководство по E2E тестированию

Подробное руководство по запуску E2E тестов на платформе конкурсов спортивных прогнозов.

## Обзор

Набор E2E тестов проверяет полные пользовательские сценарии работы всех микросервисов, обеспечивая корректную работу платформы от начала до конца.

## Архитектура тестов

### Структура тестов
```
tests/e2e/
├── main_test.go          # Запуск тестов и настройка
├── helpers.go            # Утилиты и помощники тестов
├── types.go              # Структуры данных тестов
├── auth_test.go          # Сценарии аутентификации
├── contest_test.go       # Тесты управления конкурсами
├── prediction_test.go    # Тесты подачи прогнозов
├── scoring_test.go       # Тесты подсчета очков и таблиц лидеров
├── sports_test.go        # Тесты управления спортивными данными
└── workflow_test.go      # Полные пользовательские сценарии
```

### Тестовая среда
- **Изолированная Docker среда**: Тесты выполняются против контейнеризованных сервисов
- **Свежая база данных**: Каждый запуск тестов начинается с чистой базы данных
- **Реальное взаимодействие сервисов**: Тесты используют фактическое gRPC/HTTP взаимодействие
- **Автоматическая настройка/очистка**: Инфраструктура управляется автоматически

## Предварительные требования

### Необходимое программное обеспечение
- **Docker** и **Docker Compose**
- **Go** 1.21+ для запуска тестов
- **curl** для ручного тестирования API
- **jq** для обработки JSON (опционально)

### Настройка окружения
```bash
# Убедитесь, что вы в корне проекта
cd /path/to/sports-prediction-contests

# Проверьте, что Docker запущен
docker --version
docker-compose --version

# Проверьте установку Go
go version
```

## Запуск E2E тестов

### Автоматическое выполнение тестов

#### Полный набор E2E тестов
```bash
# Запуск полного набора E2E тестов с настройкой инфраструктуры
make e2e-test

# Эта команда выполнит:
# 1. Запуск PostgreSQL и Redis
# 2. Ожидание готовности базы данных
# 3. Запуск всех микросервисов
# 4. Ожидание проверок здоровья сервисов
# 5. Запуск E2E тестов
# 6. Очистка всех сервисов
```

#### Только E2E тесты (сервисы уже запущены)
```bash
# Если сервисы уже запущены
make e2e-test-only

# Или запуск напрямую
cd tests/e2e
go test -tags=e2e -v -timeout 5m ./...
```

### Ручное выполнение тестов

#### Шаг 1: Запуск инфраструктуры
```bash
# Запуск PostgreSQL и Redis
docker-compose up -d postgres redis

# Ожидание готовности базы данных
docker-compose exec postgres pg_isready -U sports_user -d sports_prediction
```

#### Шаг 2: Запуск всех сервисов
```bash
# Запуск всех микросервисов
docker-compose --profile services up -d

# Ожидание готовности сервисов (15-30 секунд)
sleep 15

# Проверка здоровья API Gateway
curl http://localhost:8080/health
```

#### Шаг 3: Запуск тестов
```bash
# Запуск конкретных файлов тестов
cd tests/e2e

# Тесты аутентификации
go test -tags=e2e -v -run TestAuth

# Тесты управления конкурсами
go test -tags=e2e -v -run TestContest

# Полные сценарии работы
go test -tags=e2e -v -run TestWorkflow

# Все тесты с подробным выводом
go test -tags=e2e -v -timeout 5m ./...
```

#### Шаг 4: Очистка
```bash
# Остановка всех сервисов
docker-compose --profile services down

# Остановка инфраструктуры
docker-compose down
```

## Тестовые сценарии

### Тесты сценариев аутентификации

#### Регистрация и вход пользователя
```go
func TestUserRegistrationAndLogin(t *testing.T) {
    // Тест регистрации пользователя
    user := registerTestUser(t, "testuser", "test@example.com", "password123")
    
    // Тест входа пользователя
    token := loginUser(t, "test@example.com", "password123")
    
    // Проверка действительности токена
    profile := getUserProfile(t, token)
    assert.Equal(t, "testuser", profile.Username)
}
```

**Ручной тест:**
```bash
# Регистрация пользователя
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com", 
    "password": "password123",
    "full_name": "Тестовый пользователь"
  }'

# Вход пользователя
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Сохранить токен для последующих запросов
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Тесты управления конкурсами

#### Создание конкурса и участие
```go
func TestContestCreationAndParticipation(t *testing.T) {
    // Создание пользователя и получение токена
    token := setupTestUser(t)
    
    // Создание конкурса
    contest := createTestContest(t, token, "Тестовый конкурс")
    
    // Присоединение к конкурсу
    joinContest(t, token, contest.ID)
    
    // Проверка участия
    participants := getContestParticipants(t, token, contest.ID)
    assert.Len(t, participants, 1)
}
```

**Ручной тест:**
```bash
# Создание конкурса
curl -X POST http://localhost:8080/v1/contests \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Тестовый конкурс",
    "description": "E2E тестовый конкурс",
    "sport_type": "football",
    "rules": "{\"scoring_system\": \"standard\"}",
    "start_date": "2026-01-20T00:00:00Z",
    "end_date": "2026-05-20T23:59:59Z",
    "max_participants": 100
  }'

# Присоединение к конкурсу (сохранить contest_id из предыдущего ответа)
curl -X POST http://localhost:8080/v1/contests/1/join \
  -H "Authorization: Bearer $JWT_TOKEN"
```

## Управление тестовыми данными

### Создание тестового пользователя
```go
type TestUser struct {
    ID       uint32
    Username string
    Email    string
    Token    string
}

func setupTestUser(t *testing.T) *TestUser {
    username := fmt.Sprintf("testuser_%d", time.Now().Unix())
    email := fmt.Sprintf("%s@example.com", username)
    
    // Регистрация пользователя
    user := registerTestUser(t, username, email, "password123")
    
    // Вход и получение токена
    token := loginUser(t, email, "password123")
    
    return &TestUser{
        ID:       user.ID,
        Username: username,
        Email:    email,
        Token:    token,
    }
}
```

## Конфигурация тестов

### Переменные окружения
```bash
# Конфигурация тестов в tests/e2e/.env
BASE_URL=http://localhost:8080
DATABASE_URL=postgres://sports_user:sports_password@localhost:5432/sports_prediction?sslmode=disable
REDIS_URL=redis://localhost:6379
TEST_TIMEOUT=5m
LOG_LEVEL=debug
```

### Флаги тестов
```bash
# Запуск с подробным выводом
go test -tags=e2e -v

# Запуск конкретного паттерна тестов
go test -tags=e2e -run TestAuth

# Запуск с таймаутом
go test -tags=e2e -timeout 10m

# Запуск с обнаружением гонок
go test -tags=e2e -race

# Генерация отчета о покрытии
go test -tags=e2e -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Отладка неудачных тестов

### Частые проблемы и решения

#### Сервисы не готовы
```bash
# Проверить статус сервисов
docker-compose ps

# Проверить логи сервисов
docker-compose logs api-gateway
docker-compose logs user-service

# Подождать дольше запуска сервисов
sleep 30
curl http://localhost:8080/health
```

#### Проблемы подключения к базе данных
```bash
# Проверить статус PostgreSQL
docker-compose logs postgres

# Тест подключения к базе данных
docker-compose exec postgres pg_isready -U sports_user -d sports_prediction

# Сброс базы данных
docker-compose down -v
docker-compose up -d postgres
```

## Контрольный список устранения неполадок

- [ ] Docker и Docker Compose установлены и запущены
- [ ] Все необходимые порты доступны (8080-8089, 5432, 6379)
- [ ] Сервисы запущены в правильном порядке (сначала инфраструктура)
- [ ] База данных инициализирована и доступна
- [ ] Проверка здоровья API Gateway проходит
- [ ] JWT токены действительны и не истекли
- [ ] Очистка тестовых данных между запусками
- [ ] Достаточно системных ресурсов (память, диск)
- [ ] Сетевое подключение между сервисами
- [ ] Переменные окружения правильно установлены

---

**E2E тесты обеспечивают уверенность в том, что вся платформа работает корректно. Запускайте их регулярно во время разработки и всегда перед развертыванием.**
