package httphandler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/vit6556/ozon-internship-assignment/internal/config"
	httpdto "github.com/vit6556/ozon-internship-assignment/internal/delivery/http/dto"
	"github.com/vit6556/ozon-internship-assignment/internal/service"
)

type UrlHandler struct {
	urlService   service.UrlService
	validate     *validator.Validate
	aliasLength  int
	aliasCharset map[byte]struct{}
}

func NewUrlHandler(urlService service.UrlService, aliasConfig config.AliasConfig) *UrlHandler {
	aliasCharset := make(map[byte]struct{}, len(aliasConfig.Charset))
	for i := 0; i < len(aliasConfig.Charset); i++ {
		aliasCharset[aliasConfig.Charset[i]] = struct{}{}
	}

	return &UrlHandler{
		urlService:   urlService,
		validate:     validator.New(),
		aliasLength:  aliasConfig.Length,
		aliasCharset: aliasCharset,
	}
}

func (h *UrlHandler) isValidAlias(alias string) bool {
	if len(alias) != h.aliasLength {
		return false
	}

	for i := 0; i < len(alias); i++ {
		if _, exists := h.aliasCharset[alias[i]]; !exists {
			return false
		}
	}

	return true

}

func (h *UrlHandler) AddAlias(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return c.JSON(http.StatusUnsupportedMediaType, httpdto.AddAliasResponse{Error: "content-type must be application/json"})
	}

	var request httpdto.AddAliasRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpdto.AddAliasResponse{Error: "invalid json data"})
	}

	if err := h.validate.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, httpdto.AddAliasResponse{Error: "invalid request data"})
	}

	url, err := h.urlService.AddAlias(c.Request().Context(), request.Url)
	if err != nil {
		switch err {
		case service.ErrSourceUrlAlreadyExists:
			return c.JSON(http.StatusBadRequest, httpdto.AddAliasResponse{Error: "source url already exists"})
		default:
			return c.JSON(http.StatusInternalServerError, httpdto.AddAliasResponse{Error: "failed to create alias"})
		}
	}

	return c.JSON(http.StatusOK, httpdto.AddAliasResponse{Alias: url.Alias})
}

func (h *UrlHandler) GetUrl(c echo.Context) error {
	alias := c.QueryParam("alias")
	if alias == "" {
		return c.JSON(http.StatusBadRequest, httpdto.GetUrlResponse{Error: "alias required"})
	} else if !h.isValidAlias(alias) {
		return c.JSON(http.StatusBadRequest, httpdto.GetUrlResponse{Error: "invalid alias"})
	}

	url, err := h.urlService.GetUrl(c.Request().Context(), alias)
	if err != nil {
		switch err {
		case service.ErrAliasNotFound:
			return c.JSON(http.StatusNotFound, httpdto.GetUrlResponse{Error: "alias not found"})
		default:
			return c.JSON(http.StatusInternalServerError, httpdto.GetUrlResponse{Error: "failed to get url by alias"})
		}
	}

	return c.JSON(http.StatusOK, httpdto.GetUrlResponse{Url: url.SourceUrl})
}
