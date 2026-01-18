# Руководство по быстрому старту

Запустите платформу конкурсов спортивных прогнозов за несколько минут с помощью этого подробного руководства.

## Предварительные требования

Перед началом убедитесь, что в вашей системе установлено следующее:

### Необходимое программное обеспечение
- **Docker** (20.10+) и **Docker Compose** (2.0+)
- **Go** (1.21+) для разработки бэкенда
- **Node.js** (18+) и **npm** (8+) для разработки фронтенда
- **Git** для контроля версий

### Системные требования
- **ОЗУ**: Минимум 4ГБ, рекомендуется 8ГБ
- **Хранилище**: Не менее 2ГБ свободного места
- **ОС**: Linux, macOS или Windows с WSL2

### Команды проверки
```bash
# Проверка установки Docker
docker --version
docker-compose --version

# Проверка установки Go
go version

# Проверка установки Node.js
node --version
npm --version

# Проверка установки Git
git --version
```

## Шаги установки

### Шаг 1: Клонирование репозитория

```bash
# Клонирование репозитория
git clone https://github.com/coleam00/dynamous-kiro-hackathon
cd dynamous-kiro-hackathon

# Проверка, что вы в правильной директории
ls -la
```

### Шаг 2: Настройка окружения

```bash
# Запуск автоматизированного скрипта настройки
make setup

# Эта команда выполнит:
# - Копирование .env.example в .env
# - Установку зависимостей Go
# - Установку зависимостей Node.js
# - Настройку среды разработки
```

### Шаг 3: Запуск инфраструктурных сервисов

```bash
# Запуск PostgreSQL и Redis
make docker-up

# Проверка работы сервисов
docker-compose ps
```

Ожидаемый вывод:
```
NAME                COMMAND                  SERVICE             STATUS              PORTS
sports-postgres     "docker-entrypoint.s…"   postgres            running             0.0.0.0:5432->5432/tcp
sports-redis        "docker-entrypoint.s…"   redis               running             0.0.0.0:6379->6379/tcp
```

### Шаг 4: Инициализация базы данных

База данных будет автоматически инициализирована с необходимой схемой при запуске PostgreSQL. Вы можете проверить соединение:

```bash
# Тест соединения с базой данных
docker-compose exec postgres pg_isready -U sports_user -d sports_prediction
```

### Шаг 5: Запуск всех сервисов (режим разработки)

```bash
# Запуск всех микросервисов и фронтенда
make dev

# Это запустит:
# - Все 7 микросервисов
# - React фронтенд
# - Telegram бот (если настроен)
```

### Шаг 6: Проверка установки

#### Проверка здоровья сервисов
```bash
# Проверка здоровья API Gateway
curl http://localhost:8080/health

# Ожидаемый ответ: {"status": "ok", "timestamp": "..."}
```

#### Проверка отдельных сервисов
```bash
# Пользовательский сервис
curl http://localhost:8080/v1/auth/health

# Сервис конкурсов
curl http://localhost:8080/v1/contests/health

# Сервис прогнозов
curl http://localhost:8080/v1/predictions/health

# Сервис подсчета очков
curl http://localhost:8080/v1/scores/health

# Спортивный сервис
curl http://localhost:8080/v1/sports/health

# Сервис уведомлений
curl http://localhost:8080/v1/notifications/health
```

#### Доступ к фронтенду
Откройте браузер и перейдите по адресу:
- **Фронтенд**: http://localhost:3000
- **API Gateway**: http://localhost:8080

## Базовое использование

### Создание первого пользователя

```bash
# Регистрация нового пользователя
# ⚠️ БЕЗОПАСНОСТЬ: Используйте надежные пароли в продакшн (мин 12 символов, разный регистр, цифры, символы)
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "SecureP@ssw0rd2026!",
    "full_name": "Тестовый пользователь"
  }'
```

### Вход и получение токена

```bash
# Вход для получения JWT токена
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecureP@ssw0rd2026!"
  }'

# Сохраните полученный токен для последующих запросов
export JWT_TOKEN="your_jwt_token_here"
```

### Создание первого конкурса

```bash
# Создание конкурса
curl -X POST http://localhost:8080/v1/contests \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "title": "Прогнозы Премьер-лиги",
    "description": "Прогнозируйте результаты матчей Премьер-лиги",
    "sport_type": "football",
    "rules": "{\"scoring_system\": \"standard\", \"max_predictions_per_user\": 10}",
    "start_date": "2026-01-20T00:00:00Z",
    "end_date": "2026-05-20T23:59:59Z",
    "max_participants": 100
  }'
```

## Команды разработки

### Основные Make команды
```bash
# Показать все доступные команды
make help

# Настройка среды разработки
make setup

# Запуск среды разработки (только инфраструктура)
make dev

# Запуск всех сервисов включая микросервисы
make docker-services

# Сборка всех сервисов
make build

# Запуск тестов
make test

# Запуск E2E тестов
make e2e-test

# Очистка
make clean

# Показать статус сервисов
make status

# Просмотр логов
make logs
```

### Docker команды
```bash
# Запуск только инфраструктуры
docker-compose up -d postgres redis

# Запуск всех сервисов
docker-compose --profile services up -d

# Остановка всех сервисов
docker-compose --profile services down

# Просмотр логов
docker-compose logs -f [имя-сервиса]

# Перезапуск сервиса
docker-compose restart [имя-сервиса]
```

## Конфигурация

### Переменные окружения

Ключевые переменные окружения в `.env`:

```bash
# База данных
DATABASE_URL=postgres://sports_user:sports_password@localhost:5432/sports_prediction?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# API Gateway
API_GATEWAY_PORT=8080

# JWT
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRATION=24h

# Telegram бот (опционально)
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
TELEGRAM_ENABLED=false
```

Для полного списка см. [Справочник переменных окружения](environment-variables.md).

## Устранение неполадок

### Частые проблемы

#### Порт уже используется
```bash
# Проверить, что использует порт
lsof -i :8080

# Завершить процесс при необходимости
kill -9 <PID>
```

#### Проблемы с подключением к базе данных
```bash
# Проверить статус PostgreSQL
docker-compose logs postgres

# Перезапустить PostgreSQL
docker-compose restart postgres

# Сбросить базу данных
docker-compose down -v
docker-compose up -d postgres
```

#### Сервис не запускается
```bash
# Проверить логи сервиса
docker-compose logs [имя-сервиса]

# Пересобрать сервис
docker-compose build [имя-сервиса]
docker-compose up -d [имя-сервиса]
```

#### Фронтенд не загружается
```bash
# Проверить логи фронтенда
docker-compose logs frontend

# Пересобрать фронтенд
cd frontend
npm install
npm run build
```

### Получение помощи

Если у вас возникли проблемы:

1. Проверьте [Руководство по частым проблемам](../troubleshooting/common-issues.md)
2. Просмотрите логи сервисов: `docker-compose logs [имя-сервиса]`
3. Проверьте переменные окружения в `.env`
4. Убедитесь, что все предварительные требования установлены
5. Попробуйте перезапустить сервисы: `make clean && make dev`

## Следующие шаги

После запуска платформы:

1. **Изучите API**: Ознакомьтесь с [Документацией API](../api/services-overview.md)
2. **Запустите тесты**: Выполните набор тестов с помощью `make test`
3. **Настройте Telegram бота**: Настройте интеграцию с Telegram
4. **Настройте параметры**: Измените переменные окружения под ваши нужды
5. **Развертывание в продакшн**: Следуйте [Руководству по продакшн развертыванию](production.md)

## Контрольный список проверки

- [ ] Все предварительные требования установлены и проверены
- [ ] Репозиторий успешно клонирован
- [ ] Настройка окружения завершена (`make setup`)
- [ ] Инфраструктурные сервисы работают (PostgreSQL, Redis)
- [ ] Все микросервисы здоровы
- [ ] Фронтенд доступен по адресу http://localhost:3000
- [ ] API Gateway отвечает по адресу http://localhost:8080
- [ ] Регистрация и вход пользователя работают
- [ ] Создание конкурса успешно

---

**Поздравляем!** Теперь у вас есть полностью функциональная платформа конкурсов спортивных прогнозов, работающая локально. Изучите [документацию API](../api/services-overview.md), чтобы узнать больше о доступных функциях.
