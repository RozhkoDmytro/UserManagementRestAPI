package handlers

/*import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/database"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/repositories"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	"go.uber.org/zap"
)

func InitializeMockVotes(t *testing.T) (*services.MockUserServiceInterface, *votesHandler) {
	os.Setenv("CONFIG_PATH", "../../configs/.sample.env")
	defer os.Unsetenv("CONFIG_PATH")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
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
	votesHandler := NewVotesHandler(userService, logger.Sugar(), cfg)

	return mockUserService, votesHandler
}

func TestVotes(t *testing.T) {
	mockUserService, handler := InitializeMockVotes(t)

	t.Run("Like success", func(t *testing.T) {
		token := auth.GenerateTokenHandler("admin@example.com", "admin", 2, []byte(handler.cfg.JwtKey))
		mockUserService.EXPECT().Vote(gomock.Any(), gomock.Any()).Return("vote_id_123", nil)

		req := httptest.NewRequest(http.MethodPost, "/votes/like/1", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req = req.WithContext(contextWithUserID(req.Context(), 2)) // Mock user ID 2 voting

		w := httptest.NewRecorder()
		handler.Like(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Dislike success", func(t *testing.T) {
		token := auth.GenerateTokenHandler("admin@example.com", "admin", 2, []byte(handler.cfg.JwtKey))
		mockUserService.EXPECT().Vote(gomock.Any(), gomock.Any()).Return("vote_id_123", nil)

		req := httptest.NewRequest(http.MethodPost, "/votes/dislike/3", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		req = req.WithContext(contextWithUserID(req.Context(), 2)) // Mock user ID 2 voting

		w := httptest.NewRecorder()
		handler.Dislike(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Cannot vote for self", func(t *testing.T) {
		token := auth.GenerateTokenHandler("admin@example.com", "admin", 2, []byte(handler.cfg.JwtKey))
		req := httptest.NewRequest(http.MethodPost, "/votes/like/2", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		req = req.WithContext(contextWithUserID(req.Context(), 2)) // Mock user ID 2 voting for profile 2

		w := httptest.NewRecorder()
		handler.Like(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Invalid profile ID", func(t *testing.T) {
		token := auth.GenerateTokenHandler("admin@example.com", "admin", 2, []byte(handler.cfg.JwtKey))
		req := httptest.NewRequest(http.MethodPost, "/votes/like/abc", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		req = req.WithContext(contextWithUserID(req.Context(), 2)) // Mock user ID 2 voting

		w := httptest.NewRecorder()
		handler.Like(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UserService error", func(t *testing.T) {
		mockUserService.EXPECT().Vote(gomock.Any(), gomock.Any()).Return("", errors.New("service error"))
		token := auth.GenerateTokenHandler("admin@example.com", "admin", 2, []byte(handler.cfg.JwtKey))
		req := httptest.NewRequest(http.MethodPost, "/votes/like/1", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req = req.WithContext(contextWithUserID(req.Context(), 2)) // Mock user ID 2 voting

		w := httptest.NewRecorder()
		handler.Like(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func contextWithUserID(ctx context.Context, userID int) context.Context {
	ctx = context.WithValue(ctx, models.IDContextKey, strconv.Itoa(userID))
	return ctx
}
*/
