package server

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/handlers"
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
	validator   *validator.Validate
	cfg         *config.Config
	userService services.UserServiceInterface
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHttp(w, r)
}

func (srv *server) initializeRoutes() {
	userHandler := handlers.NewUserHandler(srv.userService, srv.logger, srv.validator, srv.cfg)
	loginHandler := handlers.NewLoginHandler(srv.userService, srv.logger, srv.cfg)
	votesHandler := handlers.NewVotesHandler(srv.userService, srv.logger, srv.cfg)

	srv.router.Post("/users", srv.contextExpire(userHandler.CreateUserHandler))
	srv.router.Get("/users/{id:[0-9]+}", srv.contextExpire(userHandler.GetUser))
	srv.router.Delete("/users/{id:[0-9]+}", srv.jwtMiddleware(userHandler.DeleteUser))
	srv.router.Update("/users/{id:[0-9]+}", srv.jwtMiddleware(userHandler.UpdateUser))
	srv.router.Get("/users", srv.contextExpire(userHandler.ListUsers))
	srv.router.Get("/users/count", srv.contextExpire(userHandler.CountUsers))

	srv.router.Post("/login", srv.contextExpire(loginHandler.Login))

	srv.router.Post("/like/{id:[0-9]+}", srv.jwtMiddleware(votesHandler.Like))
	srv.router.Post("/dislike/{id:[0-9]+}", srv.jwtMiddleware(votesHandler.DisLike))
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
		validator:   validate,
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
