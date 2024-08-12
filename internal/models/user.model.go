package models

import (
	"time"
)

var schemaUsers = `
CREATE TABLE public.users (
	id serial4 NOT NULL,
	"firstName" varchar(30) NOT NULL,
	"lastName" varchar(30) NULL,
	gender varchar(10) NULL,
	email varchar(30) NOT NULL,
	"password" varchar(100) NOT NULL,
	address text NULL,
	"deliveryAddress" text NULL,
	image text NULL,
	"phoneNumber" varchar(15) NULL,
	"role" varchar(20) NULL,
	"createdAt" timestamp NULL DEFAULT now(),
	"updatedAt" timestamp NULL,
	"uuid" uuid NULL DEFAULT uuid_generate_v4(),
	birthday date NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
`

func init() {
	_ = schemaUsers // Menghindari peringatan U1000
}

type Users struct {
	Id              int        `db:"id" json:"id,omitempty" valid:"-"`
	UsersUuid       string     `db:"uuid" json:"uuid" valid:"-"`
	FirstName       string     `db:"firstName" json:"firstName" form:"firstName" valid:"stringlength(2|256)~First Name minimal 2 karakter"`
	LastName        *string    `db:"lastName" json:"lastName" form:"lastName" valid:"stringlength(2|256)~Last Name minimal 2 karakter"`
	Gender          *string    `db:"gender" json:"gender" form:"gender" valid:"type(string)"`
	Email           string     `db:"email" json:"email" form:"email" valid:"email"`
	Password        string     `db:"password" json:"password,omitempty" form:"password" valid:"stringlength(6|256)~Password minimal 6 karakter"`
	Image           *string    `db:"image" json:"image" form:"image" valid:"-"`
	Address         *string    `db:"address" json:"address" form:"address" valid:"-"`
	PhoneNumber     *string    `db:"phoneNumber" json:"phoneNumber" form:"phoneNumber" valid:"numeric,optional"`
	Birthday        *string    `db:"birthday" json:"birthday" form:"birthday" valid:"-"`
	DeliveryAddress *string    `db:"deliveryAddress" json:"deliveryAddress" form:"deliveryAddress" valid:"-"`
	Role            string     `db:"role" json:"role" form:"role" valid:"in(customer|admin|staff)"`
	CreatedAt       *time.Time `db:"createdAt" json:"createdAt" valid:"-"`
	UpdatedAt       *time.Time `db:"updatedAt" json:"updatedAt,omitempty" valid:"-"`
}

type UsersRes []Users

type UsersQuery struct {
	Page   int     `form:"page"`
	Limit  int     `form:"limit"`
	Search *string `form:"search"`
	Sort   *string `form:"sort"`
}
