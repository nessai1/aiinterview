package service

import "net/http"

func (s *Service) handleAPIGetInterviewList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
