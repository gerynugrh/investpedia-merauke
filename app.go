package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	"github.com/gerywahyu/investpedia/merauke/internal"
	"github.com/line/line-bot-sdk-go/linebot"
)


var bot *linebot.Client

func initFlags(args *internal.Args) {
	port := flag.Int("port", 3000, "port number for your apps")
	args.Port = *port
}

func initHandler(handler *internal.Handler) error {

	// Initialize SQL DB
	// NOTE: change localhost:5432 to postgres:5432 if you use docker
	db, err := sql.Open("postgres", "postgres://postgres:postgres@172.32.4.208:5432/?sslmode=disable")
	if err != nil {
		return err
	}
	handler.DB = db

	return nil
}

func initRouter(router *httprouter.Router, handler *internal.Handler) {

	router.GET("/", handler.Index)

	// Single user API
	router.GET("/user/:userID", handler.GetUserByID)
	router.POST("/user", handler.InsertUser)
	router.PUT("/user/:userID", handler.EditUserByID)
	router.DELETE("/user/:userID", handler.DeleteUserByID)

	// Single book API
	router.GET("/book/:bookID", handler.GetBookByID)
	router.POST("/book", handler.InsertBook)
	router.PUT("/book/:bookID", handler.EditBook)
	router.DELETE("/book/:bookID", handler.DeleteBookByID)

	// Batch book API
	router.POST("/books", handler.InsertMultipleBooks)

	// Lending API
	router.POST("/lend", handler.LendBook)

	// `httprouter` library uses `ServeHTTP` method for it's 404 pages
	router.NotFound = handler
}

func main() {
	args := new(internal.Args)
	initFlags(args)

	handler := new(internal.Handler)
	if err := initHandler(handler); err != nil {
		panic(err)
	}

	router := httprouter.New()
	initRouter(router, handler)

	var err error

	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)

	//// Line Handler
	http.HandleFunc("/callback", callbackHandler)
	//http.HandleFunc("/", greet)

	fmt.Printf("Apps served on :%d\n", args.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.Port), router))
}


func test(ctx string, msg string) {
	log.Println("Tested : ", ctx, msg)
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
				test(id, message.Text)
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
