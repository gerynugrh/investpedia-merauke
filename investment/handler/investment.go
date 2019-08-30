package handler

import (
	"github.com/gerywahyu/investpedia/merauke/model"
	"github.com/jinzhu/gorm"
)

type InvestmentHandler struct {
	Conn *gorm.DB
}

func (i *InvestmentHandler) Create(name string, goal int64, year int, current int64, product *model.Product) (*model.Investment, error) {
	investment := model.Investment{
		Name: name,
		Goal: goal,
		Year: year,
		Current: current,
		Product: product,
	}
	i.Conn.Create(&investment)
	return &investment, nil
}

func (i *InvestmentHandler) AddPerson(investment *model.Investment, username string) {
	var user model.User
	i.Conn.Where("username = ?", username).First(&user)
	i.Conn.Model(&user).Association("Investments").Append(investment)
}

func (i *InvestmentHandler) GetById(id int) *model.Investment {
	var investment model.Investment
	i.Conn.Where("id = ?", id).First(&investment)

	return &investment
}

func (i *InvestmentHandler) GetProductById(id int) *model.Product {
	var product model.Product
	i.Conn.Where("id = ?", id).First(&product)

	return &product
}
