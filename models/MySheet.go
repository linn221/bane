package models

type MySheet struct {
	Id           int    `gorm:"primaryKey"`
	Title        string `gorm:"not null"`
	Body         string `gorm:"type:text"`
	Created      MyDate `gorm:"not null"`
	NextDate     MyDate `gorm:"index;not null"`
	PreviousDate MyDate `gorm:"index;not null"`
	LabelId      int    `gorm:"index;default:null"`
	Index        int    `gorm:"not null;default:0"`
	// Age is calculated, not stored
}

type MySheetInput struct {
	Title string  `json:"title"`
	Body  string  `json:"body"`
	Alias string  `json:"alias,omitempty"`
	Date  *MyDate `json:"date,omitempty"`
	Label *string
}

type MySheetFilter struct {
	Title        string  `json:"title,omitempty"`
	Search       string  `json:"search,omitempty"`
	NextDate     *MyDate `json:"nextDate,omitempty"`
	PreviousDate *MyDate `json:"previousDate,omitempty"`
	Label        *string
}

type MySheetLabel struct {
	Id   int    `gorm:"primaryKey" json:"id"`
	Name string `gorm:"index;not null" json:"name"`
}

type MySheetLabelInput struct {
	Name  string `json:"name"`
	Alias string `json:"alias,omitempty"`
}
