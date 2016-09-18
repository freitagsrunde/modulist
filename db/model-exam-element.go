package db

// Structs

type ExamElement struct {
	ID          int    `gorm:"primary_key"`
	ModuleID    int    `gorm:"index;not null"`
	Description string `gorm:"not null"`
	Points      int    `gorm:"not null"`
}

type SQLiteExamElement struct {
	ID          int    `gorm:"column:id"`
	ModuleID    int    `gorm:"column:module_id"`
	Description string `gorm:"column:description"`
	Points      int    `gorm:"column:points"`
}

func (sqliteExamElement *SQLiteExamElement) TableName() string {
	return "modulecrawler_examelement"
}

func (sqliteExamElement SQLiteExamElement) ToExamElement() ExamElement {

	return ExamElement{
		ID:          sqliteExamElement.ID,
		ModuleID:    sqliteExamElement.ModuleID,
		Description: sqliteExamElement.Description,
		Points:      sqliteExamElement.Points,
	}
}
