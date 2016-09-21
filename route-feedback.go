package main

import (
	"fmt"
	"strconv"
	"strings"

	"html/template"
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Structs

type ReviewModulePayload struct {
	ID int `conform:"trim,num" validate:"required,min=1"`
}

type AddFeedbackPayload struct {
	Category int    `form:"category" conform:"trim,num" validate:"min=0"`
	Comment  string `form:"comment" conform:"trim" validate:"required"`
}

// Functions

func (app *App) ReviewModule(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Extract ID of module to review from URL.
	id, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	Payload := ReviewModulePayload{ID: id}

	// Check supplied ID for conformity and validity.
	if errs := app.ConformAndValidate(&Payload); errs != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Retrieve module information for supplied ID from database.
	var Module db.Module
	app.DB.Preload("Courses").Preload("WorkingEfforts").Preload("ExamElements").First(&Module, "\"id\" = ?", Payload.ID)

	// Check if a module for supplied ID exists.
	if Module.URL == "" {

		// If no module exists, redirect to modules list view.
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Convert the working effort elements from database to
	// structure better suited for displaying in HTML template.
	Module.WorkingEffortsHTML = db.WorkingEffortsConvert(Module.WorkingEfforts)

	// Only use this field if it contains a valid value.
	if Module.ReferencePersonID.Valid {
		app.DB.Model(&Module).Related(&Module.ReferencePerson, "ReferencePersonID")
	}

	// Only use this field if it contains a valid value.
	if Module.ResponsiblePersonID.Valid {
		app.DB.Model(&Module).Related(&Module.ResponsiblePerson, "ResponsiblePersonID")
	}

	// For each of the elements containing free text with linebreaks,
	// convert each of those to HTML linebreaks in a safely manner.
	Module.LearningOutcomesHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.LearningOutcomes.String), "\n", "<br />", -1))
	Module.LearningOutcomesEnglishHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.LearningOutcomesEnglish.String), "\n", "<br />", -1))
	Module.TeachingContentsHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.TeachingContents.String), "\n", "<br />", -1))
	Module.TeachingContentsEnglishHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.TeachingContentsEnglish.String), "\n", "<br />", -1))
	Module.InstructiveFormHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.InstructiveForm), "\n", "<br />", -1))
	Module.OptionalRequirementsHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.OptionalRequirements), "\n", "<br />", -1))
	Module.MandatoryRequirementsHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.MandatoryRequirements.String), "\n", "<br />", -1))
	Module.ExaminationDescriptionHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.ExaminationDescription.String), "\n", "<br />", -1))
	Module.MiscellaneousHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.Miscellaneous.String), "\n", "<br />", -1))
	Module.LiteratureHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.Literature), "\n", "<br />", -1))
	Module.RegistrationFormalitiesHTML = template.HTML(strings.Replace(template.HTMLEscapeString(Module.RegistrationFormalities.String), "\n", "<br />", -1))

	c.HTML(http.StatusOK, "module-feedback.html", gin.H{
		"PageTitle":  fmt.Sprintf("Feedback zu Modul #%d", Module.ModuleID),
		"User":       User,
		"Module":     Module,
		"Categories": db.CategoriesByName(),
	})
}

func (app *App) AddFeedback(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Extract ID of module to add a review for from URL.
	id, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	IDPayload := ReviewModulePayload{ID: id}

	// Check supplied ID for conformity and validity.
	if errs := app.ConformAndValidate(&IDPayload); errs != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	var Module db.Module
	app.DB.First(&Module, "\"id\" = ?", IDPayload.ID)

	// Check if a module for supplied ID exists.
	if Module.URL == "" {

		// If it does not, send error message and return.
		c.JSON(http.StatusBadRequest, gin.H{
			"Reason": "Module does not exist.",
		})

		return
	}

	var FeedbackPayload AddFeedbackPayload

	err = c.BindWith(&FeedbackPayload, binding.FormPost)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"Reason": "Internal error. Please try again later.",
		})

		return
	}

	// Check sent content for validity.
	ErrorDesc := app.ConformAndValidate(&FeedbackPayload)
	if ErrorDesc != nil {

		// If payload did not pass, report errors to user.
		c.JSON(http.StatusBadRequest, gin.H{
			"Reason":            "Malformed input. Please check your values for validity and try again.",
			"ErrorDescriptions": ErrorDesc,
		})

		return
	}

	// With supplied content, create new feedback.
	var NewFeedback db.Feedback

	NewFeedback.ModuleID = IDPayload.ID
	NewFeedback.UserID = User.ID
	NewFeedback.Category = FeedbackPayload.Category
	NewFeedback.Comment = FeedbackPayload.Comment

	// Save feedback to database.
	app.DB.Create(&NewFeedback)

	// Request all feedback comments for submitted
	// category and include count of those.
	var AllFeedback []db.Feedback
	app.DB.Order("\"id\" asc").Find(&AllFeedback, "\"module_id\" = ? AND \"category\" = ?", IDPayload.ID, FeedbackPayload.Category)

	// Return success as JSON to user.
	c.JSON(http.StatusOK, gin.H{
		"Success":  true,
		"Feedback": AllFeedback,
		"Count":    len(AllFeedback),
	})
}

func (app *App) DeleteFeedback(c *gin.Context) {}

func (app *App) ListFeedback(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_REVIEWER)
	if err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Extract ID of module to fetch comments for from URL.
	id, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	Payload := ReviewModulePayload{ID: id}

	// Check supplied ID for conformity and validity.
	if errs := app.ConformAndValidate(&Payload); errs != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	var Module db.Module
	app.DB.First(&Module, "\"id\" = ?", Payload.ID)

	// Check if a module for supplied ID exists.
	if Module.URL == "" {

		// If it does not, send error message and return.
		c.JSON(http.StatusBadRequest, gin.H{
			"Reason": "Module does not exist.",
		})

		return
	}

	// Fetch all available feedback for that module,
	// calculate counts and return data to user.
	var AllFeedback []db.Feedback
	app.DB.Order("\"category\" asc").Order("\"id\" asc").Find(&AllFeedback, "\"module_id\" = ?", Payload.ID)

	// TODO: Add calculation of counters per category.

	c.JSON(http.StatusOK, gin.H{
		"Success":  true,
		"Feedback": AllFeedback,
	})
}
