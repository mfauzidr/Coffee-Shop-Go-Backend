package models

import "time"

var schemaProduct = `
CREATE TABLE public.products (
	id serial4 NOT NULL,
	"name" varchar(30) NOT NULL,
	"description" text NOT NULL,
	"category" varchar(30) NULL,
	"price" int4 NOT NULL,
	"image" text NULL,
	"discountPrice" int4 NULL,
	"createdAt" timestamp NULL DEFAULT now(),
	"updatedAt" timestamp NULL,
	"uuid" uuid NULL DEFAULT uuid_generate_v4(),
	CONSTRAINT products_pkey PRIMARY KEY (id),
	CONSTRAINT unique_name UNIQUE (name)
);
`

func init() {
	_ = schemaProduct // Menghindari peringatan U1000
}

type Product struct {
	Id          int        `db:"id" json:"id"`
	Uuid        string     `db:"uuid" json:"uuid"`
	Name        string     `db:"name" json:"name" form:"name"`
	Description string     `db:"description" json:"description" form:"description"`
	Price       int        `db:"price" json:"price" form:"price"`
	Category    string     `db:"category" json:"category" form:"category"`
	Image       *string    `db:"image,omitempty" json:"image,omitempty" form:"image,omitempty"`
	CreatedAt   *time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updatedAt" json:"updatedAt,omitempty"`
}

type Products []Product

type ProductQuery struct {
	Page     int     `form:"page"`
	Limit    int     `form:"limit"`
	Name     *string `form:"name"`
	MinPrice *int    `form:"minPrice"`
	MaxPrice *int    `form:"maxPrice"`
	Category *string `form:"category"`
	Sort     *string `form:"sort"`
}
