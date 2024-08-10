package models

import "time"

var schemaFavorite = `
create table public.favorite (
	favorite_id serial,
	favorite_uuid uuid unique default gen_random_uuid(),
	product_id int,
	created_at timestamp without time zone default now(),
	updated_at timestamp without time zone,
	constraint favorite_pk primary key(favorite_id),
	constraint product_fk foreign key (product_id) references public.product(product_id) on delete set null
);
`

func init() {
	_ = schemaFavorite // Menghindari peringatan U1000
}

type Favorite struct {
	Favorite_id   int        `db:"favorite_id" json:"favorite_id"`
	Favorite_uuid string     `db:"favorite_uuid" json:"favorite_uuid"`
	Name          string     `db:"name" json:"name"`
	Price         int        `db:"price" json:"price"`
	Category      string     `db:"category" json:"category"`
	Description   string     `db:"description" json:"description"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

type PostFavorite struct {
	ID        int        `db:"product_id" json:"product_id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

type Favorites []Favorite

type UpdateFavorite struct {
	Favorite_id   int        `db:"favorite_id" json:"favorite_id"`
	Favorite_uuid string     `db:"favorite_uuid" json:"favorite_uuid"`
	ID            int        `db:"product_id" json:"product_id"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}
