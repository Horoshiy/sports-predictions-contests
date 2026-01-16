# Code Review: Frontend Sports Management UI

**Date**: 2026-01-16
**Reviewer**: Kiro CLI
**Scope**: New Sports Management UI feature (13 new files, 1 modified file)

---

## Stats

- **Files Modified**: 1
- **Files Added**: 13
- **Files Deleted**: 0
- **New lines**: ~1,768
- **Deleted lines**: 3

---

## Issues Found

### Issue 1
```
severity: high
file: frontend/src/components/sports/SportForm.tsx
line: 35-38
issue: Slug auto-generation creates infinite loop risk
detail: The useEffect watches both `name` and `slug`, and calls setValue('slug'). When slug is empty and name changes, it sets slug, which triggers the effect again. While the condition `!slug` prevents infinite loop, it also prevents manual slug editing from being preserved if user clears it.
suggestion: Track whether slug was manually edited with a separate ref:
  const slugManuallyEdited = useRef(false);
  // In slug field onChange: slugManuallyEdited.current = true
  // In useEffect: if (!isEditing && name && !slugManuallyEdited.current)
```

### Issue 2
```
severity: high
file: frontend/src/components/sports/MatchForm.tsx
line: 48-51
issue: Team selection not reset when league changes
detail: When user selects a league, then selects teams, then changes to a different league with different sport, the previously selected teams remain selected even though they may not be valid for the new league's sport. This causes validation to pass with invalid team/league combinations.
suggestion: Add useEffect to reset team selections when leagueId changes:
  React.useEffect(() => {
    if (!isEditing) {
      setValue('homeTeamId', 0);
      setValue('awayTeamId', 0);
    }
  }, [selectedLeagueId, isEditing, setValue]);
```

### Issue 3
```
severity: medium
file: frontend/src/utils/sports-validation.ts
line: 7
issue: Slug regex rejects valid empty string before .optional() check
detail: The regex validation runs before the .optional().or(z.literal('')) check, causing validation errors when slug is empty string. The regex /^[a-z0-9]+(?:-[a-z0-9]+)*$/ requires at least one character.
suggestion: Reorder validation to check empty first:
  slug: z.string().optional().or(z.literal('')).refine(
    (val) => !val || slugRegex.test(val),
    'Slug must be lowercase alphanumeric with hyphens'
  ),
```

### Issue 4
```
severity: medium
file: frontend/src/pages/SportsPage.tsx
line: 56-68
issue: Missing error handling for async form submissions
detail: The handleSportSubmit, handleLeagueSubmit, handleTeamSubmit, and handleMatchSubmit functions use mutateAsync but don't have try/catch. If mutation fails, closeForm() is never called and form stays open with stale state.
suggestion: Wrap in try/catch or use mutate() instead of mutateAsync() since toast notifications already handle errors:
  const handleSportSubmit = async (data: SportFormData) => {
    try {
      const slug = data.slug || generateSlug(data.name)
      if (selectedSport) {
        await updateSport.mutateAsync({ ... })
      } else {
        await createSport.mutateAsync({ ... })
      }
      closeForm()
    } catch {
      // Error already shown via toast in hook
    }
  }
```

### Issue 5
```
severity: medium
file: frontend/src/components/sports/MatchList.tsx
line: 28
issue: Fetching 500 teams on every render is inefficient
detail: useTeams({ pagination: { page: 1, limit: 500 } }) fetches up to 500 teams just to build a lookup map. This is wasteful for large datasets and may cause performance issues.
suggestion: Consider fetching teams only when needed, or use a more efficient approach like fetching team names on-demand or including team names in the match response from backend.
```

### Issue 6
```
severity: medium
file: frontend/src/services/sports-service.ts
line: 37-44
issue: Missing null safety for pagination response
detail: The listSports method returns response.pagination directly, but if the API returns null/undefined pagination, this will cause runtime errors in components that access pagination.total.
suggestion: Add default pagination:
  return { 
    sports: response.sports || [], 
    pagination: response.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 } 
  }
```

### Issue 7
```
severity: low
file: frontend/src/components/sports/LeagueForm.tsx
line: 58
issue: Missing error message display for sportId field
detail: Unlike other form fields, the sportId FormControl doesn't display the error message from Zod validation. User won't see "Sport is required" message.
suggestion: Add helperText or error message display:
  <FormControl fullWidth required error={!!errors.sportId}>
    <InputLabel>Sport</InputLabel>
    <Select {...field} label="Sport" disabled={loading}>
      {sportsData?.sports?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}
    </Select>
    {errors.sportId && <FormHelperText>{errors.sportId.message}</FormHelperText>}
  </FormControl>
```

### Issue 8
```
severity: low
file: frontend/src/components/sports/SportList.tsx
line: 22-28
issue: URL sync effect runs on every searchParams change
detail: The useEffect that syncs pagination to URL has searchParams in its dependency array, which causes it to run whenever URL changes (including from this effect itself). This creates unnecessary URL updates.
suggestion: Remove searchParams from dependency array or use a ref to track if update came from this component:
  useEffect(() => {
    const params = new URLSearchParams()
    params.set('page', pagination.pageIndex.toString())
    params.set('limit', pagination.pageSize.toString())
    setSearchParams(params, { replace: true })
  }, [pagination.pageIndex, pagination.pageSize, setSearchParams])
```

### Issue 9
```
severity: low
file: frontend/src/types/sports.types.ts
line: 1
issue: Unused import PaginationRequest
detail: PaginationRequest is imported but only used in request interfaces that are passed to service methods. The service methods build URLSearchParams manually, so this type is only used for type checking, which is fine, but could be simplified.
suggestion: This is acceptable as-is for type safety. No change needed.
```

---

## Security Analysis

✅ No hardcoded credentials or API keys
✅ No SQL injection vulnerabilities (using gRPC/REST API)
✅ No XSS vulnerabilities (React handles escaping)
✅ Authentication properly enforced via ProtectedRoute
✅ Delete operations require user confirmation

---

## Performance Analysis

⚠️ **Medium concern**: MatchList fetches 500 teams on mount
⚠️ **Low concern**: Multiple API calls for lookup data (sports, leagues, teams) could be optimized with caching or combined endpoints

---

## Code Quality Summary

**Strengths:**
- Consistent patterns across all components
- Good TypeScript typing throughout
- Proper use of React Query for data fetching
- Clean separation of concerns (types, services, hooks, components)
- Follows existing codebase conventions

**Areas for Improvement:**
- Error handling in async form submissions
- Team selection reset on league change
- Slug validation order

---

## Recommendation

**Proceed with fixes for issues #2 and #4** before merging, as they affect user experience:
- Issue #2: Invalid team/league combinations can be submitted
- Issue #4: Form stays open on error, confusing users

Other issues are lower priority and can be addressed in follow-up commits.
