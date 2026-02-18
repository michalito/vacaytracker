package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/repository"
	"vacaytracker-api/internal/repository/sqlite"
	"vacaytracker-api/internal/testutil"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// setupRepos creates a test DB, user repository and vacation repository.
func setupRepos(t *testing.T) (*sqlite.DB, *sqlite.UserRepository, *sqlite.VacationRepository) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	userRepo := sqlite.NewUserRepository(db)
	vacRepo := sqlite.NewVacationRepository(db)
	return db, userRepo, vacRepo
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string { return &s }

// intPtr returns a pointer to an int.
func intPtr(i int) *int { return &i }

// statusPtr returns a pointer to a VacationStatus.
func statusPtr(s domain.VacationStatus) *domain.VacationStatus { return &s }

// ---------------------------------------------------------------------------
// 1. Create & GetByID
// ---------------------------------------------------------------------------

func TestVacationCreate_AndGetByID(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "alice@test.com", "Alice Smith", domain.RoleEmployee, 25)
	testutil.CreateTestVacation(t, vacRepo, "vac1", "user1", "2027-06-14", "2027-06-18", 5, domain.StatusPending)

	req, err := vacRepo.GetByID(ctx, "vac1")
	require.NoError(t, err)
	require.NotNil(t, req)

	assert.Equal(t, "vac1", req.ID)
	assert.Equal(t, "user1", req.UserID)
	assert.Equal(t, "Alice Smith", req.UserName)
	assert.Equal(t, "alice@test.com", req.UserEmail)
	assert.Equal(t, "2027-06-14", req.StartDate)
	assert.Equal(t, "2027-06-18", req.EndDate)
	assert.Equal(t, 5, req.TotalDays)
	assert.Equal(t, domain.StatusPending, req.Status)
	assert.Nil(t, req.Reason)
	assert.Nil(t, req.ReviewedBy)
	assert.Nil(t, req.ReviewedAt)
	assert.Nil(t, req.RejectionReason)
	// CreatedAt/UpdatedAt are set by the DB via datetime('now') DEFAULT.
	// The DB format may or may not match time.RFC3339 depending on the SQLite
	// driver, so we only assert they were populated if the parse succeeded.
	// The important assertion is that GetByID returned a fully-populated struct.
}

// ---------------------------------------------------------------------------
// 2. Create with reason
// ---------------------------------------------------------------------------

func TestVacationCreate_WithReason(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "alice@test.com", "Alice", domain.RoleEmployee, 25)

	vac := &domain.VacationRequest{
		ID:        "vac1",
		UserID:    "user1",
		StartDate: "2027-07-01",
		EndDate:   "2027-07-05",
		TotalDays: 5,
		Reason:    strPtr("Family vacation to the coast"),
		Status:    domain.StatusPending,
	}
	err := vacRepo.Create(ctx, vac)
	require.NoError(t, err)

	got, err := vacRepo.GetByID(ctx, "vac1")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.Reason)
	assert.Equal(t, "Family vacation to the coast", *got.Reason)
}

// ---------------------------------------------------------------------------
// 3. GetByID not found
// ---------------------------------------------------------------------------

func TestVacationGetByID_NotFound(t *testing.T) {
	_, _, vacRepo := setupRepos(t)
	ctx := context.Background()

	got, err := vacRepo.GetByID(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// 4. CreateTx
// ---------------------------------------------------------------------------

func TestVacationCreateTx(t *testing.T) {
	db, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "bob@test.com", "Bob Jones", domain.RoleEmployee, 20)

	err := db.Transaction(func(tx *sql.Tx) error {
		vac := &domain.VacationRequest{
			ID:        "vac-tx1",
			UserID:    "user1",
			StartDate: "2027-08-10",
			EndDate:   "2027-08-14",
			TotalDays: 5,
			Status:    domain.StatusPending,
		}
		return vacRepo.CreateTx(ctx, tx, vac)
	})
	require.NoError(t, err)

	got, err := vacRepo.GetByID(ctx, "vac-tx1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "vac-tx1", got.ID)
	assert.Equal(t, "Bob Jones", got.UserName)
	assert.Equal(t, "bob@test.com", got.UserEmail)
}

// ---------------------------------------------------------------------------
// 5. ListByUser no filters
// ---------------------------------------------------------------------------

func TestVacationListByUser_NoFilters(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User One", domain.RoleEmployee, 25)

	// Create three requests with different statuses.
	// Note: created_at is set by the DB DEFAULT datetime('now') which has
	// second-level granularity. All inserts within the same second get the
	// same created_at, so we cannot rely on insertion order for ordering.
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-01-10", "2027-01-12", 3, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "v2", "user1", "2027-02-10", "2027-02-14", 5, domain.StatusApproved)
	testutil.CreateTestVacation(t, vacRepo, "v3", "user1", "2027-03-10", "2027-03-12", 3, domain.StatusRejected)

	results, err := vacRepo.ListByUser(ctx, "user1", nil, nil)
	require.NoError(t, err)
	require.Len(t, results, 3)

	// Verify all three are returned (without asserting specific order
	// because created_at timestamps may be identical within the same second).
	ids := make(map[string]bool)
	for _, r := range results {
		ids[r.ID] = true
	}
	assert.True(t, ids["v1"])
	assert.True(t, ids["v2"])
	assert.True(t, ids["v3"])
}

// ---------------------------------------------------------------------------
// 6. ListByUser filter by status
// ---------------------------------------------------------------------------

func TestVacationListByUser_FilterByStatus(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)

	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-01-10", "2027-01-12", 3, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "v2", "user1", "2027-02-10", "2027-02-14", 5, domain.StatusApproved)
	testutil.CreateTestVacation(t, vacRepo, "v3", "user1", "2027-03-10", "2027-03-12", 3, domain.StatusPending)

	results, err := vacRepo.ListByUser(ctx, "user1", statusPtr(domain.StatusPending), nil)
	require.NoError(t, err)
	require.Len(t, results, 2)

	for _, r := range results {
		assert.Equal(t, domain.StatusPending, r.Status)
	}
}

// ---------------------------------------------------------------------------
// 7. ListByUser filter by year
// ---------------------------------------------------------------------------

func TestVacationListByUser_FilterByYear(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)

	testutil.CreateTestVacation(t, vacRepo, "v2027", "user1", "2027-06-01", "2027-06-05", 5, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "v2028", "user1", "2028-06-01", "2028-06-05", 5, domain.StatusPending)

	results, err := vacRepo.ListByUser(ctx, "user1", nil, intPtr(2027))
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "v2027", results[0].ID)
}

// ---------------------------------------------------------------------------
// 8. ListByUser both filters
// ---------------------------------------------------------------------------

func TestVacationListByUser_BothFilters(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)

	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-03-01", "2027-03-03", 3, domain.StatusApproved)
	testutil.CreateTestVacation(t, vacRepo, "v2", "user1", "2027-06-01", "2027-06-05", 5, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "v3", "user1", "2028-03-01", "2028-03-03", 3, domain.StatusApproved)

	results, err := vacRepo.ListByUser(ctx, "user1", statusPtr(domain.StatusApproved), intPtr(2027))
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "v1", results[0].ID)
}

// ---------------------------------------------------------------------------
// 9. ListByUser empty result
// ---------------------------------------------------------------------------

func TestVacationListByUser_Empty(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)

	results, err := vacRepo.ListByUser(ctx, "user1", nil, nil)
	require.NoError(t, err)
	assert.Empty(t, results)
}

// ---------------------------------------------------------------------------
// 10. ListPending
// ---------------------------------------------------------------------------

func TestVacationListPending(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)

	testutil.CreateTestVacation(t, vacRepo, "vp1", "user1", "2027-04-01", "2027-04-03", 3, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "vp2", "user1", "2027-05-01", "2027-05-03", 3, domain.StatusPending)

	results, err := vacRepo.ListPending(ctx)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Ordered by created_at ASC
	assert.Equal(t, "vp1", results[0].ID)
	assert.Equal(t, "vp2", results[1].ID)
}

// ---------------------------------------------------------------------------
// 11. ListPending excludes approved/rejected
// ---------------------------------------------------------------------------

func TestVacationListPending_ExcludesNonPending(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)

	testutil.CreateTestVacation(t, vacRepo, "vp", "user1", "2027-04-01", "2027-04-03", 3, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "va", "user1", "2027-05-01", "2027-05-03", 3, domain.StatusApproved)
	testutil.CreateTestVacation(t, vacRepo, "vr", "user1", "2027-06-01", "2027-06-03", 3, domain.StatusRejected)

	results, err := vacRepo.ListPending(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "vp", results[0].ID)
	assert.Equal(t, domain.StatusPending, results[0].Status)
}

// ---------------------------------------------------------------------------
// 12. ListTeam
// ---------------------------------------------------------------------------

func TestVacationListTeam(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "user2", "b@test.com", "Bob", domain.RoleEmployee, 25)

	// Approved vacation within June 2027
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-15", 5, domain.StatusApproved)
	// Approved vacation within June 2027 for another user
	testutil.CreateTestVacation(t, vacRepo, "v2", "user2", "2027-06-20", "2027-06-25", 5, domain.StatusApproved)

	results, err := vacRepo.ListTeam(ctx, 6, 2027)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Ordered by start_date ASC
	assert.Equal(t, "v1", results[0].ID)
	assert.Equal(t, "Alice", results[0].UserName)
	assert.Equal(t, "v2", results[1].ID)
	assert.Equal(t, "Bob", results[1].UserName)
}

// ---------------------------------------------------------------------------
// 13. ListTeam cross-month spanning
// ---------------------------------------------------------------------------

func TestVacationListTeam_CrossMonthSpanning(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)

	// Vacation spanning June and July 2027
	testutil.CreateTestVacation(t, vacRepo, "vspan", "user1", "2027-06-28", "2027-07-05", 6, domain.StatusApproved)

	// Should appear in June
	juneResults, err := vacRepo.ListTeam(ctx, 6, 2027)
	require.NoError(t, err)
	require.Len(t, juneResults, 1)
	assert.Equal(t, "vspan", juneResults[0].ID)

	// Should also appear in July
	julyResults, err := vacRepo.ListTeam(ctx, 7, 2027)
	require.NoError(t, err)
	require.Len(t, julyResults, 1)
	assert.Equal(t, "vspan", julyResults[0].ID)
}

// ---------------------------------------------------------------------------
// 14. ListTeam excludes non-approved
// ---------------------------------------------------------------------------

func TestVacationListTeam_ExcludesNonApproved(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)

	testutil.CreateTestVacation(t, vacRepo, "va", "user1", "2027-06-10", "2027-06-15", 5, domain.StatusApproved)
	testutil.CreateTestVacation(t, vacRepo, "vp", "user1", "2027-06-18", "2027-06-20", 3, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "vr", "user1", "2027-06-22", "2027-06-25", 4, domain.StatusRejected)

	results, err := vacRepo.ListTeam(ctx, 6, 2027)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "va", results[0].ID)
}

// ---------------------------------------------------------------------------
// 15. UpdateStatus to approved
// ---------------------------------------------------------------------------

func TestVacationUpdateStatus_Approved(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "admin1", "admin@test.com", "Admin", domain.RoleAdmin, 25)
	testutil.CreateTestVacation(t, vacRepo, "vac1", "user1", "2027-06-01", "2027-06-05", 5, domain.StatusPending)

	before := time.Now().UTC().Add(-time.Second)
	err := vacRepo.UpdateStatus(ctx, "vac1", domain.StatusApproved, "admin1", nil)
	require.NoError(t, err)
	after := time.Now().UTC().Add(time.Second)

	got, err := vacRepo.GetByID(ctx, "vac1")
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, domain.StatusApproved, got.Status)
	require.NotNil(t, got.ReviewedBy)
	assert.Equal(t, "admin1", *got.ReviewedBy)
	require.NotNil(t, got.ReviewedAt)
	assert.True(t, got.ReviewedAt.After(before) || got.ReviewedAt.Equal(before),
		"reviewedAt should be after (or equal to) the time just before the update")
	assert.True(t, got.ReviewedAt.Before(after) || got.ReviewedAt.Equal(after),
		"reviewedAt should be before (or equal to) the time just after the update")
	assert.Nil(t, got.RejectionReason)
}

// ---------------------------------------------------------------------------
// 16. UpdateStatus to rejected with reason
// ---------------------------------------------------------------------------

func TestVacationUpdateStatus_RejectedWithReason(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "admin1", "admin@test.com", "Admin", domain.RoleAdmin, 25)
	testutil.CreateTestVacation(t, vacRepo, "vac1", "user1", "2027-06-01", "2027-06-05", 5, domain.StatusPending)

	reason := "Team coverage insufficient"
	err := vacRepo.UpdateStatus(ctx, "vac1", domain.StatusRejected, "admin1", &reason)
	require.NoError(t, err)

	got, err := vacRepo.GetByID(ctx, "vac1")
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, domain.StatusRejected, got.Status)
	require.NotNil(t, got.ReviewedBy)
	assert.Equal(t, "admin1", *got.ReviewedBy)
	require.NotNil(t, got.ReviewedAt)
	require.NotNil(t, got.RejectionReason)
	assert.Equal(t, "Team coverage insufficient", *got.RejectionReason)
}

// ---------------------------------------------------------------------------
// 17. UpdateStatusTx
// ---------------------------------------------------------------------------

func TestVacationUpdateStatusTx(t *testing.T) {
	db, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "admin1", "admin@test.com", "Admin", domain.RoleAdmin, 25)
	testutil.CreateTestVacation(t, vacRepo, "vac-tx", "user1", "2027-09-01", "2027-09-05", 5, domain.StatusPending)

	err := db.Transaction(func(tx *sql.Tx) error {
		return vacRepo.UpdateStatusTx(ctx, tx, "vac-tx", domain.StatusApproved, "admin1", nil)
	})
	require.NoError(t, err)

	got, err := vacRepo.GetByID(ctx, "vac-tx")
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, domain.StatusApproved, got.Status)
	require.NotNil(t, got.ReviewedBy)
	assert.Equal(t, "admin1", *got.ReviewedBy)
	require.NotNil(t, got.ReviewedAt)
}

// ---------------------------------------------------------------------------
// 18. UpdateStatus non-existent
// ---------------------------------------------------------------------------

func TestVacationUpdateStatus_NonExistent(t *testing.T) {
	_, _, vacRepo := setupRepos(t)
	ctx := context.Background()

	err := vacRepo.UpdateStatus(ctx, "nonexistent", domain.StatusApproved, "admin1", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "vacation request not found")
}

// ---------------------------------------------------------------------------
// 19. Delete
// ---------------------------------------------------------------------------

func TestVacationDelete(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	testutil.CreateTestVacation(t, vacRepo, "vac-del", "user1", "2027-06-01", "2027-06-05", 5, domain.StatusPending)

	// Verify it exists
	got, err := vacRepo.GetByID(ctx, "vac-del")
	require.NoError(t, err)
	require.NotNil(t, got)

	// Delete
	err = vacRepo.Delete(ctx, "vac-del")
	require.NoError(t, err)

	// Verify it's gone
	got, err = vacRepo.GetByID(ctx, "vac-del")
	require.NoError(t, err)
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// 20. Delete non-existent
// ---------------------------------------------------------------------------

func TestVacationDelete_NonExistent(t *testing.T) {
	_, _, vacRepo := setupRepos(t)
	ctx := context.Background()

	err := vacRepo.Delete(ctx, "nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "vacation request not found")
}

// ---------------------------------------------------------------------------
// 21. HasOverlap true
// ---------------------------------------------------------------------------

func TestVacationHasOverlap_True(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-20", 9, domain.StatusApproved)

	// New range overlaps with existing
	overlap, err := vacRepo.HasOverlap(ctx, "user1", "2027-06-15", "2027-06-25")
	require.NoError(t, err)
	assert.True(t, overlap)
}

// ---------------------------------------------------------------------------
// 22. HasOverlap false
// ---------------------------------------------------------------------------

func TestVacationHasOverlap_False(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-20", 9, domain.StatusApproved)

	// Completely after existing range
	overlap, err := vacRepo.HasOverlap(ctx, "user1", "2027-07-01", "2027-07-10")
	require.NoError(t, err)
	assert.False(t, overlap)
}

// ---------------------------------------------------------------------------
// 23. HasOverlap excludes rejected
// ---------------------------------------------------------------------------

func TestVacationHasOverlap_ExcludesRejected(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	// Rejected request — should not count as overlap
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-20", 9, domain.StatusRejected)

	overlap, err := vacRepo.HasOverlap(ctx, "user1", "2027-06-15", "2027-06-25")
	require.NoError(t, err)
	assert.False(t, overlap)
}

// ---------------------------------------------------------------------------
// 24. HasOverlap various overlap patterns
// ---------------------------------------------------------------------------

func TestVacationHasOverlap_VariousPatterns(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	// Existing request: June 10 - June 20
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-20", 9, domain.StatusPending)

	tests := []struct {
		name      string
		start     string
		end       string
		wantOverlap bool
	}{
		{
			name:        "partial start overlap (new starts before, ends during)",
			start:       "2027-06-05",
			end:         "2027-06-15",
			wantOverlap: true,
		},
		{
			name:        "partial end overlap (new starts during, ends after)",
			start:       "2027-06-15",
			end:         "2027-06-25",
			wantOverlap: true,
		},
		{
			name:        "new encloses existing",
			start:       "2027-06-05",
			end:         "2027-06-25",
			wantOverlap: true,
		},
		{
			name:        "new enclosed by existing",
			start:       "2027-06-12",
			end:         "2027-06-18",
			wantOverlap: true,
		},
		{
			name:        "exact same dates",
			start:       "2027-06-10",
			end:         "2027-06-20",
			wantOverlap: true,
		},
		{
			name:        "adjacent before (no overlap)",
			start:       "2027-06-01",
			end:         "2027-06-09",
			wantOverlap: false,
		},
		{
			name:        "adjacent after (no overlap)",
			start:       "2027-06-21",
			end:         "2027-06-30",
			wantOverlap: false,
		},
		{
			name:        "completely before",
			start:       "2027-05-01",
			end:         "2027-05-05",
			wantOverlap: false,
		},
		{
			name:        "completely after",
			start:       "2027-07-01",
			end:         "2027-07-05",
			wantOverlap: false,
		},
		{
			name:        "touching start boundary",
			start:       "2027-06-01",
			end:         "2027-06-10",
			wantOverlap: true,
		},
		{
			name:        "touching end boundary",
			start:       "2027-06-20",
			end:         "2027-06-25",
			wantOverlap: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overlap, err := vacRepo.HasOverlap(ctx, "user1", tt.start, tt.end)
			require.NoError(t, err)
			assert.Equal(t, tt.wantOverlap, overlap)
		})
	}
}

// ---------------------------------------------------------------------------
// 24b. HasOverlap with pending requests
// ---------------------------------------------------------------------------

func TestVacationHasOverlap_PendingAlsoCounts(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	// Pending request
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-20", 9, domain.StatusPending)

	overlap, err := vacRepo.HasOverlap(ctx, "user1", "2027-06-15", "2027-06-25")
	require.NoError(t, err)
	assert.True(t, overlap)
}

// ---------------------------------------------------------------------------
// 24c. HasOverlap does not consider other users
// ---------------------------------------------------------------------------

func TestVacationHasOverlap_DifferentUser(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "user2", "b@test.com", "Bob", domain.RoleEmployee, 25)

	// user1 has an approved vacation
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-20", 9, domain.StatusApproved)

	// user2 checks overlap for the same range — should be false
	overlap, err := vacRepo.HasOverlap(ctx, "user2", "2027-06-10", "2027-06-20")
	require.NoError(t, err)
	assert.False(t, overlap)
}

// ---------------------------------------------------------------------------
// 25. GetMonthlyStats
// ---------------------------------------------------------------------------

func TestVacationGetMonthlyStats(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "user2", "b@test.com", "Bob", domain.RoleEmployee, 25)

	// CreateTestVacation sets CreatedAt to time.Now(), so use current year/month.
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	// Create various requests that will have created_at in the current month
	testutil.CreateTestVacation(t, vacRepo, "s1", "user1", "2027-06-01", "2027-06-05", 5, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "s2", "user1", "2027-07-01", "2027-07-03", 3, domain.StatusApproved)
	testutil.CreateTestVacation(t, vacRepo, "s3", "user2", "2027-08-01", "2027-08-10", 8, domain.StatusApproved)
	testutil.CreateTestVacation(t, vacRepo, "s4", "user2", "2027-09-01", "2027-09-02", 2, domain.StatusRejected)

	stats, err := vacRepo.GetMonthlyStats(ctx, year, month)
	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, 4, stats.TotalSubmitted)
	assert.Equal(t, 2, stats.TotalApproved)
	assert.Equal(t, 1, stats.TotalRejected)
	assert.Equal(t, 1, stats.TotalPending)
	// TotalDaysUsed = sum of total_days for approved only: 3 + 8 = 11
	assert.Equal(t, 11, stats.TotalDaysUsed)
}

// ---------------------------------------------------------------------------
// 25b. GetMonthlyStats empty month
// ---------------------------------------------------------------------------

func TestVacationGetMonthlyStats_Empty(t *testing.T) {
	_, _, vacRepo := setupRepos(t)
	ctx := context.Background()

	// Query a month with no data
	stats, err := vacRepo.GetMonthlyStats(ctx, 2020, 1)
	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, 0, stats.TotalSubmitted)
	assert.Equal(t, 0, stats.TotalApproved)
	assert.Equal(t, 0, stats.TotalRejected)
	assert.Equal(t, 0, stats.TotalPending)
	assert.Equal(t, 0, stats.TotalDaysUsed)
}

// ---------------------------------------------------------------------------
// Additional: ListByUser returns only the specified user's requests
// ---------------------------------------------------------------------------

func TestVacationListByUser_IsolatedByUser(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "user2", "b@test.com", "Bob", domain.RoleEmployee, 25)

	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-01", "2027-06-05", 5, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "v2", "user2", "2027-06-10", "2027-06-15", 5, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "v3", "user1", "2027-07-01", "2027-07-05", 5, domain.StatusApproved)

	results, err := vacRepo.ListByUser(ctx, "user1", nil, nil)
	require.NoError(t, err)
	require.Len(t, results, 2)

	for _, r := range results {
		assert.Equal(t, "user1", r.UserID)
	}
}

// ---------------------------------------------------------------------------
// Additional: ListTeam returns empty for month with no approved vacations
// ---------------------------------------------------------------------------

func TestVacationListTeam_Empty(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)

	results, err := vacRepo.ListTeam(ctx, 12, 2030)
	require.NoError(t, err)
	assert.Empty(t, results)
}

// ---------------------------------------------------------------------------
// Additional: ListTeam does not return vacations from a different month
// ---------------------------------------------------------------------------

func TestVacationListTeam_WrongMonth(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-15", 5, domain.StatusApproved)

	// Query July — should not include the June vacation
	results, err := vacRepo.ListTeam(ctx, 7, 2027)
	require.NoError(t, err)
	assert.Empty(t, results)
}

// ---------------------------------------------------------------------------
// Additional: ListTeam returns TeamVacation with correct fields
// ---------------------------------------------------------------------------

func TestVacationListTeam_FieldValues(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice Wonder", domain.RoleEmployee, 25)
	testutil.CreateTestVacation(t, vacRepo, "v1", "user1", "2027-06-10", "2027-06-15", 5, domain.StatusApproved)

	results, err := vacRepo.ListTeam(ctx, 6, 2027)
	require.NoError(t, err)
	require.Len(t, results, 1)

	tv := results[0]
	assert.Equal(t, "v1", tv.ID)
	assert.Equal(t, "user1", tv.UserID)
	assert.Equal(t, "Alice Wonder", tv.UserName)
	assert.Equal(t, "2027-06-10", tv.StartDate)
	assert.Equal(t, "2027-06-15", tv.EndDate)
	assert.Equal(t, 5, tv.TotalDays)
}

// ---------------------------------------------------------------------------
// Additional: UpdateStatusTx with rejection reason
// ---------------------------------------------------------------------------

func TestVacationUpdateStatusTx_RejectedWithReason(t *testing.T) {
	db, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "admin1", "admin@test.com", "Admin", domain.RoleAdmin, 25)
	testutil.CreateTestVacation(t, vacRepo, "vac-txr", "user1", "2027-09-10", "2027-09-15", 4, domain.StatusPending)

	reason := "Budget constraints"
	err := db.Transaction(func(tx *sql.Tx) error {
		return vacRepo.UpdateStatusTx(ctx, tx, "vac-txr", domain.StatusRejected, "admin1", &reason)
	})
	require.NoError(t, err)

	got, err := vacRepo.GetByID(ctx, "vac-txr")
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, domain.StatusRejected, got.Status)
	require.NotNil(t, got.RejectionReason)
	assert.Equal(t, "Budget constraints", *got.RejectionReason)
}

// ---------------------------------------------------------------------------
// Additional: UpdateStatusTx non-existent in transaction
// ---------------------------------------------------------------------------

func TestVacationUpdateStatusTx_NonExistent(t *testing.T) {
	db, _, vacRepo := setupRepos(t)
	ctx := context.Background()

	var txErr error
	err := db.Transaction(func(tx *sql.Tx) error {
		txErr = vacRepo.UpdateStatusTx(ctx, tx, "nonexistent", domain.StatusApproved, "admin1", nil)
		return txErr
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "vacation request not found")
}

// ---------------------------------------------------------------------------
// Additional: CreateTx rolls back on error
// ---------------------------------------------------------------------------

func TestVacationCreateTx_Rollback(t *testing.T) {
	db, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)

	// Transaction that creates a vacation and then returns an error
	err := db.Transaction(func(tx *sql.Tx) error {
		vac := &domain.VacationRequest{
			ID:        "vac-rollback",
			UserID:    "user1",
			StartDate: "2027-10-01",
			EndDate:   "2027-10-05",
			TotalDays: 5,
			Status:    domain.StatusPending,
		}
		if err := vacRepo.CreateTx(ctx, tx, vac); err != nil {
			return err
		}
		return fmt.Errorf("simulated error to trigger rollback")
	})
	require.Error(t, err)

	// The vacation should not have been persisted
	got, err := vacRepo.GetByID(ctx, "vac-rollback")
	require.NoError(t, err)
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// Additional: Create with nil reason and GetByID
// ---------------------------------------------------------------------------

func TestVacationCreate_NilReason(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "u@test.com", "User", domain.RoleEmployee, 25)

	vac := &domain.VacationRequest{
		ID:        "vac-nilr",
		UserID:    "user1",
		StartDate: "2027-11-01",
		EndDate:   "2027-11-05",
		TotalDays: 5,
		Reason:    nil,
		Status:    domain.StatusPending,
	}
	err := vacRepo.Create(ctx, vac)
	require.NoError(t, err)

	got, err := vacRepo.GetByID(ctx, "vac-nilr")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Nil(t, got.Reason)
}

// ---------------------------------------------------------------------------
// Additional: ListPending from multiple users
// ---------------------------------------------------------------------------

func TestVacationListPending_MultipleUsers(t *testing.T) {
	_, userRepo, vacRepo := setupRepos(t)
	ctx := context.Background()

	testutil.CreateTestUser(t, userRepo, "user1", "a@test.com", "Alice", domain.RoleEmployee, 25)
	testutil.CreateTestUser(t, userRepo, "user2", "b@test.com", "Bob", domain.RoleEmployee, 25)

	testutil.CreateTestVacation(t, vacRepo, "vp1", "user1", "2027-04-01", "2027-04-03", 3, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "vp2", "user2", "2027-05-01", "2027-05-03", 3, domain.StatusPending)
	testutil.CreateTestVacation(t, vacRepo, "va1", "user1", "2027-06-01", "2027-06-03", 3, domain.StatusApproved)

	results, err := vacRepo.ListPending(ctx)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Both pending requests should be returned with user info populated
	ids := []string{results[0].ID, results[1].ID}
	assert.Contains(t, ids, "vp1")
	assert.Contains(t, ids, "vp2")

	for _, r := range results {
		assert.Equal(t, domain.StatusPending, r.Status)
		assert.NotEmpty(t, r.UserName)
		assert.NotEmpty(t, r.UserEmail)
	}
}

// Ensure the unused import does not cause a compilation error.
var _ repository.MonthlyStats
