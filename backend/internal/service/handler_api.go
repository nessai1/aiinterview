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
