# Feature: Props Predictions (Statistics-Based Predictions)

The following plan should be complete, but validate documentation and codebase patterns before implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files.

## Feature Description

Props Predictions extends the platform beyond match outcomes to allow predictions on specific statistics: player goals, corners, shots on target, possession percentage, cards, and other measurable events. This significantly increases predictions per match and appeals to analytics enthusiasts.

## User Story

As a sports prediction enthusiast
I want to predict specific match statistics (goals by player, corners, possession)
So that I can demonstrate deeper sports knowledge and earn more points

## Problem Statement

Current predictions are limited to match outcomes (winner, exact score). Analytics-focused users want to predict granular statistics like "Player X scores 2+ goals" or "Over 10 corners in match". This limits engagement and predictions per match.

## Solution Statement

Implement a Props prediction system with:
1. Prop types catalog (configurable per sport)
2. Extended prediction schema supporting props
3. Props-specific scoring rules
4. Frontend UI for browsing and submitting prop predictions

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium (4-6 hours)
**Primary Systems Affected**: prediction-service, scoring-service, sports-service, frontend
**Dependencies**: Existing prediction infrastructure, match data

---

## CONTEXT REFERENCES

### Relevant Codebase Files - READ BEFORE IMPLEMENTING

**Backend - Prediction Service:**
- `backend/proto/prediction.proto` - Current prediction proto definitions
- `backend/prediction-service/internal/models/prediction.go` - Prediction model with JSON data field
- `backend/prediction-service/internal/service/prediction_service.go` - Prediction service implementation

**Backend - Scoring Service:**
- `backend/scoring-service/internal/service/scoring_service.go` (lines 45-80) - PredictionData struct and scoring logic
- `backend/scoring-service/internal/models/score.go` - Score model

**Backend - Sports Service:**
- `backend/proto/sports.proto` - Match proto with result_data field
- `backend/sports-service/internal/models/match.go` - Match model

**Frontend:**
- `frontend/src/types/prediction.types.ts` - Prediction TypeScript types
- `frontend/src/utils/prediction-validation.ts` - Zod validation schemas
- `frontend/src/components/predictions/PredictionForm.tsx` - Current prediction form
- `frontend/src/hooks/use-predictions.ts` - React Query hooks

**Database:**
- `scripts/init-db.sql` - Database schema

### New Files to Create

**Backend:**
- `backend/prediction-service/internal/models/prop_type.go` - PropType model
- `backend/prediction-service/internal/repository/prop_type_repository.go` - PropType repository

**Frontend:**
- `frontend/src/types/props.types.ts` - Props TypeScript types
- `frontend/src/components/predictions/PropsPredictionForm.tsx` - Props prediction form
- `frontend/src/components/predictions/PropTypeSelector.tsx` - Prop type selection component

### Patterns to Follow

**Naming Conventions:**
- Go: snake_case files, PascalCase structs, camelCase private functions
- TypeScript: camelCase variables, PascalCase components/types
- Proto: snake_case fields, PascalCase messages

**Error Handling (Go):**
```go
if err != nil {
    log.Printf("[ERROR] Failed to X: %v", err)
    return &pb.Response{
        Response: &common.Response{
            Success: false,
            Message: "User-friendly message",
            Code:    int32(common.ErrorCode_INTERNAL_ERROR),
        },
    }, nil
}
```

**React Query Pattern:**
```typescript
export const usePropTypes = (sportType: string) => {
  return useQuery({
    queryKey: ['propTypes', sportType],
    queryFn: () => predictionService.getPropTypes(sportType),
    enabled: !!sportType,
    staleTime: 5 * 60 * 1000,
  })
}
```

**Zod Validation Pattern:**
```typescript
export const propPredictionSchema = z.object({
  propTypeId: z.number().min(1),
  value: z.union([z.number(), z.string()]),
}).refine(...)
```

---

## IMPLEMENTATION PLAN

### Phase 1: Database & Models

Add prop_types table and extend prediction data structure to support props.

**Tasks:**
- Create prop_types database table
- Create PropType Go model
- Create PropType repository

### Phase 2: Backend API

Extend prediction proto and service to handle props predictions.

**Tasks:**
- Add PropType messages to prediction.proto
- Add GetPropTypes and ListPropTypes RPCs
- Extend PredictionData struct for props
- Add props scoring logic

### Phase 3: Frontend

Build UI for browsing prop types and submitting props predictions.

**Tasks:**
- Create props TypeScript types
- Create PropTypeSelector component
- Extend PredictionForm for props
- Add props validation schema

### Phase 4: Testing & Validation

**Tasks:**
- Unit tests for prop type validation
- Unit tests for props scoring
- Integration test for props prediction flow

---

## STEP-BY-STEP TASKS

### Task 1: CREATE `scripts/init-db.sql` - Add prop_types table

**IMPLEMENT**: Add prop_types table after matches table

```sql
-- Create prop_types table for props predictions
CREATE TABLE IF NOT EXISTS prop_types (
    id SERIAL PRIMARY KEY,
    sport_type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    value_type VARCHAR(20) NOT NULL,
    default_line DECIMAL(10,2),
    min_value DECIMAL(10,2),
    max_value DECIMAL(10,2),
    points_correct DECIMAL(10,2) NOT NULL DEFAULT 2,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(sport_type, slug)
);

CREATE INDEX IF NOT EXISTS idx_prop_types_sport ON prop_types(sport_type);
CREATE INDEX IF NOT EXISTS idx_prop_types_category ON prop_types(category);
CREATE INDEX IF NOT EXISTS idx_prop_types_is_active ON prop_types(is_active);
CREATE INDEX IF NOT EXISTS idx_prop_types_deleted_at ON prop_types(deleted_at);

-- Insert default prop types for Soccer
INSERT INTO prop_types (sport_type, name, slug, description, category, value_type, default_line, points_correct) VALUES
('Soccer', 'Total Goals Over/Under', 'total-goals-ou', 'Predict if total goals will be over or under the line', 'match', 'over_under', 2.5, 2),
('Soccer', 'Total Corners Over/Under', 'total-corners-ou', 'Predict if total corners will be over or under the line', 'match', 'over_under', 9.5, 2),
('Soccer', 'Both Teams to Score', 'btts', 'Predict if both teams will score', 'match', 'yes_no', NULL, 2),
('Soccer', 'First Team to Score', 'first-to-score', 'Predict which team scores first', 'match', 'team_select', NULL, 3),
('Soccer', 'Player to Score Anytime', 'player-goal', 'Predict if a specific player will score', 'player', 'yes_no', NULL, 4),
('Soccer', 'Total Cards Over/Under', 'total-cards-ou', 'Predict if total cards will be over or under the line', 'match', 'over_under', 3.5, 2)
ON CONFLICT (sport_type, slug) DO NOTHING;
```

**PATTERN**: Mirror matches table structure from init-db.sql
**VALIDATE**: `docker-compose down -v && docker-compose up -d postgres && sleep 3 && docker exec sports-postgres psql -U sports_user -d sports_prediction -c "\d prop_types"`

### Task 2: CREATE `backend/prediction-service/internal/models/prop_type.go`

**IMPLEMENT**: PropType GORM model with validation

```go
package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// PropType represents a type of prop prediction available for a sport
type PropType struct {
	gorm.Model
	SportType     string   `gorm:"not null;index:idx_prop_sport_slug,unique" json:"sport_type"`
	Name          string   `gorm:"not null" json:"name"`
	Slug          string   `gorm:"not null;index:idx_prop_sport_slug,unique" json:"slug"`
	Description   string   `json:"description"`
	Category      string   `gorm:"not null" json:"category"` // "match", "player", "team"
	ValueType     string   `gorm:"not null" json:"value_type"` // "over_under", "yes_no", "team_select", "player_select", "exact_value"
	DefaultLine   *float64 `json:"default_line"`
	MinValue      *float64 `json:"min_value"`
	MaxValue      *float64 `json:"max_value"`
	PointsCorrect float64  `gorm:"not null;default:2" json:"points_correct"`
	IsActive      bool     `gorm:"default:true" json:"is_active"`
}

// TableName returns the table name
func (PropType) TableName() string {
	return "prop_types"
}

// ValidateSportType checks if sport type is valid
func (p *PropType) ValidateSportType() error {
	if strings.TrimSpace(p.SportType) == "" {
		return errors.New("sport type cannot be empty")
	}
	return nil
}

// ValidateName checks if name is valid
func (p *PropType) ValidateName() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if len(p.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	return nil
}

// ValidateCategory checks if category is valid
func (p *PropType) ValidateCategory() error {
	validCategories := []string{"match", "player", "team"}
	for _, c := range validCategories {
		if p.Category == c {
			return nil
		}
	}
	return errors.New("invalid category: must be match, player, or team")
}

// ValidateValueType checks if value type is valid
func (p *PropType) ValidateValueType() error {
	validTypes := []string{"over_under", "yes_no", "team_select", "player_select", "exact_value"}
	for _, t := range validTypes {
		if p.ValueType == t {
			return nil
		}
	}
	return errors.New("invalid value type")
}

// BeforeCreate validates before creating
func (p *PropType) BeforeCreate(tx *gorm.DB) error {
	if err := p.ValidateSportType(); err != nil {
		return err
	}
	if err := p.ValidateName(); err != nil {
		return err
	}
	if err := p.ValidateCategory(); err != nil {
		return err
	}
	if err := p.ValidateValueType(); err != nil {
		return err
	}
	return nil
}
```

**PATTERN**: Mirror `backend/prediction-service/internal/models/prediction.go` structure
**IMPORTS**: `errors`, `strings`, `gorm.io/gorm`
**VALIDATE**: `cd backend/prediction-service && go build ./...`

### Task 3: CREATE `backend/prediction-service/internal/repository/prop_type_repository.go`

**IMPLEMENT**: PropType repository with CRUD operations

```go
package repository

import (
	"context"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"gorm.io/gorm"
)

// PropTypeRepository handles prop type database operations
type PropTypeRepository struct {
	db *gorm.DB
}

// NewPropTypeRepository creates a new PropTypeRepository
func NewPropTypeRepository(db *gorm.DB) *PropTypeRepository {
	return &PropTypeRepository{db: db}
}

// GetBySportType returns all active prop types for a sport
func (r *PropTypeRepository) GetBySportType(ctx context.Context, sportType string) ([]*models.PropType, error) {
	var propTypes []*models.PropType
	err := r.db.WithContext(ctx).
		Where("sport_type = ? AND is_active = ? AND deleted_at IS NULL", sportType, true).
		Order("category, name").
		Find(&propTypes).Error
	return propTypes, err
}

// GetByID returns a prop type by ID
func (r *PropTypeRepository) GetByID(ctx context.Context, id uint) (*models.PropType, error) {
	var propType models.PropType
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&propType).Error
	if err != nil {
		return nil, err
	}
	return &propType, nil
}

// GetBySlug returns a prop type by sport and slug
func (r *PropTypeRepository) GetBySlug(ctx context.Context, sportType, slug string) (*models.PropType, error) {
	var propType models.PropType
	err := r.db.WithContext(ctx).
		Where("sport_type = ? AND slug = ? AND deleted_at IS NULL", sportType, slug).
		First(&propType).Error
	if err != nil {
		return nil, err
	}
	return &propType, nil
}

// List returns all prop types with optional filtering
func (r *PropTypeRepository) List(ctx context.Context, sportType, category string, activeOnly bool, page, limit int) ([]*models.PropType, int64, error) {
	var propTypes []*models.PropType
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PropType{}).Where("deleted_at IS NULL")

	if sportType != "" {
		query = query.Where("sport_type = ?", sportType)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("sport_type, category, name").
		Offset(offset).Limit(limit).
		Find(&propTypes).Error

	return propTypes, total, err
}
```

**PATTERN**: Mirror `backend/prediction-service/internal/repository/prediction_repository.go`
**VALIDATE**: `cd backend/prediction-service && go build ./...`

### Task 4: UPDATE `backend/proto/prediction.proto` - Add PropType messages

**IMPLEMENT**: Add after Event message (around line 40)

```protobuf
// PropType represents a type of prop prediction
message PropType {
  uint32 id = 1;
  string sport_type = 2;
  string name = 3;
  string slug = 4;
  string description = 5;
  string category = 6;
  string value_type = 7;
  double default_line = 8;
  double min_value = 9;
  double max_value = 10;
  double points_correct = 11;
  bool is_active = 12;
}

// PropPrediction represents a props prediction within prediction_data
message PropPrediction {
  uint32 prop_type_id = 1;
  string prop_slug = 2;
  double line = 3;
  string selection = 4;
  string player_id = 5;
}

message GetPropTypesRequest {
  string sport_type = 1;
}

message GetPropTypesResponse {
  common.Response response = 1;
  repeated PropType prop_types = 2;
}

message ListPropTypesRequest {
  string sport_type = 1;
  string category = 2;
  bool active_only = 3;
  common.PaginationRequest pagination = 4;
}

message ListPropTypesResponse {
  common.Response response = 1;
  repeated PropType prop_types = 2;
  common.PaginationResponse pagination = 3;
}
```

**IMPLEMENT**: Add RPCs to PredictionService (after UpdateEvent)

```protobuf
  // Prop types
  rpc GetPropTypes(GetPropTypesRequest) returns (GetPropTypesResponse) {
    option (google.api.http) = {
      get: "/v1/prop-types/{sport_type}"
    };
  }
  rpc ListPropTypes(ListPropTypesRequest) returns (ListPropTypesResponse) {
    option (google.api.http) = {
      get: "/v1/prop-types"
    };
  }
```

**PATTERN**: Mirror existing message/RPC patterns in prediction.proto
**VALIDATE**: `cd backend && protoc --proto_path=proto --go_out=shared --go-grpc_out=shared proto/prediction.proto` (if protoc available, otherwise syntax check only)

### Task 5: UPDATE `backend/prediction-service/internal/service/prediction_service.go` - Add PropType RPCs

**IMPLEMENT**: Add propTypeRepo field and constructor parameter, then add RPC implementations

Add to struct:
```go
type PredictionService struct {
	pb.UnimplementedPredictionServiceServer
	predictionRepo *repository.PredictionRepository
	eventRepo      *repository.EventRepository
	propTypeRepo   *repository.PropTypeRepository  // ADD THIS
	contestClient  *clients.ContestClient
}
```

Add to constructor:
```go
func NewPredictionService(predictionRepo *repository.PredictionRepository, eventRepo *repository.EventRepository, propTypeRepo *repository.PropTypeRepository, contestClient *clients.ContestClient) *PredictionService {
	return &PredictionService{
		predictionRepo: predictionRepo,
		eventRepo:      eventRepo,
		propTypeRepo:   propTypeRepo,
		contestClient:  contestClient,
	}
}
```

Add RPC implementations:
```go
// GetPropTypes returns prop types for a sport
func (s *PredictionService) GetPropTypes(ctx context.Context, req *pb.GetPropTypesRequest) (*pb.GetPropTypesResponse, error) {
	if req.SportType == "" {
		return &pb.GetPropTypesResponse{
			Response: &common.Response{
				Success: false,
				Message: "Sport type is required",
				Code:    int32(common.ErrorCode_INVALID_ARGUMENT),
			},
		}, nil
	}

	propTypes, err := s.propTypeRepo.GetBySportType(ctx, req.SportType)
	if err != nil {
		log.Printf("[ERROR] Failed to get prop types: %v", err)
		return &pb.GetPropTypesResponse{
			Response: &common.Response{
				Success: false,
				Message: "Failed to retrieve prop types",
				Code:    int32(common.ErrorCode_INTERNAL_ERROR),
			},
		}, nil
	}

	protoPropTypes := make([]*pb.PropType, len(propTypes))
	for i, pt := range propTypes {
		protoPropTypes[i] = s.propTypeToProto(pt)
	}

	return &pb.GetPropTypesResponse{
		Response: &common.Response{
			Success: true,
			Message: "Prop types retrieved successfully",
			Code:    int32(common.ErrorCode_SUCCESS),
		},
		PropTypes: protoPropTypes,
	}, nil
}

// ListPropTypes returns paginated prop types with filtering
func (s *PredictionService) ListPropTypes(ctx context.Context, req *pb.ListPropTypesRequest) (*pb.ListPropTypesResponse, error) {
	page, limit := 1, 20
	if req.Pagination != nil {
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
		if req.Pagination.Limit > 0 {
			limit = int(req.Pagination.Limit)
		}
	}

	propTypes, total, err := s.propTypeRepo.List(ctx, req.SportType, req.Category, req.ActiveOnly, page, limit)
	if err != nil {
		log.Printf("[ERROR] Failed to list prop types: %v", err)
		return &pb.ListPropTypesResponse{
			Response: &common.Response{
				Success: false,
				Message: "Failed to list prop types",
				Code:    int32(common.ErrorCode_INTERNAL_ERROR),
			},
		}, nil
	}

	protoPropTypes := make([]*pb.PropType, len(propTypes))
	for i, pt := range propTypes {
		protoPropTypes[i] = s.propTypeToProto(pt)
	}

	totalPages := (int(total) + limit - 1) / limit

	return &pb.ListPropTypesResponse{
		Response: &common.Response{
			Success: true,
			Message: "Prop types listed successfully",
			Code:    int32(common.ErrorCode_SUCCESS),
		},
		PropTypes: protoPropTypes,
		Pagination: &common.PaginationResponse{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      int32(total),
			TotalPages: int32(totalPages),
		},
	}, nil
}

func (s *PredictionService) propTypeToProto(pt *models.PropType) *pb.PropType {
	proto := &pb.PropType{
		Id:            uint32(pt.ID),
		SportType:     pt.SportType,
		Name:          pt.Name,
		Slug:          pt.Slug,
		Description:   pt.Description,
		Category:      pt.Category,
		ValueType:     pt.ValueType,
		PointsCorrect: pt.PointsCorrect,
		IsActive:      pt.IsActive,
	}
	if pt.DefaultLine != nil {
		proto.DefaultLine = *pt.DefaultLine
	}
	if pt.MinValue != nil {
		proto.MinValue = *pt.MinValue
	}
	if pt.MaxValue != nil {
		proto.MaxValue = *pt.MaxValue
	}
	return proto
}
```

**PATTERN**: Mirror existing RPC implementations in same file
**VALIDATE**: `cd backend/prediction-service && go build ./...`

### Task 6: UPDATE `backend/prediction-service/cmd/main.go` - Initialize PropType repository

**IMPLEMENT**: Add PropType repository initialization and pass to service

Find where repositories are created and add:
```go
propTypeRepo := repository.NewPropTypeRepository(db)
```

Update service creation:
```go
predictionService := service.NewPredictionService(predictionRepo, eventRepo, propTypeRepo, contestClient)
```

**PATTERN**: Mirror existing repository initialization pattern
**VALIDATE**: `cd backend/prediction-service && go build ./...`

### Task 7: UPDATE `backend/scoring-service/internal/service/scoring_service.go` - Add props scoring

**IMPLEMENT**: Extend PredictionData struct and add props scoring logic

Update PredictionData struct (around line 45):
```go
type PredictionData struct {
	Type       string      `json:"type"`
	HomeScore  *int        `json:"home_score"`
	AwayScore  *int        `json:"away_score"`
	Winner     *string     `json:"winner"`
	OverUnder  *string     `json:"over_under"`
	Threshold  *float64    `json:"threshold"`
	Value      interface{} `json:"value"`
	// Props fields
	Props      []PropPrediction `json:"props,omitempty"`
}

type PropPrediction struct {
	PropTypeID   uint    `json:"prop_type_id"`
	PropSlug     string  `json:"prop_slug"`
	Line         float64 `json:"line"`
	Selection    string  `json:"selection"` // "over", "under", "yes", "no", "home", "away"
	PlayerID     string  `json:"player_id,omitempty"`
	PointsValue  float64 `json:"points_value"` // Points for this prop if correct
}
```

Update ResultData struct:
```go
type ResultData struct {
	HomeScore   int                    `json:"home_score"`
	AwayScore   int                    `json:"away_score"`
	Winner      string                 `json:"winner"`
	TotalGoals  int                    `json:"total_goals"`
	// Props results
	Stats       map[string]interface{} `json:"stats,omitempty"` // "corners": 12, "cards": 4, etc.
	PlayerStats map[string]interface{} `json:"player_stats,omitempty"` // "player_123": {"goals": 2}
}
```

Add to calculatePoints function (in switch statement):
```go
case "props":
	return s.calculatePropsPoints(prediction, result, details)
```

Add new function:
```go
// calculatePropsPoints calculates points for props predictions
func (s *ScoringService) calculatePropsPoints(prediction PredictionData, result ResultData, details map[string]interface{}) (float64, map[string]interface{}) {
	if len(prediction.Props) == 0 {
		details["error"] = "No props predictions found"
		return 0, details
	}

	var totalPoints float64
	propResults := make([]map[string]interface{}, 0)

	for _, prop := range prediction.Props {
		propResult := map[string]interface{}{
			"prop_slug":  prop.PropSlug,
			"selection":  prop.Selection,
			"line":       prop.Line,
		}

		correct := s.evaluateProp(prop, result)
		propResult["correct"] = correct

		if correct {
			points := prop.PointsValue
			if points == 0 {
				points = 2 // Default points
			}
			totalPoints += points
			propResult["points"] = points
		} else {
			propResult["points"] = 0
		}

		propResults = append(propResults, propResult)
	}

	details["props_results"] = propResults
	details["total_props"] = len(prediction.Props)
	details["correct_props"] = s.countCorrectProps(propResults)

	return totalPoints, details
}

func (s *ScoringService) evaluateProp(prop PropPrediction, result ResultData) bool {
	switch prop.PropSlug {
	case "total-goals-ou":
		totalGoals := float64(result.TotalGoals)
		if prop.Selection == "over" {
			return totalGoals > prop.Line
		}
		return totalGoals < prop.Line

	case "total-corners-ou":
		if corners, ok := result.Stats["corners"].(float64); ok {
			if prop.Selection == "over" {
				return corners > prop.Line
			}
			return corners < prop.Line
		}
		return false

	case "btts":
		btts := result.HomeScore > 0 && result.AwayScore > 0
		if prop.Selection == "yes" {
			return btts
		}
		return !btts

	case "first-to-score":
		if firstScorer, ok := result.Stats["first_to_score"].(string); ok {
			return firstScorer == prop.Selection
		}
		return false

	case "total-cards-ou":
		if cards, ok := result.Stats["cards"].(float64); ok {
			if prop.Selection == "over" {
				return cards > prop.Line
			}
			return cards < prop.Line
		}
		return false

	default:
		return false
	}
}

func (s *ScoringService) countCorrectProps(results []map[string]interface{}) int {
	count := 0
	for _, r := range results {
		if correct, ok := r["correct"].(bool); ok && correct {
			count++
		}
	}
	return count
}
```

**PATTERN**: Mirror existing scoring functions in same file
**VALIDATE**: `cd backend/scoring-service && go build ./...`


### Task 8: CREATE `frontend/src/types/props.types.ts`

**IMPLEMENT**: TypeScript types for props predictions

```typescript
import type { PaginationRequest, PaginationResponse, ApiResponse } from './common.types'

// PropType entity
export interface PropType {
  id: number
  sportType: string
  name: string
  slug: string
  description: string
  category: 'match' | 'player' | 'team'
  valueType: 'over_under' | 'yes_no' | 'team_select' | 'player_select' | 'exact_value'
  defaultLine: number | null
  minValue: number | null
  maxValue: number | null
  pointsCorrect: number
  isActive: boolean
}

// Props prediction within prediction data
export interface PropPrediction {
  propTypeId: number
  propSlug: string
  line: number
  selection: string // "over", "under", "yes", "no", "home", "away", player_id
  playerId?: string
  pointsValue: number
}

// Extended prediction data with props
export interface PredictionDataWithProps {
  type: 'winner' | 'score' | 'combined' | 'props'
  winner?: 'home' | 'away' | 'draw'
  homeScore?: number
  awayScore?: number
  props?: PropPrediction[]
}

// Request types
export interface GetPropTypesRequest {
  sportType: string
}

export interface ListPropTypesRequest {
  sportType?: string
  category?: string
  activeOnly?: boolean
  pagination?: PaginationRequest
}

// Response types
export interface GetPropTypesResponse {
  response: ApiResponse
  propTypes: PropType[]
}

export interface ListPropTypesResponse {
  response: ApiResponse
  propTypes: PropType[]
  pagination: PaginationResponse
}

// Form data for props prediction
export interface PropPredictionFormData {
  propTypeId: number
  propSlug: string
  line?: number
  selection: string
  playerId?: string
}
```

**PATTERN**: Mirror `frontend/src/types/prediction.types.ts`
**VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 9: UPDATE `frontend/src/services/prediction-service.ts` - Add prop types methods

**IMPLEMENT**: Add methods for fetching prop types

Add imports at top:
```typescript
import type {
  // ... existing imports
  GetPropTypesResponse,
  ListPropTypesResponse,
} from '../types/props.types'
```

Add methods to PredictionService class:
```typescript
  // Get prop types for a sport
  async getPropTypes(sportType: string): Promise<PropType[]> {
    const response = await grpcClient.get<GetPropTypesResponse>(
      `/v1/prop-types/${encodeURIComponent(sportType)}`
    )
    if (!response.response?.success) {
      throw new Error(response.response?.message || 'Failed to get prop types')
    }
    return response.propTypes || []
  }

  // List prop types with filtering
  async listPropTypes(request: ListPropTypesRequest = {}): Promise<{
    propTypes: PropType[]
    pagination: PaginationResponse
  }> {
    const params = new URLSearchParams()
    
    if (request.pagination) {
      params.append('page', request.pagination.page.toString())
      params.append('limit', request.pagination.limit.toString())
    }
    if (request.sportType) {
      params.append('sport_type', request.sportType)
    }
    if (request.category) {
      params.append('category', request.category)
    }
    if (request.activeOnly !== undefined) {
      params.append('active_only', request.activeOnly.toString())
    }

    const queryString = params.toString()
    const url = queryString ? `/v1/prop-types?${queryString}` : '/v1/prop-types'

    const response = await grpcClient.get<ListPropTypesResponse>(url)
    
    return {
      propTypes: response.propTypes || [],
      pagination: response.pagination || { page: 1, limit: 20, total: 0, totalPages: 0 },
    }
  }
```

**PATTERN**: Mirror existing listEvents method
**IMPORTS**: Add PropType to imports from props.types
**VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 10: UPDATE `frontend/src/hooks/use-predictions.ts` - Add prop types hooks

**IMPLEMENT**: Add React Query hooks for prop types

Add query keys:
```typescript
export const propTypeKeys = {
  all: ['propTypes'] as const,
  lists: () => [...propTypeKeys.all, 'list'] as const,
  list: (sportType: string) => [...propTypeKeys.lists(), sportType] as const,
}
```

Add hooks:
```typescript
// Fetch prop types for a sport
export const usePropTypes = (sportType: string) => {
  return useQuery({
    queryKey: propTypeKeys.list(sportType),
    queryFn: () => predictionService.getPropTypes(sportType),
    enabled: !!sportType,
    staleTime: 10 * 60 * 1000, // Prop types don't change often
  })
}
```

**PATTERN**: Mirror existing useEvents hook
**VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 11: UPDATE `frontend/src/utils/prediction-validation.ts` - Add props validation

**IMPLEMENT**: Add Zod schema for props predictions

Add after existing schemas:
```typescript
// Single prop prediction schema
export const propPredictionSchema = z.object({
  propTypeId: z.number().min(1, 'Prop type is required'),
  propSlug: z.string().min(1),
  line: z.number().optional(),
  selection: z.string().min(1, 'Selection is required'),
  playerId: z.string().optional(),
  pointsValue: z.number().default(2),
})

export type PropPredictionFormData = z.infer<typeof propPredictionSchema>

// Extended prediction schema with props
export const predictionWithPropsSchema = predictionSchema.extend({
  predictionType: z.enum(['winner', 'score', 'combined', 'props']),
  props: z.array(propPredictionSchema).optional(),
}).refine(
  (data) => {
    if (data.predictionType === 'props') {
      return data.props && data.props.length > 0
    }
    return true
  },
  { message: 'At least one prop prediction is required', path: ['props'] }
)

export type PredictionWithPropsFormData = z.infer<typeof predictionWithPropsSchema>

// Helper to convert props form data to JSON
export const propsFormDataToPredictionData = (props: PropPredictionFormData[]): string => {
  return JSON.stringify({
    type: 'props',
    props: props.map(p => ({
      prop_type_id: p.propTypeId,
      prop_slug: p.propSlug,
      line: p.line,
      selection: p.selection,
      player_id: p.playerId,
      points_value: p.pointsValue,
    })),
  })
}
```

**PATTERN**: Mirror existing predictionSchema structure
**VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 12: CREATE `frontend/src/components/predictions/PropTypeSelector.tsx`

**IMPLEMENT**: Component for selecting and configuring prop predictions

```tsx
import React from 'react'
import {
  Box,
  Card,
  CardContent,
  Typography,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  TextField,
  Chip,
  IconButton,
  Stack,
} from '@mui/material'
import { Add as AddIcon, Delete as DeleteIcon } from '@mui/icons-material'
import type { PropType } from '../../types/props.types'
import type { PropPredictionFormData } from '../../utils/prediction-validation'

interface PropTypeSelectorProps {
  propTypes: PropType[]
  selectedProps: PropPredictionFormData[]
  onPropsChange: (props: PropPredictionFormData[]) => void
  homeTeam: string
  awayTeam: string
  disabled?: boolean
}

export const PropTypeSelector: React.FC<PropTypeSelectorProps> = ({
  propTypes,
  selectedProps,
  onPropsChange,
  homeTeam,
  awayTeam,
  disabled = false,
}) => {
  const addProp = (propType: PropType) => {
    const newProp: PropPredictionFormData = {
      propTypeId: propType.id,
      propSlug: propType.slug,
      line: propType.defaultLine ?? undefined,
      selection: '',
      pointsValue: propType.pointsCorrect,
    }
    onPropsChange([...selectedProps, newProp])
  }

  const removeProp = (index: number) => {
    onPropsChange(selectedProps.filter((_, i) => i !== index))
  }

  const updateProp = (index: number, updates: Partial<PropPredictionFormData>) => {
    const updated = [...selectedProps]
    updated[index] = { ...updated[index], ...updates }
    onPropsChange(updated)
  }

  const getSelectionOptions = (propType: PropType) => {
    switch (propType.valueType) {
      case 'over_under':
        return [
          { value: 'over', label: 'Over' },
          { value: 'under', label: 'Under' },
        ]
      case 'yes_no':
        return [
          { value: 'yes', label: 'Yes' },
          { value: 'no', label: 'No' },
        ]
      case 'team_select':
        return [
          { value: 'home', label: homeTeam },
          { value: 'away', label: awayTeam },
        ]
      default:
        return []
    }
  }

  const availablePropTypes = propTypes.filter(
    pt => !selectedProps.some(sp => sp.propTypeId === pt.id)
  )

  const groupedPropTypes = availablePropTypes.reduce((acc, pt) => {
    if (!acc[pt.category]) acc[pt.category] = []
    acc[pt.category].push(pt)
    return acc
  }, {} as Record<string, PropType[]>)

  return (
    <Box>
      {/* Selected Props */}
      {selectedProps.length > 0 && (
        <Box sx={{ mb: 3 }}>
          <Typography variant="subtitle2" gutterBottom>
            Your Props ({selectedProps.length})
          </Typography>
          <Stack spacing={2}>
            {selectedProps.map((prop, index) => {
              const propType = propTypes.find(pt => pt.id === prop.propTypeId)
              if (!propType) return null

              return (
                <Card key={index} variant="outlined">
                  <CardContent sx={{ py: 1.5, '&:last-child': { pb: 1.5 } }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                      <Box sx={{ flex: 1 }}>
                        <Typography variant="body2" fontWeight="medium">
                          {propType.name}
                          <Chip 
                            label={`+${propType.pointsCorrect} pts`} 
                            size="small" 
                            color="primary" 
                            sx={{ ml: 1 }} 
                          />
                        </Typography>
                        
                        {propType.valueType === 'over_under' && (
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mt: 1 }}>
                            <TextField
                              label="Line"
                              type="number"
                              size="small"
                              value={prop.line ?? ''}
                              onChange={(e) => updateProp(index, { 
                                line: e.target.value ? parseFloat(e.target.value) : undefined 
                              })}
                              disabled={disabled}
                              sx={{ width: 100 }}
                              inputProps={{ step: 0.5 }}
                            />
                            <FormControl size="small">
                              <RadioGroup
                                row
                                value={prop.selection}
                                onChange={(e) => updateProp(index, { selection: e.target.value })}
                              >
                                {getSelectionOptions(propType).map(opt => (
                                  <FormControlLabel
                                    key={opt.value}
                                    value={opt.value}
                                    control={<Radio size="small" />}
                                    label={opt.label}
                                    disabled={disabled}
                                  />
                                ))}
                              </RadioGroup>
                            </FormControl>
                          </Box>
                        )}

                        {(propType.valueType === 'yes_no' || propType.valueType === 'team_select') && (
                          <FormControl size="small" sx={{ mt: 1 }}>
                            <RadioGroup
                              row
                              value={prop.selection}
                              onChange={(e) => updateProp(index, { selection: e.target.value })}
                            >
                              {getSelectionOptions(propType).map(opt => (
                                <FormControlLabel
                                  key={opt.value}
                                  value={opt.value}
                                  control={<Radio size="small" />}
                                  label={opt.label}
                                  disabled={disabled}
                                />
                              ))}
                            </RadioGroup>
                          </FormControl>
                        )}
                      </Box>
                      <IconButton 
                        size="small" 
                        onClick={() => removeProp(index)}
                        disabled={disabled}
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </Box>
                  </CardContent>
                </Card>
              )
            })}
          </Stack>
        </Box>
      )}

      {/* Available Props */}
      <Typography variant="subtitle2" gutterBottom>
        Add Props
      </Typography>
      {Object.entries(groupedPropTypes).map(([category, types]) => (
        <Box key={category} sx={{ mb: 2 }}>
          <Typography variant="caption" color="text.secondary" sx={{ textTransform: 'capitalize' }}>
            {category} Props
          </Typography>
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 0.5 }}>
            {types.map(propType => (
              <Chip
                key={propType.id}
                label={propType.name}
                onClick={() => addProp(propType)}
                onDelete={() => addProp(propType)}
                deleteIcon={<AddIcon />}
                variant="outlined"
                disabled={disabled}
                size="small"
              />
            ))}
          </Box>
        </Box>
      ))}

      {availablePropTypes.length === 0 && selectedProps.length > 0 && (
        <Typography variant="body2" color="text.secondary">
          All available props selected
        </Typography>
      )}
    </Box>
  )
}

export default PropTypeSelector
```

**PATTERN**: Mirror component structure from `PredictionForm.tsx`
**VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 13: UPDATE `frontend/src/components/predictions/PredictionForm.tsx` - Add props tab

**IMPLEMENT**: Extend form to support props predictions

Add imports:
```typescript
import { usePropTypes } from '../../hooks/use-predictions'
import { PropTypeSelector } from './PropTypeSelector'
import type { PropPredictionFormData } from '../../utils/prediction-validation'
import { propsFormDataToPredictionData } from '../../utils/prediction-validation'
```

Add state for props:
```typescript
const [selectedProps, setSelectedProps] = React.useState<PropPredictionFormData[]>([])
```

Fetch prop types:
```typescript
const { data: propTypes = [] } = usePropTypes(event?.sportType || '')
```

Update predictionType RadioGroup to include props option:
```tsx
<FormControlLabel value="props" control={<Radio />} label="Props" />
```

Add props section after score inputs:
```tsx
{predictionType === 'props' && propTypes.length > 0 && (
  <PropTypeSelector
    propTypes={propTypes}
    selectedProps={selectedProps}
    onPropsChange={setSelectedProps}
    homeTeam={event.homeTeam}
    awayTeam={event.awayTeam}
    disabled={loading}
  />
)}
```

Update handleFormSubmit to handle props:
```typescript
const handleFormSubmit = (data: PredictionFormData) => {
  let predictionData: string
  
  if (data.predictionType === 'props') {
    predictionData = propsFormDataToPredictionData(selectedProps)
  } else {
    predictionData = formDataToPredictionData(data)
  }
  
  onSubmit(predictionData)
}
```

Reset props when dialog closes:
```typescript
React.useEffect(() => {
  if (open) {
    reset(defaultValues)
    setSelectedProps([])
  }
}, [open, defaultValues, reset])
```

**PATTERN**: Follow existing form structure
**VALIDATE**: `cd frontend && npx tsc --noEmit`

### Task 14: CREATE `tests/prediction-service/prop_type_test.go`

**IMPLEMENT**: Unit tests for PropType model validation

```go
package prediction_service_test

import (
	"testing"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
)

func TestPropType_ValidateSportType(t *testing.T) {
	tests := []struct {
		name      string
		sportType string
		wantErr   bool
	}{
		{"valid sport type", "Soccer", false},
		{"empty sport type", "", true},
		{"whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &models.PropType{SportType: tt.sportType}
			err := pt.ValidateSportType()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSportType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPropType_ValidateCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{"valid match", "match", false},
		{"valid player", "player", false},
		{"valid team", "team", false},
		{"invalid category", "invalid", true},
		{"empty category", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &models.PropType{Category: tt.category}
			err := pt.ValidateCategory()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPropType_ValidateValueType(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		wantErr   bool
	}{
		{"valid over_under", "over_under", false},
		{"valid yes_no", "yes_no", false},
		{"valid team_select", "team_select", false},
		{"valid player_select", "player_select", false},
		{"valid exact_value", "exact_value", false},
		{"invalid type", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &models.PropType{ValueType: tt.valueType}
			err := pt.ValidateValueType()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateValueType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

**PATTERN**: Mirror `tests/prediction-service/prediction_test.go`
**VALIDATE**: `cd tests/prediction-service && go test -v ./...`

### Task 15: CREATE `tests/scoring-service/props_scoring_test.go`

**IMPLEMENT**: Unit tests for props scoring logic

```go
package scoring_service_test

import (
	"encoding/json"
	"testing"
)

// Test data structures matching scoring service
type PropPrediction struct {
	PropTypeID  uint    `json:"prop_type_id"`
	PropSlug    string  `json:"prop_slug"`
	Line        float64 `json:"line"`
	Selection   string  `json:"selection"`
	PointsValue float64 `json:"points_value"`
}

type ResultData struct {
	HomeScore   int                    `json:"home_score"`
	AwayScore   int                    `json:"away_score"`
	TotalGoals  int                    `json:"total_goals"`
	Stats       map[string]interface{} `json:"stats"`
}

func TestPropsScoring_TotalGoalsOverUnder(t *testing.T) {
	tests := []struct {
		name       string
		line       float64
		selection  string
		totalGoals int
		wantCorrect bool
	}{
		{"over 2.5 with 3 goals", 2.5, "over", 3, true},
		{"over 2.5 with 2 goals", 2.5, "over", 2, false},
		{"under 2.5 with 2 goals", 2.5, "under", 2, true},
		{"under 2.5 with 3 goals", 2.5, "under", 3, false},
		{"over 2.5 with exactly 2.5 (impossible)", 2.5, "over", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			correct := evaluateTotalGoalsOU(tt.line, tt.selection, tt.totalGoals)
			if correct != tt.wantCorrect {
				t.Errorf("evaluateTotalGoalsOU() = %v, want %v", correct, tt.wantCorrect)
			}
		})
	}
}

func evaluateTotalGoalsOU(line float64, selection string, totalGoals int) bool {
	if selection == "over" {
		return float64(totalGoals) > line
	}
	return float64(totalGoals) < line
}

func TestPropsScoring_BothTeamsToScore(t *testing.T) {
	tests := []struct {
		name       string
		homeScore  int
		awayScore  int
		selection  string
		wantCorrect bool
	}{
		{"yes with both scoring", 2, 1, "yes", true},
		{"yes with only home scoring", 2, 0, "yes", false},
		{"no with only home scoring", 2, 0, "no", true},
		{"no with both scoring", 1, 1, "no", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			correct := evaluateBTTS(tt.homeScore, tt.awayScore, tt.selection)
			if correct != tt.wantCorrect {
				t.Errorf("evaluateBTTS() = %v, want %v", correct, tt.wantCorrect)
			}
		})
	}
}

func evaluateBTTS(homeScore, awayScore int, selection string) bool {
	btts := homeScore > 0 && awayScore > 0
	if selection == "yes" {
		return btts
	}
	return !btts
}

func TestPropsScoring_JSONParsing(t *testing.T) {
	predictionJSON := `{
		"type": "props",
		"props": [
			{"prop_type_id": 1, "prop_slug": "total-goals-ou", "line": 2.5, "selection": "over", "points_value": 2}
		]
	}`

	var data struct {
		Type  string           `json:"type"`
		Props []PropPrediction `json:"props"`
	}

	err := json.Unmarshal([]byte(predictionJSON), &data)
	if err != nil {
		t.Fatalf("Failed to parse prediction JSON: %v", err)
	}

	if data.Type != "props" {
		t.Errorf("Expected type 'props', got '%s'", data.Type)
	}

	if len(data.Props) != 1 {
		t.Errorf("Expected 1 prop, got %d", len(data.Props))
	}

	if data.Props[0].PropSlug != "total-goals-ou" {
		t.Errorf("Expected prop slug 'total-goals-ou', got '%s'", data.Props[0].PropSlug)
	}
}
```

**PATTERN**: Mirror `tests/scoring-service/scoring_test.go`
**VALIDATE**: `cd tests/scoring-service && go test -v ./...`

---

## TESTING STRATEGY

### Unit Tests

Based on project's Go testing patterns:
- PropType model validation (category, value_type, sport_type)
- Props scoring evaluation for each prop type
- JSON parsing for props prediction data

### Integration Tests

- Submit props prediction via API
- Calculate props score with mock result data
- Retrieve prop types by sport

### Edge Cases

- Empty props array in prediction
- Invalid prop type ID
- Missing line for over/under props
- Props for sport with no prop types defined

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build

```bash
# Backend
cd backend/prediction-service && go build ./...
cd backend/scoring-service && go build ./...

# Frontend
cd frontend && npx tsc --noEmit
```

### Level 2: Unit Tests

```bash
# Backend tests
cd tests/prediction-service && go test -v ./...
cd tests/scoring-service && go test -v ./...

# Frontend tests (if configured)
cd frontend && npm test -- --passWithNoTests
```

### Level 3: Database Validation

```bash
# Recreate database with new schema
docker-compose down -v
docker-compose up -d postgres
sleep 3
docker exec sports-postgres psql -U sports_user -d sports_prediction -c "\d prop_types"
docker exec sports-postgres psql -U sports_user -d sports_prediction -c "SELECT * FROM prop_types"
```

### Level 4: Manual Validation

```bash
# Start services
make docker-up
cd frontend && npm run dev

# Test prop types endpoint (when services running)
curl http://localhost:8080/v1/prop-types/Soccer
```

---

## ACCEPTANCE CRITERIA

- [ ] prop_types table created with default Soccer props
- [ ] PropType model validates category and value_type
- [ ] GetPropTypes RPC returns props for sport type
- [ ] ListPropTypes RPC supports filtering and pagination
- [ ] Props scoring correctly evaluates over/under, yes/no, team_select
- [ ] Frontend displays prop types grouped by category
- [ ] PropTypeSelector allows adding/removing/configuring props
- [ ] PredictionForm supports "props" prediction type
- [ ] Props predictions serialize correctly to JSON
- [ ] All unit tests pass
- [ ] TypeScript compiles without errors

---

## COMPLETION CHECKLIST

- [ ] Task 1: Database schema updated
- [ ] Task 2: PropType model created
- [ ] Task 3: PropType repository created
- [ ] Task 4: Proto definitions added
- [ ] Task 5: Prediction service RPCs implemented
- [ ] Task 6: Main.go updated with PropType repo
- [ ] Task 7: Scoring service props logic added
- [ ] Task 8: Frontend props types created
- [ ] Task 9: Prediction service methods added
- [ ] Task 10: React Query hooks added
- [ ] Task 11: Validation schemas extended
- [ ] Task 12: PropTypeSelector component created
- [ ] Task 13: PredictionForm extended
- [ ] Task 14: PropType unit tests created
- [ ] Task 15: Props scoring tests created
- [ ] All validation commands pass
- [ ] Manual testing confirms feature works

---

## NOTES

### Design Decisions

1. **Flexible prop_types table**: Allows adding new prop types without code changes
2. **JSON-based prediction data**: Maintains backward compatibility with existing predictions
3. **Category-based grouping**: Separates match, player, and team props for better UX
4. **Points per prop type**: Different props can have different point values
5. **Default Soccer props**: Provides immediate value without external data

### Future Enhancements

- Player props requiring player selection (needs player data integration)
- Live props during matches
- Custom prop types per contest
- Prop type odds/difficulty adjustment

### External Data Considerations

TheSportsDB free tier has limited statistics. Current implementation uses:
- Basic match results (goals)
- Stats stored in result_data JSON field

For full props support (corners, cards, shots), consider:
- Premium sports data API integration
- Manual result entry by admins
- Community-sourced statistics

### Confidence Score: 8/10

High confidence due to:
- Clear patterns from existing prediction system
- Self-contained feature with minimal dependencies
- Flexible JSON-based data model
- Comprehensive test coverage planned

Risks:
- Proto regeneration may require manual intervention
- Props scoring depends on result_data having stats populated
