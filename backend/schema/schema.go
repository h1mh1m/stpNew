package schemas

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Customer struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Nama        string         `gorm:"column:nama" json:"nama"`
	Email       string         `gorm:"column:email;unique" json:"email"`
	NomerTelpon string         `gorm:"column:nomer_telpon" json:"nomer_telpon"`
	Password    string         `gorm:"column:password" json:"-"` // Disembunyikan saat di-parse ke JSON
	Status      string         `gorm:"column:status;default:ACTIVE" json:"status"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz;index;column:deleted_at" json:"-"` // Disembunyikan
}

// Admin struct merepresentasikan tabel admins
type Admin struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Nama        string         `gorm:"column:nama" json:"nama"`
	Email       string         `gorm:"column:email;unique" json:"email"`
	NomorTelpon string         `gorm:"column:nomor_telpon" json:"nomor_telpon"`
	Password    string         `gorm:"column:password" json:"-"` // Disembunyikan
	CreatedAt   time.Time      `gorm:"type:timestamptz;column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz;index;column:deleted_at" json:"-"` // Disembunyikan
}

// Booking struct merepresentasikan tabel bookings
type Booking struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CustomerID  string `gorm:"type:uuid;column:customer_id" json:"customer_id"`
	TipeBooking string `gorm:"column:tipe_booking;not null;default:'HARIAN'" json:"tipe_booking"`

	// Gunakan pointer (*string dan *int) agar bisa merepresentasikan nilai NULL di database
	Sesi         *string   `gorm:"column:sesi" json:"sesi,omitempty"`
	JumlahHari   *int      `gorm:"column:jumlah_hari" json:"jumlah_hari,omitempty"`
	WaktuMulai   time.Time `gorm:"type:timestamptz;column:waktu_mulai" json:"waktu_mulai"`
	WaktuSelesai time.Time `gorm:"type:timestamptz;column:waktu_selesai" json:"waktu_selesai"`

	Status    string         `gorm:"column:status;default:PENDING" json:"status"`
	CreatedAt time.Time      `gorm:"type:timestamptz;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index;column:deleted_at" json:"-"`

	// Relasi Eager Loading (Opsional saat query .Preload())
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

// Payment struct merepresentasikan tabel payments
type Payment struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BookingID string `gorm:"type:uuid;column:booking_id" json:"booking_id"`

	BankName string          `gorm:"column:bank_name" json:"bank_name"`
	VANumber string          `gorm:"column:va_number" json:"va_number"`
	Amount   decimal.Decimal `gorm:"type:numeric(15,2);column:amount" json:"amount"` // Mencegah presisi bug finansial

	Status    string         `gorm:"column:status;default:UNPAID" json:"status"`
	CreatedAt time.Time      `gorm:"type:timestamptz;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index;column:deleted_at" json:"-"`

	// Relasi
	Booking Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
}

type Claims struct {
	IdUser string `json:id`
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
type BookingInput struct {
	TanggalAwal  string `json:"tanggal_awal" binding:"required"`  // Format: YYYY-MM-DD
	TanggalAkhir string `json:"tanggal_akhir" binding:"required"` // Format: YYYY-MM-DD
	TipeBooking  string `json:"tipe_booking" binding:"required"`  // "SESI" atau "HARIAN"
	Sesi         *int   `json:"sesi"`                             // 1, 2, atau 3 (Pointer karena opsional)
}

type LoginAdminSchema struct {
	Nama     string `form:"nama" json:"nama" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required, email"`
	Password string `form:"password" json:"password" binding:"required"`
}
