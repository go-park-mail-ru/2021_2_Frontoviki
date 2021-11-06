package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"yula/internal/config"
	"yula/internal/database"

	_ "github.com/jackc/pgx/stdlib"

	imageloaderRepo "yula/internal/pkg/image_loader/repository"
	imageloaderUse "yula/internal/pkg/image_loader/usecase"
	"yula/internal/pkg/logging"
	userHttp "yula/internal/pkg/user/delivery/http"
	userRep "yula/internal/pkg/user/repository"
	userUse "yula/internal/pkg/user/usecase"

	"yula/internal/pkg/middleware"
	sessHttp "yula/internal/pkg/session/delivery/http"
	sessRep "yula/internal/pkg/session/repository"
	sessUse "yula/internal/pkg/session/usecase"

	advtHttp "yula/internal/pkg/advt/delivery/http"
	advtRep "yula/internal/pkg/advt/repository"
	advtUse "yula/internal/pkg/advt/usecase"

	cartHttp "yula/internal/pkg/cart/delivery/http"
	cartRep "yula/internal/pkg/cart/repository"
	cartUse "yula/internal/pkg/cart/usecase"

	categoryHttp "yula/internal/pkg/category/delivery/http"
	categoryRep "yula/internal/pkg/category/repository"
	categoryUse "yula/internal/pkg/category/usecase"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "yula/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	govalidator.SetFieldsRequiredByDefault(true)
}

func getPostgres() *sql.DB {
	dsn := "user=postgres dbname=yula password=password host=127.0.0.1 port=5432 sslmode=disable"
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

	cnfg := config.NewConfig()
	postgres, err := database.NewPostgres(cnfg.DbConfig.DatabaseUrl)
	if err != nil {
		logger.Fatalf("db error instance", err.Error())
		return
	}
	defer postgres.Close()

	sqlDB := getPostgres()

	r := mux.NewRouter()

	r.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)

	api := r.PathPrefix("").Subrouter()

	api.Use(middleware.CorsMiddleware)
	api.Use(middleware.ContentTypeMiddleware)
	api.Use(middleware.LoggerMiddleware)
	//api.Use(middleware.CSRFMiddleWare())

	ilr := imageloaderRepo.NewImageLoaderRepository()
	ar := advtRep.NewAdvtRepository(postgres.GetDbPool())
	ur := userRep.NewUserRepository(postgres.GetDbPool())
	rr := userRep.NewRatingRepository(sqlDB)
	sr := sessRep.NewSessionRepository(&cnfg.TarantoolCfg)
	cr := cartRep.NewCartRepository(postgres.GetDbPool())
	catr := categoryRep.NewCategoryRepository(sqlDB)

	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)
	au := advtUse.NewAdvtUsecase(ar, ilu)
	uu := userUse.NewUserUsecase(ur, rr, ilu)
	su := sessUse.NewSessionUsecase(sr)
	cu := cartUse.NewCartUsecase(cr)
	catu := categoryUse.NewCategoryUsecase(catr)

	ah := advtHttp.NewAdvertHandler(au, uu)
	uh := userHttp.NewUserHandler(uu, su)
	sh := sessHttp.NewSessionHandler(su, uu)
	ch := cartHttp.NewCartHandler(cu, uu, au)
	cath := categoryHttp.NewCategoryHandler(catu)

	sm := middleware.NewSessionMiddleware(su)

	ah.Routing(api, sm)
	uh.Routing(api, sm)
	sh.Routing(api)
	ch.Routing(api, sm)
	cath.Routing(api)
	middleware.Routing(api)

	//http
	fmt.Println("start serving ::8080")
	error := http.ListenAndServe(":8080", r)

	// //https
	// fmt.Println("start serving ::5000")
	// error := http.ListenAndServeTLS(":5000", "certificate.crt", "key.key", r)

	logger.Errorf("http serve error %v", error)
}
