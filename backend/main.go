package main

import (
	"os"
	"server/database"
	hd "server/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var (
	user     string
	dbname   string
	password string
	host     string
)

func main() {
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

	_, err = database.Init("./database/database.db")
	if err != nil {
		sugar.Fatal("Error with connect to database:", err)
	}

	if err != nil {
		sugar.Fatal("Failed to ping the database:", err)
	}

	r := echo.New()
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
    }))

	r.POST("/api/map", hd.Map)            //return map
	r.POST("/api/coords", hd.Coords) //
	r.Logger.Fatal(r.Start(":8080"))
}
