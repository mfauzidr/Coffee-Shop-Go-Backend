package repository

import (
	"fmt"
	"strings"

	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepoInterface interface {
	CreateUser(data *models.Users) (*models.Users, error)
	GetAllUsers(query *models.UsersQuery) (*models.UsersRes, int, error)
	GetDetailsUser(uuid string) (*models.Users, error)
	UpdateUser(uuid string, data map[string]interface{}) (*models.Users, error)
	DeleteUser(uuid string) (*models.Users, error)
}

type UsersRepo struct {
	*sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db}
}

func (r *UsersRepo) CreateUser(data *models.Users) (*models.Users, error) {
	query := `
        INSERT INTO users (
    			"firstName", 
    			"email", 
    			"password", 
    			"role"
				) VALUES (
    			:firstName, 
    			:email, 
    			:password, 
    			:role
				)
				RETURNING *;
    		`

	var result models.Users
	err := r.DB.QueryRowx(query, data).StructScan(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UsersRepo) GetAllUsers(query *models.UsersQuery) (*models.UsersRes, int, error) {
	baseQuery := `SELECT * FROM public.users`
	countQuery := `SELECT COUNT(*) FROM public.users`
	whereClauses := []string{}
	var values []interface{}

	if query.Search != nil {
		searchTerm := "%" + *query.Search + "%"
		whereClauses = append(whereClauses, fmt.Sprintf(`("firstName" ILIKE $%d OR "lastName" ILIKE $%d)`, len(values)+1, len(values)+2))
		values = append(values, searchTerm, searchTerm)
	}

	if len(whereClauses) > 0 {
		whereQuery := " WHERE " + strings.Join(whereClauses, " AND ")
		baseQuery += whereQuery
		countQuery += whereQuery
	}

	if query.Sort != nil && (*query.Sort == "ASC" || *query.Sort == "DESC") {
		baseQuery += ` ORDER BY "createdAt" ` + *query.Sort
	} else {
		baseQuery += ` ORDER BY "createdAt" DESC`
	}

	if query.Page > 0 && query.Limit > 0 {
		limit := query.Limit
		offset := (query.Page - 1) * limit
		baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(values)+1, len(values)+2)
		values = append(values, limit, offset)
	}

	data := models.UsersRes{}
	if err := r.Select(&data, baseQuery, values...); err != nil {
		return nil, 0, err
	}

	var total int
	if err := r.Get(&total, countQuery, values[:len(values)-2]...); err != nil {
		return nil, 0, err
	}

	return &data, total, nil
}

func (r *UsersRepo) GetDetailsUser(uuid string) (*models.Users, error) {
	query := `SELECT
							"uuid", 
							"firstName", 
    	        "lastName", 
    	        "gender", 
    	        "email",
    	        "image", 
    	        "address", 
    	        "phoneNumber", 
    	        "birthday",
    	        "deliveryAddress", 
    	        "role", 
    	        "createdAt", 
    	        "updatedAt" 
						FROM 
							public.users 
						WHERE uuid = :uuid`
	data := models.Users{}

	rows, err := r.DB.NamedQuery(query, map[string]interface{}{
		"uuid": uuid,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}

	return nil, nil
}

func (r *UsersRepo) UpdateUser(uuid string, data map[string]interface{}) (*models.Users, error) {
	var setClauses []string
	var values []interface{}

	values = append(values, uuid)

	i := 2
	for key, value := range data {
		setClauses = append(setClauses, fmt.Sprintf(`"%s" = $%d`, key, i))
		values = append(values, value)
		i++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
        UPDATE public.users
        SET %s, "updatedAt" = now()
        WHERE "uuid" = $1
        RETURNING *
    `, strings.Join(setClauses, ", "))

	fmt.Println(query)
	fmt.Println(values)

	var updatedUser models.Users
	err := r.DB.QueryRowx(query, values...).StructScan(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *UsersRepo) DeleteUser(uuid string) (*models.Users, error) {
	query := `DELETE FROM public.users WHERE uuid = $1 RETURNING *`

	var deletedUser models.Users
	if err := r.DB.QueryRowx(query, uuid).StructScan(&deletedUser); err != nil {
		return nil, err
	}

	return &deletedUser, nil
}
