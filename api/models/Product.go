package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Product struct for Product
type Product struct {
	ID          uuid.UUID `gorm:"primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Brand       string    `gorm:"size:255;not null" json:"brand"`
	Image       string    `gorm:"size:2000;null" json:"image"`
	Size        string    `gorm:"size:200;null" json:"size"`
	Model       string    `gorm:"size:255;null" json:"model"`
	Price       float64   `gorm:"default:0.00;null" json:"price"`
	OwnerID     uuid.UUID `gorm:"not null" json:"owner_id"`
	ExpDate     time.Time `gorm:"null" json:"exp_date"`
	Description string    `gorm:"size:2000;null" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ResponseProduct return for the struct Product
type ResponseProduct struct {
	ID          uuid.UUID
	Name        string
	Brand       string
	Image       string
	Size        string
	Model       string
	Price       float64
	ExpDate     time.Time
	OwnerID     uuid.UUID
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SanitizeProduct the Product response
func SanitizeProduct(p Product) ResponseProduct {
	return ResponseProduct{
		p.ID,
		p.Name,
		p.Brand,
		p.Image,
		p.Size,
		p.Model,
		p.Price,
		p.ExpDate,
		p.OwnerID,
		p.Description,
		p.CreatedAt,
		p.UpdatedAt,
	}
}

// Prepare set value for Product
func (p *Product) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Brand = html.EscapeString(strings.TrimSpace(p.Brand))
	p.Image = html.EscapeString(strings.TrimSpace(p.Image))
	p.Size = html.EscapeString(strings.TrimSpace(p.Size))
	p.Model = html.EscapeString(strings.TrimSpace(p.Model))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate validations on actions
func (p *Product) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if p.Name == "" {
			return errors.New("Required Name")
		}
		if p.Brand == "" {
			return errors.New("Required Brand")
		}
		if p.Image == "" {
			return errors.New("Required Image")
		}
		if p.Price == 0.00 {
			return errors.New("Required Price")
		}
		return nil

	default:
		if p.Name == "" {
			return errors.New("Required Name")
		}
		if p.Brand == "" {
			return errors.New("Required Brand")
		}
		if p.Image == "" {
			return errors.New("Required Image")
		}
		if p.Price == 0.00 {
			return errors.New("Required Price")
		}
		return nil
	}
}

// SaveProduct save Product
func (p *Product) SaveProduct(db *gorm.DB) (*Product, error) {

	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	return p, nil
}

// FindAllProducts get all Products
func (p *Product) FindAllProducts(db *gorm.DB, oid uuid.UUID) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Where("owner_id = ?", oid).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}

	return &products, err
}

// FindProductByID fin Product by id
func (p *Product) FindProductByID(db *gorm.DB, pid uuid.UUID, oid uuid.UUID) (*Product, error) {
	err := db.Debug().Model(Product{}).Where("id = ?", pid).Where("owner_id = ?", oid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &Product{}, errors.New("Product Not Found")
	}
	return p, err
}

// UpdateAProduct update Product
func (p *Product) UpdateAProduct(db *gorm.DB, pid uuid.UUID) (*Product, error) {

	db = db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&Product{}).UpdateColumns(
		map[string]interface{}{
			"name":        p.Name,
			"brand":       p.Brand,
			"image":       p.Image,
			"size":        p.Size,
			"model":       p.Model,
			"price":       p.Price,
			"exp_date":    p.ExpDate,
			"description": p.Description,
			"updated_at":  time.Now(),
		},
	)
	if db.Error != nil {
		return &Product{}, db.Error
	}
	// This is the display the updated Product
	err := db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}

	return p, nil
}

// DeleteAProduct delete Product by id
func (p *Product) DeleteAProduct(db *gorm.DB, pid uuid.UUID, oid uuid.UUID) (int64, error) {

	db = db.Debug().Model(&Product{}).Where("id = ? and owner_id = ?", pid, oid).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// FindAllOpenProducts get all Products
func (p *Product) FindAllOpenProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}

	return &products, err
}