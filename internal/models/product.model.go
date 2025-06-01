package models

import "time"

var schemaProduct = `
CREATE TABLE public.products (
	"id" serial4 NOT NULL,
	"name" varchar(30) NOT NULL,
	"description" text NOT NULL,
	"price" int4 NOT NULL,
	"image" text NULL,
	"discountPrice" int4 NULL,
	"isRecommended" bool NULL,
	"createdAt" timestamp DEFAULT now() NULL,
	"updatedAt" timestamp NULL,
	"uuid" uuid DEFAULT uuid_generate_v4() NULL,
	CONSTRAINT products_pkey PRIMARY KEY (id),
	CONSTRAINT products_uuid_unique UNIQUE (uuid),
	CONSTRAINT unique_name UNIQUE (name)
);
`

func init() {
	_ = schemaProduct // Menghindari peringatan U1000
}

type Product struct {
	Id            int        `db:"id" json:"id" valid:"-"`
	Uuid          string     `db:"uuid" json:"uuid" valid:"-"`
	Name          string     `db:"name" json:"name" form:"name" valid:"stringlength(4|256)~Product Name minimal 4 karakter"`
	Description   string     `db:"description" json:"description" form:"description" valid:"-"`
	Price         int        `db:"price" json:"price" form:"price" valid:"int"`
	DiscountPrice *int       `db:"discountPrice" json:"discountPrice,omitempty" form:"discountPrice" valid:"int"`
	IsRecommended *bool      `db:"isRecommended" json:"isRecommended,omitempty" form:"isRecommended" valid:"int"`
	Image         *string    `db:"image" json:"image" valid:"-"`
	CreatedAt     *time.Time `db:"createdAt" json:"createdAt" valid:"-"`
	UpdatedAt     *time.Time `db:"updatedAt" json:"updatedAt,omitempty" valid:"-"`
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
