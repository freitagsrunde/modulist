package main

import (
	"github.com/leebenson/conform"
	"gopkg.in/go-playground/validator.v9"
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

			if err.Tag() == "required" {
				errResp[err.Field()] = "Das folgende Feld muss ausgefüllt sein"
			} else if err.Tag() == "min" {
				errResp[err.Field()] = "Das folgende Feld enthält zu wenig Zeichen"
			} else if err.Tag() == "excludesall" {
				errResp[err.Field()] = "Das folgende Feld enthält unerlaubte Zeichen"
			} else if err.Tag() == "containsany" {
				errResp[err.Field()] = "Das folgende Feld enthält keine Zahlen oder Sonderzeichen"
			} else if err.Tag() == "email" {
				errResp[err.Field()] = "Das folgende Feld enthält keine valide Mail-Adresse"
			}
		}

		return errResp
	} else {
		return nil
	}
}
