package models

import (
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

type Product struct {
	Id          int     `gorm:"auto_increment;primary_key;" json:"id"`
	ProductId   string  `gorm:"size:20;not null;index:idx_products;" json:"product_id"`
	ProductName string  `gorm:"size:100;not null;index:idx_products;" json:"product_name"`
	Hna         float64 `gorm:"not null;" json:"hna"`
	Unit        string  `gorm:"size:20;not null;" json:"unit"`
}

func (p *Product) Prepare() {
	p.ProductId = html.EscapeString(strings.TrimSpace(p.ProductId))
	p.ProductName = html.EscapeString(strings.TrimSpace(p.ProductName))
	p.Unit = html.EscapeString(strings.TrimSpace(p.Unit))
}

func (p *Product) Validate() error {
	if p.ProductId == "" {
		return errors.New("required product id")
	}
	if p.ProductName == "" {
		return errors.New("required product name")
	}
	if p.Hna == 0 {
		return errors.New("required hna")
	}
	if p.Unit == "" {
		return errors.New("required unit")
	}
	return nil
}

func (p *Product) SaveProduct(db *gorm.DB) CrudResult {
	err := db.Debug().Create(&p).Error
	if err != nil {
		return CrudResult{
			Status: "0",
			Note:   err,
		}
	}

	return CrudResult{
		Status: "1",
		Note:   nil,
	}
}

func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	var products []Product

	err := db.Debug().Model(&Product{}).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}

	return &products, nil
}

func (p *Product) FindProductByID(db *gorm.DB, productid string) (*Product, error) {
	err := db.Debug().Model(&Product{}).Where(&Product{ProductId: productid}).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}

	return p, err
}

func (p *Product) UpdateAProduct(db *gorm.DB, productid string) CrudResult {
	db = db.Debug().Where(&Product{ProductId: productid}).Take(&Product{}).Updates(
		Product{
			ProductName: p.ProductName,
			Hna:         p.Hna,
			Unit:        p.Unit,
		},
	)
	if db.Error != nil {
		return CrudResult{
			Status: "0",
			Note:   db.Error,
		}
	}

	return CrudResult{
		Status: "1",
		Note:   nil,
	}
}

func (p *Product) DeleteAProduct(db *gorm.DB) CrudResult {
	db = db.Debug().Model(&Product{}).Where(&Product{ProductId: p.ProductId}).Take(&Product{}).Delete(&Product{})
	if db.Error != nil {
		return CrudResult{
			Status: "0",
			Note:   db.Error,
		}
	}

	return CrudResult{
		Status: "1",
		Note:   nil,
	}
}
