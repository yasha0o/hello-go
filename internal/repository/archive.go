package repository

import (
	"context"
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
	if req == nil {
		return nil, fmt.Errorf("Передан некорректный запрос")
	}

	tx := r.tx(ctx)

	var id *uuid.UUID

	query := `INSERT INTO archive.requests (document_id, region, sender) 
	VALUES (@documentID, @region, @sender) RETURNING id`
	args := pgx.NamedArgs{
		"documentID": req.DocumentID,
		"region":     req.Region,
		"sender":     req.Sender,
	}

	err := tx.QueryRow(ctx, query, args).Scan(id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (r *ArchiveRepository) UpsertDocument(
	ctx context.Context, req *domain.Document,
) error {
	if req == nil {
		return fmt.Errorf("Передан некорректный запрос")
	}

	tx := r.tx(ctx)

	query := `INSERT INTO archive.documents (id, data) 
	VALUES (@id, jsonb_build_array(@data)) ON CONFLICT id DO 
	UPDATE archive.documents 
	SET data = data || jsonb_build_array(@data) 
	WHERE id = @id`
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

func (r *ArchiveRepository) GetDocument(ctx context.Context, id *uuid.UUID) (
	*domain.Document, error,
) {
	if id == nil {
		return nil, fmt.Errorf("Передан некорректный id")
	}

	tx := r.tx(ctx)

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
		return nil, err
	}

	return &doc, nil
}

func (r *ArchiveRepository) GetDocuments(
	ctx context.Context, pagination *dto.PageParams) (
	[]*domain.Document, int64, error,
) {
	if pagination == nil {
		return nil, 0, fmt.Errorf("Переданы некорректные параметры пагинации")
	}

	tx := r.tx(ctx)

	pageQuery := `SELECT id, data, created_at, updated_at 
	FROM archive.documentsORDER BY id OFFSET @offset LIMIT @size`
	args := pgx.NamedArgs{
		"offset": (pagination.Page * pagination.Size) + 1,
		"size":   pagination.Size,
	}

	countQuery := `SELECT count(*) FROM archive.documents`
	var total int64
	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := tx.Query(ctx, pageQuery, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	docs := make([]*domain.Document, 0)
	for rows.Next() {
		var doc domain.Document
		err := rows.Scan(
			&doc.ID,
			&doc.Document,
			&doc.CreatedAt,
			&doc.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("Ошибка считывания строки: %w", err)
		}
		docs = append(docs, &doc)
	}

	return docs, total, nil
}

func (r *ArchiveRepository) tx(ctx context.Context) *pgxpool.Pool {
	tx := ctx.Value("tx").(*pgxpool.Pool)
	if tx == nil {
		tx = r.db
	}
	return tx
}
