package models

import "time"

var schemaProduct = `
CREATE TABLE public.products (
	id serial4 NOT NULL,
	"name" varchar(30) NOT NULL,
	description text NOT NULL,
	price int4 NOT NULL,
	image text NULL,
	"discountPrice" int4 NULL,
	"createdAt" timestamp NULL DEFAULT now(),
	"updatedAt" timestamp NULL,
	"uuid" uuid NULL DEFAULT uuid_generate_v4(),
	"category" varchar(30) NULL,
	CONSTRAINT products_pkey PRIMARY KEY (id),
	CONSTRAINT unique_name UNIQUE (name)
);
`

func init() {
	_ = schemaProduct // Menghindari peringatan U1000
}

type Product struct {
	ID          int        `db:"id" json:"id"`
	UUID        string     `db:"uuid" json:"uuid"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	Price       int        `db:"price" json:"price"`
	Category    string     `db:"category" json:"category"`
	Image       string     `db:"image,omitempty"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type Products []Product

type ProductQuery struct {
	Name     string `form:"name"`
	MinPrice int    `form:"minPrice"`
	MaxPrice int    `form:"maxPrice"`
	Category string `form:"category"`
	Sort     string `form:"sort"`
	Page     int    `form:"page"`
}
