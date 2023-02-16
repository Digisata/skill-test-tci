package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Digisata/skill-test-tci/app"
	"github.com/Digisata/skill-test-tci/controller"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/swaggo/swag/example/basic/docs"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	schemes := []string{"https"}

	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbUsername := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")

	env := os.Getenv("ENV")
	if env == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		port = os.Getenv("PORT")
		host = fmt.Sprintf("%s:%s", os.Getenv("HOST"), port)
		schemes = []string{"http"}

		dbHost = os.Getenv("MYSQL_HOST")
		dbPort = os.Getenv("MYSQL_PORT")
		dbName = os.Getenv("MYSQL_DATABASE")
		dbUsername = os.Getenv("MYSQL_USERNAME")
		dbPassword = os.Getenv("MYSQL_PASSWORD")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	docs.SwaggerInfo.Version = os.Getenv("DOCS_VERSION")
	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.Schemes = schemes

	validate := validator.New()
	c := controller.NewController(validate, db)
	router := app.NewRouter(c)

	router.Logger.Fatal(router.Start(fmt.Sprintf(":%s", port)))
}
