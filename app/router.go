package app

import (
	"github.com/Digisata/skill-test-tci/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter(c *controller.Controller) *echo.Echo {
	router := echo.New()

	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.GET("/swagger/*", echoSwagger.WrapHandler)
	router.GET("/cars", c.GetCarsHandler)
	router.POST("/is-contain-letters", c.ContainLettersHandler)

	football := router.Group("/football")

	football.POST("/recordgame", c.RecordGameHandler)
	football.GET("/leaguestanding", c.AllLeagueStandingsHandler)
	football.GET("/rank", c.ClubStandingsHandler)


	return router
}
