package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
)

type AuthRepoInterface interface {
	GetByEmail(email string) (*models.UsersRes, error)
	RegisterUser(data *models.Users) (string, error)
}

type AuthRepository struct {
	*sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}
func (r *AuthRepository) RegisterUser(data *models.Users) (string, error) {
	query := `
        INSERT INTO users (
    			"email", 
    			"password",
					"role"
				) VALUES (
    			:email, 
    			:password,
					:role
				)
				RETURNING "uuid", "email", "role";
    		`

	rows, err := r.DB.NamedQuery(query, data)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	return "Register user is success", nil
}

func (r *AuthRepository) GetByEmail(email string) (*models.UsersRes, error) {
	result := models.UsersRes{}
	query := `SELECT * FROM public.users WHERE email=$1`
	err := r.Select(&result, query, email)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
