package model

type Subscriber struct {
	ID        uint64 `gorm:"primaryKey"`
	URL       string `gorm:"unique;not null"`
	Main      bool   `gorm:"not null"`
	Library   bool   `gorm:"not null"`
	Instagram bool   `gorm:"not null"`
	Student   bool   `gorm:"not null"`
	Sanhak    bool   `gorm:"not null"`
	SW        bool   `gorm:"not null"`
}

func (s *Subscriber) TableName() string {
	return "subscriber"
}
