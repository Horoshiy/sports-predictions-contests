# Обзор API сервисов

Полный справочник всех микросервисов платформы конкурсов спортивных прогнозов.

## Обзор архитектуры

Платформа состоит из 7 основных микросервисов, взаимодействующих через gRPC, с доступом к HTTP REST API через API Gateway.

| Сервис | Порт | Базовый URL | Назначение |
|---------|------|-------------|------------|
| **API Gateway** | 8080 | `http://localhost:8080` | HTTP REST точка входа |
| **Пользовательский сервис** | 8084 | `/v1/auth/*`, `/v1/users/*` | Аутентификация и управление пользователями |
| **Сервис конкурсов** | 8085 | `/v1/contests/*` | Управление конкурсами и командами |
| **Сервис прогнозов** | 8086 | `/v1/predictions/*`, `/v1/events/*` | Прогнозы и события |
| **Сервис подсчета очков** | 8087 | `/v1/scores/*`, `/v1/leaderboard/*` | Подсчет очков и таблицы лидеров |
| **Спортивный сервис** | 8088 | `/v1/sports/*`, `/v1/leagues/*` | Спортивные данные и синхронизация |
| **Сервис уведомлений** | 8089 | `/v1/notifications/*` | Многоканальные уведомления |

## Аутентификация

Все API эндпоинты (кроме регистрации и входа) требуют JWT аутентификации:

```bash
# Включите JWT токен в заголовок Authorization
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/v1/contests
```

## API Gateway (Порт 8080)

### Проверка здоровья
```bash
GET /health
```

**Ответ:**
```json
{
  "status": "ok",
  "timestamp": "2026-01-18T12:00:00Z"
}
```

## Пользовательский сервис (Порт 8084)

### Эндпоинты аутентификации

#### Регистрация пользователя
```bash
POST /v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecureP@ssw0rd2026!",
  "full_name": "Иван Иванов"
}
```

#### Вход пользователя
```bash
POST /v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecureP@ssw0rd2026!"
}
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "Иван Иванов"
  }
}
```

## Сервис конкурсов (Порт 8085)

### Управление конкурсами

#### Создание конкурса
```bash
POST /v1/contests
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "title": "Прогнозы Премьер-лиги",
  "description": "Прогнозируйте результаты матчей Премьер-лиги",
  "sport_type": "football",
  "rules": "{\"scoring_system\": \"standard\"}",
  "start_date": "2026-01-20T00:00:00Z",
  "end_date": "2026-05-20T23:59:59Z",
  "max_participants": 100
}
```

#### Список конкурсов
```bash
GET /v1/contests?page=1&limit=10&status=active&sport_type=football
```

#### Присоединение к конкурсу
```bash
POST /v1/contests/{contest_id}/join
Authorization: Bearer JWT_TOKEN
```

## Сервис прогнозов (Порт 8086)

### Управление прогнозами

#### Подача прогноза
```bash
POST /v1/predictions
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "contest_id": 1,
  "event_id": 1,
  "prediction_type": "match_outcome",
  "prediction_value": "home_win"
}
```

## Сервис подсчета очков (Порт 8087)

### Таблица лидеров

#### Получение таблицы лидеров конкурса
```bash
GET /v1/contests/{contest_id}/leaderboard?page=1&limit=50
```

**Ответ:**
```json
{
  "leaderboard": [
    {
      "rank": 1,
      "user_id": 1,
      "username": "johndoe",
      "total_score": 150,
      "correct_predictions": 15,
      "total_predictions": 20,
      "accuracy": 75.0,
      "current_streak": 5
    }
  ]
}
```

## Обработка ошибок

Все сервисы возвращают согласованные ответы об ошибках:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Неверные входные данные",
    "details": {
      "field": "email",
      "reason": "Неверный формат email"
    }
  }
}
```

### Общие коды ошибок

| Код | HTTP статус | Описание |
|-----|-------------|----------|
| `VALIDATION_ERROR` | 400 | Неверные данные запроса |
| `UNAUTHORIZED` | 401 | Отсутствует или неверный JWT токен |
| `FORBIDDEN` | 403 | Недостаточно прав |
| `NOT_FOUND` | 404 | Ресурс не найден |
| `CONFLICT` | 409 | Ресурс уже существует |
| `INTERNAL_ERROR` | 500 | Ошибка сервера |

---

Для подробной документации отдельных сервисов см.:
- [API пользовательского сервиса](user-service.md)
- [API сервиса конкурсов](contest-service.md)
- [API сервиса прогнозов](prediction-service.md)
- [API сервиса подсчета очков](scoring-service.md)
- [API спортивного сервиса](sports-service.md)
- [API сервиса уведомлений](notification-service.md)
