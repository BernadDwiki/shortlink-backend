package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BernadDwiki/shortlink-backend/internal/model"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) CreateLink(ctx context.Context, userID int, originalURL, slug string) (*model.Link, error) {
	var link model.Link
	query := `INSERT INTO links (user_id, original_url, slug) VALUES ($1, $2, $3) RETURNING id, user_id, original_url, slug, created_at, deleted_at`

	err := r.db.QueryRowContext(ctx, query, userID, originalURL, slug).Scan(
		&link.ID,
		&link.UserID,
		&link.OriginalURL,
		&link.Slug,
		&link.CreatedAt,
		&link.DeletedAt,
	)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"links_slug_key\"" {
			return nil, errors.New("slug already taken")
		}
		return nil, err
	}

	return &link, nil
}

func (r *LinkRepository) GetLinkBySlug(ctx context.Context, slug string) (*model.Link, error) {
	var link model.Link
	query := `SELECT id, user_id, original_url, slug, created_at, deleted_at FROM links WHERE slug = $1 AND deleted_at IS NULL LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&link.ID,
		&link.UserID,
		&link.OriginalURL,
		&link.Slug,
		&link.CreatedAt,
		&link.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &link, nil
}

func (r *LinkRepository) GetLinksByUserID(ctx context.Context, userID int) ([]model.Link, error) {
	query := `SELECT id, user_id, original_url, slug, created_at, deleted_at FROM links WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []model.Link
	for rows.Next() {
		var link model.Link
		if err := rows.Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.Slug, &link.CreatedAt, &link.DeletedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, rows.Err()
}

func (r *LinkRepository) DeleteLink(ctx context.Context, linkID, userID int) error {
	query := `UPDATE links SET deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, linkID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("link not found or unauthorized")
	}

	return nil
}
