# Scoring Algorithms Implementation Guide

## Overview

This document provides implementation details for various scoring algorithms used in sports prediction contests, from basic point systems to advanced machine learning-based approaches.

## 1. Basic Scoring Systems

### Simple Points System
```go
type BasicScoring struct {
    CorrectPoints int
    WrongPoints   int
}

func (bs *BasicScoring) Calculate(prediction Prediction, result EventResult) int {
    if bs.IsCorrect(prediction, result) {
        return bs.CorrectPoints
    }
    return bs.WrongPoints
}
```

### Tiered Scoring System
```go
type TieredScoring struct {
    ExactMatch    int // Exact score prediction
    CorrectResult int // Correct outcome only
    CloseScore    int // Within 1 goal difference
    WrongResult   int // Completely wrong
}

func (ts *TieredScoring) CalculateScore(prediction ExactScore, result EventResult) int {
    predHome, predAway := prediction.HomeScore, prediction.AwayScore
    actualHome, actualAway := result.HomeScore, result.AwayScore
    
    // Exact match
    if predHome == actualHome && predAway == actualAway {
        return ts.ExactMatch
    }
    
    // Correct result (win/draw/loss)
    if ts.sameOutcome(predHome, predAway, actualHome, actualAway) {
        // Check if close score
        if abs(predHome-actualHome) <= 1 && abs(predAway-actualAway) <= 1 {
            return ts.CloseScore
        }
        return ts.CorrectResult
    }
    
    return ts.WrongResult
}

func (ts *TieredScoring) sameOutcome(ph, pa, ah, aa int) bool {
    predOutcome := getOutcome(ph, pa)
    actualOutcome := getOutcome(ah, aa)
    return predOutcome == actualOutcome
}
```

## 2. Confidence-Based Scoring

### Linear Confidence Scaling
```go
type ConfidenceScoring struct {
    BasePoints    int
    MinMultiplier float64 // 0.5 = 50% of base points minimum
    MaxMultiplier float64 // 2.0 = 200% of base points maximum
}

func (cs *ConfidenceScoring) Calculate(prediction Prediction, result EventResult) int {
    baseScore := cs.getBaseScore(prediction, result)
    if baseScore <= 0 {
        return 0 // No points for wrong predictions
    }
    
    confidence := prediction.Confidence / 100.0
    multiplier := cs.MinMultiplier + (cs.MaxMultiplier-cs.MinMultiplier)*confidence
    
    return int(float64(baseScore) * multiplier)
}
```

### Risk-Reward Confidence System
```go
type RiskRewardScoring struct {
    BasePoints int
}

func (rr *RiskRewardScoring) Calculate(prediction Prediction, result EventResult) int {
    confidence := prediction.Confidence / 100.0
    
    if rr.IsCorrect(prediction, result) {
        // Reward high confidence correct predictions
        return int(float64(rr.BasePoints) * (1.0 + confidence))
    } else {
        // Penalize high confidence wrong predictions
        penalty := int(float64(rr.BasePoints) * confidence * 0.5)
        return -penalty
    }
}
```

## 3. Difficulty-Based Scoring

### Odds-Based Difficulty
```go
type OddsBasedScoring struct {
    BasePoints int
}

func (obs *OddsBasedScoring) Calculate(prediction Prediction, result EventResult, odds EventOdds) int {
    if !obs.IsCorrect(prediction, result) {
        return 0
    }
    
    // Convert odds to implied probability
    impliedProb := obs.getImpliedProbability(prediction, odds)
    
    // Higher points for less likely outcomes
    difficultyMultiplier := 1.0 / impliedProb
    
    // Cap the multiplier to prevent extreme scores
    if difficultyMultiplier > 5.0 {
        difficultyMultiplier = 5.0
    }
    
    return int(float64(obs.BasePoints) * difficultyMultiplier)
}

func (obs *OddsBasedScoring) getImpliedProbability(prediction Prediction, odds EventOdds) float64 {
    switch prediction.Type {
    case MatchOutcome:
        outcome := prediction.Value.(MatchOutcome)
        switch outcome.Result {
        case "home_win":
            return 1.0 / odds.HomeWin
        case "away_win":
            return 1.0 / odds.AwayWin
        case "draw":
            return 1.0 / odds.Draw
        }
    }
    return 0.5 // Default 50% probability
}
```

### Historical Performance Difficulty
```go
type HistoricalDifficulty struct {
    BasePoints int
}

func (hd *HistoricalDifficulty) Calculate(prediction Prediction, result EventResult) int {
    if !hd.IsCorrect(prediction, result) {
        return 0
    }
    
    // Get historical accuracy for this prediction type
    accuracy := hd.getHistoricalAccuracy(prediction.Type, prediction.EventID)
    
    // Lower accuracy = higher difficulty = more points
    difficultyMultiplier := 1.0 / accuracy
    
    return int(float64(hd.BasePoints) * difficultyMultiplier)
}

func (hd *HistoricalDifficulty) getHistoricalAccuracy(predType PredictionType, eventID string) float64 {
    // Query database for historical accuracy of this prediction type
    // for similar events (same teams, competition, etc.)
    return 0.6 // Placeholder: 60% historical accuracy
}
```

## 4. Progressive Scoring Systems

### Streak Bonuses
```go
type StreakScoring struct {
    BasePoints   int
    StreakBonus  map[int]float64 // streak length -> multiplier
}

func NewStreakScoring() *StreakScoring {
    return &StreakScoring{
        BasePoints: 10,
        StreakBonus: map[int]float64{
            3:  1.1, // 10% bonus for 3 in a row
            5:  1.2, // 20% bonus for 5 in a row
            10: 1.5, // 50% bonus for 10 in a row
        },
    }
}

func (ss *StreakScoring) Calculate(prediction Prediction, result EventResult, userStreak int) int {
    baseScore := ss.getBaseScore(prediction, result)
    if baseScore <= 0 {
        return 0
    }
    
    multiplier := 1.0
    for streakLen, bonus := range ss.StreakBonus {
        if userStreak >= streakLen {
            multiplier = bonus
        }
    }
    
    return int(float64(baseScore) * multiplier)
}
```

### Category Mastery System
```go
type MasteryScoring struct {
    BasePoints      int
    MasteryLevels   map[string]float64 // category -> multiplier
}

func (ms *MasteryScoring) Calculate(prediction Prediction, result EventResult, userMastery UserMastery) int {
    baseScore := ms.getBaseScore(prediction, result)
    if baseScore <= 0 {
        return 0
    }
    
    category := ms.getCategory(prediction)
    masteryLevel := userMastery.GetLevel(category)
    
    multiplier := 1.0 + (masteryLevel * 0.1) // 10% per mastery level
    
    return int(float64(baseScore) * multiplier)
}

type UserMastery struct {
    Categories map[string]int // category -> level (0-10)
}

func (um *UserMastery) GetLevel(category string) int {
    if level, exists := um.Categories[category]; exists {
        return level
    }
    return 0
}
```

## 5. Advanced Scoring Algorithms

### Bayesian Scoring System
```go
type BayesianScoring struct {
    BasePoints int
    PriorBelief float64 // Prior probability of correct prediction
}

func (bs *BayesianScoring) Calculate(prediction Prediction, result EventResult, userHistory UserHistory) int {
    if !bs.IsCorrect(prediction, result) {
        return 0
    }
    
    // Calculate user's historical accuracy for this prediction type
    userAccuracy := userHistory.GetAccuracy(prediction.Type)
    
    // Bayesian update: combine prior belief with user's track record
    posterior := bs.bayesianUpdate(bs.PriorBelief, userAccuracy, prediction.Confidence/100.0)
    
    // Score based on how much the prediction improved our belief
    informationGain := posterior - bs.PriorBelief
    
    return int(float64(bs.BasePoints) * (1.0 + informationGain*2.0))
}

func (bs *BayesianScoring) bayesianUpdate(prior, likelihood, evidence float64) float64 {
    // Simplified Bayesian update
    numerator := likelihood * evidence
    denominator := likelihood*evidence + (1-likelihood)*(1-evidence)
    return numerator / denominator
}
```

### Machine Learning-Based Scoring
```go
type MLScoring struct {
    BasePoints int
    Model      PredictionModel
}

type PredictionModel interface {
    PredictDifficulty(event Event, prediction Prediction) float64
    PredictUserPerformance(user User, prediction Prediction) float64
}

func (ml *MLScoring) Calculate(prediction Prediction, result EventResult, event Event, user User) int {
    if !ml.IsCorrect(prediction, result) {
        return 0
    }
    
    // Use ML model to assess prediction difficulty
    difficulty := ml.Model.PredictDifficulty(event, prediction)
    
    // Use ML model to assess expected user performance
    expectedPerformance := ml.Model.PredictUserPerformance(user, prediction)
    
    // Score based on difficulty and user's expected performance
    difficultyMultiplier := 1.0 + difficulty
    performanceBonus := 1.0 / expectedPerformance
    
    totalMultiplier := difficultyMultiplier * performanceBonus
    
    return int(float64(ml.BasePoints) * totalMultiplier)
}
```

## 6. Composite Scoring System

### Multi-Factor Scoring Engine
```go
type CompositeScoring struct {
    Factors []ScoringFactor
    Weights map[string]float64
}

type ScoringFactor interface {
    Name() string
    Calculate(prediction Prediction, result EventResult, context ScoringContext) float64
}

type ScoringContext struct {
    User         User
    Event        Event
    Contest      Contest
    UserHistory  UserHistory
    EventOdds    EventOdds
}

func (cs *CompositeScoring) Calculate(prediction Prediction, result EventResult, context ScoringContext) int {
    if !cs.IsCorrect(prediction, result) {
        return 0
    }
    
    totalScore := 0.0
    
    for _, factor := range cs.Factors {
        factorScore := factor.Calculate(prediction, result, context)
        weight := cs.Weights[factor.Name()]
        totalScore += factorScore * weight
    }
    
    return int(totalScore)
}

// Example factors
type DifficultyFactor struct{}
func (df *DifficultyFactor) Name() string { return "difficulty" }
func (df *DifficultyFactor) Calculate(prediction Prediction, result EventResult, context ScoringContext) float64 {
    // Calculate difficulty based on odds, historical data, etc.
    return 1.0 // Placeholder
}

type ConfidenceFactor struct{}
func (cf *ConfidenceFactor) Name() string { return "confidence" }
func (cf *ConfidenceFactor) Calculate(prediction Prediction, result EventResult, context ScoringContext) float64 {
    return prediction.Confidence / 100.0
}

type TimingFactor struct{}
func (tf *TimingFactor) Name() string { return "timing" }
func (tf *TimingFactor) Calculate(prediction Prediction, result EventResult, context ScoringContext) float64 {
    // Bonus for early predictions
    timeBeforeEvent := context.Event.StartTime.Sub(prediction.SubmittedAt)
    hours := timeBeforeEvent.Hours()
    
    if hours > 24 {
        return 1.2 // 20% bonus for predictions made >24h early
    } else if hours > 1 {
        return 1.1 // 10% bonus for predictions made >1h early
    }
    return 1.0
}
```

## 7. Implementation Best Practices

### Scoring Pipeline
```go
type ScoringPipeline struct {
    scorers []Scorer
    cache   ScoringCache
}

func (sp *ScoringPipeline) ProcessEvent(eventID string) error {
    event, err := sp.getEvent(eventID)
    if err != nil {
        return err
    }
    
    result, err := sp.getEventResult(eventID)
    if err != nil {
        return err
    }
    
    predictions, err := sp.getPredictions(eventID)
    if err != nil {
        return err
    }
    
    for _, prediction := range predictions {
        score := sp.calculateScore(prediction, result, event)
        if err := sp.updatePredictionScore(prediction.ID, score); err != nil {
            log.Printf("Failed to update score for prediction %s: %v", prediction.ID, err)
        }
    }
    
    return sp.updateLeaderboards(event.ContestIDs)
}
```

### Caching Strategy
```go
type ScoringCache struct {
    redis *redis.Client
}

func (sc *ScoringCache) GetUserScore(userID, contestID string) (int, error) {
    key := fmt.Sprintf("user_score:%s:%s", userID, contestID)
    return sc.redis.Get(key).Int()
}

func (sc *ScoringCache) UpdateUserScore(userID, contestID string, score int) error {
    key := fmt.Sprintf("user_score:%s:%s", userID, contestID)
    return sc.redis.Set(key, score, time.Hour).Err()
}
```

### Audit Trail
```go
type ScoreAudit struct {
    PredictionID string    `json:"prediction_id"`
    OldScore     int       `json:"old_score"`
    NewScore     int       `json:"new_score"`
    Algorithm    string    `json:"algorithm"`
    Factors      []string  `json:"factors"`
    Timestamp    time.Time `json:"timestamp"`
    ProcessedBy  string    `json:"processed_by"`
}

func (sp *ScoringPipeline) auditScoreChange(predictionID string, oldScore, newScore int, algorithm string) {
    audit := ScoreAudit{
        PredictionID: predictionID,
        OldScore:     oldScore,
        NewScore:     newScore,
        Algorithm:    algorithm,
        Timestamp:    time.Now(),
        ProcessedBy:  "scoring-service",
    }
    
    sp.saveAuditRecord(audit)
}
```

This comprehensive scoring system allows for flexible, fair, and engaging prediction contests that can adapt to different sports, user preferences, and contest formats.