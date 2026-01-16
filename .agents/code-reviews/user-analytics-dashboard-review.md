# Code Review: User Analytics Dashboard Implementation

**Date**: 2026-01-16
**Feature**: User Analytics Dashboard
**Reviewer**: Kiro CLI

## Stats

- Files Modified: 4
- Files Added: 12
- Files Deleted: 0
- New lines: ~1,100
- Deleted lines: 2

---

## Issues Found

### HIGH

```
severity: high
file: backend/scoring-service/internal/repository/analytics_repository.go
line: 213
issue: SQL query parameter binding issue with TO_CHAR
detail: The GORM Select() method with positional parameters for TO_CHAR may not bind correctly. The dateFormat variable is passed as a parameter but GORM may interpret it as a column name rather than a string literal for the format.
suggestion: Use raw SQL with proper escaping or construct the query string directly:
  Select(fmt.Sprintf("TO_CHAR(scored_at, '%s') as period, ...", dateFormat))
  Or use GORM's Raw() method for complex queries.
```

```
severity: high
file: backend/scoring-service/internal/repository/analytics_repository.go
line: 216-217
issue: GROUP BY clause has same parameter binding issue
detail: The Group() clause also uses the dateFormat parameter which may not bind correctly, causing SQL syntax errors at runtime.
suggestion: Same fix as above - construct the GROUP BY clause with the format string embedded:
  Group(fmt.Sprintf("TO_CHAR(scored_at, '%s')", dateFormat))
```

```
severity: high
file: frontend/src/components/analytics/PlatformComparison.tsx
line: 14
issue: Unused import - UserAnalytics type imported but never used
detail: The UserAnalytics type is imported in the type declaration but the interface only uses PlatformStats. This will cause TypeScript unused import warnings.
suggestion: Remove the unused import:
  import type { PlatformStats } from '../../types/analytics.types'
```

### MEDIUM

```
severity: medium
file: backend/scoring-service/internal/repository/analytics_repository.go
line: 119
issue: Incorrect JOIN for league analytics - events.id != matches.id
detail: The query joins "LEFT JOIN matches m ON e.id = m.id" which assumes events and matches share the same ID. This is likely incorrect - events should reference matches via a foreign key, not by ID equality.
suggestion: Review the data model. If events have a match_id column, use:
  Joins("LEFT JOIN matches m ON e.match_id = m.id")
  Or if matches are events, the join should be on the correct relationship.
```

```
severity: medium
file: backend/scoring-service/internal/repository/analytics_repository.go
line: 75-76
issue: Time filter applied after GROUP BY in query builder
detail: When !since.IsZero(), the WHERE clause is added after the GROUP BY is already set up. While GORM may handle this correctly, the query order could cause issues with some database drivers.
suggestion: Restructure to add all WHERE clauses before GROUP BY:
  if !since.IsZero() {
    query = query.Where("s.scored_at >= ?", since)
  }
  query = query.Group("e.sport_type")
```

```
severity: medium
file: frontend/src/services/analytics-service.ts
line: 21
issue: No error handling for failed API response
detail: The getUserAnalytics method returns response.analytics directly without checking if response.response.success is true. If the API returns an error response, analytics may be undefined.
suggestion: Add response validation:
  if (!response.response?.success) {
    throw new Error(response.response?.message || 'Failed to fetch analytics')
  }
  return response.analytics
```

```
severity: medium
file: frontend/src/services/analytics-service.ts
line: 35
issue: Same missing error handling in exportAnalytics
detail: exportAnalytics also doesn't validate the response success status before returning data.
suggestion: Add the same response validation as suggested above.
```

```
severity: medium
file: frontend/src/pages/AnalyticsPage.tsx
line: 63
issue: Query runs with userId=0 when user is null
detail: useUserAnalytics is called with user?.id || 0, and while enabled: !!userId prevents the query from running, passing 0 as a fallback is misleading and could cause issues if the enabled check is removed.
suggestion: Use a more explicit pattern:
  const userId = user?.id
  const { data, isLoading, error } = useUserAnalytics(userId ?? 0, timeRange)
  Or handle the null case before the hook call.
```

### LOW

```
severity: low
file: backend/scoring-service/internal/service/scoring_service.go
line: 505
issue: CSV export doesn't validate time_range parameter
detail: ExportAnalytics passes req.TimeRange directly to GetUserAnalytics without validation. If TimeRange is empty, the filename will be "analytics_1_.csv" which looks odd.
suggestion: Default the time range before generating filename:
  timeRange := req.TimeRange
  if timeRange == "" {
    timeRange = "30d"
  }
  filename := fmt.Sprintf("analytics_%d_%s.csv", req.UserId, timeRange)
```

```
severity: low
file: frontend/src/components/analytics/AccuracyChart.tsx
line: 60
issue: Tooltip formatter has incorrect name matching
detail: The formatter checks for name === 'accuracy' but the Line component uses name="Accuracy %" which won't match.
suggestion: Match the actual name or use dataKey for comparison:
  if (name === 'Accuracy %') return [`${value}%`, 'Accuracy']
```

```
severity: low
file: frontend/src/components/analytics/SportBreakdown.tsx
line: 56
issue: Same tooltip formatter name mismatch
detail: Checks for name === 'accuracy' but Bar uses name="Accuracy %"
suggestion: Update to match: if (name === 'Accuracy %')
```

```
severity: low
file: tests/scoring-service/analytics_test.go
line: 1
issue: Test package name doesn't match directory
detail: Package is "scoring_service" but file is in tests/scoring-service/. Go convention is for test packages to match the directory name or use _test suffix.
suggestion: Either rename package to "scoring_service_test" for black-box testing or move tests to the actual package directory.
```

---

## Security Analysis

No critical security issues found. The implementation:
- ✅ Uses parameterized queries (GORM) preventing SQL injection
- ✅ Validates user ID is non-zero before processing
- ✅ Uses soft delete filters (deleted_at IS NULL)
- ✅ No sensitive data exposed in CSV export

---

## Performance Analysis

**Potential Concerns:**
1. **Multiple sequential queries**: GetUserAnalytics makes 6 separate database queries. Consider combining into fewer queries or using parallel execution.
2. **No pagination on trends**: GetAccuracyTrends returns all periods without limit, which could be large for "all" time range.
3. **Platform stats query scans entire scores table**: GetPlatformStats aggregates all scores which could be slow on large datasets.

**Recommendations:**
- Add caching for platform stats (changes infrequently)
- Add LIMIT to trends query (e.g., last 30 periods)
- Consider materialized views for analytics if performance becomes an issue

---

## Summary

The implementation is generally well-structured and follows existing codebase patterns. The main issues are:

1. **HIGH**: SQL parameter binding for TO_CHAR in trends query needs fixing
2. **HIGH**: Unused TypeScript import
3. **MEDIUM**: Incorrect JOIN relationship for league analytics
4. **MEDIUM**: Missing API response validation in frontend service

Recommend fixing HIGH and MEDIUM issues before merging.
