# Code Review: Props Predictions Feature

**Date**: 2026-01-16
**Feature**: Props Predictions (Statistics-Based Predictions)
**Reviewer**: Kiro CLI

## Stats

- Files Modified: 9
- Files Added: 6
- Files Deleted: 0
- New lines: ~500
- Deleted lines: ~10

---

## Issues Found

### HIGH

```
severity: high
file: frontend/src/components/predictions/PredictionForm.tsx
line: 91-95
issue: Props prediction can be submitted with empty selections
detail: When predictionType is 'props', the form calls propsFormDataToPredictionData(selectedProps) but there's no validation that selectedProps has items with valid selections. A user can add props but leave selection empty and submit.
suggestion: Add validation before submit:
  if (data.predictionType === 'props') {
    const invalidProps = selectedProps.filter(p => !p.selection)
    if (selectedProps.length === 0 || invalidProps.length > 0) {
      // Show error toast or set form error
      return
    }
    predictionData = propsFormDataToPredictionData(selectedProps)
  }
```

```
severity: high
file: backend/scoring-service/internal/service/scoring_service.go
line: 703-710
issue: evaluateProp returns false for unknown prop slugs without logging
detail: When a prop slug doesn't match any known case, the function silently returns false. This makes debugging difficult when new prop types are added but scoring logic isn't updated.
suggestion: Add logging for unknown prop slugs:
  default:
    log.Printf("[WARN] Unknown prop slug: %s", prop.PropSlug)
    return false
```

```
severity: high
file: frontend/src/components/predictions/PropTypeSelector.tsx
line: 97
issue: Using array index as React key in list
detail: Using `key={index}` for the Card component can cause rendering issues when props are reordered or removed. React may not properly update the DOM.
suggestion: Use a stable unique identifier:
  <Card key={`${prop.propTypeId}-${index}`} variant="outlined">
Or better, use propTypeId alone since each prop type can only be selected once:
  <Card key={prop.propTypeId} variant="outlined">
```

### MEDIUM

```
severity: medium
file: backend/prediction-service/internal/repository/prop_type_repository.go
line: 18-22
issue: Redundant deleted_at check with GORM soft delete
detail: GORM automatically adds `deleted_at IS NULL` when using soft delete (gorm.Model includes DeletedAt). The explicit check is redundant and could cause issues if GORM behavior changes.
suggestion: Remove explicit deleted_at check, rely on GORM's soft delete:
  err := r.db.WithContext(ctx).
    Where("sport_type = ? AND is_active = ?", sportType, true).
    Order("category, name").
    Find(&propTypes).Error
```

```
severity: medium
file: backend/scoring-service/internal/service/scoring_service.go
line: 697-700
issue: Type assertion without handling integer type from JSON
detail: When JSON is unmarshaled, numbers can be float64 or int depending on the source. The code only checks for float64, but `result.Stats["corners"]` could be an int if set programmatically.
suggestion: Handle both types:
  case "total-corners-ou":
    var corners float64
    switch v := result.Stats["corners"].(type) {
    case float64:
      corners = v
    case int:
      corners = float64(v)
    default:
      return false
    }
    if prop.Selection == "over" {
      return corners > prop.Line
    }
    return corners < prop.Line
```

```
severity: medium
file: frontend/src/utils/prediction-validation.ts
line: 57-70
issue: propsFormDataToPredictionData doesn't validate props array
detail: The function accepts any array and converts it to JSON without validation. Empty arrays or props with missing required fields will be serialized.
suggestion: Add validation or use the Zod schema:
  export const propsFormDataToPredictionData = (props: PropPredictionFormData[]): string => {
    if (props.length === 0) {
      throw new Error('At least one prop is required')
    }
    const invalidProps = props.filter(p => !p.selection || !p.propSlug)
    if (invalidProps.length > 0) {
      throw new Error('All props must have a selection')
    }
    return JSON.stringify({...})
  }
```

```
severity: medium
file: backend/prediction-service/internal/models/prop_type.go
line: 64-73
issue: BeforeCreate doesn't validate Slug field
detail: The Slug field is required (uniqueIndex) but not validated in BeforeCreate. An empty slug would fail at the database level with a less helpful error.
suggestion: Add ValidateSlug method and call it in BeforeCreate:
  func (p *PropType) ValidateSlug() error {
    if strings.TrimSpace(p.Slug) == "" {
      return errors.New("slug cannot be empty")
    }
    return nil
  }
```

### LOW

```
severity: low
file: frontend/src/types/props.types.ts
line: 17-24
issue: PropPrediction interface has line as required but it's optional for some prop types
detail: The `line` field is required in the interface but yes_no and team_select props don't use it. This could cause confusion.
suggestion: Make line optional to match actual usage:
  export interface PropPrediction {
    propTypeId: number
    propSlug: string
    line?: number  // Optional - only used for over_under types
    selection: string
    playerId?: string
    pointsValue: number
  }
```

```
severity: low
file: backend/proto/prediction.proto
line: 56-60
issue: PropPrediction message defined but not used in any RPC
detail: The PropPrediction message is defined but predictions still use string prediction_data. This message could be used for stronger typing.
suggestion: Consider using PropPrediction in a future iteration for type-safe props predictions, or remove if not planned.
```

```
severity: low
file: frontend/src/components/predictions/PropTypeSelector.tsx
line: 79-82
issue: groupedPropTypes computed on every render
detail: The reduce operation to group prop types runs on every render even when propTypes hasn't changed.
suggestion: Wrap in useMemo:
  const groupedPropTypes = React.useMemo(() => 
    availablePropTypes.reduce((acc, pt) => {
      if (!acc[pt.category]) acc[pt.category] = []
      acc[pt.category].push(pt)
      return acc
    }, {} as Record<string, PropType[]>),
    [availablePropTypes]
  )
```

---

## Security Analysis

No critical security issues found. The implementation:
- Uses parameterized queries via GORM (no SQL injection)
- Validates input on backend before database operations
- Properly encodes URL parameters in frontend

---

## Summary

The Props Predictions feature is well-implemented overall with good separation of concerns. The main issues to address before commit:

1. **HIGH**: Add validation for empty props selections before form submission
2. **HIGH**: Add logging for unknown prop slugs in scoring
3. **HIGH**: Fix React key usage in PropTypeSelector

The medium and low issues are improvements that can be addressed in a follow-up commit.
