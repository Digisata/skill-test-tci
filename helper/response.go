package helper

import (
	"net/http"

	"github.com/Digisata/skill_test_tci/model"
	"github.com/labstack/echo/v4"
)

func SuccessResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, model.SuccessResult{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   data,
	})
}

func FailResponse(c echo.Context, code int) error {
	switch code {
	case http.StatusBadRequest:
		return c.JSON(http.StatusBadRequest, model.BadRequestResult{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
		})
	case http.StatusNotFound:
		return c.JSON(http.StatusNotFound, model.NotFoundResult{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
		})
	case http.StatusInternalServerError:
		return c.JSON(http.StatusInternalServerError, model.InternalServerErrorResult{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
		})
	}
	return nil
}
