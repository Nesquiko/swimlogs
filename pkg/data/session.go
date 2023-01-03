package data

type Session struct {
	Base
	DurationMin int    `gorm:"not null"`
	StartTime   string `gorm:"not null;default:null"`
}
