package service

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nessai1/aiinterview/internal/domain"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (s *Service) handlePublicAuthorize(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	invitationCode := vars["invitation"]
	if invitationCode != s.config.InvitationCode {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("invalid invitation code"))
		if err != nil {
			s.logger.Error("Error while write answer on authorize method", zap.Error(err))
		}
		return
	}

	jwtCookie, err := r.Cookie(jwtCookieName)
	if err != nil && !errors.Is(http.ErrNoCookie, err) {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("Cannot access to user cookies", zap.Error(err))
		return
	}

	var authorizeReason string

	createNewUserAndSetCookie := func() (domain.User, error) {
		user, err := s.storage.RegisterUser(r.Context())
		if err != nil {
			return domain.User{}, fmt.Errorf("cannot register new user in storage: %w", err)
		}

		token, err := s.authService.BuildTokenByUser(user)
		if err != nil {
			return domain.User{}, fmt.Errorf("error while create token for new user: %w", err)
		}

		sameSitePolicy := http.SameSiteLaxMode
		if s.config.IsDev {
			sameSitePolicy = http.SameSiteNoneMode
		}

		http.SetCookie(w, &http.Cookie{
			Name:     jwtCookieName,
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(tokenExp),
			HttpOnly: true,
			SameSite: sameSitePolicy,
		})

		return user, nil
	}

	var userUUID string

	if err != nil {
		// No cookie: register user; set cookie
		authorizeReason = "new_user"
		user, err := createNewUserAndSetCookie()
		if err != nil {
			s.logger.Error("Cannot create new user and set it in cookie", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userUUID = user.UUID
	} else {
		// Cookie exists: check cookie on valid jwt; set cookie if non-valid
		user, err := s.authService.FetchUserFromToken(jwtCookie.Value)
		if err != nil {
			s.logger.Info("Cannot fetch user uuid from token, create new token", zap.Error(err))
			user, err = createNewUserAndSetCookie()
			if err != nil {
				s.logger.Error("cannot create new user and set it to cookie for user with invalid cookie", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userUUID = user.UUID
			authorizeReason = "invalid_old_token"
		} else {
			userUUID = user.UUID
			authorizeReason = "user_already_has_token"
		}
	}

	s.logger.Debug("User successfuly authorized, redirected", zap.String("user_uuid", userUUID), zap.String("authorize_reason", authorizeReason))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
