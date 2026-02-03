# Feature: Тотализатор (Totalizator Contest Type)

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Новый тип конкурса **"Тотализатор"** — админ выбирает 15 матчей из разных чемпионатов/лиг для прогнозирования. Правила начисления очков едины для всех матчей тура. Участник получает сумму очков за все свои прогнозы.

**Ключевые отличия от Standard:**
- Матчи выбираются вручную из любых источников (не привязаны к одному sport_type)
- Фиксированное количество матчей в туре (по умолчанию 15, настраивается)
- Один матч (Event) может участвовать в нескольких конкурсах одновременно

## User Story

As an **admin**
I want to create a "Totalizator" contest with 15 hand-picked matches from different leagues
So that users can predict results across multiple championships in one contest

As a **user**
I want to make predictions for all matches in the Totalizator contest
So that I can compete for the highest total score

## Problem Statement

Сейчас конкурсы привязаны к одному виду спорта (sport_type), и матчи фильтруются по нему. Для Тотализатора нужна возможность выбрать матчи из разных чемпионатов вручную.

## Solution Statement

1. Добавить новый тип конкурса `totalizator` в `scoring/rules.go`
2. Создать `TotalizatorRules` с настройкой количества матчей
3. Обновить UI создания конкурса для выбора матчей из всех доступных
4. Использовать существующую таблицу `contest_events` для связи матч ↔ конкурс
5. Валидировать количество матчей при создании/активации конкурса

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium
**Primary Systems Affected**: contest-service, prediction-service, frontend, telegram-bot
**Dependencies**: Существующая таблица `contest_events`, `scoring` package

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `backend/shared/scoring/rules.go` (lines 1-150) - Why: Здесь определены типы конкурсов и правила подсчёта очков
- `backend/shared/scoring/calculator.go` - Why: Логика подсчёта очков, нужно добавить поддержку totalizator
- `backend/contest-service/internal/models/contest.go` - Why: Модель Contest с полем Rules (JSON)
- `backend/prediction-service/internal/repository/event_repository.go` (lines 80-110) - Why: Метод ListByContest использует contest_events
- `backend/prediction-service/internal/models/event.go` - Why: Модель Event
- `frontend/src/components/contests/ScoringRulesEditor.tsx` - Why: UI редактор правил, нужно добавить totalizator
- `frontend/src/components/contests/ContestForm.tsx` - Why: Форма создания конкурса
- `bots/telegram/bot/contests.go` - Why: Обработка конкурсов в боте

### New Files to Create

- `frontend/src/components/contests/EventSelector.tsx` - Компонент выбора матчей для Тотализатора
- (Опционально) `frontend/src/components/contests/TotalizatorRulesEditor.tsx` - Если нужен отдельный редактор

### Relevant Documentation

- Существующий план: `.agents/plans/contest-scoring-rules.md` - паттерны реализации правил
- Proto: `backend/proto/prediction.proto` - API методы для events

### Patterns to Follow

**Contest Rules JSON Structure** (из rules.go):
```go
type ContestRules struct {
    Type     ContestType           `json:"type"`
    Standard *StandardScoringRules `json:"scoring,omitempty"`
    Risky    *RiskyScoringRules    `json:"risky,omitempty"`
    // Добавить: Totalizator *TotalizatorRules
}
```

**Связь Event ↔ Contest** (через contest_events):
```sql
-- Уже существует:
CREATE TABLE contest_events (
    contest_id INT REFERENCES contests(id),
    event_id INT REFERENCES events(id),
    PRIMARY KEY (contest_id, event_id)
);
```

---

## IMPLEMENTATION PLAN

### Phase 1: Backend — Scoring Rules

**Задачи:**
- Добавить `ContestTypeTotalizator` в rules.go
- Создать структуру `TotalizatorRules`
- Обновить `ParseRules` и `Validate`
- Обновить calculator.go для подсчёта очков

### Phase 2: Backend — Contest Validation

**Задачи:**
- Добавить валидацию количества матчей в contest_events
- Создать endpoint для привязки матчей к конкурсу
- Обновить proto если нужны новые методы

### Phase 3: Frontend — Contest Creation

**Задачи:**
- Добавить тип "Тотализатор" в ScoringRulesEditor
- Создать EventSelector для выбора матчей
- Интегрировать в ContestForm

### Phase 4: Telegram Bot

**Задачи:**
- Показывать матчи Тотализатора корректно
- Поддержка прогнозов (уже работает через contest_events)

---

## STEP-BY-STEP TASKS

### Task 1: UPDATE backend/shared/scoring/rules.go — Add Totalizator Type

- **IMPLEMENT**: Добавить константу `ContestTypeTotalizator ContestType = "totalizator"`
- **IMPLEMENT**: Создать структуру:
  ```go
  type TotalizatorRules struct {
      EventCount int                   `json:"event_count"` // default 15
      Scoring    StandardScoringRules  `json:"scoring"`     // правила подсчёта
  }
  ```
- **IMPLEMENT**: Добавить поле в ContestRules:
  ```go
  Totalizator *TotalizatorRules `json:"totalizator,omitempty"`
  ```
- **IMPLEMENT**: Обновить `ParseRules()` для totalizator
- **IMPLEMENT**: Обновить `Validate()` для totalizator (event_count 5-30)
- **IMPLEMENT**: Добавить `DefaultTotalizatorRules()`:
  ```go
  func DefaultTotalizatorRules() TotalizatorRules {
      return TotalizatorRules{
          EventCount: 15,
          Scoring:    DefaultStandardRules(),
      }
  }
  ```
- **VALIDATE**: `cd backend && go build ./...`

### Task 2: UPDATE backend/shared/scoring/calculator.go — Support Totalizator

- **IMPLEMENT**: В `CalculateScore()` добавить case для totalizator
- **MIRROR**: Использовать ту же логику что и для standard (StandardScoringRules)
- **GOTCHA**: Totalizator использует `rules.Totalizator.Scoring` вместо `rules.Standard`
- **VALIDATE**: `cd backend && go test ./shared/scoring/...`

### Task 3: CREATE database migration — contest_events если не существует

- **CHECK**: Таблица `contest_events` уже существует
- **IMPLEMENT**: Если нет, создать миграцию:
  ```sql
  CREATE TABLE IF NOT EXISTS contest_events (
      contest_id BIGINT NOT NULL,
      event_id BIGINT NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (contest_id, event_id),
      FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
      FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
  );
  ```
- **VALIDATE**: `docker exec -it sports-predictions-db psql -U postgres -d predictions -c "\d contest_events"`

### Task 4: UPDATE backend/contest-service — Validation

- **IMPLEMENT**: При создании/активации конкурса типа totalizator проверять:
  - Количество матчей в contest_events == rules.Totalizator.EventCount
  - Все матчи существуют и scheduled/live
- **PATTERN**: Добавить метод `ValidateTotalizatorEvents(contestID, expectedCount)` в service
- **VALIDATE**: `cd backend && go build ./...`

### Task 5: UPDATE proto — AddEventsToContest RPC (если нужен)

- **CHECK**: Есть ли метод для привязки events к contest
- **IMPLEMENT**: Если нет, добавить в prediction.proto:
  ```protobuf
  rpc AddEventsToContest(AddEventsToContestRequest) returns (AddEventsToContestResponse);
  
  message AddEventsToContestRequest {
      uint64 contest_id = 1;
      repeated uint64 event_ids = 2;
  }
  ```
- **VALIDATE**: `cd backend && make proto`

### Task 6: UPDATE frontend ScoringRulesEditor.tsx — Add Totalizator Option

- **IMPLEMENT**: Добавить "totalizator" в Select типов конкурса
- **IMPLEMENT**: При выборе totalizator показать:
  - Input для количества матчей (default 15, min 5, max 30)
  - StandardScoringRules поля (exact_score, goal_difference и т.д.)
- **PATTERN**: Следовать паттерну из risky секции
- **VALIDATE**: `cd frontend && npm run build`

### Task 7: CREATE frontend EventSelector.tsx — Match Selection Component

- **IMPLEMENT**: Компонент для выбора матчей:
  - Загружает все events (не фильтруя по sport_type)
  - Группировка по дате / лиге
  - Checkbox selection
  - Показывает выбрано X из Y
  - Валидация лимита
- **PATTERN**: Использовать Ant Design Table с rowSelection
- **IMPORTS**: `import { Table, Input, DatePicker, Tag } from 'antd'`
- **VALIDATE**: `cd frontend && npm run build`

### Task 8: UPDATE frontend ContestForm.tsx — Integrate EventSelector

- **IMPLEMENT**: При type="totalizator" показать EventSelector
- **IMPLEMENT**: Сохранять выбранные event_ids
- **IMPLEMENT**: При submit отправлять на AddEventsToContest
- **GOTCHA**: Порядок: сначала создать contest, потом привязать events
- **VALIDATE**: Manual testing in browser

### Task 9: UPDATE bots/telegram/bot/contests.go — Show Totalizator

- **IMPLEMENT**: При показе конкурса типа totalizator:
  - Показать "🎰 Тотализатор" вместо sport_type
  - Показать количество матчей "15 матчей"
- **PATTERN**: Следовать паттерну из handleContestList
- **VALIDATE**: Manual testing in Telegram

### Task 10: UPDATE bots/telegram/bot/predictions.go — Totalizator Support

- **CHECK**: Текущий код уже использует contest_events для фильтрации
- **IMPLEMENT**: Если нужно, добавить сортировку матчей по дате
- **VALIDATE**: Manual testing in Telegram

---

## TESTING STRATEGY

### Unit Tests

- `backend/shared/scoring/rules_test.go`:
  - TestParseTotalizatorRules
  - TestValidateTotalizatorRules
  - TestDefaultTotalizatorRules

- `backend/shared/scoring/calculator_test.go`:
  - TestCalculateScoreTotalizator

### Integration Tests

- Создание конкурса типа totalizator
- Привязка 15 матчей
- Создание прогноза
- Подсчёт очков

### Edge Cases

- [ ] Меньше 15 матчей выбрано
- [ ] Больше максимума матчей
- [ ] Матч уже завершён при создании
- [ ] Один матч в двух конкурсах

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build

```bash
cd ~/sports-predictions-contests/backend && go build ./...
cd ~/sports-predictions-contests/frontend && npm run build
```

### Level 2: Unit Tests

```bash
cd ~/sports-predictions-contests/backend && go test ./shared/scoring/... -v
```

### Level 3: Integration

```bash
docker-compose -f docker-compose.yml logs -f contest-service prediction-service
```

### Level 4: Manual Validation

1. Открыть https://forecasts.dinamchiki.ru/admin/contests/new
2. Выбрать тип "Тотализатор"
3. Настроить правила
4. Выбрать 15 матчей
5. Создать конкурс
6. Проверить в Telegram боте

---

## ACCEPTANCE CRITERIA

- [ ] Новый тип конкурса "totalizator" доступен в админке
- [ ] Можно выбрать матчи из разных лиг
- [ ] Валидация количества матчей работает
- [ ] Прогнозы создаются и считаются корректно
- [ ] Telegram бот показывает матчи Тотализатора
- [ ] Один матч может быть в нескольких конкурсах

---

## COMPLETION CHECKLIST

- [ ] Task 1-2: Backend scoring rules
- [ ] Task 3-5: Backend validation & API
- [ ] Task 6-8: Frontend UI
- [ ] Task 9-10: Telegram bot
- [ ] All tests pass
- [ ] Manual E2E testing complete

---

## NOTES

### Архитектурное решение

Тотализатор использует **те же правила подсчёта** что и Standard, только обёрнутые в TotalizatorRules. Это позволяет:
1. Переиспользовать calculator.go
2. Не дублировать логику
3. В будущем добавить специфичные для Тотализатора поля

### Связь Event ↔ Contest

Таблица `contest_events` уже поддерживает many-to-many. Один Event может быть в нескольких Contest. Это ключевое требование для Тотализатора.

### UI Flow

1. Admin создаёт Contest (type=totalizator, rules JSON)
2. Admin выбирает Events через EventSelector
3. Submit → CreateContest → AddEventsToContest
4. Users видят матчи и делают прогнозы

### Отличие от Standard

| Аспект | Standard | Totalizator |
|--------|----------|-------------|
| Sport Type | Один | Любой/смешанный |
| Выбор матчей | Автоматически | Вручную |
| Кол-во матчей | Не ограничено | Фиксированное (15) |
| Правила подсчёта | StandardScoringRules | StandardScoringRules |
