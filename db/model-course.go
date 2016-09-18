package db

import (
	"database/sql"
)

// Structs

type Course struct {
	ID                  int    `gorm:"primary_key"`
	Title               string `gorm:"not null"`
	CourseType          sql.NullString
	CourseID            sql.NullString
	Cycle               sql.NullString
	CreditHours         sql.NullInt64
	Annotation          sql.NullString
	Content             sql.NullString
	CourseURL           sql.NullString
	DetailedDescription sql.NullString
	Requirements        sql.NullString
	Audience            sql.NullString
	Comment             sql.NullString
	CourseAssessment    sql.NullString
	Literature          sql.NullString
	TeachingContents    sql.NullString
}

type SQLiteCourse struct {
	ID                  int            `gorm:"column:id"`
	Title               string         `gorm:"column:title"`
	CourseType          sql.NullString `gorm:"column:courseType"`
	CourseID            sql.NullString `gorm:"column:courseID"`
	Cycle               sql.NullString `gorm:"column:cycle"`
	CreditHours         sql.NullInt64  `gorm:"column:creditHours"`
	Annotation          sql.NullString `gorm:"column:annotation"`
	Content             sql.NullString `gorm:"column:content"`
	CourseURL           sql.NullString `gorm:"column:courseURL"`
	DetailedDescription sql.NullString `gorm:"column:detailedDescription"`
	Requirements        sql.NullString `gorm:"column:requirements"`
	Audience            sql.NullString `gorm:"column:audience"`
	Comment             sql.NullString `gorm:"column:comment"`
	CourseAssessment    sql.NullString `gorm:"column:courseAssessment"`
	Literature          sql.NullString `gorm:"column:literature"`
	TeachingContents    sql.NullString `gorm:"column:teachingContents"`
}

func (sqliteCourse *SQLiteCourse) TableName() string {
	return "modulecrawler_course"
}

func (sqliteCourse SQLiteCourse) ToCourse() Course {

	return Course{
		ID:                  sqliteCourse.ID,
		Title:               sqliteCourse.Title,
		CourseType:          sqliteCourse.CourseType,
		CourseID:            sqliteCourse.CourseID,
		Cycle:               sqliteCourse.Cycle,
		CreditHours:         sqliteCourse.CreditHours,
		Annotation:          sqliteCourse.Annotation,
		Content:             sqliteCourse.Content,
		CourseURL:           sqliteCourse.CourseURL,
		DetailedDescription: sqliteCourse.DetailedDescription,
		Requirements:        sqliteCourse.Requirements,
		Audience:            sqliteCourse.Audience,
		Comment:             sqliteCourse.Comment,
		CourseAssessment:    sqliteCourse.CourseAssessment,
		Literature:          sqliteCourse.Literature,
		TeachingContents:    sqliteCourse.TeachingContents,
	}
}

type SQLiteModuleCourses struct {
	ID       int `gorm:"column:id"`
	ModuleID int `gorm:"column:mtsmodule_id"`
	CourseID int `gorm:"column:course_id"`
}

func (sqliteModuleCourses *SQLiteModuleCourses) TableName() string {
	return "modulecrawler_mtsmodule_courses"
}
