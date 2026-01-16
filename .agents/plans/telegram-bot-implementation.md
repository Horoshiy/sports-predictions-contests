# Feature: Telegram Bot Implementation

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Полноценный Telegram бот для платформы Sports Prediction Contests, позволяющий пользователям:
- Регистрироваться и привязывать Telegram аккаунт к платформе
- Просматривать активные конкурсы и события
- Делать прогнозы прямо из Telegram
- Получать уведомления о результатах и обновлениях лидерборда
- Просматривать свою статистику и позицию в рейтинге

Бот интегрируется с существующим notification-service (который уже имеет TelegramChannel для отправки сообщений) и использует gRPC клиенты для взаимодействия с другими микросервисами.

## User Story

As a sports prediction platform user
I want to interact with the platform via Telegram bot
So that I can make predictions and track results without opening the web app

## Problem Statement

Директория `bots/telegram/` пустая, хотя Telegram бот заявлен как ключевая фича платформы в README. Существующий notification-service умеет только отправлять сообщения через Telegram, но не обрабатывает входящие команды от пользователей.

## Solution Statement

Создать standalone Telegram бот как отдельный Go сервис в `bots/telegram/`, который:
1. Использует long polling для получения обновлений от Telegram API
2. Обрабатывает команды пользователей через command handler pattern
3. Интегрируется с микросервисами через gRPC клиенты
4. Автоматически привязывает Telegram chat_id к аккаунту пользователя

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium (4-6 hours)
**Primary Systems Affected**: bots/telegram, notification-service, user-service
**Dependencies**: github.com/go-telegram-bot-api/telegram-bot-api/v5 (уже в notification-service)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/notification-service/internal/channels/telegram.go` - Why: Существующая интеграция с Telegram API, паттерн использования tgbotapi
- `backend/notification-service/internal/config/config.go` - Why: Паттерн конфигурации с env переменными
- `backend/notification-service/go.mod` - Why: Зависимости включая telegram-bot-api v5
- `backend/user-service/internal/service/user_service.go` - Why: Паттерн gRPC сервиса и аутентификации
- `backend/contest-service/internal/service/contest_service.go` (lines 1-80) - Why: Паттерн работы с контекстом и репозиториями
- `backend/proto/user.proto` - Why: gRPC API для аутентификации
- `backend/proto/contest.proto` - Why: gRPC API для конкурсов
- `backend/proto/prediction.proto` - Why: gRPC API для прогнозов
- `backend/proto/scoring.proto` - Why: gRPC API для лидерборда
- `backend/proto/notification.proto` - Why: gRPC API для настроек уведомлений
- `backend/notification-service/internal/models/notification.go` (lines 60-90) - Why: Модель NotificationPreference с TelegramChatID

### New Files to Create

- `bots/telegram/main.go` - Entry point с инициализацией бота и gRPC клиентов
- `bots/telegram/config/config.go` - Конфигурация из env переменных
- `bots/telegram/bot/bot.go` - Основная логика бота и обработка updates
- `bots/telegram/bot/handlers.go` - Command handlers (/start, /contests, /predict, etc.)
- `bots/telegram/bot/keyboards.go` - Inline keyboards для интерактивного UI
- `bots/telegram/bot/messages.go` - Шаблоны сообщений
- `bots/telegram/clients/clients.go` - gRPC клиенты к микросервисам
- `bots/telegram/go.mod` - Go module definition
- `bots/telegram/Dockerfile` - Docker image для бота
- `tests/telegram-bot/bot_test.go` - Unit тесты для handlers

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [go-telegram-bot-api GitHub](https://github.com/go-telegram-bot-api/telegram-bot-api)
  - Specific section: README example with GetUpdatesChan
  - Why: Паттерн long polling и обработки updates
- [Telegram Bot API](https://core.telegram.org/bots/api)
  - Specific section: Available methods, Inline keyboards
  - Why: Понимание возможностей API

### Patterns to Follow

**Naming Conventions:**
- Go files: `snake_case.go`
- Packages: `lowercase`
- Structs/Functions: `PascalCase` (public), `camelCase` (private)
- Config env vars: `TELEGRAM_BOT_TOKEN`, `USER_SERVICE_ENDPOINT`

**Error Handling (from notification-service):**
```go
if err != nil {
    log.Printf("[ERROR] Failed to X: %v", err)
    return fmt.Errorf("failed to X: %w", err)
}
```

**Config Pattern (from notification-service/internal/config/config.go):**
```go
func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

**gRPC Client Pattern:**
```go
conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    return nil, fmt.Errorf("failed to connect to service: %w", err)
}
client := pb.NewServiceClient(conn)
```

**Telegram Message Pattern (from channels/telegram.go):**
```go
msg := tgbotapi.NewMessage(chatID, text)
msg.ParseMode = "HTML"
_, err := bot.Send(msg)
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Создание базовой структуры бота с конфигурацией и инициализацией.

**Tasks:**
- Создать go.mod с зависимостями
- Создать конфигурацию из env переменных
- Создать базовую структуру Bot с инициализацией tgbotapi

### Phase 2: Core Implementation

Реализация command handlers и gRPC клиентов.

**Tasks:**
- Создать gRPC клиенты к микросервисам
- Реализовать command handlers (/start, /help, /contests, /leaderboard, /mystats, /link)
- Создать inline keyboards для навигации
- Реализовать callback query handlers для интерактивных кнопок

### Phase 3: Integration

Интеграция с существующей инфраструктурой.

**Tasks:**
- Добавить бота в docker-compose.yml
- Обновить .env.example с новыми переменными
- Создать Dockerfile для бота

### Phase 4: Testing & Validation

**Tasks:**
- Написать unit тесты для handlers
- Протестировать интеграцию с микросервисами
- Проверить работу inline keyboards

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### Task 1: CREATE `bots/telegram/go.mod`

- **IMPLEMENT**: Go module с зависимостями для Telegram бота
- **PATTERN**: Mirror `backend/notification-service/go.mod`
- **IMPORTS**: telegram-bot-api/v5, grpc, protobuf
- **VALIDATE**: `cd bots/telegram && go mod tidy` (после создания main.go)

```go
module github.com/sports-prediction-contests/telegram-bot

go 1.21

require (
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
    google.golang.org/grpc v1.59.0
    google.golang.org/protobuf v1.31.0
    github.com/sports-prediction-contests/shared v0.0.0
)

replace github.com/sports-prediction-contests/shared => ../../backend/shared
```

### Task 2: CREATE `bots/telegram/config/config.go`

- **IMPLEMENT**: Конфигурация бота из env переменных
- **PATTERN**: Mirror `backend/notification-service/internal/config/config.go`
- **IMPORTS**: os, strconv
- **GOTCHA**: Все endpoints должны быть настраиваемыми для Docker networking
- **VALIDATE**: `cd bots/telegram && go build ./config/...`

Конфигурация должна включать:
- `TELEGRAM_BOT_TOKEN` - токен бота от BotFather
- `USER_SERVICE_ENDPOINT` - адрес user-service (default: localhost:8084)
- `CONTEST_SERVICE_ENDPOINT` - адрес contest-service (default: localhost:8085)
- `PREDICTION_SERVICE_ENDPOINT` - адрес prediction-service (default: localhost:8086)
- `SCORING_SERVICE_ENDPOINT` - адрес scoring-service (default: localhost:8087)
- `NOTIFICATION_SERVICE_ENDPOINT` - адрес notification-service (default: localhost:8089)
- `LOG_LEVEL` - уровень логирования (default: info)

### Task 3: CREATE `bots/telegram/clients/clients.go`

- **IMPLEMENT**: gRPC клиенты ко всем микросервисам
- **PATTERN**: Standard gRPC client initialization
- **IMPORTS**: grpc, insecure credentials, proto packages
- **GOTCHA**: Использовать insecure credentials для внутренней сети Docker
- **VALIDATE**: `cd bots/telegram && go build ./clients/...`

Структура Clients должна содержать:
- UserClient для аутентификации
- ContestClient для списка конкурсов
- PredictionClient для создания прогнозов
- ScoringClient для лидерборда
- NotificationClient для привязки chat_id

### Task 4: CREATE `bots/telegram/bot/messages.go`

- **IMPLEMENT**: Константы и шаблоны сообщений на русском и английском
- **PATTERN**: Простые string constants с HTML форматированием
- **GOTCHA**: Использовать HTML parse mode (уже используется в notification-service)
- **VALIDATE**: `cd bots/telegram && go build ./bot/...`

Сообщения:
- Welcome message с инструкциями
- Help message со списком команд
- Error messages
- Success messages
- Форматы для конкурсов, лидерборда, статистики

### Task 5: CREATE `bots/telegram/bot/keyboards.go`

- **IMPLEMENT**: Inline keyboards для интерактивного UI
- **PATTERN**: tgbotapi.NewInlineKeyboardMarkup
- **IMPORTS**: tgbotapi
- **VALIDATE**: `cd bots/telegram && go build ./bot/...`

Keyboards:
- Main menu (Contests, Leaderboard, My Stats, Help)
- Contest list с кнопками выбора
- Confirmation keyboards (Yes/No)
- Back button

### Task 6: CREATE `bots/telegram/bot/handlers.go`

- **IMPLEMENT**: Command и callback handlers
- **PATTERN**: Switch-case по command/callback data
- **IMPORTS**: tgbotapi, clients, context
- **GOTCHA**: Всегда проверять nil для Message и CallbackQuery
- **VALIDATE**: `cd bots/telegram && go build ./bot/...`

Commands:
- `/start` - Welcome + main menu + инструкция по привязке аккаунта
- `/help` - Список команд
- `/contests` - Список активных конкурсов
- `/leaderboard [contest_id]` - Топ-10 лидерборда
- `/mystats` - Личная статистика (требует привязки)
- `/link <email> <password>` - Привязка Telegram к аккаунту платформы

Callbacks:
- `contest_<id>` - Детали конкурса
- `leaderboard_<id>` - Лидерборд конкурса
- `back_main` - Возврат в главное меню

### Task 7: CREATE `bots/telegram/bot/bot.go`

- **IMPLEMENT**: Основная структура Bot с методами Start и handleUpdate
- **PATTERN**: Long polling с GetUpdatesChan
- **IMPORTS**: tgbotapi, config, clients, log
- **GOTCHA**: Graceful shutdown при SIGINT/SIGTERM
- **VALIDATE**: `cd bots/telegram && go build ./bot/...`

Структура Bot:
- api *tgbotapi.BotAPI
- clients *clients.Clients
- config *config.Config
- userSessions map[int64]*UserSession (для хранения состояния диалога)

### Task 8: CREATE `bots/telegram/main.go`

- **IMPLEMENT**: Entry point с инициализацией всех компонентов
- **PATTERN**: Mirror `backend/notification-service/cmd/main.go`
- **IMPORTS**: config, bot, clients, log, os/signal
- **GOTCHA**: Graceful shutdown для gRPC connections
- **VALIDATE**: `cd bots/telegram && go build .`

### Task 9: CREATE `bots/telegram/Dockerfile`

- **IMPLEMENT**: Multi-stage build для минимального образа
- **PATTERN**: Mirror `backend/notification-service/Dockerfile`
- **GOTCHA**: Копировать shared module для replace directive
- **VALIDATE**: `docker build -t telegram-bot ./bots/telegram`

### Task 10: UPDATE `docker-compose.yml`

- **IMPLEMENT**: Добавить telegram-bot сервис
- **PATTERN**: Mirror notification-service configuration
- **GOTCHA**: Зависимости от user-service, contest-service, scoring-service
- **VALIDATE**: `docker-compose config`

Добавить сервис:
```yaml
telegram-bot:
  build:
    context: .
    dockerfile: bots/telegram/Dockerfile
  container_name: sports-telegram-bot
  environment:
    - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
    - USER_SERVICE_ENDPOINT=user-service:8084
    - CONTEST_SERVICE_ENDPOINT=contest-service:8085
    - PREDICTION_SERVICE_ENDPOINT=prediction-service:8086
    - SCORING_SERVICE_ENDPOINT=scoring-service:8087
    - NOTIFICATION_SERVICE_ENDPOINT=notification-service:8089
  depends_on:
    - user-service
    - contest-service
    - scoring-service
    - notification-service
  networks:
    - sports-network
  profiles:
    - services
```

### Task 11: UPDATE `.env.example`

- **IMPLEMENT**: Добавить TELEGRAM_BOT_TOKEN
- **VALIDATE**: `cat .env.example | grep TELEGRAM`

### Task 12: CREATE `tests/telegram-bot/bot_test.go`

- **IMPLEMENT**: Unit тесты для message formatting и keyboard generation
- **PATTERN**: Mirror `tests/notification-service/channels_test.go`
- **IMPORTS**: testing
- **VALIDATE**: `cd tests/telegram-bot && go test -v ./...`

Тесты:
- TestFormatContestMessage
- TestFormatLeaderboardMessage
- TestMainMenuKeyboard
- TestParseCommand

### Task 13: UPDATE `backend/go.work`

- **IMPLEMENT**: Добавить telegram-bot в Go workspace
- **VALIDATE**: `cd backend && go work sync`

---

## TESTING STRATEGY

### Unit Tests

- Тестирование форматирования сообщений
- Тестирование генерации keyboards
- Тестирование парсинга команд
- Mock gRPC clients для изолированного тестирования handlers

### Integration Tests

- Тестирование подключения к gRPC сервисам
- Тестирование полного flow: команда → gRPC call → ответ

### Edge Cases

- Пользователь не привязал аккаунт (команды требующие авторизации)
- Неверный формат команды /link
- Несуществующий contest_id
- Пустой лидерборд
- Сервис недоступен (graceful error handling)
- Слишком длинное сообщение (Telegram limit 4096 chars)

---

## VALIDATION COMMANDS

### Level 1: Syntax & Style

```bash
cd bots/telegram && go fmt ./...
cd bots/telegram && go vet ./...
```

### Level 2: Build

```bash
cd bots/telegram && go build .
```

### Level 3: Unit Tests

```bash
cd tests/telegram-bot && go test -v ./...
```

### Level 4: Docker Build

```bash
docker build -t telegram-bot -f bots/telegram/Dockerfile .
```

### Level 5: Integration Test

```bash
# Start services
docker-compose --profile services up -d

# Check bot logs
docker-compose logs telegram-bot
```

### Level 6: Manual Validation

1. Получить токен бота от @BotFather
2. Установить TELEGRAM_BOT_TOKEN в .env
3. Запустить бота: `cd bots/telegram && go run .`
4. Отправить /start боту в Telegram
5. Проверить inline keyboard
6. Проверить /contests (должен показать список или "нет активных")
7. Проверить /help

---

## ACCEPTANCE CRITERIA

- [ ] Бот запускается и отвечает на /start
- [ ] /help показывает список доступных команд
- [ ] /contests показывает список активных конкурсов (или сообщение что их нет)
- [ ] /leaderboard показывает топ-10 для указанного конкурса
- [ ] /link позволяет привязать Telegram к аккаунту платформы
- [ ] /mystats показывает статистику привязанного пользователя
- [ ] Inline keyboards работают корректно
- [ ] Graceful error handling при недоступности сервисов
- [ ] Docker image собирается успешно
- [ ] Бот добавлен в docker-compose.yml
- [ ] Unit тесты проходят

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms feature works
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

### Архитектурные решения

1. **Standalone сервис vs интеграция в notification-service**: Выбран standalone для separation of concerns. Notification-service отвечает за исходящие уведомления, telegram-bot — за входящие команды.

2. **Long polling vs Webhooks**: Выбран long polling для простоты развертывания (не требует публичного HTTPS endpoint).

3. **Хранение сессий**: In-memory map для простоты. Для production рекомендуется Redis.

4. **Привязка аккаунта**: Через команду /link с email/password. Альтернатива — deep link с токеном из веб-приложения.

### Безопасность

- Пароль передается в открытом виде в Telegram (не идеально, но приемлемо для MVP)
- Для production: использовать одноразовые токены или OAuth flow
- Не логировать пароли и токены

### Расширения (post-MVP)

- Inline mode для быстрых прогнозов
- Групповые чаты с ботом
- Напоминания о предстоящих матчах
- Webhook mode для production
