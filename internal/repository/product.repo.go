package repository

import (
	"fmt"
	"strings"

	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepoInterface interface {
	CreateProduct(data *models.Product) (*models.Product, error)
	GetAllProduct(query *models.ProductQuery) (*models.Products, int, error)
	GetDetailProduct(uuid string) (*models.Product, error)
	UpdateProduct(uuid string, data *models.Product) (*models.Product, error)
	DeleteProduct(uuid string) (*models.Product, error)
}

type RepoProduct struct {
	*sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *RepoProduct {
	return &RepoProduct{db}
}

func (r *RepoProduct) CreateProduct(data *models.Product) (*models.Product, error) {
	query := `
        INSERT INTO public.products (
    			"name",
					"description", 
    			"price", 
					"image",
    			"category"
				) VALUES (
    			:name, 
					:description,
    			:price,
					:image,
    			:category
				)
				RETURNING *;
    		`

	var result models.Product
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

func (r *RepoProduct) UpdateProduct(uuid string, data *models.Product) (*models.Product, error) {
	query := `
		UPDATE public.products
		SET
    	"name" = COALESCE(NULLIF(:name, ''), "name"),
    	"description" = COALESCE(NULLIF(:description, ''), "description"),
    	"category" = COALESCE(NULLIF(:category, ''), "category"),
    	"price" = COALESCE(:price, "price"),
    	"image" = COALESCE(NULLIF(:image, ''), "image"),
    	"updatedAt" = now()
		WHERE "uuid" = :uuid
		RETURNING *;
		`

	data.Uuid = uuid

	var updatedProduct models.Product
	rows, err := r.DB.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&updatedProduct)
		if err != nil {
			return nil, err
		}
	}

	return &updatedProduct, nil
}

func (r *RepoProduct) DeleteProduct(uuid string) (*models.Product, error) {
	query := `DELETE FROM public.products WHERE uuid = $1 RETURNING *`

	var deletedProduct models.Product
	if err := r.DB.QueryRowx(query, uuid).StructScan(&deletedProduct); err != nil {
		return nil, err
	}

	return &deletedProduct, nil
}
