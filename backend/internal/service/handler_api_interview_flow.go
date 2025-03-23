package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/nessai1/aiinterview/internal/domain"
	"github.com/nessai1/aiinterview/internal/interview"
	"go.uber.org/zap"
	"net/http"
)

type AnswerQuestionRequest struct {
	Answer       string `json:"answer"`
	QuestionUUID string `json:"question_uuid"`
}

// 205 - required for document.location.reload()
// 220 - next section
func (s *Service) handleAPIAnswerQuestion(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextUserKey).(domain.User)
	if !ok {
		s.logger.Error("User come to API without user in context", zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	reader := bytes.Buffer{}
	_, err := reader.ReadFrom(r.Body)

	if err != nil {
		s.logger.Error("Cannot read from request body", zap.String("req_uri", r.RequestURI), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var answerRequest AnswerQuestionRequest
	err = json.Unmarshal(reader.Bytes(), &answerRequest)
	if err != nil {
		s.logger.Error("Cannot unmarshal request body", zap.String("req_uri", r.RequestURI), zap.Error(err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	question, err := s.interviewService.AnswerQuestion(r.Context(), user, answerRequest.QuestionUUID, answerRequest.Answer)

	if err == nil || errors.Is(err, interview.ErrSectionOver) {
		code := 200
		if errors.Is(err, interview.ErrSectionOver) {
			code = 220
		}

		jsoned, err := json.Marshal(&question)
		if err != nil {
			s.logger.Error("Cannot marshal question result", zap.String("req_uri", r.RequestURI), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(code)
		_, err = w.Write(jsoned)
		if err != nil {
			s.logger.Error("Cannot write result to user", zap.String("req_uri", r.RequestURI), zap.Error(err))
			return
		}
	}

	if errors.Is(err, interview.ErrAlreadyAnswered) || errors.Is(err, interview.ErrInterviewOver) {
		w.WriteHeader(205)
		return
	}

	s.logger.Error("Error while answer question", zap.Error(err), zap.String("req_uri", r.RequestURI))
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Service) handleAPIGetInterview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	interviewID := vars["interviewID"]
	user, ok := r.Context().Value(contextUserKey).(domain.User)
	if !ok {
		s.logger.Error("User come to API without user in context", zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	s.logger.Debug("Load interview", zap.String("interview_uuid", interviewID), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))

	i, err := s.interviewService.GetInterview(r.Context(), user, interviewID)
	if err != nil {
		s.logger.Error("Error while load interview", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
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
}
