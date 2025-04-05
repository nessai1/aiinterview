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

		isUserExists, err := s.storage.IsUserExists(request.Context(), user.UUID)
		if err != nil {
			s.logger.Error("Cannot check user exists in DB", zap.Error(err))

			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !isUserExists {
			c.MaxAge = -1
			http.SetCookie(writer, c)
			s.logger.Info("User come with nonexists userUUID, clear cookie")

			writer.WriteHeader(http.StatusForbidden)
			_, err = writer.Write([]byte("access denied"))

			return
		}

		request = request.WithContext(context.WithValue(request.Context(), contextUserKey, user))
		s.logger.Debug("authorized request", zap.String("user_uuid", user.UUID), zap.String("ip", request.RemoteAddr), zap.String("uri", request.RequestURI))
		next.ServeHTTP(writer, request)
	})
}

func (s *Service) corsAllowMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		//origin := "http://localhost:5173

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обрабатываем preflight-запросы
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
