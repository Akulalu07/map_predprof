package main

import (
	"fmt"
	"os"
	"server/database"
	hd "server/handlers"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	user     string
	dbname   string
	password string
	host     string
)

func main() {
	time.Sleep(10 * time.Second)
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	user = os.Getenv("DB_USER")
	dbname = os.Getenv("DB_NAME")
	password = os.Getenv("DB_PASSWORD")
	host = os.Getenv("DB_HOST")

	conn := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable host=%s", user, dbname, password, host)
	//fmt.Println("Andrey Rabotay")
	db, err := database.Init(conn)
	if err != nil {
		sugar.Fatal("Error with connect to database:", err)
	}
	err = db.Ping()
	if err != nil {
		sugar.Fatal("Failed to ping the database:", err)
	}

	r := echo.New()

	r.GET("/api/hello", hd.Hello)
	r.GET("/api/ping", hd.Ping)
	r.Logger.Fatal(r.Start(":8080"))
}
