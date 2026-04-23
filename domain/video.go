package domain

import "time"

type Video struct {
	ID        uint
	TeacherID uint
	FileKey   string
	Title     string
	Duration  int
	CreatedAt time.Time
	Deadline  time.Time
	ReviewStatus string
	ReviewedAt   *time.Time
	PublishedAt  *time.Time
}

type Segment struct {
	ID          uint
	SegmentID   string
	VideoID     uint
	Start       int
	End         int
	KnowledgeID string
	Text        string
}
