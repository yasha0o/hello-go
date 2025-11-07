package mapper

import (
	"hello-go/internal/domain"
	"hello-go/internal/dto"
)

type ArchiveMapper struct{}

func NewArchiveMapper() *ArchiveMapper {
	return &ArchiveMapper{}
}

func (m *ArchiveMapper) MapSlice(doc []*domain.Document) []*dto.DocumentDto {
	if doc == nil {
		return nil
	}

	result := make([]*dto.DocumentDto, len(doc))

	for i, d := range doc {
		result[i] = m.MapDocument(d)
	}
	return result
}

func (m *ArchiveMapper) MapDocument(doc *domain.Document) *dto.DocumentDto {
	if doc == nil {
		return nil
	}

	return &dto.DocumentDto{
		ID:        doc.ID,
		Document:  doc.Document,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

func (m *ArchiveMapper) MapRequest(req *dto.LoadDocumentRequest) *domain.Request {
	if req == nil {
		return nil
	}

	return &domain.Request{
		DocumentID: req.DocumentID,
		Region:     req.Region,
		Sender:     req.Sender,
	}
}
