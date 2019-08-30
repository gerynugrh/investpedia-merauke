package delivery

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gerywahyu/investpedia/merauke/model"
	"github.com/gerywahyu/investpedia/merauke/user/handler"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"time"
)

type UserDelivery struct {
	Handler *handler.UserHandler
}

func NewUserDelivery(e *echo.Echo, handler *handler.UserHandler) {
	delivery := &UserDelivery{
		Handler:handler,
	}
	e.POST("/login", delivery.Login)
	e.POST("/register", delivery.Register)
	e.POST("/user", delivery.ShowInvestment)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (u *UserDelivery) Login(c echo.Context) error {
	var loginRequest LoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return err
	}

	success, err := u.Handler.Login(loginRequest.Username, loginRequest.Password)
	if !success {
		return err
	}

	claims := &model.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: "investpedia",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
		Username: loginRequest.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}
	loginResponse := LoginResponse{Token: tokenString}

	return c.JSON(http.StatusOK, loginResponse)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	LineId	 string `json:"line_id"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

func (u *UserDelivery) Register(c echo.Context) error {
	var request RegisterRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	user, err := u.Handler.Register(request.Username, request.Password, request.LineId)
	if user != nil {
		return err
	}

	response := RegisterResponse{Success: true}
	return c.JSON(http.StatusOK, response)
}

type ShowInvestmentRequest struct {
	Username 	string `json:"username"`
}

type ShowInvestmentResponse struct {
	Investments	[]model.Investment `json:"payload"`
}

func (u *UserDelivery) ShowInvestment(c echo.Context) error {
	var req ShowInvestmentRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	investments := u.Handler.ShowInvestment(req.Username)
	res := ShowInvestmentResponse{Investments: investments}

	return c.JSON(http.StatusOK, res)
}