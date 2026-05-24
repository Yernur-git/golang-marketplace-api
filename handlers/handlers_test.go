package handlers

import (
	"Marketplace-API/config"
	"Marketplace-API/middleware"
	"Marketplace-API/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	db.AutoMigrate(&models.User{}, &models.Category{}, &models.Listing{})
	config.DB = db
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/api/register", CreateUser)
	r.POST("/api/login", Login)
	r.GET("/api/categories", GetCategories)
	r.GET("/api/listings", GetListings)
	r.GET("/api/listings/:id", GetListingByID)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/categories", CreateCategory)
		protected.POST("/listings", CreateListing)
	}

	return r
}

func seedTestData() {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:         "Test User",
		Email:        "test@test.com",
		Phone:        "+7777777777",
		PasswordHash: string(hashed),
	}
	config.DB.Create(&user)

	category := models.Category{Name: "Cameras"}
	config.DB.Create(&category)

	listing := models.Listing{
		Title:       "Canon EOS R5",
		Description: "Great camera",
		Price:       500000,
		Status:      models.StatusActive,
		Location:    "Almaty",
		UserID:      user.ID,
		CategoryID:  category.ID,
	}
	config.DB.Create(&listing)
}

func TestRegisterSuccess(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

	body := `{"name":"John","email":"john@example.com","password":"123456","phone":"+7701234567"}`
	req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["email"] != "john@example.com" {
		t.Errorf("expected email john@example.com, got %v", resp["email"])
	}
}

func TestRegisterDuplicate(t *testing.T) {
	setupTestDB()
	seedTestData()
	router := setupTestRouter()

	body := `{"name":"Duplicate","email":"test@test.com","password":"123456"}`
	req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 for duplicate email, got %d", w.Code)
	}
}

func TestLoginSuccess(t *testing.T) {
	setupTestDB()
	seedTestData()
	router := setupTestRouter()

	body := `{"email":"test@test.com","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == nil || resp["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	setupTestDB()
	seedTestData()
	router := setupTestRouter()

	body := `{"email":"test@test.com","password":"wrongpassword"}`
	req, _ := http.NewRequest("POST", "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestGetCategories(t *testing.T) {
	setupTestDB()
	seedTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var categories []models.Category
	json.Unmarshal(w.Body.Bytes(), &categories)
	if len(categories) == 0 {
		t.Error("expected at least one category")
	}
}

func TestCreateCategory(t *testing.T) {
	setupTestDB()
	seedTestData()
	router := setupTestRouter()

	loginBody := `{"email":"test@test.com","password":"password123"}`
	loginReq, _ := http.NewRequest("POST", "/api/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	token := loginResp["token"].(string)

	body := `{"name":"Lenses"}`
	req, _ := http.NewRequest("POST", "/api/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestGetListings(t *testing.T) {
	setupTestDB()
	seedTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/listings?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["data"] == nil {
		t.Error("expected data field in response")
	}
	if resp["total"] == nil {
		t.Error("expected total field in response")
	}
}

func TestGetListingByID(t *testing.T) {
	setupTestDB()
	seedTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/listings/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetListingNotFound(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/listings/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestCreateListingUnauth(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

	body := `{"title":"Camera","price":100}`
	req, _ := http.NewRequest("POST", "/api/listings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}
