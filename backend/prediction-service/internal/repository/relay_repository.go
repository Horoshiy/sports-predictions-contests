package repository

import (
	"errors"

	"github.com/sports-prediction-contests/prediction-service/internal/models"
	"gorm.io/gorm"
)

// RelayRepositoryInterface defines the contract for relay repository
type RelayRepositoryInterface interface {
	// SetTeamAssignments replaces all assignments for a team (captain action)
	SetTeamAssignments(contestID, teamID, captainID uint, assignments []RelayAssignmentInput) error
	// GetTeamAssignments retrieves all assignments for a team
	GetTeamAssignments(contestID, teamID uint) ([]*models.RelayEventAssignment, error)
	// GetUserAssignments retrieves assignments for a specific user
	GetUserAssignments(contestID, userID uint) ([]*models.RelayEventAssignment, error)
	// GetUserAssignmentsForTeam retrieves assignments for a user within a team
	GetUserAssignmentsForTeam(contestID, teamID, userID uint) ([]*models.RelayEventAssignment, error)
	// ValidateUserCanPredict checks if user is assigned to predict this event
	ValidateUserCanPredict(contestID, teamID, userID, eventID uint) (bool, error)
	// GetAssignmentStats returns stats about assignments for a team
	GetAssignmentStats(contestID, teamID uint) (*RelayAssignmentStats, error)
}

// RelayAssignmentInput represents input for creating/updating an assignment
type RelayAssignmentInput struct {
	UserID  uint
	EventID uint
}

// RelayAssignmentStats contains statistics about assignments
type RelayAssignmentStats struct {
	TotalEvents    int64
	AssignedEvents int64
	MemberCounts   map[uint]int64 // userID -> assigned count
}

// RelayRepository implements RelayRepositoryInterface
type RelayRepository struct {
	db *gorm.DB
}

// NewRelayRepository creates a new relay repository instance
func NewRelayRepository(db *gorm.DB) RelayRepositoryInterface {
	return &RelayRepository{db: db}
}

// SetTeamAssignments replaces all assignments for a team
func (r *RelayRepository) SetTeamAssignments(contestID, teamID, captainID uint, assignments []RelayAssignmentInput) error {
	if contestID == 0 || teamID == 0 || captainID == 0 {
		return errors.New("contestID, teamID, and captainID are required")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete all existing assignments for this team in this contest
		if err := tx.Where("contest_id = ? AND team_id = ?", contestID, teamID).
			Delete(&models.RelayEventAssignment{}).Error; err != nil {
			return err
		}

		// Create new assignments
		for _, input := range assignments {
			assignment := models.RelayEventAssignment{
				ContestID:  contestID,
				TeamID:     teamID,
				UserID:     input.UserID,
				EventID:    input.EventID,
				AssignedBy: captainID,
			}
			if err := tx.Create(&assignment).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetTeamAssignments retrieves all assignments for a team
func (r *RelayRepository) GetTeamAssignments(contestID, teamID uint) ([]*models.RelayEventAssignment, error) {
	var assignments []*models.RelayEventAssignment
	err := r.db.Preload("Event").
		Where("contest_id = ? AND team_id = ?", contestID, teamID).
		Order("user_id, event_id").
		Find(&assignments).Error
	return assignments, err
}

// GetUserAssignments retrieves assignments for a specific user across all teams
func (r *RelayRepository) GetUserAssignments(contestID, userID uint) ([]*models.RelayEventAssignment, error) {
	var assignments []*models.RelayEventAssignment
	err := r.db.Preload("Event").
		Where("contest_id = ? AND user_id = ?", contestID, userID).
		Order("event_id").
		Find(&assignments).Error
	return assignments, err
}

// GetUserAssignmentsForTeam retrieves assignments for a user within a specific team
func (r *RelayRepository) GetUserAssignmentsForTeam(contestID, teamID, userID uint) ([]*models.RelayEventAssignment, error) {
	var assignments []*models.RelayEventAssignment
	err := r.db.Preload("Event").
		Where("contest_id = ? AND team_id = ? AND user_id = ?", contestID, teamID, userID).
		Order("event_id").
		Find(&assignments).Error
	return assignments, err
}

// ValidateUserCanPredict checks if user is assigned to predict this event
func (r *RelayRepository) ValidateUserCanPredict(contestID, teamID, userID, eventID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.RelayEventAssignment{}).
		Where("contest_id = ? AND team_id = ? AND user_id = ? AND event_id = ?",
			contestID, teamID, userID, eventID).
		Count(&count).Error
	return count > 0, err
}

// GetAssignmentStats returns stats about assignments for a team
func (r *RelayRepository) GetAssignmentStats(contestID, teamID uint) (*RelayAssignmentStats, error) {
	stats := &RelayAssignmentStats{
		MemberCounts: make(map[uint]int64),
	}

	// Get total events for this contest
	err := r.db.Raw(`
		SELECT COUNT(*) FROM contest_events WHERE contest_id = ?
	`, contestID).Scan(&stats.TotalEvents).Error
	if err != nil {
		return nil, err
	}

	// Get assigned events count
	err = r.db.Model(&models.RelayEventAssignment{}).
		Where("contest_id = ? AND team_id = ?", contestID, teamID).
		Count(&stats.AssignedEvents).Error
	if err != nil {
		return nil, err
	}

	// Get per-member counts
	type memberCount struct {
		UserID uint
		Count  int64
	}
	var counts []memberCount
	err = r.db.Model(&models.RelayEventAssignment{}).
		Select("user_id, COUNT(*) as count").
		Where("contest_id = ? AND team_id = ?", contestID, teamID).
		Group("user_id").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	for _, mc := range counts {
		stats.MemberCounts[mc.UserID] = mc.Count
	}

	return stats, nil
}
