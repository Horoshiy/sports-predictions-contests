# Частые проблемы и решения

Это руководство охватывает наиболее часто встречающиеся проблемы при работе с платформой конкурсов спортивных прогнозов и их решения.

## Проблемы запуска сервисов

### Проблема: Сервисы не запускаются
**Симптомы**: Docker контейнеры немедленно завершаются или показывают сообщения об ошибках

**Решения**:
```bash
# Проверить, что Docker daemon запущен
docker --version
sudo systemctl start docker  # Linux
# или перезапустить Docker Desktop на macOS/Windows

# Проверить доступность портов
lsof -i :8080  # Проверить, используется ли порт
kill -9 <PID>  # Завершить процесс, использующий порт

# Перезапустить сервисы
make clean
make dev
```

### Проблема: Не удается подключиться к базе данных
**Симптомы**: Сервисы не могут подключиться к PostgreSQL

**Решения**:
```bash
# Проверить статус PostgreSQL
docker-compose logs postgres

# Перезапустить PostgreSQL
docker-compose restart postgres

# Полностью сбросить базу данных
docker-compose down -v
docker-compose up -d postgres
```

## Проблемы конфигурации

### Проблема: Переменные окружения не загружаются
**Симптомы**: Сервисы используют значения по умолчанию вместо настроенных

**Решения**:
```bash
# Проверить, что файл .env существует
ls -la .env

# Скопировать из примера, если отсутствует
cp .env.example .env

# Проверить, что переменные окружения загружены
docker-compose config
```

### Проблема: JWT аутентификация не работает
**Симптомы**: Ошибки 401 Unauthorized для аутентифицированных запросов

**Решения**:
```bash
# Проверить, что JWT секрет установлен
grep JWT_SECRET .env

# Сгенерировать новый JWT секрет
openssl rand -base64 32

# Обновить файл .env с новым секретом
JWT_SECRET=your_new_secret_here
```

## Проблемы сети и подключения

### Проблема: API Gateway не отвечает
**Симптомы**: Отказ в соединении на порту 8080

**Решения**:
```bash
# Проверить статус API Gateway
curl http://localhost:8080/health

# Проверить статус контейнера
docker-compose ps api-gateway

# Просмотреть логи API Gateway
docker-compose logs api-gateway

# Перезапустить API Gateway
docker-compose restart api-gateway
```

### Проблема: Проблемы обнаружения сервисов
**Симптомы**: Сервисы не могут взаимодействовать друг с другом

**Решения**:
```bash
# Проверить Docker сеть
docker network ls
docker network inspect sports-network

# Перезапустить все сервисы
docker-compose down
docker-compose up -d
```

## Проблемы производительности

### Проблема: Медленные ответы API
**Симптомы**: Запросы выполняются дольше ожидаемого

**Решения**:
```bash
# Проверить системные ресурсы
docker stats

# Проверить производительность базы данных
docker-compose exec postgres pg_stat_activity

# Перезапустить кэш Redis
docker-compose restart redis
```

### Проблема: Высокое использование памяти
**Симптомы**: Система становится медленной, контейнеры перезапускаются

**Решения**:
```bash
# Проверить использование памяти
free -h
docker stats --no-stream

# Увеличить лимиты памяти Docker
# Отредактировать docker-compose.yml для добавления лимитов памяти

# Очистить неиспользуемые ресурсы Docker
docker system prune -f
```

## Проблемы разработки

### Проблема: Горячая перезагрузка не работает
**Симптомы**: Изменения кода не отражаются немедленно

**Решения**:
```bash
# Для Go сервисов - перезапустить конкретный сервис
docker-compose restart user-service

# Для фронтенда - проверить Vite dev server
cd frontend
npm run dev

# Проверить лимиты наблюдения за файлами (Linux)
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

### Проблема: Ошибки сборки
**Симптомы**: Ошибки Docker build или компиляции Go

**Решения**:
```bash
# Очистить кэш сборки
docker builder prune -f

# Пересобрать с нуля
docker-compose build --no-cache

# Проверить проблемы Go модулей
cd backend/user-service
go mod tidy
go mod download
```

## Проблемы тестирования

### Проблема: E2E тесты не проходят
**Симптомы**: Набор тестов сообщает о сбоях

**Решения**:
```bash
# Убедиться, что сервисы здоровы перед тестированием
./scripts/e2e-test.sh

# Проверить, что тестовая база данных чистая
docker-compose down -v
docker-compose up -d postgres redis

# Запустить тесты с подробным выводом
cd tests/e2e
go test -v -timeout 10m ./...
```

## Диагностические команды

### Проверки здоровья
```bash
# Проверить здоровье всех сервисов
curl http://localhost:8080/health
curl http://localhost:8080/v1/auth/health
curl http://localhost:8080/v1/contests/health

# Проверить подключение к базе данных
docker-compose exec postgres pg_isready -U sports_user

# Проверить подключение к Redis
docker-compose exec redis redis-cli ping
```

### Анализ логов
```bash
# Просмотреть логи всех сервисов
docker-compose logs -f

# Просмотреть логи конкретного сервиса
docker-compose logs -f user-service

# Поиск ошибок в логах
docker-compose logs | grep -i error

# Следить за логами в реальном времени
docker-compose logs -f --tail=100
```

## Получение дополнительной помощи

Если вы столкнулись с проблемами, не описанными здесь:

1. **Проверьте логи сервисов**: `docker-compose logs [имя-сервиса]`
2. **Проверьте конфигурацию**: `docker-compose config`
3. **Протестируйте подключение**: Используйте curl команды для тестирования API эндпоинтов
4. **Проверьте системные ресурсы**: Убедитесь в достаточности памяти и дискового пространства
5. **Просмотрите документацию**: Проверьте [документацию API](../api/services-overview.md)

## Экстренное восстановление

### Полный сброс системы
```bash
# Остановить все сервисы
docker-compose down -v

# Удалить все контейнеры и тома
docker system prune -a -f --volumes

# Пересобрать все
make clean
make setup
make dev
```

---

*Для дополнительных диагностических инструментов и процедур отладки см. [Диагностические инструменты](diagnostic-tools.md).*
