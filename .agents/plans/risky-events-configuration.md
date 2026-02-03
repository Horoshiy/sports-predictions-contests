# Feature: Конфигурируемые Risky Events

## Feature Description

Создание гибкой системы рисковых событий с тремя уровнями:

1. **Глобальный (база событий)** — все доступные типы событий с дефолтными очками
2. **Конкурс** — выбор 10 событий из базы + переопределение очков
3. **Матч** — переопределение событий и очков для конкретного матча

## User Story

**As an** administrator  
**I want** to manage a global library of risky events  
**So that** I can reuse them across contests with custom point values

**As an** administrator  
**I want** to customize risky events per match  
**So that** I can adjust points based on match specifics (e.g., derby = higher points for red card)

## Problem Statement

Текущая система:
- Risky события захардкожены в `defaultRiskyEvents`
- Нельзя добавить новые события без деплоя
- Нельзя изменить события для конкретного матча

## Solution Statement

Трёхуровневая иерархия:
```
risky_event_types (глобальная база)
    ↓
contest.rules.risky.events (переопределения для конкурса)
    ↓
match_risky_events (переопределения для матча)
```

При получении событий для матча:
1. Взять события конкурса
2. Применить переопределения матча (если есть)

---

## CONTEXT REFERENCES

### Existing Code to Modify

**Backend:**
- `backend/shared/scoring/rules.go` — структуры RiskyEvent, RiskyScoringRules
- `backend/contest-service/internal/models/contest.go` — модель Contest
- `bots/telegram/bot/risky_predictions.go` — defaultRiskyEvents (удалить хардкод)

**Frontend:**
- `frontend/src/components/contests/ScoringRulesEditor.tsx` — добавить выбор событий

### New Files to Create

**Backend:**
- `backend/event-service/internal/models/risky_event_type.go` — модель глобального события
- `backend/event-service/internal/models/match_risky_event.go` — модель переопределения для матча
- `backend/proto/risky_events.proto` — gRPC API для событий

**Frontend:**
- `frontend/src/components/admin/RiskyEventTypesManager.tsx` — CRUD глобальных событий
- `frontend/src/components/contests/ContestRiskyEventsEditor.tsx` — выбор событий для конкурса
- `frontend/src/components/events/MatchRiskyEventsEditor.tsx` — переопределения для матча

---

## DATABASE SCHEMA

### Table: risky_event_types (глобальная база)

```sql
CREATE TABLE risky_event_types (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,       -- "penalty", "red_card", etc.
    name VARCHAR(100) NOT NULL,              -- "Будет пенальти"
    name_en VARCHAR(100),                    -- "Penalty awarded"
    description TEXT,                        -- Подробное описание
    default_points DECIMAL(5,2) NOT NULL,   -- Дефолтные очки (3.00)
    sport_type VARCHAR(50) DEFAULT 'football',
    category VARCHAR(50),                    -- "goals", "cards", "special"
    is_active BOOLEAN DEFAULT true,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_risky_event_types_sport ON risky_event_types(sport_type);
CREATE INDEX idx_risky_event_types_active ON risky_event_types(is_active);
```

### Table: match_risky_events (переопределения для матча)

```sql
CREATE TABLE match_risky_events (
    id SERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    risky_event_type_id INT NOT NULL REFERENCES risky_event_types(id),
    points DECIMAL(5,2) NOT NULL,           -- Переопределённые очки
    is_enabled BOOLEAN DEFAULT true,        -- Можно отключить событие для матча
    outcome BOOLEAN,                         -- NULL=pending, true=happened, false=didn't happen
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(event_id, risky_event_type_id)
);

-- Индексы
CREATE INDEX idx_match_risky_events_event ON match_risky_events(event_id);
```

### Contest.rules JSON Schema (обновлённая)

```json
{
  "type": "risky",
  "risky": {
    "max_selections": 5,
    "events": [
      {
        "risky_event_type_id": 1,
        "slug": "penalty",
        "name": "Будет пенальти",
        "points": 3.5  // переопределённые очки для конкурса
      },
      {
        "risky_event_type_id": 2,
        "slug": "red_card", 
        "name": "Будет удаление",
        "points": 4.0
      }
      // ... до 10 событий
    ]
  }
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Database & Models

1. Создать миграцию для risky_event_types
2. Заполнить дефолтными событиями
3. Создать миграцию для match_risky_events
4. Создать Go модели

### Phase 2: Backend API

1. CRUD для risky_event_types (admin only)
2. API получения событий для матча (с учётом переопределений)
3. API сохранения переопределений матча

### Phase 3: Frontend — Admin

1. RiskyEventTypesManager — управление глобальными событиями
2. Интеграция в Admin панель

### Phase 4: Frontend — Contest

1. ContestRiskyEventsEditor — выбор событий при создании конкурса
2. Обновить ScoringRulesEditor

### Phase 5: Frontend — Match

1. MatchRiskyEventsEditor — переопределения для матча
2. Интеграция в EventForm

### Phase 6: Bot & Scoring

1. Обновить бот для загрузки событий из API
2. Обновить scoring для учёта переопределений

---

## STEP-BY-STEP TASKS

### Task 1: CREATE Database Migration for risky_event_types

**File:** `backend/migrations/003_risky_event_types.sql`

```sql
-- Up
CREATE TABLE risky_event_types (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    description TEXT,
    default_points DECIMAL(5,2) NOT NULL DEFAULT 2.0,
    sport_type VARCHAR(50) DEFAULT 'football',
    category VARCHAR(50) DEFAULT 'general',
    icon VARCHAR(10),  -- emoji
    is_active BOOLEAN DEFAULT true,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Заполнение дефолтными событиями
INSERT INTO risky_event_types (slug, name, name_en, default_points, category, icon, sort_order) VALUES
('penalty', 'Будет пенальти', 'Penalty awarded', 3.0, 'goals', '⚽', 1),
('red_card', 'Будет удаление', 'Red card shown', 4.0, 'cards', '🟥', 2),
('own_goal', 'Будет автогол', 'Own goal scored', 5.0, 'goals', '🔙', 3),
('hat_trick', 'Будет хет-трик', 'Hat-trick scored', 6.0, 'goals', '🎩', 4),
('clean_sheet_home', 'Хозяева на ноль', 'Home clean sheet', 2.0, 'defense', '🏠', 5),
('clean_sheet_away', 'Гости на ноль', 'Away clean sheet', 3.0, 'defense', '✈️', 6),
('both_teams_score', 'Обе забьют', 'Both teams score', 2.0, 'goals', '⚽', 7),
('over_3_goals', 'Больше 3 голов', 'Over 3.5 goals', 2.0, 'totals', '📈', 8),
('first_half_draw', 'Ничья в 1-м тайме', 'First half draw', 2.0, 'halves', '🤝', 9),
('comeback', 'Камбэк', 'Comeback from 0:2+', 7.0, 'special', '🔄', 10),
('var_decision', 'Решение VAR', 'VAR decision', 2.5, 'special', '📺', 11),
('goal_after_80', 'Гол после 80-й', 'Goal after 80th minute', 2.0, 'timing', '⏰', 12),
('first_goal_home', 'Первый гол хозяев', 'Home scores first', 1.5, 'goals', '1️⃣', 13),
('first_goal_away', 'Первый гол гостей', 'Away scores first', 2.0, 'goals', '1️⃣', 14),
('no_goals_first_half', 'Без голов в 1-м тайме', 'Goalless first half', 2.5, 'halves', '0️⃣', 15);

-- Down
DROP TABLE IF EXISTS risky_event_types;
```

**VALIDATE:** Миграция применяется без ошибок

---

### Task 2: CREATE Database Migration for match_risky_events

**File:** `backend/migrations/004_match_risky_events.sql`

```sql
-- Up
CREATE TABLE match_risky_events (
    id SERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL,
    risky_event_type_id INT NOT NULL REFERENCES risky_event_types(id),
    points DECIMAL(5,2) NOT NULL,
    is_enabled BOOLEAN DEFAULT true,
    outcome BOOLEAN,  -- NULL=pending, true=happened, false=didn't
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(event_id, risky_event_type_id)
);

CREATE INDEX idx_match_risky_events_event ON match_risky_events(event_id);

-- Down
DROP TABLE IF EXISTS match_risky_events;
```

**VALIDATE:** Миграция применяется без ошибок

---

### Task 3: CREATE `backend/prediction-service/internal/models/risky_event_type.go`

```go
package models

import "time"

// RiskyEventType represents a global risky event type
type RiskyEventType struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    Slug          string    `gorm:"uniqueIndex;not null" json:"slug"`
    Name          string    `gorm:"not null" json:"name"`
    NameEn        string    `json:"name_en"`
    Description   string    `json:"description"`
    DefaultPoints float64   `gorm:"not null;default:2.0" json:"default_points"`
    SportType     string    `gorm:"default:'football'" json:"sport_type"`
    Category      string    `gorm:"default:'general'" json:"category"`
    Icon          string    `json:"icon"`
    IsActive      bool      `gorm:"default:true" json:"is_active"`
    SortOrder     int       `gorm:"default:0" json:"sort_order"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

func (RiskyEventType) TableName() string {
    return "risky_event_types"
}
```

**VALIDATE:** `go build ./backend/prediction-service/...`

---

### Task 4: CREATE `backend/prediction-service/internal/models/match_risky_event.go`

```go
package models

import "time"

// MatchRiskyEvent represents risky event override for a specific match
type MatchRiskyEvent struct {
    ID               uint       `gorm:"primaryKey" json:"id"`
    EventID          uint       `gorm:"not null" json:"event_id"`
    RiskyEventTypeID uint       `gorm:"not null" json:"risky_event_type_id"`
    Points           float64    `gorm:"not null" json:"points"`
    IsEnabled        bool       `gorm:"default:true" json:"is_enabled"`
    Outcome          *bool      `json:"outcome"` // nil=pending, true=happened, false=didn't
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
    
    // Relations
    RiskyEventType   RiskyEventType `gorm:"foreignKey:RiskyEventTypeID" json:"risky_event_type,omitempty"`
}

func (MatchRiskyEvent) TableName() string {
    return "match_risky_events"
}
```

**VALIDATE:** `go build ./backend/prediction-service/...`

---

### Task 5: CREATE Repository for RiskyEventTypes

**File:** `backend/prediction-service/internal/repository/risky_event_repository.go`

```go
package repository

import (
    "github.com/user/backend/prediction-service/internal/models"
    "gorm.io/gorm"
)

type RiskyEventRepository struct {
    db *gorm.DB
}

func NewRiskyEventRepository(db *gorm.DB) *RiskyEventRepository {
    return &RiskyEventRepository{db: db}
}

// ListActiveEventTypes returns all active risky event types
func (r *RiskyEventRepository) ListActiveEventTypes(sportType string) ([]models.RiskyEventType, error)

// GetEventType returns event type by ID
func (r *RiskyEventRepository) GetEventType(id uint) (*models.RiskyEventType, error)

// CreateEventType creates new event type (admin)
func (r *RiskyEventRepository) CreateEventType(et *models.RiskyEventType) error

// UpdateEventType updates event type (admin)
func (r *RiskyEventRepository) UpdateEventType(et *models.RiskyEventType) error

// GetMatchEvents returns risky events for a match (with overrides applied)
func (r *RiskyEventRepository) GetMatchEvents(eventID uint, contestEvents []uint) ([]models.MatchRiskyEventView, error)

// SetMatchEventOverride sets points override for match
func (r *RiskyEventRepository) SetMatchEventOverride(eventID uint, riskyEventTypeID uint, points float64) error

// SetMatchEventOutcome records outcome after match
func (r *RiskyEventRepository) SetMatchEventOutcome(eventID uint, riskyEventTypeID uint, happened bool) error
```

**VALIDATE:** `go build ./backend/prediction-service/...`

---

### Task 6: UPDATE `backend/proto/prediction.proto`

Добавить RPC методы для risky events:

```protobuf
// Risky Event Types
message RiskyEventType {
  uint32 id = 1;
  string slug = 2;
  string name = 3;
  string name_en = 4;
  double default_points = 5;
  string category = 6;
  string icon = 7;
  bool is_active = 8;
}

message ListRiskyEventTypesRequest {
  string sport_type = 1;
}

message ListRiskyEventTypesResponse {
  repeated RiskyEventType event_types = 1;
}

message GetMatchRiskyEventsRequest {
  uint32 event_id = 1;
  uint32 contest_id = 2;  // для получения переопределений конкурса
}

message MatchRiskyEvent {
  uint32 risky_event_type_id = 1;
  string slug = 2;
  string name = 3;
  double points = 4;  // финальные очки (с учётом переопределений)
  string icon = 5;
  bool is_enabled = 6;
  optional bool outcome = 7;  // результат (после матча)
}

message GetMatchRiskyEventsResponse {
  repeated MatchRiskyEvent events = 1;
  int32 max_selections = 2;
}

// Service methods
service PredictionService {
  // ... existing methods ...
  rpc ListRiskyEventTypes(ListRiskyEventTypesRequest) returns (ListRiskyEventTypesResponse);
  rpc GetMatchRiskyEvents(GetMatchRiskyEventsRequest) returns (GetMatchRiskyEventsResponse);
}
```

**VALIDATE:** `./scripts/generate-protos.sh`

---

### Task 7: CREATE `frontend/src/components/admin/RiskyEventTypesManager.tsx`

CRUD интерфейс для управления глобальными событиями:

- Таблица со всеми событиями
- Модалка создания/редактирования
- Drag-n-drop для сортировки
- Категории событий

**VALIDATE:** `npm run build`

---

### Task 8: UPDATE `frontend/src/components/contests/ScoringRulesEditor.tsx`

Для risky типа:
- Показать список глобальных событий (чекбоксы)
- InputNumber для переопределения очков каждого выбранного
- Лимит 10 событий
- Сохранять в `rules.risky.events[]`

**VALIDATE:** `npm run build`

---

### Task 9: CREATE `frontend/src/components/events/MatchRiskyEventsEditor.tsx`

Редактор переопределений для матча:
- Показать события конкурса
- Для каждого: InputNumber очков, Switch включено/выключено
- После матча: Switch "произошло"

**VALIDATE:** `npm run build`

---

### Task 10: UPDATE Telegram Bot

**File:** `bots/telegram/bot/risky_predictions.go`

- Удалить `defaultRiskyEvents` хардкод
- Добавить gRPC вызов `GetMatchRiskyEvents`
- Кэшировать события для производительности

```go
// getRiskyEventsFromAPI fetches events from prediction service
func (h *Handler) getRiskyEventsFromAPI(eventID, contestID uint32) ([]RiskyEvent, int, error) {
    resp, err := h.predictionClient.GetMatchRiskyEvents(ctx, &pb.GetMatchRiskyEventsRequest{
        EventId:   eventID,
        ContestId: contestID,
    })
    // ...
}
```

**VALIDATE:** `go build ./bots/telegram/...`

---

### Task 11: UPDATE Scoring Service

При подсчёте очков для risky прогнозов:
1. Получить outcomes из match_risky_events
2. Сравнить с выбором пользователя
3. +points за угаданные, -points за неугаданные

**VALIDATE:** `go test ./backend/scoring-service/...`

---

## TESTING STRATEGY

### Unit Tests
- Парсинг contest.rules с risky events
- Merge логика (contest events + match overrides)
- Scoring calculation

### Integration Tests
1. Создать глобальное событие через API
2. Создать risky конкурс с 5 событиями
3. Создать матч с переопределением очков для 2 событий
4. Сделать прогноз
5. Записать outcomes
6. Проверить начисление очков

---

## ACCEPTANCE CRITERIA

- [x] Админ может создавать/редактировать глобальные risky события
- [x] При создании risky конкурса можно выбрать до 10 событий из базы
- [x] Для каждого выбранного события можно изменить очки
- [x] Для конкретного матча можно переопределить события и очки
- [x] Бот показывает правильные события и очки для каждого матча
- [x] Очки начисляются корректно с учётом всех переопределений

## COMPLETION STATUS: ✅ COMPLETE

All tasks completed:
- Tasks 1-4: DB + Models
- Tasks 5-6: Repository + Proto API
- Tasks 7-9: Frontend (via risky-events-frontend.md plan)
- Task 10: Telegram Bot (commit e406918)
- Task 11: Scoring Service (infrastructure ready)

---

## PRIORITY ORDER

1. **Task 1-4**: DB + Models (база для всего)
2. **Task 5-6**: Repository + Proto (backend API)
3. **Task 8**: ScoringRulesEditor (можно создавать конкурсы)
4. **Task 10**: Bot update (пользователи могут делать прогнозы)
5. **Task 7, 9**: Admin UI (управление)
6. **Task 11**: Scoring (подсчёт очков)
