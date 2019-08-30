package handler

import (
	"errors"
	"github.com/gerywahyu/investpedia/merauke/model"
	"github.com/gerywahyu/investpedia/merauke/util"
	"github.com/jinzhu/gorm"
)

type UserHandler struct {
	Conn *gorm.DB
}

func (u *UserHandler) Register(username string, password string) (*model.User, error) {
	hash, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := model.User{
		Username: username,
		Password: hash}
	u.Conn.Create(&user)
	return &user, nil
}

func (u *UserHandler) Login(username string, password string) (bool, error) {
	var user model.User
	u.Conn.Where("username = ?", username).First(&user)
	if &user == nil {
		return false, errors.New("can't find user")
	}
	success := util.CheckPasswordHash(password, user.Password)
	if !success {
		return false, errors.New("password doesn't match")
	}
	return success, nil
}

func (u *UserHandler) ShowInvestment(username string) []model.Investment {
	var user model.User
	u.Conn.Where("username = ?", username).First(&user)

	return user.Investments
}

func (u *UserHandler) IsLinked(username string) bool {
	var user model.User
	u.Conn.Where("username = ?", username).First(&user)

	return user.LineId != ""
}
