package worker

import (
	"context"
	"encoding/json"
	"fmt"

	knowledgev1 "github.com/geminik12/krag/internal/apiserver/biz/v1/knowledge"
	"github.com/hibiken/asynq"
)

type Worker struct {
	server       *asynq.Server
	knowledgeBiz knowledgev1.KnowledgeBiz
}

func NewWorker(redisAddr string, password string, db int, knowledgeBiz knowledgev1.KnowledgeBiz) *Worker {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: password, DB: db},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &Worker{
		server:       server,
		knowledgeBiz: knowledgeBiz,
	}
}

func (w *Worker) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(knowledgev1.TaskTypeProcessDocument, w.HandleDocumentProcess)
	return w.server.Run(mux)
}

func (w *Worker) HandleDocumentProcess(ctx context.Context, t *asynq.Task) error {
	var p map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	docIDFloat, ok := p["document_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid document_id type")
	}
	return w.knowledgeBiz.ProcessDocument(ctx, uint(docIDFloat))
}
