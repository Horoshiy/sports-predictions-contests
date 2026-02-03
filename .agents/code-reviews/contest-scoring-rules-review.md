# Code Review: Contest Scoring Rules Implementation

**Date:** 2026-02-03
**Commits:** `758cb35`, `a2914dd`
**Reviewer:** Дин (AI)

## Stats

- **Files Modified:** 4
- **Files Added:** 6
- **New lines:** ~1700
- **Deleted lines:** ~30

---

## Issues Found

### HIGH: Memory Leak in Bot Risky Selections

**severity:** high
**file:** bots/telegram/bot/handlers.go
**line:** 571
**issue:** Global map `riskySelections` never cleaned up
**detail:** The `riskySelections` map stores user selections but entries are only deleted on successful submit. If user abandons the flow, the entry remains forever, causing memory leak over time.
**suggestion:** Add TTL-based cleanup or store in session instead:
```go
// Option 1: Store in session
type UserSession struct {
    // ...existing fields...
    RiskySelections map[uint32][]string // matchID -> selections
}

// Option 2: Add cleanup in session cleanup goroutine
```

---

### MEDIUM: Missing Input Validation in Calculator

**severity:** medium
**file:** backend/shared/scoring/calculator.go
**line:** 31
**issue:** No nil check for rules pointer
**detail:** `CalculateStandard` accesses `c.rules.Standard` but `c.rules` could be nil if NewCalculator was called with nil.
**suggestion:** Add nil check:
```go
func (c *Calculator) CalculateStandard(...) CalculationResult {
    if c.rules == nil {
        return CalculationResult{Points: 0, Details: map[string]interface{}{"error": "nil rules"}}
    }
    // ...
}
```

---

### MEDIUM: Duplicate Type Definitions

**severity:** medium
**file:** bots/telegram/bot/risky_predictions.go
**line:** 12-28
**issue:** RiskyEvent, RiskyScoringRules, ContestRules duplicated from shared/scoring
**detail:** Same types defined in both packages. This creates maintenance burden and potential inconsistency.
**suggestion:** Import from shared package or create a shared types package that both can use. For now, acceptable since bot is separate module, but document the duplication.

---

### MEDIUM: Hardcoded Contest Rules in Bot

**severity:** medium
**file:** bots/telegram/bot/handlers.go
**line:** 580
**issue:** `rulesJSON := ""` hardcoded instead of fetching from contest
**detail:** handleRiskyToggle uses empty rulesJSON, so always defaults. Should fetch actual rules from contest-service.
**suggestion:** Add gRPC call to get contest rules:
```go
contestResp, err := h.clients.Contest.GetContest(ctx, &contestpb.GetContestRequest{Id: contestID})
if err == nil && contestResp.Contest != nil {
    rulesJSON = contestResp.Contest.Rules
}
```

---

### LOW: Missing Unit Tests

**severity:** low
**file:** backend/shared/scoring/
**issue:** No test files for rules.go and calculator.go
**detail:** Critical scoring logic has no automated tests.
**suggestion:** Add `rules_test.go` and `calculator_test.go` with test cases:
- ParseRules with valid/invalid JSON
- CalculateStandard for all match types
- CalculateRisky with various outcomes

---

### LOW: Callback Data Length Risk

**severity:** low
**file:** bots/telegram/bot/risky_predictions.go
**line:** 95
**issue:** Callback data could exceed Telegram's 64-byte limit
**detail:** `risky_{matchID}_{slug}` with long slugs could exceed limit.
**suggestion:** Use shorter format or hash: `r_{matchID}_{slugIndex}`

---

### LOW: Frontend Missing Error Boundary

**severity:** low
**file:** frontend/src/components/contests/ScoringRulesEditor.tsx
**line:** 72
**issue:** JSON.parse can throw but no try-catch in ContestForm integration
**detail:** If rules JSON is malformed, parse will throw and crash component.
**suggestion:** Already handled with try-catch in ContestForm.tsx ✓ (false positive)

---

## Positive Observations

1. ✅ Good separation of concerns (rules parsing, calculation, UI)
2. ✅ Backward compatibility maintained (default rules when empty)
3. ✅ Validation logic in rules.go is thorough
4. ✅ Frontend components follow existing patterns (Ant Design)
5. ✅ Bilingual support (Russian + English) in events

---

## Recommendations

### Critical (Before Deploy)
1. Fix memory leak in `riskySelections` map

### Important (Soon)
2. Add unit tests for scoring package
3. Fetch actual contest rules in bot handlers

### Nice to Have
4. Consolidate type definitions
5. Add E2E test for risky prediction flow

---

## Verdict

**Status:** ⚠️ Conditional Pass

Code is well-structured and follows patterns, but the memory leak issue should be addressed before production use. The missing test coverage is acceptable for MVP but should be added soon.

**Confidence for production:** 7/10 (after fixing HIGH issue: 9/10)
