# Тестирование фронтенда с Playwright

## Обзор

Платформа Sports Prediction Contests использует Playwright для комплексного end-to-end тестирования React фронтенда. Playwright обеспечивает надежное, быстрое и кроссбраузерное тестирование с возможностью AI-генерации тестов через интеграцию MCP.

## Установка

### Требования
- Node.js 18+
- Docker и Docker Compose (для интеграции с сервисами)

### Установка Playwright

```bash
# Установить Playwright и браузеры
make playwright-install

# Или вручную
cd frontend
npm install
npx playwright install --with-deps
```

## Запуск тестов

### Полный набор тестов с сервисами

```bash
# Запустить все тесты с автоматической оркестрацией сервисов
make playwright-test

# Запустить конкретный тестовый файл
./scripts/playwright-test.sh auth.spec.ts
```

### Быстрое тестирование (сервисы уже запущены)

```bash
# Сначала запустить сервисы
make docker-services

# Запустить только тесты
make playwright-test-only

# Или конкретные тесты
cd frontend
npm run test:e2e -- auth.spec.ts
```

### Интерактивные режимы

```bash
# UI режим - интерактивный запуск тестов
make playwright-test-ui

# Headed режим - видеть браузер во время тестов
make playwright-test-headed

# Debug режим - пошаговое выполнение тестов
cd frontend
npm run test:e2e:debug auth.spec.ts
```

## Структура тестов

### Организация тестов

```
frontend/tests/
├── e2e/              # End-to-end тестовые наборы
│   ├── auth.spec.ts
│   ├── contests.spec.ts
│   ├── predictions.spec.ts
│   ├── teams.spec.ts
│   ├── analytics.spec.ts
│   ├── profile.spec.ts
│   ├── navigation.spec.ts
│   └── workflows.spec.ts
├── visual/           # Визуальные регрессионные тесты
│   └── snapshots.spec.ts
├── fixtures/         # Тестовые фикстуры
│   ├── auth.fixture.ts
│   ├── api.fixture.ts
│   └── data.fixture.ts
└── helpers/          # Тестовые утилиты
    ├── test-utils.ts
    ├── selectors.ts
    └── assertions.ts
```

## Написание тестов

### Базовый пример теста

```typescript
import { test, expect } from '@playwright/test'
import { SELECTORS } from '../helpers/selectors'

test.describe('Название функции', () => {
  test('должен делать что-то', async ({ page }) => {
    await page.goto('/page')
    await page.click(SELECTORS.someButton)
    await expect(page.locator('.result')).toBeVisible()
  })
})
```

### Использование фикстуры аутентификации

```typescript
import { test } from '../fixtures/auth.fixture'
import { expect } from '@playwright/test'

test('аутентифицированный пользователь может получить доступ к защищенной странице', async ({ authenticatedPage }) => {
  await authenticatedPage.goto('/contests')
  await expect(authenticatedPage).toHaveURL('/contests')
})
```

## Интеграция MCP

Сервер Playwright MCP настроен в `.kiro/settings/mcp.json` и обеспечивает AI-тестирование через Kiro CLI.

### Использование MCP для генерации тестов

В Kiro CLI вы можете использовать сервер Playwright MCP для:
- Генерации тестового кода из описаний на естественном языке
- Отладки падающих тестов с помощью AI
- Захвата скриншотов и анализа проблем UI
- Автоматизации поддержки тестов

Пример рабочего процесса:
```
# В Kiro CLI
> Используй Playwright MCP для создания теста потока входа пользователя
> Отладь падающий тест создания конкурса
> Захвати скриншот панели аналитики
```

## Визуальное регрессионное тестирование

### Создание базовых скриншотов

```bash
cd frontend
npm run test:e2e -- visual/snapshots.spec.ts --update-snapshots
```

### Запуск визуальных тестов

```bash
npm run test:e2e -- visual/snapshots.spec.ts
```

## Устранение неполадок

### Тесты падают с ошибкой "Connection Refused"
- Убедитесь, что сервисы запущены: `make docker-services`
- Проверьте здоровье API Gateway: `curl http://localhost:8080/health`

### Проблемы с установкой браузеров
- Запустите: `npx playwright install --with-deps`
- На Linux могут потребоваться системные зависимости

### Нестабильные тесты
- Увеличьте таймаут в playwright.config.ts
- Добавьте явные ожидания: `await page.waitForLoadState('networkidle')`
- Используйте `test.retry(2)` для конкретных тестов

## Интеграция CI/CD

Тесты Playwright можно интегрировать в GitHub Actions или другие CI/CD пайплайны:

```yaml
- name: Запуск тестов Playwright
  run: |
    make playwright-install
    make playwright-test

- name: Загрузка отчета о тестах
  if: always()
  uses: actions/upload-artifact@v3
  with:
    name: playwright-report
    path: frontend/playwright-report/
```

## Отчеты о тестах

```bash
# Просмотр HTML отчета
npm run test:e2e:report

# Просмотр трейса для упавших тестов
npx playwright show-trace playwright-report/trace.zip
```

## Доступные тестовые наборы

- **auth.spec.ts** - Потоки аутентификации (вход, регистрация, выход)
- **contests.spec.ts** - Управление конкурсами (создание, просмотр, участие)
- **predictions.spec.ts** - Отправка и просмотр прогнозов
- **teams.spec.ts** - Рабочие процессы командных турниров
- **analytics.spec.ts** - Взаимодействие с панелью аналитики
- **profile.spec.ts** - Управление профилем пользователя
- **navigation.spec.ts** - Навигация и маршрутизация
- **workflows.spec.ts** - Полные пользовательские рабочие процессы
- **visual/snapshots.spec.ts** - Визуальные регрессионные тесты
