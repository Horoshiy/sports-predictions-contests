package repository

import (
	"errors"

	"github.com/sports-prediction-contests/contest-service/internal/models"
	"gorm.io/gorm"
)

type TeamRepositoryInterface interface {
	Create(team *models.Team) error
	CreateWithMember(team *models.Team, member *models.TeamMember) error
	GetByID(id uint) (*models.Team, error)
	GetByInviteCode(code string) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uint) error
	List(limit, offset int, userID uint, myTeamsOnly bool) ([]*models.Team, int64, error)
}

type TeamMemberRepositoryInterface interface {
	Create(member *models.TeamMember) error
	GetByID(id uint) (*models.TeamMember, error)
	GetByTeamAndUser(teamID, userID uint) (*models.TeamMember, error)
	Update(member *models.TeamMember) error
	Delete(id uint) error
	DeleteByTeamAndUser(teamID, userID uint) error
	ListByTeam(teamID uint, limit, offset int) ([]*models.TeamMember, int64, error)
	CountByTeam(teamID uint) (int64, error)
}

type TeamContestEntryRepositoryInterface interface {
	Create(entry *models.TeamContestEntry) error
	GetByTeamAndContest(teamID, contestID uint) (*models.TeamContestEntry, error)
	Delete(teamID, contestID uint) error
	ListByContest(contestID uint, limit int) ([]*models.TeamContestEntry, error)
	UpdatePoints(teamID, contestID uint, points float64) error
}

type TeamRepository struct{ db *gorm.DB }
type TeamMemberRepository struct{ db *gorm.DB }
type TeamContestEntryRepository struct{ db *gorm.DB }

func NewTeamRepository(db *gorm.DB) TeamRepositoryInterface             { return &TeamRepository{db: db} }
func NewTeamMemberRepository(db *gorm.DB) TeamMemberRepositoryInterface { return &TeamMemberRepository{db: db} }
func NewTeamContestEntryRepository(db *gorm.DB) TeamContestEntryRepositoryInterface {
	return &TeamContestEntryRepository{db: db}
}

func (r *TeamRepository) Create(team *models.Team) error {
	if team == nil {
		return errors.New("team cannot be nil")
	}
	return r.db.Create(team).Error
}

func (r *TeamRepository) CreateWithMember(team *models.Team, member *models.TeamMember) error {
	if team == nil {
		return errors.New("team cannot be nil")
	}
	if member == nil {
		return errors.New("member cannot be nil")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}
		member.TeamID = team.ID
		if err := tx.Create(member).Error; err != nil {
			return err
		}
		team.CurrentMembers = 1
		return tx.Save(team).Error
	})
}

func (r *TeamRepository) GetByID(id uint) (*models.Team, error) {
	if id == 0 {
		return nil, errors.New("invalid team ID")
	}
	var team models.Team
	if err := r.db.First(&team, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) GetByInviteCode(code string) (*models.Team, error) {
	if code == "" {
		return nil, errors.New("invite code cannot be empty")
	}
	var team models.Team
	if err := r.db.Where("invite_code = ? AND is_active = ?", code, true).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid invite code")
		}
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) Update(team *models.Team) error {
	if team == nil || team.ID == 0 {
		return errors.New("invalid team")
	}
	return r.db.Save(team).Error
}

func (r *TeamRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid team ID")
	}
	tx := r.db.Begin()
	defer tx.Rollback() // GORM ignores rollback on committed transactions
	if err := tx.Where("team_id = ?", id).Delete(&models.TeamMember{}).Error; err != nil {
		return err
	}
	if err := tx.Where("team_id = ?", id).Delete(&models.TeamContestEntry{}).Error; err != nil {
		return err
	}
	if err := tx.Delete(&models.Team{}, id).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

func (r *TeamRepository) List(limit, offset int, userID uint, myTeamsOnly bool) ([]*models.Team, int64, error) {
	var teams []*models.Team
	var total int64
	query := r.db.Model(&models.Team{}).Where("is_active = ? AND deleted_at IS NULL", true)
	if myTeamsOnly && userID > 0 {
		query = query.Where("id IN (SELECT team_id FROM user_team_members WHERE user_id = ? AND status = ? AND deleted_at IS NULL)", userID, "active")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&teams).Error; err != nil {
		return nil, 0, err
	}
	return teams, total, nil
}

func (r *TeamMemberRepository) Create(member *models.TeamMember) error {
	if member == nil {
		return errors.New("member cannot be nil")
	}
	return r.db.Create(member).Error
}

func (r *TeamMemberRepository) GetByID(id uint) (*models.TeamMember, error) {
	if id == 0 {
		return nil, errors.New("invalid member ID")
	}
	var member models.TeamMember
	if err := r.db.Preload("Team").First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, err
	}
	return &member, nil
}

func (r *TeamMemberRepository) GetByTeamAndUser(teamID, userID uint) (*models.TeamMember, error) {
	if teamID == 0 || userID == 0 {
		return nil, errors.New("invalid team or user ID")
	}
	var member models.TeamMember
	if err := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, err
	}
	return &member, nil
}

func (r *TeamMemberRepository) Update(member *models.TeamMember) error {
	if member == nil || member.ID == 0 {
		return errors.New("invalid member")
	}
	return r.db.Save(member).Error
}

func (r *TeamMemberRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid member ID")
	}
	return r.db.Delete(&models.TeamMember{}, id).Error
}

func (r *TeamMemberRepository) DeleteByTeamAndUser(teamID, userID uint) error {
	if teamID == 0 || userID == 0 {
		return errors.New("invalid team or user ID")
	}
	result := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&models.TeamMember{})
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return result.Error
}

func (r *TeamMemberRepository) ListByTeam(teamID uint, limit, offset int) ([]*models.TeamMember, int64, error) {
	if teamID == 0 {
		return nil, 0, errors.New("invalid team ID")
	}
	var members []*models.TeamMember
	var total int64
	query := r.db.Where("team_id = ?", teamID)
	if err := query.Model(&models.TeamMember{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("role DESC, joined_at ASC").Limit(limit).Offset(offset).Find(&members).Error; err != nil {
		return nil, 0, err
	}
	return members, total, nil
}

func (r *TeamMemberRepository) CountByTeam(teamID uint) (int64, error) {
	if teamID == 0 {
		return 0, errors.New("invalid team ID")
	}
	var count int64
	return count, r.db.Model(&models.TeamMember{}).Where("team_id = ? AND status = ?", teamID, "active").Count(&count).Error
}

func (r *TeamContestEntryRepository) Create(entry *models.TeamContestEntry) error {
	if entry == nil {
		return errors.New("entry cannot be nil")
	}
	return r.db.Create(entry).Error
}

func (r *TeamContestEntryRepository) GetByTeamAndContest(teamID, contestID uint) (*models.TeamContestEntry, error) {
	if teamID == 0 || contestID == 0 {
		return nil, errors.New("invalid team or contest ID")
	}
	var entry models.TeamContestEntry
	if err := r.db.Where("team_id = ? AND contest_id = ?", teamID, contestID).First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entry not found")
		}
		return nil, err
	}
	return &entry, nil
}

func (r *TeamContestEntryRepository) Delete(teamID, contestID uint) error {
	if teamID == 0 || contestID == 0 {
		return errors.New("invalid team or contest ID")
	}
	result := r.db.Where("team_id = ? AND contest_id = ?", teamID, contestID).Delete(&models.TeamContestEntry{})
	if result.RowsAffected == 0 {
		return errors.New("entry not found")
	}
	return result.Error
}

func (r *TeamContestEntryRepository) ListByContest(contestID uint, limit int) ([]*models.TeamContestEntry, error) {
	if contestID == 0 {
		return nil, errors.New("invalid contest ID")
	}
	if limit <= 0 {
		limit = 10
	}
	var entries []*models.TeamContestEntry
	return entries, r.db.Preload("Team").Where("contest_id = ?", contestID).Order("total_points DESC").Limit(limit).Find(&entries).Error
}

func (r *TeamContestEntryRepository) UpdatePoints(teamID, contestID uint, points float64) error {
	if teamID == 0 || contestID == 0 {
		return errors.New("invalid team or contest ID")
	}
	return r.db.Model(&models.TeamContestEntry{}).Where("team_id = ? AND contest_id = ?", teamID, contestID).Update("total_points", points).Error
}
