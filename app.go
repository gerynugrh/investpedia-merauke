package main

import (
	"github.com/gerywahyu/investpedia/merauke/model"
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

func main() {
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

	port := os.Getenv("PORT")
	err = e.Start(":"+port)
	if err != nil {
		log.Println(err)
		return
	}
}

func callback(c echo.Context) error{
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
				if _,err = bot.ReplyMessage(event.ReplyToken,linebot.NewTextMessage(message.Text)).Do();err!=nil{
					log.Println(err)
				}
			}
		}
	}
	return nil
}
