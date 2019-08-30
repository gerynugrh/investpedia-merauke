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

func (i *InvestmentHandler) AddFund(id int, username string, income int64) {
	investment := i.GetById(id)
	var user model.User
	i.Conn.Where("line_id = ?", username).First(&user)
	investment.Current += income
	i.Conn.Save(investment)
}

func (i *InvestmentHandler) AddPerson(investment *model.Investment, username string) {
	var user model.User
	i.Conn.Where("line_id = ?", username).First(&user)
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
