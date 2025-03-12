package urlservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/vit6556/ozon-internship-assignment/internal/config"
	"github.com/vit6556/ozon-internship-assignment/internal/database"
	"github.com/vit6556/ozon-internship-assignment/internal/database/mocks"
	"github.com/vit6556/ozon-internship-assignment/internal/entity"
	"github.com/vit6556/ozon-internship-assignment/internal/service"
	urlservice "github.com/vit6556/ozon-internship-assignment/internal/service/url"
)

func TestGetUrl(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	urlService := urlservice.NewUrlService(mockRepo, config.AliasConfig{})

	ctx := context.Background()

	tests := []struct {
		name        string
		alias       string
		mockReturn  *entity.Url
		mockError   error
		expectedURL *entity.Url
		expectedErr error
	}{
		{
			name:       "Success - valid alias",
			alias:      "abc123",
			mockReturn: &entity.Url{ID: 1, SourceUrl: "https://example.com", Alias: "abc123"},
			mockError:  nil,
			expectedURL: &entity.Url{
				ID:        1,
				SourceUrl: "https://example.com",
				Alias:     "abc123",
			},
			expectedErr: nil,
		},
		{
			name:        "Failure - alias not found",
			alias:       "notfound",
			mockReturn:  nil,
			mockError:   database.ErrAliasNotFound,
			expectedURL: nil,
			expectedErr: service.ErrAliasNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetUrl", ctx, tt.alias).Return(tt.mockReturn, tt.mockError)

			result, err := urlService.GetUrl(ctx, tt.alias)

			assert.Equal(t, tt.expectedURL, result)
			assert.Equal(t, tt.expectedErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAddAlias(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	urlService := urlservice.NewUrlService(mockRepo, config.AliasConfig{
		Length:            10,
		Charset:           "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_",
		GenerationRetries: 5,
	})

	ctx := context.Background()

	tests := []struct {
		name        string
		sourceUrl   string
		mockReturn  *entity.Url
		mockError   error
		expectedURL *entity.Url
		expectedErr error
	}{
		{
			name:      "Success - Unique alias generated",
			sourceUrl: "https://example.com",
			mockReturn: &entity.Url{
				ID:        1,
				SourceUrl: "https://example.com",
				Alias:     "abc123",
			},
			mockError:   nil,
			expectedURL: &entity.Url{ID: 1, SourceUrl: "https://example.com", Alias: "abc123"},
			expectedErr: nil,
		},
		{
			name:        "Failure - Source URL already exists",
			sourceUrl:   "https://duplicate.com",
			mockReturn:  nil,
			mockError:   database.ErrSourceUrlAlreadyExists,
			expectedURL: nil,
			expectedErr: service.ErrSourceUrlAlreadyExists,
		},
		{
			name:        "Failure - Alias generation exhausted",
			sourceUrl:   "https://collision.com",
			mockReturn:  nil,
			mockError:   database.ErrAliasAlreadyExists,
			expectedURL: nil,
			expectedErr: service.ErrAliasCreationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("AddAlias", ctx, tt.sourceUrl, testifyMock.Anything).Return(tt.mockReturn, tt.mockError)

			result, err := urlService.AddAlias(ctx, tt.sourceUrl)

			assert.Equal(t, tt.expectedURL, result)
			assert.Equal(t, tt.expectedErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
