package services

import (
	"testing"

	"github.com/levferril/racenotes/backend/internal/models"
)

func TestSetupService_Create(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSetupService(db)

	setup := &models.Setup{
		UserID:   1,
		Name:     "Race Setup",
		BikeName: "Canyon Aeroad",
		Tires:    "GP5000 28mm",
	}
	if err := svc.Create(setup); err != nil {
		t.Fatalf("create setup failed: %v", err)
	}
	if setup.ID == 0 {
		t.Error("expected setup ID to be set")
	}
}

func TestSetupService_List(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSetupService(db)

	svc.Create(&models.Setup{UserID: 1, Name: "Setup A", BikeName: "Bike A", Tires: "Tires A"})
	svc.Create(&models.Setup{UserID: 1, Name: "Setup B", BikeName: "Bike B", Tires: "Tires B"})
	svc.Create(&models.Setup{UserID: 2, Name: "Setup C", BikeName: "Bike C", Tires: "Tires C"})

	setups, err := svc.List(1)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(setups) != 2 {
		t.Errorf("expected 2 setups for user 1, got %d", len(setups))
	}
}

func TestSetupService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSetupService(db)

	setup := &models.Setup{UserID: 1, Name: "Test", BikeName: "Bike", Tires: "Tires"}
	svc.Create(setup)

	found, err := svc.GetByID(1, setup.ID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if found.Name != "Test" {
		t.Errorf("expected name Test, got %s", found.Name)
	}
}

func TestSetupService_GetByID_WrongUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSetupService(db)

	setup := &models.Setup{UserID: 1, Name: "Test", BikeName: "Bike", Tires: "Tires"}
	svc.Create(setup)

	_, err := svc.GetByID(2, setup.ID)
	if err == nil {
		t.Error("expected error accessing other user's setup")
	}
}

func TestSetupService_Update(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSetupService(db)

	setup := &models.Setup{UserID: 1, Name: "Old", BikeName: "Bike", Tires: "Tires"}
	svc.Create(setup)

	updated, err := svc.Update(1, setup.ID, map[string]interface{}{"name": "New"})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "New" {
		t.Errorf("expected name New, got %s", updated.Name)
	}
}

func TestSetupService_Delete(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSetupService(db)

	setup := &models.Setup{UserID: 1, Name: "ToDelete", BikeName: "Bike", Tires: "Tires"}
	svc.Create(setup)

	if err := svc.Delete(1, setup.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := svc.GetByID(1, setup.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestSetupService_Delete_WrongUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSetupService(db)

	setup := &models.Setup{UserID: 1, Name: "Test", BikeName: "Bike", Tires: "Tires"}
	svc.Create(setup)

	err := svc.Delete(2, setup.ID)
	if err == nil {
		t.Error("expected error deleting other user's setup")
	}
}
