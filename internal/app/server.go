package app

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/vit6556/ozon-internship-assignment/internal/config"
	"github.com/vit6556/ozon-internship-assignment/internal/database"
	inmemoryrepo "github.com/vit6556/ozon-internship-assignment/internal/database/inmemory"
	postgresrepo "github.com/vit6556/ozon-internship-assignment/internal/database/postgres"
	httphandler "github.com/vit6556/ozon-internship-assignment/internal/delivery/http/handler"
	urlservice "github.com/vit6556/ozon-internship-assignment/internal/service/url"
)

func InitServer(cfg *config.HTTPServerConfig) (*echo.Echo, *pgxpool.Pool) {
	var dbPool *pgxpool.Pool
	var databaseRepo database.UrlRepository
	switch cfg.DBType {
	case "postgres":
		dbPool = InitPostgres()
		databaseRepo = postgresrepo.NewUrlRepository(dbPool)
	case "inmemory":
		databaseRepo = inmemoryrepo.NewUrlRepository()
	default:
		log.Fatalf("invalid db_type: %q", cfg.DBType)
	}

	urlService := urlservice.NewUrlService(databaseRepo, cfg.Alias)

	urlHandler := httphandler.NewUrlHandler(urlService, cfg.Alias)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/url", urlHandler.GetUrl)
	e.POST("/url", urlHandler.AddAlias)

	return e, dbPool
}
