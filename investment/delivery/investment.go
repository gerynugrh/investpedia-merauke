package delivery

import (
	"github.com/gerywahyu/investpedia/merauke/investment/handler"
	"github.com/gerywahyu/investpedia/merauke/model"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

type InvestmentDelivery struct {
	Handler *handler.InvestmentHandler
}

func NewInvestmentDelivery(e *echo.Echo, handler *handler.InvestmentHandler) {
	delivery := &InvestmentDelivery{
		Handler: handler,
	}
	e.GET("/investment", delivery.Show)
	e.POST("/investment", delivery.Create)
	e.POST("/add_fund", delivery.AddFund)
}

type ShowResponse struct {
	Investment model.Investment `json:"investment"`
}

func (i *InvestmentDelivery) Show(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		log.Println(err)
		return err
	}
	investment := i.Handler.GetById(id)
	response := ShowResponse{Investment: *investment}

	return c.JSON(http.StatusOK, response)
}

type CreateRequest struct {
	Name		string `json:"name"`
	Goal 		int64 `json:"goal"`
	Year		int `json:"year"`
	Current		int64 `json:"current"`
	ProductId	string `json:"productId"`
	LineId 		string `json:"line_id"`
}



func (i *InvestmentDelivery) Create(c echo.Context) error {
	var request CreateRequest
	err := c.Bind(&request)
	if err != nil {
		log.Println(err)
		return err
	}
	var product *model.Product
	var value int64
	var name string
	name = request.Name
	value = request.Goal
	if request.ProductId != "" {
		id, err := strconv.Atoi(request.ProductId)
		if err != nil {
			log.Println(err)
			return err
		}
		product = i.Handler.GetProductById(id)
		value = product.Price
		name = product.Name
	}
	investment, err := i.Handler.Create(name, value, request.Year, request.Current, product)
	i.Handler.AddPerson(investment, request.LineId)
	response := ShowResponse{Investment: *investment}

	return c.JSON(http.StatusOK, response)
}

type AddFundRequest struct {
	LineId		string `json:"line_id"`
	Amount		int64 `json:"amount"`
}

type AddFundResponse struct {
	Success	bool `json:"success"`
}

func (i *InvestmentDelivery) AddFund(c echo.Context) error {
	var request AddFundRequest
	err := c.Bind(&request)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		log.Println(err)
		return err
	}

	i.Handler.AddFund(id, request.LineId, request.Amount)
	response := AddFundResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}
