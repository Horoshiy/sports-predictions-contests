# Screenshots

This directory contains screenshots of the Sports Prediction Contests platform.

## Автоматическая генерация | Automatic Generation

Для автоматической генерации всех скриншотов:

To automatically generate all screenshots:

```bash
# 1. Запустите все сервисы | Start all services
make dev
make docker-services

# 2. Заполните тестовыми данными | Populate with test data
make seed-small

# 3. Сгенерируйте скриншоты | Generate screenshots
make generate-screenshots
```

Скриншоты будут сохранены в этой директории.

Screenshots will be saved in this directory.

## Ручная генерация | Manual Generation

Вы также можете сделать скриншоты вручную:

You can also take screenshots manually:

1. Запустите приложение | Start the application
2. Откройте браузер на http://localhost:3000
3. Войдите с тестовыми данными | Login with test credentials:
   - Email: `user1@example.com`
   - Password: `password123`
4. Сделайте скриншоты каждой страницы | Take screenshots of each page

## Список скриншотов | Screenshot List

- `login-page.png` - Страница входа | Login screen
- `register-page.png` - Страница регистрации | Registration screen
- `contests-list.png` - Список конкурсов | Contests listing page
- `contest-details.png` - Детали конкурса | Contest details view
- `predictions-page.png` - Интерфейс прогнозов | Predictions interface
- `leaderboard.png` - Таблица лидеров | Leaderboard view
- `profile-page.png` - Профиль пользователя | User profile
- `analytics-dashboard.png` - Дашборд аналитики | Analytics dashboard
- `sports-management.png` - Управление спортом | Sports management interface
- `teams-page.png` - Командные турниры | Teams and tournaments

## Требования к скриншотам | Screenshot Guidelines

- Разрешение | Resolution: 1920x1080
- Формат | Format: PNG
- Полная страница | Full page capture
- Реалистичные данные | Realistic test data
- Чистый интерфейс | Clean interface (no debug tools visible)

