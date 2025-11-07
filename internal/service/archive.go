package service

import (
	"context"

	"hello-go/internal/dto"
	"hello-go/internal/mapper"
	"hello-go/internal/repository"

	"github.com/google/uuid"
)

type ArchiveService struct {
	repository *repository.ArchiveRepository
	mapper     *mapper.ArchiveMapper
}

func NewArchiveService(r *repository.ArchiveRepository,
	m *mapper.ArchiveMapper,
) *ArchiveService {
	return &ArchiveService{
		repository: r,
		mapper:     m,
	}
}

func (s *ArchiveService) GetDocumentsPage(
	ctx context.Context,
	pageParams *dto.PageParams,
) (*dto.RequestedDocumentsPage, error) {
	docs, total, err := s.repository.GetDocumentsPage(ctx, pageParams)
	if err != nil {
		return nil, err
	}

	response := &dto.RequestedDocumentsPage{
		Data:  s.mapper.MapSlice(docs),
		Total: total,
		Page:  pageParams.Page,
		Size:  pageParams.Size,
	}

	return response, nil
}

func (s *ArchiveService) GetDocument(
	ctx context.Context,
	id uuid.UUID,
) (*dto.DocumentDto, error) {
	doc, err := s.repository.GetDocument(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDocument(doc), nil
}

func (s *ArchiveService) LoadDocument(
	ctx context.Context,
	req *dto.LoadDocumentRequest,
) (*dto.CreatedRequest, error) {
	request := s.mapper.MapRequest(req)

	ctx, err := s.repository.CreateTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = s.repository.RollbackTx(ctx)
	}()

	id, err := s.repository.InsertRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// TODO здесь будет реализован запуск долгой операции
	// по запросам документа в различные архивы

	response := &dto.CreatedRequest{
		ID: *id,
	}

	if err := s.repository.CommitTx(ctx); err != nil {
		return nil, err
	}

	return response, err
}
