package agent

import (
	"context"
	"strconv"

	httpAgent "github.com/Knowpals/Knowpals-be-go/api/http/agent"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
)

func (s *agentService) saveChatPair(ctx context.Context, studentID uint, req httpAgent.ChatReq, reply string) error {
	var vidPtr *uint
	if req.VideoID != "" {
		vid64, err := strconv.ParseUint(req.VideoID, 10, 64)
		if err != nil {
			return err
		}
		vid := uint(vid64)
		vidPtr = &vid
	}
	// user message
	if err := s.chatDao.Save(ctx, model.ChatMessage{
		StudentID:   studentID,
		Role:        "user",
		Text:        req.Text,
		VideoID:     vidPtr,
		KnowledgeID: req.KnowledgeID,
	}); err != nil {
		return err
	}
	// assistant message
	if err := s.chatDao.Save(ctx, model.ChatMessage{
		StudentID:   studentID,
		Role:        "assistant",
		Text:        reply,
		VideoID:     vidPtr,
		KnowledgeID: req.KnowledgeID,
	}); err != nil {
		return err
	}
	return nil
}


