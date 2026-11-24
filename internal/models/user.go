package models

const (
	RoleAdmin    = "admin"
	RoleEmployee = "employee"
)

func IsValidRole(role string) bool {
	return role == RoleAdmin || role == RoleEmployee
}

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         string
}
