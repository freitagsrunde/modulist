package main

import (
	"github.com/go-playground/validator"
	"github.com/leebenson/conform"
)

// Functions

// ConformAndValidate checks supplied data p on
// conformity and validity to requirements.
func (app *App) ConformAndValidate(Payload interface{}) map[string]string {

	// Validate data in payload.
	conform.Strings(Payload)
	errs := app.Validator.Struct(Payload)

	if errs != nil {

		errResp := make(map[string]string)

		// Iterate over all validation errors.
		for _, err := range errs.(validator.ValidationErrors) {

			if err.Tag == "required" {
				errResp[err.Field] = "The following field is required"
			} else if err.Tag == "email" {
				errResp[err.Field] = "The following field does not contain a valid mail address"
			}
		}

		return errResp
	} else {
		return nil
	}
}
