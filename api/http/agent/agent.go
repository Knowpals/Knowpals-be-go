package agent

import "time"

type ChatReq struct {
	Text        string `json:"text" binding:"required"`
	VideoID     string `json:"video_id"`
	KnowledgeID string `json:"knowledge_id"`
}

type ChatResp struct {
	Reply   string `json:"reply"`
	Context string `json:"context"`
	VideoID string `json:"video_id"`
}

type GetChatHistoryReq struct {
	VideoID string `form:"video_id"`
	Limit   int    `form:"limit"`
}

type ChatMessage struct {
	Role        string    `json:"role"`
	Text        string    `json:"text"`
	VideoID     string    `json:"video_id"`
	KnowledgeID string    `json:"knowledge_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetChatHistoryResp struct {
	Messages []ChatMessage `json:"messages"`
}

type GenerateQuizReq struct {
	VideoID       string `json:"video_id" binding:"required"`
	NumQuestions  int32  `json:"num_questions" binding:"required"`
}

type QuizItem struct {
	KnowledgeID string   `json:"knowledge_id"`
	Type        string   `json:"type"`
	Question    string   `json:"question"`
	Options     []string `json:"options"`
	Answer      string   `json:"answer"`
	Analysis    string   `json:"analysis"`
	Difficulty  string   `json:"difficulty"`
}

type GenerateQuizResp struct {
	Quizzes []QuizItem `json:"quizzes"`
}

type GenerateReportReq struct {
	VideoID string `json:"video_id" binding:"required"`
	ForceRegen bool `json:"force_regen"`
}

type ReportItem struct {
	KnowledgeID         string       `json:"knowledge_id"`
	Mastery             float64      `json:"mastery"`
	Summary             string       `json:"summary"`
	Weakness            []string     `json:"weakness"`
	BehaviorPattern     []string     `json:"behavior_pattern"`
	Trend               string       `json:"trend"`
	RecommendedSegments []SegmentRef `json:"recommended_segments"`
}

type SegmentRef struct {
	VideoID   string `json:"video_id"`
	SegmentID string `json:"segment_id"`
	StartMs   int64  `json:"start_ms"`
	EndMs     int64  `json:"end_ms"`
}

type GenerateReportResp struct {
	VideoID         string       `json:"video_id"`
	Items           []ReportItem  `json:"items"`
	OverallSummary  string       `json:"overall_summary"`
}

type GetReportReq struct {
	VideoID string `form:"video_id" binding:"required"`
}

