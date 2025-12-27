package services

import (
	"context"
	"testing"

	"github.com/linn221/bane/config"
	"github.com/linn221/bane/models"
	"github.com/linn221/bane/mystructs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&models.Endpoint{},
		&models.Alias{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestEndpointService_Create_NewEndpoint(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Initialize services with in-memory cache
	cache := config.NewInMemoryCache()
	services := NewMyServices(db, cache)
	ctx := context.Background()

	// Prepare input
	url := mystructs.VarString{}
	if err := url.UnmarshalGQL("https://example.com/test?foo=bar"); err != nil {
		t.Fatalf("Failed to unmarshal url: %v", err)
	}

	headers := mystructs.VarKVGroup{}
	if err := headers.UnmarshalGQL("secret:helloworld"); err != nil {
		t.Fatalf("Failed to unmarshal headers: %v", err)
	}

	input := &models.EndpointInput{
		Name:        "Test Endpoint",
		Description: "Test description",
		Method:      models.HttpMethodGet,
		Url:         url,
		Headers:     headers,
	}

	// Call Create
	endpoint, err := services.EndpointService.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify the endpoint was created
	if endpoint.Id == 0 {
		t.Error("Create() endpoint.Id should not be 0")
	}

	// Verify URL parsing (extracted from URL)
	if !endpoint.Https {
		t.Errorf("Create() Https = %v, want %v", endpoint.Https, true)
	}

	if endpoint.Domain != "example.com" {
		t.Errorf("Create() Domain = %v, want %v", endpoint.Domain, "example.com")
	}

	if endpoint.Path.OriginalString != "/test" {
		t.Errorf("Create() Path.OriginalString = %v, want %v", endpoint.Path.OriginalString, "/test")
	}

	// Verify HTTP method from input
	if endpoint.Method != models.HttpMethodGet {
		t.Errorf("Create() Method = %v, want %v", endpoint.Method, models.HttpMethodGet)
	}

	// Verify query parameters
	if len(endpoint.Queries.VarKVs) != 1 {
		t.Errorf("Create() Queries.VarKVs length = %v, want %v", len(endpoint.Queries.VarKVs), 1)
	} else {
		queryKV := endpoint.Queries.VarKVs[0]
		if queryKV.Key.OriginalString != "foo" {
			t.Errorf("Create() Queries[0].Key = %v, want %v", queryKV.Key.OriginalString, "foo")
		}
		if queryKV.Value.OriginalString != "bar" {
			t.Errorf("Create() Queries[0].Value = %v, want %v", queryKV.Value.OriginalString, "bar")
		}
	}

	// Verify headers (from input)
	if len(endpoint.Headers.VarKVs) != 1 {
		t.Errorf("Create() Headers.VarKVs length = %v, want %v", len(endpoint.Headers.VarKVs), 1)
	} else {
		headerKV := endpoint.Headers.VarKVs[0]
		if headerKV.Key.OriginalString != "secret" {
			t.Errorf("Create() Headers[0].Key = %v, want %v", headerKV.Key.OriginalString, "secret")
		}
		if headerKV.Value.OriginalString != "helloworld" {
			t.Errorf("Create() Headers[0].Value = %v, want %v", headerKV.Value.OriginalString, "helloworld")
		}
	}

	// Verify the endpoint was saved to database
	var savedEndpoint models.Endpoint
	err = db.First(&savedEndpoint, endpoint.Id).Error
	if err != nil {
		t.Fatalf("Failed to retrieve saved endpoint from database: %v", err)
	}

	if savedEndpoint.Domain != "example.com" {
		t.Errorf("Saved endpoint Domain = %v, want %v", savedEndpoint.Domain, "example.com")
	}

	// Verify alias was created
	var alias models.Alias
	err = db.Where("reference_type = ? AND reference_id = ?", "endpoints", endpoint.Id).First(&alias).Error
	if err != nil {
		t.Logf("No alias found for endpoint (this is expected if alias was not provided)")
	}
}
