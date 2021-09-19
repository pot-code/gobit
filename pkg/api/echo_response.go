package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/validate"
)

func ValidateFailedResponse(c echo.Context, err *validate.ValidationErrors) error {
	return c.JSON(http.StatusBadRequest, NewRESTValidationError(ErrFailedValidate.Error(), err))
}

func BadRequestResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, NewRESTStandardError(err.Error()))
}

func JsonResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func StatusResponse(c echo.Context, code int) error {
	return c.String(code, http.StatusText(code))
}

func BindErrorResponse(c echo.Context, err error) error {
	if he, ok := err.(*echo.HTTPError); ok {
		return c.JSON(http.StatusBadRequest, NewRESTBindingError(ErrFailedBinding.Error(), he.Message))
	}
	if be, ok := err.(*echo.BindingError); ok {
		return c.JSON(http.StatusBadRequest, NewRESTBindingError(ErrFailedBinding.Error(), be.Message))
	}
	return c.JSON(http.StatusBadRequest, NewRESTStandardError(ErrFailedBinding.Error()))
}
