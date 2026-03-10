package services

import (
	"testing"

	"github.com/levferril/racenotes/backend/internal/middleware"
	"github.com/levferril/racenotes/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Setup{}, &models.Race{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	middleware.JWTSecret = []byte("test-secret")
	return db
}

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB(t)
	svc := NewAuthService(db)

	user := &models.User{
		Name:     "Test User",
		Username: "testuser",
		Email:    "test@example.com",
	}
	err := svc.Register(user, "password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected user ID to be set")
	}
	if user.Password == "password123" {
		t.Error("password should be hashed, not plaintext")
	}
}

func TestAuthService_RegisterDuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	svc := NewAuthService(db)

	user1 := &models.User{Name: "User1", Username: "testuser", Email: "test1@example.com"}
	if err := svc.Register(user1, "password123"); err != nil {
		t.Fatal(err)
	}

	user2 := &models.User{Name: "User2", Username: "testuser", Email: "test2@example.com"}
	err := svc.Register(user2, "password123")
	if err == nil {
		t.Error("expected error for duplicate username")
	}
	if err.Error() != "username already exists" {
		t.Errorf("expected 'username already exists', got '%s'", err.Error())
	}
}

func TestAuthService_RegisterDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	svc := NewAuthService(db)

	user1 := &models.User{Name: "User1", Username: "user1", Email: "same@example.com"}
	if err := svc.Register(user1, "password123"); err != nil {
		t.Fatal(err)
	}

	user2 := &models.User{Name: "User2", Username: "user2", Email: "same@example.com"}
	err := svc.Register(user2, "password123")
	if err == nil {
		t.Error("expected error for duplicate email")
	}
	if err.Error() != "email already exists" {
		t.Errorf("expected 'email already exists', got '%s'", err.Error())
	}
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB(t)
	svc := NewAuthService(db)

	user := &models.User{Name: "Test", Username: "testuser", Email: "test@example.com"}
	if err := svc.Register(user, "password123"); err != nil {
		t.Fatal(err)
	}

	token, err := svc.Login("testuser", "password123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if token == "" {
		t.Error("expected token, got empty string")
	}
}

func TestAuthService_LoginWrongPassword(t *testing.T) {
	db := setupTestDB(t)
	svc := NewAuthService(db)

	user := &models.User{Name: "Test", Username: "testuser", Email: "test@example.com"}
	if err := svc.Register(user, "password123"); err != nil {
		t.Fatal(err)
	}

	_, err := svc.Login("testuser", "wrongpassword")
	if err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestAuthService_LoginWrongUsername(t *testing.T) {
	db := setupTestDB(t)
	svc := NewAuthService(db)

	_, err := svc.Login("nonexistent", "password123")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestAuthService_GetProfile(t *testing.T) {
	db := setupTestDB(t)
	svc := NewAuthService(db)

	user := &models.User{Name: "Test", Username: "testuser", Email: "test@example.com", Height: 180, Weight: 75.5}
	if err := svc.Register(user, "pass"); err != nil {
		t.Fatal(err)
	}

	profile, err := svc.GetProfile(user.ID)
	if err != nil {
		t.Fatalf("get profile failed: %v", err)
	}
	if profile.Name != "Test" {
		t.Errorf("expected name Test, got %s", profile.Name)
	}
	if profile.Height != 180 {
		t.Errorf("expected height 180, got %d", profile.Height)
	}
}

func TestAuthService_UpdateProfile(t *testing.T) {
	db := setupTestDB(t)
	svc := NewAuthService(db)

	user := &models.User{Name: "Test", Username: "testuser", Email: "test@example.com"}
	if err := svc.Register(user, "pass"); err != nil {
		t.Fatal(err)
	}

	updated, err := svc.UpdateProfile(user.ID, map[string]interface{}{"name": "Updated Name"})
	if err != nil {
		t.Fatalf("update profile failed: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", updated.Name)
	}
}
