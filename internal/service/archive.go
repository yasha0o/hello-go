package service

import (
	"context"

	"hello-go/internal/dto"
)

type ArchiveService struct{}

func NewArchiveService() *ArchiveService {
	return &ArchiveService{}
}

func (service *ArchiveService) GetDocuments(
	ctx context.Context,
	pageParams *dto.PageParams,
) (dto.RequestedDocumentsPage, error) {
	return dto.RequestedDocumentsPage{}, nil
}

func (service *ArchiveService) GetDocument(
	ctx context.Context,
	req *dto.DocumentRequest,
) (dto.DocumentDto, error) {
	return dto.DocumentDto{}, nil
}
