package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/Knowpals/Knowpals-be-go/api/message"
	"github.com/Knowpals/Knowpals-be-go/events/consumer"
	"github.com/Knowpals/Knowpals-be-go/service/pipeline"
)

// PipelineWorker 消费 result topic，将 Python 产出的结果落库并推进下一阶段 task。
type PipelineWorker struct {
	svc pipeline.PipelineService
}

func NewPipelineWorker(svc pipeline.PipelineService) *PipelineWorker {
	return &PipelineWorker{svc: svc}
}

type pipelineHandler struct {
	svc pipeline.PipelineService
}

func (pipelineHandler) Setup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("kafka 成功建立")
	return nil
}
func (pipelineHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h pipelineHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("start consuming topic:", claim.Topic())
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			var res message.ResultMessage
			if err := json.Unmarshal(msg.Value, &res); err != nil {
				log.Printf("pipeline worker: invalid message: %v", err)
				sess.MarkMessage(msg, "")
				continue
			}
			if err := h.svc.ProcessResult(context.Background(), &res); err != nil {
				log.Printf("pipeline worker: job=%s stage=%s err=%v", res.JobID, res.Stage, err)
			}
			sess.MarkMessage(msg, "")
		case <-sess.Context().Done():
			return nil
		}
	}
}

// Run 阻塞消费 topics，直到 ctx 取消。
func (w *PipelineWorker) Run(ctx context.Context, c consumer.Consumer, topics []string) error {
	h := pipelineHandler{svc: w.svc}
	for {
		if err := c.Consume(ctx, topics, h); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
