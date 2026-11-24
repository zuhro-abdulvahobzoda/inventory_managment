package user

import "inventory-api/internal/models"

type dbModel struct {
	ID           int64  `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"`
}

func (m dbModel) toModel() models.User {
	return models.User{
		ID:           m.ID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		Role:         m.Role,
	}
}
