package db

import (
	"database/sql"
)

// Structs

type WorkingEffort struct {
	ID          int           `gorm:"primary_key"`
	ModuleID    sql.NullInt64 `gorm:"index"`
	CourseID    sql.NullInt64 `gorm:"index"`
	Description string        `gorm:"not null"`
	Category    string        `gorm:"index;not null"`
	Multiplier  float32       `gorm:"not null"`
	Hours       float32       `gorm:"not null"`
	Total       float32       `gorm:"not null"`
}

type SQLiteWorkingEffort struct {
	ID          int           `gorm:"column:id"`
	ModuleID    sql.NullInt64 `gorm:"column:module_id"`
	CourseID    sql.NullInt64 `gorm:"column:course_id"`
	Description string        `gorm:"column:description"`
	Category    string        `gorm:"column:category"`
	Multiplier  float32       `gorm:"column:multiplier"`
	Hours       float32       `gorm:"column:hours"`
	Total       float32       `gorm:"column:total"`
}

func (sqliteWorkingEffort *SQLiteWorkingEffort) TableName() string {
	return "modulecrawler_workingeffort"
}

func (sqliteWorkingEffort SQLiteWorkingEffort) ToWorkingEffort() WorkingEffort {

	return WorkingEffort{
		ID:          sqliteWorkingEffort.ID,
		ModuleID:    sqliteWorkingEffort.ModuleID,
		CourseID:    sqliteWorkingEffort.CourseID,
		Description: sqliteWorkingEffort.Description,
		Category:    sqliteWorkingEffort.Category,
		Multiplier:  sqliteWorkingEffort.Multiplier,
		Hours:       sqliteWorkingEffort.Hours,
		Total:       sqliteWorkingEffort.Total,
	}
}

// WorkingEffortTemplate represents the most useful
// structure of working effort information for use
// in HTML templates.
type WorkingEffortTemplate struct {
	Category    string
	Efforts     []WorkingEffortEfforts
	CourseTotal float32
}

// WorkingEffortEfforts is a sub type for more
// fittingly handle working efforts in HTML templates.
type WorkingEffortEfforts struct {
	Description string
	Multiplier  float32
	Hours       float32
	Total       float32
}

// (*WorkingEffort).Convert transfers the working effort
// elements from the database representation to another one
// better suited to be used in HTML templates.
func WorkingEffortsConvert(workingEfforts []WorkingEffort) []WorkingEffortTemplate {

	// Reserve space for final slice and helper map.
	finalElements := make([]WorkingEffortTemplate, 0)
	seen := make(map[string]bool)

	// Range over all working effort elements available.
	for _, w := range workingEfforts {

		// Prepare effort element to add in any following case.
		addEffort := WorkingEffortEfforts{
			Description: w.Description,
			Multiplier:  w.Multiplier,
			Hours:       w.Hours,
			Total:       w.Total,
		}

		// Check if category of element is a new one.
		if seen[w.Category] != true {

			// If so, create a new element and add effort element to.
			var addElement WorkingEffortTemplate
			addElement.Category = w.Category
			addElement.Efforts = make([]WorkingEffortEfforts, 1)
			addElement.Efforts[0] = addEffort

			// Add this new element to the final list.
			finalElements = append(finalElements, addElement)

			// Mark category as seen.
			seen[w.Category] = true
		} else {

			// If not yet seen, range over all already added elements
			// in final list and compare for matching category.
			for i, _ := range finalElements {

				// If found, add the effort element.
				if finalElements[i].Category == w.Category {
					finalElements[i].Efforts = append(finalElements[i].Efforts, addEffort)
				}
			}
		}
	}

	// Finally, range over all elements in order
	// to calculate per-course total of projected
	// working effort hours.
	for i, _ := range finalElements {

		finalElements[i].CourseTotal = 0

		for _, e := range finalElements[i].Efforts {
			finalElements[i].CourseTotal += e.Total
		}
	}

	return finalElements
}
