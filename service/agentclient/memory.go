package agentclient

import (
	"context"
	"strconv"
	"time"

	"github.com/Knowpals/Knowpals-be-go/api/grpc/memorypb"
)

type MemoryWriter interface {
	WritePause(ctx context.Context, studentID uint, knowledgeID string, videoID uint, segmentID uint) error
	WriteReplay(ctx context.Context, studentID uint, knowledgeID string, videoID uint, segmentID uint) error
	WriteQuestion(ctx context.Context, studentID uint, knowledgeID string, videoID uint, segmentID *uint, questionID uint, isCorrect bool, userAnswer string, rightAnswer string, contentIfWrong string) error
}

type memoryWriter struct {
	client memorypb.MemoryServiceClient
}

func NewMemoryWriter(client memorypb.MemoryServiceClient) MemoryWriter {
	return &memoryWriter{client: client}
}

func (w *memoryWriter) WritePause(ctx context.Context, studentID uint, knowledgeID string, videoID uint, segmentID uint) error {
	return w.writeSimple(ctx, studentID, knowledgeID, videoID, segmentID, memorypb.EventType_PAUSE)
}

func (w *memoryWriter) WriteReplay(ctx context.Context, studentID uint, knowledgeID string, videoID uint, segmentID uint) error {
	return w.writeSimple(ctx, studentID, knowledgeID, videoID, segmentID, memorypb.EventType_REPLAY)
}

func (w *memoryWriter) WriteQuestion(ctx context.Context, studentID uint, knowledgeID string, videoID uint, segmentID *uint, questionID uint, isCorrect bool, userAnswer string, rightAnswer string, contentIfWrong string) error {
	segID := ""
	if segmentID != nil {
		segID = strconv.FormatUint(uint64(*segmentID), 10)
	}
	ev := &memorypb.Event{
		StudentId:  strconv.FormatUint(uint64(studentID), 10),
		KnowledgeId: knowledgeID,
		VideoId:    strconv.FormatUint(uint64(videoID), 10),
		SegmentId:  segID,
		Ts:         time.Now().Unix(),
		EventType:  memorypb.EventType_QUESTION,
		Detail: &memorypb.Event_Question{
			Question: &memorypb.QuestionEvent{
				QuestionId:  strconv.FormatUint(uint64(questionID), 10),
				IsCorrect:   isCorrect,
				Content:     contentIfWrong,
				UserAnswer:  userAnswer,
				RightAnswer: rightAnswer,
			},
		},
	}
	_, err := w.client.Write(ctx, &memorypb.WriteRequest{Event: ev})
	return err
}

func (w *memoryWriter) writeSimple(ctx context.Context, studentID uint, knowledgeID string, videoID uint, segmentID uint, typ memorypb.EventType) error {
	ev := &memorypb.Event{
		StudentId:  strconv.FormatUint(uint64(studentID), 10),
		KnowledgeId: knowledgeID,
		VideoId:    strconv.FormatUint(uint64(videoID), 10),
		SegmentId:  strconv.FormatUint(uint64(segmentID), 10),
		Ts:         time.Now().Unix(),
		EventType:  typ,
	}
	_, err := w.client.Write(ctx, &memorypb.WriteRequest{Event: ev})
	return err
}

