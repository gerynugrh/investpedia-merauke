package main

import (
	"fmt"
	id "github.com/gerywahyu/investpedia/merauke/investment/delivery"
	ih "github.com/gerywahyu/investpedia/merauke/investment/handler"
	"github.com/gerywahyu/investpedia/merauke/model"
	processor "github.com/gerywahyu/investpedia/merauke/processor"
	pd "github.com/gerywahyu/investpedia/merauke/product/delivery"
	ph "github.com/gerywahyu/investpedia/merauke/product/handler"
	ud "github.com/gerywahyu/investpedia/merauke/user/delivery"
	uh "github.com/gerywahyu/investpedia/merauke/user/handler"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
)


var bot *linebot.Client
var dp *processor.DialogFlowProcessor

func createApp(e *echo.Echo, db *gorm.DB) {
	investmentHandler := ih.InvestmentHandler{
		Conn: db,
	}
	userHandler := uh.UserHandler{
		Conn: db,
	}
	productHandler := ph.ProductHandler{
		Conn: db,
	}
	id.NewInvestmentDelivery(e, &investmentHandler)
	ud.NewUserDelivery(e, &userHandler)
	pd.NewProductDelivery(e, &productHandler)
}

func main() {
	var err error
	cs := os.Getenv("ChannelSecret")
	cat := os.Getenv("ChannelAcccessToken")
	if cs == "" || cat == "" {
		cs = "0b4dbc6ae9724c1d43ab598aebccb02a"
		cat = "dgd8Hs4t9FnKB7KN4MTuP4x4R3I6mMfuPePLKzrCRWeOtnmdKwDO7KQIG9WztDC9VcwISk34XBbB2w38aOvB6SJrqqe2tix0QgO1Id7c88FSBXoFaKTGcAPW3Vdigy7OSWeWeLbXF129AZ7sxP5FvwdB04t89/1O/w1cDnyilFU="
	}

	dp = new(processor.DialogFlowProcessor)
	err = dp.Init("investpedia-chjbgd","en","Asia/Bangkok")
	bot, err = linebot.New(cs, cat)
	var connStr string
	if os.Getenv("GO_ENV") == "heroku" {
		connStr = os.Getenv("DATABASE_URL")
	} else {
		connStr = "host=localhost port=5432 user=merauke dbname=merauke password=merauke sslmode=disable"
	}
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Print(err)
		return
	}

	defer db.Close()
	model.Migrate(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/callback",callback)

	createApp(e, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	err = e.Start(":" + port)
	if err != nil {
		log.Println(err)
		return
	}
}


func callback(c echo.Context) error	{

	jsonListInvestasi :=  `
		{
		  "type": "carousel",
			  "contents": [
				%s
			  ]
			}
		`

	jsonInvestasiTemplate := `{
      "type": "bubble",
      "size": "mega",
      "header": {
        "type": "box",
        "layout": "vertical",
        "contents": [
          {
            "type": "text",
            "text": "%s",
            "color": "#3d3d3d",
            "align": "start",
            "size": "xl",
            "gravity": "center",
            "wrap": true,
            "weight": "bold"
          }
        ],
        "backgroundColor": "#42b549",
        "paddingTop": "19px",
        "paddingAll": "12px",
        "paddingBottom": "16px"
      },
      "body": {
        "type": "box",
        "layout": "vertical",
        "contents": [
          {
            "type": "text",
            "text": "Rp. %s,-",
            "color": "#3d3d3d",
            "align": "start",
            "size": "lg",
            "gravity": "center",
            "weight": "bold"
          },
          {
            "type": "text",
            "text": "dari Rp. %s,-",
            "color": "#8b8b8b",
            "align": "start",
            "size": "xs",
            "gravity": "center"
          },
          {
            "type": "box",
            "layout": "vertical",
            "contents": [
              {
                "type": "box",
                "layout": "vertical",
                "contents": [
                  {
                    "type": "filler"
                  }
                ],
                "width": "%s",
                "backgroundColor": "#fa591d",
                "height": "16px"
              }
            ],
            "backgroundColor": "#9FD8E36E",
            "cornerRadius": "4px",
            "height": "16px",
            "margin": "lg"
          },
          {
            "type": "text",
            "text": "Target %s dari sekarang",
            "color": "#fa591d",
            "weight": "bold",
            "align": "start",
            "size": "xs",
            "gravity": "center",
            "margin": "sm"
          },
          {
            "type": "box",
            "layout": "horizontal",
            "contents": [
              {
                "type": "text",
                "text": "Rekomendasi: Rp. %s /bulan",
                "color": "#8C8C8C",
                "size": "xs",
                "wrap": true
              }
            ],
            "flex": 1
          }
        ],
        "spacing": "md",
        "paddingAll": "16px"
      },
      "footer": {
        "type": "box",
        "layout": "vertical",
        "spacing": "sm",
        "contents": [
          {
            "type": "button",
            "style": "primary",
            "color": "#42b549",
            "action": {
              "type": "message",
              "label": "Tabung disini!",
              "text": "%s"
            }
          }
        ]
      }
    }`

	events,err := bot.ParseRequest(c.Request())
	if err != nil{
		if err == linebot.ErrInvalidSignature {
			c.Response().WriteHeader(400)
		} else{
			c.Response().WriteHeader(500)
		}
		return c.JSON(http.StatusBadRequest,err.Error())
	}
	for _,event := range events{
		if event.Type == linebot.EventTypeMessage{
			switch message:=event.Message.(type) {
			case *linebot.TextMessage:
				var id string
				if event.Source.GroupID == "" && event.Source.RoomID == ""{
					id = event.Source.UserID
				} else{
					if event.Source.GroupID != "" {
						id = event.Source.GroupID
					} else {
						id = event.Source.RoomID
					}
				}
				log.Println(id,message.Text)
				resp := dp.ProcessNLP(message.Text,id)
				fmt.Printf("%#v",resp)
				if _,err = bot.ReplyMessage(event.ReplyToken,linebot.NewTextMessage(resp.Fullfillment)).Do();err!=nil{
					log.Println(err)
				}
				if resp.Fullfillment == "*Daftar Progress Investasi*"{

					productWithData := fmt.Sprintf(jsonInvestasiTemplate, "Macbook", "15.000.000", "24.000.000", "10%", "7 bulan", "300.000", "789")

					data := []byte(fmt.Sprintf(jsonListInvestasi, productWithData))

					container, err := linebot.UnmarshalFlexMessageJSON(data)

					if _,err = bot.PushMessage(id,linebot.NewFlexMessage("list investas", container)).Do(); err!=nil{
						log.Println(err)
					}
				}
				if resp.Fullfillment == "*Pilihan Investasi*"{
					productWithData := fmt.Sprintf(jsonInvestasiTemplate, "Macbook", "15.000.000", "24.000.000", "10", "7 bulan", "300.000", "789")
					productWithData += ","
					productWithData += fmt.Sprintf(jsonInvestasiTemplate, "Macbook", "15.000.000", "24.000.000", "10", "7 bulan", "300.000", "789")
					productWithData += ","
					productWithData += fmt.Sprintf(jsonInvestasiTemplate, "Macbook", "15.000.000", "24.000.000", "10", "7 bulan", "300.000", "789")

					listWithData := fmt.Sprintf(jsonListInvestasi, productWithData)
					log.Println(listWithData)
					data := []byte(listWithData)

					container, err := linebot.UnmarshalFlexMessageJSON(data)

					if _,err = bot.PushMessage(id,linebot.NewFlexMessage("list investas", container)).Do(); err!=nil{
						log.Println(err)
					}
				}
				//if _,err = bot.PushMessage(id,linebot.NewFlexMessage("")).Do(); err!=nil{
				//	log.Println(err)
				//}
			}
		}
	}
	return nil
}
