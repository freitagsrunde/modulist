package db

import (
	"time"

	"database/sql"
	"html/template"

	"github.com/jinzhu/gorm"
)

// Constants

const (
	// Privileges have to be kept monotonic.
	// That means, a user with a lower integer
	// privilege value will also have all privileges
	// numerically greater than that value.
	// CAUTION: Changes here will need to be reflected
	// to other places, e.g. template and handler
	// functions of admin's users site.
	PRIVILEGE_ADMIN = iota
	PRIVILEGE_REVIEWER

	// Status groups as increasing integer.
	// CAUTION: Changes here will need to be reflected
	// to other places, e.g. template and handler
	// functions of admin's users site and initial user.
	STATUS_GROUP_PROF = iota
	STATUS_GROUP_WIMI
	STATUS_GROUP_STUDI
	STATUS_GROUP_OTHER
)

// Structs

type User struct {
	ID           string `gorm:"primary_key"`
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Mail         string `gorm:"index;not null;unique"`
	MailVerified bool   `gorm:"not null"`
	PasswordHash string `gorm:"not null;unique"`
	StatusGroup  int    `gorm:"not null"`
	Privileges   int    `gorm:"not null"`
	Enabled      bool   `gorm:"not null"`
}

type PasswordLink struct {
	ID          string    `gorm:"primary_key"`
	UserID      string    `gorm:"index;not null"`
	User        User      `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer;"`
	SecretToken string    `gorm:"not null;unique"`
	Expires     time.Time `gorm:"not null"`
}

type Module struct {
	ID                          int            `gorm:"primary_key"`
	ModuleID                    int            `gorm:"not null"`
	Version                     int            `gorm:"not null"`
	Title                       sql.NullString `gorm:"index"`
	TitleEnglish                sql.NullString `gorm:"index"`
	ECTS                        int            `gorm:"not null"`
	Effective                   *time.Time
	Validity                    string `gorm:"not null"`
	Lang                        string `gorm:"not null"`
	MailAddress                 sql.NullString
	Website                     sql.NullString
	AdministrationOffice        sql.NullString
	URL                         string `gorm:"not null;unique"`
	LearningOutcomes            sql.NullString
	LearningOutcomesHTML        template.HTML `gorm:"-"`
	LearningOutcomesEnglish     sql.NullString
	LearningOutcomesEnglishHTML template.HTML `gorm:"-"`
	TeachingContents            sql.NullString
	TeachingContentsHTML        template.HTML `gorm:"-"`
	TeachingContentsEnglish     sql.NullString
	TeachingContentsEnglishHTML template.HTML `gorm:"-"`
	Courses                     []Course      `gorm:"many2many:module_courses;"`
	InstructiveForm             string        `gorm:"not null"`
	InstructiveFormHTML         template.HTML `gorm:"-"`
	OptionalRequirements        string        `gorm:"not null"`
	OptionalRequirementsHTML    template.HTML `gorm:"-"`
	MandatoryRequirements       sql.NullString
	MandatoryRequirementsHTML   template.HTML `gorm:"-"`
	Graded                      bool          `gorm:"not null"`
	TypeOfExamination           string        `gorm:"not null"`
	ExaminationDescription      sql.NullString
	ExaminationDescriptionHTML  template.HTML `gorm:"-"`
	NumberOfTerms               int           `gorm:"not null"`
	ParticipantLimitation       sql.NullInt64
	RegistrationFormalities     sql.NullString
	RegistrationFormalitiesHTML template.HTML `gorm:"-"`
	Script                      bool          `gorm:"not null"`
	ScriptEnglish               bool          `gorm:"not null"`
	Literature                  string        `gorm:"not null"`
	LiteratureHTML              template.HTML `gorm:"-"`
	Miscellaneous               sql.NullString
	MiscellaneousHTML           template.HTML `gorm:"-"`
	ReferencePersonID           sql.NullInt64
	ReferencePerson             Person `gorm:"ForeignKey:ReferencePersonID;AssociationForeignKey:Refer;"`
	ResponsiblePersonID         sql.NullInt64
	ResponsiblePerson           Person `gorm:"ForeignKey:ResponsiblePersonID;AssociationForeignKey:Refer;"`
}

type SQLiteModule struct {
	ID                      int            `gorm:"column:id"`
	Title                   sql.NullString `gorm:"column:title"`
	TitleEnglish            sql.NullString `gorm:"column:titleEnglish"`
	ECTS                    int            `gorm:"column:ects"`
	ModuleID                int            `gorm:"column:moduleID"`
	Version                 int            `gorm:"column:version"`
	Effective               *time.Time     `gorm:"column:effective"`
	Validity                string         `gorm:"column:validity"`
	Lang                    string         `gorm:"column:lang"`
	MailAddress             sql.NullString `gorm:"column:mailAddress"`
	Website                 sql.NullString `gorm:"column:website"`
	AdministrationOffice    sql.NullString `gorm:"column:administrationOffice"`
	LearningOutcomes        sql.NullString `gorm:"column:learningOutcomes"`
	LearningOutcomesEnglish sql.NullString `gorm:"column:learningOutcomesEnglish"`
	TeachingContents        sql.NullString `gorm:"column:teachingContents"`
	TeachingContentsEnglish sql.NullString `gorm:"column:teachingContentsEnglish"`
	URL                     string         `gorm:"column:url"`
	InstructiveForm         string         `gorm:"column:instructiveForm"`
	OptionalRequirements    string         `gorm:"column:optionalRequirements"`
	MandatoryRequirements   sql.NullString `gorm:"column:mandatoryRequirements"`
	Graded                  bool           `gorm:"column:graded"`
	TypeOfExamination       string         `gorm:"column:typeOfExamination"`
	ExaminationDescription  sql.NullString `gorm:"column:examinationDescription"`
	NumberOfTerms           int            `gorm:"column:numberOfTerms"`
	ParticipantLimitation   sql.NullInt64  `gorm:"column:participantLimitation"`
	Miscellaneous           sql.NullString `gorm:"column:miscellaneous"`
	Script                  bool           `gorm:"column:script"`
	ScriptEnglish           bool           `gorm:"column:scriptEnglish"`
	Literature              string         `gorm:"column:literature"`
	ReferencePersonID       sql.NullInt64  `gorm:"column:referencePerson_id"`
	ReferencePerson         Person         `gorm:"ForeignKey:ReferencePersonID;AssociationForeignKey:Refer;"`
	ResponsiblePersonID     sql.NullInt64  `gorm:"column:responsiblePerson_id"`
	ResponsiblePerson       Person         `gorm:"ForeignKey:ResponsiblePersonID;AssociationForeignKey:Refer;"`
	RegistrationFormalities sql.NullString `gorm:"column:registrationFormalities"`
}

func (sqliteModule *SQLiteModule) TableName() string {
	return "modulecrawler_mtsmodule"
}

func (sqliteModule SQLiteModule) ToModule(db *gorm.DB) Module {

	// Find in module involved persons in database.
	var RefPerson Person
	var ResPerson Person
	db.First(&RefPerson, "\"id\" = ?", sqliteModule.ReferencePersonID)
	db.First(&ResPerson, "\"id\" = ?", sqliteModule.ResponsiblePersonID)

	return Module{
		ID:                      sqliteModule.ID,
		Title:                   sqliteModule.Title,
		TitleEnglish:            sqliteModule.TitleEnglish,
		ECTS:                    sqliteModule.ECTS,
		ModuleID:                sqliteModule.ModuleID,
		Version:                 sqliteModule.Version,
		Effective:               sqliteModule.Effective,
		Validity:                sqliteModule.Validity,
		Lang:                    sqliteModule.Lang,
		MailAddress:             sqliteModule.MailAddress,
		Website:                 sqliteModule.Website,
		AdministrationOffice:    sqliteModule.AdministrationOffice,
		LearningOutcomes:        sqliteModule.LearningOutcomes,
		LearningOutcomesEnglish: sqliteModule.LearningOutcomesEnglish,
		TeachingContents:        sqliteModule.TeachingContents,
		TeachingContentsEnglish: sqliteModule.TeachingContentsEnglish,
		URL:                     sqliteModule.URL,
		InstructiveForm:         sqliteModule.InstructiveForm,
		OptionalRequirements:    sqliteModule.OptionalRequirements,
		MandatoryRequirements:   sqliteModule.MandatoryRequirements,
		Graded:                  sqliteModule.Graded,
		TypeOfExamination:       sqliteModule.TypeOfExamination,
		ExaminationDescription:  sqliteModule.ExaminationDescription,
		NumberOfTerms:           sqliteModule.NumberOfTerms,
		ParticipantLimitation:   sqliteModule.ParticipantLimitation,
		Miscellaneous:           sqliteModule.Miscellaneous,
		Script:                  sqliteModule.Script,
		ScriptEnglish:           sqliteModule.ScriptEnglish,
		Literature:              sqliteModule.Literature,
		ReferencePersonID:       sqliteModule.ReferencePersonID,
		ReferencePerson:         RefPerson,
		ResponsiblePersonID:     sqliteModule.ResponsiblePersonID,
		ResponsiblePerson:       ResPerson,
		RegistrationFormalities: sqliteModule.RegistrationFormalities,
	}
}

type Person struct {
	ID        int    `gorm:"primary_key"`
	FirstName string `gorm:"index;not null"`
	LastName  string `gorm:"index;not null"`
}

type SQLitePerson struct {
	ID        int    `gorm:"column:id"`
	FirstName string `gorm:"column:firstname"`
	LastName  string `gorm:"column:lastname"`
}

func (sqlitePerson *SQLitePerson) TableName() string {
	return "modulecrawler_person"
}

func (sqlitePerson SQLitePerson) ToPerson() Person {

	return Person{
		ID:        sqlitePerson.ID,
		FirstName: sqlitePerson.FirstName,
		LastName:  sqlitePerson.LastName,
	}
}

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

type ModuleCourses struct {
	ModuleID int `gorm:"primary_key"`
	CourseID int `gorm:"primary_key"`
}

type SQLiteModuleCourses struct {
	ID       int `gorm:"column:id"`
	ModuleID int `gorm:"column:mtsmodule_id"`
	CourseID int `gorm:"column:course_id"`
}

func (sqliteModuleCourses *SQLiteModuleCourses) TableName() string {
	return "modulecrawler_mtsmodule_courses"
}

func (sqliteModuleCourses SQLiteModuleCourses) ToModuleCourses() ModuleCourses {

	return ModuleCourses{
		ModuleID: sqliteModuleCourses.ModuleID,
		CourseID: sqliteModuleCourses.CourseID,
	}
}
