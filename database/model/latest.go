package model

const MainNotice = "main"
const LibraryNotice = "library"
const InstagramNotice = "instagram"
const StudentNotice = "student"
const SWNotice = "sw"

type Latest struct {
	NoticeType string `gorm:"primaryKey;size:20"`
	URL        string `gorm:"not null;size:100"`
}

func (l *Latest) TableName() string {
	return "latest"
}
