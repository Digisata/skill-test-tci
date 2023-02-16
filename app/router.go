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

	router.POST("/football/recordgame", c.RecordGameHandler)
	router.GET("/football/leaguestanding", c.AllLeagueStandingsHandler)
	router.GET("/football/rank", c.ClubStandingsHandler)

	return router
}
