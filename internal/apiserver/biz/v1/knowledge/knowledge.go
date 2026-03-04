package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/geminik12/krag/internal/apiserver/model"
	"github.com/geminik12/krag/internal/apiserver/store"
	"github.com/geminik12/krag/internal/pkg/contextx"
	"github.com/geminik12/krag/internal/pkg/llm"
	"github.com/geminik12/krag/internal/pkg/storage"
	v1 "github.com/geminik12/krag/pkg/api/apiserver/v1"
	genericoptions "github.com/geminik12/krag/pkg/options"
	"github.com/hibiken/asynq"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

const TaskTypeProcessDocument = "document:process"

type KnowledgeBiz interface {
	Upload(ctx context.Context, req *v1.UploadKnowledgeRequest) (*v1.UploadKnowledgeResponse, error)
	Search(ctx context.Context, req *v1.SearchKnowledgeRequest) (*v1.SearchKnowledgeResponse, error)
	ProcessDocument(ctx context.Context, docID uint) error
}

type knowledgeBiz struct {
	llmClient     llm.Client
	storage       storage.Storage
	store         store.IStore
	qdrantOptions *genericoptions.QdrantOptions
	asynqClient   *asynq.Client
}

func New(llmClient llm.Client, storage storage.Storage, store store.IStore, qdrantOptions *genericoptions.QdrantOptions, asynqClient *asynq.Client) KnowledgeBiz {
	return &knowledgeBiz{
		llmClient:     llmClient,
		storage:       storage,
		store:         store,
		qdrantOptions: qdrantOptions,
		asynqClient:   asynqClient,
	}
}

func (b *knowledgeBiz) Upload(ctx context.Context, req *v1.UploadKnowledgeRequest) (*v1.UploadKnowledgeResponse, error) {
	// 1. Upload to MinIO
	key := fmt.Sprintf("%d_%s", time.Now().Unix(), req.File.Filename)
	src, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	if err := b.storage.Upload(ctx, key, src, req.File.Size, "text/plain"); err != nil {
		return nil, err
	}

	// 2. Create DB record
	doc := &model.Document{
		UserID:   contextx.UserID(ctx),
		Filename: req.File.Filename,
		Key:      key,
		Status:   model.DocStatusPending,
	}
	if err := b.store.DB(ctx).Create(doc).Error; err != nil {
		return nil, err
	}

	// 3. Async Process
	payload, _ := json.Marshal(map[string]interface{}{"document_id": doc.ID})
	task := asynq.NewTask(TaskTypeProcessDocument, payload)
	if _, err := b.asynqClient.Enqueue(task); err != nil {
		return nil, err
	}

	return &v1.UploadKnowledgeResponse{
		DocumentID: fmt.Sprintf("%d", doc.ID),
		Chunks:     0, // Async
	}, nil
}

func (b *knowledgeBiz) ProcessDocument(ctx context.Context, docID uint) error {
	var doc model.Document
	if err := b.store.DB(ctx).First(&doc, docID).Error; err != nil {
		return err
	}

	doc.Status = model.DocStatusProcessing
	b.store.DB(ctx).Save(&doc)

	rc, err := b.storage.Download(ctx, doc.Key)
	if err != nil {
		doc.Status = model.DocStatusFailed
		doc.Error = err.Error()
		b.store.DB(ctx).Save(&doc)
		return err
	}
	defer rc.Close()

	// 1. Load Document
	loader := documentloaders.NewText(rc)
	docs, err := loader.Load(ctx)
	if err != nil {
		doc.Status = model.DocStatusFailed
		doc.Error = err.Error()
		b.store.DB(ctx).Save(&doc)
		return err
	}

	// 2. Split Text
	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkSize = 500
	splitter.ChunkOverlap = 50
	splitDocs, err := textsplitter.SplitDocuments(splitter, docs)
	if err != nil {
		doc.Status = model.DocStatusFailed
		doc.Error = err.Error()
		b.store.DB(ctx).Save(&doc)
		return err
	}

	// Add Metadata
	for i := range splitDocs {
		if splitDocs[i].Metadata == nil {
			splitDocs[i].Metadata = make(map[string]any)
		}
		splitDocs[i].Metadata["user_id"] = doc.UserID
		splitDocs[i].Metadata["doc_id"] = doc.ID
		splitDocs[i].Metadata["filename"] = doc.Filename
	}

	// 3. Get Embedder
	embedder, err := b.llmClient.GetEmbedder("nomic-embed-text")
	if err != nil {
		doc.Status = model.DocStatusFailed
		doc.Error = fmt.Sprintf("Failed to get embedder: %v", err)
		b.store.DB(ctx).Save(&doc)
		return err
	}

	// 4. Store in Qdrant
	store, err := b.newQdrantStore(embedder)
	if err != nil {
		doc.Status = model.DocStatusFailed
		doc.Error = fmt.Sprintf("Failed to connect to Qdrant: %v", err)
		b.store.DB(ctx).Save(&doc)
		return err
	}

	_, err = store.AddDocuments(ctx, splitDocs)
	if err != nil {
		doc.Status = model.DocStatusFailed
		doc.Error = fmt.Sprintf("Failed to add documents: %v", err)
		b.store.DB(ctx).Save(&doc)
		return err
	}

	doc.Status = model.DocStatusSuccess
	b.store.DB(ctx).Save(&doc)
	return nil
}

func (b *knowledgeBiz) Search(ctx context.Context, req *v1.SearchKnowledgeRequest) (*v1.SearchKnowledgeResponse, error) {
	// 1. Get Embedder
	embedder, err := b.llmClient.GetEmbedder("nomic-embed-text")
	if err != nil {
		return nil, err
	}

	// 2. Connect to Qdrant
	store, err := b.newQdrantStore(embedder)
	if err != nil {
		return nil, err
	}

	// 3. Search
	k := req.K
	if k <= 0 {
		k = 3
	}

	// Add UserID filter
	filter := map[string]any{
		"user_id": contextx.UserID(ctx),
	}

	results, err := store.SimilaritySearch(ctx, req.Query, k, vectorstores.WithFilters(filter))
	if err != nil {
		return nil, err
	}

	// 4. Format Response
	var searchResults []v1.SearchResult
	for _, doc := range results {
		searchResults = append(searchResults, v1.SearchResult{
			Content: doc.PageContent,
			Score:   doc.Score,
		})
	}

	return &v1.SearchKnowledgeResponse{
		Results: searchResults,
	}, nil
}

func (b *knowledgeBiz) newQdrantStore(embedder embeddings.Embedder) (vectorstores.VectorStore, error) {
	qdrantURL, err := url.Parse(b.qdrantOptions.URL)
	if err != nil {
		return nil, err
	}

	opts := []qdrant.Option{
		qdrant.WithURL(*qdrantURL),
		qdrant.WithEmbedder(embedder),
		qdrant.WithCollectionName(b.qdrantOptions.CollectionName),
	}
	if b.qdrantOptions.APIKey != "" {
		opts = append(opts, qdrant.WithAPIKey(b.qdrantOptions.APIKey))
	}

	return qdrant.New(opts...)
}
