package message

// Kafka pipeline stages
const (
	StageReady     = "ready"     //准备阶段（刚创建任务）
	StageKnowledge = "knowledge" // 分段 + 知识点
	StageQuiz      = "quiz"      // 题目生成
)

type TaskMessage struct {
	JobID   string                 `json:"job_id"`
	Stage   string                 `json:"stage"`
	Payload map[string]interface{} `json:"payload"`
	Retry   int                    `json:"retry"`
}

type ResultMessage struct {
	JobID  string                 `json:"job_id"`
	Stage  string                 `json:"stage"`
	Status string                 `json:"status"`
	Result map[string]interface{} `json:"result"`
	Error  string                 `json:"error"`
}

type KnowledgeSegmentPayload struct {
	Concepts   []Concept `json:"concepts"`
	Segments   []Segment `json:"segments"`
	DurationMs int       `json:"duration_ms"`
}

type Concept struct {
	ConceptID string `json:"concept_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}

type Segment struct {
	SegmentID string `json:"segment_id"`
	ConceptID string `json:"concept_id"`
	Text      string `json:"text"`
	StartMs   int    `json:"start_ms"`
	EndMs     int    `json:"end_ms"`
}

type QuizPayload struct {
	Type       string   `json:"type"`
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	Answer     string   `json:"answer"`
	Analysis   string   `json:"analysis"`
	Difficulty string   `json:"difficulty"`
	ConceptID  string   `json:"concept_id"`
}
