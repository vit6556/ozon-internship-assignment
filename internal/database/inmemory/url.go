package inmemoryrepo

import (
	"context"
	"log"
	"sync"

	"github.com/vit6556/ozon-internship-assignment/internal/database"
	"github.com/vit6556/ozon-internship-assignment/internal/entity"
)

type UrlRepository struct {
	mutex    sync.RWMutex
	aliasUrl map[string]*entity.Url
	urlAlias map[string]string
}

func NewUrlRepository() *UrlRepository {
	return &UrlRepository{
		aliasUrl: make(map[string]*entity.Url),
		urlAlias: make(map[string]string),
	}
}

func (r *UrlRepository) GetUrl(ctx context.Context, alias string) (*entity.Url, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	url, exists := r.aliasUrl[alias]
	if !exists {
		log.Printf("failed to get source url by alias %q: not found", alias)
		return nil, database.ErrAliasNotFound
	}

	return url, nil
}

func (r *UrlRepository) AddAlias(ctx context.Context, sourceUrl string, alias string) (*entity.Url, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.urlAlias[sourceUrl]; exists {
		return nil, database.ErrSourceUrlAlreadyExists
	}

	if _, exists := r.aliasUrl[alias]; exists {
		return nil, database.ErrAliasAlreadyExists
	}

	url := &entity.Url{
		ID:        len(r.aliasUrl) + 1,
		SourceUrl: sourceUrl,
		Alias:     alias,
	}

	r.aliasUrl[alias] = url
	r.urlAlias[sourceUrl] = alias

	return url, nil
}
