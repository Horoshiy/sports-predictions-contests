# Code Review: Dynamic Point Coefficients - Post-Fix Review

**Date**: 2026-01-16
**Feature**: Dynamic Point Coefficients (after fixes)
**Reviewer**: Kiro CLI

---

## Stats

- Files Modified: 12
- Files Added: 8
- Files Deleted: 0
- New lines: ~183
- Deleted lines: ~23

---

## Issues Found

### Issue 1
```
severity: medium
file: backend/scoring-service/internal/models/coefficient.go
line: 11-18
issue: Double calculation when getting both coefficient and tier
detail: The wrapper functions `CalculateTimeCoefficient` and `GetCoefficientTier` each call `coefficient.Calculate()` separately. If both are needed, this results in duplicate computation.
suggestion: Either use the shared package directly where both values are needed, or add a combined wrapper:
  func CalculateWithTier(submittedAt, eventDate time.Time) coefficient.CoefficientResult {
      return coefficient.Calculate(submittedAt, eventDate)
  }
```

### Issue 2
```
severity: low
file: frontend/src/services/prediction-service.ts
line: 169-180
issue: Missing error handling for API response
detail: The `getPotentialCoefficient` method doesn't check `response.response?.success` before returning data, unlike other methods in the service (e.g., `getPropTypes` on line 108).
suggestion: Add response validation:
  if (!response.response?.success) {
    throw new Error(response.response?.message || 'Failed to get coefficient')
  }
```

### Issue 3
```
severity: low
file: tests/frontend/src/components/CoefficientIndicator.test.tsx
line: 3
issue: Fragile import path with multiple parent directory traversals
detail: Import path `../../../../frontend/src/components/predictions/CoefficientIndicator` uses 4 levels of parent traversal which is fragile and hard to maintain.
suggestion: Move test to `frontend/src/components/predictions/__tests__/CoefficientIndicator.test.tsx` or configure path aliases in vitest config.
```

---

## Positive Observations

1. **DRY principle followed**: Shared coefficient package eliminates code duplication between services.

2. **Explicit negative hours handling**: The shared package explicitly handles `hoursUntilEvent < 0` case.

3. **Combined result struct**: `CoefficientResult` returns both coefficient and tier together, avoiding float comparison issues.

4. **Good test coverage**: Both backend packages have comprehensive tests covering edge cases.

5. **Safe query key handling**: Frontend hook uses `eventId ?? 0` fallback for query key.

6. **Proper null checks**: Backend correctly handles nil `SubmittedAt` and `EventDate`.

---

## Summary

The code is in good shape after the previous fixes. The remaining issues are minor:
- One medium issue about potential double calculation (optimization)
- Two low issues about error handling and test organization

**Recommendation**: Code is ready for merge. The medium issue can be addressed in a follow-up if performance profiling shows it's needed.

---

## Code Review Passed

No critical or high severity issues detected. Ready for commit.
