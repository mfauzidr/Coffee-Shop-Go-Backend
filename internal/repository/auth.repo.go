package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
)

type AuthRepoInterface interface {
	GetByEmail(email string) (*models.Users, error)
	RegisterUser(data *models.Users) (*models.Users, error)
}

type AuthRepository struct {
	*sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}
func (r *AuthRepository) RegisterUser(data *models.Users) (*models.Users, error) {
	query := `
        INSERT INTO users (
    			"email", 
    			"password",
					"role"
				) VALUES (
    			:email, 
    			:password,
					"customer"
				)
				RETURNING "uud", "email", "role";
    		`

	var result models.Users
	err := r.DB.QueryRowx(query, data).StructScan(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AuthRepository) GetByEmail(email string) (*models.Users, error) {
	result := models.Users{}
	query := `SELECT * FROM public.users WHERE email=$1`
	err := r.Get(&result, query, email)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
