package api

import (
	api_utils "simplebank/api/utils"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return api_utils.IsSupportedCurrency(currency)
	}
	return false
}
