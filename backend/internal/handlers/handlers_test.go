package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/levferril/racenotes/backend/internal/middleware"
	"github.com/levferril/racenotes/backend/internal/models"
	"github.com/levferril/racenotes/backend/internal/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

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

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()

	authService := services.NewAuthService(db)
	authHandler := NewAuthHandler(authService)
	setupService := services.NewSetupService(db)
	setupHandler := NewSetupHandler(setupService)
	raceService := services.NewRaceService(db)
	raceHandler := NewRaceHandler(raceService)
	calcHandler := NewCalculatorHandler()

	api := r.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	user := api.Group("/user")
	user.Use(middleware.AuthRequired())
	user.GET("/profile", authHandler.GetProfile)
	user.PUT("/profile", authHandler.UpdateProfile)

	setups := api.Group("/setups")
	setups.Use(middleware.AuthRequired())
	setups.POST("", setupHandler.Create)
	setups.GET("", setupHandler.List)
	setups.GET("/:id", setupHandler.GetByID)
	setups.PUT("/:id", setupHandler.Update)
	setups.DELETE("/:id", setupHandler.Delete)

	races := api.Group("/races")
	races.Use(middleware.AuthRequired())
	races.POST("", raceHandler.Create)
	races.GET("", raceHandler.List)
	races.GET("/:id", raceHandler.GetByID)
	races.PUT("/:id", raceHandler.Update)
	races.DELETE("/:id", raceHandler.Delete)

	calculator := api.Group("/calculator")
	calculator.Use(middleware.AuthRequired())
	calculator.POST("/tire-pressure", calcHandler.Calculate)

	return r
}

func registerAndLogin(t *testing.T, router *gin.Engine) string {
	t.Helper()
	body, _ := json.Marshal(map[string]interface{}{
		"name": "Test User", "username": "testuser",
		"email": "test@example.com", "password": "password123",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	token := resp["token"]
	if token == "" {
		t.Fatalf("failed to get token from register: %s", w.Body.String())
	}
	return token
}

func authRequest(method, path, token string, body interface{}) *http.Request {
	var buf *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewReader(b)
	} else {
		buf = bytes.NewReader(nil)
	}
	req, _ := http.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// ===== Auth Integration Tests =====

func TestAPI_Register(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "Test", "username": "newuser",
		"email": "new@example.com", "password": "password123",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestAPI_RegisterValidation(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	body, _ := json.Marshal(map[string]interface{}{"name": "Test"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAPI_RegisterDuplicate(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "Test", "username": "dup",
		"email": "dup@example.com", "password": "password123",
	})
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w2.Code)
	}
}

func TestAPI_Login(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	registerAndLogin(t, router)

	body, _ := json.Marshal(map[string]string{"username": "testuser", "password": "password123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAPI_LoginWrongPassword(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	registerAndLogin(t, router)

	body, _ := json.Marshal(map[string]string{"username": "testuser", "password": "wrong"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAPI_GetProfile(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/user/profile", token, nil))

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var user map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &user)
	if user["username"] != "testuser" {
		t.Errorf("expected username testuser, got %v", user["username"])
	}
}

func TestAPI_GetProfileNoAuth(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/profile", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAPI_UpdateProfile(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("PUT", "/api/user/profile", token, map[string]interface{}{"name": "Updated"}))

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var user map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &user)
	if user["name"] != "Updated" {
		t.Errorf("expected name Updated, got %v", user["name"])
	}
}

// ===== Setup Integration Tests =====

func TestAPI_SetupCRUD(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	token := registerAndLogin(t, router)

	// Create
	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("POST", "/api/setups", token, map[string]interface{}{
		"name": "Race Setup", "bike_name": "Canyon", "tires": "GP5000",
	}))
	if w.Code != http.StatusCreated {
		t.Fatalf("create setup: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var setup map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &setup)
	setupID := setup["id"]

	// List
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/setups", token, nil))
	if w.Code != http.StatusOK {
		t.Errorf("list setups: expected 200, got %d", w.Code)
	}
	var setups []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &setups)
	if len(setups) != 1 {
		t.Errorf("expected 1 setup, got %d", len(setups))
	}

	// GetByID
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/setups/1", token, nil))
	if w.Code != http.StatusOK {
		t.Errorf("get setup: expected 200, got %d", w.Code)
	}

	// Update
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("PUT", "/api/setups/1", token, map[string]interface{}{"name": "Updated Setup"}))
	if w.Code != http.StatusOK {
		t.Errorf("update setup: expected 200, got %d", w.Code)
	}
	json.Unmarshal(w.Body.Bytes(), &setup)
	if setup["name"] != "Updated Setup" {
		t.Errorf("expected name 'Updated Setup', got '%v'", setup["name"])
	}

	// Delete
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("DELETE", "/api/setups/1", token, nil))
	if w.Code != http.StatusOK {
		t.Errorf("delete setup: expected 200, got %d", w.Code)
	}

	// Verify deleted
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/setups/1", token, nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("get deleted setup: expected 404, got %d", w.Code)
	}

	_ = setupID
}

func TestAPI_SetupCreateValidation(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("POST", "/api/setups", token, map[string]interface{}{"name": "No Bike"}))
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing bike_name/tires, got %d", w.Code)
	}
}

// ===== Race Integration Tests =====

func TestAPI_RaceCRUD(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	token := registerAndLogin(t, router)

	// Create
	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("POST", "/api/races", token, map[string]interface{}{
		"name": "Spring Race", "date": "2026-04-01", "type": "Road",
		"bike_name": "Canyon", "tires": "GP5000",
	}))
	if w.Code != http.StatusCreated {
		t.Fatalf("create race: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// List
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/races", token, nil))
	if w.Code != http.StatusOK {
		t.Errorf("list races: expected 200, got %d", w.Code)
	}
	var races []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &races)
	if len(races) != 1 {
		t.Errorf("expected 1 race, got %d", len(races))
	}

	// GetByID
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/races/1", token, nil))
	if w.Code != http.StatusOK {
		t.Errorf("get race: expected 200, got %d", w.Code)
	}

	// Update
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("PUT", "/api/races/1", token, map[string]interface{}{"name": "Updated Race"}))
	if w.Code != http.StatusOK {
		t.Errorf("update race: expected 200, got %d", w.Code)
	}

	// Delete
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("DELETE", "/api/races/1", token, nil))
	if w.Code != http.StatusOK {
		t.Errorf("delete race: expected 200, got %d", w.Code)
	}

	// Verify deleted
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/races/1", token, nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("get deleted race: expected 404, got %d", w.Code)
	}
}

func TestAPI_RaceCreateValidation(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	token := registerAndLogin(t, router)

	// Missing required fields
	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("POST", "/api/races", token, map[string]interface{}{"name": "Test"}))
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}

	// Missing bike_name/tires without setup_id
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("POST", "/api/races", token, map[string]interface{}{
		"name": "Test", "date": "2026-04-01", "type": "Road",
	}))
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing bike/tires, got %d", w.Code)
	}
}

func TestAPI_RaceListFilters(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	token := registerAndLogin(t, router)

	// Create two races
	router.ServeHTTP(httptest.NewRecorder(), authRequest("POST", "/api/races", token, map[string]interface{}{
		"name": "Road Race", "date": "2026-04-01", "type": "Road",
		"bike_name": "Canyon", "tires": "GP5000",
	}))
	router.ServeHTTP(httptest.NewRecorder(), authRequest("POST", "/api/races", token, map[string]interface{}{
		"name": "MTB Race", "date": "2026-04-02", "type": "MTB",
		"bike_name": "Trek", "tires": "Maxxis", "is_completed": true,
	}))

	// Filter by type
	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/races?type=Road", token, nil))
	var races []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &races)
	if len(races) != 1 {
		t.Errorf("type filter: expected 1 race, got %d", len(races))
	}

	// Filter by completed
	w = httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/races?is_completed=true", token, nil))
	json.Unmarshal(w.Body.Bytes(), &races)
	if len(races) != 1 {
		t.Errorf("completed filter: expected 1 race, got %d", len(races))
	}
}

// ===== Calculator Integration Test =====

func TestAPI_Calculator(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	token := registerAndLogin(t, router)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("POST", "/api/calculator/tire-pressure", token, map[string]interface{}{
		"rider_weight": 70, "bike_weight": 8, "tire_width": 28,
		"tire_type": "clincher", "surface": "road", "conditions": "dry",
	}))

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	if result["front_pressure"] == nil {
		t.Error("expected front_pressure in response")
	}
	if result["rear_pressure"] == nil {
		t.Error("expected rear_pressure in response")
	}
}

func TestAPI_CalculatorNoAuth(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/calculator/tire-pressure", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// ===== Middleware Test =====

func TestAPI_InvalidToken(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, authRequest("GET", "/api/user/profile", "invalid-token", nil))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAPI_MalformedAuthHeader(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user/profile", nil)
	req.Header.Set("Authorization", "NotBearer token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
