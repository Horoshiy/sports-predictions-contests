# Sports Prediction Contests Platform

🏆 **Мультиязычная, мультиспортивная платформа для создания и проведения конкурсов спортивных прогнозов** | **Multilingual, multi-sport platform for creating and running sports prediction competitions**

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18.2-61DAFB?logo=react)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.2-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> **🏆 Built for Dynamous Kiro Hackathon** | Comprehensive microservices architecture with 10 services, 7 innovative features, and full E2E testing

---

## 📖 Содержание | Table of Contents

- [🌟 О проекте | About](#-о-проекте--about)
- [✨ Реализованные возможности | Implemented Features](#-реализованные-возможности--implemented-features)
- [🖼️ Скриншоты | Screenshots](#️-скриншоты--screenshots)
- [🏗️ Архитектура | Architecture](#️-архитектура--architecture)
- [🚀 Быстрый старт | Quick Start](#-быстрый-старт--quick-start)
- [📚 Документация | Documentation](#-документация--documentation)
- [🛠️ Технологии | Tech Stack](#️-технологии--tech-stack)
- [💡 Инновации | Innovations](#-инновации--innovations)
- [🧪 Тестирование | Testing](#-тестирование--testing)
- [📦 Развертывание | Deployment](#-развертывание--deployment)

---

## 🌟 О проекте | About

**Sports Prediction Contests** — это полнофункциональная платформа-конструктор для создания и проведения конкурсов спортивных прогнозов. Платформа превращает конкурсы прогнозов из узкоспециализированного продукта в универсальный движок вовлечения спортивной аудитории.

**Sports Prediction Contests** is a full-featured platform constructor for creating and running sports prediction competitions. The platform transforms prediction contests from a niche product into a universal engagement engine for sports communities.

### 🎯 Ключевые преимущества | Key Advantages

- **🚀 Быстрый запуск** | **Quick Launch**: Создание конкурса за минуты, а не дни | Create contests in minutes, not days
- **🌍 Мультиязычность** | **Multilingual**: Полная поддержка русского и английского языков | Full Russian and English support
- **⚽ Мультиспорт** | **Multi-Sport**: Футбол, баскетбол, хоккей и другие виды спорта | Football, basketball, hockey, and more
- **🎮 Геймификация** | **Gamification**: Серии прогнозов, достижения, турниры команд | Prediction streaks, achievements, team tournaments
- **📱 Мультиплатформенность** | **Multi-Platform**: Web, Telegram бот, API для интеграций | Web, Telegram bot, API for integrations
- **📊 Аналитика** | **Analytics**: Детальная статистика точности прогнозов | Detailed prediction accuracy statistics
- **🏆 Инновации** | **Innovations**: 7 уникальных функций геймификации | 7 unique gamification features

---

## ✨ Реализованные возможности | Implemented Features

### 🎯 Основной функционал | Core Functionality

#### 👤 Управление пользователями | User Management
- ✅ Регистрация и аутентификация с JWT токенами | Registration and JWT authentication
- ✅ Профили пользователей с настройками | User profiles with preferences
- ✅ Управление предпочтениями и уведомлениями | Notification and preference management
- ✅ Статистика и история прогнозов | Prediction statistics and history

#### 🏆 Система конкурсов | Contest System
- ✅ Конструктор конкурсов с гибкими правилами | Contest constructor with flexible rules
- ✅ Поддержка различных типов прогнозов | Multiple prediction types support
- ✅ Управление участниками | Participant management
- ✅ Настраиваемые системы подсчета очков | Customizable scoring systems
- ✅ Публичные и приватные конкурсы | Public and private contests

#### 🎲 Прогнозы | Predictions
- ✅ Прогнозы исходов матчей | Match outcome predictions
- ✅ Прогнозы точного счета | Exact score predictions
- ✅ Прогнозы статистики (Props) | Statistical predictions (Props)
- ✅ Валидация и ограничения по времени | Validation and time restrictions
- ✅ История и редактирование прогнозов | Prediction history and editing

#### 📊 Подсчет очков и рейтинги | Scoring & Leaderboards
- ✅ Автоматический подсчет очков | Automatic score calculation
- ✅ Динамические коэффициенты по времени | Dynamic time-based coefficients
- ✅ Таблицы лидеров с кэшированием | Cached leaderboards
- ✅ Серии правильных прогнозов с множителями | Prediction streaks with multipliers
- ✅ Детальная статистика по пользователям | Detailed user statistics

#### ⚽ Управление спортивными данными | Sports Data Management
- ✅ Управление видами спорта | Sports management
- ✅ Лиги и турниры | Leagues and tournaments
- ✅ Команды и составы | Teams and rosters
- ✅ Матчи и результаты | Matches and results
- ✅ Интеграция с внешними API | External API integration

#### 🔔 Уведомления | Notifications
- ✅ Email уведомления | Email notifications
- ✅ Telegram бот интеграция | Telegram bot integration
- ✅ Уведомления о начале матчей | Match start notifications
- ✅ Уведомления о результатах | Results notifications
- ✅ Система очередей для массовых рассылок | Queue system for bulk notifications

### 🚀 Инновационные функции | Innovative Features

#### 1. 🔥 Серии прогнозов с множителями | Prediction Streaks with Multipliers
Серия правильных прогнозов увеличивает множитель очков, но сбрасывается при ошибке.

A series of correct predictions increases the point multiplier, but resets on failure.

- Отслеживание текущей и максимальной серии | Current and max streak tracking
- Прогрессивные множители (1.1x, 1.2x, 1.5x, 2.0x) | Progressive multipliers
- Визуализация серий в профиле | Streak visualization in profile

#### 2. ⏰ Динамические коэффициенты | Dynamic Point Coefficients
Очки за прогнозы меняются в зависимости от времени подачи — ранние прогнозы дают больше очков.

Points for predictions change based on submission time — earlier predictions earn more points.

- Формула затухания по времени | Time-decay formula
- Максимум очков за ранние прогнозы | Maximum points for early predictions
- Отображение потенциальных очков | Potential points display

#### 3. ⚔️ Дуэли один-на-один | Head-to-Head Challenges
Прямые поединки между двумя пользователями на конкретный матч или серию матчей.

Direct duels between two users on a specific match or series of matches.

- Система вызовов и принятия | Challenge invitation system
- Отдельный подсчет очков для дуэлей | Dedicated H2H scoring
- Интеграция с Telegram ботом | Telegram bot integration
- История дуэлей | Challenge history

#### 4. 👥 Командные турниры | Team Tournaments
Создание команд из нескольких участников с общим рейтингом.

Create teams of multiple participants with shared ranking.

- Создание и управление командами | Team creation and management
- Система приглашений | Invitation system
- Командные таблицы лидеров | Team leaderboards
- Роли в команде (капитан, участник) | Team roles (captain, member)

#### 5. 📈 Дашборд аналитики | Analytics Dashboard
Детальная статистика прогнозов: точность по лигам, командам, типам ставок, тренды.

Detailed prediction statistics: accuracy by league, team, bet type, trends over time.

- Графики точности по времени | Accuracy trends over time
- Анализ по лигам и командам | League and team analysis
- Сравнение со средним по платформе | Platform average comparison
- Экспорт данных | Data export functionality

#### 6. 📊 Прогнозы статистики (Props) | Props Predictions
Прогнозы не только на исход, но и на статистику: голы игроков, угловые, владение мячом.

Predictions not just on outcome, but on statistics: player goals, corners, possession.

- Расширенные типы событий | Extended event types
- Специфичные правила подсчета | Props-specific scoring rules
- Интеграция с детальной статистикой | Detailed stats integration

#### 7. 🤖 Telegram бот | Telegram Bot
Полнофункциональный бот для участия в конкурсах через Telegram.

Full-featured bot for participating in contests via Telegram.

- Просмотр конкурсов и прогнозов | View contests and predictions
- Подача прогнозов | Submit predictions
- Проверка рейтингов | Check leaderboards
- Уведомления о матчах | Match notifications

---

## 🖼️ Скриншоты | Screenshots

> **Примечание**: Скриншоты будут добавлены после запуска приложения. Для генерации скриншотов используйте команду `make playwright-test`.
>
> **Note**: Screenshots will be added after running the application. Use `make playwright-test` to generate screenshots.

### Основные экраны | Main Screens

#### 🔐 Аутентификация | Authentication
![Login Page](docs/screenshots/login-page.png)
*Страница входа с JWT аутентификацией | Login page with JWT authentication*

![Register Page](docs/screenshots/register-page.png)
*Регистрация нового пользователя | New user registration*

#### 🏆 Конкурсы | Contests
![Contests List](docs/screenshots/contests-list.png)
*Список доступных конкурсов | Available contests list*

![Contest Details](docs/screenshots/contest-details.png)
*Детали конкурса с правилами и участниками | Contest details with rules and participants*

#### 🎲 Прогнозы | Predictions
![Predictions Page](docs/screenshots/predictions-page.png)
*Интерфейс подачи прогнозов | Prediction submission interface*

![Leaderboard](docs/screenshots/leaderboard.png)
*Таблица лидеров с очками и рейтингами | Leaderboard with scores and rankings*

#### 👤 Профиль и аналитика | Profile & Analytics
![Profile Page](docs/screenshots/profile-page.png)
*Профиль пользователя со статистикой | User profile with statistics*

![Analytics Dashboard](docs/screenshots/analytics-dashboard.png)
*Дашборд аналитики с графиками | Analytics dashboard with charts*

#### ⚽ Управление данными | Data Management
![Sports Management](docs/screenshots/sports-management.png)
*Управление спортивными данными | Sports data management*

![Teams Page](docs/screenshots/teams-page.png)
*Командные турниры | Team tournaments*

---
## 🏗️ Архитектура | Architecture

### Микросервисная архитектура | Microservices Architecture

Платформа построена на основе микросервисной архитектуры с 10 независимыми сервисами, взаимодействующими через gRPC.

The platform is built on a microservices architecture with 10 independent services communicating via gRPC.

```
backend/
├── api-gateway/           # API Gateway (порт 8080) | API Gateway (port 8080)
├── user-service/          # Пользователи и аутентификация (8084) | Users & auth (8084)
├── contest-service/       # Управление конкурсами (8085) | Contest management (8085)
├── prediction-service/    # Прогнозы пользователей (8086) | User predictions (8086)
├── scoring-service/       # Подсчет очков (8087) | Scoring calculation (8087)
├── sports-service/        # Спортивные данные (8088) | Sports data (8088)
├── notification-service/  # Уведомления (8089) | Notifications (8089)
├── challenge-service/     # Дуэли 1-на-1 (8090) | H2H challenges (8090)
├── proto/                 # gRPC схемы | gRPC schemas
└── shared/                # Общие библиотеки | Shared libraries
```

### Компоненты системы | System Components

#### Backend Services
- **API Gateway**: HTTP/REST точка входа, маршрутизация к микросервисам | HTTP/REST entry point, routing to microservices
- **User Service**: JWT аутентификация, профили, настройки | JWT authentication, profiles, preferences
- **Contest Service**: CRUD конкурсов, правила, участники | Contest CRUD, rules, participants
- **Prediction Service**: Подача и валидация прогнозов | Prediction submission and validation
- **Scoring Service**: Алгоритмы подсчета, кэширование рейтингов | Scoring algorithms, leaderboard caching
- **Sports Service**: Виды спорта, лиги, команды, матчи | Sports, leagues, teams, matches
- **Notification Service**: Email, Telegram, очереди уведомлений | Email, Telegram, notification queues
- **Challenge Service**: Дуэли между пользователями | User-to-user challenges

#### Frontend Application
```
frontend/
├── src/
│   ├── pages/            # 8 основных страниц | 8 main pages
│   ├── components/       # Переиспользуемые компоненты | Reusable components
│   ├── services/         # gRPC-Web клиенты | gRPC-Web clients
│   ├── hooks/            # Custom React hooks
│   ├── contexts/         # Управление состоянием | State management
│   └── types/            # TypeScript определения | TypeScript definitions
└── tests/
    └── e2e/              # Playwright E2E тесты | Playwright E2E tests
```

#### Bot Integration
```
bots/
├── telegram/             # Telegram бот с gRPC клиентами | Telegram bot with gRPC clients
└── facebook/             # Facebook Messenger бот (запланирован) | Facebook bot (planned)
```

### Хранение данных | Data Storage

- **PostgreSQL 15**: Основная база данных для всех сервисов | Primary database for all services
- **Redis 7**: Кэширование рейтингов, сессии пользователей | Leaderboard caching, user sessions
- **Docker Volumes**: Персистентное хранение данных | Persistent data storage

---

## 🚀 Быстрый старт | Quick Start

### Предварительные требования | Prerequisites

Перед установкой убедитесь, что у вас установлено:

Before installation, ensure you have installed:

- **Go 1.21+** - [Руководство по установке](https://golang.org/doc/install) | [Installation Guide](https://golang.org/doc/install)
- **Node.js 18+** - [Руководство по установке](https://nodejs.org/en/download/) | [Installation Guide](https://nodejs.org/en/download/)
- **Docker & Docker Compose** - [Руководство по установке](https://docs.docker.com/get-docker/) | [Installation Guide](https://docs.docker.com/get-docker/)
- **Protocol Buffers Compiler** - [Руководство по установке](https://grpc.io/docs/protoc-installation/) | [Installation Guide](https://grpc.io/docs/protoc-installation/)

### Шаг 1: Клонирование и настройка | Step 1: Clone and Setup

```bash
# Клонировать репозиторий | Clone repository
git clone https://github.com/yourusername/sports-prediction-contests
cd sports-prediction-contests

# Запустить автоматическую настройку | Run automatic setup
make setup
```

Команда `make setup` выполнит:
- Создание файла `.env` из `.env.example`
- Установку зависимостей Go и Node.js
- Проверку наличия необходимых инструментов

The `make setup` command will:
- Create `.env` file from `.env.example`
- Install Go and Node.js dependencies
- Check for required tools

### Шаг 2: Запуск окружения разработки | Step 2: Start Development Environment

```bash
# Запустить PostgreSQL и Redis | Start PostgreSQL and Redis
make dev
```

Это запустит:
- PostgreSQL базу данных (localhost:5432)
- Redis кэш (localhost:6379)
- Инициализацию схемы базы данных

This will start:
- PostgreSQL database (localhost:5432)
- Redis cache (localhost:6379)
- Database schema initialization

### Шаг 3: Запуск всех сервисов | Step 3: Start All Services

```bash
# Запустить все микросервисы и фронтенд | Start all microservices and frontend
make docker-services
```

Сервисы будут доступны по адресам:
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Все микросервисы**: порты 8084-8090

Services will be available at:
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **All microservices**: ports 8084-8090

### Шаг 4: Заполнение тестовыми данными (опционально) | Step 4: Populate with Test Data (Optional)

```bash
# Небольшой набор данных (20 пользователей, 8 конкурсов) | Small dataset (20 users, 8 contests)
make seed-small

# Средний набор данных (100 пользователей, 25 конкурсов) | Medium dataset (100 users, 25 contests)
make seed-medium

# Большой набор данных (500 пользователей, 50 конкурсов) | Large dataset (500 users, 50 contests)
make seed-large
```

Тестовые данные включают:
- Реалистичные пользователи с профилями
- Конкурсы с различными правилами
- Спортивные данные (команды, лиги, матчи)
- Прогнозы и результаты
- Командные турниры

Test data includes:
- Realistic users with profiles
- Contests with various rules
- Sports data (teams, leagues, matches)
- Predictions and results
- Team tournaments

### Шаг 5: Проверка работы | Step 5: Verify Installation

```bash
# Проверить статус всех сервисов | Check status of all services
make status

# Просмотреть логи | View logs
make logs
```

### Первый вход | First Login

После запуска откройте браузер и перейдите на http://localhost:3000

After starting, open your browser and navigate to http://localhost:3000

**Тестовые учетные данные** | **Test Credentials**:
- Email: `user1@example.com`
- Password: `password123`

---

## 📚 Документация | Documentation

### Полная документация | Complete Documentation

📚 **Доступна комплексная двуязычная документация на русском и английском языках:**

📚 **Comprehensive bilingual documentation is available in English and Russian:**

#### Русская документация | Russian Documentation
- [📖 Полная документация](docs/ru/README.md) - Главная страница документации
- [🚀 Быстрый старт](docs/ru/deployment/quick-start.md) - Запуск за несколько минут
- [📋 Справочник API](docs/ru/api/services-overview.md) - Полная документация API
- [🧪 Руководство по тестированию](docs/ru/testing/e2e-testing.md) - Процедуры тестирования
- [🔧 Устранение неполадок](docs/ru/troubleshooting/common-issues.md) - Частые проблемы и решения

#### English Documentation
- [📖 Complete Documentation](docs/en/README.md) - Documentation home page
- [🚀 Quick Start Guide](docs/en/deployment/quick-start.md) - Get running in minutes
- [📋 API Reference](docs/en/api/services-overview.md) - Complete API documentation
- [🧪 Testing Guide](docs/en/testing/e2e-testing.md) - Testing procedures
- [🔧 Troubleshooting](docs/en/troubleshooting/common-issues.md) - Common issues and solutions

#### Архитектура | Architecture
- [🏗️ Диаграммы архитектуры](docs/assets/architecture-diagram.md) - Визуальный обзор системы | Visual system overview

---

## 🛠️ Технологии | Tech Stack

### Backend
- **Язык | Language**: Go 1.24
- **Фреймворк | Framework**: gRPC с Protocol Buffers v3 | gRPC with Protocol Buffers v3
- **База данных | Database**: PostgreSQL 15 (GORM ORM)
- **Кэш | Cache**: Redis 7
- **Аутентификация | Authentication**: JWT (golang-jwt/jwt/v5)
- **Тестирование | Testing**: Go testing framework
- **Сборка | Build**: Go workspaces, Docker multi-stage builds

### Frontend
- **Фреймворк | Framework**: React 18.2 с TypeScript 5.2 | React 18.2 with TypeScript 5.2
- **Инструмент сборки | Build Tool**: Vite 5.0
- **UI Библиотека | UI Library**: Ant Design 5.22
- **Управление состоянием | State Management**: TanStack Query (React Query) v5
- **Формы | Forms**: React Hook Form с Zod валидацией | React Hook Form with Zod validation
- **Маршрутизация | Routing**: React Router v6
- **Графики | Charts**: Recharts 2.8
- **API Клиент | API Client**: gRPC-Web 1.4.2
- **Тестирование | Testing**: Playwright 1.48 (E2E), Vitest (unit)

### Инфраструктура | Infrastructure
- **Контейнеризация | Containerization**: Docker с Docker Compose | Docker with Docker Compose
- **Оркестрация | Orchestration**: Docker Compose (dev), готово к Kubernetes | Docker Compose (dev), Kubernetes-ready
- **CI/CD**: Скрипты для автоматизированного тестирования | Scripts for automated testing
- **Мониторинг | Monitoring**: Структурированное логирование, health checks | Structured logging, health checks

---

## 💡 Инновации | Innovations

Платформа включает 7 уникальных инновационных функций, выходящих за рамки базового функционала конкурсов прогнозов:

The platform includes 7 unique innovative features beyond basic prediction contest functionality:

### Реализованные инновации | Implemented Innovations

| # | Функция | Feature | Сложность | Complexity | Время | Time |
|---|---------|---------|-----------|------------|-------|------|
| 1 | 🔥 Серии прогнозов | Prediction Streaks | Низкая | Low | 2-4ч | 2-4h |
| 2 | ⏰ Динамические коэффициенты | Dynamic Coefficients | Низкая | Low | 2-4ч | 2-4h |
| 3 | ⚔️ Дуэли 1-на-1 | H2H Challenges | Низкая | Low | 2-4ч | 2-4h |
| 4 | 👥 Командные турниры | Team Tournaments | Средняя | Medium | 4-8ч | 4-8h |
| 5 | 📈 Дашборд аналитики | Analytics Dashboard | Средняя | Medium | 4-8ч | 4-8h |
| 6 | 📊 Props прогнозы | Props Predictions | Средняя | Medium | 4-8ч | 4-8h |
| 7 | 🤖 Telegram бот | Telegram Bot | Средняя | Medium | 4-8ч | 4-8h |

### Будущие инновации | Future Innovations

Запланированные функции для будущих версий:

Planned features for future versions:

- **🎯 Мультиспортивные комбо** | **Multi-Sport Combos**: Прогнозы на несколько видов спорта | Predictions across multiple sports
- **👥 Социальные прогнозы** | **Social Predictions**: Копирование прогнозов экспертов | Copy expert predictions
- **🎮 Сезонный пропуск** | **Season Pass**: Battle Pass система наград | Battle Pass reward system
- **🤖 AI ассистент** | **AI Assistant**: LLM для анализа и рекомендаций | LLM for analysis and recommendations
- **⚡ Live прогнозы** | **Live Predictions**: Прогнозы во время матчей | In-match predictions
- **📺 Интеграция стриминга** | **Streaming Integration**: Виджеты для Twitch/YouTube | Widgets for Twitch/YouTube

---

## 🧪 Тестирование | Testing

### Комплексное тестирование | Comprehensive Testing

Платформа включает полный набор тестов для обеспечения качества и надежности:

The platform includes a complete test suite to ensure quality and reliability:

#### E2E тестирование с Playwright | E2E Testing with Playwright

```bash
# Запустить все E2E тесты | Run all E2E tests
make playwright-test

# Запустить в UI режиме | Run in UI mode
make playwright-test-ui

# Запустить в headed режиме | Run in headed mode
make playwright-test-headed

# Показать отчет | Show report
make playwright-report
```

**Покрытие тестами** | **Test Coverage**:
- ✅ Аутентификация (вход, регистрация, выход) | Authentication (login, register, logout)
- ✅ Управление конкурсами | Contest management
- ✅ Подача прогнозов | Prediction submission
- ✅ Просмотр рейтингов | Leaderboard viewing
- ✅ Профиль пользователя | User profile
- ✅ Дашборд аналитики | Analytics dashboard
- ✅ Управление спортивными данными | Sports data management
- ✅ Командные турниры | Team tournaments
- ✅ Навигация и маршрутизация | Navigation and routing
- ✅ Полные пользовательские сценарии | Complete user workflows

#### Кросс-браузерное тестирование | Cross-Browser Testing

Тесты выполняются на трех браузерах:
- ✅ Chromium (Chrome, Edge)
- ✅ Firefox
- ✅ WebKit (Safari)

Tests run on three browsers:
- ✅ Chromium (Chrome, Edge)
- ✅ Firefox
- ✅ WebKit (Safari)

#### Unit тестирование | Unit Testing

```bash
# Запустить unit тесты backend | Run backend unit tests
cd backend && go test ./...

# Запустить unit тесты frontend | Run frontend unit tests
cd frontend && npm test
```

### Автоматизированная проверка | Automated Validation

```bash
# Проверить все сервисы | Check all services
make check-services

# Запустить полный набор тестов | Run full test suite
make test
```

---

## 📦 Развертывание | Deployment

### Разработка | Development

```bash
# Запустить окружение разработки | Start development environment
make dev

# Запустить все сервисы | Start all services
make docker-services

# Просмотреть логи | View logs
make logs

# Остановить все сервисы | Stop all services
make docker-down
```

### Production

Платформа готова к развертыванию в production с использованием:

The platform is ready for production deployment using:

- **Docker Compose**: Для простого развертывания | For simple deployment
- **Kubernetes**: Манифесты готовы в `k8s/` | Manifests ready in `k8s/`
- **Cloud Providers**: AWS, GCP, Azure совместимы | AWS, GCP, Azure compatible

#### Переменные окружения | Environment Variables

Скопируйте `.env.example` в `.env` и настройте:

Copy `.env.example` to `.env` and configure:

```bash
# База данных | Database
DB_PASSWORD=your_secure_password
DB_SSLMODE=require  # Для production | For production

# JWT
JWT_SECRET=your_secure_random_string

# Telegram (опционально) | Telegram (optional)
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_ENABLED=true

# Email (опционально) | Email (optional)
SMTP_HOST=smtp.example.com
SMTP_USER=your_email
SMTP_PASSWORD=your_password
EMAIL_ENABLED=true
```

#### Безопасность | Security

⚠️ **Важно для production** | **Important for production**:
- Используйте сильные пароли | Use strong passwords
- Включите SSL/TLS для базы данных | Enable SSL/TLS for database
- Настройте CORS правильно | Configure CORS properly
- Используйте secrets management | Use secrets management
- Включите rate limiting | Enable rate limiting

---

## 🎯 Использование | Usage

### Для организаторов конкурсов | For Contest Organizers

1. **Создание конкурса** | **Create Contest**
   - Войдите в систему | Log in to the system
   - Перейдите на страницу "Конкурсы" | Navigate to "Contests" page
   - Нажмите "Создать конкурс" | Click "Create Contest"
   - Настройте правила и систему очков | Configure rules and scoring system

2. **Добавление матчей** | **Add Matches**
   - Перейдите в "Управление спортом" | Go to "Sports Management"
   - Создайте виды спорта, лиги, команды | Create sports, leagues, teams
   - Добавьте матчи с датами | Add matches with dates

3. **Управление участниками** | **Manage Participants**
   - Просматривайте список участников | View participant list
   - Отслеживайте активность | Monitor activity
   - Управляйте доступом | Manage access

### Для участников | For Participants

1. **Регистрация** | **Registration**
   - Создайте аккаунт на странице регистрации | Create account on registration page
   - Подтвердите email (если настроено) | Confirm email (if configured)

2. **Присоединение к конкурсу** | **Join Contest**
   - Просмотрите доступные конкурсы | Browse available contests
   - Присоединитесь к интересующему конкурсу | Join contest of interest

3. **Подача прогнозов** | **Submit Predictions**
   - Перейдите на страницу "Прогнозы" | Go to "Predictions" page
   - Выберите матч | Select match
   - Сделайте прогноз до начала матча | Make prediction before match starts

4. **Отслеживание результатов** | **Track Results**
   - Просматривайте таблицу лидеров | View leaderboard
   - Проверяйте свою статистику в профиле | Check your statistics in profile
   - Анализируйте точность в дашборде аналитики | Analyze accuracy in analytics dashboard

### Через Telegram бота | Via Telegram Bot

1. Найдите бота в Telegram | Find bot in Telegram
2. Отправьте `/start` для начала | Send `/start` to begin
3. Используйте команды для взаимодействия | Use commands to interact:
   - `/contests` - Список конкурсов | Contest list
   - `/predict` - Сделать прогноз | Make prediction
   - `/leaderboard` - Таблица лидеров | Leaderboard
   - `/profile` - Ваш профиль | Your profile

---

## 🤝 Вклад в проект | Contributing

Мы приветствуем вклад в развитие платформы!

We welcome contributions to the platform!

### Как внести вклад | How to Contribute

1. Fork репозиторий | Fork the repository
2. Создайте feature branch | Create a feature branch
3. Внесите изменения | Make your changes
4. Напишите тесты | Write tests
5. Отправьте pull request | Submit a pull request

### Разработка с Kiro CLI

Проект оптимизирован для разработки с Kiro CLI:

The project is optimized for development with Kiro CLI:

- **`@prime`** - Загрузить контекст проекта | Load project context
- **`@plan-feature`** - Спланировать новую функцию | Plan new feature
- **`@execute`** - Реализовать план | Implement plan
- **`@code-review`** - Проверить качество кода | Review code quality

---

## 📄 Лицензия | License

MIT License - см. файл [LICENSE](LICENSE) для деталей

MIT License - see [LICENSE](LICENSE) file for details

---

## 🙏 Благодарности | Acknowledgments

- **Dynamous Kiro Hackathon** - За возможность создать этот проект | For the opportunity to create this project
- **Go Community** - За отличные инструменты и библиотеки | For excellent tools and libraries
- **React Community** - За мощный фронтенд фреймворк | For powerful frontend framework
- **Open Source Contributors** - За все используемые библиотеки | For all the libraries used

---

## 📞 Контакты | Contact

- **GitHub**: [github.com/yourusername/sports-prediction-contests](https://github.com/yourusername/sports-prediction-contests)
- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/yourusername/sports-prediction-contests/issues)

---

**Готовы начать?** | **Ready to start?** 🚀

```bash
git clone https://github.com/yourusername/sports-prediction-contests
cd sports-prediction-contests
make setup
make dev
make docker-services
```

Откройте http://localhost:3000 и начните создавать конкурсы!

Open http://localhost:3000 and start creating contests!
