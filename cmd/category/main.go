package main

import (
	"database/sql"
	"log"
	"yula/internal/config"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"

	"yula/internal/pkg/logging"

	categoryRep "yula/internal/services/category/repository"
	categoryUse "yula/internal/services/category/usecase"

	categoryServer "yula/internal/services/category/server"
	// _ "yula/docs"
)

func getPostgres(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)
	return db
}

// @title Volchock's API
// @version 1.0
// @description Advert placement service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /
func main() {
	logger := logging.GetLogger()

	if err := config.LoadConfig(); err != nil {
		logger.Errorf("error with load config: %s", err.Error())
		return
	}

	// cnfg := config.NewConfig()

	sqlDB := getPostgres(config.Cfg.GetPostgresUrl())
	defer sqlDB.Close()

	catr := categoryRep.NewCategoryRepository(sqlDB)

	catu := categoryUse.NewCategoryUsecase(catr)

	grpcCategory := categoryServer.NewCategoryGRPCServer(logrus.New(), catu)
	err := grpcCategory.NewGRPCServer(config.Cfg.GetCategoryEndPoint())
	if err != nil {
		logger.Errorf("error with grpc server: %s", err.Error())
	}
}
