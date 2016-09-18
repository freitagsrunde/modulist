package db

// Structs

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
