package delivery

import (
	"github.com/gerywahyu/investpedia/merauke/model"
	"github.com/gerywahyu/investpedia/merauke/product/handler"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

type ProductDelivery struct {
	Handler *handler.ProductHandler
}

func NewProductDelivery(e *echo.Echo, handler *handler.ProductHandler) {
	delivery := &ProductDelivery{
		Handler: handler,
	}
	e.GET("/product", delivery.Show)
	e.GET("/products", delivery.ShowAll)
	e.POST("/products", delivery.Create)
}

type CreateRequest struct {
	Name 	string
	Image	string
	Price	int64
}

type CreateResponse struct {
	Product model.Product `json:"payload"`
}

func (p *ProductDelivery) Create(e echo.Context) error {
	var request CreateRequest
	err := e.Bind(&request)
	if err != nil {
		log.Println(err)
		return err
	}

	product := p.Handler.Create(request.Name, request.Image, request.Price)

	response := CreateResponse{Product: *product}
	return e.JSON(http.StatusOK, response)
}

type ShowAllResponse struct {
	Products []model.Product `json:"payload"`
}

func (p *ProductDelivery) ShowAll(e echo.Context) error {
	products := p.Handler.ShowAll()
	response := ShowAllResponse{Products: *products}
	return e.JSON(http.StatusOK, response)
}

type ShowResponse struct {
	Product model.Product `json:"payload"`
}

func (p *ProductDelivery) Show(e echo.Context) error {
	id, err := strconv.Atoi(e.QueryParam("id"))
	if err != nil {
		log.Println(err)
		return err
	}
	product := p.Handler.GetById(id)
	response := ShowResponse{Product: *product}
	return e.JSON(http.StatusOK, response)
}