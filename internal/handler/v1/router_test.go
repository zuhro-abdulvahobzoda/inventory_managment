package v1

import (
	"context"
	"inventory-api/internal/models"
	"inventory-api/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeAuth struct{ service.AuthService }

func (fakeAuth) Authenticate(_ context.Context, token string) (*models.User, error) {
	switch token {
	case "admin-token":
		return &models.User{ID: 1, Username: "admin", Role: models.RoleAdmin}, nil
	case "employee-token":
		return &models.User{ID: 2, Username: "employee", Role: models.RoleEmployee}, nil
	}
	return nil, models.ErrInvalidToken
}

type fakeItems struct{ service.ItemService }

func (fakeItems) ListItems(context.Context) ([]models.Item, error) {
	return []models.Item{{ID: 1, Name: "Laptop", Quantity: 5}}, nil
}

func (fakeItems) CreateItem(context.Context, int64, models.Item) (int64, error) { return 1, nil }

func (fakeItems) UpdateItem(context.Context, int64, models.Item) error { return nil }

func (fakeItems) DeleteItem(context.Context, int64, int64) error { return nil }

type fakeLocations struct{ service.LocationService }

func (fakeLocations) CreateLocation(context.Context, models.Location) (int64, error) { return 1, nil }

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	h := NewHandler(&service.Services{
		AuthService:     fakeAuth{},
		ItemService:     fakeItems{},
		LocationService: fakeLocations{},
	})

	router := gin.New()
	h.InitRoutes(router.Group("/api/v1"))
	return router
}

func TestAccess(t *testing.T) {
	const itemBody = `{"name":"Laptop","quantity":5}`
	const locationBody = `{"name":"Warehouse","latitude":38.5598,"longitude":68.787}`

	tests := []struct {
		name       string
		method     string
		path       string
		token      string
		body       string
		wantStatus int
	}{
		// 401: token not found or invalid token
		{"list without token", http.MethodGet, "/items", "", "", http.StatusUnauthorized},
		{"list with bad token", http.MethodGet, "/items", "bad", "", http.StatusUnauthorized},
		{"create without token", http.MethodPost, "/items", "", itemBody, http.StatusUnauthorized},

		// employee only sees the inventory
		{"employee lists items", http.MethodGet, "/items", "employee-token", "", http.StatusOK},
		{"employee cannot create item", http.MethodPost, "/items", "employee-token", itemBody, http.StatusForbidden},
		{"employee cannot update item", http.MethodPut, "/items/1", "employee-token", itemBody, http.StatusForbidden},
		{"employee cannot delete item", http.MethodDelete, "/items/1", "employee-token", "", http.StatusForbidden},
		{"employee cannot create location", http.MethodPost, "/locations", "employee-token", locationBody, http.StatusForbidden},

		// admin, well, can do anything
		{"admin lists items", http.MethodGet, "/items", "admin-token", "", http.StatusOK},
		{"admin creates item", http.MethodPost, "/items", "admin-token", itemBody, http.StatusCreated},
		{"admin creates item with zero quantity", http.MethodPost, "/items", "admin-token", `{"name":"Mouse","quantity":0}`, http.StatusCreated},
		{"admin updates item", http.MethodPut, "/items/1", "admin-token", itemBody, http.StatusNoContent},
		{"admin deletes item", http.MethodDelete, "/items/1", "admin-token", "", http.StatusNoContent},
		{"admin creates location", http.MethodPost, "/locations", "admin-token", locationBody, http.StatusCreated},
		{"admin creates location on the equator", http.MethodPost, "/locations", "admin-token", `{"name":"Warehouse","latitude":0,"longitude":0}`, http.StatusCreated},

		// invalid requests
		{"quantity is required", http.MethodPost, "/items", "admin-token", `{"name":"Laptop"}`, http.StatusBadRequest},
		{"bad item id", http.MethodDelete, "/items/abc", "admin-token", "", http.StatusBadRequest},
		{"location without coordinates", http.MethodPost, "/locations", "admin-token", `{"name":"Warehouse"}`, http.StatusBadRequest},
		{"location without longitude", http.MethodPost, "/locations", "admin-token", `{"name":"Warehouse","latitude":10}`, http.StatusBadRequest},
	}

	router := newTestRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/v1"+tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d, body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}
