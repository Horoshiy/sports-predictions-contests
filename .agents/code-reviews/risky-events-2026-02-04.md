# Code Review: Risky Events Feature

**Date:** 2026-02-04  
**Reviewer:** Дин  
**Commits:** ca809e8, be28a8c, e406918, 1aa9c4d  
**Files Reviewed:** 16 files (frontend + telegram bot)

---

## Summary

Overall code quality is **GOOD**. The implementation follows existing patterns and is well-structured. Found **3 issues** requiring attention.

---

## Issues Found

### 🔴 HIGH: Cache Memory Leak (Telegram Bot)

**File:** `bots/telegram/bot/risky_predictions.go:55-63`

**Problem:** The `matchEvents` map in `RiskyEventsCache` grows indefinitely. Old entries are never cleaned up.

```go
type RiskyEventsCache struct {
    matchEvents    map[string]matchEventsCacheEntry // key: "eventId:contestId"
    // No cleanup mechanism!
}
```

**Impact:** Memory will grow unbounded over time as new matches are queried.

**Fix Required:** Add periodic cache cleanup or use LRU cache.

```go
// Option 1: Add cleanup goroutine
func (c *RiskyEventsCache) startCleanup() {
    ticker := time.NewTicker(10 * time.Minute)
    go func() {
        for range ticker.C {
            c.mu.Lock()
            now := time.Now()
            for key, entry := range c.matchEvents {
                if now.After(entry.expiry) {
                    delete(c.matchEvents, key)
                }
            }
            c.mu.Unlock()
        }
    }()
}

// Option 2: Use github.com/hashicorp/golang-lru
```

---

### 🟡 MEDIUM: Missing Error Handling in Frontend

**File:** `frontend/src/services/risky-events-service.ts:73-78`

**Problem:** `getMatchEvents` doesn't handle API errors gracefully - returns empty array which may hide errors.

```typescript
async getMatchEvents(request: GetMatchRiskyEventsRequest): Promise<GetMatchRiskyEventsResponse> {
    const response = await grpcClient.get<GetMatchRiskyEventsResponse>(...)
    return {
      events: response.events || [],  // Silent fallback
      maxSelections: response.maxSelections || 5,
    }
}
```

**Impact:** If API returns malformed response, errors are silently swallowed.

**Fix:** Log warning or throw when response is unexpected.

---

### 🟡 MEDIUM: Unused Import in MatchRiskyEventsEditor

**File:** `frontend/src/components/events/MatchRiskyEventsEditor.tsx:23`

**Problem:** `RISKY_EVENT_CATEGORIES` is imported but never used.

```typescript
import { RISKY_EVENT_CATEGORIES } from '../../types/risky-events.types'
// Never used in the component
```

**Impact:** Minor - dead code, bundle size slightly larger.

**Fix:** Remove unused import.

---

## Positive Observations ✅

### Good Practices Found:

1. **Proper cache invalidation** in React Query hooks after mutations
2. **Fallback mechanism** in bot when API fails
3. **Type-safe API** with proper TypeScript interfaces
4. **Separation of concerns** - service, hooks, components properly separated
5. **Consistent snake_case conversion** for API payloads
6. **Timeouts on gRPC calls** (5 seconds) in bot

### Code Quality:

- Follows existing patterns in codebase
- Good error messages in Russian for toasts
- Proper use of React Query for caching
- Clean component structure with Ant Design

---

## Test Coverage

- ✅ TypeScript compilation: PASS
- ✅ Go vet: PASS
- ⚠️ Unit tests: Not run (no test files for new code)
- ⚠️ E2E tests: Need verification

---

## Recommendations

1. **Add unit tests** for:
   - `risky-events-service.ts`
   - `use-risky-events.ts` hooks
   - `risky_predictions.go` cache logic

2. **Add E2E test** for:
   - Admin creating/editing risky event types
   - Contest with risky events flow

3. **Consider** adding OpenTelemetry tracing for cache hits/misses in bot

---

## Action Items

| Priority | Issue | Fix | Status |
|----------|-------|-----|--------|
| 🔴 HIGH | Cache memory leak | Add cleanup goroutine | ✅ Fixed (4fa380c) |
| 🟡 MEDIUM | Error handling | Add logging/validation | ✅ Fixed (4fa380c) |
| 🟡 MEDIUM | Unused import | Remove import | ✅ Fixed (4fa380c) |

---

## Approval Status

**✅ APPROVED** - All issues fixed in commit 4fa380c.
