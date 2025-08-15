package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "segment_service/docs"
	"segment_service/internal/repository/postgresql"
	"segment_service/internal/rest"
	"segment_service/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			User segment service
//	@version		1.0
//	@description	This is a educational project for vk

//	@contact.name	Creator
//	@contact.email	demid.shumakher@gmail.com

// @host		localhost:8080
// @BasePath	/
func main() {
	e := echo.New()

	// Читаем конфиг подключения к БД
	dbUser := getenv("DB_USER", "postgres")
	dbPass := getenv("DB_PASSWORD", "postgres")
	dbHost := getenv("DB_HOST", "localhost")
	dbPort := getenv("DB_PORT", "5432")
	dbName := getenv("DB_NAME", "segment_service")
	sslMode := getenv("DB_SSLMODE", "disable")
	appPort := getenv("APP_PORT", "8080")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, sslMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}

	// Инициализируем репозитории и сервисы
	userRepo := postgresql.NewUserRepository(db)
	segmentRepo := postgresql.NewSegmentRepository(db)

	rest.NewUserHandler(e, service.NewUserService(userRepo))
	rest.NewSegmentHandler(e, service.NewSegmentService(segmentRepo))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":" + appPort))
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
