package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"yula/internal/config"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"

	imageloaderRepo "yula/internal/pkg/image_loader/repository"
	imageloaderUse "yula/internal/pkg/image_loader/usecase"
	"yula/internal/pkg/logging"
	userHttp "yula/internal/pkg/user/delivery/http"
	userRep "yula/internal/pkg/user/repository"
	userUse "yula/internal/pkg/user/usecase"

	"yula/internal/pkg/middleware"
	sessHttp "yula/internal/pkg/session/delivery/http"
	sessRep "yula/services/auth/repository"
	sessUse "yula/services/auth/usecase"

	advtHttp "yula/internal/pkg/advt/delivery/http"
	advtRep "yula/internal/pkg/advt/repository"
	advtUse "yula/internal/pkg/advt/usecase"

	cartHttp "yula/internal/pkg/cart/delivery/http"
	cartRep "yula/internal/pkg/cart/repository"
	cartUse "yula/internal/pkg/cart/usecase"

	srchHttp "yula/internal/pkg/search/delivery/http"
	srchRep "yula/internal/pkg/search/repository"
	srchUse "yula/internal/pkg/search/usecase"

	categoryHttp "yula/internal/pkg/category/delivery/http"
	categoryRep "yula/internal/pkg/category/repository"
	categoryUse "yula/internal/pkg/category/usecase"

	chatHttp "yula/internal/pkg/chat/delivery"
	metrics "yula/internal/pkg/metrics"
	metricsHttp "yula/internal/pkg/metrics/delivery"
	chatRep "yula/services/chat/repository"
	chatUse "yula/services/chat/usecase"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	authProto "yula/proto/generated/auth"
	chatProto "yula/proto/generated/chat"

	authServer "yula/services/auth/delivery"
	chatServer "yula/services/chat/delivery"

	"google.golang.org/grpc"

	// _ "yula/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	govalidator.SetFieldsRequiredByDefault(true)
}

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

	r := mux.NewRouter()

	r.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)

	api := r.PathPrefix("").Subrouter()

	// ставим мидлварину с метриками
	m := metrics.NewMetrics(r)
	mmw := metricsHttp.NewMetricsMiddleware(m)
	r.Use(mmw.ScanMetrics)

	api.Use(middleware.CorsMiddleware)
	api.Use(middleware.ContentTypeMiddleware)
	api.Use(middleware.LoggerMiddleware)
	//api.Use(middleware.CSRFMiddleWare())
	ilr := imageloaderRepo.NewImageLoaderRepository()
	ar := advtRep.NewAdvtRepository(sqlDB)
	ur := userRep.NewUserRepository(sqlDB)
	rr := userRep.NewRatingRepository(sqlDB)
	sr := sessRep.NewSessionRepository(config.Cfg.GetTarantoolCfg())
	cr := cartRep.NewCartRepository(sqlDB)
	serr := srchRep.NewSearchRepository(sqlDB)
	catr := categoryRep.NewCategoryRepository(sqlDB)
	chr := chatRep.NewChatRepository(sqlDB)

	ilu := imageloaderUse.NewImageLoaderUsecase(ilr)
	au := advtUse.NewAdvtUsecase(ar, ilu)
	uu := userUse.NewUserUsecase(ur, rr, ilu)
	su := sessUse.NewSessionUsecase(sr)
	cu := cartUse.NewCartUsecase(cr)
	seru := srchUse.NewSearchUsecase(serr, ar)
	catu := categoryUse.NewCategoryUsecase(catr)
	chu := chatUse.NewChatUsecase(chr)

	ah := advtHttp.NewAdvertHandler(au, uu)
	uh := userHttp.NewUserHandler(uu, su)

	grpcAuthClient, err := grpc.Dial(
		"127.0.0.1:8180",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("cant open grpc conn")
	}
	defer grpcAuthClient.Close()

	grpcChatClient, err := grpc.Dial(
		"127.0.0.1:8280",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("cant open grpc conn")
	}
	defer grpcChatClient.Close()

	sh := sessHttp.NewSessionHandler(authProto.NewAuthClient(grpcAuthClient), uu)
	ch := cartHttp.NewCartHandler(cu, uu, au)
	serh := srchHttp.NewSearchHandler(seru)
	cath := categoryHttp.NewCategoryHandler(catu)
	chth := chatHttp.NewChatHandler(chatProto.NewChatClient(grpcChatClient))

	sm := middleware.NewSessionMiddleware(su)

	ah.Routing(api, sm)
	uh.Routing(api, sm)
	sh.Routing(api)
	ch.Routing(api, sm)
	serh.Routing(api)
	cath.Routing(api)
	middleware.Routing(api)
	chth.Routing(api, sm)

	grpcAuth := authServer.NewAuthGRPCServer(logrus.New(), su)
	go grpcAuth.NewGRPCServer("127.0.0.1:8180")

	grpcChat := chatServer.NewChatGRPCServer(logrus.New(), chu)
	go grpcChat.NewGRPCServer("127.0.0.1:8280")

	port := config.Cfg.GetMainPort()
	fmt.Printf("start serving ::%s\n", port)

	var error error
	secure := config.Cfg.IsSecure()
	if secure {
		error = http.ListenAndServeTLS(fmt.Sprintf(":%s", port), "certificate.crt", "key.key", r)
	} else {
		error = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	}

	logger.Errorf("http serve error %v", error)
}
