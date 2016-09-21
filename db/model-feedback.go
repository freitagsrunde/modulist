package db

// Constants

const (
	// This category counter assings each part
	// of a module description an integer identifier
	// to map feedback to those individual parts
	// of the description.
	CATEGORY_HEADER = iota
	CATEGORY_LEARNING_OUTCOMES
	CATEGORY_TEACHING_CONTENTS
	CATEGORY_COURSES
	CATEGORY_WORKING_EFFORT
	CATEGORY_INSTRUCTIVE_FORM
	CATEGORY_REQUIREMENTS
	CATEGORY_EXAMINATION
	CATEGORY_NUMBER_TERMS
	CATEGORY_PARTICIPANT_LIMITATION
	CATEGORY_REGISTRATION_FORMALITIES
	CATEGORY_SCRIPT
	CATEGORY_LITERATURE
	CATEGORY_MISCELLANEOUS
)

// Structs

type Feedback struct {
	ID       int    `gorm:"primary_key"`
	ModuleID int    `gorm:"index;not null"`
	UserID   string `gorm:"index;not null"`
	Category int    `gorm:"not null"`
	Comment  string `gorm:"not null"`
}

// Functions

// CategoriesByName returns a map of all
// categories queryable by name of category.
func CategoriesByName() map[string]int {

	Categories := make(map[string]int)

	Categories["Header"] = CATEGORY_HEADER
	Categories["LearningOutcomes"] = CATEGORY_LEARNING_OUTCOMES
	Categories["TeachingContents"] = CATEGORY_TEACHING_CONTENTS
	Categories["Courses"] = CATEGORY_COURSES
	Categories["WorkingEffort"] = CATEGORY_WORKING_EFFORT
	Categories["InstructiveForm"] = CATEGORY_INSTRUCTIVE_FORM
	Categories["Requirements"] = CATEGORY_REQUIREMENTS
	Categories["Examination"] = CATEGORY_EXAMINATION
	Categories["NumberOfTerms"] = CATEGORY_NUMBER_TERMS
	Categories["ParticipantLimitation"] = CATEGORY_PARTICIPANT_LIMITATION
	Categories["RegistrationFormalities"] = CATEGORY_REGISTRATION_FORMALITIES
	Categories["Script"] = CATEGORY_SCRIPT
	Categories["Literature"] = CATEGORY_LITERATURE
	Categories["Miscellaneous"] = CATEGORY_MISCELLANEOUS

	return Categories
}
