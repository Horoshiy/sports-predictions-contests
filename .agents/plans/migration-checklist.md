# Frontend Migration Checklist: Material-UI → Ant Design

## Phase 1: Foundation & Setup ✅
- [x] Update package.json dependencies
- [x] Install Ant Design packages
- [x] Create theme configuration
- [x] Create helper utilities
- [x] Create migration checklist

## Phase 2: Core Layout & Navigation ✅
- [x] App.tsx - Main application structure
- [x] components/auth/ProtectedRoute.tsx

## Phase 3: Authentication Components ✅
- [x] components/auth/LoginForm.tsx
- [x] components/auth/RegisterForm.tsx
- [x] pages/LoginPage.tsx
- [x] pages/RegisterPage.tsx

## Phase 4: Data Display Components (In Progress - 6/16 Complete)
- [x] components/contests/ContestCard.tsx
- [ ] components/contests/ContestList.tsx
- [ ] components/contests/ParticipantList.tsx
- [ ] components/leaderboard/LeaderboardTable.tsx (Large - 408 lines - requires Table migration)
- [x] components/leaderboard/UserScore.tsx
- [ ] components/teams/TeamList.tsx
- [x] components/teams/TeamMembers.tsx
- [x] components/teams/TeamInvite.tsx
- [ ] components/teams/TeamLeaderboard.tsx
- [ ] components/challenges/ChallengeCard.tsx
- [ ] components/challenges/ChallengeList.tsx
- [ ] components/challenges/ChallengeDialog.tsx
- [ ] components/sports/SportList.tsx
- [ ] components/sports/LeagueList.tsx
- [ ] components/sports/TeamList.tsx (sports)
- [ ] components/sports/MatchList.tsx
- [x] components/analytics/ExportButton.tsx
- [x] components/predictions/CoefficientIndicator.tsx

## Phase 5: Form Components
- [ ] components/contests/ContestForm.tsx
- [ ] components/teams/TeamForm.tsx
- [ ] components/sports/SportForm.tsx
- [ ] components/sports/LeagueForm.tsx
- [ ] components/sports/TeamForm.tsx (sports)
- [ ] components/sports/MatchForm.tsx
- [ ] components/predictions/PredictionForm.tsx
- [ ] components/predictions/PropTypeSelector.tsx
- [ ] components/profile/ProfileForm.tsx
- [ ] components/profile/AvatarUpload.tsx
- [ ] components/profile/PrivacySettings.tsx

## Phase 6: Pages & Complex Components
- [ ] pages/ContestsPage.tsx
- [ ] pages/SportsPage.tsx
- [ ] pages/TeamsPage.tsx
- [ ] pages/PredictionsPage.tsx
- [ ] pages/AnalyticsPage.tsx
- [ ] pages/ProfilePage.tsx
- [ ] components/predictions/EventList.tsx
- [ ] components/predictions/EventCard.tsx
- [ ] components/predictions/PredictionList.tsx
- [ ] components/predictions/CoefficientIndicator.tsx
- [ ] components/analytics/AccuracyChart.tsx
- [ ] components/analytics/SportBreakdown.tsx
- [ ] components/analytics/PlatformComparison.tsx
- [ ] components/analytics/ExportButton.tsx
- [ ] components/profile/ProfileCompletion.tsx

## Phase 7: Utilities & Context ✅
- [x] contexts/ToastContext.tsx
- [ ] contexts/AuthContext.tsx (no changes needed - no MUI dependencies)

## Phase 8: Testing & Cleanup
- [ ] Remove all @mui imports
- [ ] Remove MUI dependencies from package.json
- [ ] Update global styles
- [ ] Test all pages and components
- [ ] Update documentation

## Component Mapping Reference

| MUI Component | Ant Design | Status |
|---------------|-----------|--------|
| Box | Space/Flex | - |
| Container | Layout.Content | - |
| Paper | Card | - |
| Typography | Typography | - |
| Button | Button | - |
| TextField | Input | - |
| Select | Select | - |
| Dialog | Modal | - |
| Snackbar | message | - |
| Alert | Alert | - |
| Chip | Tag | - |
| Avatar | Avatar | - |
| Tabs | Tabs | - |
| Table | Table | - |
| Card | Card | - |
| AppBar | Layout.Header | - |
| Grid | Row/Col | - |
| DatePicker | DatePicker | - |

## Notes

- Keep react-hook-form for validation
- Migrate date-fns to dayjs
- Test each phase before proceeding
- Document any issues or deviations
