# Feature: Score Prediction Schema System with Seeded Predictions

## Feature Description

Implement a prediction schema system that allows contests to define how predictions are made. Initially support "exact score" schema where users select from predefined score options (1-0, 0-1, 2-0, etc.) or enter custom scores. Add prediction seeding to `make seed` commands so the frontend displays realistic prediction data.

## User Story

As a contest participant
I want to predict exact match scores from predefined options or custom input
So that I can make predictions in a structured, user-friendly way

As a developer
I want the seeding system to generate realistic predictions
So that the frontend displays populated data during development

## Problem Statement

1. **No predictions visible on frontend** - Seeding system doesn't generate predictions
2. **No prediction schema system** - Contests don't define how predictions should be made
3. **Limited prediction UX** - Users need a structured way to predict scores

## Solution Statement

1. Add `prediction_schema` field to Contest model defining prediction type and options
2. Implement "exact_score" schema with predefined score options
3. Update seeder to generate predictions for all seeded contests and users
4. Update frontend prediction form to show score selection UI based on contest schema

## Feature Metadata

**Feature Type**: Enhancement + Bug Fix
**Estimated Complexity**: Medium
**Primary Systems Affected**: 
- Contest Service (schema definition)
- Prediction Service (schema validation)
- Seeder (prediction generation)
- Frontend (prediction form UI)

**Dependencies**: Existing prediction and contest services

---

## CONTEXT REFERENCES

### Relevant Codebase Files - MUST READ BEFORE IMPLEMENTING!

**Backend Models:**
- `backend/contest-service/internal/models/contest.go` (lines 1-186) - Contest model structure, add prediction_schema field
- `backend/prediction-service/internal/models/prediction.go` (lines 1-132) - Prediction model, PredictionData JSON field
- `backend/prediction-service/internal/models/event.go` (lines 1-160) - Event model with ResultData

**Backend Services:**
- `backend/contest-service/internal/service/contest_service.go` (lines 1-595) - Contest CRUD, update to handle schema
- `backend/prediction-service/internal/service/prediction_service.go` (lines 1-735) - Prediction submission logic

**Seeder:**
- `backend/shared/seeder/coordinator.go` (lines 1-1047) - Main seeding orchestration
- `backend/shared/seeder/factory.go` (lines 1-239) - Data generation functions
- `backend/shared/seeder/models.go` (lines 1-296) - Seeder model definitions

**Frontend:**
- `frontend/src/components/predictions/PredictionForm.tsx` (lines 1-150) - Prediction submission form
- `frontend/src/pages/PredictionsPage.tsx` (lines 1-149) - Predictions page
- `frontend/src/types/contest.types.ts` (lines 1-149) - Contest TypeScript types
- `frontend/src/types/prediction.types.ts` (lines 1-111) - Prediction TypeScript types

**Proto Definitions:**
- `backend/proto/contest/contest.proto` - Contest message definitions
- `backend/proto/prediction/prediction.proto` - Prediction message definitions

### New Files to Create

- None (all modifications to existing files)

### Relevant Documentation - READ BEFORE IMPLEMENTING!

- [GORM JSON Fields](https://gorm.io/docs/data_types.html#JSON)
  - Section: JSON data type usage
  - Why: prediction_schema will be stored as JSON
  
- [React Hook Form](https://react-hook-form.com/docs/useform)
  - Section: Dynamic form fields
  - Why: Score selection UI needs dynamic rendering

- [Ant Design Select](https://ant.design/components/select)
  - Section: Select component API
  - Why: Score dropdown selection

### Patterns to Follow

**JSON Field Pattern (from existing code):**
```go
// backend/contest-service/internal/models/contest.go
type Contest struct {
    Rules datatypes.JSON `gorm:"type:jsonb" json:"rules"`
}

// backend/prediction-service/internal/models/prediction.go
type Prediction struct {
    PredictionData datatypes.JSON `gorm:"type:jsonb" json:"prediction_data"`
}
```

**Seeder Pattern (from coordinator.go):**
```go
func (c *Coordinator) seedPredictions(ctx context.Context) error {
    // Batch insert pattern
    predictions := make([]*Prediction, 0, batchSize)
    for i := 0; i < count; i++ {
        predictions = append(predictions, generatePrediction())
        if len(predictions) >= batchSize {
            if err := c.db.Create(&predictions).Error; err != nil {
                return err
            }
            predictions = predictions[:0]
        }
    }
    return nil
}
```

**Frontend Form Pattern (from PredictionForm.tsx):**
```tsx
const { control, handleSubmit } = useForm<PredictionFormData>({
  resolver: zodResolver(predictionSchema),
});

<Controller
  name="predictionData"
  control={control}
  render={({ field }) => (
    <Form.Item label="Prediction">
      <Input {...field} />
    </Form.Item>
  )}
/>
```

---

## IMPLEMENTATION PLAN

### Phase 1: Backend Schema Foundation

Add prediction schema support to Contest model and proto definitions. Define schema structure for exact score predictions.

**Tasks:**
- Add `prediction_schema` JSON field to Contest model
- Update Contest proto with prediction_schema field
- Define schema structure: `{"type": "exact_score", "options": ["1-0", "0-1", ...]}`

### Phase 2: Seeder Enhancement

Update seeder to generate realistic predictions for all contests and users.

**Tasks:**
- Add prediction generation to seeder coordinator
- Generate predictions for 60-80% of users per contest
- Use realistic score distributions (common scores more frequent)
- Link predictions to existing events

### Phase 3: Frontend Schema-Based UI

Update prediction form to render score selection based on contest schema.

**Tasks:**
- Add schema field to Contest TypeScript types
- Update PredictionForm to detect schema type
- Render score selector for exact_score schema
- Add custom score input option

### Phase 4: Testing & Validation

Verify seeded data appears correctly and prediction submission works.

**Tasks:**
- Test seeding generates predictions
- Verify frontend displays predictions
- Test score selection and submission
- Validate prediction data structure

---

## STEP-BY-STEP TASKS

### UPDATE backend/contest-service/internal/models/contest.go

- **ADD**: `PredictionSchema datatypes.JSON` field to Contest struct after Rules field
- **PATTERN**: Mirror Rules field pattern (line 23): `Rules datatypes.JSON \`gorm:"type:jsonb" json:"rules"\``
- **IMPORTS**: Already has `gorm.io/datatypes`
- **GOTCHA**: Use `jsonb` type for PostgreSQL, not `json`
- **VALIDATE**: `cd backend/contest-service && go build ./...`

```go
// Add after Rules field (around line 24)
PredictionSchema datatypes.JSON `gorm:"type:jsonb" json:"prediction_schema"`
```

### UPDATE backend/proto/contest/contest.proto

- **ADD**: `string prediction_schema = 14;` to Contest message
- **PATTERN**: Follow existing JSON field pattern (rules field)
- **GOTCHA**: Use next available field number (14)
- **VALIDATE**: `cd backend && make proto`

### UPDATE backend/shared/seeder/factory.go

- **ADD**: `GenerateDefaultPredictionSchema()` function
- **IMPLEMENT**: Return exact_score schema with predefined options
- **PATTERN**: Mirror other Generate functions in file
- **VALIDATE**: `cd backend/shared && go build ./...`

```go
// Add after GenerateContests function
func GenerateDefaultPredictionSchema() map[string]interface{} {
    return map[string]interface{}{
        "type": "exact_score",
        "options": []string{
            "1-0", "0-1", "2-0", "0-2", "2-1", "1-2",
            "3-0", "0-3", "3-1", "1-3", "3-2", "2-3",
            "0-0", "1-1", "2-2", "3-3",
        },
        "allow_custom": true,
    }
}
```

### UPDATE backend/shared/seeder/factory.go - GenerateContests

- **UPDATE**: Add prediction_schema to generated contests
- **IMPLEMENT**: Marshal schema to JSON and assign to PredictionSchema field
- **PATTERN**: Mirror how Rules field is generated (lines 150-160)
- **VALIDATE**: `cd backend/shared && go build ./...`

```go
// In GenerateContests function, after rules generation
schemaData := GenerateDefaultPredictionSchema()
schemaJSON, _ := json.Marshal(schemaData)
contest.PredictionSchema = datatypes.JSON(schemaJSON)
```

### UPDATE backend/shared/seeder/coordinator.go - seedPredictions

- **REFACTOR**: Current seedPredictions function (lines 600-700)
- **IMPLEMENT**: Generate predictions for 60-80% of users per contest
- **PATTERN**: Use batch insert pattern from seedContests (lines 400-450)
- **IMPORTS**: Add `math/rand` for random selection
- **GOTCHA**: Ensure predictions link to valid event_id, user_id, contest_id
- **VALIDATE**: `cd backend/shared && go test ./...`

```go
// Replace existing seedPredictions implementation
func (c *Coordinator) seedPredictions(ctx context.Context) error {
    log.Println("Seeding predictions...")
    
    // Get all users, contests, and events
    var users []User
    var contests []Contest
    var events []Event
    
    if err := c.db.Find(&users).Error; err != nil {
        return err
    }
    if err := c.db.Find(&contests).Error; err != nil {
        return err
    }
    if err := c.db.Find(&events).Error; err != nil {
        return err
    }
    
    if len(users) == 0 || len(contests) == 0 || len(events) == 0 {
        return fmt.Errorf("need users, contests, and events before seeding predictions")
    }
    
    predictions := make([]*Prediction, 0, 1000)
    scoreOptions := []string{"1-0", "0-1", "2-0", "0-2", "2-1", "1-2", "3-0", "0-3", "3-1", "1-3", "3-2", "2-3", "0-0", "1-1", "2-2"}
    
    // For each contest, 60-80% of users make predictions
    for _, contest := range contests {
        // Get events for this contest's sport
        contestEvents := make([]Event, 0)
        for _, event := range events {
            if event.SportType == contest.SportType {
                contestEvents = append(contestEvents, event)
            }
        }
        
        if len(contestEvents) == 0 {
            continue
        }
        
        // Random 60-80% of users participate
        participationRate := 0.6 + rand.Float64()*0.2
        numParticipants := int(float64(len(users)) * participationRate)
        
        // Shuffle users
        shuffledUsers := make([]User, len(users))
        copy(shuffledUsers, users)
        rand.Shuffle(len(shuffledUsers), func(i, j int) {
            shuffledUsers[i], shuffledUsers[j] = shuffledUsers[j], shuffledUsers[i]
        })
        
        for i := 0; i < numParticipants && i < len(shuffledUsers); i++ {
            user := shuffledUsers[i]
            
            // Each user predicts 3-8 random events
            numPredictions := 3 + rand.Intn(6)
            if numPredictions > len(contestEvents) {
                numPredictions = len(contestEvents)
            }
            
            // Shuffle events
            shuffledEvents := make([]Event, len(contestEvents))
            copy(shuffledEvents, contestEvents)
            rand.Shuffle(len(shuffledEvents), func(i, j int) {
                shuffledEvents[i], shuffledEvents[j] = shuffledEvents[j], shuffledEvents[i]
            })
            
            for j := 0; j < numPredictions; j++ {
                event := shuffledEvents[j]
                
                // Generate score prediction
                score := scoreOptions[rand.Intn(len(scoreOptions))]
                parts := strings.Split(score, "-")
                homeScore, _ := strconv.Atoi(parts[0])
                awayScore, _ := strconv.Atoi(parts[1])
                
                predictionData := map[string]interface{}{
                    "home_score": homeScore,
                    "away_score": awayScore,
                    "score_string": score,
                }
                predictionJSON, _ := json.Marshal(predictionData)
                
                prediction := &Prediction{
                    UserID:         user.ID,
                    ContestID:      contest.ID,
                    EventID:        event.ID,
                    PredictionData: datatypes.JSON(predictionJSON),
                    Status:         "pending",
                    SubmittedAt:    time.Now().Add(-time.Duration(rand.Intn(72)) * time.Hour),
                }
                
                predictions = append(predictions, prediction)
                
                // Batch insert
                if len(predictions) >= 500 {
                    if err := c.db.Create(&predictions).Error; err != nil {
                        return fmt.Errorf("failed to insert predictions batch: %w", err)
                    }
                    log.Printf("Inserted %d predictions", len(predictions))
                    predictions = predictions[:0]
                }
            }
        }
    }
    
    // Insert remaining
    if len(predictions) > 0 {
        if err := c.db.Create(&predictions).Error; err != nil {
            return fmt.Errorf("failed to insert final predictions batch: %w", err)
        }
        log.Printf("Inserted final %d predictions", len(predictions))
    }
    
    log.Printf("✓ Predictions seeded successfully")
    return nil
}
```

### UPDATE frontend/src/types/contest.types.ts

- **ADD**: `predictionSchema` field to Contest interface
- **PATTERN**: Follow existing optional JSON fields pattern
- **VALIDATE**: `cd frontend && npm run build`

```typescript
// Add to Contest interface (around line 20)
predictionSchema?: {
  type: string;
  options?: string[];
  allow_custom?: boolean;
};
```

### UPDATE frontend/src/components/predictions/PredictionForm.tsx

- **REFACTOR**: Replace generic input with schema-based rendering
- **ADD**: Score selector component for exact_score schema
- **PATTERN**: Use Ant Design Select component with Controller
- **IMPORTS**: Add `Select` from 'antd'
- **GOTCHA**: Parse score string to home_score/away_score for submission
- **VALIDATE**: `cd frontend && npm run build`

```tsx
// Add after imports
import { Select } from 'antd';

// Inside form, replace prediction input with:
{contest?.predictionSchema?.type === 'exact_score' ? (
  <>
    <Controller
      name="scoreString"
      control={control}
      render={({ field }) => (
        <Form.Item label="Predict Score" required>
          <Select
            {...field}
            placeholder="Select score"
            options={[
              ...(contest.predictionSchema.options || []).map(score => ({
                label: score,
                value: score,
              })),
              ...(contest.predictionSchema.allow_custom ? [
                { label: 'Custom...', value: 'custom' }
              ] : []),
            ]}
          />
        </Form.Item>
      )}
    />
    {watchScoreString === 'custom' && (
      <Space>
        <Controller
          name="homeScore"
          control={control}
          render={({ field }) => (
            <Form.Item label="Home">
              <InputNumber {...field} min={0} max={20} />
            </Form.Item>
          )}
        />
        <Controller
          name="awayScore"
          control={control}
          render={({ field }) => (
            <Form.Item label="Away">
              <InputNumber {...field} min={0} max={20} />
            </Form.Item>
          )}
        />
      </Space>
    )}
  </>
) : (
  <Controller
    name="predictionData"
    control={control}
    render={({ field }) => (
      <Form.Item label="Prediction Data">
        <Input.TextArea {...field} placeholder="Enter prediction data as JSON" />
      </Form.Item>
    )}
  />
)}
```

### UPDATE frontend/src/components/predictions/PredictionForm.tsx - onSubmit

- **UPDATE**: Transform score selection to prediction_data format
- **IMPLEMENT**: Parse scoreString or custom scores to JSON structure
- **PATTERN**: Mirror existing form submission pattern
- **VALIDATE**: Test in browser

```tsx
// Update onSubmit handler
const onSubmit = async (data: PredictionFormData) => {
  try {
    let predictionData: any;
    
    if (contest?.predictionSchema?.type === 'exact_score') {
      if (data.scoreString === 'custom') {
        predictionData = {
          home_score: data.homeScore,
          away_score: data.awayScore,
          score_string: `${data.homeScore}-${data.awayScore}`,
        };
      } else {
        const [home, away] = data.scoreString.split('-').map(Number);
        predictionData = {
          home_score: home,
          away_score: away,
          score_string: data.scoreString,
        };
      }
    } else {
      predictionData = JSON.parse(data.predictionData);
    }
    
    await submitPrediction({
      contestId: data.contestId,
      eventId: data.eventId,
      predictionData: JSON.stringify(predictionData),
    });
    
    message.success('Prediction submitted successfully');
    onSuccess?.();
  } catch (error) {
    message.error('Failed to submit prediction');
  }
};
```

### REGENERATE backend proto files

- **EXECUTE**: Generate updated gRPC code with new prediction_schema field
- **PATTERN**: Use existing proto generation script
- **VALIDATE**: `cd backend && make proto && go build ./...`

```bash
cd backend
make proto
```

### UPDATE backend/contest-service - rebuild

- **EXECUTE**: Rebuild contest service with new schema field
- **VALIDATE**: Service compiles and starts
- **GOTCHA**: May need to restart Docker containers

```bash
cd backend/contest-service
go mod tidy
go build ./...
```

---

## TESTING STRATEGY

### Unit Tests

**Backend:**
- Test `GenerateDefaultPredictionSchema()` returns correct structure
- Test prediction seeding generates valid data
- Test contest creation with prediction_schema

**Frontend:**
- Test PredictionForm renders score selector when schema present
- Test score selection transforms to correct prediction_data format
- Test custom score input works

### Integration Tests

**Seeding:**
```bash
# Clear database and reseed
make docker-down
make dev
make seed-small
```

**Frontend Display:**
1. Start services: `make docker-services`
2. Open http://localhost:3000/predictions
3. Verify predictions are visible
4. Verify score selector appears for contests with schema

### Edge Cases

- Contest without prediction_schema (fallback to generic input)
- Custom score input validation (0-20 range)
- Empty prediction options array
- Invalid score string format

---

## VALIDATION COMMANDS

### Level 1: Syntax & Build

```bash
# Backend build
cd backend/shared && go build ./...
cd backend/contest-service && go build ./...
cd backend/prediction-service && go build ./...

# Frontend build
cd frontend && npm run build
```

### Level 2: Proto Generation

```bash
cd backend && make proto
```

### Level 3: Database Seeding

```bash
# Full reseed
make docker-down
make dev
sleep 5
make seed-small
```

### Level 4: Manual Validation

```bash
# Start all services
make docker-services

# Check predictions exist in database
docker exec -it sports-postgres psql -U sports_user -d sports_prediction -c "SELECT COUNT(*) FROM predictions;"

# Expected: > 0 predictions

# Frontend check:
# 1. Open http://localhost:3000/predictions
# 2. Verify predictions list shows data
# 3. Click "Make Prediction" 
# 4. Verify score selector appears
# 5. Select score and submit
# 6. Verify success message
```

### Level 5: API Validation

```bash
# Get predictions via API
curl http://localhost:8080/api/predictions/user/1 | jq

# Expected: Array of predictions with prediction_data containing scores
```

---

## ACCEPTANCE CRITERIA

- [ ] Contest model has prediction_schema JSON field
- [ ] Seeder generates predictions for 60-80% of users per contest
- [ ] Seeded predictions have realistic score distributions
- [ ] Frontend displays seeded predictions on Predictions page
- [ ] Prediction form shows score selector for exact_score schema
- [ ] Score selector includes predefined options (1-0, 0-1, etc.)
- [ ] Custom score input available when allow_custom is true
- [ ] Score selection transforms to correct prediction_data JSON format
- [ ] All validation commands pass
- [ ] No regressions in existing prediction functionality

---

## COMPLETION CHECKLIST

- [ ] Contest model updated with prediction_schema field
- [ ] Proto definitions regenerated
- [ ] Seeder generates predictions with realistic data
- [ ] Frontend types updated with schema field
- [ ] PredictionForm renders schema-based UI
- [ ] Score selection works and submits correctly
- [ ] Database seeding produces visible predictions
- [ ] All build commands pass
- [ ] Manual testing confirms predictions visible
- [ ] API returns predictions with correct data structure

---

## NOTES

**Design Decisions:**
- Using JSON field for schema flexibility (future schema types: over/under, props, etc.)
- Predefined score options cover most common football scores
- Custom input allows any score for edge cases
- 60-80% participation rate creates realistic data distribution

**Future Enhancements:**
- Additional schema types (over/under, both_teams_to_score, etc.)
- Schema validation in prediction service
- Schema-specific scoring algorithms
- UI for contest creators to define custom schemas

**Performance Considerations:**
- Batch insert for predictions (500 per batch)
- Index on (user_id, contest_id, event_id) for fast lookups
- JSON field allows flexible schema without migrations

**Confidence Score: 8/10**
- Clear implementation path with existing patterns
- Seeder complexity moderate but well-defined
- Frontend changes straightforward with Ant Design
- Risk: Seeder performance with large datasets (mitigated by batching)
