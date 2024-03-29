package controller

import (
	"database/sql"
	"fmt"
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
// @Param RecordGameRequest body []model.RecordGameRequest true "Record Game"
// @Success 200 {object} model.SuccessResult
// @Failure 400 {object} model.BadRequestResult
// @Failure 404 {object} model.NotFoundResult
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

		stmt := fmt.Sprintf("SELECT * FROM football WHERE clubname IN ('%s','%s')", game.ClubHomeName, game.ClubAwayName)
		rows, err := ctr.DB.Query(stmt)
		if err != nil {
			return helper.FailResponse(c, http.StatusInternalServerError)
		}
		defer rows.Close()

		counter := 0
		for rows.Next() {
			counter++
		}

		if counter != 2 {
			return helper.FailResponse(c, http.StatusNotFound)
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
// @Summary All league standings
// @Description get league standings
// @Accept json
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
// @Summary Club standings
// @Description get club by name
// @Accept json
// @Produce json
// @Param clubname query string true "club search by clubname"
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

// GetCarsHandler godoc
// @ID get-cars-handler
// @Summary Get cars
// @Description get cars with the corresponding type
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessResult
// @Failure 500 {object} model.InternalServerErrorResult
// @Router /cars [get]
func (ctr *Controller) GetCarsHandler(c echo.Context) error {
	rows, err := ctr.DB.Query(`
		SELECT
		brand,
		MAX(CASE WHEN row_num = 1 THEN CONCAT(type, ' : ', price) END) AS TYPE1,
		MAX(CASE WHEN row_num = 2 THEN CONCAT(type, ' : ', price) END) AS TYPE2,
		MAX(CASE WHEN row_num = 3 THEN CONCAT(type, ' : ', price) END) AS TYPE3
		FROM (
		SELECT 
			brand,
			type,
			price,
			@row_num := IF(@current_brand = brand, @row_num + 1, 1) AS row_num,
			@current_brand := brand
		FROM 
			cars
			CROSS JOIN (SELECT @current_brand := '', @row_num := 0) AS vars
		ORDER BY 
			brand, price DESC
		) AS ranked_cars
		WHERE row_num <= 3
		GROUP BY brand
		ORDER BY brand;
	`)

	if err != nil {
		return helper.FailResponse(c, http.StatusInternalServerError)
	}

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	cars := []map[string]interface{}{}

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		car := make(map[string]interface{})
		
		for i, col := range columns {
			val := values[i]

			b, ok := val.([]byte)
			var v interface{}
			if ok {
				v = string(b)
			} else {
				v = val
			}

			car[col] = v
		}

		cars = append(cars, car)
	}

	return helper.SuccessResponse(c, cars)
}

// ContainLettersHandler godoc
// @ID contain-letters-handler
// @Accept json
// @Produce json
// @Param ContainLettersRequest body model.ContainLettersRequest true "Contain Letters"
// @Success 200 {object} model.SuccessResult
// @Failure 400 {object} model.BadRequestResult
// @Failure 500 {object} model.InternalServerErrorResult
// @Router /is-contain-letters [post]
func (ctr *Controller) ContainLettersHandler(c echo.Context) error {
	requestBody := &model.ContainLettersRequest{}
	if err := c.Bind(requestBody); err != nil {
		return helper.FailResponse(c, http.StatusBadRequest)
	}

	err := ctr.Validate.Struct(requestBody)
	if err != nil {
		return helper.FailResponse(c, http.StatusBadRequest)
	}

	for _, letter := range requestBody.FirstWord {
		if !strings.Contains(strings.ToLower(requestBody.SecondWord), strings.ToLower(string(letter))) {
			return helper.SuccessResponse(c, false)
		}
	}

	return helper.SuccessResponse(c, true)
}
