package urlservice

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/vit6556/ozon-internship-assignment/internal/config"
	"github.com/vit6556/ozon-internship-assignment/internal/database"
	"github.com/vit6556/ozon-internship-assignment/internal/entity"
	"github.com/vit6556/ozon-internship-assignment/internal/service"
)

type UrlService struct {
	urlRepo     database.UrlRepository
	aliasConfig config.AliasConfig
}

func NewUrlService(urlRepo database.UrlRepository, aliasConfig config.AliasConfig) *UrlService {
	return &UrlService{
		urlRepo:     urlRepo,
		aliasConfig: aliasConfig,
	}
}

func (s *UrlService) generateAlias() string {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	alias := make([]byte, s.aliasConfig.Length)
	for i := range alias {
		alias[i] = s.aliasConfig.Charset[rng.Intn(len(s.aliasConfig.Charset))]
	}

	return string(alias)
}

func (s *UrlService) GetUrl(ctx context.Context, alias string) (*entity.Url, error) {
	url, err := s.urlRepo.GetUrl(ctx, alias)
	if err != nil {
		log.Printf("failed to get source url by alias %q: %v", alias, err)
		return nil, service.ErrAliasNotFound
	}

	return url, nil
}

func (s *UrlService) AddAlias(ctx context.Context, sourceUrl string) (*entity.Url, error) {
	var url *entity.Url
	var err error

	for i := 0; i < s.aliasConfig.GenerationRetries; i++ {
		alias := s.generateAlias()
		url, err = s.urlRepo.AddAlias(ctx, sourceUrl, alias)

		switch err {
		case nil:
			return url, nil
		case database.ErrSourceUrlAlreadyExists:
			return nil, service.ErrSourceUrlAlreadyExists
		case database.ErrAliasAlreadyExists:
			log.Printf("alias %q already exists, retrying (%d/%d)...", alias, i+1, s.aliasConfig.GenerationRetries)
		default:
			log.Printf("failed to create alias for url %q: %v", sourceUrl, err)
			return nil, service.ErrAliasCreationFailed
		}
	}

	log.Printf("failed to create alias for url %q: all %d retries exhausted", sourceUrl, s.aliasConfig.GenerationRetries)
	return nil, service.ErrAliasCreationFailed
}
