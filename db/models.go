package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
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
	User        User      `gorm:"ForeignKey:UserID;"`
	SecretToken string    `gorm:"not null;unique"`
	Expires     time.Time `gorm:"not null"`
}

type Module struct {
	ID                      string `gorm:"primary_key"`
	Title                   string `gorm:"index;not null"`
	TitleEnglish            string `gorm:"index"`
	ECTS                    int
	ModuleID                int `gorm:"not null"`
	Version                 int `gorm:"not null"`
	Effective               time.Time
	Validity                string `gorm:"not null"`
	Lang                    string `gorm:"not null"`
	MailAddress             string
	Website                 string
	AdministrationOffice    string
	LearningOutcomes        string
	LearningOutcomesEnglish string
	TeachingContents        string
	TeachingContentsEnglish string
	URL                     string `gorm:"not null;unique"`
	InstructiveForm         string `gorm:"not null"`
	OptionalRequirements    string `gorm:"not null"`
	MandatoryRequirements   string
	Graded                  bool   `gorm:"not null"`
	TypeOfExamination       string `gorm:"not null"`
	ExaminationDescription  string
	NumberOfTerms           int `gorm:"not null"`
	ParticipantLimitation   int
	Miscellaneous           string
	Script                  bool   `gorm:"not null"`
	ScriptEnglish           bool   `gorm:"not null"`
	Literature              string `gorm:"not null"`
	ReferencePersonID       int
	ReferencePerson         Person `gorm:"ForeignKey:ReferencePersonID;AssociationForeignKey:Refer;"`
	ResponsiblePersonID     int
	ResponsiblePerson       Person `gorm:"ForeignKey:ResponsiblePersonID;AssociationForeignKey:Refer;"`
	RegistrationFormalities string
}

type SQLiteModule struct {
	ID                      string    `gorm:"column:id"`
	Title                   string    `gorm:"column:title"`
	TitleEnglish            string    `gorm:"column:titleEnglish"`
	ECTS                    int       `gorm:"column:ects"`
	ModuleID                int       `gorm:"column:moduleID"`
	Version                 int       `gorm:"column:version"`
	Effective               time.Time `gorm:"column:effective"`
	Validity                string    `gorm:"column:validity"`
	Lang                    string    `gorm:"column:lang"`
	MailAddress             string    `gorm:"column:mailAddress"`
	Website                 string    `gorm:"column:website"`
	AdministrationOffice    string    `gorm:"column:administrationOffice"`
	LearningOutcomes        string    `gorm:"column:learningOutcomes"`
	LearningOutcomesEnglish string    `gorm:"column:learningOutcomesEnglish"`
	TeachingContents        string    `gorm:"column:teachingContents"`
	TeachingContentsEnglish string    `gorm:"column:teachingContentsEnglish"`
	URL                     string    `gorm:"column:url"`
	InstructiveForm         string    `gorm:"column:instructiveForm"`
	OptionalRequirements    string    `gorm:"column:optionalRequirements"`
	MandatoryRequirements   string    `gorm:"column:mandatoryRequirements"`
	Graded                  bool      `gorm:"column:graded"`
	TypeOfExamination       string    `gorm:"column:typeOfExamination"`
	ExaminationDescription  string    `gorm:"column:examinationDescription"`
	NumberOfTerms           int       `gorm:"column:numberOfTerms"`
	ParticipantLimitation   int       `gorm:"column:participantLimitation"`
	Miscellaneous           string    `gorm:"column:miscellaneous"`
	Script                  bool      `gorm:"column:script"`
	ScriptEnglish           bool      `gorm:"column:scriptEnglish"`
	Literature              string    `gorm:"column:literature"`
	ReferencePersonID       int       `gorm:"column:referencePerson_id"`
	ReferencePerson         Person    `gorm:"ForeignKey:ReferencePersonID;AssociationForeignKey:Refer;"`
	ResponsiblePersonID     int       `gorm:"column:responsiblePerson_id"`
	ResponsiblePerson       Person    `gorm:"ForeignKey:ResponsiblePersonID;AssociationForeignKey:Refer;"`
	RegistrationFormalities string    `gorm:"column:registrationFormalities"`
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
		ID:                      fmt.Sprintf("%s", uuid.NewV4()),
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
