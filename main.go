package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Max-Jordan/ArticlesBot/bot"
	"github.com/Max-Jordan/ArticlesBot/handler"
	"github.com/Max-Jordan/ArticlesBot/services"
	"github.com/lib/pq"
)

const (
	updateLimits = 100
)

func main() {
	b := bot.NewBot(os.Getenv("bot_token"), updateLimits)

	cnf := bot.UpdateConfig{
		Offset:  0,
		Limit:   updateLimits,
		Timeout: 60,
	}
	dbConf := pq.Config{
		Host:     "localhost",
		Port:     5432,
		Database: "article",
		User:     "admin",
		Password: "01226457",
		SSLMode:  "disable",
	}

	c, err := pq.NewConnectorConfig(dbConf)
	if err != nil {
		log.Fatal(err)
	}

	db := sql.OpenDB(c)
	defer db.Close()

	b.StartPolling(cnf)

	storage := services.NewStorage(db)

	handler := handler.NewHandler(b, storage)
	handler.StartHandling()
}
