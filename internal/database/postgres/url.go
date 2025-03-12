package postgresrepo

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vit6556/ozon-internship-assignment/internal/database"
	"github.com/vit6556/ozon-internship-assignment/internal/entity"
)

type UrlRepository struct {
	db *pgxpool.Pool
}

func NewUrlRepository(db *pgxpool.Pool) *UrlRepository {
	return &UrlRepository{
		db: db,
	}
}

func (r *UrlRepository) GetUrl(ctx context.Context, alias string) (*entity.Url, error) {
	var url entity.Url
	err := r.db.QueryRow(ctx, "SELECT id, source_url, alias FROM urls WHERE alias = $1", alias).
		Scan(&url.ID, &url.SourceUrl, &url.Alias)
	if err != nil {
		log.Printf("failed to get source url by alias %q: %v", alias, err)
		return nil, database.ErrAliasNotFound
	}

	return &url, nil
}

func (r *UrlRepository) AddAlias(ctx context.Context, sourceUrl string, alias string) (*entity.Url, error) {
	var url entity.Url
	err := r.db.QueryRow(ctx, "INSERT INTO urls (source_url, alias) VALUES ($1, $2) RETURNING id, source_url, alias",
		sourceUrl, alias).Scan(&url.ID, &url.SourceUrl, &url.Alias)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "urls_source_url_key":
				return nil, database.ErrSourceUrlAlreadyExists
			case "urls_alias_key":
				return nil, database.ErrAliasAlreadyExists
			}
		}

		log.Printf("failed to create alias %q for url %q: %v", alias, sourceUrl, err)
		return nil, database.ErrAliasCreationFailed
	}

	return &url, nil
}
