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
