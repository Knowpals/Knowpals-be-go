package question

type GenerateVideoExerciseReq struct {
	VideoID uint `json:"video_id"`
}

type KnowledgePoint struct {
	KnowledgeID uint   `json:"knowledge_id"`
	Title       string `json:"title"`
}

type Question struct {
	ID              uint             `json:"id"`
	Type            string           `json:"type"`
	Content         string           `json:"content"`
	Options         []string         `json:"options"`
	Answer          string           `json:"answer"`
	Analysis        string           `json:"analysis"`
	SegmentID       *uint            `json:"segment_id"`
	KnowledgePoints []KnowledgePoint `json:"knowledge_points"`
}

type GenerateVideoExerciseResp struct {
	Questions []Question `json:"questions"`
}

// Review (teacher)
type ReviewListReq struct {
	VideoID uint `uri:"video_id" binding:"required"`
}

type ReviewAddReq struct {
	VideoID  uint     `json:"video_id" binding:"required"`
	TimeMs   int64    `json:"time_ms" binding:"required"` // 用于定位 segment
	Type     string   `json:"type" binding:"required"`    // choice|fill|judge
	Content  string   `json:"content" binding:"required"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer" binding:"required"`
	Analysis string   `json:"analysis"`
}

type ReviewUpdateReq struct {
	QuestionID   uint     `uri:"question_id" binding:"required"`
	Type         *string  `json:"type"`
	Content      *string  `json:"content"`
	Options      []string `json:"options"`
	Answer       *string  `json:"answer"`
	Analysis     *string  `json:"analysis"`
	KnowledgePKs []uint   `json:"knowledge_pks"`
}

type ReviewDeleteReq struct {
	QuestionID uint `uri:"question_id" binding:"required"`
}

type StudentAnswer struct {
	QuestionID uint   `json:"question_id"`
	Answer     string `json:"answer"`
	TimeCost   int    `json:"time_cost"`
}

type Result struct {
	QuestionID uint   `json:"question_id"`
	IsCorrect  bool   `json:"is_correct"`
	Answer     string `json:"answer"`
	Analysis   string `json:"analysis"`
}

type AnswerQuestionReq struct {
	VideoID        uint            `json:"video_id"`
	StudentAnswers []StudentAnswer `json:"studentanswers"`
}

type AnswerQuestionResp struct {
	Results []Result `json:"results"`
}
