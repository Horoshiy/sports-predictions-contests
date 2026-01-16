# Code Review: Sports Data Integration

**Date**: 2026-01-16
**Feature**: Sports Data Integration with TheSportsDB API
**Reviewer**: Kiro CLI

## Stats

- Files Modified: 14
- Files Added: 4
- Files Deleted: 0
- New lines: ~777
- Deleted lines: ~30

---

## Issues Found

### CRITICAL

```
severity: critical
file: backend/sports-service/internal/sync/sync_service.go
line: 192
issue: N+1 query problem in SyncMatchResults
detail: For each scheduled match (up to 100), the code makes an individual HTTP request to GetEventByID. With 100 matches, this results in 100 sequential HTTP calls, which will be extremely slow and may trigger rate limiting from TheSportsDB API.
suggestion: Batch process matches or add rate limiting/throttling. Consider fetching events by league instead of individual lookups:
```go
// Add delay between requests to avoid rate limiting
time.Sleep(100 * time.Millisecond)
```
Or better, track which leagues have scheduled matches and fetch events by league.
```

---

### HIGH

```
severity: high
file: backend/sports-service/internal/sync/sync_service.go
line: 23-24
issue: lastSyncAt field is not thread-safe
detail: The lastSyncAt field is accessed and modified without synchronization. When the background worker runs concurrently with manual sync triggers via TriggerSync, this can cause data races.
suggestion: Add a mutex to protect lastSyncAt access:
```go
type SyncService struct {
    // ...
    lastSyncAt *time.Time
    mu         sync.RWMutex
}

func (s *SyncService) GetLastSyncAt() *time.Time {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.lastSyncAt
}
```
```

```
severity: high
file: backend/sports-service/internal/sync/worker.go
line: 68
issue: Initial sync runs immediately on Start() without error handling
detail: When the worker starts, it immediately calls runSync() which makes external API calls. If the API is down or rate-limited at startup, this could cause issues. The error is logged but the worker continues, which is fine, but there's no backoff mechanism.
suggestion: Add exponential backoff for failed syncs or make initial sync optional:
```go
func (w *SyncWorker) Start() {
    // ...
    go w.run()
    // Consider: go w.runInitialSyncWithDelay() to avoid startup thundering herd
}
```
```

```
severity: high
file: backend/sports-service/internal/sync/sync_service.go
line: 275-278
issue: Regex compilation on every slugify call
detail: Two regex patterns are compiled on every call to slugify(), which is called for every sport, league, and team during sync. This is inefficient.
suggestion: Compile regexes once at package level:
```go
var (
    slugInvalidChars = regexp.MustCompile(`[^a-z0-9-]+`)
    slugMultipleDash = regexp.MustCompile(`-+`)
)

func (s *SyncService) slugify(name string) string {
    slug := strings.ToLower(name)
    slug = strings.ReplaceAll(slug, " ", "-")
    slug = slugInvalidChars.ReplaceAllString(slug, "")
    slug = slugMultipleDash.ReplaceAllString(slug, "-")
    return strings.Trim(slug, "-")
}
```
```

---

### MEDIUM

```
severity: medium
file: backend/sports-service/internal/external/thesportsdb.go
line: 137
issue: URL parameter not escaped
detail: The leagueID parameter in GetTeamsByLeague is directly interpolated into the URL without escaping. While IDs from TheSportsDB are typically numeric strings, this could be a security issue if user input ever reaches this function.
suggestion: Use url.QueryEscape or url.Values for query parameters:
```go
url := fmt.Sprintf("%s/lookup_all_teams.php?id=%s", c.baseURL, url.QueryEscape(leagueID))
```
```

```
severity: medium
file: backend/sports-service/internal/sync/sync_service.go
line: 186-188
issue: Status comparison logic is incorrect
detail: The condition `event.StrStatus != "" && event.StrStatus != match.Status` compares the raw API status (e.g., "FT") with the mapped status (e.g., "completed"). This will always be true for completed matches, causing unnecessary updates.
suggestion: Compare mapped statuses:
```go
mappedStatus := s.mapStatus(event.StrStatus)
if mappedStatus != match.Status {
    match.Status = mappedStatus
    // ...
}
```
```

```
severity: medium
file: backend/sports-service/internal/sync/sync_service.go
line: 248-252
issue: Silent fallback to tomorrow for unparseable dates
detail: When both date parsing attempts fail, the code silently defaults to tomorrow. This could create matches with incorrect dates that are hard to debug.
suggestion: Return an error or log a warning when date parsing fails:
```go
if err != nil {
    log.Printf("[WARN] Failed to parse date for event %s, using fallback", api.IDEvent)
    scheduledAt = time.Now().Add(24 * time.Hour)
}
```
```

```
severity: medium
file: backend/sports-service/internal/service/sports_service.go
line: 530-532
issue: GetSyncStatus doesn't return actual last sync time
detail: The lastSyncAt variable is always empty string because the code to get it from sync service is commented out/incomplete.
suggestion: Implement the actual retrieval:
```go
if s.syncWorker != nil {
    if lastSync := s.syncWorker.GetLastSyncAt(); lastSync != nil {
        lastSyncAt = lastSync.Format(time.RFC3339)
    }
}
```
Note: This requires adding GetLastSyncAt() method to SyncWorker that delegates to SyncService.
```

---

### LOW

```
severity: low
file: backend/sports-service/internal/sync/worker.go
line: 84-96
issue: Duplicate logging in runSync
detail: The runSync method logs "Synced %d sports" and "Synced %d leagues" but SyncSports and SyncLeagues already log the same information. This results in duplicate log entries.
suggestion: Remove the duplicate logs in runSync or remove them from the sync service methods.
```

```
severity: low
file: backend/sports-service/internal/external/thesportsdb.go
line: 91-92
issue: Duplicate response types
detail: EventsResponse and EventResponse are identical structs. Only one is needed.
suggestion: Remove EventResponse and use EventsResponse for both:
```go
// Remove EventResponse, use EventsResponse instead
```
```

```
severity: low
file: tests/sports-service/sync_test.go
line: 1
issue: Test package naming inconsistency
detail: The test package is named `sports_service_test` but the file is in `tests/sports-service/`. This may cause import issues when running tests.
suggestion: Either move tests to `backend/sports-service/internal/external/` with `_test.go` suffix, or update the package name to match the directory structure.
```

---

## Security Analysis

- **No hardcoded secrets**: API client uses configurable base URL ✅
- **No SQL injection**: Uses GORM parameterized queries ✅
- **URL parameters**: Should be escaped (medium issue noted)
- **Rate limiting**: No protection against API rate limits (noted in critical issue)

## Performance Analysis

- **N+1 HTTP calls**: Critical issue in SyncMatchResults
- **Regex compilation**: High issue with repeated compilation
- **Database queries**: Upsert pattern is efficient (check-then-update)

## Code Quality

- Clean separation of concerns (external client, sync service, worker)
- Good error handling with logging
- Follows existing codebase patterns
- Tests cover happy path and error cases

---

## Summary

The implementation is well-structured and follows existing patterns. The main concerns are:

1. **Critical**: N+1 HTTP requests in SyncMatchResults will cause performance issues and potential rate limiting
2. **High**: Thread safety issue with lastSyncAt field
3. **High**: Inefficient regex compilation on every slugify call
4. **Medium**: Status comparison bug will cause unnecessary database updates

Recommend fixing critical and high issues before merging.
