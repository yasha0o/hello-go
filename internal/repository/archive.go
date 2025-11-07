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

func (r *ArchiveRepository) CreateTx(ctx context.Context) (context.Context, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, txKey{}, tx), nil
}

func (r *ArchiveRepository) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		return fmt.Errorf("транзакция не найдена в контексте")
	}
	return tx.Commit(ctx)
}

func (r *ArchiveRepository) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		return fmt.Errorf("транзакция не найдена в контексте")
	}
	return tx.Rollback(ctx)
}

func (r *ArchiveRepository) InsertRequest(
	ctx context.Context, req *domain.Request,
) (*uuid.UUID, error) {
	if req == nil {
		return nil, fmt.Errorf("передан некорректный запрос")
	}

	querier := r.getTx(ctx)
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
	if req == nil {
		return fmt.Errorf("передан некорректный запрос")
	}

	querier := r.getTx(ctx)
	query := `INSERT INTO archive.documents (id, data) 
	VALUES (@id, jsonb_build_array(@data)) ON CONFLICT id DO 
	UPDATE archive.documents 
	SET data = data || jsonb_build_array(@data) 
	WHERE id = @id`
	args := pgx.NamedArgs{
		"id":   req.ID,
		"data": req.Document,
	}

	_, err := querier.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *ArchiveRepository) GetDocument(ctx context.Context, id uuid.UUID) (
	*domain.Document, error,
) {
	querier := r.getTx(ctx)
	var doc domain.Document
	query := `SELECT id, data, created_at, updated_at 
	FROM archive.documents
	WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := querier.QueryRow(ctx, query, args).Scan(
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

func (r *ArchiveRepository) GetDocumentsPage(
	ctx context.Context, pagination *dto.PageParams) (
	[]*domain.Document, int64, error,
) {
	if pagination == nil {
		return nil, 0, fmt.Errorf("переданы некорректные параметры пагинации")
	}

	querier := r.getTx(ctx)
	pageQuery := `SELECT id, data, created_at, updated_at 
	FROM archive.documents ORDER BY id OFFSET @offset LIMIT @size`
	args := pgx.NamedArgs{
		"offset": (pagination.Page * pagination.Size) + 1,
		"size":   pagination.Size,
	}

	countQuery := `SELECT count(*) FROM archive.documents`
	var total int64
	err := querier.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := querier.Query(ctx, pageQuery, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	docs := make([]*domain.Document, 0, pagination.Size)
	for rows.Next() {
		doc := new(domain.Document)
		err := rows.Scan(
			&doc.ID,
			&doc.Document,
			&doc.CreatedAt,
			&doc.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка считывания строки: %w", err)
		}
		docs = append(docs, doc)
	}

	return docs, total, nil
}

func (r *ArchiveRepository) getTx(ctx context.Context) Querier {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if ok {
		return tx
	}
	return r.db
}
