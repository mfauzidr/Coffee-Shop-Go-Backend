package models

import "time"

var schemaFavorite = `
CREATE TABLE public.favorite (
	id serial4 NOT NULL,
	"userId" int4 NULL,
	"productId" int4 NULL,
	"createdAt" timestamp NULL DEFAULT now(),
	"updatedAt" timestamp NULL,
	CONSTRAINT favorite_pkey PRIMARY KEY (id),
	CONSTRAINT "favorite_productId_fkey" FOREIGN KEY ("productId") REFERENCES public.products(id) ON DELETE SET NULL,
	CONSTRAINT "favorite_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) ON DELETE SET NULL
);
`

func init() {
	_ = schemaFavorite // Menghindari peringatan U1000
}

type Favorite struct {
	Id          int        `db:"id" json:"id"`
	UserName    *string    `db:"userName" json:"userName,omitempty"`
	ProductName string     `db:"productName" json:"productName"`
	Price       int        `db:"price" json:"price"`
	Category    string     `db:"category" json:"category"`
	Description *string    `db:"description" json:"description,omitempty"`
	CreatedAt   *time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updatedAt" json:"updatedAt"`
}

type PostFavorite struct {
	UserId    int        `db:"userId" json:"userId" form:"userId"`
	ProductId int        `db:"productId" json:"productId" form:"productId"`
	CreatedAt *time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt *time.Time `db:"updatedAt" json:"updatedAt"`
}

type Favorites []Favorite

type UpdateFavorite struct {
	Id        int        `db:"id" json:"id"`
	UserId    *int       `db:"userId" json:"userId"`
	ProductId int        `db:"productId" json:"productId" form:"productId"`
	CreatedAt *time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt *time.Time `db:"updatedAt" json:"updatedAt"`
}
