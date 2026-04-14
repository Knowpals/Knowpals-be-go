package domain

// Question 流水线 quiz 阶段落库用（知识点已解析为 knowledge_points 主键）
type Question struct {
	ID           uint
	SegmentID    *uint
	Type         string
	Content      string
	Options      []string
	Answer       string
	Analysis     string
	KnowledgePKs []uint
}
