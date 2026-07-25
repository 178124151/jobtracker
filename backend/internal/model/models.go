package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	DisplayName  string         `gorm:"type:varchar(100)" json:"display_name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Company struct {
	ID            string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name          string         `gorm:"type:varchar(200);not null" json:"name"`
	NameEN        string         `gorm:"type:varchar(200)" json:"name_en"`
	LogoURL       string         `gorm:"type:varchar(512)" json:"logo_url"`
	Website       string         `gorm:"type:varchar(512);not null" json:"website"`
	Industry      string         `gorm:"type:varchar(100)" json:"industry"`
	Group         string         `gorm:"type:varchar(50)" json:"group"`
	Description   string         `gorm:"type:text" json:"description"`
	IsPreset      bool           `gorm:"default:true" json:"is_preset"`
	HealthStatus  string         `gorm:"type:varchar(10);default:UNKNOWN" json:"health_status"`
	LastChecked   *time.Time     `json:"last_checked"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Application struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      string         `gorm:"type:uuid;not null" json:"user_id"`
	CompanyID   *string        `gorm:"type:uuid" json:"company_id"`
	CompanyName string         `gorm:"type:varchar(200);not null" json:"company_name"`
	Position    string         `gorm:"type:varchar(200);not null" json:"position"`
	AppliedAt   time.Time      `gorm:"type:date;not null" json:"applied_at"`
	Status      string         `gorm:"type:varchar(20);default:RESUME" json:"status"`
	Notes       string         `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Resume struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"type:uuid" json:"user_id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Template  string         `gorm:"type:varchar(50);default:classic" json:"template"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (c *Company) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (a *Application) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

func (r *Resume) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}