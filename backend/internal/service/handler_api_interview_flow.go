package service

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nessai1/aiinterview/internal/domain"
	"go.uber.org/zap"
	"net/http"
)

type AnswerQuestionRequest struct {
	Answer       string `json:"answer"`
	QuestionUUID string `json:"question_uuid"`
}

type AnswerQuestionResponse struct {
}

func (s *Service) handleAPIAnswerQuestion(w http.ResponseWriter, r *http.Request) {

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

	//s.interviewService.AnswerQuestionAndGetNewQuestion(r.Context(), answerRequest.QuestionUUID, answerRequest.Answer)
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

	interview, err := s.interviewService.GetInterview(r.Context(), user, interviewID)
	if err != nil {
		s.logger.Error("Error while load interview", zap.Error(err), zap.String("user_uuid", user.UUID), zap.String("req_uri", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	jsoned, err := json.Marshal(&interview)
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
