package handler

import (
	"github.com/gerywahyu/investpedia/merauke/model"
	"github.com/jinzhu/gorm"
)

type ProductHandler struct {
	Conn *gorm.DB
}

func (p *ProductHandler) GetById(productId int) *model.Product {
	var product model.Product
	p.Conn.Where("id = ?", productId).Find(&product)

	return &product
}

func (p *ProductHandler) Create(name string, image string, price int64) *model.Product {
	product := model.Product{
		Name: name,
		ImageURL: image,
		Price: price,
	}
	p.Conn.Create(&product)

	return &product
}

func (p *ProductHandler) ShowAll() *[]model.Product {
	var products []model.Product
	p.Conn.Find(&products)

	return &products
}