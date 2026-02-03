# Feature: Эстафета (Relay Contest Type)

## Feature Description

Новый тип конкурса **"Эстафета"** — командный турнир:
- Админ настраивает количество участников в команде
- Админ определяет матчи и правила начисления очков
- Каждая команда имеет капитана
- **Капитан распределяет матчи** между участниками команды (один прогнозирует Италию, другой Францию и т.д.)
- Очки команды суммируются от всех участников

## User Story

**Как админ:**
- Хочу создать командный конкурс с X участниками в команде
- Хочу выбрать матчи из разных лиг для прогнозирования
- Хочу настроить правила начисления очков

**Как капитан команды:**
- Хочу распределить матчи между участниками моей команды
- Хочу видеть кто какие матчи прогнозирует
- Хочу видеть суммарные очки команды

**Как участник команды:**
- Хочу видеть только назначенные мне матчи
- Хочу делать прогнозы на свои матчи
- Хочу видеть вклад в общие очки команды

## Problem Statement

Сейчас все конкурсы индивидуальные. Нет механизма командного соревнования с распределением ответственности между участниками.

## Solution Statement

1. Добавить тип конкурса `relay` с настройками команды
2. Создать таблицу `relay_event_assignments` для распределения матчей
3. Добавить API для капитана (назначение матчей)
4. Ограничить прогнозы участников только назначенными матчами
5. Суммировать очки команды от всех участников

## Feature Metadata

**Feature Type**: New Capability  
**Estimated Complexity**: High  
**Primary Systems Affected**: contest-service, prediction-service, frontend, telegram-bot  
**Dependencies**: Существующие модели Team, TeamMember, TeamContestEntry

---

## CONTEXT REFERENCES

### Relevant Codebase Files

- `backend/shared/scoring/rules.go` — добавить ContestTypeRelay и RelayRules
- `backend/contest-service/internal/models/team.go` — структура Team с капитаном
- `backend/contest-service/internal/models/team_member.go` — TeamMember (role: captain/member)
- `backend/contest-service/internal/models/team_contest_entry.go` — TeamContestEntry с TotalPoints
- `backend/prediction-service/internal/repository/event_repository.go` — contest_events
- `backend/prediction-service/internal/service/prediction_service.go` — проверка прогнозов
- `frontend/src/components/contests/ScoringRulesEditor.tsx` — добавить relay тип
- `frontend/src/components/contests/EventSelector.tsx` — переиспользовать для выбора матчей

### New Files to Create

**Backend:**
- `backend/prediction-service/internal/models/relay_assignment.go` — модель распределения матчей
- `backend/prediction-service/internal/repository/relay_repository.go` — CRUD для assignments

**Frontend:**
- `frontend/src/components/relay/RelayAssignmentEditor.tsx` — UI для капитана
- `frontend/src/components/relay/TeamEventsList.tsx` — список матчей команды
- `frontend/src/pages/RelayManagement.tsx` — страница управления эстафетой

### Database Schema

```sql
-- Распределение матчей между участниками команды
CREATE TABLE relay_event_assignments (
    id SERIAL PRIMARY KEY,
    contest_id BIGINT NOT NULL REFERENCES contests(id),
    team_id BIGINT NOT NULL REFERENCES user_teams(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    event_id BIGINT NOT NULL REFERENCES events(id),
    assigned_by BIGINT NOT NULL REFERENCES users(id),  -- капитан
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(contest_id, team_id, event_id),  -- один матч = один участник в команде
    FOREIGN KEY (team_id, user_id) REFERENCES user_team_members(team_id, user_id)
);

CREATE INDEX idx_relay_assignments_contest ON relay_event_assignments(contest_id);
CREATE INDEX idx_relay_assignments_team ON relay_event_assignments(team_id);
CREATE INDEX idx_relay_assignments_user ON relay_event_assignments(user_id);
```

### Patterns to Follow

**RelayRules структура:**
```go
type RelayRules struct {
    TeamSize    int                  `json:"team_size"`    // участников в команде
    EventCount  int                  `json:"event_count"`  // матчей в конкурсе
    Scoring     StandardScoringRules `json:"scoring"`      // правила подсчёта
    AllowReassign bool               `json:"allow_reassign"` // можно ли переназначать
}
```

**Валидация при прогнозе:**
```go
// Проверить что матч назначен этому участнику в этой команде
func (s *PredictionService) validateRelayAssignment(contestID, teamID, userID, eventID uint) error {
    // SELECT 1 FROM relay_event_assignments WHERE ...
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Backend — Rules & Models

**Задачи:**
- Добавить `ContestTypeRelay` в rules.go
- Создать `RelayRules` структуру
- Создать модель `RelayEventAssignment`
- Создать репозиторий для assignments
- Добавить миграцию БД

### Phase 2: Backend — API

**Задачи:**
- Добавить proto методы для relay assignments
- Реализовать SetRelayAssignments (капитан)
- Реализовать GetTeamAssignments
- Реализовать GetUserAssignments (для участника)
- Добавить валидацию в SubmitPrediction

### Phase 3: Backend — Scoring

**Задачи:**
- Добавить CalculateRelay в calculator
- Агрегировать очки команды в TeamContestEntry
- Обновлять ранг команды после каждого прогноза

### Phase 4: Frontend — Admin

**Задачи:**
- Добавить "Эстафета" в ScoringRulesEditor
- Настройка team_size и event_count
- Выбор матчей (EventSelector)

### Phase 5: Frontend — Captain

**Задачи:**
- Создать RelayAssignmentEditor
- UI распределения матчей (drag & drop или select)
- Показать кто какие матчи прогнозирует
- Валидация (все матчи распределены)

### Phase 6: Frontend — Member

**Задачи:**
- Показывать только назначенные матчи
- Показывать вклад в очки команды
- Показывать лидерборд команд

### Phase 7: Telegram Bot

**Задачи:**
- Показывать командные конкурсы
- Капитан: команда для распределения матчей
- Участник: только свои матчи
- Лидерборд команд

---

## STEP-BY-STEP TASKS

### Task 1: UPDATE backend/shared/scoring/rules.go — Add Relay Type

- **IMPLEMENT**: Добавить `ContestTypeRelay ContestType = "relay"`
- **IMPLEMENT**: Создать структуру:
  ```go
  type RelayRules struct {
      TeamSize      int                  `json:"team_size"`       // 2-10 участников
      EventCount    int                  `json:"event_count"`     // 5-50 матчей
      Scoring       StandardScoringRules `json:"scoring"`
      AllowReassign bool                 `json:"allow_reassign"`  // переназначение до старта
  }
  ```
- **IMPLEMENT**: Добавить поле `Relay *RelayRules` в ContestRules
- **IMPLEMENT**: Обновить ParseRules, Validate, DefaultRelayRules
- **VALIDATE**: `cd backend/shared/scoring && go build .`

### Task 2: CREATE migration — relay_event_assignments table

- **IMPLEMENT**: Создать миграцию в `backend/shared/migrations/`
- **PATTERN**: Следовать паттерну из existing migrations
- **VALIDATE**: `docker exec sports-postgres psql -U sports_user -d sports_prediction -c "\d relay_event_assignments"`

### Task 3: CREATE backend/prediction-service/internal/models/relay_assignment.go

- **IMPLEMENT**: Модель RelayEventAssignment
  ```go
  type RelayEventAssignment struct {
      ID         uint      `gorm:"primaryKey"`
      ContestID  uint      `gorm:"not null;index"`
      TeamID     uint      `gorm:"not null;index"`
      UserID     uint      `gorm:"not null;index"`
      EventID    uint      `gorm:"not null"`
      AssignedBy uint      `gorm:"not null"`
      AssignedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
      gorm.Model
  }
  ```
- **VALIDATE**: `go build ./...`

### Task 4: CREATE backend/prediction-service/internal/repository/relay_repository.go

- **IMPLEMENT**: RelayRepository с методами:
  - SetTeamAssignments(contestID, teamID, assignments []Assignment) error
  - GetTeamAssignments(contestID, teamID) ([]RelayEventAssignment, error)
  - GetUserAssignments(contestID, userID) ([]RelayEventAssignment, error)
  - ValidateUserCanPredict(contestID, teamID, userID, eventID) (bool, error)
- **PATTERN**: Следовать event_repository.go
- **VALIDATE**: `go build ./...`

### Task 5: UPDATE backend/proto/prediction.proto — Relay API

- **IMPLEMENT**: Добавить RPC методы:
  ```protobuf
  rpc SetRelayAssignments(SetRelayAssignmentsRequest) returns (SetRelayAssignmentsResponse);
  rpc GetTeamAssignments(GetTeamAssignmentsRequest) returns (GetTeamAssignmentsResponse);
  rpc GetUserRelayEvents(GetUserRelayEventsRequest) returns (GetUserRelayEventsResponse);
  ```
- **IMPLEMENT**: Добавить message типы
- **VALIDATE**: `./scripts/generate-protos.sh`

### Task 6: UPDATE prediction_service.go — Implement Relay Methods

- **IMPLEMENT**: SetRelayAssignments — только капитан может вызывать
- **IMPLEMENT**: GetTeamAssignments — список распределения
- **IMPLEMENT**: GetUserRelayEvents — матчи участника
- **VALIDATE**: `go build ./...`

### Task 7: UPDATE prediction_service.go — Validate Relay Predictions

- **IMPLEMENT**: В SubmitPrediction проверять:
  - Если contest type = relay
  - Проверить что user принадлежит команде в этом contest
  - Проверить что event назначен этому user
- **GOTCHA**: Не сломать существующую логику для standard/risky/totalizator
- **VALIDATE**: Manual testing

### Task 8: UPDATE scoring/calculator.go — Team Score Aggregation

- **IMPLEMENT**: CalculateRelay (то же что standard)
- **IMPLEMENT**: Функция агрегации очков команды:
  ```go
  func AggregateTeamScore(predictions []Prediction) float64
  ```
- **VALIDATE**: Unit tests

### Task 9: UPDATE frontend ScoringRulesEditor.tsx — Add Relay

- **IMPLEMENT**: Добавить "🏃 Эстафета" в Radio.Group
- **IMPLEMENT**: При выборе relay показать:
  - team_size (2-10, default 5)
  - event_count (5-50, default 15)
  - StandardScoringRules поля
- **VALIDATE**: `npm run build`

### Task 10: CREATE frontend RelayAssignmentEditor.tsx

- **IMPLEMENT**: UI для капитана:
  - Список участников команды (слева)
  - Список матчей (справа)
  - Drag & drop или multi-select для назначения
  - Валидация: все матчи распределены
  - Кнопка "Сохранить распределение"
- **PATTERN**: Использовать Ant Design Transfer или custom drag-drop
- **VALIDATE**: `npm run build`

### Task 11: UPDATE bots/telegram — Relay Support

- **IMPLEMENT**: Показывать тип "🏃 Эстафета"
- **IMPLEMENT**: Для капитана: команда `/assign` для распределения
- **IMPLEMENT**: Для участника: показывать только назначенные матчи
- **IMPLEMENT**: Лидерборд команд
- **VALIDATE**: Manual testing

---

## TESTING STRATEGY

### Unit Tests

- `rules_test.go`: TestParseRelayRules, TestValidateRelayRules
- `calculator_test.go`: TestCalculateRelay, TestAggregateTeamScore
- `relay_repository_test.go`: TestSetTeamAssignments, TestValidateUserCanPredict

### Integration Tests

- Создание relay конкурса
- Регистрация команды
- Капитан распределяет матчи
- Участники делают прогнозы
- Очки команды суммируются

### Edge Cases

- [ ] Команда не полная (меньше team_size)
- [ ] Не все матчи распределены
- [ ] Участник пытается прогнозировать чужой матч
- [ ] Капитан переназначает матч после старта (если allow_reassign=false)
- [ ] Участник выходит из команды

---

## VALIDATION COMMANDS

### Level 1: Build

```bash
cd backend/shared/scoring && go build .
cd backend && go build ./prediction-service/...
cd frontend && npm run build
cd bots/telegram && go build ./...
```

### Level 2: Tests

```bash
cd backend && go test ./shared/scoring/... -v
cd backend && go test ./prediction-service/... -v
```

### Level 3: Integration

```bash
docker-compose up -d
# Test via API / Telegram bot
```

---

## ACCEPTANCE CRITERIA

- [ ] Админ может создать конкурс типа "Эстафета" с настройками команды
- [ ] Команды могут регистрироваться в конкурсе
- [ ] Капитан может распределять матчи между участниками
- [ ] Участники видят и прогнозируют только свои матчи
- [ ] Очки команды суммируются от всех участников
- [ ] Лидерборд показывает команды, не индивидуалов
- [ ] Telegram бот поддерживает relay конкурсы

---

## NOTES

### Архитектурные решения

**Распределение матчей:**
- Капитан назначает каждый матч конкретному участнику
- Один матч = один участник (в рамках команды)
- Все матчи должны быть распределены до первого прогноза

**Суммирование очков:**
- Каждый прогноз считается индивидуально (как в standard)
- TeamContestEntry.TotalPoints = SUM(всех прогнозов участников команды)
- Ранг считается по TotalPoints команды

### Отличие от Тотализатора

| Аспект | Тотализатор | Эстафета |
|--------|-------------|----------|
| Участие | Индивидуальное | Командное |
| Матчи | Все прогнозирует один | Распределены между участниками |
| Очки | Индивидуальные | Суммируются в командные |
| Лидерборд | По участникам | По командам |

### Сложность

Эта фича значительно сложнее Тотализатора из-за:
1. Командной логики
2. Распределения матчей
3. Валидации прав на прогноз
4. Агрегации очков

**Рекомендуется разбить на 2-3 PR:**
1. Backend: rules + models + repository
2. Backend: API + validation
3. Frontend + Telegram bot

---

## CONFIDENCE SCORE

**7/10** — Высокая сложность, много точек интеграции. Требуется тщательное тестирование командной логики.
