package services

import (
	"testing"

	"github.com/levferril/racenotes/backend/internal/models"
)

func TestRaceService_Create(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Canyon"
	tires := "GP5000"
	race := &models.Race{
		UserID:   1,
		Name:     "Spring Race",
		Date:     "2026-04-01",
		Type:     "Road",
		BikeName: &bike,
		Tires:    &tires,
	}
	if err := svc.Create(race); err != nil {
		t.Fatalf("create race failed: %v", err)
	}
	if race.ID == 0 {
		t.Error("expected race ID to be set")
	}
}

func TestRaceService_CreateWithSetup(t *testing.T) {
	db := setupTestDB(t)
	raceSvc := NewRaceService(db)
	setupSvc := NewSetupService(db)

	setup := &models.Setup{UserID: 1, Name: "Race Setup", BikeName: "Canyon Aeroad", Tires: "GP5000 28mm"}
	setupSvc.Create(setup)

	race := &models.Race{
		UserID:  1,
		Name:    "Race with Setup",
		Date:    "2026-04-01",
		Type:    "Road",
		SetupID: &setup.ID,
	}
	if err := raceSvc.Create(race); err != nil {
		t.Fatalf("create race with setup failed: %v", err)
	}
	if race.BikeName == nil || *race.BikeName != "Canyon Aeroad" {
		t.Error("expected bike_name to be filled from setup")
	}
	if race.Tires == nil || *race.Tires != "GP5000 28mm" {
		t.Error("expected tires to be filled from setup")
	}
}

func TestRaceService_List(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Bike"
	tires := "Tires"
	svc.Create(&models.Race{UserID: 1, Name: "Race A", Date: "2026-04-01", Type: "Road", BikeName: &bike, Tires: &tires})
	svc.Create(&models.Race{UserID: 1, Name: "Race B", Date: "2026-04-02", Type: "MTB", BikeName: &bike, Tires: &tires, IsCompleted: true})
	svc.Create(&models.Race{UserID: 2, Name: "Race C", Date: "2026-04-03", Type: "Road", BikeName: &bike, Tires: &tires})

	races, err := svc.List(1, "", "", "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(races) != 2 {
		t.Errorf("expected 2 races for user 1, got %d", len(races))
	}
}

func TestRaceService_ListFilterByType(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Bike"
	tires := "Tires"
	svc.Create(&models.Race{UserID: 1, Name: "Race A", Date: "2026-04-01", Type: "Road", BikeName: &bike, Tires: &tires})
	svc.Create(&models.Race{UserID: 1, Name: "Race B", Date: "2026-04-02", Type: "MTB", BikeName: &bike, Tires: &tires})

	races, err := svc.List(1, "Road", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(races) != 1 {
		t.Errorf("expected 1 Road race, got %d", len(races))
	}
}

func TestRaceService_ListFilterByCompleted(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Bike"
	tires := "Tires"
	svc.Create(&models.Race{UserID: 1, Name: "Planned", Date: "2026-04-01", Type: "Road", BikeName: &bike, Tires: &tires})
	svc.Create(&models.Race{UserID: 1, Name: "Done", Date: "2026-04-02", Type: "Road", BikeName: &bike, Tires: &tires, IsCompleted: true})

	races, err := svc.List(1, "", "true", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(races) != 1 {
		t.Errorf("expected 1 completed race, got %d", len(races))
	}
	if races[0].Name != "Done" {
		t.Errorf("expected 'Done', got '%s'", races[0].Name)
	}
}

func TestRaceService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Bike"
	tires := "Tires"
	race := &models.Race{UserID: 1, Name: "Test", Date: "2026-04-01", Type: "Road", BikeName: &bike, Tires: &tires}
	svc.Create(race)

	found, err := svc.GetByID(1, race.ID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if found.Name != "Test" {
		t.Errorf("expected name Test, got %s", found.Name)
	}
}

func TestRaceService_GetByID_WrongUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Bike"
	tires := "Tires"
	race := &models.Race{UserID: 1, Name: "Test", Date: "2026-04-01", Type: "Road", BikeName: &bike, Tires: &tires}
	svc.Create(race)

	_, err := svc.GetByID(2, race.ID)
	if err == nil {
		t.Error("expected error accessing other user's race")
	}
}

func TestRaceService_Update(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Bike"
	tires := "Tires"
	race := &models.Race{UserID: 1, Name: "Old", Date: "2026-04-01", Type: "Road", BikeName: &bike, Tires: &tires}
	svc.Create(race)

	updated, err := svc.Update(1, race.ID, map[string]interface{}{"name": "New"})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "New" {
		t.Errorf("expected name New, got %s", updated.Name)
	}
}

func TestRaceService_Delete(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Bike"
	tires := "Tires"
	race := &models.Race{UserID: 1, Name: "ToDelete", Date: "2026-04-01", Type: "Road", BikeName: &bike, Tires: &tires}
	svc.Create(race)

	if err := svc.Delete(1, race.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := svc.GetByID(1, race.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestRaceService_Delete_WrongUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewRaceService(db)

	bike := "Bike"
	tires := "Tires"
	race := &models.Race{UserID: 1, Name: "Test", Date: "2026-04-01", Type: "Road", BikeName: &bike, Tires: &tires}
	svc.Create(race)

	err := svc.Delete(2, race.ID)
	if err == nil {
		t.Error("expected error deleting other user's race")
	}
}
