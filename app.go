package main

import (
	"fmt"
	"github.com/gerywahyu/investpedia/merauke/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strconv"
)


var bot *linebot.Client

func main() {
	var cs string
	if os.Getenv("GO_ENV") == "heroku" {
		cs = os.Getenv("DATABASE_URL")
	} else {
		cs = "host=localhost port=5432 user=merauke dbname=merauke password=merauke sslmode=disable"
	}
	db, err := gorm.Open("postgres", cs)
	if err != nil {
		log.Print(err)
		return
	}

	defer db.Close()
	model.Migrate(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	err = e.Start(":8123")
	if err != nil {
		return
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var id string
				if event.Source.GroupID == "" && event.Source.RoomID == "" {
					id = event.Source.UserID
				} else {
					if event.Source.GroupID != "" {
						id = event.Source.GroupID
					} else {
						id = event.Source.RoomID
					}
				}
				fmt.Print(id)
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
