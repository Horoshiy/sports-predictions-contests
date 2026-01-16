# Code Review: Dynamic Point Coefficients Implementation

**Date**: 2026-01-16
**Feature**: Dynamic Point Coefficients
**Reviewer**: Kiro CLI

---

## Stats

- Files Modified: 12
- Files Added: 5
- Files Deleted: 0
- New lines: ~216
- Deleted lines: ~23

---

## Issues Found

### Issue 1
```
severity: high
file: backend/prediction-service/internal/service/prediction_service.go
line: 701-765
issue: Code duplication - coefficient functions duplicated across services
detail: The functions `calculateTimeCoefficient` and `getCoefficientTier` are defined in both prediction_service.go and coefficient.go (scoring-service). This violates DRY principle and creates maintenance burden - if the coefficient tiers change, both files must be updated.
suggestion: Import and use the functions from the scoring-service models package, or move to a shared package. Example:
  - Create `backend/shared/coefficient/coefficient.go` with these functions
  - Import from both services
```

### Issue 2
```
severity: medium
file: backend/scoring-service/internal/models/coefficient.go
line: 25-36
issue: Float comparison using equality operator is fragile
detail: Using `switch coefficient { case 2.0: ... }` for float comparison can fail due to floating-point precision issues. While unlikely with these specific values, it's a code smell.
suggestion: Use epsilon comparison or refactor to return both coefficient and tier from a single function:
  func CalculateTimeCoefficientWithTier(submittedAt, eventDate time.Time) (float64, string) {
      hoursUntilEvent := eventDate.Sub(submittedAt).Hours()
      switch {
      case hoursUntilEvent >= 168:
          return 2.0, "Early Bird"
      // ...
      }
  }
```

### Issue 3
```
severity: medium
file: frontend/src/hooks/use-predictions.ts
line: 41
issue: Non-null assertion on potentially undefined eventId
detail: `coefficientKeys.detail(eventId!)` uses non-null assertion even though eventId could be undefined. While `enabled: !!eventId` prevents the query from running, the queryKey is still evaluated.
suggestion: Use a fallback value for the query key:
  queryKey: coefficientKeys.detail(eventId ?? 0),
```

### Issue 4
```
severity: medium
file: frontend/src/components/predictions/EventCard.tsx
line: 46
issue: Unnecessary API calls for non-predictable events
detail: The hook `usePotentialCoefficient(isPredictable ? event.id : undefined)` correctly prevents calls for non-predictable events, but when rendering a list of many events, this still creates many React Query instances. Consider batching or lifting state.
suggestion: For now this is acceptable, but consider adding a batch endpoint `/v1/events/coefficients?ids=1,2,3` if performance becomes an issue with large event lists.
```

### Issue 5
```
severity: low
file: backend/prediction-service/internal/service/prediction_service.go
line: 717
issue: Negative hours not explicitly handled
detail: When `hoursUntilEvent` is negative (event in past), the function returns 1.0 via default case. This is correct behavior but not explicitly documented or handled.
suggestion: Add explicit handling for clarity:
  case hoursUntilEvent < 0:
      return 1.0 // Event already started
  default:
      return 1.0
```

### Issue 6
```
severity: low
file: tests/scoring-service/coefficient_test.go
line: 1
issue: Test file in wrong package
detail: The test file declares `package models` but is located in `tests/scoring-service/` directory. Go tests should be in the same directory as the code they test, or use `_test` suffix package.
suggestion: Move to `backend/scoring-service/internal/models/coefficient_test.go` for proper test discovery and execution.
```

### Issue 7
```
severity: low
file: tests/frontend/src/components/CoefficientIndicator.test.tsx
line: 3
issue: Incorrect import path
detail: Import path `../../src/components/predictions/CoefficientIndicator` is relative to test location but may not resolve correctly depending on test runner configuration.
suggestion: Verify the import path works with the project's vitest configuration. Consider using path aliases like `@/components/predictions/CoefficientIndicator`.
```

### Issue 8
```
severity: low
file: frontend/src/components/predictions/CoefficientIndicator.tsx
line: 34
issue: Potential division issue with very small hours
detail: `Math.floor((hours % 1) * 60)` could produce unexpected results for very small fractional hours due to floating-point precision.
suggestion: Use `Math.round` instead of `Math.floor` for minutes, or handle edge case:
  const minutes = Math.round((hours % 1) * 60)
```

---

## Positive Observations

1. **Good separation of concerns**: The coefficient calculation is properly isolated in its own model file.

2. **Consistent error handling**: Backend follows established patterns for error responses with proper error codes.

3. **Smart query optimization**: Frontend only fetches coefficients for predictable events.

4. **Proper null checks**: Backend correctly handles nil `SubmittedAt` and `EventDate` with fallback to 1.0 coefficient.

5. **Good test coverage**: Unit tests cover main scenarios and edge cases for coefficient calculation.

6. **Type safety**: TypeScript types properly extended for new coefficient fields.

---

## Summary

The implementation is solid overall with no critical security issues. The main concerns are:

1. **Code duplication** (high) - Should be addressed before merge
2. **Float comparison** (medium) - Low risk but worth refactoring
3. **Test file location** (low) - Should be moved for proper Go test discovery

**Recommendation**: Fix the code duplication issue by creating a shared coefficient package, then the code is ready for merge.
