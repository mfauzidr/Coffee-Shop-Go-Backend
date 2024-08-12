package repository

import (
	"fmt"
	"strings"

	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepoInterface interface {
	CreateProduct(data map[string]interface{}) ([]models.Product, error)
	GetAllProduct(query *models.ProductQuery) (*models.Products, int, error)
	GetDetailProduct(uuid string) (*models.Product, error)
	UpdateProduct(uuid string, data map[string]interface{}) (*models.Product, error)
	DeleteProduct(uuid string) (*models.Product, error)
}

type RepoProduct struct {
	*sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *RepoProduct {
	return &RepoProduct{db}
}

func (r *RepoProduct) CreateProduct(data map[string]interface{}) ([]models.Product, error) {
	var columns []string
	var values []interface{}
	var placeholders []string

	i := 1
	for key, value := range data {
		columns = append(columns, fmt.Sprintf(`"%s"`, key))
		values = append(values, value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		i++
	}

	query := fmt.Sprintf(`
        INSERT INTO public.products
        (%s)
        VALUES
        (%s)
        RETURNING *
    `, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	var users []models.Product
	err := r.DB.Select(&users, query, values...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *RepoProduct) GetAllProduct(query *models.ProductQuery) (*models.Products, int, error) {
	baseQuery := `SELECT * FROM public.products`
	countQuery := `SELECT COUNT(*) FROM public.products`
	whereClauses := []string{}
	var values []interface{}

	if query.Name != nil {
		searchTerm := "%" + *query.Name + "%"
		whereClauses = append(whereClauses, fmt.Sprintf(`("name" ILIKE $%d)`, len(values)+1))
		values = append(values, searchTerm)
	}

	if query.Category != nil {
		searchTerm := "%" + *query.Category + "%"
		whereClauses = append(whereClauses, fmt.Sprintf(`("category" = $%d)`, len(values)+1))
		values = append(values, searchTerm)
	}

	if query.MinPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf(`("price" >= $%d)`, len(values)+1))
		values = append(values, *query.MinPrice)
	}

	if query.MaxPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf(`("price" <= $%d)`, len(values)+1))
		values = append(values, *query.MaxPrice)
	}

	if len(whereClauses) > 0 {
		whereQuery := " WHERE " + strings.Join(whereClauses, " AND ")
		baseQuery += whereQuery
		countQuery += whereQuery
	}

	if query.Name != nil {
		baseQuery += ` ORDER BY "name"`
	} else if query.MinPrice != nil || query.MaxPrice != nil {
		baseQuery += ` ORDER BY "price"`
	} else {
		baseQuery += ` ORDER BY "createdAt"`
	}

	if query.Sort != nil && (*query.Sort == "ASC" || *query.Sort == "DESC") {
		baseQuery += ` ` + *query.Sort
	} else {
		baseQuery += ` ASC`
	}

	if query.Page > 0 && query.Limit > 0 {
		limit := query.Limit
		offset := (query.Page - 1) * limit
		baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(values)+1, len(values)+2)
		values = append(values, limit, offset)
	}

	data := models.Products{}
	if err := r.Select(&data, baseQuery, values...); err != nil {
		return nil, 0, err
	}

	var total int
	if err := r.Get(&total, countQuery, values[:len(values)-2]...); err != nil {
		return nil, 0, err
	}

	return &data, total, nil
}

func (r *RepoProduct) GetDetailProduct(uuid string) (*models.Product, error) {
	query := `SELECT * FROM public.products WHERE uuid = :uuid`
	data := models.Product{}

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

func (r *RepoProduct) UpdateProduct(uuid string, data map[string]interface{}) (*models.Product, error) {
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
        UPDATE public.products
        SET %s, "updatedAt" = now()
        WHERE "uuid" = $1
        RETURNING *
    `, strings.Join(setClauses, ", "))

	fmt.Println(query)
	fmt.Println(values)

	var updatedUser models.Product
	err := r.DB.QueryRowx(query, values...).StructScan(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *RepoProduct) DeleteProduct(uuid string) (*models.Product, error) {
	query := `DELETE FROM public.products WHERE uuid = $1 RETURNING *`

	var deletedProduct models.Product
	if err := r.DB.QueryRowx(query, uuid).StructScan(&deletedProduct); err != nil {
		return nil, err
	}

	return &deletedProduct, nil
}
