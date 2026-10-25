package schemas

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Customer struct {
	ID          string         `gorm:"primaryKey;column:id" json:"id"`
	Nama        string         `gorm:"column:nama" json:"nama"`
	Email       string         `gorm:"column:email" json:"email"`
	NomerTelpon string         `gorm:"column:nomer_telpon" json:"nomer_telpon"`
	Password    string         `gorm:"column:password" json:"-"`
	Status      string         `gorm:"column:status" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

type Admin struct {
	ID          string         `gorm:"primaryKey;column:id" json:"id"`
	Nama        string         `gorm:"column:nama" json:"nama"`
	Email       string         `gorm:"column:email;unique" json:"email"`
	NomorTelpon string         `gorm:"column:nomor_telpon" json:"nomor_telpon"`
	Password    string         `gorm:"column:password" json:"-"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

// Order status constants
const (
	StatusMenungguPembayaran = "Menunggu Pembayaran"
	StatusSelesai            = "Selesai"
	StatusDibatalkan         = "Dibatalkan"
)

// BlockingStatuses lists statuses that hold room reservation and block overlapping bookings.
var BlockingStatuses = []string{
	StatusMenungguPembayaran,
	StatusSelesai,
}

// IsBlockingStatus checks if a given order status reserves availability.
func IsBlockingStatus(status string) bool {
	for _, s := range BlockingStatuses {
		if s == status {
			return true
		}
	}
	return false
}

type Claims struct {
	IdUser string `json:"id"`
	jwt.RegisteredClaims
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
	RoomID     string `json:"room_id" binding:"required"`
	DateStart  string `json:"date_start" binding:"required"`
	DateFinish string `json:"date_finish" binding:"required"`
}

type Room struct {
	ID          string    `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Capacity    int       `gorm:"column:capacity" json:"capacity"`
	Price       float64   `gorm:"column:price" json:"price"`
	ImageURL    string    `gorm:"column:image_url" json:"image_url"`
	Features    string    `gorm:"column:features" json:"features"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Order struct {
	ID         string    `gorm:"primaryKey;column:id" json:"id"`
	UserID     string    `gorm:"column:user_id;not null" json:"user_id"`
	RoomID     string    `gorm:"column:room_id;not null" json:"room_id"`
	DateStart  time.Time `gorm:"column:date_start;type:timestamptz;not null" json:"date_start"`
	DateFinish time.Time `gorm:"column:date_finish;type:timestamptz;not null" json:"date_finish"`
	Status     string    `gorm:"column:status;not null" json:"status"`
	TotalPrice float64   `gorm:"column:total_price;not null" json:"total_price"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type LoginAdminSchema struct {
	Nama     string `form:"nama" json:"nama" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

// AvailabilitySlot represents a privacy-safe occupied schedule interval for a room.
// Customer identities, user IDs, prices, and internal payment data are strictly omitted.
type AvailabilitySlot struct {
	RoomID     string    `gorm:"column:room_id" json:"room_id"`
	DateStart  time.Time `gorm:"column:date_start" json:"date_start"`
	DateFinish time.Time `gorm:"column:date_finish" json:"date_finish"`
	Status     string    `gorm:"column:status" json:"status"`
}
