package service

import (
	"go.uber.org/zap"
	"net/http"
)

func (s *Service) handlePublicAuthorize(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		s.logger.Debug("Client send empty token", zap.String("req", r.RequestURI))
	}

}
