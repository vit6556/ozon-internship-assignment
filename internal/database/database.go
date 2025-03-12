package database

import (
	"context"
	"errors"

	"github.com/vit6556/ozon-internship-assignment/internal/entity"
)

var (
	ErrAliasNotFound          = errors.New("alias not found")
	ErrAliasAlreadyExists     = errors.New("alias already exists")
	ErrSourceUrlAlreadyExists = errors.New("source url already exists")
	ErrAliasCreationFailed    = errors.New("failed to create url alias")
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=UrlRepository --output=mocks --filename=url.go --structname=MockUrlRepository
type UrlRepository interface {
	GetUrl(ctx context.Context, alias string) (*entity.Url, error)
	AddAlias(ctx context.Context, sourceUrl string, alias string) (*entity.Url, error)
}
