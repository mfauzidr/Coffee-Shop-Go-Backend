package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
)

type AuthRepoInterface interface {
	GetByEmail(email string) (*models.UsersRes, error)
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
					'customer'
				)
				RETURNING "uuid", "email", "role";
    		`

	var result models.Users
	rows, err := r.DB.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
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
