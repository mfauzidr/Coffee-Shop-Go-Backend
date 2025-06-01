package models

import (
	"time"
)

var schemaUsers = `
CREATE TABLE public.users (
	"id" serial4 NOT NULL,
	"fullName" varchar(30) NOT NULL,
	"email" varchar(30) NOT NULL,
	"password" varchar(100) NOT NULL,
	"address" text NULL,
	"image" text NULL,
	"phoneNumber" varchar(15) NULL,
	"role" varchar(20) NULL,
	"createdAt" timestamp DEFAULT now() NULL,
	"updatedAt" timestamp NULL,
	"uuid" uuid DEFAULT uuid_generate_v4() NULL,
	CONSTRAINT user_uuid_unique UNIQUE (uuid),
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
`

func init() {
	_ = schemaUsers // Menghindari peringatan U1000
}

type Users struct {
	Id          int        `db:"id" json:"id,omitempty" valid:"-"`
	UUID        string     `db:"uuid" json:"uuid" valid:"-"`
	Fullname    *string    `db:"fullName" json:"fullName" form:"fullName" valid:"stringlength(2|256)~First Name minimal 2 karakter"`
	Email       string     `db:"email" json:"email" form:"email" valid:"email"`
	Password    string     `db:"password" json:"password,omitempty" form:"password" valid:"stringlength(6|256)~Password minimal 6 karakter"`
	Image       *string    `db:"image" json:"image" valid:"-"`
	Address     *string    `db:"address" json:"address" form:"address" valid:"-"`
	PhoneNumber *string    `db:"phoneNumber" json:"phoneNumber" form:"phoneNumber" valid:"numeric,optional"`
	Role        string     `db:"role" json:"role" form:"role" valid:"in(customer|admin|staff)"`
	CreatedAt   *time.Time `db:"createdAt" json:"createdAt" valid:"-"`
	UpdatedAt   *time.Time `db:"updatedAt" json:"updatedAt,omitempty" valid:"-"`
}

type UsersRes []Users

type UsersQuery struct {
	Page   int     `form:"page"`
	Limit  int     `form:"limit"`
	Search *string `form:"search"`
	Sort   *string `form:"sort"`
}
