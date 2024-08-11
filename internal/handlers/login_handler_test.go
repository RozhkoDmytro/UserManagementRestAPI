package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/database"
	"gitlab.com/jkozhemiaka/web-layout/internal/repositories"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	myValidate "gitlab.com/jkozhemiaka/web-layout/internal/validate"
	"go.uber.org/zap"
)

func InitializeMockLogin(t *testing.T) (*services.MockUserServiceInterface, *loginHandler) {
	os.Setenv("CONFIG_PATH", "../../configs/.sample.env")
	defer os.Unsetenv("CONFIG_PATH")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)
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
	validator := validator.New()
	validator.RegisterValidation("password", myValidate.Password)

	loginHandler := NewLoginHandler(userService, logger.Sugar(), cfg)

	return mockUserService, loginHandler
}

func TestLogin(t *testing.T) {
	_, srv := InitializeMockLogin(t)

	t.Run("success", func(t *testing.T) {
		newEmail := "admin@example.com" // Use a fixed email for predictability

		// Encode the form values
		form := url.Values{}
		form.Add("email", newEmail)
		form.Add("password", "securePassword7!")
		reqBody := form.Encode()

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		srv.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
