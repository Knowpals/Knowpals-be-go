package cache

import "context"

type KnowledgeCache interface {
	SetKnowledgeIDToKnowledge(ctx context.Context, knowledgeID string)
}
