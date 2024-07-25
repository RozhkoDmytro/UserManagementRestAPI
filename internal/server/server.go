package server

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/database"
	"gorm.io/gorm"
)

type server struct {
	db     *gorm.DB
	router Router
	logger *zap.SugaredLogger
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHttp(w, r)
}

func (srv *server) initializeRoutes() {
	srv.router.Post("/user", srv.contextExpire(srv.createUserHandler()))
	srv.router.Get("/user/{id}", srv.contextExpire(srv.getUser()))
	srv.router.Delete("/user/{id}", srv.contextExpire(srv.deleteUser()))
	srv.router.Update("/user/{id}", srv.contextExpire(srv.updateUser()))
	srv.router.Get("/users", srv.contextExpire(srv.listUsers()))
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

	srvRouter := &router{mux: mux.NewRouter()}
	srv := &server{
		db:     db,
		router: srvRouter,
		logger: logger.Sugar(),
	}
	srv.initializeRoutes()

	logger.Sugar().Infof("Listening HTTP service on %s port", cfg.AppPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), srv)
	if err != nil {
		logger.Sugar().Fatal(err)
	}
}
