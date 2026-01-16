# Code Review: Predictions UI Implementation

**Date**: 2026-01-16
**Reviewer**: Kiro CLI
**Feature**: Predictions UI - Frontend for submitting and managing sports predictions

## Stats

- Files Modified: 1
- Files Added: 9
- Files Deleted: 0
- New lines: ~650
- Deleted lines: 0

---

## Issues Found

### HIGH

```
severity: high
file: frontend/src/pages/PredictionsPage.tsx
line: 49-60
issue: Placeholder event data when editing prediction loses real event information
detail: When editing a prediction from PredictionList, handleEditPrediction creates a fake Event object with placeholder data ("Home", "Away", empty strings). This means the PredictionForm displays incorrect team names and event details, confusing users.
suggestion: Fetch the actual event data using useEvent hook or pass events data from EventList to PredictionList, or store events in a map when loading predictions.
```

```
severity: high
file: frontend/src/pages/PredictionsPage.tsx
line: 28
issue: Hardcoded pagination limit of 100 may cause performance issues
detail: Fetching up to 100 predictions at once without pagination could cause slow load times and memory issues for users with many predictions. The limit should match the UI pagination or use infinite scroll.
suggestion: Either implement proper pagination in the predictions fetch or reduce the limit to match typical page sizes (10-20).
```

### MEDIUM

```
severity: medium
file: frontend/src/components/predictions/EventList.tsx
line: 22-23
issue: Hardcoded sport types don't match backend dynamically
detail: The sportTypes array is hardcoded and may not match the actual sports available in the backend. This could show filters for non-existent sports or miss new ones.
suggestion: Fetch sport types from the sports-service API or use the existing useSports hook to get available sports dynamically.
```

```
severity: medium
file: frontend/src/hooks/use-predictions.ts
line: 29-34
issue: Query key doesn't include pagination, causing stale data
detail: useUserPredictions uses only contestId in the query key but the request includes pagination. Changing pages won't trigger a refetch because React Query sees the same key.
suggestion: Include pagination in the query key: `queryKey: predictionKeys.list(request.contestId, request.pagination)`
```

```
severity: medium
file: frontend/src/components/predictions/PredictionForm.tsx
line: 140-141
issue: Score input converts empty string to 0, preventing clearing
detail: The onChange handler `parseInt(e.target.value) || 0` converts empty input to 0, making it impossible for users to clear a score field. Users might want to switch from "combined" to "winner" type.
suggestion: Use `parseInt(e.target.value)` without the `|| 0` fallback, and handle NaN in validation instead.
```

```
severity: medium
file: frontend/src/utils/prediction-validation.ts
line: 21-24
issue: Validation refinement doesn't check for score type when combined
detail: When predictionType is 'combined', the second refinement only checks if scores exist but doesn't validate that both are numbers >= 0. The schema allows undefined scores to pass through.
suggestion: Add explicit number validation in the refinement or ensure the base schema handles this case.
```

### LOW

```
severity: low
file: frontend/src/components/predictions/PredictionList.tsx
line: 65
issue: handleDelete uses window.confirm which is not accessible
detail: Using window.confirm() for delete confirmation is not accessible and doesn't match the Material-UI design system used elsewhere. Other parts of the codebase may use MUI Dialog for confirmations.
suggestion: Consider using a MUI Dialog or Snackbar with undo action for better UX consistency.
```

```
severity: low
file: frontend/src/components/predictions/EventCard.tsx
line: 29
issue: Date comparison doesn't account for timezone differences
detail: `new Date(event.eventDate) > new Date()` compares dates without considering that eventDate might be in a different timezone than the user's local time.
suggestion: Use a date library like date-fns with proper timezone handling, or ensure backend sends UTC timestamps consistently.
```

```
severity: low
file: frontend/src/services/prediction-service.ts
line: 1-15
issue: Unused imports - Event and Prediction types imported but only used in return types
detail: The Prediction and Event types are imported but TypeScript can infer them from the response types. This is minor but adds unnecessary imports.
suggestion: Remove unused direct imports if they're only needed for type inference from response types.
```

---

## Summary

The Predictions UI implementation follows established patterns from the codebase (contest-service, use-contests hooks) and provides a functional interface for managing predictions. The code is well-structured with proper TypeScript types, React Query integration, and Material-UI components.

**Key concerns to address before merge:**
1. **HIGH**: Fix the placeholder event data issue in handleEditPrediction - users will see wrong team names
2. **HIGH**: Address the hardcoded pagination limit of 100 predictions
3. **MEDIUM**: Fix query key to include pagination parameters

The remaining medium/low issues are improvements that can be addressed in follow-up PRs.
