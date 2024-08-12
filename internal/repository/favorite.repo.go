package repository

import (
	"database/sql"
	"fmt"

	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type FavoriteRepoInterface interface {
	CreateFavorite(data *models.PostFavorite) (*models.PostFavorite, error)
	GetAllFavorite() (*models.Favorites, error)
	GetDetailFavorite(id int) (*models.Favorite, error)
	UpdateFavorite(id int, data *models.UpdateFavorite) (*models.UpdateFavorite, error)
	DeleteFavorite(id int) error
}

type FavoriteRepo struct {
	*sqlx.DB
}

func NewFavoriteRepository(db *sqlx.DB) *FavoriteRepo {
	return &FavoriteRepo{db}
}

func (r *FavoriteRepo) CreateFavorite(data *models.PostFavorite) (*models.PostFavorite, error) {
	query := `INSERT INTO public.favorite("userId", "productId") VALUES(:userId, :productId) RETURNING *`

	var result models.PostFavorite
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

func (r *FavoriteRepo) GetAllFavorite() (*models.Favorites, error) {
	query := `SELECT 
							"f"."id", 
							"u"."firstName" || ' ' || COALESCE("u"."lastName", '') AS "userName",
							"p"."name" AS "productName", 
							"p"."price", 
							"p"."category",
							"f"."createdAt", 
							"f"."updatedAt" FROM public.favorite "f"
						JOIN public.users "u" ON "f"."userId" = "u"."id"
						JOIN public.products "p" ON "f"."productId" = "p"."id"
						ORDER BY "f"."createdAt" DESC`
	data := models.Favorites{}

	if err := r.Select(&data, query); err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *FavoriteRepo) GetDetailFavorite(id int) (*models.Favorite, error) {
	query := `SELECT 
							"f"."id", 
							"u"."firstName" || ' ' || COALESCE("u"."lastName", '') AS "userName",
							"p"."name" AS "productName", 
							"p"."price", 
							"p"."category", 
							"p"."description", 
							"f"."createdAt", 
							"f"."updatedAt" FROM public.favorite "f"
						JOIN public.users "u" ON "f"."userId" = "u"."id"
						JOIN public.products "p" ON "f"."productId" = "p"."id"
						WHERE "f"."id" = :id`
	data := models.Favorite{}

	rows, err := r.DB.NamedQuery(query, map[string]interface{}{
		"id": id,
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

func (r *FavoriteRepo) UpdateFavorite(id int, data *models.UpdateFavorite) (*models.UpdateFavorite, error) {
	query := `
		UPDATE public.favorite 
		SET "productId" = :productId, "updatedAt" = now() 
		WHERE id = :id
		RETURNING *
	`
	rows, err := r.DB.NamedQuery(query, map[string]interface{}{
		"productId": data.ProductId,
		"id":        data.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var favorite models.UpdateFavorite
		err := rows.StructScan(&favorite)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		return &favorite, nil
	}

	return nil, sql.ErrNoRows
}

func (r *FavoriteRepo) DeleteFavorite(id int) error {
	query := `DELETE FROM public.favorite WHERE id = :id`

	_, err := r.DB.NamedExec(query, map[string]interface{}{
		"id": id,
	})
	return err
}
