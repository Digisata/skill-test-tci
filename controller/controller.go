package controller

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/Digisata/skill-test-tci/helper"
	"github.com/Digisata/skill-test-tci/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	Validate *validator.Validate
	DB       *sql.DB
}

func NewController(validate *validator.Validate, db *sql.DB) *Controller {
	return &Controller{
		Validate: validate,
		DB:       db,
	}
}

// RecordGameHandler godoc
// @ID record-game-handler
// @Accept json
// @Produce json
// @Param RecordGameRequest body model.RecordGameRequest true "Record Game"
// @Success 200 {object} model.SuccessResult
// @Failure 400 {object} model.BadRequestResult
// @Failure 500 {object} model.InternalServerErrorResult
// @Router /football/recordgame [post]
func (ctr *Controller) RecordGameHandler(c echo.Context) error {
	requestBody := &[]model.RecordGameRequest{}
	if err := c.Bind(requestBody); err != nil {
		return helper.FailResponse(c, http.StatusBadRequest)
	}

	for _, game := range *requestBody {
		err := ctr.Validate.Struct(game)
		if err != nil {
			return helper.FailResponse(c, http.StatusBadRequest)
		}
	}

	for _, game := range *requestBody {
		score := strings.Split(game.Score, ":")
		scoreHome, scoreAway := strings.TrimSpace(score[0]), strings.TrimSpace(score[1])
		scoreHomeInt, err := strconv.Atoi(scoreHome)
		if err != nil {
			return helper.FailResponse(c, http.StatusInternalServerError)
		}

		scoreAwayInt, err := strconv.Atoi(scoreAway)
		if err != nil {
			return helper.FailResponse(c, http.StatusInternalServerError)
		}

		pointsHome, pointsAway := 0, 0
		if scoreHomeInt > scoreAwayInt {
			pointsHome = 3
		} else if scoreHomeInt == scoreAwayInt {
			pointsHome = 1
			pointsAway = 1
		} else {
			pointsAway = 3
		}

		_, err = ctr.DB.Exec("UPDATE football SET points = points + ? WHERE clubname = ?", pointsHome, game.ClubHomeName)
		if err != nil {
			return helper.FailResponse(c, http.StatusInternalServerError)
		}

		_, err = ctr.DB.Exec("UPDATE football SET points = points + ? WHERE clubname = ?", pointsAway, game.ClubAwayName)
		if err != nil {
			return helper.FailResponse(c, http.StatusInternalServerError)
		}
	}

	return helper.SuccessResponse(c, "Success")

}

// AllLeagueStandingsHandler godoc
// @ID all-league-standings-handler
// @Produce json
// @Success 200 {object} model.SuccessResult
// @Failure 400 {object} model.BadRequestResult
// @Failure 500 {object} model.InternalServerErrorResult
// @Router /football/leaguestanding [get]
func (ctr *Controller) AllLeagueStandingsHandler(c echo.Context) error {
	rows, err := ctr.DB.Query("SELECT clubname, points FROM football ORDER BY points DESC")
	if err != nil {
		return helper.FailResponse(c, http.StatusInternalServerError)
	}
	defer rows.Close()

	var clubs []model.Club
	for rows.Next() {
		var club model.Club
		if err := rows.Scan(&club.ClubName, &club.Points); err != nil {
			return helper.FailResponse(c, http.StatusInternalServerError)
		}
		clubs = append(clubs, club)
	}
	if err := rows.Err(); err != nil {
		return helper.FailResponse(c, http.StatusInternalServerError)
	}

	return helper.SuccessResponse(c, clubs)
}

// ClubStandingsHandler godoc
// @ID club-standings-handler
// @Produce json
// @Success 200 {object} model.SuccessResult
// @Failure 400 {object} model.BadRequestResult
// @Failure 404 {object} model.NotFoundResult
// @Failure 500 {object} model.InternalServerErrorResult
// @Router /football/rank [get]
func (ctr *Controller) ClubStandingsHandler(c echo.Context) error {
	clubName := c.QueryParam("clubname")
	if clubName == "" {
		return helper.FailResponse(c, http.StatusBadRequest)
	}

	row := ctr.DB.QueryRow("SELECT id, clubname, points FROM football WHERE clubname=?", clubName)

	var club model.Club
	if err := row.Scan(&club.ID, &club.ClubName, &club.Points); err != nil {
		if err == sql.ErrNoRows {
			return helper.FailResponse(c, http.StatusNotFound)
		}

		return helper.FailResponse(c, http.StatusInternalServerError)
	}

	var standing int
	row = ctr.DB.QueryRow("SELECT COUNT(*)+1 FROM football WHERE points > ?", club.Points)
	if err := row.Scan(&standing); err != nil {
		return helper.FailResponse(c, http.StatusInternalServerError)
	}

	response := []model.Club{{ClubName: club.ClubName, Standing: standing}}

	return helper.SuccessResponse(c, response)
}
