package repository

import (
	"context"
	"errors"
	"fmt"

	"hello-go/internal/domain"
	"hello-go/internal/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArchiveRepository struct {
	db *pgxpool.Pool
}

func NewArchiveRepository(db *pgxpool.Pool) *ArchiveRepository {
	return &ArchiveRepository{
		db: db,
	}
}

func (r *ArchiveRepository) InsertRequest(
	ctx context.Context, req *domain.Request,
) (*uuid.UUID, error) {
	return r.InsertRequestTx(ctx, r.db, req)
}

func (r *ArchiveRepository) InsertRequestTx(
	ctx context.Context, querier Querier, req *domain.Request,
) (*uuid.UUID, error) {
	if req == nil {
		return nil, fmt.Errorf("incorrect request")
	}

	query := `INSERT INTO archive.requests (document_id, region, sender) 
	VALUES (@documentID, @region, @sender) RETURNING id`
	args := pgx.NamedArgs{
		"documentID": req.DocumentID,
		"region":     req.Region,
		"sender":     req.Sender,
	}

	var id uuid.UUID
	err := querier.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *ArchiveRepository) UpsertDocument(
	ctx context.Context, req *domain.Document,
) error {
	return r.UpsertDocumentTx(ctx, r.db, req)
}

func (r *ArchiveRepository) UpsertDocumentTx(
	ctx context.Context, tx Querier, req *domain.Document,
) error {
	if req == nil {
		return fmt.Errorf("incorrect document")
	}

	query := `INSERT INTO archive.documents (id, data) 
	VALUES (@id, jsonb_build_array(@data)) ON CONFLICT (id) DO 
	UPDATE SET data = data || jsonb_build_array(@data)`
	args := pgx.NamedArgs{
		"id":   req.ID,
		"data": req.Document,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepository) GetDocument(
	ctx context.Context, id uuid.UUID) (
	*domain.Document, error,
) {
	return r.GetDocumentTx(ctx, r.db, id)
}

func (r *ArchiveRepository) GetDocumentTx(
	ctx context.Context, tx Querier, id uuid.UUID) (
	*domain.Document, error,
) {
	var doc domain.Document
	query := `SELECT id, data, created_at, updated_at 
	FROM archive.documents
	WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := tx.QueryRow(ctx, query, args).Scan(
		&doc.ID,
		&doc.Document,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.NotFoundError{
				Err: err,
			}
		}
		return nil, err
	}

	return &doc, nil
}

func (r *ArchiveRepository) GetDocumentsPage(
	ctx context.Context, pagination *dto.PageParams) (
	[]*domain.Document, int64, error,
) {
	return r.GetDocumentsPageTx(ctx, r.db, pagination)
}

func (r *ArchiveRepository) GetDocumentsPageTx(
	ctx context.Context, tx Querier, pagination *dto.PageParams) (
	[]*domain.Document, int64, error,
) {
	if pagination == nil {
		return nil, 0, fmt.Errorf("incorrect pagination")
	}

	pageQuery := `SELECT id, data, created_at, updated_at, 
	COUNT(*) OVER() as total 
	FROM archive.documents ORDER BY id OFFSET @offset LIMIT @size`
	args := pgx.NamedArgs{
		"offset": (pagination.Page * pagination.Size),
		"size":   pagination.Size,
	}

	rows, err := tx.Query(ctx, pageQuery, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var total int64

	docs := make([]*domain.Document, 0, pagination.Size)
	for rows.Next() {
		doc := new(domain.Document)
		err := rows.Scan(
			&doc.ID,
			&doc.Document,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&total,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("row scan error: %w", err)
		}
		docs = append(docs, doc)
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	return docs, total, nil
}

func (r *ArchiveRepository) WithTx(
	ctx context.Context,
	fn func(Querier) error,
) error {
	return WithTx(ctx, r.db, fn)
}
