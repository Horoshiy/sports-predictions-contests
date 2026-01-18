# Диагностические инструменты и отладка

Подробное руководство по диагностическим инструментам и процедурам отладки для платформы конкурсов спортивных прогнозов.

## Диагностика Docker

### Мониторинг здоровья контейнеров

#### Проверка статуса контейнеров
```bash
# Список всех контейнеров со статусом
docker-compose ps

# Показать детальную информацию о контейнере
docker inspect <имя_контейнера>

# Проверить использование ресурсов контейнера
docker stats --no-stream

# Мониторинг использования ресурсов в реальном времени
docker stats
```

#### Анализ логов контейнеров
```bash
# Просмотр логов всех сервисов
docker-compose logs

# Просмотр логов конкретного сервиса
docker-compose logs user-service

# Следить за логами в реальном времени
docker-compose logs -f api-gateway

# Просмотр последних N строк логов
docker-compose logs --tail=50 postgres

# Поиск логов по конкретным паттернам
docker-compose logs | grep -i "error\|warning\|failed"

# Экспорт логов в файл
docker-compose logs > system-logs.txt
```

### Диагностика сети

#### Анализ Docker сети
```bash
# Список Docker сетей
docker network ls

# Проверка конфигурации сети
docker network inspect sports-network

# Проверка сетевого подключения между контейнерами
docker-compose exec api-gateway ping user-service
docker-compose exec user-service ping postgres
```

#### Порты и обнаружение сервисов
```bash
# Проверка, какие порты открыты
docker-compose port api-gateway 8080
docker-compose port postgres 5432

# Тест подключения к сервисам
docker-compose exec api-gateway curl http://user-service:8084/health
docker-compose exec user-service curl http://postgres:5432
```

## Диагностика базы данных

### Проверки здоровья PostgreSQL

#### Тестирование подключения
```bash
# Тест подключения к базе данных
docker-compose exec postgres pg_isready -U sports_user -d sports_prediction

# Подключение к оболочке базы данных
docker-compose exec postgres psql -U sports_user -d sports_prediction

# Проверка версии и статуса базы данных
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "SELECT version();"
```

#### Анализ производительности базы данных
```bash
# Проверка активных подключений
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT count(*) as active_connections 
FROM pg_stat_activity 
WHERE state = 'active';"

# Проверка размера базы данных
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT pg_size_pretty(pg_database_size('sports_prediction')) as db_size;"

# Проверка размеров таблиц
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT schemaname,tablename,pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size 
FROM pg_tables 
WHERE schemaname = 'public' 
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
```

### Диагностика Redis

#### Здоровье и производительность Redis
```bash
# Тест подключения к Redis
docker-compose exec redis redis-cli ping

# Проверка информации Redis
docker-compose exec redis redis-cli info

# Мониторинг команд Redis в реальном времени
docker-compose exec redis redis-cli monitor

# Проверка использования памяти Redis
docker-compose exec redis redis-cli info memory

# Список всех ключей (осторожно использовать в продакшн)
docker-compose exec redis redis-cli keys "*"

# Проверка конкретных типов ключей и значений
docker-compose exec redis redis-cli type "leaderboard:contest:<CONTEST_ID>"
docker-compose exec redis redis-cli get "session:user:<USER_ID>"
```

## Диагностика конкретных сервисов

### Диагностика API Gateway

#### Проверки здоровья и статуса
```bash
# Проверка здоровья API Gateway
curl -v http://localhost:8080/health

# Тест маршрутизации к сервисам
curl -v http://localhost:8080/v1/auth/health
curl -v http://localhost:8080/v1/contests/health

# Проверка конфигурации API Gateway
docker-compose exec api-gateway env | grep -E "(PORT|ENDPOINT|JWT)"
```

### Проверки здоровья микросервисов

#### Тестирование отдельных сервисов
```bash
# Тест эндпоинта здоровья каждого сервиса
services=("auth" "contests" "predictions" "scores" "sports" "notifications")

for service in "${services[@]}"; do
    echo "Тестирование сервиса $service:"
    curl -s http://localhost:8080/v1/$service/health | jq '.'
    echo "---"
done
```

## Диагностика производительности

### Мониторинг системных ресурсов

#### Анализ CPU и памяти
```bash
# Проверка системных ресурсов
free -h
top -bn1 | head -20

# Мониторинг ресурсов Docker контейнеров
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"

# Проверка использования диска
df -h
docker system df
```

## Инструменты анализа логов

### Анализ структурированных логов

#### Парсинг и фильтрация логов
```bash
# Парсинг JSON логов (если используется структурированное логирование)
docker-compose logs user-service | jq -r 'select(.level == "error") | .message'

# Фильтрация логов по времени
docker-compose logs --since="2026-01-18T10:00:00" api-gateway

# Подсчет вхождений ошибок
docker-compose logs | grep -c "ERROR"

# Извлечение уникальных сообщений об ошибках
docker-compose logs | grep "ERROR" | sort | uniq -c | sort -nr
```

## Рабочие процессы отладки

### Процесс исследования проблем

#### Шаг 1: Первоначальная оценка
```bash
# Быстрая проверка здоровья системы
echo "=== Проверка здоровья системы ==="
docker-compose ps
curl -s http://localhost:8080/health | jq '.'
docker stats --no-stream | head -10
```

#### Шаг 2: Исследование конкретного сервиса
```bash
# Проверка конкретного сервиса, который не работает
SERVICE_NAME="user-service"

echo "=== Исследование $SERVICE_NAME ==="
docker-compose logs --tail=50 $SERVICE_NAME
docker-compose exec $SERVICE_NAME env | grep -E "(DATABASE|REDIS|JWT)"
docker inspect $(docker-compose ps -q $SERVICE_NAME)
```

### Общие сценарии отладки

#### Проблемы подключения к базе данных
```bash
# Отладка подключения к базе данных
echo "=== Отладка базы данных ==="
docker-compose logs postgres | tail -20
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "\l"
docker-compose exec user-service nc -zv postgres 5432
```

#### Проблемы аутентификации
```bash
# Отладка JWT аутентификации
echo "=== Отладка JWT ==="
echo "Проверка JWT_SECRET:"
docker-compose exec api-gateway env | grep JWT_SECRET
docker-compose exec user-service env | grep JWT_SECRET

# Тест генерации токена
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -v
```

## Автоматизированные диагностические скрипты

### Скрипт проверки здоровья
```bash
#!/bin/bash
# сохранить как scripts/health-check.sh

echo "=== Проверка здоровья платформы спортивных прогнозов ==="
echo "Время: $(date)"
echo

# Проверка Docker
echo "Статус Docker:"
docker --version
docker-compose --version
echo

# Проверка контейнеров
echo "Статус контейнеров:"
docker-compose ps
echo

# Проверка сервисов
echo "Здоровье сервисов:"
services=("health" "v1/auth/health" "v1/contests/health" "v1/predictions/health" "v1/scores/health" "v1/sports/health" "v1/notifications/health")

for endpoint in "${services[@]}"; do
    echo -n "$endpoint: "
    if curl -s -f http://localhost:8080/$endpoint > /dev/null; then
        echo "✅ OK"
    else
        echo "❌ FAIL"
    fi
done
echo

# Проверка базы данных
echo "Статус базы данных:"
if docker-compose exec -T postgres pg_isready -U sports_user -d sports_prediction > /dev/null 2>&1; then
    echo "PostgreSQL: ✅ OK"
else
    echo "PostgreSQL: ❌ FAIL"
fi

# Проверка Redis
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo "Redis: ✅ OK"
else
    echo "Redis: ❌ FAIL"
fi

echo
echo "=== Проверка здоровья завершена ==="
```

## Контрольный список устранения неполадок

### Контрольный список перед отладкой
- [ ] Все контейнеры запущены (`docker-compose ps`)
- [ ] Нет конфликтов портов (`lsof -i :8080`)
- [ ] Достаточно системных ресурсов (`free -h`, `df -h`)
- [ ] Переменные окружения установлены (файл `.env` существует)
- [ ] Сетевое подключение между сервисами
- [ ] База данных и Redis доступны

### Во время исследования
- [ ] Проверить логи сервисов на ошибки
- [ ] Проверить конфигурацию сервиса
- [ ] Протестировать эндпоинты отдельных сервисов
- [ ] Проверить подключение и производительность базы данных
- [ ] Мониторить системные ресурсы
- [ ] Отследить поток запросов через сервисы

---

*Эти диагностические инструменты помогают быстро выявлять и решать проблемы. Используйте их систематически при исследовании проблем с платформой конкурсов спортивных прогнозов.*
