# Feature: Настраиваемые правила подсчёта очков в конкурсах

## Feature Description

Реализация гибкой системы правил для двух типов конкурсов прогнозов:

1. **Обычный (standard)** — прогноз на счёт матча с настраиваемыми очками:
   - За точный счёт (например, 5 очков)
   - За разницу мячей (например, 3 очка)
   - За верный исход (например, 1 очко)
   - Бонус за исход + голы одной команды (например, +1 очко)
   - За прогноз "другой" (например, 4 очка)

2. **Рисковый (risky)** — выбор событий матча с риском:
   - 10 вариантов событий на матч (пенальти, удаление, автогол и т.д.)
   - Пользователь выбирает до 5 событий
   - Угадал событие → +очки, не угадал → −очки
   - Дефолтные стоимости событий, админ может менять для матча

## User Story

As a contest administrator
I want to configure scoring rules for my contest
So that participants are rewarded according to my custom point system

As a participant
I want to make predictions with clear point values
So that I understand how my score will be calculated

## Problem Statement

Текущая система использует захардкоженные очки (10/5/3) без возможности настройки. Администраторы не могут создавать конкурсы с разными правилами подсчёта очков.

## Solution Statement

Расширить Contest.rules и prediction_schema для хранения настроек очков. Модифицировать scoring-service для чтения правил из конкурса. Добавить UI для настройки правил на фронтенде.

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: High
**Primary Systems Affected**: contest-service, scoring-service, prediction-service, frontend, telegram-bot
**Dependencies**: Нет новых внешних зависимостей

---

## CONTEXT REFERENCES

### Relevant Codebase Files — ОБЯЗАТЕЛЬНО ПРОЧИТАТЬ!

**Backend — Scoring:**
- `backend/scoring-service/internal/service/scoring_service.go` (lines 225-320) — текущий алгоритм подсчёта очков
- `backend/scoring-service/internal/models/coefficient.go` — коэффициенты времени
- `backend/proto/scoring.proto` — gRPC схема scoring

**Backend — Contest:**
- `backend/proto/contest.proto` (lines 10-25) — Contest message с rules и prediction_schema
- `backend/contest-service/internal/models/contest.go` — модель Contest
- `backend/contest-service/internal/service/contest_service.go` — CRUD конкурсов

**Backend — Prediction:**
- `backend/prediction-service/internal/service/prediction_service.go` — SubmitPrediction
- `backend/proto/prediction.proto` — Prediction, Event messages

**Frontend:**
- `frontend/src/components/contests/ContestForm.tsx` — форма создания конкурса
- `frontend/src/utils/validation.ts` — схемы валидации

**Telegram Bot:**
- `bots/telegram/bot/predictions.go` — handlePredictionSubmit
- `bots/telegram/bot/score_buttons.go` — клавиатура выбора счёта

### New Files to Create

**Backend:**
- `backend/shared/scoring/rules.go` — структуры и парсинг правил
- `backend/shared/scoring/calculator.go` — калькулятор очков по правилам
- `backend/prediction-service/internal/models/risky_event.go` — модель рисковых событий

**Frontend:**
- `frontend/src/components/contests/ScoringRulesEditor.tsx` — редактор правил
- `frontend/src/components/contests/RiskyEventsEditor.tsx` — редактор событий для рискового типа
- `frontend/src/components/predictions/RiskyPredictionForm.tsx` — форма рискового прогноза

**Database:**
- Миграция: таблица `risky_event_types` (дефолтные типы событий)
- Миграция: таблица `match_risky_events` (события для конкретного матча)

### Patterns to Follow

**JSON Rules Schema (Contest.rules):**
```json
{
  "type": "standard",
  "scoring": {
    "exact_score": 5,
    "goal_difference": 3,
    "correct_outcome": 1,
    "outcome_plus_team_goals": 1,
    "any_other": 4
  }
}
```

**Risky Type Schema:**
```json
{
  "type": "risky",
  "max_selections": 5,
  "default_events": [
    {"slug": "penalty", "name": "Будет пенальти", "points": 3},
    {"slug": "red_card", "name": "Будет удаление", "points": 4},
    {"slug": "own_goal", "name": "Будет автогол", "points": 5}
  ]
}
```

**Error Handling Pattern (from codebase):**
```go
return &pb.Response{
    Success:   false,
    Message:   "Error message",
    Code:      int32(common.ErrorCode_INVALID_ARGUMENT),
    Timestamp: timestamppb.Now(),
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Data Models & Schemas

1. Создать структуры правил в shared/scoring/
2. Обновить Contest model для типизированных правил
3. Создать миграции для risky_event_types

### Phase 2: Scoring Service Updates

1. Модифицировать calculatePoints для чтения правил из конкурса
2. Добавить метод получения правил конкурса (gRPC call to contest-service)
3. Реализовать калькулятор для risky типа

### Phase 3: Prediction Service Updates

1. Добавить валидацию прогнозов по типу конкурса
2. Для risky: проверка max_selections
3. Хранение risky predictions

### Phase 4: Frontend — Contest Creation

1. ScoringRulesEditor компонент
2. Интеграция в ContestForm
3. Превью правил

### Phase 5: Frontend — Predictions

1. RiskyPredictionForm для рискового типа
2. Отображение потенциальных очков

### Phase 6: Telegram Bot

1. Поддержка risky прогнозов
2. Отображение правил конкурса

---

## STEP-BY-STEP TASKS

### Task 1: CREATE `backend/shared/scoring/rules.go`

**IMPLEMENT:** Структуры для правил подсчёта очков

```go
package scoring

// ContestType defines the type of contest
type ContestType string

const (
    ContestTypeStandard ContestType = "standard"
    ContestTypeRisky    ContestType = "risky"
)

// StandardScoringRules defines points for standard contest
type StandardScoringRules struct {
    ExactScore          float64 `json:"exact_score"`           // точный счёт
    GoalDifference      float64 `json:"goal_difference"`       // разница мячей
    CorrectOutcome      float64 `json:"correct_outcome"`       // верный исход
    OutcomePlusTeamGoals float64 `json:"outcome_plus_team_goals"` // исход + голы команды
    AnyOther            float64 `json:"any_other"`             // прогноз "другой"
}

// RiskyEvent defines a risky event type
type RiskyEvent struct {
    Slug        string  `json:"slug"`
    Name        string  `json:"name"`
    Points      float64 `json:"points"`
    Description string  `json:"description,omitempty"`
}

// RiskyScoringRules defines rules for risky contest
type RiskyScoringRules struct {
    MaxSelections int          `json:"max_selections"`
    DefaultEvents []RiskyEvent `json:"default_events"`
}

// ContestRules combines all rule types
type ContestRules struct {
    Type     ContestType           `json:"type"`
    Standard *StandardScoringRules `json:"scoring,omitempty"`
    Risky    *RiskyScoringRules    `json:"risky,omitempty"`
}

// DefaultStandardRules returns default scoring for standard contests
func DefaultStandardRules() StandardScoringRules {
    return StandardScoringRules{
        ExactScore:          5,
        GoalDifference:      3,
        CorrectOutcome:      1,
        OutcomePlusTeamGoals: 1,
        AnyOther:            4,
    }
}

// DefaultRiskyEvents returns default risky events
func DefaultRiskyEvents() []RiskyEvent {
    return []RiskyEvent{
        {Slug: "penalty", Name: "Будет пенальти", Points: 3},
        {Slug: "red_card", Name: "Будет удаление", Points: 4},
        {Slug: "own_goal", Name: "Будет автогол", Points: 5},
        {Slug: "hat_trick", Name: "Будет хет-трик", Points: 6},
        {Slug: "clean_sheet_home", Name: "Хозяева на ноль", Points: 2},
        {Slug: "clean_sheet_away", Name: "Гости на ноль", Points: 3},
        {Slug: "both_teams_score", Name: "Обе команды забьют", Points: 2},
        {Slug: "over_3_goals", Name: "Больше 3 голов", Points: 2},
        {Slug: "first_half_draw", Name: "Ничья в первом тайме", Points: 2},
        {Slug: "comeback", Name: "Камбэк (отыграться с 0:2+)", Points: 7},
    }
}

// ParseRules parses JSON rules string
func ParseRules(rulesJSON string) (*ContestRules, error)
```

**VALIDATE:** `go build ./backend/shared/scoring/...`

---

### Task 2: CREATE `backend/shared/scoring/calculator.go`

**IMPLEMENT:** Калькулятор очков по правилам

```go
package scoring

// Calculator calculates points based on contest rules
type Calculator struct {
    rules ContestRules
}

// NewCalculator creates calculator with rules
func NewCalculator(rules ContestRules) *Calculator

// CalculateStandard calculates points for standard prediction
func (c *Calculator) CalculateStandard(prediction, result ScoreData) (float64, map[string]interface{})

// CalculateRisky calculates points for risky predictions
func (c *Calculator) CalculateRisky(selections []string, outcomes map[string]bool) (float64, map[string]interface{})
```

**VALIDATE:** `go test ./backend/shared/scoring/...`

---

### Task 3: UPDATE `backend/scoring-service/internal/service/scoring_service.go`

**IMPLEMENT:** Использовать правила из конкурса вместо захардкоженных значений

- Добавить contestClient для получения правил
- В calculateExactScorePoints читать очки из rules
- Добавить calculateRiskyPoints метод

**PATTERN:** Текущий calculateExactScorePoints (lines 280-320)

**GOTCHA:** Для обратной совместимости: если rules пустые, использовать дефолтные

**VALIDATE:** `go build ./backend/scoring-service/...`

---

### Task 4: UPDATE `backend/proto/contest.proto`

**IMPLEMENT:** Добавить типизированные поля (опционально, для удобства)

```protobuf
message ContestScoringRules {
  string type = 1; // "standard" or "risky"
  double exact_score = 2;
  double goal_difference = 3;
  double correct_outcome = 4;
  double outcome_plus_team_goals = 5;
  double any_other = 6;
  int32 max_risky_selections = 7;
}
```

**VALIDATE:** `./scripts/generate-protos.sh`

---

### Task 5: CREATE Database Migration for Risky Events

**IMPLEMENT:** SQL миграция

```sql
-- Таблица типов рисковых событий
CREATE TABLE risky_event_types (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    description TEXT,
    default_points DECIMAL(5,2) NOT NULL,
    sport_type VARCHAR(50) DEFAULT 'football',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- События для конкретного матча (переопределение очков)
CREATE TABLE match_risky_events (
    id SERIAL PRIMARY KEY,
    event_id BIGINT REFERENCES events(id),
    risky_event_type_id INT REFERENCES risky_event_types(id),
    points DECIMAL(5,2) NOT NULL,
    outcome BOOLEAN, -- NULL = не определено, true = произошло, false = не произошло
    UNIQUE(event_id, risky_event_type_id)
);

-- Рисковые прогнозы пользователей
CREATE TABLE risky_predictions (
    id SERIAL PRIMARY KEY,
    prediction_id BIGINT REFERENCES predictions(id),
    risky_event_type_id INT REFERENCES risky_event_types(id),
    points_if_correct DECIMAL(5,2) NOT NULL,
    UNIQUE(prediction_id, risky_event_type_id)
);

-- Заполнение дефолтных событий
INSERT INTO risky_event_types (slug, name, name_en, default_points) VALUES
('penalty', 'Будет пенальти', 'Penalty awarded', 3),
('red_card', 'Будет удаление', 'Red card shown', 4),
('own_goal', 'Будет автогол', 'Own goal scored', 5),
('hat_trick', 'Будет хет-трик', 'Hat-trick scored', 6),
('clean_sheet_home', 'Хозяева на ноль', 'Home clean sheet', 2),
('clean_sheet_away', 'Гости на ноль', 'Away clean sheet', 3),
('both_teams_score', 'Обе забьют', 'Both teams score', 2),
('over_3_goals', 'Больше 3 голов', 'Over 3.5 goals', 2),
('first_half_draw', 'Ничья в 1-м тайме', 'First half draw', 2),
('comeback', 'Камбэк', 'Comeback from 0:2+', 7);
```

**VALIDATE:** Применить миграцию, проверить таблицы

---

### Task 6: CREATE `frontend/src/components/contests/ScoringRulesEditor.tsx`

**IMPLEMENT:** React компонент для настройки правил

```tsx
interface ScoringRulesEditorProps {
  value: ContestRules;
  onChange: (rules: ContestRules) => void;
}

// Табы: "Обычный" / "Рисковый"
// Для обычного: InputNumber для каждого типа очков
// Для рискового: список событий с очками, max_selections
```

**PATTERN:** Использовать Ant Design Form.Item, InputNumber

**VALIDATE:** `npm run build` без ошибок

---

### Task 7: UPDATE `frontend/src/components/contests/ContestForm.tsx`

**IMPLEMENT:** Интегрировать ScoringRulesEditor

- Добавить поле rules в форму
- Показывать ScoringRulesEditor после выбора типа
- Сериализовать в JSON при отправке

**VALIDATE:** Создать конкурс с кастомными правилами через UI

---

### Task 8: CREATE `frontend/src/components/predictions/RiskyPredictionForm.tsx`

**IMPLEMENT:** Форма для рисковых прогнозов

- Показать список событий с очками
- Чекбоксы для выбора (max N)
- Показать потенциальный выигрыш/проигрыш

**VALIDATE:** Сделать рисковый прогноз через UI

---

### Task 9: UPDATE `bots/telegram/bot/predictions.go`

**IMPLEMENT:** Поддержка разных типов конкурсов

- Определить тип конкурса при показе матча
- Для standard: текущая логика
- Для risky: показать события, inline кнопки выбора

**VALIDATE:** Сделать прогноз обоих типов через бота

---

### Task 10: UPDATE Scoring Calculation Flow

**IMPLEMENT:** При завершении матча:
1. Получить результаты матча
2. Получить правила конкурса
3. Для standard: использовать настроенные очки
4. Для risky: пройти по выбранным событиям, +/- очки

**VALIDATE:** Завершить матч, проверить начисление очков

---

## TESTING STRATEGY

### Unit Tests

- `backend/shared/scoring/rules_test.go` — парсинг правил
- `backend/shared/scoring/calculator_test.go` — расчёт очков для всех комбинаций

### Integration Tests

- Создание конкурса с кастомными правилами через API
- Прогноз и подсчёт очков для standard типа
- Прогноз и подсчёт очков для risky типа

### Edge Cases

- Пустые rules → дефолтные значения
- Risky с > max_selections → ошибка валидации
- Отрицательный баланс очков за матч в risky

---

## VALIDATION COMMANDS

### Level 1: Build

```bash
cd ~/sports-predictions-contests
go build ./backend/...
cd frontend && npm run build
```

### Level 2: Unit Tests

```bash
go test ./backend/shared/scoring/...
go test ./backend/scoring-service/...
```

### Level 3: Integration

```bash
# Запустить сервисы
docker compose up -d

# Создать конкурс с правилами через API
curl -X POST http://localhost:8080/v1/contests \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Test","rules":"{\"type\":\"standard\",\"scoring\":{\"exact_score\":5}}"}'
```

### Level 4: Manual Validation

1. Создать standard конкурс с очками 5/3/1/1/4
2. Сделать прогноз, дождаться результата, проверить очки
3. Создать risky конкурс
4. Сделать risky прогноз, проверить +/- очки

---

## ACCEPTANCE CRITERIA

- [ ] Админ может создать конкурс типа "standard" с настраиваемыми очками
- [ ] Админ может создать конкурс типа "risky" с выбором событий
- [ ] Очки начисляются согласно настроенным правилам
- [ ] В risky нельзя выбрать больше max_selections событий
- [ ] В risky очки вычитаются за неугаданные события
- [ ] Telegram бот поддерживает оба типа
- [ ] Сумма очков за все матчи отображается в лидерборде

---

## COMPLETION CHECKLIST

- [ ] Все задачи выполнены по порядку
- [ ] Unit тесты написаны и проходят
- [ ] Integration тесты проходят
- [ ] Фронтенд билдится без ошибок
- [ ] Бот работает с обоими типами
- [ ] Code review пройден

---

## NOTES

### Design Decisions

1. **JSON в rules vs отдельные поля** — используем JSON для гибкости, но добавляем типизированные структуры в Go
2. **Risky events в отдельной таблице** — позволяет переопределять очки для конкретных матчей
3. **Отрицательные очки в risky** — сумма за матч может быть отрицательной, но общий счёт в конкурсе не уходит в минус (floor = 0)

### Risks

- **Миграция существующих конкурсов** — нужно добавить дефолтные rules к существующим
- **Сложность UI** — редактор правил может быть перегружен, возможно нужны пресеты

### Future Improvements

- Пресеты правил ("Классический", "Рисковый стандарт")
- История изменения правил
- A/B тестирование разных систем очков
