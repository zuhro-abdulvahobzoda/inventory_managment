package auth

import (
	"context"
	"errors"
	"inventory-api/internal/models"
	"inventory-api/internal/repository"
	jwttoken "inventory-api/internal/service/token/jwt"
	"strings"
	"testing"
	"time"
)

type fakeUserRepo struct {
	users map[string]models.User
}

var _ repository.UserRepository = (*fakeUserRepo)(nil)

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[string]models.User)}
}

func (r *fakeUserRepo) CreateUser(_ context.Context, user *models.User) (int64, error) {
	if _, ok := r.users[user.Username]; ok {
		return 0, models.ErrUsernameAlreadyExists
	}
	if user.Role == models.RoleAdmin {
		for _, u := range r.users {
			if u.Role == models.RoleAdmin {
				return 0, models.ErrAdminAlreadyExists
			}
		}
	}
	user.ID = int64(len(r.users) + 1)
	r.users[user.Username] = *user
	return user.ID, nil
}

func (r *fakeUserRepo) GetByUsername(_ context.Context, username string) (*models.User, error) {
	user, ok := r.users[username]
	if !ok {
		return nil, models.ErrUserNotFound
	}
	return &user, nil
}

func (r *fakeUserRepo) GetByID(_ context.Context, id int64) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, models.ErrUserNotFound
}

func TestSignUp(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{
			name:     "success",
			username: "bob",
			password: "12345678",
		},
		{
			name:     "password counted in runes, not bytes",
			username: "bob",
			password: "пароль12",
		},
		{
			name:     "password too short",
			username: "bob",
			password: "1234567",
			wantErr:  models.ErrPasswordTooShort,
		},
		{
			name:     "multibyte password too short",
			username: "bob",
			password: "пароль1",
			wantErr:  models.ErrPasswordTooShort,
		},
		{
			name:     "password at max length",
			username: "bob",
			password: strings.Repeat("a", 72),
		},
		{
			name:     "password too long",
			username: "bob",
			password: strings.Repeat("a", 73),
			wantErr:  models.ErrPasswordTooLong,
		},
		{
			name:     "multibyte password too long in bytes",
			username: "bob",
			password: strings.Repeat("я", 37),
			wantErr:  models.ErrPasswordTooLong,
		},
		{
			name:     "blank username",
			username: "   ",
			password: "12345678",
			wantErr:  models.ErrInvalidInput,
		},
		{
			name:     "empty password",
			username: "bob",
			wantErr:  models.ErrInvalidInput,
		},
		{
			name:     "username too long",
			username: strings.Repeat("a", 101),
			password: "12345678",
			wantErr:  models.ErrInvalidUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAuthService(newFakeUserRepo(), jwttoken.NewTokenManager("test-secret", time.Hour, time.Hour))

			_, err := svc.SignUp(context.Background(), tt.username, tt.password, "")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestSignUpDuplicateUsernameIgnoresCase(t *testing.T) {
	ctx := context.Background()
	svc := NewAuthService(newFakeUserRepo(), jwttoken.NewTokenManager("test-secret", time.Hour, time.Hour))

	if _, err := svc.SignUp(ctx, "Bob", "12345678", ""); err != nil {
		t.Fatalf("first sign up: %v", err)
	}

	_, err := svc.SignUp(ctx, "  bOB ", "12345678", "")
	if !errors.Is(err, models.ErrUsernameAlreadyExists) {
		t.Fatalf("err = %v, want %v", err, models.ErrUsernameAlreadyExists)
	}
}

func TestSignUpRole(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		role     string
		wantRole string
		wantErr  error
	}{
		{name: "empty role defaults to employee", role: "", wantRole: models.RoleEmployee},
		{name: "employee", role: "employee", wantRole: models.RoleEmployee},
		{name: "admin", role: "admin", wantRole: models.RoleAdmin},
		{name: "role is case insensitive and trimmed", role: "  ADMIN ", wantRole: models.RoleAdmin},
		{name: "unknown role", role: "superuser", wantErr: models.ErrInvalidRole},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeUserRepo()
			svc := NewAuthService(repo, jwttoken.NewTokenManager("test-secret", time.Hour, time.Hour))

			_, err := svc.SignUp(ctx, "bob", "12345678", tt.role)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got := repo.users["bob"].Role; got != tt.wantRole {
				t.Fatalf("role = %q, want %q", got, tt.wantRole)
			}
		})
	}
}

func TestFirstUserIsNotAutomaticallyAdmin(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewAuthService(repo, jwttoken.NewTokenManager("test-secret", time.Hour, time.Hour))

	if _, err := svc.SignUp(context.Background(), "first", "12345678", ""); err != nil {
		t.Fatalf("sign up: %v", err)
	}
	if got := repo.users["first"].Role; got != models.RoleEmployee {
		t.Fatalf("role = %q, want %q", got, models.RoleEmployee)
	}
}

func TestOnlyOneAdmin(t *testing.T) {
	ctx := context.Background()
	svc := NewAuthService(newFakeUserRepo(), jwttoken.NewTokenManager("test-secret", time.Hour, time.Hour))

	if _, err := svc.SignUp(ctx, "admin", "12345678", models.RoleAdmin); err != nil {
		t.Fatalf("first admin: %v", err)
	}

	_, err := svc.SignUp(ctx, "admin2", "12345678", models.RoleAdmin)
	if !errors.Is(err, models.ErrAdminAlreadyExists) {
		t.Fatalf("err = %v, want %v", err, models.ErrAdminAlreadyExists)
	}

	if _, err := svc.SignUp(ctx, "worker", "12345678", models.RoleEmployee); err != nil {
		t.Fatalf("employee after admin: %v", err)
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	svc := NewAuthService(newFakeUserRepo(), jwttoken.NewTokenManager("test-secret", time.Hour, time.Hour))

	if _, err := svc.SignUp(ctx, "bob", "12345678", ""); err != nil {
		t.Fatalf("sign up: %v", err)
	}

	tests := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{name: "success", username: "bob", password: "12345678"},
		{name: "username is case insensitive", username: "BOB", password: "12345678"},
		{name: "wrong password", username: "bob", password: "wrong-password", wantErr: models.ErrInvalidCredentials},
		{name: "unknown user", username: "alice", password: "12345678", wantErr: models.ErrInvalidCredentials},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Login(ctx, tt.username, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()
	tokens := jwttoken.NewTokenManager("test-secret", time.Hour, time.Hour)
	svc := NewAuthService(newFakeUserRepo(), tokens)

	adminID, err := svc.SignUp(ctx, "admin", "12345678", models.RoleAdmin)
	if err != nil {
		t.Fatalf("sign up admin: %v", err)
	}
	employeeID, err := svc.SignUp(ctx, "employee", "12345678", "")
	if err != nil {
		t.Fatalf("sign up employee: %v", err)
	}

	adminTokens, err := svc.Login(ctx, "admin", "12345678")
	if err != nil {
		t.Fatalf("login admin: %v", err)
	}
	employeeTokens, err := svc.Login(ctx, "employee", "12345678")
	if err != nil {
		t.Fatalf("login employee: %v", err)
	}

	// The token belongs to a user that is not in the database (e.g. the user was deleted).
	ghostToken, err := tokens.GenerateAccessToken(ctx, 999)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	tests := []struct {
		name     string
		token    string
		wantID   int64
		wantRole string
		wantErr  error
	}{
		{
			name:     "user who signed up as admin",
			token:    adminTokens.AccessToken,
			wantID:   adminID,
			wantRole: models.RoleAdmin,
		},
		{
			name:     "user who signed up without a role is an employee",
			token:    employeeTokens.AccessToken,
			wantID:   employeeID,
			wantRole: models.RoleEmployee,
		},
		{
			name:    "refresh token is not accepted",
			token:   adminTokens.RefreshToken,
			wantErr: models.ErrInvalidToken,
		},
		{
			name:    "token of a missing user",
			token:   ghostToken,
			wantErr: models.ErrInvalidToken,
		},
		{
			name:    "unknown token",
			token:   "unknown",
			wantErr: models.ErrInvalidToken,
		},
		{
			name:    "blank token",
			token:   "   ",
			wantErr: models.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.Authenticate(ctx, tt.token)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if user.ID != tt.wantID || user.Role != tt.wantRole {
				t.Fatalf("user = {id: %d, role: %q}, want {id: %d, role: %q}", user.ID, user.Role, tt.wantID, tt.wantRole)
			}
		})
	}
}
