# Code Review: Totalizator Contest Type

**Commit:** `c7bb40e`
**Date:** 2026-02-03
**Reviewer:** Din (AI)

## Summary

Добавлен новый тип конкурса "Тотализатор" — админ выбирает матчи из разных лиг, правила подсчёта очков как в стандартном конкурсе.

## Files Changed

- **28 files**, +1167/-395 lines
- Основные изменения в: `rules.go`, `calculator.go`, `event_repository.go`, `prediction_service.go`, `ScoringRulesEditor.tsx`, `EventSelector.tsx`, `ContestForm.tsx`

---

## ✅ Positive Findings

### 1. Clean Type Extension
```go
const (
    ContestTypeStandard    ContestType = "standard"
    ContestTypeRisky       ContestType = "risky"
    ContestTypeTotalizator ContestType = "totalizator"  // ✅ Clean addition
)
```
Новый тип добавлен без изменения существующей логики.

### 2. Code Reuse
`TotalizatorRules` переиспользует `StandardScoringRules` — нет дублирования логики подсчёта очков.

### 3. Proper Transaction Handling
```go
func (r *EventRepository) SetContestEvents(contestID uint, eventIDs []uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {  // ✅ Atomic
        // DELETE + INSERT in transaction
    })
}
```

### 4. Input Validation
- `event_count` валидируется (5-30)
- Negative points validation
- Frontend показывает min/max selected events

### 5. Frontend UX
- EventSelector с фильтрами (поиск, дата, лига)
- Visual feedback для выбранных матчей
- Disabled checkboxes при достижении лимита

---

## ⚠️ Issues Found

### Issue 1: Code Duplication in Calculator (Medium)

**Location:** `backend/shared/scoring/calculator.go`

**Problem:** `CalculateTotalizator()` почти идентичен `CalculateStandard()` — ~70 строк дублированного кода.

**Recommendation:** Рефакторинг — извлечь общую логику:
```go
func (c *Calculator) calculateWithScoring(prediction, result ScoreData, isAnyOther bool, scoring *StandardScoringRules, contestType string) CalculationResult {
    // Shared logic
}

func (c *Calculator) CalculateStandard(...) CalculationResult {
    return c.calculateWithScoring(prediction, result, isAnyOther, c.rules.Standard, "standard")
}

func (c *Calculator) CalculateTotalizator(...) CalculationResult {
    scoring := &c.rules.Totalizator.Scoring
    return c.calculateWithScoring(prediction, result, isAnyOther, scoring, "totalizator")
}
```

**Severity:** Medium (tech debt, not blocking)

---

### Issue 2: Missing Error Handling in EventSelector (Low)

**Location:** `frontend/src/components/contests/EventSelector.tsx`

**Problem:** Fetch error только логируется, пользователь не видит ошибку.

```tsx
} catch (error) {
    console.error('Failed to fetch events:', error)  // ❌ Silent fail
}
```

**Recommendation:**
```tsx
const [error, setError] = useState<string | null>(null)

} catch (error) {
    setError('Не удалось загрузить матчи')
}

// В render:
{error && <Alert type="error" message={error} />}
```

**Severity:** Low

---

### Issue 3: No Unit Tests (Medium)

**Problem:** Нет тестов для:
- `CalculateTotalizator()`
- `DefaultTotalizatorRules()`
- `SetContestEvents()`, `GetContestEventCount()`

**Recommendation:** Добавить тесты в:
- `backend/shared/scoring/calculator_test.go`
- `backend/shared/scoring/rules_test.go`
- `backend/prediction-service/internal/repository/event_repository_test.go`

**Severity:** Medium

---

### Issue 4: Potential N+1 in AddEventsToContest (Low)

**Location:** `backend/prediction-service/internal/repository/event_repository.go`

**Problem:** Loop с отдельными INSERT:
```go
for _, eventID := range eventIDs {
    err := r.db.Exec("INSERT INTO contest_events ...")  // N queries
}
```

**Recommendation:** Bulk insert:
```go
values := make([]string, len(eventIDs))
for i, id := range eventIDs {
    values[i] = fmt.Sprintf("(%d, %d)", contestID, id)
}
query := fmt.Sprintf("INSERT INTO contest_events (contest_id, event_id) VALUES %s ON CONFLICT DO NOTHING", 
    strings.Join(values, ","))
```

**Severity:** Low (usually < 30 events)

---

### Issue 5: Missing Proto Generated Files in shared/proto (Info)

**Problem:** Proto файлы генерируются в `backend/proto/` но код импортирует из `backend/shared/proto/`. Требуется ручное копирование.

**Recommendation:** Исправить `scripts/generate-protos.sh` или пути импорта.

**Severity:** Info (workaround exists)

---

## 🔒 Security Check

- ✅ No hardcoded secrets
- ✅ SQL injection safe (parameterized queries)
- ✅ Input validation present
- ✅ No sensitive data exposure

---

## 📊 Summary

| Category | Status |
|----------|--------|
| Logic Errors | ✅ None found |
| Security Issues | ✅ None found |
| Performance | ⚠️ Minor (N+1 in batch insert) |
| Code Quality | ⚠️ Duplication in calculator |
| Tests | ❌ Missing |

**Overall:** Code is functional and safe. Recommended to address code duplication and add tests before production.

---

## Action Items

1. **[Optional]** Refactor calculator to reduce duplication
2. **[Low]** Add error state to EventSelector
3. **[Medium]** Add unit tests for totalizator
4. **[Low]** Optimize bulk insert

**Approval:** ✅ Ready to deploy (with noted improvements for follow-up)
