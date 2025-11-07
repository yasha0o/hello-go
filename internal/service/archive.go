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
	var response *dto.RequestedDocumentsPage

	docs, total, err := s.repository.GetDocumentsPage(ctx, pageParams)
	if err != nil {
		return nil, err
	}

	response = &dto.RequestedDocumentsPage{
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

	var response *dto.CreatedRequest

	action := func(tx repository.Querier) error {
		id, err := s.repository.InsertRequestTx(ctx, tx, request)
		if err != nil {
			return err
		}

		// TODO здесь будет реализован запуск долгой операции
		// по запросам документа в различные архивы

		response = &dto.CreatedRequest{
			ID: *id,
		}
		return nil
	}

	err := s.repository.WithTx(ctx, action)
	if err != nil {
		return nil, err
	}

	return response, nil
}
