package server

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/repositories"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	myValidate "gitlab.com/jkozhemiaka/web-layout/internal/validate"

	"go.uber.org/zap"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/database"
	"gorm.io/gorm"
)

type server struct {
	db          *gorm.DB
	router      Router
	logger      *zap.SugaredLogger
	validate    *validator.Validate
	cfg         *config.Config
	userService services.UserServiceInterface
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHttp(w, r)
}

func (srv *server) initializeRoutes() {
	srv.router.Post("/users", srv.contextExpire(srv.createUserHandler))
	srv.router.Get("/users/{id}", srv.contextExpire(srv.getUser))
	srv.router.Delete("/users/{id}", srv.basicAuth(srv.deleteUser))
	srv.router.Update("/users/{id}", srv.basicAuth(srv.updateUser))
	srv.router.Get("/users", srv.contextExpire(srv.listUsers))
}

func Run() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(apperrors.LoggerInitError.AppendMessage(err))
	}
	defer logger.Sync()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	db, err := database.SetupDatabase(cfg)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	userService := services.NewUserService(repositories.NewUserRepo(db, logger.Sugar()), logger.Sugar())

	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	srvRouter := &router{mux: mux.NewRouter()}
	srv := &server{
		db:          db,
		router:      srvRouter,
		logger:      logger.Sugar(),
		validate:    validate,
		cfg:         cfg,
		userService: userService,
	}
	srv.initializeRoutes()

	logger.Sugar().Infof("Listening HTTP service on %s port", cfg.AppPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), srv)
	if err != nil {
		logger.Sugar().Fatal(err)
	}
}
