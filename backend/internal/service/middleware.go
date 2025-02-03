package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

const jwtCookieName = "AIINTERVIEW_AUTH"

const contextUserKey = "context_user"

func (s *Service) middlewareTokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie(jwtCookieName)
		if err != nil && errors.Is(http.ErrNoCookie, err) {
			s.logger.Info("Unauthorized access defend", zap.String("ip", request.RemoteAddr), zap.String("uri", request.RequestURI))

			writer.WriteHeader(http.StatusForbidden)
			_, err = writer.Write([]byte("access denied"))
			if err != nil {
				s.logger.Error("Cannot write 'access denied' message to user", zap.Error(err))
			}

			return
		} else if err != nil {
			s.logger.Error("Cannot get jwt cookie of user", zap.Error(err))
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		token := c.Value
		user, err := s.authService.FetchUserFromToken(token)
		if err != nil {
			s.logger.Info("Cannot fetch user from user token cookie", zap.Error(err))
			writer.WriteHeader(http.StatusForbidden)
			_, err = writer.Write([]byte("access denied"))
			if err != nil {
				s.logger.Error("Cannot write 'access denied' message to user", zap.Error(err))
			}

			return
		}

		request.WithContext(context.WithValue(request.Context(), contextUserKey, user))
		s.logger.Debug("authorized request", zap.String("user_uuid", user.UUID), zap.String("ip", request.RemoteAddr), zap.String("uri", request.RequestURI))
		next.ServeHTTP(writer, request)
	})
}
