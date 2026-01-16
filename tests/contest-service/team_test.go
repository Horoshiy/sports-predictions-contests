package contest_service_test

import (
	"testing"

	"github.com/sports-prediction-contests/contest-service/internal/models"
)

func TestTeamValidateName(t *testing.T) {
	tests := []struct {
		name    string
		team    models.Team
		wantErr bool
	}{
		{"valid name", models.Team{Name: "Test Team"}, false},
		{"empty name", models.Team{Name: ""}, true},
		{"whitespace only", models.Team{Name: "   "}, true},
		{"too long", models.Team{Name: string(make([]byte, 101))}, true},
		{"max length", models.Team{Name: string(make([]byte, 100))}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.team.ValidateName()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTeamValidateDescription(t *testing.T) {
	tests := []struct {
		name    string
		team    models.Team
		wantErr bool
	}{
		{"empty description", models.Team{Description: ""}, false},
		{"valid description", models.Team{Description: "A great team"}, false},
		{"too long", models.Team{Description: string(make([]byte, 501))}, true},
		{"max length", models.Team{Description: string(make([]byte, 500))}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.team.ValidateDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTeamValidateMaxMembers(t *testing.T) {
	tests := []struct {
		name    string
		team    models.Team
		wantErr bool
	}{
		{"valid 10", models.Team{MaxMembers: 10}, false},
		{"min valid", models.Team{MaxMembers: 2}, false},
		{"max valid", models.Team{MaxMembers: 50}, false},
		{"too small", models.Team{MaxMembers: 1}, true},
		{"too large", models.Team{MaxMembers: 51}, true},
		{"zero fails validation", models.Team{MaxMembers: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.team.ValidateMaxMembers()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMaxMembers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateInviteCode(t *testing.T) {
	code1, err := models.GenerateInviteCode()
	if err != nil {
		t.Fatalf("GenerateInviteCode() error = %v", err)
	}
	if len(code1) != 8 {
		t.Errorf("GenerateInviteCode() length = %d, want 8", len(code1))
	}

	// Test uniqueness
	code2, _ := models.GenerateInviteCode()
	if code1 == code2 {
		t.Error("GenerateInviteCode() generated duplicate codes")
	}
}

func TestTeamCanJoin(t *testing.T) {
	tests := []struct {
		name string
		team models.Team
		want bool
	}{
		{"active with space", models.Team{IsActive: true, CurrentMembers: 5, MaxMembers: 10}, true},
		{"active full", models.Team{IsActive: true, CurrentMembers: 10, MaxMembers: 10}, false},
		{"inactive", models.Team{IsActive: false, CurrentMembers: 5, MaxMembers: 10}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.team.CanJoin(); got != tt.want {
				t.Errorf("CanJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTeamMemberValidateRole(t *testing.T) {
	tests := []struct {
		name    string
		member  models.TeamMember
		wantErr bool
	}{
		{"captain", models.TeamMember{Role: "captain"}, false},
		{"member", models.TeamMember{Role: "member"}, false},
		{"invalid", models.TeamMember{Role: "admin"}, true},
		{"empty", models.TeamMember{Role: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.member.ValidateRole()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTeamMemberValidateStatus(t *testing.T) {
	tests := []struct {
		name    string
		member  models.TeamMember
		wantErr bool
	}{
		{"active", models.TeamMember{Status: "active"}, false},
		{"inactive", models.TeamMember{Status: "inactive"}, false},
		{"invalid", models.TeamMember{Status: "banned"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.member.ValidateStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
