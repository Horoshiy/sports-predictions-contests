# Project Structure

## Directory Layout
```
sports-prediction-contests/
├── backend/                    # Go микросервисы
│   ├── api-gateway/           # API Gateway сервис
│   ├── contest-service/       # Управление конкурсами
│   ├── prediction-service/    # Обработка прогнозов
│   ├── scoring-service/       # Подсчёт очков
│   ├── user-service/          # Управление пользователями
│   ├── sports-service/        # Виды спорта и события
│   ├── notification-service/  # Уведомления и боты
│   ├── proto/                 # Protocol Buffers схемы
│   └── shared/                # Общие библиотеки
├── frontend/                  # React приложение
│   ├── src/
│   │   ├── components/        # React компоненты
│   │   ├── pages/            # Страницы приложения
│   │   ├── services/         # API клиенты (gRPC-Web)
│   │   ├── hooks/            # Custom React hooks
│   │   ├── utils/            # Утилиты и хелперы
│   │   └── types/            # TypeScript типы
│   ├── public/               # Статические файлы
│   └── build/                # Собранное приложение
├── bots/                     # Telegram и Facebook боты
│   ├── telegram/             # Telegram Bot
│   └── facebook/             # Facebook Messenger Bot
├── docs/                     # Документация
│   ├── api/                  # API документация
│   ├── deployment/           # Инструкции по развёртыванию
│   └── architecture/         # Архитектурные диаграммы
├── scripts/                  # Скрипты автоматизации
├── docker/                   # Docker конфигурации
├── k8s/                      # Kubernetes манифесты
└── .kiro/                    # Kiro CLI конфигурация
```

## File Naming Conventions
**Go Backend:**
- Файлы: `snake_case.go`
- Пакеты: `lowercase`
- Структуры: `PascalCase`
- Функции: `PascalCase` (публичные), `camelCase` (приватные)

**React Frontend:**
- Компоненты: `PascalCase.tsx`
- Хуки: `use[Name].ts`
- Утилиты: `camelCase.ts`
- Типы: `[Name].types.ts`

**Protocol Buffers:**
- Файлы: `snake_case.proto`
- Сервисы: `PascalCase`
- Сообщения: `PascalCase`

## Module Organization
**Backend сервисы:**
- Каждый сервис в отдельной директории
- Стандартная структура: `cmd/`, `internal/`, `pkg/`
- Общие proto файлы в `proto/`
- Shared библиотеки в `shared/`

**Frontend модули:**
- Компоненты группируются по функциональности
- Переиспользуемые компоненты в `components/common/`
- Страницы в `pages/` с соответствующими компонентами
- API клиенты генерируются из proto файлов

## Configuration Files
**Backend:**
- `config.yaml` - конфигурация каждого сервиса
- `docker-compose.yml` - локальная разработка
- `Dockerfile` - для каждого сервиса

**Frontend:**
- `package.json` - зависимости и скрипты
- `tsconfig.json` - TypeScript конфигурация
- `vite.config.ts` - сборка и dev server

**Общие:**
- `.env` файлы для переменных окружения
- `Makefile` - команды сборки и развёртывания

## Documentation Structure
- **API документация** - автогенерация из proto файлов
- **Архитектурные решения** - ADR (Architecture Decision Records)
- **Руководства пользователя** - для каждого типа пользователей
- **Инструкции разработчика** - setup, testing, deployment

## Asset Organization
**Frontend ассеты:**
- `public/images/` - статические изображения
- `src/assets/` - ассеты для сборки
- `src/styles/` - глобальные стили и темы
- `src/locales/` - файлы локализации

## Build Artifacts
**Backend:**
- Бинарные файлы в `bin/`
- Docker образы с тегами версий
- Proto-generated код в `gen/`

**Frontend:**
- Собранное приложение в `build/`
- Статические ассеты с хешами
- Source maps для отладки

## Environment-Specific Files
**Конфигурация окружений:**
- `config/dev.yaml` - разработка
- `config/staging.yaml` - тестовое окружение
- `config/prod.yaml` - продакшн
- `.env.example` - шаблон переменных окружения

**Kubernetes:**
- `k8s/dev/` - манифесты для разработки
- `k8s/staging/` - манифесты для staging
- `k8s/prod/` - манифесты для продакшн
