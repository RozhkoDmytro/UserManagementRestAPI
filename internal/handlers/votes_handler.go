package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	"go.uber.org/zap"
)

type votesHandler struct {
	*BaseHandler
	userService services.UserServiceInterface
	logger      *zap.SugaredLogger
	cfg         *config.Config
}

func NewVotesHandler(userService services.UserServiceInterface, logger *zap.SugaredLogger, cfg *config.Config) *votesHandler {
	return &votesHandler{
		BaseHandler: NewBaseHandler(logger),
		userService: userService,
		logger:      logger,
		cfg:         cfg,
	}
}

func (h *votesHandler) Like(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, 1)
}

func (h *votesHandler) Dislike(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, -1)
}

func (h *votesHandler) vote(w http.ResponseWriter, r *http.Request, value int) {
	type CreateUserResponse struct {
		VoteId string `json:"vote_id"`
	}

	vars := mux.Vars(r)
	profileID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}

	// Getting the ID of the voting user (let's say it is in the context or token)
	userID, err := strconv.Atoi(h.GetAuthenticatedUserID(r.Context()))
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}

	if userID == profileID {
		err := errors.New("you cannot vote for yourself")
		h.sendError(w, err, http.StatusForbidden)
		return
	}

	ctx := r.Context()
	vote := &models.Vote{
		UserID:    uint(userID),
		ProfileID: uint(profileID),
		Value:     value,
		CreatedAt: time.Now(),
	}

	// Attempting to create or update a voice
	voteId, err := h.userService.Vote(ctx, vote)
	if err != nil {
		h.sendError(w, err, http.StatusInternalServerError)
		return
	}

	createUserResponse := &CreateUserResponse{VoteId: strconv.Itoa(int(voteId))}
	h.respond(w, createUserResponse, http.StatusCreated)
}

func (h *votesHandler) RevokeVote(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		VoteId string `json:"vote_id"`
	}

	vars := mux.Vars(r)
	profileID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}

	// Getting the ID of the voting user (let's say it is in the context or token)
	userID, err := strconv.Atoi(h.GetAuthenticatedUserID(r.Context()))
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}

	if userID == profileID {
		err := errors.New("you cannot vote for yourself")
		h.sendError(w, err, http.StatusForbidden)
		return
	}

	ctx := r.Context()

	// Attempting to create or update a voice
	err = h.userService.RevokeVote(ctx, uint(userID), uint(profileID))
	if err != nil {
		h.sendError(w, err, http.StatusInternalServerError)
		return
	}

	createUserResponse := &CreateUserResponse{}
	h.respond(w, createUserResponse, http.StatusCreated)
}
