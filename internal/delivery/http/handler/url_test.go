package httphandler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vit6556/ozon-internship-assignment/internal/config"
	httpdto "github.com/vit6556/ozon-internship-assignment/internal/delivery/http/dto"
	httphandler "github.com/vit6556/ozon-internship-assignment/internal/delivery/http/handler"
	"github.com/vit6556/ozon-internship-assignment/internal/entity"
	"github.com/vit6556/ozon-internship-assignment/internal/service"
	"github.com/vit6556/ozon-internship-assignment/internal/service/mocks"
)

func TestAddAlias(t *testing.T) {
	mockService := new(mocks.MockUrlService)
	handler := httphandler.NewUrlHandler(mockService, config.AliasConfig{
		Length:  10,
		Charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_",
	})

	e := echo.New()

	tests := []struct {
		name           string
		requestBody    httpdto.AddAliasRequest
		contentType    string
		mockReturn     *entity.Url
		mockError      error
		expectedStatus int
		expectedBody   httpdto.AddAliasResponse
	}{
		{
			name:           "Success - Alias Created",
			requestBody:    httpdto.AddAliasRequest{Url: "https://example.com"},
			contentType:    "application/json",
			mockReturn:     &entity.Url{Alias: "abc123"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   httpdto.AddAliasResponse{Alias: "abc123"},
		},
		{
			name:           "Failure - Invalid Content-Type",
			requestBody:    httpdto.AddAliasRequest{Url: "https://example.com"},
			contentType:    "text/plain",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedBody:   httpdto.AddAliasResponse{Error: "content-type must be application/json"},
		},
		{
			name:           "Failure - Invalid JSON (empty url)",
			requestBody:    httpdto.AddAliasRequest{Url: ""},
			contentType:    "application/json",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   httpdto.AddAliasResponse{Error: "invalid request data"},
		},
		{
			name:           "Failure - URL already exists",
			requestBody:    httpdto.AddAliasRequest{Url: "https://duplicate.com"},
			contentType:    "application/json",
			mockReturn:     nil,
			mockError:      service.ErrSourceUrlAlreadyExists,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   httpdto.AddAliasResponse{Error: "source url already exists"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/alias", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", tt.contentType)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.mockReturn != nil || tt.mockError != nil {
				mockService.On("AddAlias", mock.Anything, tt.requestBody.Url).
					Return(tt.mockReturn, tt.mockError)
			}

			err := handler.AddAlias(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response httpdto.AddAliasResponse
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUrl(t *testing.T) {
	mockService := new(mocks.MockUrlService)
	handler := httphandler.NewUrlHandler(mockService, config.AliasConfig{
		Length:  6,
		Charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_",
	})

	e := echo.New()

	tests := []struct {
		name           string
		alias          string
		mockReturn     *entity.Url
		mockError      error
		expectedStatus int
		expectedBody   httpdto.GetUrlResponse
	}{
		{
			name:           "Success - URL Found",
			alias:          "abc123",
			mockReturn:     &entity.Url{SourceUrl: "https://example.com"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   httpdto.GetUrlResponse{Url: "https://example.com"},
		},
		{
			name:           "Failure - Empty Alias",
			alias:          "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   httpdto.GetUrlResponse{Error: "alias required"},
		},
		{
			name:           "Failure - Invalid Alias",
			alias:          "invalid@alias",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   httpdto.GetUrlResponse{Error: "invalid alias"},
		},
		{
			name:           "Failure - Alias Not Found",
			alias:          "notfou",
			mockReturn:     nil,
			mockError:      service.ErrAliasNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   httpdto.GetUrlResponse{Error: "alias not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/url?alias="+tt.alias, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.mockReturn != nil || tt.mockError != nil {
				mockService.On("GetUrl", mock.Anything, tt.alias).
					Return(tt.mockReturn, tt.mockError)
			}

			err := handler.GetUrl(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response httpdto.GetUrlResponse
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}
