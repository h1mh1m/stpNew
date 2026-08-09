package schemas

import (
	"gorm.io/gorm"
)

type customer struct {
	gorm.Model
	id          string ``
	nama        string ``
	email       string ``
	nomerTelpon string ``
	password    string ``
}
type admin struct {
	gorm.Model
	id          string ``
	nama        string ``
	email       string ``
	nomorTelpon string ``
	password    string ``
}

type LoginSchema struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type SignupSchema struct {
	Nama        string `json:"nama"`
	Email       string `json:"email"`
	NomorTelpon string `json:"nomor_telpon"`
	Password    string `json:"password"`
}
type BookingSchema struct {
	dateStart  string ``
	dateFinish string ``
}
