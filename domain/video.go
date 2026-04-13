package domain

type Video struct {
	ID        uint
	TeacherID uint
	FileKey   string
	Title     string
	Duration  int
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
