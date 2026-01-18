# API Services Overview

Complete reference for all microservices in the Sports Prediction Contests platform.

## Architecture Overview

The platform consists of 7 core microservices communicating via gRPC, with HTTP REST API access through the API Gateway.

| Service | Port | Base URL | Purpose |
|---------|------|----------|---------|
| **API Gateway** | 8080 | `http://localhost:8080` | HTTP REST entry point |
| **User Service** | 8084 | `/v1/auth/*`, `/v1/users/*` | Authentication & user management |
| **Contest Service** | 8085 | `/v1/contests/*` | Contest & team management |
| **Prediction Service** | 8086 | `/v1/predictions/*`, `/v1/events/*` | Predictions & events |
| **Scoring Service** | 8087 | `/v1/scores/*`, `/v1/leaderboard/*` | Scoring & leaderboards |
| **Sports Service** | 8088 | `/v1/sports/*`, `/v1/leagues/*` | Sports data & sync |
| **Notification Service** | 8089 | `/v1/notifications/*` | Multi-channel notifications |

## Authentication

All API endpoints (except registration and login) require JWT authentication:

```bash
# Include JWT token in Authorization header
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/v1/contests
```

## API Gateway (Port 8080)

### Health Check
```bash
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2026-01-18T12:00:00Z"
}
```

## User Service (Port 8084)

### Authentication Endpoints

#### Register User
```bash
POST /v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecureP@ssw0rd2026!",
  "full_name": "John Doe"
}
```

#### Login User
```bash
POST /v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecureP@ssw0rd2026!"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe"
  }
}
```

### User Management Endpoints

#### Get User Profile
```bash
GET /v1/users/profile
Authorization: Bearer JWT_TOKEN
```

#### Update User Profile
```bash
PUT /v1/users/profile
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "full_name": "John Smith",
  "username": "johnsmith"
}
```

## Contest Service (Port 8085)

### Contest Management

#### Create Contest
```bash
POST /v1/contests
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "title": "Premier League Predictions",
  "description": "Predict Premier League match outcomes",
  "sport_type": "football",
  "rules": "{\"scoring_system\": \"standard\"}",
  "start_date": "2026-01-20T00:00:00Z",
  "end_date": "2026-05-20T23:59:59Z",
  "max_participants": 100
}
```

#### List Contests
```bash
GET /v1/contests?page=1&limit=10&status=active&sport_type=football
```

#### Get Contest Details
```bash
GET /v1/contests/{id}
```

#### Join Contest
```bash
POST /v1/contests/{contest_id}/join
Authorization: Bearer JWT_TOKEN
```

### Team Management

#### Create Team
```bash
POST /v1/teams
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "name": "Dream Team",
  "description": "Best prediction team"
}
```

#### Join Team by Invite Code
```bash
POST /v1/teams/join
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "invite_code": "ABC123"
}
```

## Prediction Service (Port 8086)

### Prediction Management

#### Submit Prediction
```bash
POST /v1/predictions
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "contest_id": 1,
  "event_id": 1,
  "prediction_type": "match_outcome",
  "prediction_value": "home_win",
  "prop_type_id": null
}
```

#### Get User Predictions for Contest
```bash
GET /v1/predictions/contest/{contest_id}
Authorization: Bearer JWT_TOKEN
```

### Event Management

#### Create Event
```bash
POST /v1/events
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "title": "Manchester United vs Liverpool",
  "description": "Premier League match",
  "sport_type": "football",
  "start_time": "2026-01-25T15:00:00Z",
  "end_time": "2026-01-25T17:00:00Z",
  "external_id": "match_12345"
}
```

#### Get Time Coefficient
```bash
GET /v1/events/{event_id}/coefficient
```

**Response:**
```json
{
  "coefficient": 1.5,
  "time_remaining_hours": 24,
  "base_points": 10,
  "adjusted_points": 15
}
```

### Prop Types

#### List Prop Types
```bash
GET /v1/prop-types?sport_type=football
```

**Response:**
```json
{
  "prop_types": [
    {
      "id": 1,
      "name": "Total Goals",
      "description": "Predict total goals in match",
      "sport_type": "football",
      "value_type": "number",
      "scoring_rules": "{\"exact_match\": 10, \"close_match\": 5}"
    }
  ]
}
```

## Scoring Service (Port 8087)

### Score Management

#### Calculate Score
```bash
POST /v1/scores/calculate
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "prediction_id": 1,
  "actual_result": "home_win",
  "event_id": 1
}
```

### Leaderboard

#### Get Contest Leaderboard
```bash
GET /v1/contests/{contest_id}/leaderboard?page=1&limit=50
```

**Response:**
```json
{
  "leaderboard": [
    {
      "rank": 1,
      "user_id": 1,
      "username": "johndoe",
      "total_score": 150,
      "correct_predictions": 15,
      "total_predictions": 20,
      "accuracy": 75.0,
      "current_streak": 5
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 100
  }
}
```

#### Get User Rank
```bash
GET /v1/contests/{contest_id}/users/{user_id}/rank
```

### Streak Management

#### Get User Streak
```bash
GET /v1/contests/{contest_id}/users/{user_id}/streak
```

**Response:**
```json
{
  "current_streak": 5,
  "max_streak": 8,
  "streak_multiplier": 1.5,
  "streak_type": "correct_predictions"
}
```

### Analytics

#### Get User Analytics
```bash
GET /v1/users/{user_id}/analytics?contest_id={contest_id}&period=30d
```

**Response:**
```json
{
  "total_predictions": 50,
  "correct_predictions": 35,
  "accuracy": 70.0,
  "total_score": 420,
  "average_score_per_prediction": 8.4,
  "sport_breakdown": {
    "football": {"accuracy": 75.0, "predictions": 30},
    "basketball": {"accuracy": 60.0, "predictions": 20}
  },
  "streak_stats": {
    "current_streak": 3,
    "max_streak": 8,
    "total_streaks": 5
  }
}
```

## Sports Service (Port 8088)

### Sports Management

#### Create Sport
```bash
POST /v1/sports
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "name": "Football",
  "description": "Association Football",
  "rules": "{\"match_duration\": 90, \"overtime_allowed\": true}"
}
```

#### List Sports
```bash
GET /v1/sports
```

### League Management

#### Create League
```bash
POST /v1/leagues
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "name": "Premier League",
  "sport_id": 1,
  "country": "England",
  "season": "2025-2026"
}
```

### Team Management

#### Create Team
```bash
POST /v1/teams
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "name": "Manchester United",
  "league_id": 1,
  "country": "England",
  "founded_year": 1878
}
```

### Match Management

#### Create Match
```bash
POST /v1/matches
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "home_team_id": 1,
  "away_team_id": 2,
  "league_id": 1,
  "scheduled_time": "2026-01-25T15:00:00Z",
  "venue": "Old Trafford"
}
```

#### Update Match Result
```bash
PUT /v1/matches/{id}
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "status": "completed",
  "home_score": 2,
  "away_score": 1,
  "result": "home_win"
}
```

### Data Synchronization

#### Trigger Data Sync
```bash
POST /v1/sports/sync
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "source": "thesportsdb",
  "sport_type": "football",
  "league_ids": [1, 2, 3]
}
```

#### Get Sync Status
```bash
GET /v1/sports/sync/status
```

## Notification Service (Port 8089)

### Notification Management

#### Send Notification
```bash
POST /v1/notifications
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "user_id": 1,
  "type": "contest_result",
  "title": "Contest Results Available",
  "message": "Your contest 'Premier League Predictions' has ended. Check your results!",
  "channels": ["in_app", "telegram"],
  "data": {
    "contest_id": 1,
    "final_rank": 5
  }
}
```

#### Get User Notifications
```bash
GET /v1/users/{user_id}/notifications?page=1&limit=20&unread_only=false
```

#### Mark Notification as Read
```bash
PUT /v1/notifications/{id}/read
Authorization: Bearer JWT_TOKEN
```

### Notification Preferences

#### Get User Preferences
```bash
GET /v1/users/{user_id}/notification-preferences
Authorization: Bearer JWT_TOKEN
```

#### Update Preferences
```bash
PUT /v1/users/{user_id}/notification-preferences
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "in_app_enabled": true,
  "email_enabled": true,
  "telegram_enabled": false,
  "contest_updates": true,
  "score_updates": true,
  "leaderboard_updates": false
}
```

## Error Handling

All services return consistent error responses:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Invalid request data |
| `UNAUTHORIZED` | 401 | Missing or invalid JWT token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource already exists |
| `INTERNAL_ERROR` | 500 | Server error |

## Rate Limiting

API endpoints are rate limited:

- **Authentication**: 5 requests per minute
- **General API**: 100 requests per minute
- **Bulk operations**: 10 requests per minute

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642694400
```

## Pagination

List endpoints support pagination:

```bash
GET /v1/contests?page=1&limit=20&sort=created_at&order=desc
```

**Response:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

## WebSocket Support (Future)

Real-time updates will be available via WebSocket connections:

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:8080/ws');

// Subscribe to contest updates
ws.send(JSON.stringify({
  type: 'subscribe',
  channel: 'contest_updates',
  contest_id: 1
}));
```

---

For detailed documentation of individual services, see:
- [User Service API](user-service.md)
- [Contest Service API](contest-service.md)
- [Prediction Service API](prediction-service.md)
- [Scoring Service API](scoring-service.md)
- [Sports Service API](sports-service.md)
- [Notification Service API](notification-service.md)
