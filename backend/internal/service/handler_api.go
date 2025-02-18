package service

import (
	"encoding/json"
	"github.com/nessai1/aiinterview/internal/domain"
	"go.uber.org/zap"
	"net/http"
)

func (s *Service) handleAPIGetInterviewList(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextUserKey).(domain.User)
	if !ok {
		s.logger.Error("User come to API without user in context", zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	interviews, err := s.storage.GetUserInterviewList(r.Context(), user.UUID)
	if err != nil {
		s.logger.Error("Error while load interview list", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	jsoned, err := json.Marshal(&interviews)
	if err != nil {
		s.logger.Error("Error while marshal interviews list", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write(jsoned)
	if err != nil {
		s.logger.Error("Cannot write result to user", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

type CreateInterviewRequest struct {
	Title  string         `json:"title"`
	Timing int            `json:"timing"`
	Topics []domain.Topic `json:"topics"`
}

func (s *Service) handleAPICreateInterview(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextUserKey).(domain.User)
	if !ok {
		s.logger.Error("User come to API without user in context", zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var req CreateInterviewRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.logger.Error("Error while decode request", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	i, err := s.interviewService.CreateInterview(r.Context(), user, req.Title, req.Timing, req.Topics)
	if err != nil {
		s.logger.Error("Error while create interview", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	jsoned, err := json.Marshal(&i)
	if err != nil {
		s.logger.Error("Error while marshal interview", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write(jsoned)
	if err != nil {
		s.logger.Error("Cannot write result to user", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	s.logger.Debug("User create new interview", zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI), zap.String("title", req.Title), zap.Int("timing", req.Timing), zap.Any("topics", req.Topics))
}
