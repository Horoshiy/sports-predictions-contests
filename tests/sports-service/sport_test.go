package sports_service_test

import (
	"testing"

	"github.com/sports-prediction-contests/sports-service/internal/models"
	"github.com/sports-prediction-contests/sports-service/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := db.AutoMigrate(&models.Sport{}, &models.League{}, &models.Team{}, &models.Match{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	return db
}

func TestSportValidation(t *testing.T) {
	sport := &models.Sport{Name: "Football", Slug: "football"}
	if err := sport.ValidateName(); err != nil {
		t.Errorf("Expected valid name, got error: %v", err)
	}

	sport.Name = ""
	if err := sport.ValidateName(); err == nil {
		t.Error("Expected error for empty name")
	}
}

func TestSportSlugValidation(t *testing.T) {
	sport := &models.Sport{Name: "Football", Slug: "football"}
	if err := sport.ValidateSlug(); err != nil {
		t.Errorf("Expected valid slug, got error: %v", err)
	}

	sport.Slug = "Invalid Slug!"
	if err := sport.ValidateSlug(); err == nil {
		t.Error("Expected error for invalid slug")
	}
}

func TestSportAutoSlugSanitization(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSportRepository(db)

	// Test special characters are sanitized
	sport := &models.Sport{Name: "Fútbol (Soccer)"}
	if err := repo.Create(sport); err != nil {
		t.Fatalf("Failed to create sport with special chars: %v", err)
	}

	if sport.Slug != "ftbol-soccer" {
		t.Errorf("Expected sanitized slug 'ftbol-soccer', got '%s'", sport.Slug)
	}

	// Test spaces are converted to hyphens
	sport2 := &models.Sport{Name: "Ice Hockey"}
	if err := repo.Create(sport2); err != nil {
		t.Fatalf("Failed to create sport: %v", err)
	}

	if sport2.Slug != "ice-hockey" {
		t.Errorf("Expected slug 'ice-hockey', got '%s'", sport2.Slug)
	}
}

func TestSportRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSportRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football", Description: "Association football"}
	if err := repo.Create(sport); err != nil {
		t.Fatalf("Failed to create sport: %v", err)
	}

	if sport.ID == 0 {
		t.Error("Expected sport ID to be set")
	}

	retrieved, err := repo.GetByID(sport.ID)
	if err != nil {
		t.Fatalf("Failed to get sport: %v", err)
	}

	if retrieved.Name != sport.Name {
		t.Errorf("Expected name %s, got %s", sport.Name, retrieved.Name)
	}
}

func TestSportList(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSportRepository(db)

	sports := []*models.Sport{
		{Name: "Football", Slug: "football", IsActive: true},
		{Name: "Basketball", Slug: "basketball", IsActive: true},
		{Name: "Tennis", Slug: "tennis", IsActive: false},
	}

	for _, s := range sports {
		if err := repo.Create(s); err != nil {
			t.Fatalf("Failed to create sport: %v", err)
		}
	}

	list, total, err := repo.List(10, 0, false)
	if err != nil {
		t.Fatalf("Failed to list sports: %v", err)
	}

	if total != 3 {
		t.Errorf("Expected 3 sports, got %d", total)
	}

	activeList, activeTotal, err := repo.List(10, 0, true)
	if err != nil {
		t.Fatalf("Failed to list active sports: %v", err)
	}

	if activeTotal != 2 {
		t.Errorf("Expected 2 active sports, got %d", activeTotal)
	}

	if len(activeList) != 2 {
		t.Errorf("Expected 2 active sports in list, got %d", len(activeList))
	}
}

func TestSportUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSportRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football"}
	repo.Create(sport)

	sport.Description = "Updated description"
	if err := repo.Update(sport); err != nil {
		t.Fatalf("Failed to update sport: %v", err)
	}

	retrieved, _ := repo.GetByID(sport.ID)
	if retrieved.Description != "Updated description" {
		t.Errorf("Expected updated description, got %s", retrieved.Description)
	}
}

func TestSportDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSportRepository(db)

	sport := &models.Sport{Name: "Football", Slug: "football"}
	repo.Create(sport)

	if err := repo.Delete(sport.ID); err != nil {
		t.Fatalf("Failed to delete sport: %v", err)
	}

	_, err := repo.GetByID(sport.ID)
	if err == nil {
		t.Error("Expected error when getting deleted sport")
	}
}
